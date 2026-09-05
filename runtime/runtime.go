package runtime

import (
	"github.com/go-native/go-native/ui"
	"sync"
	"sync/atomic"
	"time"
)

// App builds a declarative tree from current state.
type App func() ui.Component

// ContextApp is the production application contract. It receives immutable
// environment and mounted-path context for every committed build.
type ContextApp func(ui.BuildContext) ui.Component

// Runtime owns tree identity, handlers, reconciliation, and render scheduling.
type Runtime struct {
	app                   App
	renderer              Renderer
	layoutProvider        LayoutProvider
	events                *EventRegistry
	mu                    sync.Mutex
	tree                  *ui.Node
	geometry              map[ui.NodeID]ui.LayoutRect
	scheduled             bool
	stopped               bool
	sequence              atomic.Uint64
	timingMu              sync.Mutex
	pending               map[uint64]pendingTiming
	samples               []TimingSample
	eventAt               time.Time
	diagnostics           *Diagnostics
	environment           ui.Environment
	contextApp            ContextApp
	hooks                 *ui.HookRegistry
	focus                 *ui.FocusManager
	lifecycleObservers    map[uint64]func(ui.LifecycleState)
	nextLifecycleObserver uint64
}

type pendingTiming struct {
	sent          time.Time
	event         time.Time
	mutationCount int
}

const timingHistoryLimit uint64 = 1024

// New creates an application runtime.
func New(app App, renderer Renderer) *Runtime {
	runtime := &Runtime{app: app, renderer: renderer, events: NewEventRegistry(), pending: make(map[uint64]pendingTiming), diagnostics: NewDiagnostics(), environment: ui.DefaultEnvironment(), lifecycleObservers: make(map[uint64]func(ui.LifecycleState))}
	runtime.focus = ui.NewFocusManager(runtime)
	runtime.environment.Focus = runtime.focus
	runtime.hooks = ui.NewHookRegistry(runtime)
	return runtime
}

// NewContext creates a runtime using the context-aware application contract.
func NewContext(app ContextApp, renderer Renderer, environment ui.Environment) *Runtime {
	runtime := &Runtime{contextApp: app, renderer: renderer, events: NewEventRegistry(), pending: make(map[uint64]pendingTiming), diagnostics: NewDiagnostics(), environment: environment, lifecycleObservers: make(map[uint64]func(ui.LifecycleState))}
	if environment.Focus == nil {
		environment.Focus = ui.NewFocusManager(runtime)
	} else {
		environment.Focus.SetScheduler(runtime)
	}
	runtime.environment = environment
	runtime.focus = environment.Focus
	runtime.hooks = ui.NewHookRegistry(runtime)
	return runtime
}

// Environment returns the current immutable application environment snapshot.
func (r *Runtime) Environment() ui.Environment {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.environment
}

// UpdateEnvironment replaces runtime context and schedules a render.
func (r *Runtime) UpdateEnvironment(update func(ui.Environment) ui.Environment) {
	if update == nil {
		return
	}
	r.mu.Lock()
	r.environment = update(r.environment)
	stopped := r.stopped
	r.mu.Unlock()
	if !stopped {
		r.Schedule()
	}
}

// SetLifecycle updates the portable lifecycle state and schedules observers.
func (r *Runtime) SetLifecycle(state ui.LifecycleState) {
	r.mu.Lock()
	if r.environment.Lifecycle == state || r.stopped {
		r.mu.Unlock()
		return
	}
	r.environment.Lifecycle = state
	observers := make([]func(ui.LifecycleState), 0, len(r.lifecycleObservers))
	for _, observer := range r.lifecycleObservers {
		observers = append(observers, observer)
	}
	r.mu.Unlock()
	for _, observer := range observers {
		observer(state)
	}
	r.Schedule()
}

// ObserveLifecycle subscribes application services to native lifecycle changes.
// Observation is independent of rendering and cancellation is idempotent.
func (r *Runtime) ObserveLifecycle(observer func(ui.LifecycleState)) func() {
	if observer == nil {
		return func() {}
	}
	r.mu.Lock()
	r.nextLifecycleObserver++
	id := r.nextLifecycleObserver
	r.lifecycleObservers[id] = observer
	state := r.environment.Lifecycle
	r.mu.Unlock()
	observer(state)
	var once sync.Once
	return func() { once.Do(func() { r.mu.Lock(); delete(r.lifecycleObservers, id); r.mu.Unlock() }) }
}

// FocusManager exposes the application-scoped portable focus tree.
func (r *Runtime) FocusManager() *ui.FocusManager { return r.focus }

// SetLayoutProvider installs the Go-owned geometry pipeline used before each
// commit. Passing nil disables computed-frame mutations.
func (r *Runtime) SetLayoutProvider(provider LayoutProvider) {
	r.mu.Lock()
	r.layoutProvider = provider
	r.geometry = nil
	r.mu.Unlock()
}

// Start renders the initial tree and installs this runtime as the state scheduler.
func (r *Runtime) Start() error {
	r.mu.Lock()
	r.stopped = false
	r.mu.Unlock()
	ui.SetScheduler(r)
	r.diagnostics.Record(LogEntry{Kind: LogRuntimeStarted})
	return r.render()
}

// Stop releases handlers owned by the current tree and prevents new renders.
func (r *Runtime) Stop() {
	r.mu.Lock()
	if r.stopped {
		r.mu.Unlock()
		return
	}
	r.stopped = true
	r.scheduled = false
	var lifecycleObservers []func(ui.LifecycleState)
	if r.environment.Lifecycle != ui.LifecycleDestroyed {
		r.environment.Lifecycle = ui.LifecycleDestroyed
		lifecycleObservers = make([]func(ui.LifecycleState), 0, len(r.lifecycleObservers))
		for _, observer := range r.lifecycleObservers {
			lifecycleObservers = append(lifecycleObservers, observer)
		}
	}
	releaseTree(r.events, r.tree)
	r.tree = nil
	r.lifecycleObservers = make(map[uint64]func(ui.LifecycleState))
	r.mu.Unlock()
	for _, observer := range lifecycleObservers {
		observer(ui.LifecycleDestroyed)
	}
	r.hooks.Dispose()
	ui.SetScheduler(nil)
	r.diagnostics.Record(LogEntry{Kind: LogRuntimeStopped})
}

// Schedule coalesces synchronous/re-entrant updates. Rendering is serialized.
func (r *Runtime) Schedule() {
	r.mu.Lock()
	if r.scheduled || r.stopped {
		r.mu.Unlock()
		return
	}
	r.scheduled = true
	r.mu.Unlock()
	// Native event entrypoints invoke handlers serially. A goroutine prevents Set from
	// recursively rendering while a handler still owns application state.
	go func() { r.mu.Lock(); r.scheduled = false; r.mu.Unlock(); _ = r.render() }()
}

// Dispatch invokes a native event callback by ID.
func (r *Runtime) Dispatch(id ui.HandlerID) bool {
	r.timingMu.Lock()
	r.eventAt = time.Now()
	r.timingMu.Unlock()
	if r.events.Dispatch(id) {
		r.diagnostics.Record(LogEntry{Kind: LogEventDispatched, HandlerID: id})
		return true
	}
	r.timingMu.Lock()
	r.eventAt = time.Time{}
	r.timingMu.Unlock()
	r.diagnostics.Record(LogEntry{Kind: LogEventMissing, HandlerID: id})
	return false
}

// DispatchValue invokes a native value event callback by ID.
func (r *Runtime) DispatchValue(id ui.HandlerID, value string) bool {
	r.timingMu.Lock()
	r.eventAt = time.Now()
	r.timingMu.Unlock()
	if r.events.DispatchValue(id, value) {
		r.diagnostics.Record(LogEntry{Kind: LogEventDispatched, HandlerID: id})
		return true
	}
	r.timingMu.Lock()
	r.eventAt = time.Time{}
	r.timingMu.Unlock()
	r.diagnostics.Record(LogEntry{Kind: LogEventMissing, HandlerID: id})
	return false
}

// DispatchBool invokes a native boolean event callback by ID.
func (r *Runtime) DispatchBool(id ui.HandlerID, value bool) bool {
	r.timingMu.Lock()
	r.eventAt = time.Now()
	r.timingMu.Unlock()
	if r.events.DispatchBool(id, value) {
		return true
	}
	r.timingMu.Lock()
	r.eventAt = time.Time{}
	r.timingMu.Unlock()
	return false
}

// DispatchGesture invokes a native gesture callback by ID.
func (r *Runtime) DispatchGesture(id ui.HandlerID, event ui.GestureEvent) bool {
	r.timingMu.Lock()
	r.eventAt = time.Now()
	r.timingMu.Unlock()
	if r.events.DispatchGesture(id, event) {
		return true
	}
	r.timingMu.Lock()
	r.eventAt = time.Time{}
	r.timingMu.Unlock()
	return false
}

// DispatchFocus mirrors a native focus transition into the portable focus tree.
// Native callers identify controls only by integer NodeID values.
func (r *Runtime) DispatchFocus(id ui.NodeID, focused bool) bool {
	r.mu.Lock()
	node := findNode(r.tree, id)
	r.mu.Unlock()
	if node == nil || node.Focus == nil {
		return false
	}
	if focused {
		return r.focus.RequestFocus(node.Focus)
	}
	return r.focus.ClearFocus(node.Focus)
}

// FocusedNodeID returns the currently requested native focus target.
func (r *Runtime) FocusedNodeID() (ui.NodeID, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	focused := r.focus.Focused()
	if focused == nil {
		return 0, false
	}
	return findFocusNodeID(r.tree, focused)
}

func findNode(node *ui.Node, id ui.NodeID) *ui.Node {
	if node == nil {
		return nil
	}
	if node.ID == id {
		return node
	}
	for _, child := range node.Children {
		if found := findNode(child, id); found != nil {
			return found
		}
	}
	return nil
}

func findFocusNodeID(node *ui.Node, focus *ui.FocusNode) (ui.NodeID, bool) {
	if node == nil {
		return 0, false
	}
	if node.Focus == focus {
		return node.ID, true
	}
	for _, child := range node.Children {
		if id, ok := findFocusNodeID(child, focus); ok {
			return id, true
		}
	}
	return 0, false
}

// LogEntries returns a copy of the bounded structured runtime log.
func (r *Runtime) LogEntries() []LogEntry { return r.diagnostics.Entries() }

// ClearLogEntries clears accumulated diagnostic events.
func (r *Runtime) ClearLogEntries() { r.diagnostics.Clear() }

// TreeSnapshot returns a detached snapshot of the last rendered tree.
func (r *Runtime) TreeSnapshot() TreeSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	return SnapshotTree(r.tree)
}

// RecordNativeApply records completion timing reported by a platform renderer.
func (r *Runtime) RecordNativeApply(sequence uint64, nativeApply time.Duration) {
	now := time.Now()
	r.timingMu.Lock()
	pending, ok := r.pending[sequence]
	if ok {
		delete(r.pending, sequence)
		r.samples = append(r.samples, TimingSample{Sequence: sequence, MutationCount: pending.mutationCount, NativeApply: nativeApply, BridgeToApply: now.Sub(pending.sent), EventToApply: elapsedIfSet(now, pending.event)})
		if len(r.samples) > int(timingHistoryLimit) {
			copy(r.samples, r.samples[len(r.samples)-int(timingHistoryLimit):])
			r.samples = r.samples[:timingHistoryLimit]
		}
	}
	r.timingMu.Unlock()
}

// TimingSamples returns a copy of completed native timing samples.
func (r *Runtime) TimingSamples() []TimingSample {
	r.timingMu.Lock()
	defer r.timingMu.Unlock()
	return append([]TimingSample(nil), r.samples...)
}

func elapsedIfSet(now, start time.Time) time.Duration {
	if start.IsZero() {
		return 0
	}
	return now.Sub(start)
}

func (r *Runtime) render() error {
	r.mu.Lock()
	if r.stopped {
		r.diagnostics.Record(LogEntry{Kind: LogRenderSkipped, Message: "runtime stopped"})
		r.mu.Unlock()
		return nil
	}
	r.hooks.BeginRender()
	var component ui.Component
	context := ui.NewBuildContext(r.environment).WithHooks(r.hooks)
	if r.contextApp != nil {
		component = r.contextApp(context)
	} else if r.app != nil {
		component = r.app()
	}
	if component == nil {
		r.hooks.AbortRender()
		r.diagnostics.Record(LogEntry{Kind: LogRenderSkipped, Message: "nil component"})
		r.mu.Unlock()
		return nil
	}
	next := ui.BuildWithContext(component, context)
	r.bindHandlers(r.tree, next)
	stabilizeIDs(r.tree, next)
	syncPortableFocus(next, r.focus.Focused())
	batch := Reconcile(r.tree, next)
	provider := r.layoutProvider
	if provider == nil {
		provider, _ = r.renderer.(LayoutProvider)
	}
	if provider != nil {
		frames, err := provider.ComputeLayout(next, r.environment)
		if err != nil {
			r.hooks.AbortRender()
			r.mu.Unlock()
			return err
		}
		batch.Mutations = attachGeometry(batch.Mutations, next, r.geometry, frames)
		r.geometry = geometryMap(frames)
	}
	if len(batch.Mutations) > 0 {
		batch.Sequence = r.sequence.Add(1)
		r.timingMu.Lock()
		if batch.Sequence > timingHistoryLimit {
			delete(r.pending, batch.Sequence-timingHistoryLimit)
		}
		r.pending[batch.Sequence] = pendingTiming{sent: time.Now(), event: r.eventAt, mutationCount: len(batch.Mutations)}
		r.eventAt = time.Time{}
		r.timingMu.Unlock()
		if err := r.renderer.Apply(batch); err != nil {
			r.timingMu.Lock()
			delete(r.pending, batch.Sequence)
			r.timingMu.Unlock()
			r.diagnostics.Record(LogEntry{Kind: LogRenderFailed, Sequence: batch.Sequence, MutationCount: len(batch.Mutations), Message: err.Error()})
			r.hooks.AbortRender()
			r.mu.Unlock()
			return err
		}
		r.diagnostics.Record(LogEntry{Kind: LogBatchApplied, Sequence: batch.Sequence, MutationCount: len(batch.Mutations)})
	}
	r.releaseRemovedHandlers(r.tree, next)
	r.tree = next
	r.mu.Unlock()
	r.hooks.CommitRender()
	return nil
}

func syncPortableFocus(node *ui.Node, focused *ui.FocusNode) {
	if node == nil {
		return
	}
	if node.Focus != nil {
		node.Props.Focused = node.Focus == focused
	}
	for _, child := range node.Children {
		syncPortableFocus(child, focused)
	}
}

func geometryMap(frames []ui.LayoutFrame) map[ui.NodeID]ui.LayoutRect {
	out := make(map[ui.NodeID]ui.LayoutRect, len(frames))
	for _, frame := range frames {
		out[frame.NodeID] = frame.Rect
	}
	return out
}

func attachGeometry(mutations []Mutation, tree *ui.Node, previous map[ui.NodeID]ui.LayoutRect, frames []ui.LayoutFrame) []Mutation {
	nodes := make(map[ui.NodeID]*ui.Node)
	var visit func(*ui.Node)
	visit = func(node *ui.Node) {
		if node == nil {
			return
		}
		nodes[node.ID] = node
		for _, child := range node.Children {
			visit(child)
		}
	}
	visit(tree)
	mutationByID := make(map[ui.NodeID]int)
	for index := range mutations {
		if mutations[index].Type == MutationCreate || mutations[index].Type == MutationUpdate {
			mutationByID[mutations[index].NodeID] = index
		}
	}
	for _, frame := range frames {
		if index, ok := mutationByID[frame.NodeID]; ok {
			mutations[index].HasFrame, mutations[index].Frame = true, frame.Rect
			continue
		}
		if old, ok := previous[frame.NodeID]; ok && old == frame.Rect {
			continue
		}
		node := nodes[frame.NodeID]
		if node == nil {
			continue
		}
		mutations = append(mutations, Mutation{Type: MutationUpdate, NodeID: node.ID, NodeType: node.Type, Props: node.Props, Style: node.Style, Platform: node.Platform, HasFrame: true, Frame: frame.Rect})
	}
	return mutations
}

func stabilizeIDs(oldNode, newNode *ui.Node) {
	if oldNode == nil || newNode == nil || oldNode.Type != newNode.Type || newNode.ExplicitID {
		return
	}
	newNode.ID = oldNode.ID
	oldExplicit := make(map[ui.NodeID]*ui.Node)
	for _, child := range oldNode.Children {
		if child.ExplicitID {
			oldExplicit[child.ID] = child
		}
	}
	for i, child := range newNode.Children {
		if child.ExplicitID {
			if oldChild := oldExplicit[child.ID]; oldChild != nil {
				stabilizeDescendants(oldChild, child)
			}
			continue
		}
		if i < len(oldNode.Children) && !oldNode.Children[i].ExplicitID {
			stabilizeIDs(oldNode.Children[i], child)
		}
	}
}

func stabilizeDescendants(oldNode, newNode *ui.Node) {
	if oldNode == nil || newNode == nil || oldNode.ID != newNode.ID || oldNode.Type != newNode.Type {
		return
	}
	for i, child := range newNode.Children {
		if i < len(oldNode.Children) {
			stabilizeIDs(oldNode.Children[i], child)
		}
	}
}

func (r *Runtime) bindHandlers(oldNode, n *ui.Node) {
	oldGestureIDs := []ui.HandlerID(nil)
	if oldNode != nil && oldNode.Type == n.Type {
		oldGestureIDs = oldNode.GestureHandlerIDs
	}
	n.GestureHandlerIDs = make([]ui.HandlerID, len(n.Intents.Gestures))
	for i, intent := range n.Intents.Gestures {
		if i < len(oldGestureIDs) && oldGestureIDs[i] != 0 {
			n.GestureHandlerIDs[i] = oldGestureIDs[i]
			r.events.ReplaceGesture(oldGestureIDs[i], intent.Handler)
		} else {
			n.GestureHandlerIDs[i] = r.events.RegisterGesture(intent.Handler)
		}
	}
	for i := len(n.Intents.Gestures); i < len(oldGestureIDs); i++ {
		r.events.Release(oldGestureIDs[i])
	}
	n.Props.Interactions = marshalInteractions(n.Intents, n.GestureHandlerIDs)
	if n.Type == ui.NodeButton {
		if fn := n.Press; fn != nil {
			if oldNode != nil && oldNode.Type == n.Type && oldNode.Props.OnPress != 0 {
				n.Props.OnPress = oldNode.Props.OnPress
				r.events.Replace(n.Props.OnPress, fn)
			} else {
				n.Props.OnPress = r.events.Register(fn)
			}
		}
	}
	if n.Type == ui.NodeTextInput {
		if fn := n.Change; fn != nil {
			if oldNode != nil && oldNode.Type == n.Type && oldNode.Props.OnChange != 0 {
				n.Props.OnChange = oldNode.Props.OnChange
				r.events.ReplaceValue(n.Props.OnChange, fn)
			} else {
				n.Props.OnChange = r.events.RegisterValue(fn)
			}
		}
	}
	if n.Type == ui.NodeSwitch && n.Toggle != nil {
		if oldNode != nil && oldNode.Type == n.Type && oldNode.Props.OnToggle != 0 {
			n.Props.OnToggle = oldNode.Props.OnToggle
			r.events.ReplaceBool(n.Props.OnToggle, n.Toggle)
		} else {
			n.Props.OnToggle = r.events.RegisterBool(n.Toggle)
		}
	}
	oldExplicit := make(map[ui.NodeID]*ui.Node)
	if oldNode != nil {
		for _, child := range oldNode.Children {
			if child.ExplicitID {
				oldExplicit[child.ID] = child
			}
		}
	}
	for i, child := range n.Children {
		var oldChild *ui.Node
		if child.ExplicitID {
			oldChild = oldExplicit[child.ID]
		} else if oldNode != nil && i < len(oldNode.Children) && !oldNode.Children[i].ExplicitID {
			oldChild = oldNode.Children[i]
		}
		r.bindHandlers(oldChild, child)
	}
}

func (r *Runtime) releaseRemovedHandlers(oldNode, newNode *ui.Node) {
	if oldNode == nil {
		return
	}
	if newNode == nil || oldNode.ID != newNode.ID {
		releaseTree(r.events, oldNode)
		return
	}
	newByID := map[ui.NodeID]*ui.Node{}
	for _, n := range newNode.Children {
		newByID[n.ID] = n
	}
	for _, o := range oldNode.Children {
		r.releaseRemovedHandlers(o, newByID[o.ID])
	}
}
func releaseTree(events *EventRegistry, n *ui.Node) {
	if n == nil {
		return
	}
	events.Release(n.Props.OnPress)
	events.Release(n.Props.OnChange)
	events.Release(n.Props.OnToggle)
	for _, id := range n.GestureHandlerIDs {
		events.Release(id)
	}
	for _, c := range n.Children {
		releaseTree(events, c)
	}
}

package runtime

import (
	"github.com/go-native/go-native/ui"
	"sync"
)

// App builds a declarative tree from current state.
type App func() ui.Component

// Runtime owns tree identity, handlers, reconciliation, and render scheduling.
type Runtime struct {
	app       App
	renderer  Renderer
	events    *EventRegistry
	mu        sync.Mutex
	tree      *ui.Node
	scheduled bool
}

// New creates an application runtime.
func New(app App, renderer Renderer) *Runtime {
	return &Runtime{app: app, renderer: renderer, events: NewEventRegistry()}
}

// Start renders the initial tree and installs this runtime as the state scheduler.
func (r *Runtime) Start() error { ui.SetScheduler(r); return r.render() }

// Schedule coalesces synchronous/re-entrant updates. Rendering is serialized.
func (r *Runtime) Schedule() {
	r.mu.Lock()
	if r.scheduled {
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
func (r *Runtime) Dispatch(id ui.HandlerID) bool { return r.events.Dispatch(id) }

func (r *Runtime) render() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	component := r.app()
	if component == nil {
		return nil
	}
	next := component.Build()
	r.bindHandlers(r.tree, next)
	stabilizeIDs(r.tree, next)
	batch := Reconcile(r.tree, next)
	if len(batch.Mutations) > 0 {
		if err := r.renderer.Apply(batch); err != nil {
			return err
		}
	}
	r.releaseRemovedHandlers(r.tree, next)
	r.tree = next
	return nil
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
	events.Release(n.Props.OnPress)
	for _, c := range n.Children {
		releaseTree(events, c)
	}
}

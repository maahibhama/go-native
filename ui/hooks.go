package ui

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// Cleanup releases resources created by an effect.
type Cleanup func()

// Effect runs after a successful render commit. Its context is cancelled before
// cleanup when dependencies change or the mounted component disappears.
type Effect func(context.Context) Cleanup

type hookKind uint8

const (
	hookState hookKind = iota + 1
	hookReducer
	hookRef
	hookMemo
	hookEffect
	hookLayoutEffect
)

type hookSlot struct {
	kind    hookKind
	value   any
	deps    []any
	ready   bool
	effect  Effect
	pending bool
	cancel  context.CancelFunc
	cleanup Cleanup
}

type hookScope struct {
	path   string
	cursor int
	slots  []hookSlot
}

// HookRegistry owns persistent hook state for one renderer/runtime.
type HookRegistry struct {
	mu          sync.Mutex
	transaction sync.Mutex
	scheduler   Scheduler
	scopes      map[string]*hookScope
	active      map[string]bool
	activeOrder []string
	rendering   bool
}

func NewHookRegistry(scheduler Scheduler) *HookRegistry {
	return &HookRegistry{scheduler: scheduler, scopes: make(map[string]*hookScope), active: make(map[string]bool)}
}

// BeginRender starts a hook transaction. CommitRender must follow after a
// successful host commit; AbortRender discards pending effects.
func (r *HookRegistry) BeginRender() {
	if r == nil {
		return
	}
	r.transaction.Lock()
	r.mu.Lock()
	r.active = make(map[string]bool)
	r.activeOrder = nil
	r.rendering = true
	r.mu.Unlock()
}

func (r *HookRegistry) beginScope(path string) *hookScope {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.rendering {
		panic("ui: hooks can only run during a render")
	}
	scope := r.scopes[path]
	if scope == nil {
		scope = &hookScope{path: path}
		r.scopes[path] = scope
	}
	scope.cursor = 0
	if !r.active[path] {
		r.activeOrder = append(r.activeOrder, path)
	}
	r.active[path] = true
	return scope
}

func (r *HookRegistry) finishScope(scope *hookScope) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if scope.cursor < len(scope.slots) {
		panic(fmt.Sprintf("ui: hook count changed at %s: used %d, previously %d", scope.path, scope.cursor, len(scope.slots)))
	}
}

func (r *HookRegistry) next(scope *hookScope, kind hookKind) *hookSlot {
	r.mu.Lock()
	defer r.mu.Unlock()
	index := scope.cursor
	scope.cursor++
	if index == len(scope.slots) {
		scope.slots = append(scope.slots, hookSlot{kind: kind})
	}
	slot := &scope.slots[index]
	if slot.kind != kind {
		panic(fmt.Sprintf("ui: hook order changed at %s index %d", scope.path, index))
	}
	return slot
}

// CommitRender applies effect changes and unmounts scopes absent from this render.
func (r *HookRegistry) CommitRender() {
	if r == nil {
		return
	}
	defer r.transaction.Unlock()
	r.mu.Lock()
	var removed []*hookScope
	for path, scope := range r.scopes {
		if !r.active[path] {
			removed = append(removed, scope)
			delete(r.scopes, path)
		}
	}
	active := make([]*hookScope, 0, len(r.activeOrder))
	for _, path := range r.activeOrder {
		active = append(active, r.scopes[path])
	}
	r.rendering = false
	r.mu.Unlock()

	for _, scope := range removed {
		cleanupScope(scope)
	}
	for _, kind := range []hookKind{hookLayoutEffect, hookEffect} {
		for _, scope := range active {
			for index := range scope.slots {
				slot := &scope.slots[index]
				if slot.kind == kind && slot.pending {
					runEffect(slot)
				}
			}
		}
	}
}

// AbortRender prevents effects prepared by a failed render from running.
func (r *HookRegistry) AbortRender() {
	if r == nil {
		return
	}
	defer r.transaction.Unlock()
	r.mu.Lock()
	for path := range r.active {
		for index := range r.scopes[path].slots {
			slot := &r.scopes[path].slots[index]
			if slot.pending {
				slot.ready = false
				slot.pending = false
			}
		}
	}
	r.rendering = false
	r.mu.Unlock()
}

// Dispose unmounts every hook scope and runs all effect cleanup functions.
func (r *HookRegistry) Dispose() {
	if r == nil {
		return
	}
	r.transaction.Lock()
	defer r.transaction.Unlock()
	r.mu.Lock()
	scopes := make([]*hookScope, 0, len(r.scopes))
	for _, scope := range r.scopes {
		scopes = append(scopes, scope)
	}
	r.scopes = make(map[string]*hookScope)
	r.active = make(map[string]bool)
	r.rendering = false
	r.mu.Unlock()
	for _, scope := range scopes {
		cleanupScope(scope)
	}
}

func cleanupScope(scope *hookScope) {
	for index := range scope.slots {
		cleanupEffect(&scope.slots[index])
	}
}

func cleanupEffect(slot *hookSlot) {
	if slot.cancel != nil {
		slot.cancel()
		slot.cancel = nil
	}
	if slot.cleanup != nil {
		slot.cleanup()
		slot.cleanup = nil
	}
}

func runEffect(slot *hookSlot) {
	cleanupEffect(slot)
	slot.pending = false
	if slot.effect == nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	slot.cancel = cancel
	slot.cleanup = slot.effect(ctx)
}

type functionalComponent struct {
	key    string
	render func(BuildContext) Component
}

// Functional creates a context-aware component with persistent hook identity.
// key must be stable among siblings across renders.
func Functional(key string, render func(BuildContext) Component) Component {
	return &functionalComponent{key: key, render: render}
}

func (c *functionalComponent) Build() *Node {
	return c.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (c *functionalComponent) BuildContext(context BuildContext) *Node {
	if c == nil || c.render == nil {
		return nil
	}
	context = context.Child("component:" + c.key)
	if context.hooks == nil {
		return BuildWithContext(c.render(context), context)
	}
	scope := context.hooks.beginScope(context.Path)
	component := c.render(context)
	context.hooks.finishScope(scope)
	return BuildWithContext(component, context)
}

func hookContext(context BuildContext) (*HookRegistry, *hookScope) {
	if context.hooks == nil {
		panic("ui: hook used outside a mounted Functional component")
	}
	context.hooks.mu.Lock()
	scope := context.hooks.scopes[context.Path]
	context.hooks.mu.Unlock()
	if scope == nil {
		panic("ui: hook used outside a mounted Functional component")
	}
	return context.hooks, scope
}

func UseState[T any](context BuildContext, initial T) *State[T] {
	registry, scope := hookContext(context)
	slot := registry.next(scope, hookState)
	if !slot.ready {
		slot.value = newScheduledState(initial, registry.scheduler)
		slot.ready = true
	}
	state, ok := slot.value.(*State[T])
	if !ok {
		panic("ui: UseState type changed at " + context.Path)
	}
	return state
}

type ReducerState[S, A any] struct {
	state   *State[S]
	reducer func(S, A) S
}

func (r *ReducerState[S, A]) Get() S { return r.state.Get() }
func (r *ReducerState[S, A]) Dispatch(action A) {
	r.state.Update(func(value S) S { return r.reducer(value, action) })
}

func UseReducer[S, A any](context BuildContext, initial S, reducer func(S, A) S) *ReducerState[S, A] {
	registry, scope := hookContext(context)
	slot := registry.next(scope, hookReducer)
	if !slot.ready {
		slot.value = &ReducerState[S, A]{state: newScheduledState(initial, registry.scheduler), reducer: reducer}
		slot.ready = true
	}
	value, ok := slot.value.(*ReducerState[S, A])
	if !ok {
		panic("ui: UseReducer type changed at " + context.Path)
	}
	value.reducer = reducer
	return value
}

type Ref[T any] struct{ Current T }

func UseRef[T any](context BuildContext, initial T) *Ref[T] {
	registry, scope := hookContext(context)
	slot := registry.next(scope, hookRef)
	if !slot.ready {
		slot.value = &Ref[T]{Current: initial}
		slot.ready = true
	}
	value, ok := slot.value.(*Ref[T])
	if !ok {
		panic("ui: UseRef type changed at " + context.Path)
	}
	return value
}

func UseMemo[T any](context BuildContext, factory func() T, dependencies ...any) T {
	registry, scope := hookContext(context)
	slot := registry.next(scope, hookMemo)
	if !slot.ready || !reflect.DeepEqual(slot.deps, dependencies) {
		slot.value = factory()
		slot.deps = append([]any(nil), dependencies...)
		slot.ready = true
	}
	value, ok := slot.value.(T)
	if !ok {
		panic("ui: UseMemo type changed at " + context.Path)
	}
	return value
}

func UseCallback[T any](context BuildContext, callback T, dependencies ...any) T {
	return UseMemo(context, func() T { return callback }, dependencies...)
}

func UseEffect(context BuildContext, effect Effect, dependencies ...any) {
	useEffect(context, hookEffect, effect, dependencies)
}

func UseLayoutEffect(context BuildContext, effect Effect, dependencies ...any) {
	useEffect(context, hookLayoutEffect, effect, dependencies)
}

func useEffect(context BuildContext, kind hookKind, effect Effect, dependencies []any) {
	registry, scope := hookContext(context)
	slot := registry.next(scope, kind)
	if !slot.ready || !reflect.DeepEqual(slot.deps, dependencies) {
		slot.deps = append([]any(nil), dependencies...)
		slot.effect = effect
		slot.pending = true
		slot.ready = true
	}
}

func UseLifecycle(context BuildContext) LifecycleState { return context.Environment.Lifecycle }
func UseMediaQuery(context BuildContext) MediaQuery    { return context.Environment.MediaQuery }

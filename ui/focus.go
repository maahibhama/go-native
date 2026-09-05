package ui

import "sync"

// FocusOptions controls whether a focus node participates in keyboard traversal.
type FocusOptions struct {
	CanRequestFocus bool
	SkipTraversal   bool
}

func DefaultFocusOptions() FocusOptions { return FocusOptions{CanRequestFocus: true} }

// FocusNode is a stable, portable focus identity owned by a mounted component.
// Native renderers may mirror its state, but it never contains a native object.
type FocusNode struct {
	mu        sync.RWMutex
	manager   *FocusManager
	scope     *FocusScope
	key       string
	options   FocusOptions
	listeners map[uint64]func(bool)
	next      uint64
}

// NewFocusNode creates a focus identity. Use a stable key for diagnostics.
func NewFocusNode(key string, options FocusOptions) *FocusNode {
	return &FocusNode{key: key, options: options, listeners: make(map[uint64]func(bool))}
}

func (n *FocusNode) Key() string { n.mu.RLock(); defer n.mu.RUnlock(); return n.key }

func (n *FocusNode) HasFocus() bool {
	n.mu.RLock()
	manager := n.manager
	n.mu.RUnlock()
	return manager != nil && manager.Focused() == n
}

func (n *FocusNode) RequestFocus() bool {
	n.mu.RLock()
	manager := n.manager
	n.mu.RUnlock()
	return manager != nil && manager.RequestFocus(n)
}

func (n *FocusNode) Unfocus() {
	n.mu.RLock()
	manager := n.manager
	n.mu.RUnlock()
	if manager != nil {
		manager.ClearFocus(n)
	}
}

// Observe calls listener when this node gains or loses focus. The returned
// function removes the listener and is safe to call repeatedly.
func (n *FocusNode) Observe(listener func(bool)) func() {
	if n == nil || listener == nil {
		return func() {}
	}
	n.mu.Lock()
	n.next++
	id := n.next
	n.listeners[id] = listener
	n.mu.Unlock()
	var once sync.Once
	return func() { once.Do(func() { n.mu.Lock(); delete(n.listeners, id); n.mu.Unlock() }) }
}

func (n *FocusNode) notify(focused bool) {
	n.mu.RLock()
	listeners := make([]func(bool), 0, len(n.listeners))
	for _, listener := range n.listeners {
		listeners = append(listeners, listener)
	}
	n.mu.RUnlock()
	for _, listener := range listeners {
		listener(focused)
	}
}

// FocusScope groups focus nodes and child scopes for deterministic traversal.
type FocusScope struct {
	manager  *FocusManager
	parent   *FocusScope
	children []*FocusScope
	nodes    []*FocusNode
}

// FocusManager owns one application focus tree.
type FocusManager struct {
	mu        sync.RWMutex
	root      *FocusScope
	focused   *FocusNode
	scheduler Scheduler
}

func NewFocusManager(scheduler Scheduler) *FocusManager {
	m := &FocusManager{scheduler: scheduler}
	m.root = &FocusScope{manager: m}
	return m
}

// SetScheduler connects focus changes to a render host.
func (m *FocusManager) SetScheduler(scheduler Scheduler) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.scheduler = scheduler
	m.mu.Unlock()
}

func (m *FocusManager) RootScope() *FocusScope {
	if m == nil {
		return nil
	}
	return m.root
}

// NewScope creates a traversal group beneath parent. Passing nil uses the root.
func (m *FocusManager) NewScope(parent *FocusScope) *FocusScope {
	if m == nil {
		return nil
	}
	if parent == nil {
		parent = m.root
	}
	scope := &FocusScope{manager: m, parent: parent}
	m.mu.Lock()
	parent.children = append(parent.children, scope)
	m.mu.Unlock()
	return scope
}

// Dispose removes a scope, all descendant scopes, and their mounted nodes.
func (s *FocusScope) Dispose() {
	if s == nil || s.manager == nil || s.parent == nil {
		return
	}
	m := s.manager
	m.mu.Lock()
	nodes := traversalNodes(s)
	var lost *FocusNode
	for i, child := range s.parent.children {
		if child == s {
			s.parent.children = append(s.parent.children[:i], s.parent.children[i+1:]...)
			break
		}
	}
	for _, node := range nodes {
		if m.focused == node {
			m.focused = nil
			lost = node
		}
		node.mu.Lock()
		node.manager, node.scope = nil, nil
		node.mu.Unlock()
	}
	s.manager, s.parent, s.children, s.nodes = nil, nil, nil, nil
	m.mu.Unlock()
	if lost != nil {
		lost.notify(false)
		m.schedule()
	}
}

// Mount attaches node to a scope. The returned cleanup deterministically
// removes it and clears focus if necessary.
func (m *FocusManager) Mount(scope *FocusScope, node *FocusNode) func() {
	if m == nil || node == nil {
		return func() {}
	}
	if scope == nil {
		scope = m.root
	}
	if scope.manager != m {
		panic("ui: focus scope belongs to a different manager")
	}
	m.mu.Lock()
	for _, existing := range scope.nodes {
		if existing == node {
			m.mu.Unlock()
			return func() { m.Unmount(node) }
		}
	}
	scope.nodes = append(scope.nodes, node)
	node.mu.Lock()
	node.manager, node.scope = m, scope
	node.mu.Unlock()
	m.mu.Unlock()
	var once sync.Once
	return func() { once.Do(func() { m.Unmount(node) }) }
}

func (m *FocusManager) Unmount(node *FocusNode) {
	if m == nil || node == nil {
		return
	}
	var lost bool
	m.mu.Lock()
	node.mu.RLock()
	scope := node.scope
	owner := node.manager
	node.mu.RUnlock()
	if owner == m && scope != nil {
		for i, candidate := range scope.nodes {
			if candidate == node {
				scope.nodes = append(scope.nodes[:i], scope.nodes[i+1:]...)
				break
			}
		}
		if m.focused == node {
			m.focused = nil
			lost = true
		}
		node.mu.Lock()
		node.manager, node.scope = nil, nil
		node.mu.Unlock()
	}
	m.mu.Unlock()
	if lost {
		node.notify(false)
		m.schedule()
	}
}

func (m *FocusManager) Focused() *FocusNode {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.focused
}

func (m *FocusManager) RequestFocus(node *FocusNode) bool {
	if m == nil || node == nil {
		return false
	}
	node.mu.RLock()
	owner, options := node.manager, node.options
	node.mu.RUnlock()
	if owner != m || !options.CanRequestFocus {
		return false
	}
	m.mu.Lock()
	previous := m.focused
	if previous == node {
		m.mu.Unlock()
		return true
	}
	m.focused = node
	m.mu.Unlock()
	if previous != nil {
		previous.notify(false)
	}
	node.notify(true)
	m.schedule()
	return true
}

// ClearFocus clears current focus. If expected is non-nil, it clears only when
// that node currently owns focus.
func (m *FocusManager) ClearFocus(expected *FocusNode) bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	previous := m.focused
	if previous == nil || (expected != nil && previous != expected) {
		m.mu.Unlock()
		return false
	}
	m.focused = nil
	m.mu.Unlock()
	previous.notify(false)
	m.schedule()
	return true
}

func (m *FocusManager) FocusNext(scope *FocusScope) bool     { return m.move(scope, 1) }
func (m *FocusManager) FocusPrevious(scope *FocusScope) bool { return m.move(scope, -1) }

func (m *FocusManager) move(scope *FocusScope, delta int) bool {
	if m == nil {
		return false
	}
	if scope == nil {
		scope = m.root
	}
	m.mu.RLock()
	nodes := traversalNodes(scope)
	current := m.focused
	m.mu.RUnlock()
	eligible := nodes[:0]
	for _, node := range nodes {
		node.mu.RLock()
		o := node.options
		node.mu.RUnlock()
		if o.CanRequestFocus && !o.SkipTraversal {
			eligible = append(eligible, node)
		}
	}
	if len(eligible) == 0 {
		return false
	}
	index := -1
	for i, node := range eligible {
		if node == current {
			index = i
			break
		}
	}
	if delta > 0 {
		index = (index + 1) % len(eligible)
	} else if index <= 0 {
		index = len(eligible) - 1
	} else {
		index--
	}
	return m.RequestFocus(eligible[index])
}

func traversalNodes(scope *FocusScope) []*FocusNode {
	out := append([]*FocusNode(nil), scope.nodes...)
	for _, child := range scope.children {
		out = append(out, traversalNodes(child)...)
	}
	return out
}

func (m *FocusManager) schedule() {
	m.mu.RLock()
	scheduler := m.scheduler
	m.mu.RUnlock()
	if scheduler != nil {
		scheduler.Schedule()
	}
}

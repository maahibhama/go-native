package ui

import "testing"

type focusScheduler struct{ calls int }

func (s *focusScheduler) Schedule() { s.calls++ }

func TestFocusManagerTraversalAndCleanup(t *testing.T) {
	scheduler := &focusScheduler{}
	manager := NewFocusManager(scheduler)
	first := NewFocusNode("first", DefaultFocusOptions())
	skipped := NewFocusNode("skip", FocusOptions{CanRequestFocus: true, SkipTraversal: true})
	last := NewFocusNode("last", DefaultFocusOptions())
	unmountFirst := manager.Mount(nil, first)
	manager.Mount(nil, skipped)
	manager.Mount(nil, last)
	if !manager.FocusNext(nil) || manager.Focused() != first {
		t.Fatal("expected first node")
	}
	if !manager.FocusNext(nil) || manager.Focused() != last {
		t.Fatal("expected traversal to skip node")
	}
	if !manager.FocusPrevious(nil) || manager.Focused() != first {
		t.Fatal("expected previous traversal")
	}
	unmountFirst()
	if manager.Focused() != nil || first.HasFocus() {
		t.Fatal("unmount must clear focus")
	}
	if scheduler.calls != 4 {
		t.Fatalf("scheduler calls = %d", scheduler.calls)
	}
}

func TestFocusScopeDisposalDetachesDescendants(t *testing.T) {
	manager := NewFocusManager(nil)
	scope := manager.NewScope(nil)
	child := manager.NewScope(scope)
	node := NewFocusNode("nested", DefaultFocusOptions())
	manager.Mount(child, node)
	if !node.RequestFocus() {
		t.Fatal("request focus")
	}
	scope.Dispose()
	if node.RequestFocus() || manager.Focused() != nil {
		t.Fatal("disposed nodes must detach")
	}
}

func TestUseFocusNodeUnmountsWithHookScope(t *testing.T) {
	scheduler := &focusScheduler{}
	manager := NewFocusManager(scheduler)
	hooks := NewHookRegistry(scheduler)
	environment := DefaultEnvironment()
	environment.Focus = manager
	ctx := NewBuildContext(environment).WithHooks(hooks)
	var node *FocusNode
	hooks.BeginRender()
	BuildWithContext(Functional("field", func(ctx BuildContext) Component {
		node = UseFocusNode(ctx, "email", DefaultFocusOptions())
		return Text("email")
	}), ctx)
	hooks.CommitRender()
	if !node.RequestFocus() {
		t.Fatal("mounted hook focus node")
	}
	hooks.BeginRender()
	BuildWithContext(Text("removed"), ctx)
	hooks.CommitRender()
	if node.RequestFocus() || manager.Focused() != nil {
		t.Fatal("unmounted hook focus node")
	}
}

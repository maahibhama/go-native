package ui

import (
	"context"
	"reflect"
	"testing"
)

func TestFunctionalHooksPersistByMountedPath(t *testing.T) {
	scheduler := &testScheduler{}
	hooks := NewHookRegistry(scheduler)
	ctx := NewBuildContext(DefaultEnvironment()).WithHooks(hooks)
	var state *State[int]
	var ref *Ref[string]
	memoCalls := 0
	component := Functional("counter", func(ctx BuildContext) Component {
		state = UseState(ctx, 1)
		ref = UseRef(ctx, "mounted")
		value := UseMemo(ctx, func() int { memoCalls++; return state.Get() * 2 }, state.Get())
		return Text(string(rune('0' + value)))
	})

	hooks.BeginRender()
	if got := BuildWithContext(component, ctx).Props.Text; got != "2" {
		t.Fatalf("first text = %q", got)
	}
	hooks.CommitRender()
	firstState, firstRef := state, ref
	state.Set(2)
	if scheduler.calls != 1 {
		t.Fatalf("scheduled = %d", scheduler.calls)
	}

	hooks.BeginRender()
	if got := BuildWithContext(component, ctx).Props.Text; got != "4" {
		t.Fatalf("second text = %q", got)
	}
	hooks.CommitRender()
	if state != firstState || ref != firstRef || memoCalls != 2 {
		t.Fatalf("hooks did not persist: state=%v ref=%v memoCalls=%d", state == firstState, ref == firstRef, memoCalls)
	}
}

func TestEffectsCommitInOrderAndCleanUpOnUnmount(t *testing.T) {
	hooks := NewHookRegistry(&testScheduler{})
	ctx := NewBuildContext(DefaultEnvironment()).WithHooks(hooks)
	dependency := 1
	events := []string{}
	component := Functional("effects", func(ctx BuildContext) Component {
		UseLayoutEffect(ctx, func(context.Context) Cleanup {
			events = append(events, "layout")
			return func() { events = append(events, "layout-cleanup") }
		}, dependency)
		UseEffect(ctx, func(context.Context) Cleanup {
			events = append(events, "effect")
			return func() { events = append(events, "effect-cleanup") }
		}, dependency)
		return Text("effects")
	})

	hooks.BeginRender()
	BuildWithContext(component, ctx)
	hooks.CommitRender()
	if want := []string{"layout", "effect"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
	dependency++
	hooks.BeginRender()
	BuildWithContext(component, ctx)
	hooks.CommitRender()
	if want := []string{"layout", "effect", "layout-cleanup", "layout", "effect-cleanup", "effect"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
	hooks.BeginRender()
	hooks.CommitRender()
	if want := []string{"layout", "effect", "layout-cleanup", "layout", "effect-cleanup", "effect", "layout-cleanup", "effect-cleanup"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("unmount events = %v, want %v", events, want)
	}
}

func TestNestedContextComponentReceivesLazyChildContext(t *testing.T) {
	environment := DefaultEnvironment()
	environment.Locale = "fr"
	root := Column(contextualFixture{})
	if got := BuildWithContext(root, NewBuildContext(environment)).Children[0].Props.Text; got != "fr" {
		t.Fatalf("nested contextual text = %q", got)
	}
}

package runtime

import (
	"github.com/go-native/go-native/ui"
	"sync"
	"testing"
	"time"
)

type recordingRenderer struct {
	mu      sync.Mutex
	batches []MutationBatch
}

type rendererFunc func(MutationBatch) error

func (f rendererFunc) Apply(batch MutationBatch) error { return f(batch) }

func (r *recordingRenderer) Apply(b MutationBatch) error {
	r.mu.Lock()
	r.batches = append(r.batches, b)
	r.mu.Unlock()
	return nil
}

func TestRuntimeStateChangeOnlyUpdatesText(t *testing.T) {
	state := ui.NewState(0)
	renderer := &recordingRenderer{}
	app := func() ui.Component {
		return ui.Column(ui.Text(string(rune('0'+state.Get()))), ui.Button("Increment", func() { state.Set(state.Get() + 1) }))
	}
	rt := New(app, renderer)
	if err := rt.Start(); err != nil {
		t.Fatal(err)
	}
	renderer.mu.Lock()
	handler := renderer.batches[0].Mutations[len(renderer.batches[0].Mutations)-2].Props.OnPress
	renderer.mu.Unlock()
	if handler == 0 {
		t.Fatal("button handler not bound")
	}
	rt.Dispatch(handler)
	deadline := time.Now().Add(time.Second)
	for {
		renderer.mu.Lock()
		n := len(renderer.batches)
		renderer.mu.Unlock()
		if n >= 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("render not scheduled")
		}
		time.Sleep(time.Millisecond)
	}
	renderer.mu.Lock()
	defer renderer.mu.Unlock()
	got := renderer.batches[1].Mutations
	if len(got) != 1 || got[0].Type != MutationUpdate || got[0].NodeType != ui.NodeText || got[0].Props.Text != "1" {
		t.Fatalf("expected one text update, got %#v", got)
	}
}

func TestStabilizeIDsPreservesExplicitIdentityAcrossReorder(t *testing.T) {
	oldTree := ui.Column(
		ui.WithID(ui.Text("a"), 100),
		ui.WithID(ui.Text("b"), 200),
	).Build()
	newTree := ui.Column(
		ui.WithID(ui.Text("b"), 200),
		ui.WithID(ui.Text("a"), 100),
	).Build()
	stabilizeIDs(oldTree, newTree)
	if newTree.ID != oldTree.ID {
		t.Fatal("unkeyed root identity was not stabilized")
	}
	if newTree.Children[0].ID != 200 || newTree.Children[1].ID != 100 {
		t.Fatalf("explicit IDs were replaced: %d %d", newTree.Children[0].ID, newTree.Children[1].ID)
	}
	batch := Reconcile(oldTree, newTree)
	if len(batch.Mutations) != 1 || batch.Mutations[0].Type != MutationMove {
		t.Fatalf("expected one native move, got %#v", batch.Mutations)
	}
}

func TestKeyedButtonReorderPreservesHandlerIdentity(t *testing.T) {
	r := &Runtime{events: NewEventRegistry()}
	firstCalls, secondCalls := 0, 0
	oldTree := ui.Column(
		ui.WithID(ui.Button("first", func() { firstCalls++ }), 100),
		ui.WithID(ui.Button("second", func() { secondCalls++ }), 200),
	).Build()
	r.bindHandlers(nil, oldTree)
	firstHandler := oldTree.Children[0].Props.OnPress
	secondHandler := oldTree.Children[1].Props.OnPress

	nextTree := ui.Column(
		ui.WithID(ui.Button("second", func() { secondCalls += 10 }), 200),
		ui.WithID(ui.Button("first", func() { firstCalls += 10 }), 100),
	).Build()
	r.bindHandlers(oldTree, nextTree)
	if nextTree.Children[0].Props.OnPress != secondHandler || nextTree.Children[1].Props.OnPress != firstHandler {
		t.Fatalf("handlers followed position rather than identity: %#v", nextTree.Children)
	}
	r.events.Dispatch(nextTree.Children[0].Props.OnPress)
	r.events.Dispatch(nextTree.Children[1].Props.OnPress)
	if firstCalls != 10 || secondCalls != 10 {
		t.Fatalf("wrong callbacks dispatched: first=%d second=%d", firstCalls, secondCalls)
	}
}

func TestTextInputValueHandlerRetainsIdentity(t *testing.T) {
	r := New(func() ui.Component { return nil }, &recordingRenderer{})
	got := ""
	oldTree := ui.TextInput("a", func(value string) { got = "old:" + value }).Build()
	r.bindHandlers(nil, oldTree)
	id := oldTree.Props.OnChange
	nextTree := ui.TextInput("b", func(value string) { got = "new:" + value }).Build()
	r.bindHandlers(oldTree, nextTree)
	if id == 0 || nextTree.Props.OnChange != id {
		t.Fatalf("handler identity changed: %d -> %d", id, nextTree.Props.OnChange)
	}
	if !r.DispatchValue(id, "typed") || got != "new:typed" {
		t.Fatalf("value callback got %q", got)
	}
	r.releaseRemovedHandlers(nextTree, nil)
	if r.DispatchValue(id, "ignored") {
		t.Fatal("removed value handler dispatched")
	}
}

func TestSwitchHandlerRetainsIdentity(t *testing.T) {
	r := New(func() ui.Component { return nil }, &recordingRenderer{})
	got := false
	oldTree := ui.Switch(false, func(bool) {}).Build()
	r.bindHandlers(nil, oldTree)
	id := oldTree.Props.OnToggle
	next := ui.Switch(true, func(v bool) { got = v }).Build()
	r.bindHandlers(oldTree, next)
	if id == 0 || next.Props.OnToggle != id {
		t.Fatalf("toggle identity changed: %d -> %d", id, next.Props.OnToggle)
	}
	if !r.DispatchBool(id, true) || !got {
		t.Fatal("toggle not dispatched")
	}
	r.releaseRemovedHandlers(next, nil)
	if r.DispatchBool(id, false) {
		t.Fatal("removed toggle dispatched")
	}
}

func TestGestureHandlersRetainPerIndexIdentityAndEncodeAllIntents(t *testing.T) {
	r := New(func() ui.Component { return nil }, &recordingRenderer{})
	calls := 0
	makeNode := func(delta int) *ui.Node {
		return ui.WithAnimation(ui.WithGesture(ui.WithGesture(ui.Text("target"), ui.GestureIntent{Kind: ui.GestureTap, Handler: func(ui.GestureEvent) { calls += delta }}), ui.GestureIntent{Kind: ui.GestureLongPress, MinimumPress: 400 * time.Millisecond, Handler: func(ui.GestureEvent) { calls += delta * 10 }}), ui.AnimationIntent{Property: ui.AnimateScale, Duration: 180 * time.Millisecond, Curve: ui.CurveSpring, SpringDamping: .8, ReduceMotionOK: true}).Build()
	}
	old := makeNode(1)
	r.bindHandlers(nil, old)
	if len(old.GestureHandlerIDs) != 2 || len(old.Props.Interactions) == 0 {
		t.Fatalf("missing interaction descriptors: %#v", old)
	}
	ids := append([]ui.HandlerID(nil), old.GestureHandlerIDs...)
	next := makeNode(2)
	r.bindHandlers(old, next)
	if next.GestureHandlerIDs[0] != ids[0] || next.GestureHandlerIDs[1] != ids[1] {
		t.Fatalf("ids changed: %v -> %v", ids, next.GestureHandlerIDs)
	}
	r.DispatchGesture(ids[0], ui.GestureEvent{TranslationX: 1})
	r.DispatchGesture(ids[1], ui.GestureEvent{VelocityY: 2})
	if calls != 22 {
		t.Fatalf("calls=%d", calls)
	}
	r.releaseRemovedHandlers(next, nil)
	if r.DispatchGesture(ids[0], ui.GestureEvent{}) {
		t.Fatal("removed gesture dispatched")
	}
}

func TestStopReleasesHandlersAndSuppressesRenders(t *testing.T) {
	renderer := &recordingRenderer{}
	r := New(func() ui.Component { return ui.Button("tap", func() {}) }, renderer)
	if err := r.Start(); err != nil {
		t.Fatal(err)
	}
	r.mu.Lock()
	handler := r.tree.Props.OnPress
	r.mu.Unlock()
	if handler == 0 || !r.events.Dispatch(handler) {
		t.Fatal("expected a live handler before stop")
	}
	r.Stop()
	if r.events.Dispatch(handler) {
		t.Fatal("handler remained live after stop")
	}
	r.Schedule()
	time.Sleep(10 * time.Millisecond)
	renderer.mu.Lock()
	defer renderer.mu.Unlock()
	if len(renderer.batches) != 1 {
		t.Fatalf("stop allowed another render: %d batches", len(renderer.batches))
	}
}

func TestNativeTimingAcknowledgement(t *testing.T) {
	var r *Runtime
	r = New(func() ui.Component { return ui.Text("timed") }, rendererFunc(func(batch MutationBatch) error {
		r.RecordNativeApply(batch.Sequence, 25*time.Microsecond)
		return nil
	}))
	if err := r.Start(); err != nil {
		t.Fatal(err)
	}
	samples := r.TimingSamples()
	if len(samples) != 1 || samples[0].Sequence == 0 || samples[0].MutationCount != 1 || samples[0].NativeApply != 25*time.Microsecond || samples[0].BridgeToApply <= 0 {
		t.Fatalf("unexpected timing samples: %#v", samples)
	}
}

func TestContextRuntimeRebuildsForEnvironment(t *testing.T) {
	renderer := &recordingRenderer{}
	r := NewContext(func(ctx ui.BuildContext) ui.Component { return ui.Text(ctx.Environment.Locale) }, renderer, ui.DefaultEnvironment())
	if err := r.Start(); err != nil {
		t.Fatal(err)
	}
	r.UpdateEnvironment(func(environment ui.Environment) ui.Environment { environment.Locale = "fr"; return environment })
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		renderer.mu.Lock()
		count := len(renderer.batches)
		renderer.mu.Unlock()
		if count >= 2 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	renderer.mu.Lock()
	defer renderer.mu.Unlock()
	if len(renderer.batches) < 2 {
		t.Fatalf("batches = %d", len(renderer.batches))
	}
	last := renderer.batches[len(renderer.batches)-1]
	if len(last.Mutations) != 1 || last.Mutations[0].Props.Text != "fr" {
		t.Fatalf("last batch = %#v", last)
	}
}

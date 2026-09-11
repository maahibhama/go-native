package ui

import (
	"testing"
	"time"
)

func TestPressableCreationAndDefaults(t *testing.T) {
	child := Text("Click Me")
	p := Pressable(child)
	node := p.Build()

	if node.Type != NodeView {
		t.Fatalf("expected NodeView, got %v", node.Type)
	}
	if node.Props.AccessRole != RoleButton {
		t.Fatalf("expected RoleButton, got %v", node.Props.AccessRole)
	}
	if len(node.Children) != 1 || node.Children[0].Props.Text != "Click Me" {
		t.Fatalf("unexpected children: %#v", node.Children)
	}

	// Verify minimum touch target constraints: 44x44pt on iOS/portable, 48x48dp on Android
	if node.Style.Layout.MinWidth.Value != 44 || node.Style.Layout.MinHeight.Value != 44 {
		t.Fatalf("expected portable 44x44 touch target, got %v x %v", node.Style.Layout.MinWidth, node.Style.Layout.MinHeight)
	}
	if node.Platform.IOS.Layout.MinWidth.Value != 44 || node.Platform.IOS.Layout.MinHeight.Value != 44 {
		t.Fatalf("expected iOS 44x44 touch target, got %v x %v", node.Platform.IOS.Layout.MinWidth, node.Platform.IOS.Layout.MinHeight)
	}
	if node.Platform.Android.Layout.MinWidth.Value != 48 || node.Platform.Android.Layout.MinHeight.Value != 48 {
		t.Fatalf("expected Android 48x48 touch target, got %v x %v", node.Platform.Android.Layout.MinWidth, node.Platform.Android.Layout.MinHeight)
	}

	// Verify Compact(true) clears touch target constraints
	compactNode := Pressable(child).Compact(true).Build()
	if compactNode.Style.Layout.MinWidth.Value != 0 || compactNode.Style.Layout.MinHeight.Value != 0 {
		t.Fatalf("expected cleared portable min size on Compact(true), got %v x %v", compactNode.Style.Layout.MinWidth, compactNode.Style.Layout.MinHeight)
	}
	if compactNode.Platform.IOS.Layout.MinWidth.Value != 0 || compactNode.Platform.Android.Layout.MinWidth.Value != 0 {
		t.Fatalf("expected cleared platform min size on Compact(true)")
	}

	// Verify Compact(false) restores touch target constraints
	restoredNode := Pressable(child).Compact(true).Compact(false).Build()
	if restoredNode.Style.Layout.MinWidth.Value != 44 || restoredNode.Platform.Android.Layout.MinWidth.Value != 48 {
		t.Fatalf("expected restored touch target constraints on Compact(false)")
	}
}

func TestPressableOnPressAndGestures(t *testing.T) {
	pressed := false
	p := Pressable(Text("Press")).OnPress(func() {
		pressed = true
	})
	node := p.Build()

	if node.Press == nil {
		t.Fatal("expected non-nil node.Press")
	}
	node.Press()
	if !pressed {
		t.Fatal("node.Press() did not trigger onPress callback")
	}

	// Gesture intent verification
	if len(node.Intents.Gestures) == 0 || node.Intents.Gestures[0].Kind != GestureTap {
		t.Fatalf("expected GestureTap intent, got %#v", node.Intents.Gestures)
	}

	gestureTriggered := false
	p2 := Pressable(Text("Tap")).OnPress(func() {
		gestureTriggered = true
	})
	node2 := p2.Build()
	node2.Intents.Gestures[0].Handler(GestureEvent{})
	if !gestureTriggered {
		t.Fatal("gesture handler did not trigger onPress callback")
	}

	// Disabled pressable should not invoke callbacks
	disabledPressed := false
	disabledNode := Pressable(Text("Disabled")).OnPress(func() {
		disabledPressed = true
	}).Disabled(true).Build()

	disabledNode.Press()
	if disabledPressed {
		t.Fatal("disabled pressable should not invoke node.Press")
	}
	if len(disabledNode.Intents.Gestures) > 0 {
		disabledNode.Intents.Gestures[0].Handler(GestureEvent{})
		if disabledPressed {
			t.Fatal("disabled pressable should not invoke gesture handler")
		}
	}
}

func TestPressableOnLongPress(t *testing.T) {
	longPressed := false
	duration := 450 * time.Millisecond
	p := Pressable(Text("Long Press")).OnLongPress(duration, func() {
		longPressed = true
	})
	node := p.Build()

	if len(node.Intents.Gestures) == 0 {
		t.Fatal("expected long press gesture intent")
	}
	intent := node.Intents.Gestures[0]
	if intent.Kind != GestureLongPress || intent.MinimumPress != duration {
		t.Fatalf("unexpected gesture intent: %#v", intent)
	}

	intent.Handler(GestureEvent{})
	if !longPressed {
		t.Fatal("long press handler did not fire callback")
	}
}

func TestPressablePressInAndPressOut(t *testing.T) {
	inCalled := false
	outCalled := false
	p := Pressable(Text("Touch")).
		OnPressIn(func() { inCalled = true }).
		OnPressOut(func() { outCalled = true })
	node := p.Build()

	if node.PressIn == nil || node.PressOut == nil {
		t.Fatal("expected non-nil PressIn and PressOut")
	}
	node.PressIn()
	if !inCalled {
		t.Fatal("PressIn did not fire")
	}
	node.PressOut()
	if !outCalled {
		t.Fatal("PressOut did not fire")
	}
}

func TestPressableFeedback(t *testing.T) {
	feedbacks := []FeedbackType{
		FeedbackDefault,
		FeedbackNone,
		FeedbackRipple,
		FeedbackHighlight,
		FeedbackHaptic,
	}
	for _, fb := range feedbacks {
		p := Pressable(Text("Feedback")).Feedback(fb)
		node := p.Build()
		if node.Feedback != fb {
			t.Fatalf("expected feedback %v, got %v", fb, node.Feedback)
		}
	}
}

func TestPressableBuilder(t *testing.T) {
	builder := PressableBuilder(func(state InteractionState) Component {
		if state.Disabled {
			return Text("State: Disabled")
		}
		if state.Pressed {
			return Text("State: Pressed")
		}
		if state.Focused {
			return Text("State: Focused")
		}
		if state.Hovered {
			return Text("State: Hovered")
		}
		return Text("State: Default")
	})

	// Default build should have zero InteractionState
	node := builder.Build()
	if node.Props.Text != "State: Default" {
		t.Fatalf("expected 'State: Default', got %q", node.Props.Text)
	}

	// Test with explicit InteractionState via WithInteractionState
	ctx := NewBuildContext(DefaultEnvironment())
	ctxPressed := WithInteractionState(ctx, InteractionState{Pressed: true})
	pressedNode := BuildWithContext(builder, ctxPressed)
	if pressedNode.Props.Text != "State: Pressed" {
		t.Fatalf("expected 'State: Pressed', got %q", pressedNode.Props.Text)
	}

	ctxFocused := WithInteractionState(ctx, InteractionState{Focused: true})
	focusedNode := BuildWithContext(builder, ctxFocused)
	if focusedNode.Props.Text != "State: Focused" {
		t.Fatalf("expected 'State: Focused', got %q", focusedNode.Props.Text)
	}

	// Pressable with Disabled(true) propagating to child PressableBuilder
	parent := Pressable(builder).Disabled(true)
	parentNode := parent.Build()
	if len(parentNode.Children) != 1 || parentNode.Children[0].Props.Text != "State: Disabled" {
		t.Fatalf("expected child to receive disabled state, got %#v", parentNode.Children)
	}
}

func TestButtons(t *testing.T) {
	actionCalled := false
	onPress := func() { actionCalled = true }

	variants := []struct {
		name string
		btn  *element
	}{
		{"FilledButton", FilledButton("Filled", onPress)},
		{"OutlinedButton", OutlinedButton("Outlined", onPress)},
		{"TextButton", TextButton("Text", onPress)},
		{"ElevatedButton", ElevatedButton("Elevated", onPress)},
		{"TonalButton", TonalButton("Tonal", onPress)},
		{"IconButton", IconButton("star", onPress)},
	}

	for _, tc := range variants {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.btn.Build()
			if node.Type != NodeButton {
				t.Fatalf("expected NodeButton, got %v", node.Type)
			}
			if node.Props.AccessRole != RoleButton {
				t.Fatalf("expected RoleButton, got %v", node.Props.AccessRole)
			}
			if node.Press == nil {
				t.Fatal("expected non-nil node.Press")
			}
			// Verify accessible touch target defaults
			if node.Style.Layout.MinHeight.Value != 44 || node.Platform.Android.Layout.MinHeight.Value != 48 {
				t.Fatalf("accessible touch target not enforced: %v", node.Style.Layout.MinHeight)
			}
		})
	}

	// Verify action invocation
	btn := FilledButton("Submit", onPress)
	node := btn.Build()
	node.Press()
	if !actionCalled {
		t.Fatal("button onPress did not fire")
	}

	// Verify OutlinedButton border
	outlined := OutlinedButton("Outline", nil).Build()
	if outlined.Style.Appearance.Border.Width != 1 {
		t.Fatalf("expected border width 1, got %v", outlined.Style.Appearance.Border.Width)
	}

	// Verify ElevatedButton shadow
	elevated := ElevatedButton("Elevated", nil).Build()
	if elevated.Style.Appearance.Shadow.Blur <= 0 {
		t.Fatalf("expected elevation shadow, got %#v", elevated.Style.Appearance.Shadow)
	}

	// Verify button modifiers: Leading, Trailing, Loading, Compact
	fullBtn := FilledButton("Full", nil).
		Leading(Text("Icon")).
		Trailing(Text("Arrow")).
		Compact(true)

	fullNode := fullBtn.Build()
	if len(fullNode.Children) != 2 {
		t.Fatalf("expected 2 children (leading and trailing), got %d", len(fullNode.Children))
	}
	if fullNode.Children[0].Props.Text != "Icon" || fullNode.Children[1].Props.Text != "Arrow" {
		t.Fatalf("unexpected children text: %#v", fullNode.Children)
	}
	if fullNode.Style.Layout.MinHeight.Value != 0 {
		t.Fatalf("expected Compact button to have min-height 0, got %v", fullNode.Style.Layout.MinHeight.Value)
	}

	// Verify Loading state
	loadingBtn := FilledButton("Saving", nil).Loading(true)
	loadingNode := loadingBtn.Build()
	if !loadingNode.Style.Interaction.Disabled {
		t.Fatal("loading button should be disabled")
	}
	if len(loadingNode.Children) != 1 || loadingNode.Children[0].Type != NodeProgressIndicator {
		t.Fatalf("expected loading button to include progress indicator child, got %#v", loadingNode.Children)
	}
	if !loadingNode.Children[0].Indeterminate() {
		t.Fatal("loading indicator should be indeterminate")
	}
}

func TestLinks(t *testing.T) {
	link := Link("Visit Site", "https://example.com")
	linkNode := link.Build()

	if linkNode.Type != NodeText {
		t.Fatalf("expected NodeText for Link, got %v", linkNode.Type)
	}
	if linkNode.Props.AccessRole != RoleButton {
		t.Fatalf("expected RoleButton for Link, got %v", linkNode.Props.AccessRole)
	}
	if linkNode.Props.Text != "Visit Site" {
		t.Fatalf("expected 'Visit Site', got %q", linkNode.Props.Text)
	}
	if len(linkNode.Props.RichText) == 0 {
		t.Fatal("expected RichText payload for link underline and target")
	}

	actionFired := false
	linkAction := LinkAction("Perform Action", func() {
		actionFired = true
	})
	actionNode := linkAction.Build()

	if actionNode.Type != NodeText {
		t.Fatalf("expected NodeText for LinkAction, got %v", actionNode.Type)
	}
	if actionNode.Props.AccessRole != RoleButton {
		t.Fatalf("expected RoleButton for LinkAction, got %v", actionNode.Props.AccessRole)
	}
	if actionNode.Press == nil {
		t.Fatal("expected non-nil node.Press for LinkAction")
	}
	actionNode.Press()
	if !actionFired {
		t.Fatal("LinkAction node.Press() did not fire callback")
	}

	// Verify link callback is also wired
	actionFired2 := false
	linkAction2 := LinkAction("Action 2", func() {
		actionFired2 = true
	})
	actionNode2 := linkAction2.Build()
	if actionNode2.Link != nil {
		actionNode2.Link("action://press")
		if !actionFired2 {
			t.Fatal("LinkAction link handler did not fire callback")
		}
	}
}

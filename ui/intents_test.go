package ui

import (
	"testing"
	"time"
)

func TestInteractionIntentsComposeWithoutChangingNode(t *testing.T) {
	button := Button("Save", func() {})
	decorated := WithAnimation(WithGesture(button, GestureIntent{
		Kind: GestureLongPress, MinimumPress: 400 * time.Millisecond,
	}), AnimationIntent{
		Property: AnimateScale, Duration: 180 * time.Millisecond, Curve: CurveSpring,
		SpringDamping: .8, ReduceMotionOK: true,
		From: .8, To: 1, FromX: 1, FromY: 2, ToX: 3, ToY: 4,
	})

	intents := IntentsOf(decorated)
	if len(intents.Gestures) != 1 || intents.Gestures[0].Kind != GestureLongPress {
		t.Fatalf("unexpected gestures: %#v", intents.Gestures)
	}
	if len(intents.Animations) != 1 || intents.Animations[0].Property != AnimateScale {
		t.Fatalf("unexpected animations: %#v", intents.Animations)
	}
	if got, want := decorated.Build().Type, button.Build().Type; got != want {
		t.Fatalf("decoration changed node type: got %v want %v", got, want)
	}
}

func TestInteractionIntentsAreCopied(t *testing.T) {
	component := WithGesture(Text("Swipe"), GestureIntent{Kind: GestureSwipe, Direction: SwipeLeading, MinimumTravel: 24})
	first := IntentsOf(component)
	first.Gestures[0].Direction = SwipeTrailing
	if got := IntentsOf(component).Gestures[0].Direction; got != SwipeLeading {
		t.Fatalf("intent metadata was mutable through returned slice: %v", got)
	}
}

func TestInteractionIntentsPreserveOrder(t *testing.T) {
	component := WithGesture(Text("Target"), GestureIntent{Kind: GestureTap})
	component = WithGesture(component, GestureIntent{Kind: GestureDrag})
	component = WithAnimation(component, AnimationIntent{Property: AnimateOpacity})
	component = WithAnimation(component, AnimationIntent{Property: AnimateLayout})

	intents := IntentsOf(component)
	if intents.Gestures[0].Kind != GestureTap || intents.Gestures[1].Kind != GestureDrag {
		t.Fatalf("gesture order changed: %#v", intents.Gestures)
	}
	if intents.Animations[0].Property != AnimateOpacity || intents.Animations[1].Property != AnimateLayout {
		t.Fatalf("animation order changed: %#v", intents.Animations)
	}
}

func TestIntentsOfPlainAndNilComponents(t *testing.T) {
	if got := IntentsOf(Text("plain")); len(got.Gestures)+len(got.Animations) != 0 {
		t.Fatalf("plain component has intents: %#v", got)
	}
	if got := WithGesture(nil, GestureIntent{Kind: GestureTap}); got != nil {
		t.Fatalf("decorating nil returned %#v", got)
	}
}

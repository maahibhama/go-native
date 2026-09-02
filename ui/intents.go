package ui

import "time"

// GestureKind identifies a high-level gesture without exposing platform APIs.
type GestureKind uint8

const (
	GestureTap GestureKind = iota + 1
	GestureLongPress
	GestureSwipe
	GestureDrag
)

// SwipeDirection is a logical rather than device-orientation-specific direction.
type SwipeDirection uint8

const (
	SwipeAny SwipeDirection = iota
	SwipeUp
	SwipeDown
	SwipeLeading
	SwipeTrailing
)

// GestureEvent is the portable value delivered when a gesture is recognized.
// Translation and velocity use logical points.
type GestureEvent struct {
	TranslationX float32
	TranslationY float32
	VelocityX    float32
	VelocityY    float32
}

// GestureIntent describes native gesture recognition requested by a component.
// Handler remains Go-owned and must be registered by the runtime before an
// intent crosses a native boundary.
type GestureIntent struct {
	Kind          GestureKind
	Direction     SwipeDirection
	MinimumPress  time.Duration
	MinimumTravel float32
	Handler       func(GestureEvent)
}

// AnimationProperty identifies a native property transition.
type AnimationProperty uint8

const (
	AnimateOpacity AnimationProperty = iota + 1
	AnimateScale
	AnimatePosition
	AnimateLayout
)

// AnimationCurve selects a portable timing curve.
type AnimationCurve uint8

const (
	CurveEaseInOut AnimationCurve = iota
	CurveEaseIn
	CurveEaseOut
	CurveLinear
	CurveSpring
)

// AnimationIntent describes how a property change should be presented.
// Duration and Delay are converted to platform-native timing units by the
// renderer. Spring damping and velocity apply only to CurveSpring.
type AnimationIntent struct {
	Property       AnimationProperty
	Duration       time.Duration
	Delay          time.Duration
	Curve          AnimationCurve
	SpringDamping  float32
	SpringVelocity float32
	ReduceMotionOK bool
	// From and To are scalar targets for opacity and scale.
	From float32
	To   float32
	// FromX/FromY and ToX/ToY are logical-point targets for position.
	FromX float32
	FromY float32
	ToX   float32
	ToY   float32
}

// IntentSet is the protocol-neutral interaction metadata attached to a component.
type IntentSet struct {
	Gestures   []GestureIntent
	Animations []AnimationIntent
}

// IntentComponent is a component carrying gesture or animation metadata.
type IntentComponent interface {
	Component
	InteractionIntents() IntentSet
}

type intentComponent struct {
	component Component
	intents   IntentSet
}

func (c *intentComponent) Build() *Node {
	n := c.component.Build()
	n.Intents = cloneIntents(c.intents)
	return n
}

func (c *intentComponent) InteractionIntents() IntentSet { return cloneIntents(c.intents) }

// WithGesture attaches a gesture intent. Repeated calls preserve existing intents.
func WithGesture(component Component, intent GestureIntent) Component {
	return decorateIntents(component, func(set *IntentSet) { set.Gestures = append(set.Gestures, intent) })
}

// WithAnimation attaches an animation intent. Repeated calls preserve existing intents.
func WithAnimation(component Component, intent AnimationIntent) Component {
	return decorateIntents(component, func(set *IntentSet) { set.Animations = append(set.Animations, intent) })
}

// IntentsOf returns a detached copy of a component's interaction metadata.
func IntentsOf(component Component) IntentSet {
	if decorated, ok := component.(IntentComponent); ok {
		return decorated.InteractionIntents()
	}
	return IntentSet{}
}

func decorateIntents(component Component, update func(*IntentSet)) Component {
	if component == nil {
		return nil
	}
	base := component
	set := IntentSet{}
	if decorated, ok := component.(*intentComponent); ok {
		base = decorated.component
		set = cloneIntents(decorated.intents)
	}
	update(&set)
	return &intentComponent{component: base, intents: set}
}

func cloneIntents(set IntentSet) IntentSet {
	set.Gestures = append([]GestureIntent(nil), set.Gestures...)
	set.Animations = append([]AnimationIntent(nil), set.Animations...)
	return set
}

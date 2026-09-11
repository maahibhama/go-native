package ui

import "time"

// InteractionState describes the interactive state of a control.
type InteractionState struct {
	Pressed  bool
	Hovered  bool
	Focused  bool
	Disabled bool
}

// FeedbackType selects native touch and haptic feedback behavior.
type FeedbackType uint8

const (
	FeedbackDefault FeedbackType = iota
	FeedbackNone
	FeedbackRipple
	FeedbackHighlight
	FeedbackHaptic
)

var interactionStateKey = DependencyKey[InteractionState]{Name: "go-native.interaction_state"}

func cloneDependencies(src *Dependencies) *Dependencies {
	dst := NewDependencies()
	if src != nil {
		src.mu.RLock()
		for k, v := range src.values {
			dst.values[k] = v
		}
		src.mu.RUnlock()
	}
	return dst
}

// Pressable creates an interactive native container with RoleButton accessibility.
func Pressable(children ...Component) *element {
	e := newElement(NodeView, Props{AccessRole: RoleButton}, children...)
	e.isPressable = true
	applyButtonTouchTargets(e)
	return e
}

// OnPress registers a tap callback and ensures native tap recognition is attached.
func (e *element) OnPress(fn func()) *element {
	if fn == nil {
		e.node.Press = nil
		return e
	}
	e.node.Press = func() {
		if !e.node.Style.Interaction.Disabled {
			fn()
		}
	}
	if e.node.Type != NodeButton {
		e.node.Intents.Gestures = append(e.node.Intents.Gestures, GestureIntent{
			Kind: GestureTap,
			Handler: func(GestureEvent) {
				if !e.node.Style.Interaction.Disabled {
					fn()
				}
			},
		})
	}
	return e
}

// OnLongPress attaches a long-press gesture recognizer with minimum duration.
func (e *element) OnLongPress(duration time.Duration, onLongPress func()) *element {
	if onLongPress == nil {
		return e
	}
	e.node.Intents.Gestures = append(e.node.Intents.Gestures, GestureIntent{
		Kind:         GestureLongPress,
		MinimumPress: duration,
		Handler: func(GestureEvent) {
			if !e.node.Style.Interaction.Disabled {
				onLongPress()
			}
		},
	})
	return e
}

// OnPressIn sets a callback invoked when a press interaction begins.
func (e *element) OnPressIn(fn func()) *element {
	e.node.PressIn = fn
	return e
}

// OnPressOut sets a callback invoked when a press interaction ends.
func (e *element) OnPressOut(fn func()) *element {
	e.node.PressOut = fn
	return e
}

// Feedback configures the native touch feedback effect.
func (e *element) Feedback(feedback FeedbackType) *element {
	e.node.Feedback = feedback
	return e
}

type pressableBuilderComponent struct {
	builder func(state InteractionState) Component
}

// PressableBuilder builds dynamic components based on the current InteractionState.
func PressableBuilder(builder func(state InteractionState) Component) Component {
	return &pressableBuilderComponent{builder: builder}
}

func (p *pressableBuilderComponent) Build() *Node {
	return p.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (p *pressableBuilderComponent) BuildContext(context BuildContext) *Node {
	if p == nil || p.builder == nil {
		return nil
	}
	state, _ := Resolve(context.Environment.Dependencies, interactionStateKey)
	child := p.builder(state)
	if child == nil {
		return nil
	}
	return BuildWithContext(child, context)
}

// WithInteractionState decorates a BuildContext with a given InteractionState.
func WithInteractionState(context BuildContext, state InteractionState) BuildContext {
	deps := cloneDependencies(context.Environment.Dependencies)
	Provide(deps, interactionStateKey, state)
	env := context.Environment
	env.Dependencies = deps
	return context.WithEnvironment(env)
}

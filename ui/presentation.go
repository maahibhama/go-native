package ui

// Route describes one mounted destination in a declarative navigation stack.
// Key must be stable within the stack so native controllers can preserve state.
type Route struct {
	Key     string
	Title   string
	Content Component
}

// NavigationIntent describes the complete desired native navigation stack.
// Routes are ordered root first and visible destination last.
type NavigationIntent struct {
	Routes []Route
}

// NavigationComponent carries a navigation stack without exposing platform
// controller types.
type NavigationComponent interface {
	Component
	NavigationIntent() NavigationIntent
}

type navigationComponent struct{ intent NavigationIntent }

// NavigationStack declares a native navigation stack. Until native mounting is
// integrated, Build returns the visible route's content so it remains usable by
// the existing renderer.
func NavigationStack(routes ...Route) Component {
	intent := NavigationIntent{Routes: append([]Route(nil), routes...)}
	return &navigationComponent{intent: intent}
}

func (c *navigationComponent) Build() *Node {
	return c.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (c *navigationComponent) BuildContext(context BuildContext) *Node {
	for i := len(c.intent.Routes) - 1; i >= 0; i-- {
		if content := c.intent.Routes[i].Content; content != nil {
			return BuildWithContext(content, context.Child("route:"+c.intent.Routes[i].Key))
		}
	}
	return nil
}

func (c *navigationComponent) NavigationIntent() NavigationIntent {
	return cloneNavigation(c.intent)
}

// NavigationOf returns detached navigation metadata when component declares a stack.
func NavigationOf(component Component) (NavigationIntent, bool) {
	value, ok := component.(NavigationComponent)
	if !ok {
		return NavigationIntent{}, false
	}
	return value.NavigationIntent(), true
}

func cloneNavigation(intent NavigationIntent) NavigationIntent {
	intent.Routes = append([]Route(nil), intent.Routes...)
	return intent
}

// ModalStyle requests a platform-appropriate presentation style.
type ModalStyle uint8

const (
	ModalAutomatic ModalStyle = iota
	ModalSheet
	ModalFullscreen
)

// ModalIntent describes one presented destination. Key identifies presentation
// identity across renders. OnDismiss is invoked after native dismissal completes.
type ModalIntent struct {
	Key         string
	Content     Component
	Style       ModalStyle
	Dismissible bool
	OnDismiss   func()
}

// ModalComponent carries optional modal presentation metadata.
type ModalComponent interface {
	Component
	ModalIntent() (ModalIntent, bool)
}

type modalComponent struct {
	base   Component
	intent ModalIntent
}

// PresentModal decorates base with a desired modal presentation. Passing a nil
// base is supported for applications whose root content is entirely presented.
func PresentModal(base Component, intent ModalIntent) Component {
	return &modalComponent{base: base, intent: intent}
}

func (c *modalComponent) Build() *Node {
	return c.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (c *modalComponent) BuildContext(context BuildContext) *Node {
	if c.base == nil {
		return nil
	}
	return BuildWithContext(c.base, context.Child("modal-base"))
}

func (c *modalComponent) ModalIntent() (ModalIntent, bool) { return c.intent, c.intent.Content != nil }

// ModalOf returns presentation metadata and whether a modal is requested.
func ModalOf(component Component) (ModalIntent, bool) {
	value, ok := component.(ModalComponent)
	if !ok {
		return ModalIntent{}, false
	}
	return value.ModalIntent()
}

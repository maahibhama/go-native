package ui

import (
	"strconv"
	"sync"
	"time"
)

// ScaffoldProps defines the configuration for a Scaffold layout.
type ScaffoldProps struct {
	AppBar               Component
	Body                 Component
	BottomBar            Component
	FloatingActionButton Component
	BackgroundColor      Color
}

type scaffoldComponent struct {
	props ScaffoldProps
}

// Scaffold arranges screen-level structural components: an optional AppBar at the top,
// a flexible Body in the middle, an optional BottomBar at the bottom, and an optional
// FloatingActionButton overlaid at the bottom-trailing edge.
func Scaffold(props ScaffoldProps) Component {
	return &scaffoldComponent{props: props}
}

func (s *scaffoldComponent) Props() ScaffoldProps {
	return s.props
}

func (s *scaffoldComponent) Build() *Node {
	return s.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (s *scaffoldComponent) BuildContext(context BuildContext) *Node {
	var children []Component
	if s.props.AppBar != nil {
		children = append(children, s.props.AppBar)
	}
	if s.props.Body != nil {
		children = append(children, View(s.props.Body).Flex(1, 1))
	} else {
		children = append(children, Spacer())
	}
	if s.props.BottomBar != nil {
		children = append(children, s.props.BottomBar)
	}

	mainColumn := Column(children...).Flex(1, 1)

	var root Component
	if s.props.FloatingActionButton != nil {
		fab := View(s.props.FloatingActionButton).Absolute(EdgeInsets{Bottom: 24, Trailing: 24})
		root = Stack(mainColumn, fab)
	} else {
		root = mainColumn
	}

	node := BuildWithContext(root, context.Child("scaffold"))
	if node == nil {
		return nil
	}
	if !s.props.BackgroundColor.IsTransparent() {
		node.Style.Appearance.Background = s.props.BackgroundColor
	}
	node.Style.Layout.FlexGrow = 1
	node.Style.Layout.FlexShrink = 1
	return node
}

// AppBarProps configures an application header bar.
type AppBarProps struct {
	Title     string
	Leading   Component
	Actions   []Component
	Elevation float32
}

type appBarComponent struct {
	props AppBarProps
}

// AppBar creates a standard header bar with title and trailing action items.
func AppBar(title string, actions ...Component) Component {
	return AppBarWithProps(AppBarProps{
		Title:   title,
		Actions: actions,
	})
}

// AppBarWithProps creates a header bar with full configuration including leading component and elevation.
func AppBarWithProps(props AppBarProps) Component {
	return &appBarComponent{props: props}
}

func (a *appBarComponent) Props() AppBarProps {
	return a.props
}

func (a *appBarComponent) Build() *Node {
	return a.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (a *appBarComponent) BuildContext(context BuildContext) *Node {
	var rowChildren []Component
	if a.props.Leading != nil {
		rowChildren = append(rowChildren, a.props.Leading)
	}
	if a.props.Title != "" {
		titleElem := Text(a.props.Title).Bold().FontSize(20)
		rowChildren = append(rowChildren, titleElem)
	}
	rowChildren = append(rowChildren, Spacer())
	for _, act := range a.props.Actions {
		if act != nil {
			rowChildren = append(rowChildren, act)
		}
	}

	bar := Row(rowChildren...).
		Height(56).
		PaddingXY(16, 8).
		Align(AlignCenter)

	if a.props.Elevation > 0 {
		bar = bar.Shadow(Shadow{
			Color:   RGBA(0, 0, 0, 50),
			Offset:  Point{X: 0, Y: a.props.Elevation / 2},
			Blur:    a.props.Elevation,
			Opacity: 0.2,
		})
	}

	return BuildWithContext(bar, context.Child("appbar"))
}

// FormValidationState manages error state and validation status for form fields.
type FormValidationState struct {
	mu     sync.RWMutex
	errors map[string]string
}

// NewFormValidationState creates an empty, valid form validation container.
func NewFormValidationState() *FormValidationState {
	return &FormValidationState{errors: make(map[string]string)}
}

// SetError records an error message for a specific field name.
func (s *FormValidationState) SetError(field, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.errors == nil {
		s.errors = make(map[string]string)
	}
	s.errors[field] = message
}

// ClearError removes any error message associated with field.
func (s *FormValidationState) ClearError(field string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.errors, field)
}

// ClearAll removes all recorded field errors.
func (s *FormValidationState) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errors = make(map[string]string)
}

// Error returns the error message for field, or an empty string if valid.
func (s *FormValidationState) Error(field string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errors[field]
}

// Errors returns a shallow copy of all active field errors.
func (s *FormValidationState) Errors() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]string, len(s.errors))
	for k, v := range s.errors {
		copy[k] = v
	}
	return copy
}

// IsValid reports true when no validation errors are present.
func (s *FormValidationState) IsValid() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.errors) == 0
}

// FormValidationKey is the dependency injection key for resolving FormValidationState from context.
var FormValidationKey = DependencyKey[*FormValidationState]{Name: "ui.FormValidationState"}

// FormProps configures optional submit handling and validation state for Form.
type FormProps struct {
	OnSubmit   func()
	Validation *FormValidationState
}

// FormComponent is a container for form inputs with submit handling and validation state.
type FormComponent struct {
	children   []Component
	onSubmit   func()
	validation *FormValidationState
}

// Form creates a form container with validation state and submit handling.
func Form(children ...Component) *FormComponent {
	return &FormComponent{
		children:   append([]Component(nil), children...),
		validation: NewFormValidationState(),
	}
}

// FormWithProps creates a FormComponent with explicit submit handler and validation state.
func FormWithProps(props FormProps, children ...Component) *FormComponent {
	f := Form(children...)
	f.onSubmit = props.OnSubmit
	if props.Validation != nil {
		f.validation = props.Validation
	}
	return f
}

// OnSubmit registers a callback invoked when the form is submitted.
func (f *FormComponent) OnSubmit(handler func()) *FormComponent {
	f.onSubmit = handler
	return f
}

// Submit executes the registered submit callback.
func (f *FormComponent) Submit() {
	if f.onSubmit != nil {
		f.onSubmit()
	}
}

// ValidationState returns the validation container owned by this form.
func (f *FormComponent) ValidationState() *FormValidationState {
	return f.validation
}

// SetError records an error for field in this form's validation container.
func (f *FormComponent) SetError(field, message string) *FormComponent {
	f.validation.SetError(field, message)
	return f
}

// ClearError removes an error for field in this form's validation container.
func (f *FormComponent) ClearError(field string) *FormComponent {
	f.validation.ClearError(field)
	return f
}

// ClearErrors removes all errors in this form's validation container.
func (f *FormComponent) ClearErrors() *FormComponent {
	f.validation.ClearAll()
	return f
}

// Error retrieves the error message for field.
func (f *FormComponent) Error(field string) string {
	return f.validation.Error(field)
}

// Errors returns a snapshot of all active field errors.
func (f *FormComponent) Errors() map[string]string {
	return f.validation.Errors()
}

// IsValid reports true if this form has no active errors.
func (f *FormComponent) IsValid() bool {
	return f.validation.IsValid()
}

func (f *FormComponent) Build() *Node {
	return f.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (f *FormComponent) BuildContext(context BuildContext) *Node {
	if context.Environment.Dependencies != nil && f.validation != nil {
		Provide(context.Environment.Dependencies, FormValidationKey, f.validation)
	}
	col := Column(f.children...).Gap(8)
	node := BuildWithContext(col, context.Child("form"))
	if node != nil {
		node.Submit = f.onSubmit
	}
	return node
}

// Card creates a card container with rounded corners, padding, and elevated styling.
func Card(children ...Component) *element {
	e := Column(children...).
		Padding(16).
		CornerRadius(12)
	return e.CardElevated()
}

// CardElevated applies elevation shadow and a solid surface background to the element.
func (e *element) CardElevated() *element {
	e.node.Style.Appearance.Background = RGB(255, 255, 255)
	e.node.Style.Appearance.Border = Border{}
	e.node.Style.Appearance.Shadow = Shadow{
		Color:   RGBA(0, 0, 0, 30),
		Offset:  Point{X: 0, Y: 2},
		Blur:    8,
		Opacity: 0.15,
	}
	return e
}

// CardOutlined applies a subtle border and clears elevation shadow on the element.
func (e *element) CardOutlined() *element {
	e.node.Style.Appearance.Background = RGB(255, 255, 255)
	e.node.Style.Appearance.Border = Border{
		Width: 1,
		Color: RGB(224, 224, 224),
	}
	e.node.Style.Appearance.Shadow = Shadow{}
	return e
}

// CardFilled applies a filled neutral background with no border or elevation shadow.
func (e *element) CardFilled() *element {
	e.node.Style.Appearance.Background = RGB(245, 245, 245)
	e.node.Style.Appearance.Border = Border{}
	e.node.Style.Appearance.Shadow = Shadow{}
	return e
}

type badgeComponent struct {
	content Component
	count   int
}

// Badge overlays a numeric count pill on content. If count is <= 0, only content is rendered.
func Badge(content Component, count int) Component {
	return &badgeComponent{content: content, count: count}
}

func (b *badgeComponent) Count() int {
	return b.count
}

func (b *badgeComponent) Content() Component {
	return b.content
}

func (b *badgeComponent) Build() *Node {
	return b.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (b *badgeComponent) BuildContext(context BuildContext) *Node {
	if b.content == nil && b.count <= 0 {
		return nil
	}
	if b.count <= 0 {
		return BuildWithContext(b.content, context.Child("badge-content"))
	}
	label := strconv.Itoa(b.count)
	if b.count > 99 {
		label = "99+"
	}
	badgeLabel := Text(label).FontSize(11).Bold().Foreground(RGB(255, 255, 255))
	badgeIndicator := View(badgeLabel).
		Background(RGB(255, 59, 48)).
		CornerRadius(9).
		PaddingXY(6, 2).
		Absolute(EdgeInsets{Top: -6, Trailing: -6})

	if b.content == nil {
		return BuildWithContext(badgeIndicator, context.Child("badge-indicator"))
	}
	stack := Stack(b.content, badgeIndicator)
	return BuildWithContext(stack, context.Child("badge"))
}

type dotBadgeComponent struct {
	content Component
	visible bool
}

// DotBadge overlays a small status dot indicator on content when visible is true.
func DotBadge(content Component, visible bool) Component {
	return &dotBadgeComponent{content: content, visible: visible}
}

func (d *dotBadgeComponent) Visible() bool {
	return d.visible
}

func (d *dotBadgeComponent) Content() Component {
	return d.content
}

func (d *dotBadgeComponent) Build() *Node {
	return d.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (d *dotBadgeComponent) BuildContext(context BuildContext) *Node {
	if d.content == nil && !d.visible {
		return nil
	}
	if !d.visible {
		return BuildWithContext(d.content, context.Child("dot-badge-content"))
	}
	dot := View().
		Width(8).
		Height(8).
		CornerRadius(4).
		Background(RGB(255, 59, 48)).
		Absolute(EdgeInsets{Top: -2, Trailing: -2})

	if d.content == nil {
		return BuildWithContext(dot, context.Child("dot-badge-indicator"))
	}
	stack := Stack(d.content, dot)
	return BuildWithContext(stack, context.Child("dot-badge"))
}

type avatarComponent struct {
	source   string
	initials string
	size     float32
}

// Avatar renders a circular user avatar image, falling back to initials when source is empty.
func Avatar(source string, initials string, size float32) Component {
	if size <= 0 {
		size = 40
	}
	return &avatarComponent{
		source:   source,
		initials: initials,
		size:     size,
	}
}

func (a *avatarComponent) Source() string {
	return a.source
}

func (a *avatarComponent) Initials() string {
	return a.initials
}

func (a *avatarComponent) Size() float32 {
	return a.size
}

func (a *avatarComponent) Build() *Node {
	return a.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (a *avatarComponent) BuildContext(context BuildContext) *Node {
	radius := a.size / 2
	if a.source != "" {
		img := Image(a.source).
			Width(a.size).
			Height(a.size).
			CornerRadius(radius).
			ResizeMode(ImageFill)
		return BuildWithContext(img, context.Child("avatar-image"))
	}

	fontSize := a.size * 0.4
	if fontSize < 10 {
		fontSize = 10
	}
	initialsText := Text(a.initials).
		Bold().
		FontSize(fontSize).
		Foreground(RGB(255, 255, 255))

	placeholder := View(Center(initialsText)).
		Width(a.size).
		Height(a.size).
		CornerRadius(radius).
		Background(RGB(142, 142, 147)).
		Overflow(OverflowHidden)

	return BuildWithContext(placeholder, context.Child("avatar-fallback"))
}

// Skeleton creates a loading placeholder element styled with pulse opacity animation.
func Skeleton(width, height float32) *element {
	e := View().
		Width(width).
		Height(height).
		CornerRadius(4).
		Background(RGB(225, 228, 232))

	e.node.Intents.Animations = append(e.node.Intents.Animations, AnimationIntent{
		Property: AnimateOpacity,
		From:     0.4,
		To:       1.0,
		Duration: 1000 * time.Millisecond,
		Curve:    CurveEaseInOut,
	})
	return e
}

type emptyStateComponent struct {
	title       string
	description string
	icon        Component
	action      Component
}

// EmptyState displays a centered placeholder with title, description, and optional icon and action.
func EmptyState(title, description string, icon Component, action Component) Component {
	return &emptyStateComponent{
		title:       title,
		description: description,
		icon:        icon,
		action:      action,
	}
}

func (e *emptyStateComponent) Title() string {
	return e.title
}

func (e *emptyStateComponent) Description() string {
	return e.description
}

func (e *emptyStateComponent) Icon() Component {
	return e.icon
}

func (e *emptyStateComponent) Action() Component {
	return e.action
}

func (e *emptyStateComponent) Build() *Node {
	return e.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (e *emptyStateComponent) BuildContext(context BuildContext) *Node {
	var children []Component
	if e.icon != nil {
		children = append(children, e.icon)
	}
	if e.title != "" {
		children = append(children, Text(e.title).Bold().FontSize(20).Align(AlignCenter))
	}
	if e.description != "" {
		children = append(children, Text(e.description).FontSize(14).Foreground(RGB(128, 128, 128)).Align(AlignCenter))
	}
	if e.action != nil {
		children = append(children, e.action)
	}

	content := Column(children...).
		Gap(12).
		Align(AlignCenter).
		Padding(24)

	return BuildWithContext(Center(content), context.Child("empty-state"))
}

type errorStateComponent struct {
	title   string
	message string
	onRetry func()
}

// ErrorState displays an error notification screen with title, message, and optional retry action.
func ErrorState(title, message string, onRetry func()) Component {
	return &errorStateComponent{
		title:   title,
		message: message,
		onRetry: onRetry,
	}
}

func (e *errorStateComponent) Title() string {
	return e.title
}

func (e *errorStateComponent) Message() string {
	return e.message
}

func (e *errorStateComponent) OnRetry() func() {
	return e.onRetry
}

func (e *errorStateComponent) Build() *Node {
	return e.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (e *errorStateComponent) BuildContext(context BuildContext) *Node {
	var children []Component
	if e.title != "" {
		children = append(children, Text(e.title).Bold().FontSize(20).Foreground(RGB(218, 54, 51)).Align(AlignCenter))
	}
	if e.message != "" {
		children = append(children, Text(e.message).FontSize(14).Foreground(RGB(128, 128, 128)).Align(AlignCenter))
	}
	if e.onRetry != nil {
		children = append(children, Button("Retry", e.onRetry))
	}

	content := Column(children...).
		Gap(12).
		Align(AlignCenter).
		Padding(24)

	return BuildWithContext(Center(content), context.Child("error-state"))
}

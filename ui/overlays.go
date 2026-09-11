package ui

// DialogProps configures an alert or confirmation dialog.
type DialogProps struct {
	Title        string
	Message      string
	Content      Component
	ConfirmLabel string
	OnConfirm    func()
	CancelLabel  string
	OnCancel     func()
	OnDismiss    func()
}

type dialogComponent struct {
	props DialogProps
}

// Dialog creates an alert/confirmation dialog with backdrop, title, message, optional content, and action buttons.
func Dialog(props DialogProps) Component {
	return &dialogComponent{props: props}
}

func (d *dialogComponent) Props() DialogProps {
	return d.props
}

func (d *dialogComponent) ModalIntent() (ModalIntent, bool) {
	return ModalIntent{
		Content:     d,
		Style:       ModalAutomatic,
		Dismissible: d.props.OnDismiss != nil,
		OnDismiss:   d.props.OnDismiss,
	}, true
}

func (d *dialogComponent) Build() *Node {
	return d.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (d *dialogComponent) BuildContext(context BuildContext) *Node {
	var cardChildren []Component
	if d.props.Title != "" {
		cardChildren = append(cardChildren, Text(d.props.Title).Bold().FontSize(18))
	}
	if d.props.Message != "" {
		cardChildren = append(cardChildren, Text(d.props.Message).FontSize(14).Foreground(RGB(100, 100, 100)))
	}
	if d.props.Content != nil {
		cardChildren = append(cardChildren, d.props.Content)
	}

	var actionButtons []Component
	if d.props.CancelLabel != "" || d.props.OnCancel != nil {
		label := d.props.CancelLabel
		if label == "" {
			label = "Cancel"
		}
		actionButtons = append(actionButtons, Button(label, d.props.OnCancel))
	}
	if d.props.ConfirmLabel != "" || d.props.OnConfirm != nil {
		label := d.props.ConfirmLabel
		if label == "" {
			label = "OK"
		}
		actionButtons = append(actionButtons, Button(label, d.props.OnConfirm))
	}
	if len(actionButtons) > 0 {
		actionsRow := Row(append([]Component{Spacer()}, actionButtons...)...).Gap(8)
		cardChildren = append(cardChildren, actionsRow)
	}

	dialogCard := Column(cardChildren...).
		Gap(16).
		Padding(20).
		CornerRadius(14).
		Background(RGB(255, 255, 255)).
		Shadow(Shadow{
			Color:   RGBA(0, 0, 0, 40),
			Offset:  Point{X: 0, Y: 4},
			Blur:    16,
			Opacity: 0.25,
		}).
		MinWidth(280).
		MaxWidth(400)

	overlay := View(Center(dialogCard)).
		Absolute(EdgeInsets{Top: 0, Leading: 0, Bottom: 0, Trailing: 0}).
		Background(RGBA(0, 0, 0, 100))

	return BuildWithContext(overlay, context.Child("dialog"))
}

// SheetProps configures a bottom sheet modal.
type SheetProps struct {
	Title     string
	Content   Component
	OnDismiss func()
}

type sheetComponent struct {
	props SheetProps
}

// Sheet creates a bottom-anchored modal sheet with drag indicator, title, and content.
func Sheet(props SheetProps) Component {
	return &sheetComponent{props: props}
}

func (s *sheetComponent) Props() SheetProps {
	return s.props
}

func (s *sheetComponent) ModalIntent() (ModalIntent, bool) {
	return ModalIntent{
		Content:     s.props.Content,
		Style:       ModalSheet,
		Dismissible: s.props.OnDismiss != nil,
		OnDismiss:   s.props.OnDismiss,
	}, true
}

func (s *sheetComponent) Build() *Node {
	return s.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (s *sheetComponent) BuildContext(context BuildContext) *Node {
	var sheetChildren []Component

	// Visual handle bar indicator
	handle := View().
		Width(36).
		Height(4).
		CornerRadius(2).
		Background(RGB(200, 200, 200))
	sheetChildren = append(sheetChildren, Row(Spacer(), handle, Spacer()).Align(AlignCenter))

	if s.props.Title != "" {
		sheetChildren = append(sheetChildren, Text(s.props.Title).Bold().FontSize(18))
	}
	if s.props.Content != nil {
		sheetChildren = append(sheetChildren, s.props.Content)
	}

	sheetCard := Column(sheetChildren...).
		Gap(12).
		PaddingXY(16, 12).
		CornerRadius(16).
		Background(RGB(255, 255, 255)).
		Shadow(Shadow{
			Color:   RGBA(0, 0, 0, 30),
			Offset:  Point{X: 0, Y: -2},
			Blur:    10,
			Opacity: 0.2,
		})

	overlay := Column(
		Spacer(),
		sheetCard,
	).
		Absolute(EdgeInsets{Top: 0, Leading: 0, Bottom: 0, Trailing: 0}).
		Background(RGBA(0, 0, 0, 100))

	return BuildWithContext(overlay, context.Child("sheet"))
}

type modalOverlayComponent struct {
	visible   bool
	onDismiss func()
	content   Component
}

// Modal creates a modal overlay that is rendered when visible is true.
func Modal(visible bool, onDismiss func(), content Component) Component {
	return &modalOverlayComponent{
		visible:   visible,
		onDismiss: onDismiss,
		content:   content,
	}
}

func (m *modalOverlayComponent) Visible() bool {
	return m.visible
}

func (m *modalOverlayComponent) OnDismiss() func() {
	return m.onDismiss
}

func (m *modalOverlayComponent) Content() Component {
	return m.content
}

func (m *modalOverlayComponent) ModalIntent() (ModalIntent, bool) {
	return ModalIntent{
		Content:     m.content,
		Style:       ModalAutomatic,
		Dismissible: m.onDismiss != nil,
		OnDismiss:   m.onDismiss,
	}, m.visible
}

func (m *modalOverlayComponent) Build() *Node {
	return m.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (m *modalOverlayComponent) BuildContext(context BuildContext) *Node {
	if !m.visible || m.content == nil {
		return nil
	}
	overlay := View(Center(m.content)).
		Absolute(EdgeInsets{Top: 0, Leading: 0, Bottom: 0, Trailing: 0}).
		Background(RGBA(0, 0, 0, 100))

	return BuildWithContext(overlay, context.Child("modal"))
}

type snackbarComponent struct {
	message     string
	actionLabel string
	onAction    func()
}

// Snackbar creates a floating feedback banner displaying a message and an optional action.
func Snackbar(message string, actionLabel string, onAction func()) Component {
	return &snackbarComponent{
		message:     message,
		actionLabel: actionLabel,
		onAction:    onAction,
	}
}

func (s *snackbarComponent) Message() string {
	return s.message
}

func (s *snackbarComponent) ActionLabel() string {
	return s.actionLabel
}

func (s *snackbarComponent) OnAction() func() {
	return s.onAction
}

func (s *snackbarComponent) Build() *Node {
	return s.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (s *snackbarComponent) BuildContext(context BuildContext) *Node {
	msg := Text(s.message).FontSize(14).Foreground(RGB(255, 255, 255))
	var rowChildren []Component
	rowChildren = append(rowChildren, msg)

	if s.actionLabel != "" && s.onAction != nil {
		rowChildren = append(rowChildren, Spacer(), Button(s.actionLabel, s.onAction))
	}

	bar := Row(rowChildren...).
		PaddingXY(16, 12).
		CornerRadius(8).
		Background(RGB(50, 50, 50)).
		Shadow(Shadow{
			Color:   RGBA(0, 0, 0, 50),
			Offset:  Point{X: 0, Y: 2},
			Blur:    8,
			Opacity: 0.2,
		}).
		Align(AlignCenter)

	return BuildWithContext(bar, context.Child("snackbar"))
}

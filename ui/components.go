package ui

// View creates a generic native container.
func View(children ...Component) *element { return newElement(NodeView, Props{}, children...) }

// Column creates a vertical native container.
func Column(children ...Component) *element { return newElement(NodeColumn, Props{}, children...) }

// Row creates a horizontal native container.
func Row(children ...Component) *element { return newElement(NodeRow, Props{}, children...) }

// Grid creates a native container whose Go-owned layout uses equal-width columns.
func Grid(columns int, children ...Component) *element {
	if columns < 1 {
		columns = 1
	}
	e := newElement(NodeView, Props{}, children...)
	e.node.Style.Layout.GridColumns = columns
	return e
}

// Text creates a native text label.
func Text(value string) *element { return newElement(NodeText, Props{Text: value}) }

// TextFunc evaluates value whenever the application rebuilds its virtual tree.
func TextFunc(value func() string) *element {
	if value == nil {
		return Text("")
	}
	return Text(value())
}

// SafeArea creates a container whose children avoid platform system insets.
func SafeArea(children ...Component) *element { return newElement(NodeSafeArea, Props{}, children...) }

// Button creates a native button. BindHandlers resolves its callback to a stable ID.
func Button(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label})
	e.node.Press = onPress
	return e
}

// TextInput creates an editable native text control. onChange receives user edits.
func TextInput(value string, onChange func(string)) *element {
	e := newElement(NodeTextInput, Props{Text: value})
	e.node.Change = onChange
	return e
}

// Switch creates a native boolean control.
func Switch(value bool, onChange func(bool)) *element {
	e := newElement(NodeSwitch, Props{Checked: value})
	e.node.Toggle = onChange
	return e
}

// ProgressIndicator creates a native determinate progress control, clamped to 0...1.
func ProgressIndicator(value float32) *element {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return newElement(NodeProgressIndicator, Props{Progress: value})
}

// Image creates a native image view using a platform bundle/resource name.
func Image(source string) *element {
	return newElement(NodeImage, Props{ImageSource: source, ImageMode: ImageFit})
}

// ResizeMode controls how an Image fits its requested bounds.
func (e *element) ResizeMode(value ImageResizeMode) *element {
	e.node.Props.ImageMode = value
	return e
}

// ScrollView creates a vertical single-child native scrolling container.
func ScrollView(child Component) *element { return newElement(NodeScrollView, Props{}, child) }

// HorizontalScroll makes a ScrollView scroll horizontally.
func (e *element) HorizontalScroll() *element { e.node.Props.Horizontal = true; return e }

func (e *element) Padding(value float32) *element {
	e.node.Props.Padding = value
	e.node.Style.Layout.Padding = Insets(value)
	return e
}
func (e *element) Gap(value float32) *element {
	e.node.Props.Gap = value
	e.node.Style.Layout.Gap = value
	return e
}
func (e *element) Width(value float32) *element {
	e.node.Props.Width = value
	e.node.Style.Layout.Width = Points(value)
	return e
}
func (e *element) Height(value float32) *element {
	e.node.Props.Height = value
	e.node.Style.Layout.Height = Points(value)
	return e
}
func (e *element) Align(value AxisAlignment) *element {
	e.node.Props.Alignment = value
	e.node.Style.Layout.Alignment = value
	return e
}
func (e *element) Flex(grow, shrink float32) *element {
	e.node.Style.Layout.FlexGrow = grow
	e.node.Style.Layout.FlexShrink = shrink
	return e
}
func (e *element) FlexBasis(value Length) *element {
	e.node.Style.Layout.FlexBasis = value
	return e
}
func (e *element) Wrap() *element { e.node.Style.Layout.Wrap = Wrap; return e }
func (e *element) AspectRatio(value float32) *element {
	e.node.Style.Layout.AspectRatio = value
	return e
}
func (e *element) GridColumns(count int) *element {
	if count < 1 {
		count = 1
	}
	e.node.Style.Layout.GridColumns = count
	return e
}
func (e *element) AdaptiveGrid(minColumnWidth float32) *element {
	e.node.Style.Layout.GridMinColumnWidth = minColumnWidth
	return e
}
func (e *element) FontSize(value float32) *element {
	e.node.Props.FontSize = value
	e.node.Style.Text.FontSize = value
	return e
}
func (e *element) Bold() *element {
	e.node.Props.Bold = true
	e.node.Style.Text.FontWeight = 700
	return e
}

// Styled applies the production typed style model. Fields supported by protocol
// v7 are projected into legacy Props; remaining fields stay on Node for layout,
// inspection, and the upcoming typed protocol.
func (e *element) Styled(style Style) *element {
	e.node.Style = e.node.Style.Merge(style)
	applyLegacyStyle(&e.node.Props, e.node.Style)
	return e
}

func (e *element) PlatformStyled(style PlatformStyle) *element { e.node.Platform = style; return e }
func (e *element) Background(color Color) *element {
	e.node.Style.Appearance.Background = color
	return e
}
func (e *element) Foreground(color Color) *element {
	e.node.Style.Appearance.Foreground = color
	return e
}
func (e *element) CornerRadius(value float32) *element {
	e.node.Style.Appearance.CornerRadius = value
	return e
}
func (e *element) Opacity(value float32) *element {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	e.node.Style.Appearance.Opacity = value
	return e
}
func (e *element) Disabled(value bool) *element { e.node.Style.Interaction.Disabled = value; return e }

func applyLegacyStyle(props *Props, style Style) {
	if style.Layout.Width.Unit == LengthPoints {
		props.Width = style.Layout.Width.Value
	}
	if style.Layout.Height.Unit == LengthPoints {
		props.Height = style.Layout.Height.Value
	}
	padding := style.Layout.Padding
	if padding.Top == padding.Leading && padding.Top == padding.Bottom && padding.Top == padding.Trailing {
		props.Padding = padding.Top
	}
	props.Gap = style.Layout.Gap
	props.Alignment = style.Layout.Alignment
	if style.Text.FontSize != 0 {
		props.FontSize = style.Text.FontSize
	}
	if style.Text.FontWeight >= 600 {
		props.Bold = true
	}
}
func (e *element) AccessibilityLabel(value string) *element {
	e.node.Props.AccessLabel = value
	return e
}

// AccessibilityHint describes the result of interacting with an element.
func (e *element) AccessibilityHint(value string) *element {
	e.node.Props.AccessHint = value
	return e
}

// AccessibilityRole overrides the native semantic role inferred from node type.
func (e *element) AccessibilityRole(value AccessibilityRole) *element {
	e.node.Props.AccessRole = value
	return e
}

// AccessibilityFocused requests native accessibility focus after this update.
func (e *element) AccessibilityFocused(value bool) *element {
	e.node.Props.Focused = value
	return e
}

// ScalesText opts text into the platform's user-selected text scaling behavior.
func (e *element) ScalesText() *element {
	e.node.Props.ScalesText = true
	return e
}

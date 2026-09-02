package ui

// View creates a generic native container.
func View(children ...Component) *element { return newElement(NodeView, Props{}, children...) }

// Column creates a vertical native container.
func Column(children ...Component) *element { return newElement(NodeColumn, Props{}, children...) }

// Row creates a horizontal native container.
func Row(children ...Component) *element { return newElement(NodeRow, Props{}, children...) }

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

func (e *element) Padding(value float32) *element     { e.node.Props.Padding = value; return e }
func (e *element) Gap(value float32) *element         { e.node.Props.Gap = value; return e }
func (e *element) Width(value float32) *element       { e.node.Props.Width = value; return e }
func (e *element) Height(value float32) *element      { e.node.Props.Height = value; return e }
func (e *element) Align(value AxisAlignment) *element { e.node.Props.Alignment = value; return e }
func (e *element) FontSize(value float32) *element    { e.node.Props.FontSize = value; return e }
func (e *element) Bold() *element                     { e.node.Props.Bold = true; return e }
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

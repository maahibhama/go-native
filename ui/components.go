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

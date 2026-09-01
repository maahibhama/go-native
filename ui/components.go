package ui

// View creates a generic native container.
func View(children ...Component) *element { return newElement(NodeView, Props{}, children...) }

// Column creates a vertical native container.
func Column(children ...Component) *element { return newElement(NodeColumn, Props{}, children...) }

// Row creates a horizontal native container.
func Row(children ...Component) *element { return newElement(NodeRow, Props{}, children...) }

// Text creates a native text label.
func Text(value string) *element { return newElement(NodeText, Props{Text: value}) }

// Button creates a native button. BindHandlers resolves its callback to a stable ID.
func Button(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label})
	buttonCallbacks.Store(e.node.ID, onPress)
	return e
}

var buttonCallbacks callbackStore

// TakeButtonCallback returns the callback associated with a button node.
// Runtime uses this during tree binding; applications normally do not call it.
func TakeButtonCallback(id NodeID) func() { return buttonCallbacks.Load(id) }

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

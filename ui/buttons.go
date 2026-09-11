package ui

func applyButtonTouchTargets(e *element) {
	e.node.Style.Layout.MinHeight = Points(44)
	e.node.Style.Layout.MinWidth = Points(44)
	e.node.Platform.IOS.Layout.MinHeight = Points(44)
	e.node.Platform.IOS.Layout.MinWidth = Points(44)
	e.node.Platform.Android.Layout.MinHeight = Points(48)
	e.node.Platform.Android.Layout.MinWidth = Points(48)
}

// FilledButton creates a high-emphasis filled button.
func FilledButton(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label, AccessRole: RoleButton})
	e.OnPress(onPress)
	e.Background(RGB(0, 122, 255))
	e.Foreground(RGB(255, 255, 255))
	e.CornerRadius(8)
	e.PaddingXY(16, 10)
	e.Align(AlignCenter)
	applyButtonTouchTargets(e)
	return e
}

// OutlinedButton creates a medium-emphasis button with an outline border.
func OutlinedButton(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label, AccessRole: RoleButton})
	e.OnPress(onPress)
	e.Background(RGBA(0, 0, 0, 0))
	e.Border(1, RGB(0, 122, 255))
	e.Foreground(RGB(0, 122, 255))
	e.CornerRadius(8)
	e.PaddingXY(16, 10)
	e.Align(AlignCenter)
	applyButtonTouchTargets(e)
	return e
}

// TextButton creates a low-emphasis button without background or border.
func TextButton(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label, AccessRole: RoleButton})
	e.OnPress(onPress)
	e.Background(RGBA(0, 0, 0, 0))
	e.Foreground(RGB(0, 122, 255))
	e.CornerRadius(8)
	e.PaddingXY(12, 8)
	e.Align(AlignCenter)
	applyButtonTouchTargets(e)
	return e
}

// ElevatedButton creates a button with elevation shadow for depth.
func ElevatedButton(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label, AccessRole: RoleButton})
	e.OnPress(onPress)
	e.Background(RGB(255, 255, 255))
	e.Foreground(RGB(0, 122, 255))
	e.Shadow(Shadow{Offset: Point{X: 0, Y: 2}, Blur: 4, Opacity: 0.2, Color: RGB(0, 0, 0)})
	e.CornerRadius(8)
	e.PaddingXY(16, 10)
	e.Align(AlignCenter)
	applyButtonTouchTargets(e)
	return e
}

// TonalButton creates a medium-emphasis tonal container button.
func TonalButton(label string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: label, AccessRole: RoleButton})
	e.OnPress(onPress)
	e.Background(RGBA(0, 122, 255, 38))
	e.Foreground(RGB(0, 122, 255))
	e.CornerRadius(8)
	e.PaddingXY(16, 10)
	e.Align(AlignCenter)
	applyButtonTouchTargets(e)
	return e
}

// IconButton creates an accessible icon button.
func IconButton(icon string, onPress func()) *element {
	e := newElement(NodeButton, Props{Text: icon, AccessRole: RoleButton})
	e.OnPress(onPress)
	e.Background(RGBA(0, 0, 0, 0))
	e.Foreground(RGB(0, 122, 255))
	e.CornerRadius(22)
	e.Padding(8)
	e.Align(AlignCenter)
	applyButtonTouchTargets(e)
	return e
}

// Leading attaches a leading component before the button label.
func (e *element) Leading(component Component) *element {
	e.leading = component
	return e
}

// Trailing attaches a trailing component after the button label.
func (e *element) Trailing(component Component) *element {
	e.trailing = component
	return e
}

// Loading configures whether the button is in a loading state.
func (e *element) Loading(loading bool) *element {
	e.loading = loading
	if loading {
		e.Disabled(true)
	} else {
		e.Disabled(false)
	}
	return e
}

// Compact controls whether the accessible touch target enforcement is bypassed.
func (e *element) Compact(compact bool) *element {
	e.compact = compact
	if compact {
		e.node.Style.Layout.MinHeight = Points(0)
		e.node.Style.Layout.MinWidth = Points(0)
		e.node.Platform.IOS.Layout.MinHeight = Points(0)
		e.node.Platform.IOS.Layout.MinWidth = Points(0)
		e.node.Platform.Android.Layout.MinHeight = Points(0)
		e.node.Platform.Android.Layout.MinWidth = Points(0)
	} else {
		applyButtonTouchTargets(e)
	}
	return e
}

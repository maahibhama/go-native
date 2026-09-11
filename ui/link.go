package ui

// Link creates an accessible hyperlink element with underline styling.
func Link(label string, url string) *element {
	linkColor := RGB(0, 122, 255)
	span := Span{
		Text:      label,
		Link:      url,
		Underline: true,
		Color:     linkColor,
		HasColor:  true,
	}
	e := RichText([]Span{span}, nil)
	e.node.Props.AccessRole = RoleButton
	e.node.Style.Text.Color = linkColor
	return e
}

// LinkAction creates an accessible button that presents as an underlined text link.
func LinkAction(label string, onPress func()) *element {
	linkColor := RGB(0, 122, 255)
	span := Span{
		Text:      label,
		Link:      "action://press",
		Underline: true,
		Color:     linkColor,
		HasColor:  true,
	}
	e := RichText([]Span{span}, func(string) {
		if onPress != nil {
			onPress()
		}
	})
	e.node.Press = onPress
	e.node.Props.AccessRole = RoleButton
	e.node.Style.Text.Color = linkColor
	if onPress != nil {
		e.node.Intents.Gestures = append(e.node.Intents.Gestures, GestureIntent{
			Kind: GestureTap,
			Handler: func(GestureEvent) {
				if !e.node.Style.Interaction.Disabled {
					onPress()
				}
			},
		})
	}
	return e
}

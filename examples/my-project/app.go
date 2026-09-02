// Package app contains the application's declarative native UI.
package app

import "github.com/go-native/go-native/ui"

// App builds the root UI component.
func App() ui.Component {
	return ui.SafeArea(
		ui.Column(
			ui.Text("Hello from Go Native").FontSize(28).Bold(),
			ui.Text("Edit app.go to get started."),
		).Padding(20).Gap(12).Align(ui.AlignCenter),
	)
}

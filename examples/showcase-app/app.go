// Package app contains the application's declarative native UI.
package app

import "github.com/go-native/go-native/ui"

// App builds the root UI component.
func App() ui.Component {
	return ui.SafeArea(
		ui.ScrollView(
			ui.Column(
				ui.Text("Android native view gallery").FontSize(26).Bold(),
				ui.View(ui.Text("View")).Height(54).FillColor(ui.RGB(220, 238, 255)).CornerRadius(10).Padding(12),
				ui.Row(ui.Text("Row A"), ui.Text("Row B")).Gap(16),
				ui.Column(ui.Text("Column A"), ui.Text("Column B")).Gap(4),
				ui.Button("Button", func() {}),
				ui.TextInput("Input value", func(string) {}).Placeholder("Placeholder"),
				ui.Switch(false, func(bool) {}),
				ui.ProgressIndicator(0.65).Width(300),
				ui.Image("app_logo").Width(64).Height(64),
				ui.View(ui.Center(ui.Text("Centered text"))).Height(52).FillColor(ui.RGB(245, 245, 245)),
				ui.AspectRatio(3, ui.View()).Width(240).FillColor(ui.RGB(255, 225, 190)),
				ui.Divider(ui.RGB(80, 80, 80)).Width(240),
				ui.KeyboardAvoidingView(ui.Text("Keyboard avoiding view")),
				ui.Stack(ui.Text("Stack base"), ui.Text("Stack overlay")),
			).Padding(20).Gap(14).Align(ui.AlignCenter),
		),
	).FillColor(ui.RGB(20, 250, 250))
}

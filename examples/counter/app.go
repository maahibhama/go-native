// Package counter is the Milestone 0 application written entirely in Go.
package counter

import (
	"fmt"
	"github.com/go-native/go-native/ui"
)

var count = ui.NewState(0)
var name = ui.NewState("")
var enabled = ui.NewState(true)

// App builds the counter's native UI tree.
func App() ui.Component {
	return ui.SafeArea(
		ui.Column(
			ui.TextFunc(func() string { return fmt.Sprintf("Count: %d", count.Get()) }).FontSize(28).Bold().ScalesText().AccessibilityLabel("Counter value").AccessibilityRole(ui.RoleHeader),
			ui.Button("Increment", func() { count.Update(func(v int) int { return v + 1 }) }).AccessibilityLabel("Increment counter").AccessibilityHint("Increases the counter by one"),
			ui.TextInput(name.Get(), func(value string) { name.Set(value) }).AccessibilityLabel("Name"),
			ui.TextFunc(func() string { return "Hello, " + name.Get() }).ScalesText(),
			ui.Switch(enabled.Get(), func(value bool) { enabled.Set(value) }).AccessibilityLabel("Enable progress"),
			ui.ProgressIndicator(float32(count.Get()%11)/10).AccessibilityLabel("Counter progress"),
			ui.ScrollView(ui.Row(ui.Text("Native"), ui.Text("Go"), ui.Text("UI")).Gap(16)).HorizontalScroll().Width(180),
		).Padding(20).Gap(12).Align(ui.AlignCenter),
	)
}

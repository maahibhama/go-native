// Package counter is the Milestone 0 application written entirely in Go.
package counter

import (
	"fmt"
	"github.com/go-native/go-native/ui"
)

var count = ui.NewState(0)

// App builds the counter's native UI tree.
func App() ui.Component {
	return ui.SafeArea(
		ui.Column(
			ui.TextFunc(func() string { return fmt.Sprintf("Count: %d", count.Get()) }).FontSize(28).Bold().AccessibilityLabel("Counter value"),
			ui.Button("Increment", func() { count.Update(func(v int) int { return v + 1 }) }).AccessibilityLabel("Increment counter"),
		).Padding(20).Gap(12).Align(ui.AlignCenter),
	)
}

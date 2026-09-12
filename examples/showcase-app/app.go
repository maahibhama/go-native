// Package app contains the application's declarative native UI.
package app

import "github.com/go-native/go-native/ui"

// App builds the root UI component.
func App() ui.Component {

	return ui.View().FillColor(ui.RGB(230, 245, 245))

}

// Package app contains the application's declarative native UI.
package app

import (
	"fmt"
	"strings"

	"github.com/go-native/go-native/ui"
)

// Global state persisting across reconciliation passes
var (
	isLoggedIn    = ui.NewState(false)
	username      = ui.NewState("")
	password      = ui.NewState("")
	rememberMe    = ui.NewState(true)
	statusMessage = ui.NewState("")
	loginCount    = ui.NewState(0)
)

// App builds the root UI component with navigation between Login and Dashboard.
func App() ui.Component {
	if isLoggedIn.Get() {
		return renderDashboardScreen()
	}
	return renderLoginScreen()
}

// renderLoginScreen builds the interactive Login View with input fields, image logo, and buttons.
func renderLoginScreen() ui.Component {
	return ui.SafeArea(
		ui.Column(
			// Brand Header Image & Titles
			ui.Image("app_logo").Width(64).Height(64).ResizeMode(ui.ImageFit),
			ui.Text("Welcome Back").FontSize(28).Bold(),
			ui.Text("Sign in to continue to Go Native").FontSize(15),

			// Username Input
			ui.Text("Username / Email").FontSize(14).Bold(),
			ui.TextInput(username.Get(), func(val string) {
				username.Set(val)
			}).AccessibilityHint("Enter your username or email").Width(280),

			// Password Input
			ui.Text("Password").FontSize(14).Bold(),
			ui.TextInput(password.Get(), func(val string) {
				password.Set(val)
			}).AccessibilityHint("Enter your password").Width(280),

			// Remember Me Switch
			ui.Row(
				ui.Switch(rememberMe.Get(), func(val bool) {
					rememberMe.Set(val)
				}),
				ui.Text("Remember me").FontSize(14),
			).Gap(12).Align(ui.AlignCenter),

			// Inline Feedback/Error Message
			ui.TextFunc(func() string {
				return statusMessage.Get()
			}).FontSize(14).Bold(),

			// Primary Action Button
			ui.Button("Sign In", func() {
				user := strings.TrimSpace(username.Get())
				pass := strings.TrimSpace(password.Get())

				if user == "" || pass == "" {
					statusMessage.Set("Please enter both username and password.")
					return
				}

				statusMessage.Set("")
				loginCount.Update(func(v int) int { return v + 1 })
				isLoggedIn.Set(true)
			}).Width(280).Height(46).FontSize(16).Bold().AccessibilityLabel("Sign In Button"),

			// Demo Helper Button
			ui.Button("Fill Demo Credentials", func() {
				username.Set("alex.developer")
				password.Set("gonative2026")
				statusMessage.Set("Demo credentials filled.")
			}).Width(280).Height(38).FontSize(14),
		).Padding(24).Gap(12).Align(ui.AlignCenter),
	).Align(ui.AlignCenter)
}

// renderDashboardScreen builds the authenticated User Dashboard View.
func renderDashboardScreen() ui.Component {
	return ui.SafeArea(
		ui.Column(
			ui.Image("avatar").Width(64).Height(64).ResizeMode(ui.ImageFit),
			ui.TextFunc(func() string {
				return fmt.Sprintf("Hello, %s!", username.Get())
			}).FontSize(26).Bold(),
			ui.Text("You are successfully signed in with native UI controls.").FontSize(15),

			ui.Column(
				ui.TextFunc(func() string {
					return fmt.Sprintf("Session Logins: %d", loginCount.Get())
				}).FontSize(16).Bold(),
				ui.Text("Platform Controls: Pure UIKit / Android Views").FontSize(14),
				ui.ProgressIndicator(1.0).Width(240),
			).Padding(16).Gap(8).Align(ui.AlignCenter),

			ui.Button("Sign Out", func() {
				isLoggedIn.Set(false)
				password.Set("")
				statusMessage.Set("Signed out successfully.")
			}).Width(220).Height(44).FontSize(15).Bold(),
		).Padding(24).Gap(18).Align(ui.AlignCenter),
	).Align(ui.AlignCenter)
}

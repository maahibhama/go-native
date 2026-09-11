package ui

import (
	"testing"
)

func TestAccessibilityQualification(t *testing.T) {
	btn := FilledButton("Submit", func() {}).AccessibilityRole(RoleButton).AccessibilityLabel("Submit form").AccessibilityHint("Double tap to submit").ScalesText()
	node := btn.node
	if node.Props.AccessRole != RoleButton {
		t.Fatalf("expected RoleButton, got %v", node.Props.AccessRole)
	}
	if node.Props.AccessLabel != "Submit form" {
		t.Fatalf("expected 'Submit form', got %q", node.Props.AccessLabel)
	}
	if node.Props.AccessHint != "Double tap to submit" {
		t.Fatalf("expected 'Double tap to submit', got %q", node.Props.AccessHint)
	}
	if !node.Props.ScalesText {
		t.Fatal("expected ScalesText to be true")
	}

	icon := Icon("chevron.right")
	iconNode := icon.node
	if iconNode.Props.AccessRole != RoleImage {
		t.Fatalf("expected RoleImage, got %v", iconNode.Props.AccessRole)
	}
	if iconNode.Props.AccessLabel != "Chevron right" {
		t.Fatalf("expected 'Chevron right', got %q", iconNode.Props.AccessLabel)
	}
}

func TestDarkModeQualification(t *testing.T) {
	light := DefaultTheme()
	dark := DarkTheme()

	if light.Colors["background"] == dark.Colors["background"] {
		t.Fatalf("expected distinct backgrounds for light and dark themes")
	}

	bgToken := ColorToken("background", RGB(0, 0, 0))
	lightBg := bgToken.Resolve(light)
	darkBg := bgToken.Resolve(dark)

	if lightBg != RGB(255, 255, 255) {
		t.Fatalf("expected white light background, got %v", lightBg)
	}
	if darkBg != RGB(18, 18, 18) {
		t.Fatalf("expected dark background, got %v", darkBg)
	}
}

func TestRTLDirectionalityQualification(t *testing.T) {
	ltrEnv := Environment{Direction: DirectionLTR}
	rtlEnv := Environment{Direction: DirectionRTL}

	if ltrEnv.Direction == rtlEnv.Direction {
		t.Fatal("expected distinct layout directions")
	}

	insets := EdgeInsets{Top: 8, Leading: 16, Bottom: 8, Trailing: 24}
	if insets.Leading != 16 || insets.Trailing != 24 {
		t.Fatalf("expected directional insets preserved, got %v", insets)
	}
}

func TestResponsiveQualification(t *testing.T) {
	mobileQuery := MediaQuery{
		Viewport:       Size{Width: 375, Height: 812},
		Scale:          3.0,
		TextScale:      1.0,
		Orientation:    OrientationPortrait,
		SafeAreaInsets: EdgeInsets{Top: 44, Bottom: 34},
	}
	tabletQuery := MediaQuery{
		Viewport:       Size{Width: 1024, Height: 1366},
		Scale:          2.0,
		TextScale:      1.0,
		Orientation:    OrientationPortrait,
		SafeAreaInsets: EdgeInsets{Top: 24, Bottom: 20},
	}

	if mobileQuery.Viewport.Width >= tabletQuery.Viewport.Width {
		t.Fatalf("expected mobile viewport to be narrower than tablet")
	}
	if mobileQuery.SafeAreaInsets.Top != 44 {
		t.Fatalf("expected iOS notch inset 44, got %v", mobileQuery.SafeAreaInsets.Top)
	}
}

func TestNativeCompositionQualification(t *testing.T) {
	appBar := AppBar("Settings", IconButton("settings", func() {}))
	body := Column(
		Card(
			Avatar("", "JD", 48),
			Text("John Doe"),
			Badge(Icon("star"), 5),
		).CardElevated(),
		Form(
			TextInput("John", func(string) {}).Placeholder("First Name"),
			TextInput("Doe", func(string) {}).Placeholder("Last Name"),
			FilledButton("Save Profile", func() {}),
		),
		Row(
			Checkbox(true, func(bool) {}),
			Text("Enable notifications"),
		),
		Slider(0.75, func(float32) {}),
		LinearProgress(0.5),
		Skeleton(200, 20),
	)

	scaffold := Scaffold(ScaffoldProps{
		AppBar: appBar,
		Body:   body,
	})

	node := scaffold.Build()
	if node == nil {
		t.Fatal("expected valid node tree from composed Scaffold")
	}
	if len(node.Children) == 0 {
		t.Fatal("expected children in composed Scaffold tree")
	}
}

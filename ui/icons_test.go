package ui

import (
	"testing"
)

func TestStandardIconConstants(t *testing.T) {
	expected := map[string]string{
		"chevron.right": IconChevronRight,
		"chevron.left":  IconChevronLeft,
		"chevron.up":    IconChevronUp,
		"chevron.down":  IconChevronDown,
		"check":         IconCheck,
		"close":         IconClose,
		"plus":          IconPlus,
		"minus":         IconMinus,
		"search":        IconSearch,
		"settings":      IconSettings,
		"star":          IconStar,
		"heart":         IconHeart,
		"share":         IconShare,
		"trash":         IconTrash,
		"refresh":       IconRefresh,
		"info":          IconInfo,
		"warning":       IconWarning,
		"error":         IconError,
	}

	for name, constantVal := range expected {
		if constantVal != name {
			t.Fatalf("expected constant %q, got %q", name, constantVal)
		}
		if !IsStandardIcon(constantVal) {
			t.Fatalf("expected %q to be recognized as standard icon", constantVal)
		}
	}

	if IsStandardIcon("nonexistent.icon") {
		t.Fatal("expected 'nonexistent.icon' to not be recognized as standard icon")
	}

	allIcons := StandardIcons()
	if len(allIcons) < len(expected) {
		t.Fatalf("expected at least %d standard icons, got %d", len(expected), len(allIcons))
	}
}

func TestIconCreationAndDefaultProps(t *testing.T) {
	node := Icon(IconSearch).Build()

	if node == nil {
		t.Fatal("expected non-nil node for Icon")
	}
	if node.Type != NodeImage {
		t.Fatalf("expected NodeImage type, got %v", node.Type)
	}
	if node.Props.ImageSource != "search" {
		t.Fatalf("expected image source 'search', got %q", node.Props.ImageSource)
	}
	if node.Props.ImageMode != ImageFit {
		t.Fatalf("expected image mode ImageFit, got %v", node.Props.ImageMode)
	}
	if node.Props.AccessRole != RoleImage {
		t.Fatalf("expected access role RoleImage, got %v", node.Props.AccessRole)
	}
	if node.Props.AccessLabel != "Search" {
		t.Fatalf("expected access label 'Search', got %q", node.Props.AccessLabel)
	}
	if node.Props.Width != 24 || node.Props.Height != 24 {
		t.Fatalf("expected default size 24x24, got %vx%v", node.Props.Width, node.Props.Height)
	}
	if node.Style.Layout.Width != Points(24) || node.Style.Layout.Height != Points(24) {
		t.Fatalf("expected layout size Points(24), got %v x %v", node.Style.Layout.Width, node.Style.Layout.Height)
	}
}

func TestIconModifiers(t *testing.T) {
	color := RGB(255, 64, 128)
	node := Icon(IconChevronRight).
		IconSize(36).
		IconColor(color).
		Build()

	if node.Props.Width != 36 || node.Props.Height != 36 {
		t.Fatalf("expected size 36x36, got %vx%v", node.Props.Width, node.Props.Height)
	}
	if node.Style.Layout.Width != Points(36) || node.Style.Layout.Height != Points(36) {
		t.Fatalf("expected layout size Points(36), got %v x %v", node.Style.Layout.Width, node.Style.Layout.Height)
	}
	if node.Style.Appearance.Foreground != color {
		t.Fatalf("expected foreground color %v, got %v", color, node.Style.Appearance.Foreground)
	}
	if node.Style.Text.Color != color {
		t.Fatalf("expected text color %v, got %v", color, node.Style.Text.Color)
	}

	// Negative size clamped to 0
	clampedNode := Icon(IconStar).IconSize(-10).Build()
	if clampedNode.Props.Width != 0 || clampedNode.Props.Height != 0 {
		t.Fatalf("expected clamped size 0x0, got %vx%v", clampedNode.Props.Width, clampedNode.Props.Height)
	}
}

func TestIconAccessibleLabels(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{IconChevronRight, "Chevron right"},
		{IconChevronLeft, "Chevron left"},
		{IconChevronUp, "Chevron up"},
		{IconChevronDown, "Chevron down"},
		{IconCheck, "Check"},
		{IconClose, "Close"},
		{IconPlus, "Plus"},
		{IconMinus, "Minus"},
		{IconSearch, "Search"},
		{IconSettings, "Settings"},
		{IconStar, "Star"},
		{IconHeart, "Heart"},
		{IconShare, "Share"},
		{IconTrash, "Trash"},
		{IconRefresh, "Refresh"},
		{IconInfo, "Info"},
		{IconWarning, "Warning"},
		{IconError, "Error"},
		{"custom.download_file", "Custom download file"},
		{"arrow-forward", "Arrow forward"},
		{"", ""},
	}

	for _, tc := range testCases {
		got := IconAccessibleLabel(tc.name)
		if got != tc.expected {
			t.Errorf("IconAccessibleLabel(%q) = %q; want %q", tc.name, got, tc.expected)
		}
	}

	// Accessibility label override
	customLabelNode := Icon(IconCheck).AccessibilityLabel("Task completed").Build()
	if customLabelNode.Props.AccessLabel != "Task completed" {
		t.Fatalf("expected overridden access label 'Task completed', got %q", customLabelNode.Props.AccessLabel)
	}
}

func TestPlatformIconResolution(t *testing.T) {
	// iOS symbols
	if got := PlatformIcon(IconCheck, PlatformIOS); got != "checkmark" {
		t.Fatalf("expected iOS checkmark, got %q", got)
	}
	if got := PlatformIcon(IconClose, PlatformIOS); got != "xmark" {
		t.Fatalf("expected iOS xmark, got %q", got)
	}
	if got := PlatformIcon(IconSearch, PlatformIOS); got != "magnifyingglass" {
		t.Fatalf("expected iOS magnifyingglass, got %q", got)
	}
	if got := PlatformIcon(IconSettings, PlatformIOS); got != "gearshape" {
		t.Fatalf("expected iOS gearshape, got %q", got)
	}

	// Android drawables
	if got := PlatformIcon(IconCheck, PlatformAndroid); got != "checkbox_on_background" {
		t.Fatalf("expected Android checkbox_on_background, got %q", got)
	}
	if got := PlatformIcon(IconClose, PlatformAndroid); got != "ic_menu_close_clear_cancel" {
		t.Fatalf("expected Android ic_menu_close_clear_cancel, got %q", got)
	}
	if got := PlatformIcon(IconSearch, PlatformAndroid); got != "ic_menu_search" {
		t.Fatalf("expected Android ic_menu_search, got %q", got)
	}
	if got := PlatformIcon(IconSettings, PlatformAndroid); got != "ic_menu_preferences" {
		t.Fatalf("expected Android ic_menu_preferences, got %q", got)
	}

	// Custom/unknown passthrough
	if got := PlatformIcon("custom_icon", PlatformIOS); got != "custom_icon" {
		t.Fatalf("expected passthrough for custom icon on iOS, got %q", got)
	}
	if got := PlatformIcon("custom_icon", PlatformAndroid); got != "custom_icon" {
		t.Fatalf("expected passthrough for custom icon on Android, got %q", got)
	}
}

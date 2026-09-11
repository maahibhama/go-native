package ui

import (
	"testing"
)

func TestMediaViewCreationAndDefaults(t *testing.T) {
	node := MediaView("https://example.com/photo.jpg").Build()

	if node == nil {
		t.Fatal("expected non-nil node for MediaView")
	}
	if node.Type != NodeImage {
		t.Fatalf("expected NodeImage type, got %v", node.Type)
	}
	if node.Props.ImageSource != "https://example.com/photo.jpg" {
		t.Fatalf("expected image source 'https://example.com/photo.jpg', got %q", node.Props.ImageSource)
	}
	if node.Props.ImageMode != ImageFit {
		t.Fatalf("expected default ImageFit mode, got %v", node.Props.ImageMode)
	}
	if node.Props.AccessRole != RoleImage {
		t.Fatalf("expected access role RoleImage, got %v", node.Props.AccessRole)
	}
	if node.CachePolicy != CacheDefault {
		t.Fatalf("expected default CacheDefault policy, got %v", node.CachePolicy)
	}
	if node.Placeholder != nil {
		t.Fatalf("expected nil placeholder by default, got %#v", node.Placeholder)
	}
}

func TestMediaViewModifiers(t *testing.T) {
	placeholderText := "Loading photo..."
	placeholder := Text(placeholderText)

	node := MediaView("https://example.com/banner.png").
		ResizeMode(ImageFill).
		CornerRadius(16).
		CachePolicy(CacheReload).
		Placeholder(placeholder).
		Build()

	if node.Props.ImageMode != ImageFill {
		t.Fatalf("expected ImageFill mode, got %v", node.Props.ImageMode)
	}
	if node.Style.Appearance.CornerRadius != 16 {
		t.Fatalf("expected corner radius 16, got %v", node.Style.Appearance.CornerRadius)
	}
	if node.CachePolicy != CacheReload {
		t.Fatalf("expected CacheReload policy, got %v", node.CachePolicy)
	}
	if node.Placeholder == nil {
		t.Fatal("expected non-nil placeholder")
	}

	placeholderNode := node.Placeholder.Build()
	if placeholderNode.Props.Text != placeholderText {
		t.Fatalf("expected placeholder text %q, got %q", placeholderText, placeholderNode.Props.Text)
	}

	// Test CacheReturnCacheDataElseLoad and ImageCenter
	nodeCenter := MediaView("icon.png").
		ResizeMode(ImageCenter).
		CachePolicy(CacheReturnCacheDataElseLoad).
		Build()

	if nodeCenter.Props.ImageMode != ImageCenter {
		t.Fatalf("expected ImageCenter mode, got %v", nodeCenter.Props.ImageMode)
	}
	if nodeCenter.CachePolicy != CacheReturnCacheDataElseLoad {
		t.Fatalf("expected CacheReturnCacheDataElseLoad policy, got %v", nodeCenter.CachePolicy)
	}
}

func TestCachePolicyString(t *testing.T) {
	testCases := []struct {
		policy   CachePolicy
		expected string
	}{
		{CacheDefault, "CacheDefault"},
		{CacheReload, "CacheReload"},
		{CacheReturnCacheDataElseLoad, "CacheReturnCacheDataElseLoad"},
		{CachePolicy(99), "CachePolicy(99)"},
	}

	for _, tc := range testCases {
		if got := tc.policy.String(); got != tc.expected {
			t.Errorf("CachePolicy(%d).String() = %q; want %q", tc.policy, got, tc.expected)
		}
	}
}

func TestCircleShape(t *testing.T) {
	node := Circle(64).Build()

	if node == nil {
		t.Fatal("expected non-nil node for Circle")
	}
	if node.Type != NodeView {
		t.Fatalf("expected NodeView type, got %v", node.Type)
	}
	if node.Props.Width != 64 || node.Props.Height != 64 {
		t.Fatalf("expected size 64x64, got %vx%v", node.Props.Width, node.Props.Height)
	}
	if node.Style.Layout.Width != Points(64) || node.Style.Layout.Height != Points(64) {
		t.Fatalf("expected layout size Points(64), got %v x %v", node.Style.Layout.Width, node.Style.Layout.Height)
	}
	if node.Style.Appearance.CornerRadius != 32 {
		t.Fatalf("expected corner radius 32 (half diameter), got %v", node.Style.Appearance.CornerRadius)
	}

	// Negative diameter clamped to 0
	clamped := Circle(-20).Build()
	if clamped.Props.Width != 0 || clamped.Props.Height != 0 || clamped.Style.Appearance.CornerRadius != 0 {
		t.Fatalf("expected clamped circle dimensions 0, got %vx%v radius=%v",
			clamped.Props.Width, clamped.Props.Height, clamped.Style.Appearance.CornerRadius)
	}
}

func TestRectShape(t *testing.T) {
	node := Rect(120, 80).Build()

	if node == nil {
		t.Fatal("expected non-nil node for Rect")
	}
	if node.Type != NodeView {
		t.Fatalf("expected NodeView type, got %v", node.Type)
	}
	if node.Props.Width != 120 || node.Props.Height != 80 {
		t.Fatalf("expected size 120x80, got %vx%v", node.Props.Width, node.Props.Height)
	}
	if node.Style.Layout.Width != Points(120) || node.Style.Layout.Height != Points(80) {
		t.Fatalf("expected layout size 120x80, got %v x %v", node.Style.Layout.Width, node.Style.Layout.Height)
	}

	// Negative dimensions clamped to 0
	clamped := Rect(-10, -5).Build()
	if clamped.Props.Width != 0 || clamped.Props.Height != 0 {
		t.Fatalf("expected clamped dimensions 0x0, got %vx%v", clamped.Props.Width, clamped.Props.Height)
	}
}

func TestCapsuleShape(t *testing.T) {
	// Horizontal capsule: height is smaller dimension
	hCapsule := Capsule(120, 40).Build()
	if hCapsule.Props.Width != 120 || hCapsule.Props.Height != 40 {
		t.Fatalf("expected size 120x40, got %vx%v", hCapsule.Props.Width, hCapsule.Props.Height)
	}
	if hCapsule.Style.Appearance.CornerRadius != 20 {
		t.Fatalf("expected corner radius 20 (half of height 40), got %v", hCapsule.Style.Appearance.CornerRadius)
	}

	// Vertical capsule: width is smaller dimension
	vCapsule := Capsule(50, 150).Build()
	if vCapsule.Props.Width != 50 || vCapsule.Props.Height != 150 {
		t.Fatalf("expected size 50x150, got %vx%v", vCapsule.Props.Width, vCapsule.Props.Height)
	}
	if vCapsule.Style.Appearance.CornerRadius != 25 {
		t.Fatalf("expected corner radius 25 (half of width 50), got %v", vCapsule.Style.Appearance.CornerRadius)
	}

	// Square capsule: equal dimensions
	sqCapsule := Capsule(60, 60).Build()
	if sqCapsule.Style.Appearance.CornerRadius != 30 {
		t.Fatalf("expected corner radius 30 (half of 60), got %v", sqCapsule.Style.Appearance.CornerRadius)
	}

	// Negative dimensions clamped
	clamped := Capsule(-10, -20).Build()
	if clamped.Props.Width != 0 || clamped.Props.Height != 0 || clamped.Style.Appearance.CornerRadius != 0 {
		t.Fatalf("expected clamped capsule dimensions 0, got %vx%v radius=%v",
			clamped.Props.Width, clamped.Props.Height, clamped.Style.Appearance.CornerRadius)
	}
}

func TestVectorShapeModifiers(t *testing.T) {
	fillColor := RGB(40, 120, 200)
	strokeColor := RGB(10, 30, 60)
	shadowColor := RGBA(0, 0, 0, 120)

	node := Circle(48).
		Fill(fillColor).
		Stroke(3, strokeColor).
		DropShadow(shadowColor, 0, 4, 8).
		Build()

	if node.Style.Appearance.Background != fillColor {
		t.Fatalf("expected fill background %v, got %v", fillColor, node.Style.Appearance.Background)
	}
	if node.Style.Appearance.Border.Width != 3 {
		t.Fatalf("expected stroke width 3, got %v", node.Style.Appearance.Border.Width)
	}
	if node.Style.Appearance.Border.Color != strokeColor {
		t.Fatalf("expected stroke color %v, got %v", strokeColor, node.Style.Appearance.Border.Color)
	}
	if node.Style.Appearance.Shadow.Color != shadowColor {
		t.Fatalf("expected shadow color %v, got %v", shadowColor, node.Style.Appearance.Shadow.Color)
	}
	if node.Style.Appearance.Shadow.Offset != (Point{X: 0, Y: 4}) {
		t.Fatalf("expected shadow offset (0, 4), got %v", node.Style.Appearance.Shadow.Offset)
	}
	if node.Style.Appearance.Shadow.Blur != 8 {
		t.Fatalf("expected shadow blur 8, got %v", node.Style.Appearance.Shadow.Blur)
	}
	if node.Style.Appearance.Shadow.Opacity != 1.0 {
		t.Fatalf("expected shadow opacity 1.0, got %v", node.Style.Appearance.Shadow.Opacity)
	}

	// Additional modifier variants: FillColor, StrokeBorder, StrokeWidth, StrokeColor, ShadowColor
	altNode := Rect(80, 40).
		FillColor(RGB(200, 50, 50)).
		StrokeBorder(2, RGB(100, 0, 0)).
		StrokeWidth(4).
		StrokeColor(RGB(80, 0, 0)).
		ShadowColor(RGB(20, 20, 20)).
		Build()

	if altNode.Style.Appearance.Background != RGB(200, 50, 50) {
		t.Fatalf("expected FillColor %v, got %v", RGB(200, 50, 50), altNode.Style.Appearance.Background)
	}
	if altNode.Style.Appearance.Border.Width != 4 {
		t.Fatalf("expected StrokeWidth 4, got %v", altNode.Style.Appearance.Border.Width)
	}
	if altNode.Style.Appearance.Border.Color != RGB(80, 0, 0) {
		t.Fatalf("expected StrokeColor %v, got %v", RGB(80, 0, 0), altNode.Style.Appearance.Border.Color)
	}
	if altNode.Style.Appearance.Shadow.Color != RGB(20, 20, 20) {
		t.Fatalf("expected ShadowColor %v, got %v", RGB(20, 20, 20), altNode.Style.Appearance.Shadow.Color)
	}
}

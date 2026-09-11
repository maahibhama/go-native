package ui

import "testing"

func TestLayoutConvenienceModifiersPopulateTypedStyle(t *testing.T) {
	node := View().MinWidth(10).MinHeight(11).MaxWidth(100).MaxHeight(101).
		MarginXY(3, 4).PaddingXY(5, 6).Absolute(EdgeInsets{Top: 1}).
		Overflow(OverflowHidden).Border(2, RGB(1, 2, 3)).CornerRadius(8).
		Shadow(Shadow{Blur: 4}).Transform(Transform{ScaleX: 2}).Visibility(Hidden).
		HitSlop(Insets(12)).Build()
	if node.Style.Layout.MinWidth != Points(10) || node.Style.Layout.MaxHeight != Points(101) || node.Style.Layout.Margin != InsetsXY(3, 4) || node.Style.Layout.Padding != InsetsXY(5, 6) {
		t.Fatalf("layout style = %#v", node.Style.Layout)
	}
	if node.Style.Layout.Position != PositionAbsolute || node.Style.Layout.Overflow != OverflowHidden || node.Style.Appearance.Border.Width != 2 || node.Style.Appearance.Visibility != Hidden || node.Style.Interaction.HitSlop != Insets(12) {
		t.Fatalf("typed style = %#v", node.Style)
	}
}

func TestStackKeepsSizingChildAndOverlaysFollowingChildren(t *testing.T) {
	node := BuildWithContext(Stack(Text("base"), Text("overlay"), Text("top")), NewBuildContext(DefaultEnvironment()))
	if len(node.Children) != 3 || node.Children[0].Style.Layout.Position != PositionRelative || node.Children[1].Style.Layout.Position != PositionAbsolute || node.Children[2].Style.Layout.Position != PositionAbsolute {
		t.Fatalf("stack children = %#v", node.Children)
	}
}

func TestKeyboardAvoidingViewReadsEnvironmentInset(t *testing.T) {
	environment := DefaultEnvironment()
	environment.MediaQuery.KeyboardInsets.Bottom = 216
	node := BuildWithContext(KeyboardAvoidingView(Text("form")), NewBuildContext(environment))
	if node.Style.Layout.Padding.Bottom != 216 || len(node.Children) != 1 {
		t.Fatalf("keyboard avoiding node = %#v", node)
	}
}

func TestComponentShellHelpers(t *testing.T) {
	if got := FixedSpacer(12).Build(); got.Style.Layout.Width != Points(12) || got.Style.Layout.Height != Points(12) {
		t.Fatalf("fixed spacer = %#v", got.Style.Layout)
	}
	if got := AspectRatio(16.0/9.0, Text("media")).Build(); got.Style.Layout.AspectRatio != 16.0/9.0 {
		t.Fatalf("aspect ratio = %#v", got.Style.Layout)
	}
	if got := Divider(RGB(1, 2, 3)).Build(); got.Style.Layout.Height != Points(1) || got.Style.Appearance.Background != RGB(1, 2, 3) {
		t.Fatalf("divider = %#v", got)
	}
}

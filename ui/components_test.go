package ui

import "testing"

func TestTextFuncEvaluatesOnBuildDescription(t *testing.T) {
	value := "first"
	build := func() *Node { return TextFunc(func() string { return value }).Build() }
	if got := build().Props.Text; got != "first" {
		t.Fatalf("text=%q", got)
	}
	value = "second"
	if got := build().Props.Text; got != "second" {
		t.Fatalf("text=%q", got)
	}
}

func TestSafeAreaNode(t *testing.T) {
	n := SafeArea(Text("inside")).Build()
	if n.Type != NodeSafeArea || len(n.Children) != 1 || n.Children[0].Type != NodeText {
		t.Fatalf("unexpected safe area: %#v", n)
	}
}

func TestTextInputNode(t *testing.T) {
	called := ""
	n := TextInput("initial", func(value string) { called = value }).Build()
	if n.Type != NodeTextInput || n.Props.Text != "initial" || n.Change == nil {
		t.Fatalf("unexpected text input: %#v", n)
	}
	n.Change("edited")
	if called != "edited" {
		t.Fatalf("change value=%q", called)
	}
}

func TestSwitchAndProgressNodes(t *testing.T) {
	got := false
	toggle := Switch(true, func(value bool) { got = value }).Build()
	if toggle.Type != NodeSwitch || !toggle.Props.Checked || toggle.Toggle == nil {
		t.Fatalf("unexpected switch: %#v", toggle)
	}
	toggle.Toggle(true)
	if !got {
		t.Fatal("toggle handler not invoked")
	}
	for _, tc := range []struct{ in, want float32 }{{-1, 0}, {.4, .4}, {2, 1}} {
		n := ProgressIndicator(tc.in).Build()
		if n.Type != NodeProgressIndicator || n.Props.Progress != tc.want {
			t.Fatalf("progress(%v)=%#v", tc.in, n)
		}
	}
}

func TestImageAndScrollViewNodes(t *testing.T) {
	image := Image("logo").ResizeMode(ImageFill).Width(80).Height(80).Build()
	if image.Type != NodeImage || image.Props.ImageSource != "logo" || image.Props.ImageMode != ImageFill {
		t.Fatalf("unexpected image: %#v", image)
	}
	scroll := ScrollView(Row(Text("one"), Text("two"))).HorizontalScroll().Build()
	if scroll.Type != NodeScrollView || !scroll.Props.Horizontal || len(scroll.Children) != 1 {
		t.Fatalf("unexpected scroll: %#v", scroll)
	}
}

func TestAccessibilityModifiers(t *testing.T) {
	n := Text("Title").AccessibilityLabel("Screen title").AccessibilityHint("Introduces the screen").AccessibilityRole(RoleHeader).AccessibilityFocused(true).ScalesText().Build()
	if n.Props.AccessLabel != "Screen title" || n.Props.AccessHint != "Introduces the screen" || n.Props.AccessRole != RoleHeader || !n.Props.Focused || !n.Props.ScalesText {
		t.Fatalf("unexpected accessibility props: %#v", n.Props)
	}
}

func TestTypedStyleProjectsProtocolV7Fields(t *testing.T) {
	n := Text("styled").Styled(Style{Layout: LayoutStyle{Width: Points(120), Height: Points(44), Padding: Insets(8), Gap: 6, Alignment: AlignCenter}, Text: TextStyle{FontSize: 18, FontWeight: 700}, Appearance: AppearanceStyle{Background: RGB(1, 2, 3), CornerRadius: 10}}).Build()
	if n.Props.Width != 120 || n.Props.Height != 44 || n.Props.Padding != 8 || n.Props.Gap != 6 || n.Props.Alignment != AlignCenter || n.Props.FontSize != 18 || !n.Props.Bold {
		t.Fatalf("legacy projection = %#v", n.Props)
	}
	if n.Style.Appearance.Background != RGB(1, 2, 3) || n.Style.Appearance.CornerRadius != 10 {
		t.Fatalf("typed style lost: %#v", n.Style)
	}
}

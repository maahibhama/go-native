package layout

import (
	"testing"

	"github.com/go-native/go-native/ui"
)

var measureText = MeasureFunc(func(node *ui.Node, _ Constraints) ui.Size {
	if node.Type == ui.NodeText {
		return ui.Size{Width: float32(len(node.Props.Text) * 8), Height: 20}
	}
	return ui.Size{Width: 44, Height: 44}
})

func TestColumnLayoutUsesPaddingGapAndCenterAlignment(t *testing.T) {
	root := ui.Column(ui.Text("one"), ui.Text("longer")).Width(200).Padding(10).Gap(6).Align(ui.AlignCenter).Build()
	box := (Engine{Measurer: measureText}).Layout(root, Constraints{MaxWidth: 300, MaxHeight: 600})
	if box.Frame.Width != 200 || box.Frame.Height != 66 {
		t.Fatalf("root frame = %#v", box.Frame)
	}
	if box.Children[0].Frame.X != 88 || box.Children[1].Frame.X != 76 {
		t.Fatalf("children = %#v", box.Children)
	}
	if box.Children[1].Frame.Y != 36 {
		t.Fatalf("second y = %v", box.Children[1].Frame.Y)
	}
}

func TestRowPercentAndAbsoluteLayout(t *testing.T) {
	first := ui.View().Styled(ui.Style{Layout: ui.LayoutStyle{Width: ui.Percent(50), Height: ui.Points(30)}})
	second := ui.View().Width(20).Height(30)
	overlay := ui.View().Width(10).Height(10).Styled(ui.Style{Layout: ui.LayoutStyle{Position: ui.PositionAbsolute, Inset: ui.EdgeInsets{Top: 7, Leading: 9}}})
	root := ui.Row(first, second, overlay).Width(200).Height(50).Build()
	box := (Engine{Measurer: measureText}).Layout(root, Constraints{MaxWidth: 200, MaxHeight: 50})
	if box.Children[0].Frame.Width != 100 || box.Children[1].Frame.X != 100 {
		t.Fatalf("flow boxes = %#v", box.Children)
	}
	if box.Children[2].Frame.X != 9 || box.Children[2].Frame.Y != 7 {
		t.Fatalf("overlay = %#v", box.Children[2].Frame)
	}
}

func TestMinMaxConstraintsClampLeaf(t *testing.T) {
	leaf := ui.Text("wide").Styled(ui.Style{Layout: ui.LayoutStyle{MinWidth: ui.Points(50), MaxWidth: ui.Points(60)}}).Build()
	box := (Engine{Measurer: measureText}).Layout(leaf, Constraints{MaxWidth: 200, MaxHeight: 100})
	if box.Frame.Width != 50 {
		t.Fatalf("width = %v", box.Frame.Width)
	}
}

func TestFlexGrowAndShrink(t *testing.T) {
	grow := ui.Row(
		ui.View().Height(20).Flex(1, 1).FlexBasis(ui.Points(50)),
		ui.View().Height(20).Flex(2, 1).FlexBasis(ui.Points(50)),
	).Width(200).Build()
	box := (Engine{Measurer: measureText}).Layout(grow, Constraints{MaxWidth: 200, MaxHeight: 100})
	if !near(box.Children[0].Frame.Width, 83.333) || !near(box.Children[1].Frame.Width, 116.667) {
		t.Fatalf("grow widths = %v, %v", box.Children[0].Frame.Width, box.Children[1].Frame.Width)
	}

	shrink := ui.Row(
		ui.View().Width(100).Height(20).Flex(0, 1),
		ui.View().Width(100).Height(20).Flex(0, 1),
	).Width(150).Build()
	box = (Engine{Measurer: measureText}).Layout(shrink, Constraints{MaxWidth: 150, MaxHeight: 100})
	if box.Children[0].Frame.Width != 75 || box.Children[1].Frame.Width != 75 {
		t.Fatalf("shrink widths = %v, %v", box.Children[0].Frame.Width, box.Children[1].Frame.Width)
	}
}

func TestWrapAspectRatioAndGrid(t *testing.T) {
	wrapped := ui.Row(
		ui.View().Width(60).Height(20), ui.View().Width(60).Height(20), ui.View().Width(60).Height(20),
	).Width(130).Gap(10).Wrap().Build()
	box := (Engine{Measurer: measureText}).Layout(wrapped, Constraints{MaxWidth: 130, MaxHeight: 200})
	if box.Frame.Height != 50 || box.Children[2].Frame.Y != 30 {
		t.Fatalf("wrapped layout = %#v", box)
	}

	ratio := ui.View().Width(120).AspectRatio(1.5).Build()
	box = (Engine{Measurer: measureText}).Layout(ratio, Constraints{MaxWidth: 200, MaxHeight: 200})
	if box.Frame.Height != 80 {
		t.Fatalf("aspect frame = %#v", box.Frame)
	}

	grid := ui.Grid(2, ui.View().Height(20), ui.View().Height(20), ui.View().Height(20)).Width(220).Gap(10).Build()
	box = (Engine{Measurer: measureText}).Layout(grid, Constraints{MaxWidth: 220, MaxHeight: 200})
	if box.Children[0].Frame.Width != 105 || box.Children[1].Frame.X != 115 || box.Children[2].Frame.Y != 30 {
		t.Fatalf("grid layout = %#v", box.Children)
	}
}

func TestAdaptiveGridAndRTL(t *testing.T) {
	grid := ui.Grid(1, ui.View().Height(20), ui.View().Height(20), ui.View().Height(20)).AdaptiveGrid(90).Width(300).Gap(10).Build()
	box := (Engine{Direction: ui.DirectionRTL, Measurer: measureText}).Layout(grid, Constraints{MaxWidth: 300, MaxHeight: 200})
	if box.Children[0].Frame.X != 206.66667 || box.Children[2].Frame.X != 0 {
		t.Fatalf("RTL adaptive columns = %#v", box.Children)
	}

	row := ui.Row(ui.View().Width(20).Height(20), ui.View().Width(30).Height(20)).Width(200).Build()
	box = (Engine{Direction: ui.DirectionRTL, Measurer: measureText}).Layout(row, Constraints{MaxWidth: 200, MaxHeight: 100})
	if box.Children[0].Frame.X != 180 || box.Children[1].Frame.X != 150 {
		t.Fatalf("RTL row = %#v", box.Children)
	}
}

func near(got, want float32) bool {
	if got > want {
		return got-want < .01
	}
	return want-got < .01
}

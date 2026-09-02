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

package layout

import (
	"context"
	"testing"

	"github.com/go-native/go-native/ui"
)

func TestPipelineProducesDeterministicParentFirstGeometry(t *testing.T) {
	root := ui.WithID(ui.Column(ui.WithID(ui.Text("first"), 2), ui.WithID(ui.Text("second"), 3)), 1).Build()
	pipeline := &Pipeline{Engine: Engine{Measurer: MeasureFunc(func(_ *ui.Node, _ Constraints) ui.Size { return ui.Size{Width: 20, Height: 10} })}}
	geometry, err := pipeline.Compute(context.Background(), root, Constraints{MaxWidth: 100, MaxHeight: 100})
	if err != nil {
		t.Fatal(err)
	}
	if len(geometry.Frames) != 3 || geometry.Frames[0].NodeID != 1 || geometry.Frames[1].NodeID != 2 || geometry.Frames[2].NodeID != 3 {
		t.Fatalf("unexpected geometry: %#v", geometry.Frames)
	}
	if rect, ok := geometry.RectFor(3); !ok || rect.Y != 10 {
		t.Fatalf("second rect = %#v, %v", rect, ok)
	}
}

func TestComputeLayoutUsesTightNativeViewport(t *testing.T) {
	root := ui.WithID(ui.SafeArea(ui.WithID(ui.View().Width(40).Height(20), 2)).Align(ui.AlignCenter), 1).Build()
	pipeline := &Pipeline{}
	frames, err := pipeline.ComputeLayout(root, ui.Environment{MediaQuery: ui.MediaQuery{Viewport: ui.Size{Width: 320, Height: 640}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 2 || frames[0].Rect.Width != 320 || frames[0].Rect.Height != 640 {
		t.Fatalf("root viewport frame = %#v", frames)
	}
	if frames[1].Rect.X != 140 || frames[1].Rect.Y != 0 {
		t.Fatalf("centered child frame = %#v", frames[1].Rect)
	}
}

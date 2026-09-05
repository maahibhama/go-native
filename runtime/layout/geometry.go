package layout

import (
	"context"
	"fmt"

	"github.com/go-native/go-native/ui"
)

// Geometry is the deterministic, flattened output of Go-owned layout. It uses
// values only and is safe to encode for a native host.
type Geometry struct {
	Frames []ui.LayoutFrame
}

// Frame is retained as a source-compatible alias for portable layout output.
type Frame = ui.LayoutFrame

// Flatten returns parent-before-child geometry in deterministic tree order.
func Flatten(root *Box) Geometry {
	geometry := Geometry{}
	var visit func(*Box)
	visit = func(box *Box) {
		if box == nil {
			return
		}
		geometry.Frames = append(geometry.Frames, Frame{NodeID: box.NodeID, Rect: ui.LayoutRect{X: box.Frame.X, Y: box.Frame.Y, Width: box.Frame.Width, Height: box.Frame.Height}})
		for _, child := range box.Children {
			visit(child)
		}
	}
	visit(root)
	return geometry
}

// RectFor finds a computed frame by node identity.
func (g Geometry) RectFor(id ui.NodeID) (ui.LayoutRect, bool) {
	for _, frame := range g.Frames {
		if frame.NodeID == id {
			return frame.Rect, true
		}
	}
	return ui.LayoutRect{}, false
}

// Pipeline coordinates intrinsic native measurement with deterministic layout.
// It is host-independent: renderers provide BatchMeasurer implementations.
type Pipeline struct {
	Engine   Engine
	Measurer BatchMeasurer
	Cache    *MeasurementCache
}

func (p *Pipeline) Compute(ctx context.Context, root *ui.Node, constraints Constraints) (Geometry, error) {
	if root == nil {
		return Geometry{}, nil
	}
	if p == nil {
		return Geometry{}, fmt.Errorf("layout: nil pipeline")
	}
	var box *Box
	var err error
	if p.Measurer == nil {
		box = p.Engine.Layout(root, constraints)
	} else {
		box, err = p.Engine.LayoutMeasured(ctx, root, constraints, p.Measurer, p.Cache)
	}
	if err != nil {
		return Geometry{}, err
	}
	return Flatten(box), nil
}

func (p *Pipeline) InvalidateMeasurements() {
	if p != nil && p.Cache != nil {
		p.Cache.Clear()
	}
}

// ComputeLayout lets Pipeline be installed directly on runtime.Runtime as its
// optional LayoutProvider without coupling this package back to runtime.
func (p *Pipeline) ComputeLayout(root *ui.Node, environment ui.Environment) ([]ui.LayoutFrame, error) {
	width, height := environment.MediaQuery.Viewport.Width, environment.MediaQuery.Viewport.Height
	if width <= 0 || height <= 0 {
		return nil, nil
	}
	// Native hosts mount the root into the complete reported viewport. Tight
	// root constraints keep descendant alignment relative to that same host
	// instead of shrinking the abstract root around its content.
	geometry, err := p.Compute(context.Background(), root, Constraints{MinWidth: width, MaxWidth: width, MinHeight: height, MaxHeight: height})
	if err != nil {
		return nil, err
	}
	return geometry.Frames, nil
}

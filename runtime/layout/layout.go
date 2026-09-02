// Package layout implements deterministic Go-owned box measurement and placement.
package layout

import (
	"math"

	"github.com/go-native/go-native/ui"
)

type Constraints struct{ MinWidth, MaxWidth, MinHeight, MaxHeight float32 }
type Rect struct{ X, Y, Width, Height float32 }
type Box struct {
	NodeID   ui.NodeID
	Frame    Rect
	Children []*Box
}

type Measurer interface {
	Measure(*ui.Node, Constraints) ui.Size
}
type MeasureFunc func(*ui.Node, Constraints) ui.Size

func (f MeasureFunc) Measure(node *ui.Node, constraints Constraints) ui.Size {
	return f(node, constraints)
}

type Engine struct{ Measurer Measurer }

func (e Engine) Layout(root *ui.Node, constraints Constraints) *Box {
	if root == nil {
		return nil
	}
	return e.layout(root, normalize(constraints), constraints.MaxWidth, constraints.MaxHeight)
}

func (e Engine) layout(node *ui.Node, constraints Constraints, parentWidth, parentHeight float32) *Box {
	style := node.Style.Layout
	width, widthSet := resolve(style.Width, parentWidth)
	height, heightSet := resolve(style.Height, parentHeight)
	padding := style.Padding
	contentMaxWidth := maxZero(constraints.MaxWidth - padding.Leading - padding.Trailing)
	contentMaxHeight := maxZero(constraints.MaxHeight - padding.Top - padding.Bottom)

	box := &Box{NodeID: node.ID}
	if len(node.Children) == 0 {
		measured := ui.Size{}
		if e.Measurer != nil {
			measured = e.Measurer.Measure(node, Constraints{MaxWidth: contentMaxWidth, MaxHeight: contentMaxHeight})
		}
		if !widthSet {
			width = measured.Width + padding.Leading + padding.Trailing
		}
		if !heightSet {
			height = measured.Height + padding.Top + padding.Bottom
		}
		box.Frame.Width = bound(width, style.MinWidth, style.MaxWidth, constraints.MinWidth, constraints.MaxWidth, parentWidth)
		box.Frame.Height = bound(height, style.MinHeight, style.MaxHeight, constraints.MinHeight, constraints.MaxHeight, parentHeight)
		return box
	}

	main, cross := float32(0), float32(0)
	flowCount := 0
	for _, child := range node.Children {
		childBox := e.layout(child, Constraints{MaxWidth: contentMaxWidth, MaxHeight: contentMaxHeight}, contentMaxWidth, contentMaxHeight)
		box.Children = append(box.Children, childBox)
		if child.Style.Layout.Position == ui.PositionAbsolute {
			continue
		}
		if flowCount > 0 {
			main += style.Gap
		}
		margin := child.Style.Layout.Margin
		if style.Direction == ui.FlexRow {
			main += margin.Leading + childBox.Frame.Width + margin.Trailing
			cross = max(cross, margin.Top+childBox.Frame.Height+margin.Bottom)
		} else {
			main += margin.Top + childBox.Frame.Height + margin.Bottom
			cross = max(cross, margin.Leading+childBox.Frame.Width+margin.Trailing)
		}
		flowCount++
	}
	if !widthSet {
		if style.Direction == ui.FlexRow {
			width = main + padding.Leading + padding.Trailing
		} else {
			width = cross + padding.Leading + padding.Trailing
		}
	}
	if !heightSet {
		if style.Direction == ui.FlexRow {
			height = cross + padding.Top + padding.Bottom
		} else {
			height = main + padding.Top + padding.Bottom
		}
	}
	box.Frame.Width = bound(width, style.MinWidth, style.MaxWidth, constraints.MinWidth, constraints.MaxWidth, parentWidth)
	box.Frame.Height = bound(height, style.MinHeight, style.MaxHeight, constraints.MinHeight, constraints.MaxHeight, parentHeight)
	e.position(node, box)
	return box
}

func (e Engine) position(node *ui.Node, box *Box) {
	style := node.Style.Layout
	x, y := style.Padding.Leading, style.Padding.Top
	contentWidth := maxZero(box.Frame.Width - style.Padding.Leading - style.Padding.Trailing)
	contentHeight := maxZero(box.Frame.Height - style.Padding.Top - style.Padding.Bottom)
	flowIndex := 0
	for i, child := range node.Children {
		childBox := box.Children[i]
		childStyle := child.Style.Layout
		margin := childStyle.Margin
		if childStyle.Position == ui.PositionAbsolute {
			childBox.Frame.X = childStyle.Inset.Leading
			childBox.Frame.Y = childStyle.Inset.Top
			continue
		}
		if flowIndex > 0 {
			if style.Direction == ui.FlexRow {
				x += style.Gap
			} else {
				y += style.Gap
			}
		}
		if style.Direction == ui.FlexRow {
			x += margin.Leading
			childBox.Frame.X = x
			childBox.Frame.Y = crossPosition(style.Alignment, contentHeight, childBox.Frame.Height, margin.Top, margin.Bottom, style.Padding.Top)
			x += childBox.Frame.Width + margin.Trailing
		} else {
			y += margin.Top
			childBox.Frame.X = crossPosition(style.Alignment, contentWidth, childBox.Frame.Width, margin.Leading, margin.Trailing, style.Padding.Leading)
			childBox.Frame.Y = y
			y += childBox.Frame.Height + margin.Bottom
		}
		flowIndex++
	}
}

func crossPosition(alignment ui.AxisAlignment, available, size, leading, trailing, origin float32) float32 {
	switch alignment {
	case ui.AlignCenter:
		return origin + leading + maxZero(available-leading-trailing-size)/2
	case ui.AlignEnd:
		return origin + maxZero(available-trailing-size)
	default:
		return origin + leading
	}
}

func resolve(length ui.Length, parent float32) (float32, bool) {
	switch length.Unit {
	case ui.LengthPoints:
		return maxZero(length.Value), true
	case ui.LengthPercent:
		return maxZero(parent * length.Value / 100), true
	default:
		return 0, false
	}
}
func bound(value float32, minLength, maxLength ui.Length, constraintMin, constraintMax, parent float32) float32 {
	minValue, minSet := resolve(minLength, parent)
	if !minSet {
		minValue = constraintMin
	}
	maxValue, maxSet := resolve(maxLength, parent)
	if !maxSet {
		maxValue = constraintMax
	}
	if maxValue > 0 && value > maxValue {
		value = maxValue
	}
	if value < minValue {
		value = minValue
	}
	return maxZero(value)
}
func normalize(c Constraints) Constraints {
	if c.MaxWidth <= 0 {
		c.MaxWidth = float32(math.MaxInt16)
	}
	if c.MaxHeight <= 0 {
		c.MaxHeight = float32(math.MaxInt16)
	}
	return c
}
func maxZero(v float32) float32 {
	if v < 0 {
		return 0
	}
	return v
}
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

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

func (f MeasureFunc) Measure(n *ui.Node, c Constraints) ui.Size { return f(n, c) }

type Engine struct {
	Measurer  Measurer
	Direction ui.LayoutDirection
}

func (e Engine) Layout(root *ui.Node, c Constraints) *Box {
	if root == nil {
		return nil
	}
	c = normalize(c)
	return e.layout(root, c, c.MaxWidth, c.MaxHeight)
}

func (e Engine) layout(n *ui.Node, c Constraints, parentW, parentH float32) *Box {
	s := n.Style.Layout
	w, hasW := resolve(s.Width, parentW)
	h, hasH := resolve(s.Height, parentH)
	if s.AspectRatio > 0 {
		if hasW && !hasH {
			h, hasH = w/s.AspectRatio, true
		}
		if hasH && !hasW {
			w, hasW = h*s.AspectRatio, true
		}
	}
	maxW := maxZero(c.MaxWidth - s.Padding.Leading - s.Padding.Trailing)
	maxH := maxZero(c.MaxHeight - s.Padding.Top - s.Padding.Bottom)
	if hasW {
		maxW = maxZero(w - s.Padding.Leading - s.Padding.Trailing)
	}
	if hasH {
		maxH = maxZero(h - s.Padding.Top - s.Padding.Bottom)
	}
	b := &Box{NodeID: n.ID}
	if len(n.Children) == 0 {
		measured := ui.Size{}
		if e.Measurer != nil {
			measured = e.Measurer.Measure(n, Constraints{MaxWidth: maxW, MaxHeight: maxH})
		}
		if !hasW {
			w = measured.Width + s.Padding.Leading + s.Padding.Trailing
		}
		if !hasH {
			h = measured.Height + s.Padding.Top + s.Padding.Bottom
		}
		b.Frame.Width = bound(w, s.MinWidth, s.MaxWidth, c.MinWidth, c.MaxWidth, parentW)
		b.Frame.Height = bound(h, s.MinHeight, s.MaxHeight, c.MinHeight, c.MaxHeight, parentH)
		return b
	}
	for _, child := range n.Children {
		b.Children = append(b.Children, e.layout(child, Constraints{MaxWidth: maxW, MaxHeight: maxH}, maxW, maxH))
	}
	var contentW, contentH float32
	if s.GridColumns > 0 || s.GridMinColumnWidth > 0 {
		contentW, contentH = e.sizeGrid(n, b, maxW)
	} else {
		contentW, contentH = e.sizeFlex(n, b, maxW, maxH)
	}
	if !hasW {
		w = contentW + s.Padding.Leading + s.Padding.Trailing
	}
	if !hasH {
		h = contentH + s.Padding.Top + s.Padding.Bottom
	}
	b.Frame.Width = bound(w, s.MinWidth, s.MaxWidth, c.MinWidth, c.MaxWidth, parentW)
	b.Frame.Height = bound(h, s.MinHeight, s.MaxHeight, c.MinHeight, c.MaxHeight, parentH)
	if s.GridColumns > 0 || s.GridMinColumnWidth > 0 {
		e.placeGrid(n, b)
	} else {
		e.placeFlex(n, b)
	}
	e.placeAbsolute(n, b)
	return b
}

type flexItem struct {
	index                                int
	main, cross, mainMargin, crossMargin float32
}
type flexLine struct {
	items       []flexItem
	main, cross float32
}

func lines(n *ui.Node, b *Box, limit float32) []flexLine {
	s := n.Style.Layout
	out := []flexLine{{}}
	for i, child := range n.Children {
		if child.Style.Layout.Position == ui.PositionAbsolute {
			continue
		}
		cb, m := b.Children[i], child.Style.Layout.Margin
		main, cross, mm, cm := cb.Frame.Height, cb.Frame.Width, m.Top+m.Bottom, m.Leading+m.Trailing
		if s.Direction == ui.FlexRow {
			main, cross, mm, cm = cb.Frame.Width, cb.Frame.Height, m.Leading+m.Trailing, m.Top+m.Bottom
		}
		if basis, ok := resolve(child.Style.Layout.FlexBasis, limit); ok {
			main = basis
			if s.Direction == ui.FlexRow {
				cb.Frame.Width = basis
			} else {
				cb.Frame.Height = basis
			}
		}
		item := flexItem{i, main, cross, mm, cm}
		line := &out[len(out)-1]
		prospective := line.main + main + mm
		if len(line.items) > 0 {
			prospective += s.Gap
		}
		if s.Wrap == ui.Wrap && len(line.items) > 0 && prospective > limit {
			out = append(out, flexLine{})
			line = &out[len(out)-1]
		}
		if len(line.items) > 0 {
			line.main += s.Gap
		}
		line.items = append(line.items, item)
		line.main += main + mm
		line.cross = max(line.cross, cross+cm)
	}
	return out
}

func distribute(n *ui.Node, b *Box, line *flexLine, limit float32) {
	free := limit - line.main
	var total float32
	for _, item := range line.items {
		s := n.Children[item.index].Style.Layout
		if free >= 0 {
			total += maxZero(s.FlexGrow)
		} else {
			total += maxZero(s.FlexShrink) * item.main
		}
	}
	if total == 0 {
		return
	}
	line.main = limit
	for i := range line.items {
		item := &line.items[i]
		s := n.Children[item.index].Style.Layout
		factor := maxZero(s.FlexGrow)
		if free < 0 {
			factor = maxZero(s.FlexShrink) * item.main
		}
		item.main = maxZero(item.main + free*factor/total)
		if n.Style.Layout.Direction == ui.FlexRow {
			b.Children[item.index].Frame.Width = item.main
		} else {
			b.Children[item.index].Frame.Height = item.main
		}
	}
}

func (e Engine) sizeFlex(n *ui.Node, b *Box, maxW, maxH float32) (float32, float32) {
	limit := maxH
	if n.Style.Layout.Direction == ui.FlexRow {
		limit = maxW
	}
	ls := lines(n, b, limit)
	main, cross := float32(0), float32(0)
	for i := range ls {
		distribute(n, b, &ls[i], limit)
		if n.Style.Layout.Wrap == ui.Wrap {
			main = max(main, ls[i].main)
			if i > 0 {
				cross += n.Style.Layout.Gap
			}
			cross += ls[i].cross
		} else {
			main, cross = ls[i].main, ls[i].cross
		}
	}
	if n.Style.Layout.Direction == ui.FlexRow {
		return main, cross
	}
	return cross, main
}

func (e Engine) placeFlex(n *ui.Node, b *Box) {
	s := n.Style.Layout
	w := maxZero(b.Frame.Width - s.Padding.Leading - s.Padding.Trailing)
	h := maxZero(b.Frame.Height - s.Padding.Top - s.Padding.Bottom)
	limit := h
	if s.Direction == ui.FlexRow {
		limit = w
	}
	ls := lines(n, b, limit)
	crossCursor := float32(0)
	for li := range ls {
		distribute(n, b, &ls[li], limit)
		cursor := float32(0)
		lineCross := ls[li].cross
		if s.Wrap == ui.NoWrap {
			lineCross = w
			if s.Direction == ui.FlexRow {
				lineCross = h
			}
		}
		for ii, item := range ls[li].items {
			child, cb := n.Children[item.index], b.Children[item.index]
			m := child.Style.Layout.Margin
			if ii > 0 {
				cursor += s.Gap
			}
			if s.Direction == ui.FlexRow {
				cursor += m.Leading
				cb.Frame.X = s.Padding.Leading + cursor
				cb.Frame.Y = s.Padding.Top + crossCursor + crossOffset(s.Alignment, lineCross, cb.Frame.Height, m.Top, m.Bottom)
				cursor += cb.Frame.Width + m.Trailing
			} else {
				cursor += m.Top
				cb.Frame.X = s.Padding.Leading + crossOffset(s.Alignment, lineCross, cb.Frame.Width, m.Leading, m.Trailing)
				cb.Frame.Y = s.Padding.Top + cursor
				cursor += cb.Frame.Height + m.Bottom
			}
		}
		crossCursor += ls[li].cross
		if li < len(ls)-1 {
			crossCursor += s.Gap
		}
	}
	if e.Direction == ui.DirectionRTL {
		for i, child := range n.Children {
			if child.Style.Layout.Position != ui.PositionAbsolute {
				cb := b.Children[i]
				cb.Frame.X = b.Frame.Width - cb.Frame.X - cb.Frame.Width
			}
		}
	}
}

func (e Engine) sizeGrid(n *ui.Node, b *Box, width float32) (float32, float32) {
	s := n.Style.Layout
	cols := gridColumns(s, width)
	cw := maxZero((width - float32(cols-1)*s.Gap) / float32(cols))
	rows := (flowCount(n) + cols - 1) / cols
	heights := make([]float32, rows)
	flow := 0
	for i, ch := range n.Children {
		if ch.Style.Layout.Position == ui.PositionAbsolute {
			continue
		}
		m := ch.Style.Layout.Margin
		b.Children[i].Frame.Width = maxZero(cw - m.Leading - m.Trailing)
		row := flow / cols
		heights[row] = max(heights[row], m.Top+b.Children[i].Frame.Height+m.Bottom)
		flow++
	}
	height := float32(0)
	for i, v := range heights {
		if i > 0 {
			height += s.Gap
		}
		height += v
	}
	return width, height
}
func (e Engine) placeGrid(n *ui.Node, b *Box) {
	s := n.Style.Layout
	width := maxZero(b.Frame.Width - s.Padding.Leading - s.Padding.Trailing)
	cols := gridColumns(s, width)
	cw := maxZero((width - float32(cols-1)*s.Gap) / float32(cols))
	rows := (flowCount(n) + cols - 1) / cols
	heights := make([]float32, rows)
	flow := 0
	for i, ch := range n.Children {
		if ch.Style.Layout.Position == ui.PositionAbsolute {
			continue
		}
		m := ch.Style.Layout.Margin
		heights[flow/cols] = max(heights[flow/cols], m.Top+b.Children[i].Frame.Height+m.Bottom)
		flow++
	}
	y := s.Padding.Top
	flow = 0
	for i, ch := range n.Children {
		if ch.Style.Layout.Position == ui.PositionAbsolute {
			continue
		}
		row, col := flow/cols, flow%cols
		if e.Direction == ui.DirectionRTL {
			col = cols - 1 - col
		}
		m := ch.Style.Layout.Margin
		b.Children[i].Frame.X = s.Padding.Leading + float32(col)*(cw+s.Gap) + m.Leading
		b.Children[i].Frame.Y = y + m.Top
		flow++
		if flow%cols == 0 {
			y += heights[row] + s.Gap
		}
	}
}
func gridColumns(s ui.LayoutStyle, w float32) int {
	c := s.GridColumns
	if s.GridMinColumnWidth > 0 {
		c = int((w + s.Gap) / (s.GridMinColumnWidth + s.Gap))
	}
	if c < 1 {
		return 1
	}
	return c
}
func flowCount(n *ui.Node) int {
	c := 0
	for _, ch := range n.Children {
		if ch.Style.Layout.Position != ui.PositionAbsolute {
			c++
		}
	}
	return c
}
func (e Engine) placeAbsolute(n *ui.Node, b *Box) {
	for i, ch := range n.Children {
		if ch.Style.Layout.Position != ui.PositionAbsolute {
			continue
		}
		in := ch.Style.Layout.Inset
		if e.Direction == ui.DirectionRTL {
			b.Children[i].Frame.X = b.Frame.Width - in.Leading - b.Children[i].Frame.Width
		} else {
			b.Children[i].Frame.X = in.Leading
		}
		b.Children[i].Frame.Y = in.Top
	}
}
func crossOffset(a ui.AxisAlignment, available, size, leading, trailing float32) float32 {
	switch a {
	case ui.AlignCenter:
		return leading + maxZero(available-leading-trailing-size)/2
	case ui.AlignEnd:
		return maxZero(available - trailing - size)
	default:
		return leading
	}
}
func resolve(l ui.Length, parent float32) (float32, bool) {
	switch l.Unit {
	case ui.LengthPoints:
		return maxZero(l.Value), true
	case ui.LengthPercent:
		return maxZero(parent * l.Value / 100), true
	default:
		return 0, false
	}
}
func bound(v float32, minL, maxL ui.Length, cmin, cmax, parent float32) float32 {
	minV, ok := resolve(minL, parent)
	if !ok {
		minV = cmin
	}
	maxV, ok := resolve(maxL, parent)
	if !ok {
		maxV = cmax
	}
	if maxV > 0 && v > maxV {
		v = maxV
	}
	if v < minV {
		v = minV
	}
	return maxZero(v)
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

package ui

import "fmt"

// CachePolicy defines network and local caching behavior for media resources.
type CachePolicy uint8

const (
	// CacheDefault uses the platform's default caching policy.
	CacheDefault CachePolicy = iota

	// CacheReload ignores local cached data and reloads the resource from the origin source.
	CacheReload

	// CacheReturnCacheDataElseLoad uses cached data if available regardless of age,
	// falling back to loading from origin if not cached.
	CacheReturnCacheDataElseLoad
)

func (c CachePolicy) String() string {
	switch c {
	case CacheDefault:
		return "CacheDefault"
	case CacheReload:
		return "CacheReload"
	case CacheReturnCacheDataElseLoad:
		return "CacheReturnCacheDataElseLoad"
	default:
		return fmt.Sprintf("CachePolicy(%d)", c)
	}
}

// MediaView creates an enhanced media element supporting caching policies,
// placeholders, resize modes, and corner rounding.
func MediaView(source string) *element {
	e := newElement(NodeImage, Props{
		ImageSource: source,
		ImageMode:   ImageFit,
		AccessRole:  RoleImage,
	})
	e.node.CachePolicy = CacheDefault
	return e
}

// Placeholder sets a fallback or loading component displayed while the media loads,
// or a placeholder text string when applied to text input elements.
func (e *element) Placeholder(value any) *element {
	switch v := value.(type) {
	case string:
		e.node.Props.Placeholder = v
	case Component:
		e.node.Placeholder = v
	}
	return e
}

// CachePolicy sets the caching policy for media retrieval.
func (e *element) CachePolicy(policy CachePolicy) *element {
	e.node.CachePolicy = policy
	return e
}

// Circle creates a circular container element with equal width, height, and half-diameter corner radius.
func Circle(diameter float32) *element {
	if diameter < 0 {
		diameter = 0
	}
	return View().
		Width(diameter).
		Height(diameter).
		CornerRadius(diameter / 2)
}

// Rect creates a rectangular container element with specified width and height.
func Rect(width, height float32) *element {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return View().
		Width(width).
		Height(height)
}

// Capsule creates a pill-shaped container element with fully rounded ends.
// Its corner radius is half of the smaller dimension.
func Capsule(width, height float32) *element {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	radius := width / 2
	if height < width {
		radius = height / 2
	}
	return View().
		Width(width).
		Height(height).
		CornerRadius(radius)
}

// Fill sets the background fill color of an element or shape.
func (e *element) Fill(color Color) *element {
	e.node.Style.Appearance.Background = color
	return e
}

// FillColor is an alias for Fill.
func (e *element) FillColor(color Color) *element {
	return e.Fill(color)
}

// Stroke sets the border stroke width and color.
func (e *element) Stroke(width float32, color Color) *element {
	if width < 0 {
		width = 0
	}
	e.node.Style.Appearance.Border = Border{Width: width, Color: color}
	return e
}

// StrokeBorder is an alias for Stroke.
func (e *element) StrokeBorder(width float32, color Color) *element {
	return e.Stroke(width, color)
}

// StrokeWidth sets the border stroke width, preserving any existing stroke color.
func (e *element) StrokeWidth(width float32) *element {
	if width < 0 {
		width = 0
	}
	e.node.Style.Appearance.Border.Width = width
	return e
}

// StrokeColor sets the border stroke color, preserving any existing stroke width.
func (e *element) StrokeColor(color Color) *element {
	e.node.Style.Appearance.Border.Color = color
	return e
}

// DropShadow configures the appearance shadow with color, offset, and blur radius.
func (e *element) DropShadow(color Color, offsetX, offsetY, blur float32) *element {
	e.node.Style.Appearance.Shadow = Shadow{
		Color:   color,
		Offset:  Point{X: offsetX, Y: offsetY},
		Blur:    blur,
		Opacity: 1.0,
	}
	return e
}

// ShadowColor sets the appearance shadow color, preserving existing shadow geometry.
func (e *element) ShadowColor(color Color) *element {
	e.node.Style.Appearance.Shadow.Color = color
	return e
}

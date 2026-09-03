package ui

import "time"

// Color is a portable 8-bit RGBA color.
type Color struct{ R, G, B, A uint8 }

func RGB(r, g, b uint8) Color       { return Color{R: r, G: g, B: b, A: 255} }
func RGBA(r, g, b, a uint8) Color   { return Color{R: r, G: g, B: b, A: a} }
func (c Color) IsTransparent() bool { return c.A == 0 }

type LengthUnit uint8

const (
	LengthAuto LengthUnit = iota
	LengthPoints
	LengthPercent
)

// Length represents an automatic, logical-point, or percentage dimension.
type Length struct {
	Value float32
	Unit  LengthUnit
}

func Points(value float32) Length  { return Length{Value: value, Unit: LengthPoints} }
func Percent(value float32) Length { return Length{Value: value, Unit: LengthPercent} }

type EdgeInsets struct{ Top, Leading, Bottom, Trailing float32 }

func Insets(all float32) EdgeInsets { return EdgeInsets{all, all, all, all} }
func InsetsXY(horizontal, vertical float32) EdgeInsets {
	return EdgeInsets{Top: vertical, Leading: horizontal, Bottom: vertical, Trailing: horizontal}
}

type Size struct{ Width, Height float32 }
type Point struct{ X, Y float32 }

type FlexDirection uint8

const (
	FlexColumn FlexDirection = iota
	FlexRow
)

type WrapMode uint8

const (
	NoWrap WrapMode = iota
	Wrap
)

type PositionMode uint8

const (
	PositionRelative PositionMode = iota
	PositionAbsolute
)

type Overflow uint8

const (
	OverflowVisible Overflow = iota
	OverflowHidden
	OverflowScroll
)

type Visibility uint8

const (
	Visible Visibility = iota
	Hidden
	Gone
)

type Border struct {
	Width float32
	Color Color
}

type Shadow struct {
	Color   Color
	Offset  Point
	Blur    float32
	Spread  float32
	Opacity float32
}

type Transform struct {
	Translate Point
	ScaleX    float32
	ScaleY    float32
	Rotation  float32
}

// LayoutStyle is the portable constraint and flex model. Unsupported fields are
// retained on Node today and will become wire-visible with the typed protocol.
type LayoutStyle struct {
	Width, Height       Length
	MinWidth, MinHeight Length
	MaxWidth, MaxHeight Length
	Margin, Padding     EdgeInsets
	Gap                 float32
	Direction           FlexDirection
	Wrap                WrapMode
	Alignment           AxisAlignment
	FlexGrow            float32
	FlexShrink          float32
	FlexBasis           Length
	AspectRatio         float32
	Position            PositionMode
	Inset               EdgeInsets
	Overflow            Overflow
	GridColumns         int
	GridMinColumnWidth  float32
}

type AppearanceStyle struct {
	Background   Color
	Foreground   Color
	Border       Border
	CornerRadius float32
	Shadow       Shadow
	Opacity      float32
	Transform    Transform
	Visibility   Visibility
}

type TextStyle struct {
	FontFamily    string
	FontSize      float32
	FontWeight    uint16
	LineHeight    float32
	LetterSpacing float32
	Color         Color
}

type InteractionStyle struct {
	Disabled bool
	HitSlop  EdgeInsets
}

// Style groups portable layout, appearance, text, and interaction properties.
type Style struct {
	Layout      LayoutStyle
	Appearance  AppearanceStyle
	Text        TextStyle
	Interaction InteractionStyle
}

// Merge overlays non-zero values from override. It is intended for theme and
// component variants; explicit zero-valued properties should use constructors.
func (s Style) Merge(override Style) Style {
	if override.Layout.Width.Unit != LengthAuto {
		s.Layout.Width = override.Layout.Width
	}
	if override.Layout.Height.Unit != LengthAuto {
		s.Layout.Height = override.Layout.Height
	}
	if override.Layout.MinWidth.Unit != LengthAuto {
		s.Layout.MinWidth = override.Layout.MinWidth
	}
	if override.Layout.MinHeight.Unit != LengthAuto {
		s.Layout.MinHeight = override.Layout.MinHeight
	}
	if override.Layout.MaxWidth.Unit != LengthAuto {
		s.Layout.MaxWidth = override.Layout.MaxWidth
	}
	if override.Layout.MaxHeight.Unit != LengthAuto {
		s.Layout.MaxHeight = override.Layout.MaxHeight
	}
	if override.Layout.Margin != (EdgeInsets{}) {
		s.Layout.Margin = override.Layout.Margin
	}
	if override.Layout.Padding != (EdgeInsets{}) {
		s.Layout.Padding = override.Layout.Padding
	}
	if override.Layout.Gap != 0 {
		s.Layout.Gap = override.Layout.Gap
	}
	if override.Layout.Direction != 0 {
		s.Layout.Direction = override.Layout.Direction
	}
	if override.Layout.Wrap != 0 {
		s.Layout.Wrap = override.Layout.Wrap
	}
	if override.Layout.Alignment != 0 {
		s.Layout.Alignment = override.Layout.Alignment
	}
	if override.Layout.FlexGrow != 0 {
		s.Layout.FlexGrow = override.Layout.FlexGrow
	}
	if override.Layout.FlexShrink != 0 {
		s.Layout.FlexShrink = override.Layout.FlexShrink
	}
	if override.Layout.FlexBasis.Unit != LengthAuto {
		s.Layout.FlexBasis = override.Layout.FlexBasis
	}
	if override.Layout.AspectRatio != 0 {
		s.Layout.AspectRatio = override.Layout.AspectRatio
	}
	if override.Layout.Position != 0 {
		s.Layout.Position = override.Layout.Position
	}
	if override.Layout.Inset != (EdgeInsets{}) {
		s.Layout.Inset = override.Layout.Inset
	}
	if override.Layout.Overflow != 0 {
		s.Layout.Overflow = override.Layout.Overflow
	}
	if override.Layout.GridColumns != 0 {
		s.Layout.GridColumns = override.Layout.GridColumns
	}
	if override.Layout.GridMinColumnWidth != 0 {
		s.Layout.GridMinColumnWidth = override.Layout.GridMinColumnWidth
	}
	if override.Appearance != (AppearanceStyle{}) {
		s.Appearance = override.Appearance
	}
	if override.Text != (TextStyle{}) {
		s.Text = override.Text
	}
	if override.Interaction != (InteractionStyle{}) {
		s.Interaction = override.Interaction
	}
	return s
}

type Theme struct {
	Name              string
	Colors            map[string]Color
	Typography        map[string]TextStyle
	Spacing           map[string]float32
	Radii             map[string]float32
	Elevations        map[string]Shadow
	Motion            map[string]time.Duration
	IconSizes         map[string]float32
	ControlSizes      map[string]Size
	ComponentVariants map[string]Style
}

// Token is a typed semantic theme lookup with a deterministic fallback.
type Token[T any] struct {
	Name     string
	Fallback T
	resolve  func(Theme, string) (T, bool)
}

func (t Token[T]) Resolve(theme Theme) T {
	if t.resolve != nil {
		if value, ok := t.resolve(theme, t.Name); ok {
			return value
		}
	}
	return t.Fallback
}

func ColorToken(name string, fallback Color) Token[Color] {
	return Token[Color]{Name: name, Fallback: fallback, resolve: func(t Theme, key string) (Color, bool) { value, ok := t.Colors[key]; return value, ok }}
}
func SpacingToken(name string, fallback float32) Token[float32] {
	return Token[float32]{Name: name, Fallback: fallback, resolve: func(t Theme, key string) (float32, bool) { value, ok := t.Spacing[key]; return value, ok }}
}
func TypographyToken(name string, fallback TextStyle) Token[TextStyle] {
	return Token[TextStyle]{Name: name, Fallback: fallback, resolve: func(t Theme, key string) (TextStyle, bool) { value, ok := t.Typography[key]; return value, ok }}
}
func StyleToken(name string, fallback Style) Token[Style] {
	return Token[Style]{Name: name, Fallback: fallback, resolve: func(t Theme, key string) (Style, bool) { value, ok := t.ComponentVariants[key]; return value, ok }}
}

func DefaultTheme() Theme {
	return Theme{Name: "default", Colors: map[string]Color{"background": RGB(255, 255, 255), "foreground": RGB(17, 17, 17), "primary": RGB(0, 122, 255)}, Typography: map[string]TextStyle{"body": {FontSize: 16, FontWeight: 400}, "title": {FontSize: 28, FontWeight: 700}}, Spacing: map[string]float32{"xs": 4, "sm": 8, "md": 16, "lg": 24, "xl": 32}, Radii: map[string]float32{"sm": 4, "md": 8, "lg": 16}, Elevations: map[string]Shadow{}, Motion: map[string]time.Duration{"fast": 150 * time.Millisecond, "normal": 250 * time.Millisecond}, IconSizes: map[string]float32{"sm": 16, "md": 24, "lg": 32}, ControlSizes: map[string]Size{"minimumTouch": {Width: 44, Height: 44}}, ComponentVariants: map[string]Style{}}
}

// PlatformStyle is applied after portable style resolution.
type PlatformStyle struct{ IOS, Android Style }

// Breakpoint applies Style when the viewport is at least MinWidth points wide.
type Breakpoint struct {
	MinWidth float32
	Style    Style
}

// ResponsiveStyle resolves breakpoints in declaration order after base style.
func ResponsiveStyle(media MediaQuery, base Style, breakpoints ...Breakpoint) Style {
	resolved := base
	for _, breakpoint := range breakpoints {
		if media.Viewport.Width >= breakpoint.MinWidth {
			resolved = resolved.Merge(breakpoint.Style)
		}
	}
	return resolved
}

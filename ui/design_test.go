package ui

import "testing"

func TestSemanticTokensResolveAndFallback(t *testing.T) {
	theme := DefaultTheme()
	if got := ColorToken("primary", RGB(1, 2, 3)).Resolve(theme); got != RGB(0, 122, 255) {
		t.Fatalf("unexpected primary: %#v", got)
	}
	if got := SpacingToken("missing", 13).Resolve(theme); got != 13 {
		t.Fatalf("unexpected fallback: %v", got)
	}
}

func TestDependenciesAreTyped(t *testing.T) {
	deps := NewDependencies()
	key := DependencyKey[string]{Name: "api.baseURL"}
	Provide(deps, key, "https://example.test")
	if got, ok := Resolve(deps, key); !ok || got != "https://example.test" {
		t.Fatalf("resolve = %q, %v", got, ok)
	}
	if _, ok := Resolve(deps, DependencyKey[int]{Name: "api.baseURL"}); ok {
		t.Fatal("wrong dependency type resolved")
	}
}

type contextualFixture struct{}

func (contextualFixture) Build() *Node { return Text("legacy").Build() }
func (contextualFixture) BuildContext(ctx BuildContext) *Node {
	return Text(ctx.Environment.Locale).Build()
}

func TestBuildWithContextPrefersContextContract(t *testing.T) {
	env := DefaultEnvironment()
	env.Locale = "fr"
	if got := BuildWithContext(contextualFixture{}, NewBuildContext(env)).Props.Text; got != "fr" {
		t.Fatalf("text = %q", got)
	}
}

func TestStyleMergeOverridesPortableValues(t *testing.T) {
	base := Style{Layout: LayoutStyle{Width: Points(100), Padding: Insets(8)}, Text: TextStyle{FontSize: 14}}
	override := Style{Layout: LayoutStyle{Width: Percent(50)}, Text: TextStyle{FontSize: 18}}
	got := base.Merge(override)
	if got.Layout.Width != Percent(50) || got.Layout.Padding != Insets(8) || got.Text.FontSize != 18 {
		t.Fatalf("merged style = %#v", got)
	}
}

func TestResponsiveStyleAppliesMatchingBreakpointsInOrder(t *testing.T) {
	base := Style{Layout: LayoutStyle{Gap: 8, GridColumns: 1}}
	resolved := ResponsiveStyle(MediaQuery{Viewport: Size{Width: 900}}, base,
		Breakpoint{MinWidth: 600, Style: Style{Layout: LayoutStyle{Gap: 16, GridColumns: 2}}},
		Breakpoint{MinWidth: 840, Style: Style{Layout: LayoutStyle{GridColumns: 3}}},
	)
	if resolved.Layout.Gap != 16 || resolved.Layout.GridColumns != 3 {
		t.Fatalf("resolved = %#v", resolved.Layout)
	}
}

func TestResolvePlatformStyleFallsBackAndOverridesGroups(t *testing.T) {
	portable := Style{Layout: LayoutStyle{Gap: 8}, Appearance: AppearanceStyle{Background: RGB(1, 2, 3), CornerRadius: 4}, Text: TextStyle{FontSize: 14}}
	overrides := PlatformStyle{Android: Style{Appearance: AppearanceStyle{Background: RGB(9, 8, 7)}, Text: TextStyle{FontSize: 18}}}
	got := ResolvePlatformStyle(portable, overrides, PlatformAndroid)
	if got.Layout.Gap != 8 || got.Appearance.Background != RGB(9, 8, 7) || got.Appearance.CornerRadius != 0 || got.Text.FontSize != 18 {
		t.Fatalf("resolved style = %#v", got)
	}
}

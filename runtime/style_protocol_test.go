package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/go-native/go-native/ui"
)

func TestTypedStylesRoundTrip(t *testing.T) {
	portable := ui.Style{Layout: ui.LayoutStyle{Width: ui.Percent(75), MinHeight: ui.Points(44), Margin: ui.EdgeInsets{Top: 1, Leading: 2, Bottom: 3, Trailing: 4}, Padding: ui.Insets(8), Gap: 6, Direction: ui.FlexRow, Wrap: ui.Wrap, Alignment: ui.AlignCenter, FlexGrow: 2, FlexShrink: 1, FlexBasis: ui.Points(80), AspectRatio: 1.5, Position: ui.PositionAbsolute, Inset: ui.InsetsXY(5, 7), Overflow: ui.OverflowHidden, GridColumns: 3, GridMinColumnWidth: 120}, Appearance: ui.AppearanceStyle{Background: ui.RGBA(1, 2, 3, 4), Foreground: ui.RGB(5, 6, 7), Border: ui.Border{Width: 2, Color: ui.RGB(8, 9, 10)}, CornerRadius: 12, Opacity: .8, Visibility: ui.Hidden}, Text: ui.TextStyle{FontFamily: "Inter", FontSize: 17, FontWeight: 600, LineHeight: 22, LetterSpacing: .2, Color: ui.RGB(11, 12, 13)}, Interaction: ui.InteractionStyle{Disabled: true, HitSlop: ui.Insets(10)}}
	platform := ui.PlatformStyle{IOS: ui.Style{Layout: ui.LayoutStyle{Gap: 9}}, Android: ui.Style{Appearance: ui.AppearanceStyle{CornerRadius: 14}}}
	data, err := MarshalTypedStyles(portable, platform)
	if err != nil {
		t.Fatal(err)
	}
	got, gp, err := UnmarshalTypedStyles(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, portable) || !reflect.DeepEqual(gp, platform) {
		t.Fatalf("round trip mismatch\n%#v\n%#v", got, gp)
	}
	hash := sha256.Sum256(data)
	if value := hex.EncodeToString(hash[:]); value != "b31999329fed1cafb8c6f30e7110c6b7abf7e23c569b801903d2c1c76434f427" {
		t.Fatalf("typed style golden hash = %s", value)
	}
}

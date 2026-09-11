package layout

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/go-native/go-native/ui"
)

func TestMeasurementRequestProtocolRoundTrip(t *testing.T) {
	in := []MeasurementRequest{{ID: 7, NodeType: ui.NodeText, Text: "Hello", TextProps: TextMeasurementProps{Wrap: ui.TextWrap, Overflow: ui.TextOverflowEllipsisTail, MaxLines: 2, Selectable: true, RichText: "spans", Placeholder: "Search", InputKind: ui.InputSearch, Secure: true, Multiline: true, MaxLength: 80}, Style: ui.Style{Text: ui.TextStyle{FontFamily: "Inter", FontSize: 17, FontWeight: 600}}, Constraints: Constraints{MinWidth: 10, MaxWidth: 200, MaxHeight: 80}}}
	data, err := MarshalMeasurementRequests(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := UnmarshalMeasurementRequests(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("requests mismatch: %#v", out)
	}
	hash := sha256.Sum256(data)
	if got := hex.EncodeToString(hash[:]); got != "3007dbc336adbc1b94614e97fe76df05f0ffc4cfb453df80e953233e61f353fd" {
		t.Fatalf("measurement request golden hash = %s", got)
	}
}

func TestMeasurementResultProtocolRoundTrip(t *testing.T) {
	in := []MeasurementResult{{ID: 1, Size: ui.Size{Width: 44.5, Height: 20}}, {ID: 2, Err: "unavailable"}}
	data, err := MarshalMeasurementResults(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := UnmarshalMeasurementResults(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("results mismatch: %#v", out)
	}
}

func TestMeasurementProtocolRejectsTrailingData(t *testing.T) {
	data, err := MarshalMeasurementResults(nil)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, 1)
	if _, err = UnmarshalMeasurementResults(data); err == nil {
		t.Fatal("expected trailing data error")
	}
}

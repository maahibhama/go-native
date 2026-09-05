package layout

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/go-native/go-native/ui"
)

func TestMeasurementRequestProtocolRoundTrip(t *testing.T) {
	in := []MeasurementRequest{{ID: 7, NodeType: ui.NodeText, Text: "Hello", Style: ui.Style{Text: ui.TextStyle{FontFamily: "Inter", FontSize: 17, FontWeight: 600}}, Constraints: Constraints{MinWidth: 10, MaxWidth: 200, MaxHeight: 80}}}
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
	if got := hex.EncodeToString(hash[:]); got != "fe0d0d0d69fc9c63260bdd443d0447de0d6ec6d9082491cb71f0c10c94a163d8" {
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

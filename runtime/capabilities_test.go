package runtime

import "testing"

func TestNegotiateProtocol(t *testing.T) {
	got, err := NegotiateProtocol(ProtocolCapabilities{MinimumVersion: 7, MaximumVersion: 8, Features: CapabilityTypedStyle | CapabilityNativeMeasurement}, ProtocolCapabilities{MinimumVersion: 8, MaximumVersion: 9, Features: CapabilityTypedStyle})
	if err != nil {
		t.Fatal(err)
	}
	if got.MaximumVersion != 8 || got.Features != CapabilityTypedStyle {
		t.Fatalf("got %#v", got)
	}
	if _, err = NegotiateProtocol(ProtocolCapabilities{MinimumVersion: 7, MaximumVersion: 7}, ProtocolCapabilities{MinimumVersion: 8, MaximumVersion: 8}); err == nil {
		t.Fatal("expected incompatible version error")
	}
}

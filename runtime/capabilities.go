package runtime

import "fmt"

type Capability uint64

const (
	CapabilityTypedStyle Capability = 1 << iota
	CapabilityNativeMeasurement
	CapabilityFocus
	CapabilityLifecycle
	CapabilityComputedGeometry
)

type ProtocolCapabilities struct {
	MinimumVersion uint16
	MaximumVersion uint16
	Features       Capability
}

func CurrentProtocolCapabilities() ProtocolCapabilities {
	return ProtocolCapabilities{MinimumVersion: protocolVersion, MaximumVersion: protocolVersion, Features: CapabilityTypedStyle | CapabilityNativeMeasurement | CapabilityFocus | CapabilityLifecycle | CapabilityComputedGeometry}
}

func NegotiateProtocol(local, remote ProtocolCapabilities) (ProtocolCapabilities, error) {
	minimum := local.MinimumVersion
	if remote.MinimumVersion > minimum {
		minimum = remote.MinimumVersion
	}
	maximum := local.MaximumVersion
	if remote.MaximumVersion < maximum {
		maximum = remote.MaximumVersion
	}
	if minimum > maximum {
		return ProtocolCapabilities{}, fmt.Errorf("no compatible mutation protocol: local %d-%d, remote %d-%d", local.MinimumVersion, local.MaximumVersion, remote.MinimumVersion, remote.MaximumVersion)
	}
	return ProtocolCapabilities{MinimumVersion: maximum, MaximumVersion: maximum, Features: local.Features & remote.Features}, nil
}

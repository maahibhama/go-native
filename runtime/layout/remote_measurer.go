package layout

import (
	"context"
	"errors"

	"github.com/go-native/go-native/runtime/devtransport"
)

// RemoteMeasurer forwards intrinsic measurement batches to a persistent native
// development shell without changing the production measurement wire format.
type RemoteMeasurer struct{ Peer *devtransport.Peer }

func (measurer *RemoteMeasurer) MeasureBatch(ctx context.Context, requests []MeasurementRequest) ([]MeasurementResult, error) {
	if measurer == nil || measurer.Peer == nil {
		return nil, errors.New("remote measurer: disconnected")
	}
	payload, err := MarshalMeasurementRequests(requests)
	if err != nil {
		return nil, err
	}
	frame, err := measurer.Peer.Request(ctx, devtransport.KindMeasurementRequest, payload)
	if err != nil {
		return nil, err
	}
	if frame.Kind == devtransport.KindError {
		return nil, errors.New(string(frame.Payload))
	}
	if frame.Kind != devtransport.KindMeasurementResponse {
		return nil, errors.New("remote measurer: unexpected response")
	}
	return UnmarshalMeasurementResults(frame.Payload)
}

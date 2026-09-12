package devtransport

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestPeerRoutesResponsesAndIncomingFrames(t *testing.T) {
	left, right := net.Pipe()
	peer := NewPeer(left, 3)
	defer peer.Close()
	defer right.Close()
	go func() {
		request, _ := ReadFrame(right)
		_ = WriteFrame(right, Frame{Kind: KindMeasurementResponse, Generation: 3, RequestID: request.RequestID, Payload: []byte("result")})
		_ = WriteFrame(right, Frame{Kind: KindLifecycle, Generation: 3, Payload: []byte{2}})
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := peer.Request(ctx, KindMeasurementRequest, []byte("measure"))
	if err != nil || string(response.Payload) != "result" {
		t.Fatalf("response=%#v err=%v", response, err)
	}
	select {
	case incoming := <-peer.Incoming():
		if incoming.Kind != KindLifecycle {
			t.Fatalf("incoming=%#v", incoming)
		}
	case <-ctx.Done():
		t.Fatal("missing incoming frame")
	}
}

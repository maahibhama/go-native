package devtransport

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

type shortWriter struct{ bytes.Buffer }

func (writer *shortWriter) Write(data []byte) (int, error) {
	if len(data) > 3 {
		data = data[:3]
	}
	return writer.Buffer.Write(data)
}

func TestFrameRoundTripWithPartialWritesAndReads(t *testing.T) {
	var output shortWriter
	want := Frame{Kind: KindMutationBatch, Flags: 2, Generation: 9, RequestID: 4, Payload: []byte("payload")}
	if err := WriteFrame(&output, want); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFrame(io.LimitReader(bytes.NewReader(output.Bytes()), int64(output.Len())))
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != want.Kind || got.Generation != want.Generation || got.RequestID != want.RequestID || !bytes.Equal(got.Payload, want.Payload) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestFrameRejectsInvalidAndOversizedHeaders(t *testing.T) {
	if _, err := ReadFrame(bytes.NewReader(make([]byte, headerSize))); err == nil {
		t.Fatal("accepted invalid magic")
	}
	header := make([]byte, headerSize)
	copy(header, magic[:])
	binary.LittleEndian.PutUint16(header[4:6], Version)
	header[6] = byte(KindStatus)
	binary.LittleEndian.PutUint32(header[20:24], MaxPayload+1)
	if _, err := ReadFrame(bytes.NewReader(header)); err == nil {
		t.Fatal("accepted oversized payload")
	}
}

func TestHelloRoundTripAndAuthentication(t *testing.T) {
	var token [32]byte
	copy(token[:], []byte("secret"))
	want := Hello{Platform: PlatformAndroid, MutationMinimum: 10, MutationMaximum: 10, MeasurementVersion: 2, Capabilities: 7, Token: token}
	got, err := UnmarshalHello(want.MarshalBinary())
	if err != nil || got != want {
		t.Fatalf("got %#v, %v", got, err)
	}
	if !Authenticate(token, got.Token) {
		t.Fatal("token rejected")
	}
}

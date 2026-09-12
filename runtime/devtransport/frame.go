// Package devtransport defines the debug-only binary transport used between a
// host Go runtime, the gonative development broker, and persistent native shells.
package devtransport

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	Version    uint16 = 1
	MaxPayload        = 16 << 20
	headerSize        = 24
)

var magic = [4]byte{'G', 'N', 'D', 'R'}

type Kind uint8

const (
	KindHello Kind = iota + 1
	KindHelloAccepted
	KindMutationBatch
	KindMeasurementRequest
	KindMeasurementResponse
	KindActionEvent
	KindValueEvent
	KindBooleanEvent
	KindGestureEvent
	KindSelectionEvent
	KindFocusEvent
	KindViewport
	KindLifecycle
	KindBatchApplied
	KindResetTree
	KindStatus
	KindError
	KindPing
	KindPong
)

type Frame struct {
	Kind       Kind
	Flags      uint8
	Generation uint64
	RequestID  uint32
	Payload    []byte
}

func WriteFrame(writer io.Writer, frame Frame) error {
	if !validKind(frame.Kind) {
		return fmt.Errorf("dev transport: invalid kind %d", frame.Kind)
	}
	if len(frame.Payload) > MaxPayload {
		return errors.New("dev transport: payload exceeds limit")
	}
	header := make([]byte, headerSize)
	copy(header[:4], magic[:])
	binary.LittleEndian.PutUint16(header[4:6], Version)
	header[6], header[7] = byte(frame.Kind), frame.Flags
	binary.LittleEndian.PutUint64(header[8:16], frame.Generation)
	binary.LittleEndian.PutUint32(header[16:20], frame.RequestID)
	binary.LittleEndian.PutUint32(header[20:24], uint32(len(frame.Payload)))
	if err := writeFull(writer, header); err != nil {
		return err
	}
	return writeFull(writer, frame.Payload)
}

func ReadFrame(reader io.Reader) (Frame, error) {
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(reader, header); err != nil {
		return Frame{}, err
	}
	if string(header[:4]) != string(magic[:]) {
		return Frame{}, errors.New("dev transport: invalid magic")
	}
	if version := binary.LittleEndian.Uint16(header[4:6]); version != Version {
		return Frame{}, fmt.Errorf("dev transport: unsupported version %d", version)
	}
	kind := Kind(header[6])
	if !validKind(kind) {
		return Frame{}, fmt.Errorf("dev transport: invalid kind %d", kind)
	}
	length := binary.LittleEndian.Uint32(header[20:24])
	if length > MaxPayload {
		return Frame{}, errors.New("dev transport: payload exceeds limit")
	}
	payload := make([]byte, int(length))
	if _, err := io.ReadFull(reader, payload); err != nil {
		return Frame{}, err
	}
	return Frame{Kind: kind, Flags: header[7], Generation: binary.LittleEndian.Uint64(header[8:16]), RequestID: binary.LittleEndian.Uint32(header[16:20]), Payload: payload}, nil
}

func validKind(kind Kind) bool { return kind >= KindHello && kind <= KindPong }

func writeFull(writer io.Writer, data []byte) error {
	for len(data) > 0 {
		written, err := writer.Write(data)
		if err != nil {
			return err
		}
		if written == 0 {
			return io.ErrShortWrite
		}
		data = data[written:]
	}
	return nil
}

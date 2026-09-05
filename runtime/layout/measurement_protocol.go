package layout

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
)

const (
	measurementProtocolVersion uint16 = 1
	maxMeasurementItems               = 100_000
	maxMeasurementPayload             = 16 << 20
)

func MarshalMeasurementRequests(requests []MeasurementRequest) ([]byte, error) {
	if len(requests) > maxMeasurementItems {
		return nil, fmt.Errorf("measurement protocol: too many requests")
	}
	var out bytes.Buffer
	_ = binary.Write(&out, binary.LittleEndian, measurementProtocolVersion)
	_ = binary.Write(&out, binary.LittleEndian, uint32(len(requests)))
	for _, request := range requests {
		_ = binary.Write(&out, binary.LittleEndian, request.ID)
		out.WriteByte(byte(request.NodeType))
		for _, value := range []float32{request.Constraints.MinWidth, request.Constraints.MaxWidth, request.Constraints.MinHeight, request.Constraints.MaxHeight} {
			_ = binary.Write(&out, binary.LittleEndian, value)
		}
		for _, value := range []string{request.Text, request.ImageSource} {
			if err := writeMeasurementString(&out, value); err != nil {
				return nil, err
			}
		}
		style, err := gnruntime.MarshalTypedStyles(request.Style, ui.PlatformStyle{})
		if err != nil {
			return nil, err
		}
		_ = binary.Write(&out, binary.LittleEndian, uint32(len(style)))
		out.Write(style)
		if out.Len() > maxMeasurementPayload {
			return nil, fmt.Errorf("measurement protocol: payload exceeds limit")
		}
	}
	return out.Bytes(), nil
}

func UnmarshalMeasurementRequests(data []byte) ([]MeasurementRequest, error) {
	reader, count, err := measurementHeader(data)
	if err != nil {
		return nil, err
	}
	requests := make([]MeasurementRequest, 0, count)
	for range count {
		var request MeasurementRequest
		if err = binary.Read(reader, binary.LittleEndian, &request.ID); err != nil {
			return nil, err
		}
		kind, e := reader.ReadByte()
		if e != nil {
			return nil, e
		}
		request.NodeType = ui.NodeType(kind)
		for _, target := range []*float32{&request.Constraints.MinWidth, &request.Constraints.MaxWidth, &request.Constraints.MinHeight, &request.Constraints.MaxHeight} {
			if err = binary.Read(reader, binary.LittleEndian, target); err != nil {
				return nil, err
			}
		}
		if request.Text, err = readMeasurementString(reader); err != nil {
			return nil, err
		}
		if request.ImageSource, err = readMeasurementString(reader); err != nil {
			return nil, err
		}
		styleBytes, e := readMeasurementBytes(reader)
		if e != nil {
			return nil, e
		}
		request.Style, _, err = gnruntime.UnmarshalTypedStyles(styleBytes)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	if reader.Len() != 0 {
		return nil, fmt.Errorf("measurement protocol: trailing request bytes")
	}
	return requests, nil
}

func MarshalMeasurementResults(results []MeasurementResult) ([]byte, error) {
	if len(results) > maxMeasurementItems {
		return nil, fmt.Errorf("measurement protocol: too many results")
	}
	var out bytes.Buffer
	_ = binary.Write(&out, binary.LittleEndian, measurementProtocolVersion)
	_ = binary.Write(&out, binary.LittleEndian, uint32(len(results)))
	for _, result := range results {
		_ = binary.Write(&out, binary.LittleEndian, result.ID)
		_ = binary.Write(&out, binary.LittleEndian, result.Size.Width)
		_ = binary.Write(&out, binary.LittleEndian, result.Size.Height)
		if err := writeMeasurementString(&out, result.Err); err != nil {
			return nil, err
		}
	}
	if out.Len() > maxMeasurementPayload {
		return nil, fmt.Errorf("measurement protocol: payload exceeds limit")
	}
	return out.Bytes(), nil
}

func UnmarshalMeasurementResults(data []byte) ([]MeasurementResult, error) {
	reader, count, err := measurementHeader(data)
	if err != nil {
		return nil, err
	}
	results := make([]MeasurementResult, 0, count)
	for range count {
		var result MeasurementResult
		if err = binary.Read(reader, binary.LittleEndian, &result.ID); err != nil {
			return nil, err
		}
		if err = binary.Read(reader, binary.LittleEndian, &result.Size.Width); err != nil {
			return nil, err
		}
		if err = binary.Read(reader, binary.LittleEndian, &result.Size.Height); err != nil {
			return nil, err
		}
		if result.Err, err = readMeasurementString(reader); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if reader.Len() != 0 {
		return nil, fmt.Errorf("measurement protocol: trailing result bytes")
	}
	return results, nil
}

func measurementHeader(data []byte) (*bytes.Reader, uint32, error) {
	if len(data) > maxMeasurementPayload {
		return nil, 0, fmt.Errorf("measurement protocol: payload exceeds limit")
	}
	reader := bytes.NewReader(data)
	var version uint16
	var count uint32
	if binary.Read(reader, binary.LittleEndian, &version) != nil || version != measurementProtocolVersion {
		return nil, 0, fmt.Errorf("measurement protocol: unsupported version")
	}
	if binary.Read(reader, binary.LittleEndian, &count) != nil || count > maxMeasurementItems {
		return nil, 0, fmt.Errorf("measurement protocol: invalid item count")
	}
	return reader, count, nil
}
func writeMeasurementString(out *bytes.Buffer, value string) error {
	if len(value) > gnruntime.MaxProtocolString {
		return fmt.Errorf("measurement protocol: string exceeds limit")
	}
	_ = binary.Write(out, binary.LittleEndian, uint32(len(value)))
	out.WriteString(value)
	return nil
}
func readMeasurementString(reader *bytes.Reader) (string, error) {
	value, err := readMeasurementBytes(reader)
	return string(value), err
}
func readMeasurementBytes(reader *bytes.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	if length > gnruntime.MaxProtocolString || uint64(length) > uint64(reader.Len()) {
		return nil, fmt.Errorf("measurement protocol: invalid byte length")
	}
	value := make([]byte, length)
	_, err := io.ReadFull(reader, value)
	return value, err
}

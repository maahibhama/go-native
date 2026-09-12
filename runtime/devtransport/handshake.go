package devtransport

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Platform uint8

const (
	PlatformIOS Platform = iota + 1
	PlatformAndroid
)

type Hello struct {
	Platform           Platform
	MutationMinimum    uint16
	MutationMaximum    uint16
	MeasurementVersion uint16
	Capabilities       uint64
	Token              [32]byte
}

func (hello Hello) MarshalBinary() []byte {
	data := make([]byte, 47)
	data[0] = byte(hello.Platform)
	binary.LittleEndian.PutUint16(data[1:3], hello.MutationMinimum)
	binary.LittleEndian.PutUint16(data[3:5], hello.MutationMaximum)
	binary.LittleEndian.PutUint16(data[5:7], hello.MeasurementVersion)
	binary.LittleEndian.PutUint64(data[7:15], hello.Capabilities)
	copy(data[15:], hello.Token[:])
	return data
}

func UnmarshalHello(data []byte) (Hello, error) {
	if len(data) != 47 {
		return Hello{}, errors.New("dev transport: invalid hello length")
	}
	var hello Hello
	hello.Platform = Platform(data[0])
	if hello.Platform != PlatformIOS && hello.Platform != PlatformAndroid {
		return Hello{}, errors.New("dev transport: invalid platform")
	}
	hello.MutationMinimum = binary.LittleEndian.Uint16(data[1:3])
	hello.MutationMaximum = binary.LittleEndian.Uint16(data[3:5])
	hello.MeasurementVersion = binary.LittleEndian.Uint16(data[5:7])
	hello.Capabilities = binary.LittleEndian.Uint64(data[7:15])
	copy(hello.Token[:], data[15:])
	return hello, nil
}

func Authenticate(expected, actual [32]byte) bool { return bytes.Equal(expected[:], actual[:]) }

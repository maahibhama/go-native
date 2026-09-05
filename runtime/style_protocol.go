package runtime

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/go-native/go-native/ui"
)

const typedStyleRecordVersion uint16 = 1

// MarshalTypedStyles encodes portable, iOS, and Android styles in declaration
// order using fixed-width little-endian values. It is the nested style record
// used by mutation protocol v9.
func MarshalTypedStyles(portable ui.Style, platform ui.PlatformStyle) ([]byte, error) {
	var out bytes.Buffer
	_ = binary.Write(&out, binary.LittleEndian, typedStyleRecordVersion)
	for _, style := range []ui.Style{portable, platform.IOS, platform.Android} {
		if err := writeStyleValue(&out, reflect.ValueOf(style)); err != nil {
			return nil, err
		}
	}
	if out.Len() > MaxProtocolString {
		return nil, &ProtocolError{Kind: "limit", Detail: "typed style exceeds maximum length"}
	}
	return out.Bytes(), nil
}

func UnmarshalTypedStyles(data []byte) (ui.Style, ui.PlatformStyle, error) {
	if len(data) > MaxProtocolString {
		return ui.Style{}, ui.PlatformStyle{}, &ProtocolError{Kind: "limit", Detail: "typed style exceeds maximum length"}
	}
	r := bytes.NewReader(data)
	var version uint16
	if binary.Read(r, binary.LittleEndian, &version) != nil || version != typedStyleRecordVersion {
		return ui.Style{}, ui.PlatformStyle{}, &ProtocolError{Kind: "version", Detail: "unsupported typed style record"}
	}
	styles := make([]ui.Style, 3)
	for i := range styles {
		if err := readStyleValue(r, reflect.ValueOf(&styles[i]).Elem()); err != nil {
			return ui.Style{}, ui.PlatformStyle{}, err
		}
	}
	if r.Len() != 0 {
		return ui.Style{}, ui.PlatformStyle{}, &ProtocolError{Kind: "trailing-data", Offset: len(data) - r.Len(), Detail: "unexpected typed style bytes"}
	}
	return styles[0], ui.PlatformStyle{IOS: styles[1], Android: styles[2]}, nil
}

func writeStyleValue(w io.Writer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := writeStyleValue(w, v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		value := v.String()
		if len(value) > MaxProtocolString {
			return &ProtocolError{Kind: "limit", Detail: "typed style string exceeds maximum length"}
		}
		if err := binary.Write(w, binary.LittleEndian, uint32(len(value))); err != nil {
			return err
		}
		_, err := io.WriteString(w, value)
		return err
	case reflect.Bool:
		var value uint8
		if v.Bool() {
			value = 1
		}
		return binary.Write(w, binary.LittleEndian, value)
	case reflect.Int:
		value := v.Int()
		if value < mathMinInt32 || value > mathMaxInt32 {
			return fmt.Errorf("typed style int outside int32 range")
		}
		return binary.Write(w, binary.LittleEndian, int32(value))
	case reflect.Uint8:
		return binary.Write(w, binary.LittleEndian, uint8(v.Uint()))
	case reflect.Uint16:
		return binary.Write(w, binary.LittleEndian, uint16(v.Uint()))
	case reflect.Float32:
		return binary.Write(w, binary.LittleEndian, float32(v.Float()))
	default:
		return fmt.Errorf("unsupported typed style field %s", v.Type())
	}
	return nil
}

func readStyleValue(r *bytes.Reader, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := readStyleValue(r, v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		var length uint32
		if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
			return err
		}
		if length > MaxProtocolString {
			return &ProtocolError{Kind: "limit", Offset: int(r.Size()) - r.Len(), Detail: "typed style string exceeds maximum length"}
		}
		value := make([]byte, length)
		if _, err := io.ReadFull(r, value); err != nil {
			return err
		}
		v.SetString(string(value))
	case reflect.Bool:
		value, err := r.ReadByte()
		if err != nil {
			return err
		}
		if value > 1 {
			return &ProtocolError{Kind: "value", Offset: int(r.Size()) - r.Len() - 1, Detail: "invalid typed style boolean"}
		}
		v.SetBool(value == 1)
	case reflect.Int:
		var value int32
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetInt(int64(value))
	case reflect.Uint8:
		value, err := r.ReadByte()
		if err != nil {
			return err
		}
		v.SetUint(uint64(value))
	case reflect.Uint16:
		var value uint16
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetUint(uint64(value))
	case reflect.Float32:
		var value float32
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetFloat(float64(value))
	default:
		return fmt.Errorf("unsupported typed style field %s", v.Type())
	}
	return nil
}

const (
	mathMinInt32 = -1 << 31
	mathMaxInt32 = 1<<31 - 1
)

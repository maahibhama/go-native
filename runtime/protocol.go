package runtime

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/go-native/go-native/ui"
	"io"
)

const protocolVersion uint16 = 10

const (
	MaxProtocolPayload   = 16 << 20
	MaxProtocolMutations = 100_000
	MaxProtocolString    = 1 << 20
)

func ProtocolVersion() uint16 { return protocolVersion }

type ProtocolError struct {
	Kind   string
	Offset int
	Detail string
}

func (e *ProtocolError) Error() string {
	return fmt.Sprintf("mutation protocol %s at byte %d: %s", e.Kind, e.Offset, e.Detail)
}

// MarshalBinary encodes a batch for one coarse-grained native call.
func (b MutationBatch) MarshalBinary() ([]byte, error) {
	if len(b.Mutations) > MaxProtocolMutations {
		return nil, &ProtocolError{Kind: "limit", Detail: "too many mutations"}
	}
	var out bytes.Buffer
	_ = binary.Write(&out, binary.LittleEndian, protocolVersion)
	_ = binary.Write(&out, binary.LittleEndian, uint32(len(b.Mutations)))
	_ = binary.Write(&out, binary.LittleEndian, b.Sequence)
	for _, m := range b.Mutations {
		out.WriteByte(byte(m.Type))
		out.WriteByte(byte(m.NodeType))
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.NodeID))
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.ParentID))
		_ = binary.Write(&out, binary.LittleEndian, m.Index)
		_ = binary.Write(&out, binary.LittleEndian, m.FromIndex)
		_ = binary.Write(&out, binary.LittleEndian, m.Props.Width)
		_ = binary.Write(&out, binary.LittleEndian, m.Props.Height)
		_ = binary.Write(&out, binary.LittleEndian, m.Props.Padding)
		_ = binary.Write(&out, binary.LittleEndian, m.Props.Gap)
		out.WriteByte(byte(m.Props.Alignment))
		if m.Props.Bold {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		_ = binary.Write(&out, binary.LittleEndian, m.Props.FontSize)
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.Props.OnPress))
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.Props.OnChange))
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.Props.OnToggle))
		if m.Props.Checked {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		_ = binary.Write(&out, binary.LittleEndian, m.Props.Progress)
		for _, s := range []string{m.Props.Text, m.Props.AccessLabel, m.Props.AccessHint} {
			if len(s) > MaxProtocolString {
				return nil, &ProtocolError{Kind: "limit", Detail: "string exceeds maximum length"}
			}
			_ = binary.Write(&out, binary.LittleEndian, uint32(len(s)))
			out.WriteString(s)
		}
		out.WriteByte(byte(m.Props.AccessRole))
		if m.Props.Focused {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		if m.Props.ScalesText {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		_ = binary.Write(&out, binary.LittleEndian, uint32(len(m.Props.ImageSource)))
		if len(m.Props.ImageSource) > MaxProtocolString {
			return nil, &ProtocolError{Kind: "limit", Detail: "image source exceeds maximum length"}
		}
		out.WriteString(m.Props.ImageSource)
		out.WriteByte(byte(m.Props.ImageMode))
		if m.Props.Horizontal {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		_ = binary.Write(&out, binary.LittleEndian, uint32(len(m.Props.Interactions)))
		if len(m.Props.Interactions) > MaxProtocolString {
			return nil, &ProtocolError{Kind: "limit", Detail: "interaction payload exceeds maximum length"}
		}
		out.WriteString(m.Props.Interactions)
		out.WriteByte(byte(m.Props.TextWrap))
		out.WriteByte(byte(m.Props.TextOverflow))
		_ = binary.Write(&out, binary.LittleEndian, m.Props.MaxLines)
		writeBool(&out, m.Props.Selectable)
		if err := writeBoundedString(&out, m.Props.RichText, "rich text"); err != nil {
			return nil, err
		}
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.Props.OnLink))
		if err := writeBoundedString(&out, m.Props.Placeholder, "placeholder"); err != nil {
			return nil, err
		}
		out.WriteByte(byte(m.Props.InputMode))
		out.WriteByte(byte(m.Props.InputKind))
		out.WriteByte(byte(m.Props.ReturnKey))
		out.WriteByte(byte(m.Props.Capitalization))
		out.WriteByte(byte(m.Props.AutoCorrect))
		writeBool(&out, m.Props.Secure)
		writeBool(&out, m.Props.Multiline)
		writeBool(&out, m.Props.ReadOnly)
		out.WriteByte(byte(m.Props.Validation))
		if err := writeBoundedString(&out, m.Props.ErrorText, "input error"); err != nil {
			return nil, err
		}
		_ = binary.Write(&out, binary.LittleEndian, m.Props.SelectionStart)
		_ = binary.Write(&out, binary.LittleEndian, m.Props.SelectionEnd)
		_ = binary.Write(&out, binary.LittleEndian, m.Props.MaxLength)
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.Props.OnSubmit))
		_ = binary.Write(&out, binary.LittleEndian, uint64(m.Props.OnSelection))
		styles, err := MarshalTypedStyles(m.Style, m.Platform)
		if err != nil {
			return nil, err
		}
		_ = binary.Write(&out, binary.LittleEndian, uint32(len(styles)))
		out.Write(styles)
		if m.HasFrame {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		for _, value := range []float32{m.Frame.X, m.Frame.Y, m.Frame.Width, m.Frame.Height} {
			_ = binary.Write(&out, binary.LittleEndian, value)
		}
		if out.Len() > MaxProtocolPayload {
			return nil, &ProtocolError{Kind: "limit", Detail: "payload exceeds maximum length"}
		}
	}
	return out.Bytes(), nil
}

func writeBool(out *bytes.Buffer, value bool) {
	if value {
		out.WriteByte(1)
	} else {
		out.WriteByte(0)
	}
}

func writeBoundedString(out *bytes.Buffer, value, field string) error {
	if len(value) > MaxProtocolString {
		return &ProtocolError{Kind: "limit", Detail: field + " exceeds maximum length"}
	}
	_ = binary.Write(out, binary.LittleEndian, uint32(len(value)))
	out.WriteString(value)
	return nil
}

// UnmarshalMutationBatch decodes the renderer protocol (used by tests and non-native renderers).
func UnmarshalMutationBatch(data []byte) (MutationBatch, error) {
	if len(data) > MaxProtocolPayload {
		return MutationBatch{}, &ProtocolError{Kind: "limit", Detail: "payload exceeds maximum length"}
	}
	r := bytes.NewReader(data)
	var version uint16
	var count uint32
	if binary.Read(r, binary.LittleEndian, &version) != nil || version != protocolVersion {
		return MutationBatch{}, errors.New("unsupported mutation protocol")
	}
	if binary.Read(r, binary.LittleEndian, &count) != nil {
		return MutationBatch{}, errors.New("invalid mutation count")
	}
	if count > MaxProtocolMutations {
		return MutationBatch{}, &ProtocolError{Kind: "limit", Offset: 6, Detail: "too many mutations"}
	}
	var sequence uint64
	if binary.Read(r, binary.LittleEndian, &sequence) != nil {
		return MutationBatch{}, errors.New("invalid batch sequence")
	}
	b := MutationBatch{Sequence: sequence, Mutations: make([]Mutation, 0, count)}
	for range count {
		var m Mutation
		a, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		z, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Type = MutationType(a)
		m.NodeType = ui.NodeType(z)
		var nid, pid, h, change, toggle uint64
		fields := []any{&nid, &pid, &m.Index, &m.FromIndex, &m.Props.Width, &m.Props.Height, &m.Props.Padding, &m.Props.Gap}
		for _, f := range fields {
			if e = binary.Read(r, binary.LittleEndian, f); e != nil {
				return MutationBatch{}, e
			}
		}
		al, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		bold, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Props.Alignment = ui.AxisAlignment(al)
		m.Props.Bold = bold != 0
		if e = binary.Read(r, binary.LittleEndian, &m.Props.FontSize); e != nil {
			return MutationBatch{}, e
		}
		if e = binary.Read(r, binary.LittleEndian, &h); e != nil {
			return MutationBatch{}, e
		}
		if e = binary.Read(r, binary.LittleEndian, &change); e != nil {
			return MutationBatch{}, e
		}
		if e = binary.Read(r, binary.LittleEndian, &toggle); e != nil {
			return MutationBatch{}, e
		}
		checked, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		if e = binary.Read(r, binary.LittleEndian, &m.Props.Progress); e != nil {
			return MutationBatch{}, e
		}
		m.NodeID = ui.NodeID(nid)
		m.ParentID = ui.NodeID(pid)
		m.Props.OnPress = ui.HandlerID(h)
		m.Props.OnChange = ui.HandlerID(change)
		m.Props.OnToggle = ui.HandlerID(toggle)
		m.Props.Checked = checked != 0
		strings := []*string{&m.Props.Text, &m.Props.AccessLabel, &m.Props.AccessHint}
		for _, dst := range strings {
			var n uint32
			if e = binary.Read(r, binary.LittleEndian, &n); e != nil {
				return MutationBatch{}, e
			}
			if n > MaxProtocolString {
				return MutationBatch{}, &ProtocolError{Kind: "limit", Offset: len(data) - r.Len(), Detail: "string exceeds maximum length"}
			}
			buf := make([]byte, n)
			if _, e = io.ReadFull(r, buf); e != nil {
				return MutationBatch{}, e
			}
			*dst = string(buf)
		}
		role, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		focused, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		scales, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Props.AccessRole = ui.AccessibilityRole(role)
		m.Props.Focused = focused != 0
		m.Props.ScalesText = scales != 0
		var imageLength uint32
		if e = binary.Read(r, binary.LittleEndian, &imageLength); e != nil {
			return MutationBatch{}, e
		}
		if imageLength > MaxProtocolString {
			return MutationBatch{}, &ProtocolError{Kind: "limit", Offset: len(data) - r.Len(), Detail: "image source exceeds maximum length"}
		}
		imageBytes := make([]byte, imageLength)
		if _, e = io.ReadFull(r, imageBytes); e != nil {
			return MutationBatch{}, e
		}
		m.Props.ImageSource = string(imageBytes)
		mode, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		horizontal, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Props.ImageMode = ui.ImageResizeMode(mode)
		m.Props.Horizontal = horizontal != 0
		var interactionLength uint32
		if e = binary.Read(r, binary.LittleEndian, &interactionLength); e != nil {
			return MutationBatch{}, e
		}
		if interactionLength > MaxProtocolString {
			return MutationBatch{}, &ProtocolError{Kind: "limit", Offset: len(data) - r.Len(), Detail: "interaction payload exceeds maximum length"}
		}
		interactionBytes := make([]byte, interactionLength)
		if _, e = io.ReadFull(r, interactionBytes); e != nil {
			return MutationBatch{}, e
		}
		m.Props.Interactions = string(interactionBytes)
		textWrap, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		textOverflow, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Props.TextWrap = ui.TextWrapMode(textWrap)
		m.Props.TextOverflow = ui.TextOverflow(textOverflow)
		if e = binary.Read(r, binary.LittleEndian, &m.Props.MaxLines); e != nil {
			return MutationBatch{}, e
		}
		selectable, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Props.Selectable = selectable != 0
		if m.Props.RichText, e = readBoundedString(r, data, "rich text"); e != nil {
			return MutationBatch{}, e
		}
		var onLink uint64
		if e = binary.Read(r, binary.LittleEndian, &onLink); e != nil {
			return MutationBatch{}, e
		}
		m.Props.OnLink = ui.HandlerID(onLink)
		if m.Props.Placeholder, e = readBoundedString(r, data, "placeholder"); e != nil {
			return MutationBatch{}, e
		}
		inputMode, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		inputKind, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		returnKey, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		capitalization, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		autoCorrect, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		secure, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		multiline, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		readOnly, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		validation, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.Props.InputMode = ui.InputValueMode(inputMode)
		m.Props.InputKind = ui.InputKind(inputKind)
		m.Props.ReturnKey = ui.ReturnKey(returnKey)
		m.Props.Capitalization = ui.TextCapitalization(capitalization)
		m.Props.AutoCorrect = ui.AutoCorrectMode(autoCorrect)
		m.Props.Secure, m.Props.Multiline, m.Props.ReadOnly = secure != 0, multiline != 0, readOnly != 0
		m.Props.Validation = ui.ValidationState(validation)
		if m.Props.ErrorText, e = readBoundedString(r, data, "input error"); e != nil {
			return MutationBatch{}, e
		}
		var onSubmit, onSelection uint64
		for _, field := range []any{&m.Props.SelectionStart, &m.Props.SelectionEnd, &m.Props.MaxLength, &onSubmit, &onSelection} {
			if e = binary.Read(r, binary.LittleEndian, field); e != nil {
				return MutationBatch{}, e
			}
		}
		m.Props.OnSubmit, m.Props.OnSelection = ui.HandlerID(onSubmit), ui.HandlerID(onSelection)
		var styleLength uint32
		if e = binary.Read(r, binary.LittleEndian, &styleLength); e != nil {
			return MutationBatch{}, e
		}
		if styleLength > MaxProtocolString {
			return MutationBatch{}, &ProtocolError{Kind: "limit", Offset: len(data) - r.Len(), Detail: "typed style exceeds maximum length"}
		}
		styleBytes := make([]byte, styleLength)
		if _, e = io.ReadFull(r, styleBytes); e != nil {
			return MutationBatch{}, e
		}
		m.Style, m.Platform, e = UnmarshalTypedStyles(styleBytes)
		if e != nil {
			return MutationBatch{}, e
		}
		hasFrame, e := r.ReadByte()
		if e != nil {
			return MutationBatch{}, e
		}
		m.HasFrame = hasFrame != 0
		for _, value := range []*float32{&m.Frame.X, &m.Frame.Y, &m.Frame.Width, &m.Frame.Height} {
			if e = binary.Read(r, binary.LittleEndian, value); e != nil {
				return MutationBatch{}, e
			}
		}
		b.Mutations = append(b.Mutations, m)
	}
	if r.Len() != 0 {
		return MutationBatch{}, &ProtocolError{Kind: "trailing-data", Offset: len(data) - r.Len(), Detail: "unexpected bytes after mutation batch"}
	}
	return b, nil
}

func readBoundedString(r *bytes.Reader, data []byte, field string) (string, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	if length > MaxProtocolString {
		return "", &ProtocolError{Kind: "limit", Offset: len(data) - r.Len(), Detail: field + " exceeds maximum length"}
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

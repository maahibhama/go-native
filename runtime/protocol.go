package runtime

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/go-native/go-native/ui"
	"io"
)

const protocolVersion uint16 = 1

// MarshalBinary encodes a batch for one coarse-grained native call.
func (b MutationBatch) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer
	_ = binary.Write(&out, binary.LittleEndian, protocolVersion)
	_ = binary.Write(&out, binary.LittleEndian, uint32(len(b.Mutations)))
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
		for _, s := range []string{m.Props.Text, m.Props.AccessLabel} {
			_ = binary.Write(&out, binary.LittleEndian, uint32(len(s)))
			out.WriteString(s)
		}
	}
	return out.Bytes(), nil
}

// UnmarshalMutationBatch decodes the renderer protocol (used by tests and non-native renderers).
func UnmarshalMutationBatch(data []byte) (MutationBatch, error) {
	r := bytes.NewReader(data)
	var version uint16
	var count uint32
	if binary.Read(r, binary.LittleEndian, &version) != nil || version != protocolVersion {
		return MutationBatch{}, errors.New("unsupported mutation protocol")
	}
	if binary.Read(r, binary.LittleEndian, &count) != nil {
		return MutationBatch{}, errors.New("invalid mutation count")
	}
	b := MutationBatch{Mutations: make([]Mutation, 0, count)}
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
		var nid, pid, h uint64
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
		m.NodeID = ui.NodeID(nid)
		m.ParentID = ui.NodeID(pid)
		m.Props.OnPress = ui.HandlerID(h)
		strings := []*string{&m.Props.Text, &m.Props.AccessLabel}
		for _, dst := range strings {
			var n uint32
			if e = binary.Read(r, binary.LittleEndian, &n); e != nil {
				return MutationBatch{}, e
			}
			buf := make([]byte, n)
			if _, e = io.ReadFull(r, buf); e != nil {
				return MutationBatch{}, e
			}
			*dst = string(buf)
		}
		b.Mutations = append(b.Mutations, m)
	}
	return b, nil
}

package runtime

import (
	"bytes"
	"encoding/binary"
	"github.com/go-native/go-native/ui"
	"math"
)

// marshalInteractions encodes arbitrary ordered intent counts into comparable Props storage.
func marshalInteractions(set ui.IntentSet, handlers []ui.HandlerID) string {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.LittleEndian, uint32(len(set.Gestures)))
	for i, g := range set.Gestures {
		b.WriteByte(byte(g.Kind))
		b.WriteByte(byte(g.Direction))
		_ = binary.Write(&b, binary.LittleEndian, int64(g.MinimumPress))
		_ = binary.Write(&b, binary.LittleEndian, g.MinimumTravel)
		var id ui.HandlerID
		if i < len(handlers) {
			id = handlers[i]
		}
		_ = binary.Write(&b, binary.LittleEndian, uint64(id))
	}
	_ = binary.Write(&b, binary.LittleEndian, uint32(len(set.Animations)))
	for _, a := range set.Animations {
		from, to := a.From, a.To
		if a.Property == ui.AnimateOpacity {
			from = clamp(from, 0, 1)
			to = clamp(to, 0, 1)
		}
		if a.Property == ui.AnimateScale {
			from = maxZero(from)
			to = maxZero(to)
		}
		b.WriteByte(byte(a.Property))
		_ = binary.Write(&b, binary.LittleEndian, int64(a.Duration))
		_ = binary.Write(&b, binary.LittleEndian, int64(a.Delay))
		b.WriteByte(byte(a.Curve))
		_ = binary.Write(&b, binary.LittleEndian, a.SpringDamping)
		_ = binary.Write(&b, binary.LittleEndian, a.SpringVelocity)
		if a.ReduceMotionOK {
			b.WriteByte(1)
		} else {
			b.WriteByte(0)
		}
		for _, value := range []float32{from, to, a.FromX, a.FromY, a.ToX, a.ToY} {
			_ = binary.Write(&b, binary.LittleEndian, finite(value))
		}
	}
	return b.String()
}

func clamp(v, lo, hi float32) float32 {
	v = finite(v)
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
func maxZero(v float32) float32 {
	v = finite(v)
	if v < 0 {
		return 0
	}
	return v
}
func finite(v float32) float32 {
	if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
		return 0
	}
	return v
}

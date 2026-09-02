package runtime

import (
	"bytes"
	"encoding/binary"
	"github.com/go-native/go-native/ui"
	"math"
	"testing"
	"time"
)

func TestInteractionPayloadEncodesTargetsAndClamps(t *testing.T) {
	set := ui.IntentSet{Gestures: []ui.GestureIntent{{Kind: ui.GestureSwipe, Direction: ui.SwipeLeading, MinimumPress: 2 * time.Second, MinimumTravel: 24}}, Animations: []ui.AnimationIntent{{Property: ui.AnimateOpacity, Duration: 200 * time.Millisecond, Delay: 10 * time.Millisecond, Curve: ui.CurveEaseOut, SpringDamping: .7, SpringVelocity: 2, ReduceMotionOK: true, From: -1, To: 2, FromX: 1, FromY: 2, ToX: 3, ToY: 4}, {Property: ui.AnimateScale, From: -2, To: float32(math.Inf(1))}}}
	r := bytes.NewReader([]byte(marshalInteractions(set, []ui.HandlerID{99})))
	var gc uint32
	binary.Read(r, binary.LittleEndian, &gc)
	if gc != 1 {
		t.Fatalf("gestures=%d", gc)
	}
	k, _ := r.ReadByte()
	d, _ := r.ReadByte()
	var press int64
	var travel float32
	var id uint64
	binary.Read(r, binary.LittleEndian, &press)
	binary.Read(r, binary.LittleEndian, &travel)
	binary.Read(r, binary.LittleEndian, &id)
	if k != byte(ui.GestureSwipe) || d != byte(ui.SwipeLeading) || press != int64(2*time.Second) || travel != 24 || id != 99 {
		t.Fatal("gesture descriptor changed")
	}
	var ac uint32
	binary.Read(r, binary.LittleEndian, &ac)
	if ac != 2 {
		t.Fatalf("animations=%d", ac)
	}
	read := func() (float32, float32, [4]float32) {
		r.ReadByte()
		var duration, delay int64
		binary.Read(r, binary.LittleEndian, &duration)
		binary.Read(r, binary.LittleEndian, &delay)
		r.ReadByte()
		var damping, velocity float32
		binary.Read(r, binary.LittleEndian, &damping)
		binary.Read(r, binary.LittleEndian, &velocity)
		r.ReadByte()
		var v [6]float32
		for i := range v {
			binary.Read(r, binary.LittleEndian, &v[i])
		}
		return v[0], v[1], [4]float32{v[2], v[3], v[4], v[5]}
	}
	from, to, vector := read()
	if from != 0 || to != 1 || vector != [4]float32{1, 2, 3, 4} {
		t.Fatalf("opacity %v %v %v", from, to, vector)
	}
	from, to, _ = read()
	if from != 0 || to != 0 {
		t.Fatalf("scale %v %v", from, to)
	}
}

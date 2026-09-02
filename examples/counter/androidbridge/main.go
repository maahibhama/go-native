//go:build android

package main

/*
#include <stdint.h>
void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length);
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/go-native/go-native/examples/counter"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
)

var benchmarkOutput string

type androidRenderer struct{}

func (androidRenderer) Apply(batch gnruntime.MutationBatch) error {
	data, err := batch.MarshalBinary()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty mutation batch")
	}
	C.GNAndroidApplyMutationBatch((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
	return nil
}

var appRuntime *gnruntime.Runtime

//export GoNativeAndroidStart
func GoNativeAndroidStart() {
	appRuntime = gnruntime.New(counter.App, androidRenderer{})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
}

//export GoNativeAndroidDispatchEvent
func GoNativeAndroidDispatchEvent(handler C.uint64_t) {
	if appRuntime != nil {
		appRuntime.Dispatch(ui.HandlerID(handler))
	}
}

//export GoNativeAndroidDispatchValueEvent
func GoNativeAndroidDispatchValueEvent(handler C.uint64_t, value *C.char) {
	if appRuntime != nil {
		appRuntime.DispatchValue(ui.HandlerID(handler), C.GoString(value))
	}
}

//export GoNativeAndroidDispatchBoolEvent
func GoNativeAndroidDispatchBoolEvent(handler C.uint64_t, value C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchBool(ui.HandlerID(handler), value != 0)
	}
}

//export GoNativeAndroidDispatchGestureEvent
func GoNativeAndroidDispatchGestureEvent(handler C.uint64_t, translationX, translationY, velocityX, velocityY C.float) {
	if appRuntime != nil {
		appRuntime.DispatchGesture(ui.HandlerID(handler), ui.GestureEvent{
			TranslationX: float32(translationX), TranslationY: float32(translationY),
			VelocityX: float32(velocityX), VelocityY: float32(velocityY),
		})
	}
}

//export GoNativeAndroidStop
func GoNativeAndroidStop() {
	if appRuntime != nil {
		appRuntime.Stop()
		appRuntime = nil
	}
}

//export GoNativeAndroidReportBatchApplied
func GoNativeAndroidReportBatchApplied(sequence C.uint64_t, nativeNanos C.uint64_t) {
	if appRuntime != nil {
		appRuntime.RecordNativeApply(uint64(sequence), time.Duration(nativeNanos))
		emitTimingSample(uint64(sequence))
	}
}

func emitTimingSample(sequence uint64) {
	if benchmarkOutput != "1" {
		return
	}
	for _, sample := range appRuntime.TimingSamples() {
		if sample.Sequence == sequence {
			fmt.Printf("GONATIVE_TIMING {\"sequence\":%d,\"mutations\":%d,\"native_apply_ns\":%d,\"bridge_to_apply_ns\":%d,\"event_to_apply_ns\":%d}\n", sample.Sequence, sample.MutationCount, sample.NativeApply.Nanoseconds(), sample.BridgeToApply.Nanoseconds(), sample.EventToApply.Nanoseconds())
			return
		}
	}
}

func main() {}

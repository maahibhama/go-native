//go:build android

package main

/*
#include <stdint.h>
#include <stdlib.h>
void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length);
int32_t GNAndroidMeasureBatch(const uint8_t *bytes, int32_t length, uint8_t **result);
void GNAndroidFreeBuffer(uint8_t *bytes);
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unsafe"

	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/runtime/layout"
	"github.com/go-native/go-native/ui"
	"my-project"
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

func (androidRenderer) MeasureBatch(ctx context.Context, requests []layout.MeasurementRequest) ([]layout.MeasurementResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	payload, err := layout.MarshalMeasurementRequests(requests)
	if err != nil {
		return nil, err
	}
	if len(payload) == 0 {
		return nil, errors.New("android measurement: empty request batch")
	}
	var result *C.uint8_t
	length := C.GNAndroidMeasureBatch((*C.uint8_t)(unsafe.Pointer(&payload[0])), C.int32_t(len(payload)), &result)
	if length <= 0 || result == nil {
		return nil, fmt.Errorf("android measurement: native adapter unavailable (%d)", int32(length))
	}
	defer C.GNAndroidFreeBuffer(result)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return layout.UnmarshalMeasurementResults(C.GoBytes(unsafe.Pointer(result), C.int(length)))
}

var appRuntime *gnruntime.Runtime

//export GoNativeAndroidStart
func GoNativeAndroidStart() {
	renderer := androidRenderer{}
	appRuntime = gnruntime.New(app.App, renderer)
	appRuntime.SetLayoutProvider(&layout.Pipeline{Measurer: renderer, Cache: layout.NewMeasurementCache()})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
}

//export GoNativeAndroidUpdateViewport
func GoNativeAndroidUpdateViewport(width, height, scale C.float) {
	if appRuntime == nil || width <= 0 || height <= 0 || scale <= 0 {
		return
	}
	current := appRuntime.Environment().MediaQuery
	if current.Viewport.Width == float32(width) && current.Viewport.Height == float32(height) && current.Scale == float32(scale) {
		return
	}
	appRuntime.UpdateEnvironment(func(environment ui.Environment) ui.Environment {
		environment.MediaQuery.Viewport = ui.Size{Width: float32(width), Height: float32(height)}
		environment.MediaQuery.Scale = float32(scale)
		if width > height {
			environment.MediaQuery.Orientation = ui.OrientationLandscape
		} else {
			environment.MediaQuery.Orientation = ui.OrientationPortrait
		}
		return environment
	})
}

//export GoNativeAndroidSetLifecycle
func GoNativeAndroidSetLifecycle(state C.uint8_t) {
	if appRuntime != nil {
		appRuntime.SetLifecycle(ui.LifecycleState(state))
	}
}

//export GoNativeAndroidDispatchFocus
func GoNativeAndroidDispatchFocus(nodeID C.uint64_t, focused C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchFocus(ui.NodeID(nodeID), focused != 0)
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

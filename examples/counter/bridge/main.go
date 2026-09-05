package main

/*
#include <stdint.h>
#include <stdlib.h>
void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
int32_t GNMeasureNativeBatch(const uint8_t *bytes, int32_t length, uint8_t **results, int32_t *resultLength);
void GNFreeNativeBuffer(void *buffer);
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-native/go-native/examples/counter"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/runtime/layout"
	"github.com/go-native/go-native/ui"
	"time"
	"unsafe"
)

var benchmarkOutput string

type iosRenderer struct{}

// iosNativeMeasurer is the UIKit implementation of the portable batched
// intrinsic-measurement contract. All values cross as owned byte buffers.
type iosNativeMeasurer struct{}

func (iosNativeMeasurer) MeasureBatch(ctx context.Context, requests []layout.MeasurementRequest) ([]layout.MeasurementResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := layout.MarshalMeasurementRequests(requests)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("empty native measurement request")
	}
	var output *C.uint8_t
	var outputLength C.int32_t
	status := C.GNMeasureNativeBatch((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int32_t(len(data)), &output, &outputLength)
	if output != nil {
		defer C.GNFreeNativeBuffer(unsafe.Pointer(output))
	}
	if status != 0 {
		return nil, fmt.Errorf("UIKit measurement failed with status %d", int32(status))
	}
	if output == nil || outputLength <= 0 {
		return nil, errors.New("UIKit measurement returned an empty response")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return layout.UnmarshalMeasurementResults(C.GoBytes(unsafe.Pointer(output), C.int(outputLength)))
}

var _ layout.BatchMeasurer = iosNativeMeasurer{}

func (iosRenderer) Apply(batch gnruntime.MutationBatch) error {
	data, err := batch.MarshalBinary()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty mutation batch")
	}
	C.GNApplyMutationBatch((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
	return nil
}

var appRuntime *gnruntime.Runtime

//export GoNativeStart
func GoNativeStart() {
	appRuntime = gnruntime.New(counter.App, iosRenderer{})
	appRuntime.SetLayoutProvider(&layout.Pipeline{Measurer: iosNativeMeasurer{}, Cache: layout.NewMeasurementCache()})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
}

//export GoNativeSetViewport
func GoNativeSetViewport(width, height, scale C.float) {
	if appRuntime == nil || width <= 0 || height <= 0 {
		return
	}
	current := appRuntime.Environment().MediaQuery
	if current.Viewport.Width == float32(width) && current.Viewport.Height == float32(height) && current.Scale == float32(scale) {
		return
	}
	appRuntime.UpdateEnvironment(func(environment ui.Environment) ui.Environment {
		environment.MediaQuery.Viewport = ui.Size{Width: float32(width), Height: float32(height)}
		environment.MediaQuery.Scale = float32(scale)
		return environment
	})
}

//export GoNativeDispatchEvent
func GoNativeDispatchEvent(handler C.uint64_t) {
	if appRuntime != nil {
		appRuntime.Dispatch(ui.HandlerID(handler))
	}
}

//export GoNativeDispatchValueEvent
func GoNativeDispatchValueEvent(handler C.uint64_t, value *C.char) {
	if appRuntime != nil {
		appRuntime.DispatchValue(ui.HandlerID(handler), C.GoString(value))
	}
}

//export GoNativeDispatchBoolEvent
func GoNativeDispatchBoolEvent(handler C.uint64_t, value C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchBool(ui.HandlerID(handler), value != 0)
	}
}

//export GoNativeDispatchGestureEvent
func GoNativeDispatchGestureEvent(handler C.uint64_t, translationX, translationY, velocityX, velocityY C.float) {
	if appRuntime != nil {
		appRuntime.DispatchGesture(ui.HandlerID(handler), ui.GestureEvent{
			TranslationX: float32(translationX), TranslationY: float32(translationY),
			VelocityX: float32(velocityX), VelocityY: float32(velocityY),
		})
	}
}

//export GoNativeDispatchFocus
func GoNativeDispatchFocus(nodeID C.uint64_t, focused C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchFocus(ui.NodeID(nodeID), focused != 0)
	}
}

//export GoNativeStop
func GoNativeStop() {
	if appRuntime != nil {
		appRuntime.Stop()
		appRuntime = nil
	}
}

//export GoNativeSetLifecycle
func GoNativeSetLifecycle(state C.uint8_t) {
	if appRuntime != nil && state <= C.uint8_t(ui.LifecycleDestroyed) {
		appRuntime.SetLifecycle(ui.LifecycleState(state))
	}
}

//export GoNativeReportBatchApplied
func GoNativeReportBatchApplied(sequence C.uint64_t, nativeNanos C.uint64_t) {
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

//go:build android

package main

/*
#include <stdint.h>
void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length);
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"

	"github.com/go-native/go-native/examples/counter"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
)

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
	}
}

func main() {}

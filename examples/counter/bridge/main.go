package main

/*
#include <stdint.h>
#include <stdlib.h>
void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
*/
import "C"

import (
	"errors"
	"github.com/go-native/go-native/examples/counter"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
	"time"
	"unsafe"
)

type iosRenderer struct{}

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
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
}

//export GoNativeDispatchEvent
func GoNativeDispatchEvent(handler C.uint64_t) {
	if appRuntime != nil {
		appRuntime.Dispatch(ui.HandlerID(handler))
	}
}

//export GoNativeStop
func GoNativeStop() {
	if appRuntime != nil {
		appRuntime.Stop()
		appRuntime = nil
	}
}

//export GoNativeReportBatchApplied
func GoNativeReportBatchApplied(sequence C.uint64_t, nativeNanos C.uint64_t) {
	if appRuntime != nil {
		appRuntime.RecordNativeApply(uint64(sequence), time.Duration(nativeNanos))
	}
}

func main() {}

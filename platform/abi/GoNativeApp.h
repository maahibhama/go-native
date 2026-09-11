#ifndef GO_NATIVE_APP_H
#define GO_NATIVE_APP_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// GO_NATIVE_APP_ABI_VERSION versions native-framework calls into an
// application-specific Go library. It is independent of mutation protocol v10
// and native measurement protocol v2.
#define GO_NATIVE_APP_ABI_VERSION 1u

void GoNativeStart(void);
void GoNativeStop(void);
void GoNativeSetViewport(float width, float height, float scale);
void GoNativeSetLifecycle(uint8_t state);

void GoNativeDispatchEvent(uint64_t handler);
void GoNativeDispatchValueEvent(uint64_t handler, char *value);
void GoNativeDispatchBoolEvent(uint64_t handler, uint8_t value);
void GoNativeDispatchGestureEvent(
    uint64_t handler,
    float translation_x,
    float translation_y,
    float velocity_x,
    float velocity_y);
void GoNativeDispatchSelection(uint64_t handler, int32_t start, int32_t end);
void GoNativeDispatchFocus(uint64_t node_id, uint8_t focused);
void GoNativeReportBatchApplied(uint64_t sequence, uint64_t native_nanos);

#ifdef __cplusplus
}
#endif

#endif

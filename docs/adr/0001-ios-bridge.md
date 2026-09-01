# ADR 0001: iOS uses a c-archive and a batch-oriented C ABI

- Status: Accepted for Milestone 0
- Date: 2026-09-02

## Decision

Build the Go entry package with `-buildmode=c-archive`. Link the archive into a small Objective-C UIKit host. Cross Go-to-native with one `GNApplyMutationBatch(bytes, length)` call per render and native-to-Go with exported start and integer event functions.

## Rationale

UIKit is imperative and fits a retained `NodeID → UIView` renderer. Objective-C consumes the generated C header directly and avoids requiring Swift runtime conventions. A copied, versioned binary batch keeps the ABI small and prevents Go pointers or function pointers from being retained by native code.

## Risks

- Go runtime startup and archive size affect cold start and binary size; measure before making performance claims.
- cgo and Apple SDK/toolchain combinations can regress. CI must compile the simulator artifact on a pinned macOS/Xcode runner.
- Calls from UIKit enter the Go runtime on a native thread. Callbacks must stay short and schedule work rather than block the main thread.
- Batch memory is Go-owned only during the call. Objective-C must copy it before asynchronous main-queue use.
- UIKit operations are main-thread-only. The renderer always dispatches application to the main queue.
- App lifecycle, backgrounding, panic reporting, cancellation, and release-build signing remain future hardening work.

## Alternatives considered

`gomobile bind` offers higher-level bindings but broadens the generated surface and does not remove the need for a custom batched renderer. SwiftUI is optimized around its own declarative state model, while this framework needs direct imperative control of native view identity. Per-property cgo calls were rejected because they make boundary overhead scale with property count.

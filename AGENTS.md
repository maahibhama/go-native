# Go Native — Agent Context & Guidelines

## Overview
Go Native is a declarative Go UI runtime that renders platform-native controls on iOS (UIKit) and Android (Android Views).
- **Zero Web/Canvas Overhead**: No WebView, JavaScript, React Native bridge JSON, Flutter, or Skia canvas.
- **Binary Mutation Batch**: State changes produce virtual UI trees reconciled into a versioned, little-endian binary batch (`runtime.MutationBatch`) sent across the native boundary in **one coarse call per render pass**.
- **Pointer-Free Handler Boundary**: Native event listeners hold only 64-bit integer `HandlerID`s. No Go pointers or objects cross cgo/JNI.

---

## Repository Map
- [`ui/`](./ui): Platform-independent declarative primitives (`View`, `Column`, `Row`, `Text`, `SafeArea`, `Button`, `TextInput`, `Switch`, `ProgressIndicator`, `Image`, `ScrollView`), typed `Props`, reactive `State[T]`, gesture/animation intents, and navigation/modal abstractions.
- [`runtime/`](./runtime): Reconciler (`Reconcile`), identity stabilization (`stabilizeIDs`), binary serialization (`MarshalBinary`/`UnmarshalMutationBatch`), thread-safe `EventRegistry`, diagnostics, and performance timing instrumentation.
- [`runtime/inspector/`](./runtime/inspector): Loopback HTTP diagnostic server (`GET /v1/tree`, `GET /v1/logs`).
- [`platform/ios/`](./platform/ios): Objective-C UIKit host & renderer (`GoNativeRenderer.m`/`.h`, `main.m`).
- [`platform/android/`](./platform/android): Java Views host & renderer (`MainActivity.java`, `GapDrawable.java`, Gradle build configuration).
- [`cmd/gonative/`](./cmd/gonative): Developer CLI (`init`, `doctor`, `build`, `run`, `benchmark native`).
- [`examples/counter/`](./examples/counter): Demo app and cgo/JNI bridge entrypoints (`bridge/main.go`, `androidbridge/main.go`, `androidbridge/jni.c`).
- [`scripts/`](./scripts): Standalone build and run automation scripts for iOS and Android.
- [`docs/`](./docs): Architectural decision records (ADRs), performance baselines, and roadmaps.

---

## Core Invariants & Rules

1. **Protocol Synchronization**:
   - Binary serialization is defined in [`runtime/protocol.go`](./runtime/protocol.go) with `protocolVersion`.
   - Any change to `ui.Props` or binary serialization fields **MUST** be mirrored identically in:
     - [`runtime/protocol.go`](./runtime/protocol.go) (Go encoding & test decoding)
     - [`platform/ios/GoNativeRenderer.m`](./platform/ios/GoNativeRenderer.m) (`GNApply`, `GNStyle`)
     - [`platform/android/src/dev/gonative/counter/MainActivity.java`](./platform/android/src/dev/gonative/counter/MainActivity.java) (`applyOnUiThread`)
   - Protocol version constants in all 3 layers must match.

2. **UI Thread & Concurrency Safety**:
   - `ui.State[T]` updates can happen from any goroutine.
   - `runtime.Runtime.Schedule()` coalesces renders and serializes reconciliation.
   - Native renderers **must** execute DOM/view mutations exclusively on the platform UI thread (`dispatch_async(dispatch_get_main_queue(), ...)` on iOS, `runOnUiThread(...)` on Android).

3. **Memory & Handler Lifecycle**:
   - Native view destruction or replacement must unbind and release event handlers via `releaseTree()` / `EventRegistry.Release(id)`.
   - Native code retains zero Go pointers.

4. **Testing & Quality Gates**:
   - Run tests: `go test -race ./...`
   - Vet code: `go vet ./...`
   - Verify formatting: `gofmt -l .`
   - Check performance: `make benchmark`
   - Check toolchains: `go run ./cmd/gonative doctor`

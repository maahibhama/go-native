---
name: go-native
description: >-
  Workflows, guidelines, architecture rules, build/run procedures, and implementation steps
  for developing, extending, and testing the go-native framework across Go runtime, iOS UIKit,
  and Android Views. Use whenever developing, building, or modifying go-native components,
  protocols, bridges, or examples.
---

# Go Native Development & Extension Guide

This skill provides step-by-step runbooks for developing, extending, testing, and debugging the `go-native` declarative UI runtime across Go, iOS (UIKit), and Android (Views).

---

## 1. Quick Reference Commands

```bash
# Run all Go tests with race detection
go test -race ./...

# Run static analysis and formatting check
go vet ./...
gofmt -l .

# Run Go performance microbenchmarks
make benchmark

# Run the CLI doctor to verify local toolchains
go run ./cmd/gonative doctor

# Build & run iOS Counter
go run ./cmd/gonative build ios
go run ./cmd/gonative run ios

# Build & run Android Counter
go run ./cmd/gonative build android
go run ./cmd/gonative run android

# Run Native Sampling Benchmarks
go run ./cmd/gonative benchmark native ios
go run ./cmd/gonative benchmark native android
```

---

## 2. Core Architecture Rules & Invariants

1. **No Go Pointers Across Foreign Function Boundaries**:
   - Event callbacks use 64-bit integer `HandlerID`s registered in `runtime.EventRegistry`.
   - Native views are tracked by `NodeID` in native dictionaries (`GNViews` on iOS, `LongSparseArray<View>` on Android).
2. **Binary Protocol Synchronization**:
   - Serialization is versioned via `protocolVersion` in `runtime/protocol.go`.
   - Any layout/prop addition requires updating:
     - `ui.Props` in `ui/node.go`
     - `runtime/protocol.go` (`MarshalBinary` and `UnmarshalMutationBatch`)
     - `platform/ios/GoNativeRenderer.m` (`GNApply`, `GNStyle`)
     - `platform/android/src/dev/gonative/counter/MainActivity.java` (`applyOnUiThread`)
3. **UI Thread Guarantee**:
   - Native renderers must execute DOM mutations on the main UI thread (`dispatch_get_main_queue()` on iOS, `runOnUiThread` on Android).
   - Go state updates (`ui.State[T]`) are thread-safe and scheduled via `runtime.Runtime.Schedule()`.

---

## 3. Runbook: Adding a New Native UI Primitive

When adding a new UI primitive (e.g., `Slider`, `DatePicker`, `Canvas`):

### Step 1: UI Layer (`ui/`)
1. Define a new `NodeType` in `ui/node.go` (e.g., `NodeSlider NodeType = ...`).
2. Add necessary properties to `ui.Props` in `ui/node.go`.
3. Create constructor functions and fluent property chainers in `ui/components.go` (e.g., `Slider(...)`).
4. Add unit tests in `ui/components_test.go`.

### Step 2: Runtime Protocol & Diagnostics (`runtime/`)
1. Update `runtime/protocol.go` to serialize new fields in `MarshalBinary` and deserialize in `UnmarshalMutationBatch`.
2. Bump `protocolVersion` if serialization byte layout changes.
3. Update `nodeTypeName()` in `runtime/diagnostics.go` to include the new node type.
4. Add tests in `runtime/protocol_test.go` and `runtime/diagnostics_test.go`.

### Step 3: iOS UIKit Renderer (`platform/ios/`)
1. Update `GNNode` enum in `platform/ios/GoNativeRenderer.m`.
2. Update `GNMake` to instantiate the platform view (e.g., `UISlider`).
3. Update `GNStyle` to apply properties, accessibility traits, and target-actions.
4. If linking new frameworks (e.g., `CoreGraphics`, `QuartzCore`), ensure `scripts/build-ios.sh` and `scripts/build-ios-device.sh` include `-framework <Name>`.

### Step 4: Android Views Renderer (`platform/android/`)
1. Update node constants in `platform/android/src/dev/gonative/counter/MainActivity.java`.
2. Update view creation, styling, and event listener attachment.
3. If new JNI native methods are needed, mirror them in `examples/counter/androidbridge/jni.c` and `examples/counter/androidbridge/main.go`.

### Step 5: Verification
1. Run `go test -race ./...` and `go vet ./...`.
2. Build iOS: `go run ./cmd/gonative build ios`.
3. Build Android: `go run ./cmd/gonative build android`.

---

## 4. Runbook: Standalone Project Scaffolding

Running `gonative init <name>` creates a full React Native-style standalone project directory containing:
- `app.go`: Pure Go declarative UI component.
- `ios/`: Native UIKit host (`main.m`, `GoNativeRenderer.m/.h`, `Info.plist`, `bridge/main.go`).
- `android/`: Native Android Views host (`AndroidManifest.xml`, Gradle project, `MainActivity.java`, `GapDrawable.java`, `bridge/main.go`, `bridge/jni.c`).

To create and run a new project:
```bash
gonative init my-new-app
cd my-new-app
gonative doctor
gonative run ios
gonative run android
```

---

## 5. Runbook: Diagnostics & Inspector

The inspector provides read-only HTTP endpoints on loopback:
- `GET /v1/tree`: Returns the JSON snapshot of the virtual UI tree.
- `GET /v1/logs`: Returns bounded chronological diagnostic events (`runtime.started`, `render.batch_applied`, `event.dispatched`, etc.).

To start the inspector in an application or test:
```go
import "github.com/go-native/go-native/runtime/inspector"

srv := inspector.New(appRuntime, "127.0.0.1:8080")
addr, err := srv.Start()
defer srv.Stop(ctx)
```

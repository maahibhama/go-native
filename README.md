# go-native

An experimental declarative Go UI runtime that renders genuine platform controls. The counter renders with UIKit `UILabel`/`UIButton` on iOS and Android `TextView`/`Button`—no WebView, JavaScript, HTML, React, canvas, Flutter, or Skia.

## Requirements

- macOS on Apple Silicon
- Go 1.24 or newer
- Xcode with an iOS Simulator runtime
- Xcode command-line tools selected (`xcode-select -p`)
- For Android: Android SDK platform 35, build-tools 36.0.0, NDK 28.2.13676358, and JDK 17

## Test

```bash
GOCACHE=/tmp/go-native-gocache go test ./...
GOCACHE=/tmp/go-native-gocache go vet ./...
```

## Build the iOS counter

```bash
./scripts/build-ios.sh
```

The signed simulator bundle is written to `build/ios-simulator/GoNativeCounter.app`.

## Run the iOS counter

List available devices if the default `iPhone 17 Pro` is absent:

```bash
xcrun simctl list devices available
```

Then run:

```bash
./scripts/run-ios.sh
```

To select another device:

```bash
GONATIVE_SIMULATOR="iPhone 16 Pro" ./scripts/run-ios.sh
```

Tap **Increment**. UIKit sends a numeric handler ID to Go, Go updates state and rebuilds the tree, reconciliation emits one label update, and UIKit changes the existing `UILabel` to `Count: 1`.

## Build and run the Android counter

Set `ANDROID_SDK_ROOT` if your SDK is not at `~/Library/Android/sdk`, then build:

```bash
./scripts/build-android.sh
```

The signed debug APK is written to `build/android/GoNativeCounter.apk`. With an emulator or arm64 device running:

```bash
./scripts/run-android.sh
```

If multiple devices are attached, select one explicitly:

```bash
GONATIVE_ANDROID_SERIAL=emulator-5554 ./scripts/run-android.sh
```

The Java host uses `TextView`, `Button`, and `LinearLayout`. JNI carries one binary mutation batch per render and only integer handler IDs on the callback path.

## Repository map

- `ui`: typed declarative primitives and state
- `runtime`: events, identity, reconciliation, scheduler, and binary mutation protocol
- `platform/ios`: Objective-C UIKit renderer and host
- `platform/android`: Java Android Views renderer and host
- `examples/counter`: all-Go counter app plus the narrow cgo entry bridge
- `docs`: architecture and decision records

The current milestone intentionally supports only View, Column, Row, Text, and Button. See [docs/architecture.md](docs/architecture.md) for ownership, threading, lifecycle, and deferred scope.

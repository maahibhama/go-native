# go-native

An experimental declarative Go UI runtime that renders genuine platform controls. Milestone 0 renders a counter with UIKit `UILabel`, `UIButton`, and `UIStackView`—no WebView, JavaScript, HTML, React, canvas, Flutter, or Skia.

## Requirements

- macOS on Apple Silicon
- Go 1.24 or newer
- Xcode with an iOS Simulator runtime
- Xcode command-line tools selected (`xcode-select -p`)

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

## Repository map

- `ui`: typed declarative primitives and state
- `runtime`: events, identity, reconciliation, scheduler, and binary mutation protocol
- `platform/ios`: Objective-C UIKit renderer and host
- `examples/counter`: all-Go counter app plus the narrow cgo entry bridge
- `docs`: architecture and decision records

The current milestone intentionally supports only View, Column, Row, Text, and Button. See [docs/architecture.md](docs/architecture.md) for ownership, threading, lifecycle, and deferred scope.

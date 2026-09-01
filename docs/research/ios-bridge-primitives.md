# iOS bridge research notes

Checked on 2026-09-02 before locking the Milestone 0 bridge:

- Go's documented [`c-archive` build mode](https://pkg.go.dev/cmd/go#hdr-Build_modes) builds one main package and its dependencies into a C archive; only cgo-exported functions are callable. This matches the desired tiny entry surface.
- The official [cgo documentation](https://pkg.go.dev/cmd/cgo) permits C to call exported Go functions and states that C may retain a Go pointer only while it remains pinned. The bridge therefore copies batch bytes during the C call and retains only integer node/handler IDs.
- [`gomobile bind`](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile) can generate an Apple XCFramework and Objective-C bindings. It remains a viable packaging option, but its broader generated binding layer is unnecessary for the POC's three-function ABI.
- The portable [`x/mobile/app`](https://pkg.go.dev/golang.org/x/mobile/app) path focuses on portable APIs such as OpenGL and events. This project instead needs direct UIKit ownership in a native host.
- Apple states that [UIKit UI work must run on the main thread](https://developer.apple.com/documentation/technologyoverviews/uikit-appkit). The Objective-C renderer copies each payload and dispatches mutation application to the main queue.

The Android shared-library/JNI path is intentionally not selected or implemented until the complete iOS loop is stable, per Milestone 0 scope.

# App code cleanup and native framework packaging

Status: completed  
Branch: `feat/v1/app_code_clean`

## Implementation progress

- [x] Define the native framework packaging decision.
- [x] Add the app-neutral iOS ABI header and remove the framework renderer's
  dependency on an app-generated header.
- [x] Extract bounded little-endian reading into `GNProtocolReader` and compile it
  through framework simulator/device build paths.
- [x] Extract the iOS integer-keyed native registry and lifecycle runtime host.
- [x] Extract iOS native measurement into `GNMeasurementHost`.
- [x] Split iOS control creation and styling.
- [x] Extract Android bounded protocol reading into an app-neutral runtime package.
- [x] Move Android renderer ownership behind the app-neutral `GoNativeActivity` runtime surface.
- [x] Split Android registry, controls, measurement, and runtime host internally.
- [x] Build a reproducible `GoNativeKit.xcframework` with device and simulator slices.
- [x] Build `gonative-runtime.aar` and consume the runtime module from the framework fixture.
- [x] Consume the packaged AAR from generated standalone applications.
- [x] Reduce generated projects and regenerate `examples/showcase-app` as a thin consumer.
- [x] Standardize iOS `AppDelegate` (`AppDelegate.h`, `AppDelegate.m`, `main.m`) and Android `MainActivity` host launchers mirroring React Native architecture.

## Outcome

A generated Go Native application contains application code, assets, configuration,
and minimal platform launchers. UIKit rendering, Android Views rendering, protocol
decoding, view registries, measurement, lifecycle forwarding, and event forwarding
are versioned framework dependencies maintained in this repository.

```text
application Go code
        │ compiled per application
        ▼
Go application library (.a / .so)
        │ stable app ABI
        ▼
GoNativeKit.xcframework / gonative-runtime.aar
        │ mutation batches and integer events
        ▼
UIKit / Android Views
```

The migration must preserve one binary mutation batch per render, integer-only
native references, copied payload ownership, and UI-thread-only native mutations.

## Desired generated application

```text
my-app/
├── app.go
├── go.mod
├── gonative.yaml
├── assets/
├── ios/
│   ├── AppDelegate.h
│   ├── AppDelegate.m
│   ├── main.m
│   ├── Info.plist
│   └── MyApp.xcodeproj
└── android/
    ├── settings.gradle
    └── app/
        ├── build.gradle
        └── src/main/
            ├── AndroidManifest.xml
            └── java/.../MainActivity.java
```

The generated project must not contain a mutation decoder, native view registry,
control factory, renderer implementation, or copied framework tests.

## Package boundaries

### Go application ABI

`platform/abi/` owns the C-level contract between an app-specific Go library and
the native framework. The ABI uses fixed-width values, byte buffers, and integer
IDs only. It is independent of the app name and native package name.

The ABI is not the mutation protocol. Both are versioned separately:

- app ABI version: native framework ↔ compiled Go application entrypoints;
- mutation protocol version: serialized render batch;
- measurement protocol version: native intrinsic measurement request/response.

### iOS framework

Create `platform/ios/GoNativeKit/` with focused sources:

- `GNRuntimeHost` — lifecycle and calls into the app ABI;
- `GNMutationDecoder` — bounded little-endian decoding;
- `GNRenderer` — ordered mutation application;
- `GNViewRegistry` — `NodeID` to `UIView` ownership;
- `GNEventDispatcher` — integer handler event forwarding;
- `GNMeasurementHost` — intrinsic control measurement;
- `controls/` — control creation and prop application by family.

Build one `GoNativeKit.xcframework` containing device and simulator slices. Add a
Swift Package manifest that exposes the binary artifact. The app-specific Go
archive remains linked by the generated application target.

### Android library

Create `platform/android/runtime/` as `com.android.library`:

- `GoNativeActivity` and `GoNativeHost` — minimal integration surface;
- `MutationDecoder` — bounded little-endian decoding;
- `NativeRenderer` and `ViewRegistry` — view lifecycle;
- `EventDispatcher` — calls the app ABI through JNI;
- `NativeMeasurer` — batched intrinsic measurement;
- `controls/` — control creation and props by family.

Publish `gonative-runtime.aar` locally first and later to a Maven repository. Use
`JNI_OnLoad` plus `RegisterNatives` so JNI registration is independent of the
developer's Java/Kotlin package name. The app-specific `libgonative.so` remains in
the application and is loaded by the runtime library.

## Delivery phases

### Phase 0 — Baseline and ownership

- Preserve the current passing Go and native builds.
- Record renderer/template/example equivalence checks.
- Remove build products and machine files from repository views.
- Add the app-neutral ABI header and an ADR for packaging boundaries.

Acceptance:

- Current examples still build and render unchanged.
- The framework renderer no longer imports an app-named generated header.
- No wire-visible protocol field changes.

### Phase 1 — Split native monoliths

- Split the iOS renderer by decoder, registry, controls, measurement, events, and
  runtime host.
- Split Android `MainActivity` by the same responsibilities.
- Keep behavior byte-for-byte compatible with protocol v10 and measurement v2.
- Add focused decoder and registry contract tests.

Acceptance:

- No production native source file exceeds an agreed maintainability limit.
- Framework native builds and existing examples pass.
- Decoder failures are structured and bounded.

### Phase 2 — Produce native libraries

- Add reproducible XCFramework build scripts.
- Convert Android renderer into an AAR library module.
- Define debug/release artifacts and version metadata.
- Verify minimum iOS 15 and Android API 23 targets.

Acceptance:

- A fixture application links prebuilt platform libraries.
- Platform implementation sources are absent from the fixture.
- Library/runtime/protocol incompatibility fails with a useful diagnostic.

### Phase 3 — Thin generated projects

- Replace native renderer blobs in `cmd/gonative/templates.go` with small launchers.
- Generate a manifest containing framework and protocol compatibility requirements.
- Teach `gonative init`, `build`, and `run` to resolve or build native artifacts.
- Regenerate `examples/showcase-app` only through the generator.

Acceptance:

- `examples/showcase-app` consumes the frameworks and contains no renderer copy.
- `gonative init` output is small and understandable.
- Xcode and Android Studio builds work without framework repository paths.

### Phase 4 — Dependency and module integration

- [x] Add local SPM/XCFramework and Maven/AAR dependency resolution.
- Generate permission, plist, manifest, framework, and Gradle declarations from
  `gonative.yaml`.
- Introduce typed native-module registration without modifying app launchers.

Acceptance:

- Adding a supported module changes the manifest and generated build metadata,
  not renderer source.
- Missing capabilities report explicit runtime or build errors.

### Phase 5 — Release engineering

- Build signed checksummed framework artifacts in CI.
- Add compatibility matrices and upgrade/migration tests.
- Cache artifacts for incremental local builds.
- Publish an example app created from a clean release installation.

Acceptance:

- Framework artifacts are reproducible and traceable to a release.
- Clean-machine iOS and Android builds pass.
- A framework rollback does not require editing application source.

## CLI experience

```text
gonative new my-app
gonative add camera
gonative generate
gonative run ios
gonative run android
```

The CLI owns generated build metadata. Developers own `app.go`, assets, application
configuration, and explicitly marked platform extension points. Generated files
carry a header and are replaced deterministically; user files are never overwritten.

## Testing gates

- Go race tests, vet, formatting, and protocol goldens;
- Objective-C decoder/control contract tests on iOS;
- Android instrumentation tests for decoder/control contracts;
- generated-project snapshot tests;
- build tests using only packaged XCFramework/AAR artifacts;
- screenshot parity tests for the component fixture;
- startup, mutation throughput, memory, and artifact-size regression gates.

## Migration policy

During v0, land the work in compatibility-preserving slices. Keep the old source
integration available behind an internal fallback until both packaged fixture builds
pass. Remove copied renderers only after `gonative init`, framework examples, and
`examples/showcase-app` use the packaged path. Document every generated-project change
in `docs/migrations/`.

## Completion definition

This cleanup is complete when native renderer implementation exists in exactly one
framework-owned location per platform, templates contain only launch/configuration
code, `examples/showcase-app` is an ordinary consumer, and both platforms communicate
with compiled Go application code through the versioned app-neutral ABI.

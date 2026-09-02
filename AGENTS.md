# Go Native — Agent Context

## Mission

Go Native is a declarative Go UI runtime that renders real UIKit controls on iOS and Android Views on Android. Preserve these architectural boundaries:

- no WebView, JavaScript runtime, JSON bridge, Flutter engine, or canvas renderer;
- one versioned little-endian `runtime.MutationBatch` crosses the native boundary per render pass;
- native code stores integer `NodeID` and `HandlerID` values only—never Go pointers;
- native view mutations run exclusively on the platform UI thread.

For implementation work in this repository, load the project skill at [`.agents/skills/go-native/SKILL.md`](.agents/skills/go-native/SKILL.md).

## Source-of-truth map

- `ui/node.go`: node types, IDs, `Props`, component contract, explicit identity.
- `ui/components.go`: public primitives and fluent modifiers.
- `ui/state.go`: goroutine-safe state and the process-wide render scheduler.
- `ui/intents.go`: gesture and animation contracts.
- `ui/presentation.go`: navigation/modal contracts; these currently expose metadata and fallback content, not full native mounting.
- `runtime/reconciler.go`: mutation ordering and keyed/unkeyed identity behavior.
- `runtime/runtime.go`: scheduling, handler binding/release, diagnostics, timing, and renderer calls.
- `runtime/protocol.go`: canonical binary wire layout and protocol version.
- `runtime/interactions.go`: nested gesture/animation payload embedded in `Props.Interactions`.
- `platform/ios/GoNativeRenderer.m`: UIKit decoder and renderer.
- `platform/android/src/dev/gonative/counter/MainActivity.java`: Android decoder and renderer.
- `cmd/gonative/main.go`: CLI dispatch, toolchain doctor, standalone builds.
- `cmd/gonative/templates.go`: generated standalone project sources. Changes to platform bridges or renderers often need matching template changes.
- `examples/counter/`: framework demo bridge.
- `examples/my-project/`: checked-in generated-project fixture; keep it aligned with scaffolding when generator behavior changes.
- `runtime/inspector/`: loopback-only, read-only diagnostic HTTP service.

Narrative documentation can lag implementation. Resolve discrepancies in this order: tests and live source, ADRs, then overview/roadmap prose. In particular, `docs/architecture.md` and `docs/roadmap.md` currently understate implemented gesture/animation/diagnostic contracts.

## Cross-layer invariants

### Protocol changes

Any wire-visible change—including `ui.NodeType`, mutation values, `ui.Props`, interaction payloads, or field order—must be reviewed across:

1. `ui/` declarations and tests;
2. `runtime/protocol.go`, `runtime/interactions.go`, and protocol/runtime tests;
3. iOS and Android framework renderers;
4. `cmd/gonative/templates.go` generated renderer/bridge text;
5. checked-in generated example renderers under `examples/my-project/`.

The current protocol version is `7`. Native decoders currently compare literal `7` values, so search for the old version before bumping it. Preserve exact field order, byte widths, signedness, little-endian encoding, and length-prefix handling.

### Identity and handlers

- Unkeyed identity is stabilized by type and structural position.
- `ui.WithID` provides explicit identity for reorderable children.
- Reuse handler IDs when logical nodes survive; replace their Go callbacks in the registry.
- Release action, value, boolean, and gesture handlers when nodes disappear, are replaced, or the runtime stops.
- Never retain Go pointers or closures in Objective-C/Java state.

### Threading and ownership

- `ui.State[T]` may be changed from any goroutine.
- `runtime.Runtime.Schedule` coalesces renders; reconciliation and tree ownership remain serialized.
- iOS copies bytes before dispatching to the main queue.
- Android clones the payload before `runOnUiThread`.
- Native teardown must stop the Go runtime and clear native registries/references.

## Change discipline

- Inspect `git status --short` before editing. Preserve unrelated user changes.
- Keep architecture changes synchronized across Go, iOS, Android, templates, examples, tests, and relevant ADR/docs.
- Add focused tests beside the owning package. Protocol changes require round-trip and renderer-alignment coverage where practical.
- Do not claim native behavior is complete from Go-only tests; build or device verification is required when platform behavior changes.

## Verification

Run the smallest relevant checks while iterating, then the applicable gates:

```bash
go test -race ./...
go vet ./...
gofmt -l .
make benchmark
go run ./cmd/gonative doctor
```

For native changes also run the relevant builds:

```bash
go run ./cmd/gonative build ios
go run ./cmd/gonative build android
```

Native benchmark commands are `go run ./cmd/gonative benchmark native ios` and `... android`. Toolchain-dependent failures should be reported separately from code/test failures.

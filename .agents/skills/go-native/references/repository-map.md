# Repository ownership map

Read only the sections relevant to the current task.

## Portable UI

- `ui/node.go`: wire-facing node/prop types, identity, component construction.
- `ui/components.go`: supported native primitives and modifiers.
- `ui/state.go`: state synchronization and scheduler hook.
- `ui/intents.go`: gesture and animation metadata.
- `ui/presentation.go`: navigation and modal metadata contracts.

Public API additions need package tests and an assessment of whether they change the binary protocol.

## Runtime

- `runtime/reconciler.go`: create/delete/update/insert/remove/move ordering.
- `runtime/runtime.go`: rendering lifecycle, ID stabilization, handler reuse/release, diagnostics and timing.
- `runtime/events.go`: typed callback registries.
- `runtime/protocol.go`: outer mutation-batch wire format, currently version 7.
- `runtime/interactions.go`: inner gesture/animation wire payload.
- `runtime/diagnostics.go`: bounded structured logs and detached snapshots.
- `runtime/inspector/service.go`: loopback HTTP exposure of logs/tree.

## Native hosts

- iOS: `platform/ios/GoNativeRenderer.m`, header, and `main.m`.
- Android: `platform/android/src/dev/gonative/counter/MainActivity.java`, `GapDrawable.java`, manifests, and Gradle files.
- Bridge entrypoints: `examples/counter/bridge/` and `examples/counter/androidbridge/`.

Native renderers own view construction, styling, accessibility mapping, interaction recognizers/animators, mutation application, and UI-thread delivery. They must not own Go application state.

## Tooling and generated projects

- `cmd/gonative/main.go`: CLI parsing, doctor, build/run routing, native benchmark routing, and standalone project operations.
- `cmd/gonative/templates.go`: embedded scaffold sources for iOS, Android, Go bridges, Gradle, and Xcode.
- `cmd/gonative/main_test.go`: generator and command-routing expectations.
- `examples/my-project/`: checked-in generated output used as a practical fixture; it may contain user changes, so inspect diffs before syncing.
- `scripts/`: framework-repository build/run and native benchmark automation.

## Documentation

- `docs/adr/`: architectural decisions; update/add an ADR for substantial architecture changes.
- `docs/performance.md`: benchmark baselines.
- `docs/diagnostics.md` and `docs/inspector.md`: observability behavior.
- `docs/roadmap.md` and `docs/architecture.md`: useful intent, but verify status claims against live code and tests.

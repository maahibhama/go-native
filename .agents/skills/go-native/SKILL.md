---
name: go-native
description: Develop, review, debug, or extend the Go Native declarative UI framework across its Go runtime, binary protocol, UIKit renderer, Android Views renderer, CLI scaffolding, examples, diagnostics, and tests. Use for work inside the go-native framework or generated Go Native applications; do not use for unrelated Go or generic mobile projects.
---

# Go Native repository workflow

Work from the repository's live source and tests. Read the root `AGENTS.md` first; it defines the project boundaries, source-of-truth map, synchronization surfaces, and verification gates.

## Route the task

- For UI primitives or props, read `ui/node.go`, `ui/components.go`, and the relevant UI tests.
- For reconciliation, identity, scheduling, events, or diagnostics, read the owning files in `runtime/` and their tests.
- For any wire-visible change, read [references/protocol-checklist.md](references/protocol-checklist.md) before editing.
- For iOS/Android rendering, inspect both native implementations even if only one platform is requested, because they decode the same protocol.
- For CLI initialization or standalone builds, inspect `cmd/gonative/main.go`, `cmd/gonative/templates.go`, CLI tests, and the checked-in `examples/my-project/` fixture.
- For inspector work, preserve loopback binding, read-only routes, bounded diagnostics, and detached tree snapshots.
- For navigation/modal work, distinguish the current metadata/fallback contracts in `ui/presentation.go` from native controller mounting, which is not yet integrated.

Use [references/repository-map.md](references/repository-map.md) when the task crosses packages or the correct ownership location is unclear.

## Implement with cross-platform consistency

Keep changes in the layer that owns the behavior. Avoid putting platform policy into portable `ui` types or application behavior into native renderers.

When a feature crosses the bridge, update Go encoding and decoding, both native decoders/renderers, generated templates, the checked-in generated example, and focused tests as one change. Preserve one batch per render pass, integer-only foreign references, payload ownership before asynchronous dispatch, and UI-thread-only view mutation.

Preserve stable IDs and handler lifecycle. Keyed children must keep logical identity across reorders; surviving handlers should keep their IDs while callbacks are replaced; removed subtrees must release every handler kind.

## Validate proportionally

During iteration, run focused package tests. Before handing off ordinary Go changes, run:

```bash
go test -race ./...
go vet ./...
gofmt -l .
```

Run `make benchmark` for reconciliation, serialization, event-registry, or performance-sensitive changes. Run `go run ./cmd/gonative doctor` when diagnosing local toolchains. Build each affected native platform for renderer, bridge, template, or build-workflow changes; state clearly when device or toolchain verification was unavailable.

Inspect `git diff --check` and `git status --short` before completion. Do not overwrite unrelated work in this frequently modified cross-platform tree.

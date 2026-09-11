# Code tour

Start here when changing Go Native. The repository has four layers; a change should have one clear owner.

```text
Application code
    ↓ builds
ui.Component / ui.Node tree
    ↓ reconciled, laid out, encoded by
runtime.Runtime
    ↓ one little-endian MutationBatch per render
native host (UIKit or Android Views)
    ↓ events return as integer HandlerID values
runtime.EventRegistry → Go callback → state update
```

## The four layers

### 1. Portable UI (`ui/`)

This is the public API used by applications. It must not import platform code.

- `node.go` — the small wire-facing node model, IDs, primitive kinds, and callbacks.
- `components.go` — primitive constructors and modifiers shared by components.
- `design.go` and `context.go` — styles, themes, tokens, environment, and media query.
- `layout_components.go` — layout compositions such as `Stack`, `Center`, and keyboard avoidance.
- `text_input.go`, `pressable.go`, `selection_controls.go`, `media.go` — control families.
- `hooks.go`, `state.go`, and `focus.go` — stateful application behavior.
- `presentation.go` — navigation and modal descriptions. These are not full native routers yet.

Rule: a UI component describes intent. It does not perform native work or serialize itself.

### 2. Runtime (`runtime/`)

The runtime turns descriptions into ordered native mutations.

- `runtime.go` — application lifecycle and the render transaction.
- `reconciler.go` — tree identity and create/update/move/delete ordering.
- `events.go` — typed callback registries keyed only by integer handler IDs.
- `protocol.go` — canonical mutation-batch binary format.
- `style_protocol.go` and `interactions.go` — nested wire payloads.
- `layout/` — constraint layout, intrinsic measurement, and measurement protocol.
- `headless/` — deterministic renderer for tests.
- `diagnostics.go` and `inspector/` — bounded, read-only diagnostics.

The render transaction is:

1. Build a new portable tree.
2. Preserve stable IDs from the mounted tree.
3. Bind or release integer handler IDs.
4. Compute Go-owned layout boxes.
5. Reconcile old and new trees.
6. Encode and submit exactly one mutation batch.
7. Commit hooks and retain the new tree only after successful submission.

### 3. Native hosts (`platform/`)

- `platform/ios/GoNativeRenderer.m` decodes the batch and owns UIKit objects.
- `platform/android/.../MainActivity.java` decodes the same batch and owns Android Views.

Native hosts create controls, apply properties and geometry, deliver events, and clean up native objects. They do not own application state. All mutations execute on the platform UI thread, and native registries store only `NodeID` and `HandlerID` integers.

### 4. Tooling and applications

- `cmd/gonative/` — CLI and generated project templates.
- `examples/counter/` — framework development bridge.
- `examples/my-project/` — checked-in generated application fixture.

`cmd/gonative/templates.go` is generated-project source embedded as Go strings. Native protocol changes therefore have four synchronization surfaces: framework iOS, framework Android, templates, and `examples/my-project`.

## Where a change belongs

| Change | Primary owner | Other required checks |
|---|---|---|
| New component composed from existing nodes | `ui/` component-family file | UI tests and example |
| New native primitive or wire property | `ui/node.go` + `runtime/` protocol | Both native hosts, templates, fixture, protocol tests |
| Rebuild, identity, or handler behavior | `runtime/` | Race tests and benchmarks |
| Layout algorithm | `runtime/layout/` | Headless/layout tests and native measurement compatibility |
| UIKit or Android visual mapping | Both `platform/` hosts | Templates and generated fixture |
| CLI/scaffold behavior | `cmd/gonative/` | `examples/my-project/` and CLI tests |

## Invariants to check before merging

- One `MutationBatch` crosses the native boundary per committed render.
- The binary protocol is versioned, bounded, little-endian, and decoded in identical field order.
- Native code never retains Go pointers or callbacks.
- Removed or replaced nodes release every handler kind.
- Reorderable children use explicit stable IDs.
- Native mutation application is UI-thread-only.
- Framework renderers, templates, and generated examples stay synchronized.

## Suggested reading order

For the shortest useful path through the implementation:

1. `examples/my-project/app.go`
2. `ui/node.go`
3. `runtime/runtime.go` (`render` method)
4. `runtime/reconciler.go`
5. `runtime/protocol.go`
6. One native renderer
7. `cmd/gonative/templates.go` only when working on scaffolding

Use `AGENTS.md` as the detailed source-of-truth map and cross-layer checklist.

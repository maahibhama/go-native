# Typed design foundation

Go Native v0.2 introduces the production design model without changing mutation protocol version 7.

## Public model

- `ui.Style` groups layout, appearance, text, and interaction properties.
- `ui.Theme` owns semantic colors, typography, spacing, radii, elevations, motion, icon sizes, control sizes, and component variants.
- `ui.Token[T]` resolves semantic values with deterministic fallbacks. Use `ColorToken`, `SpacingToken`, `TypographyToken`, and `StyleToken` to create typed tokens.
- `ui.Environment` carries theme, locale, layout direction, media query, lifecycle, and typed dependencies.
- `ui.BuildContext` is immutable and scoped by mounted component path.
- `runtime/layout.Engine` computes deterministic box geometry and accepts a batched/intrinsic `Measurer` implementation.
- `runtime/headless.Renderer` applies mutation batches in memory for component and integration tests.

## Protocol v7 compatibility

`Styled` retains the complete typed style on `ui.Node`. Width, height, uniform padding, gap, alignment, font size, and bold weight are projected into the existing `ui.Props` wire record. Other typed fields are available to layout, inspection, and tests but are not yet sent to native renderers.

This split is temporary. The typed protocol migration must update Go encoding, both native decoders, generated templates, examples, and golden fixtures together.

Existing modifiers such as `Width`, `Padding`, and `FontSize` now populate both typed style and legacy props. Existing applications therefore remain source- and wire-compatible.

## Context-aware applications

Legacy applications continue to use `runtime.New(func() ui.Component, renderer)`. New applications can use `runtime.NewContext(func(ui.BuildContext) ui.Component, renderer, environment)`.

Runtime environment updates are immutable replacements and schedule a coalesced render. Lifecycle changes use `Runtime.SetLifecycle`. Components implementing `ui.ContextComponent` receive context even when hosted by a legacy component tree.

## Mounted components and hooks

`ui.Functional(key, render)` creates a mounted component scope. The key and its structural parent path define hook identity, so keys must remain stable among siblings. Child components are built lazily with their inherited `BuildContext`; keyed wrappers, interaction decorators, navigation fallbacks, and modal bases preserve that context.

The initial hook set is `UseState`, `UseReducer`, `UseRef`, `UseMemo`, `UseCallback`, `UseEffect`, `UseLayoutEffect`, `UseLifecycle`, and `UseMediaQuery`. Go requires the active `BuildContext` as the first hook argument. Hooks must be called unconditionally and in the same order on every render.

State created by a hook schedules only its owning runtime. Effects run after a successful renderer commit, with layout effects first. A changed dependency cancels the prior effect context and runs cleanup before replacement. Removed mounted scopes and `Runtime.Stop` cancel and clean up all effects deterministically.

## Current layout coverage

The Go layout engine supports intrinsic leaf measurement, logical-point and percentage dimensions, minimum and maximum constraints, padding, margins, gaps, row and column flow, cross-axis alignment, absolute overlays, flex grow/shrink/basis, wrapping, aspect ratios, fixed and adaptive grids, and RTL mirroring.

Use `ui.Grid(columns, children...)` for fixed tracks or `AdaptiveGrid(minColumnWidth)` for viewport-derived columns. `ui.ResponsiveStyle` applies ordered minimum-width `ui.Breakpoint` overrides using `MediaQuery.Viewport`. `runtime/layout.Engine.Direction` controls logical RTL placement without changing the component tree.

These advanced fields remain Go-owned metadata under protocol v7. Native measurement batching and the typed layout mutation protocol are separate v0.2 milestones.

## Batched intrinsic measurement

`runtime/layout.Engine.LayoutMeasured` collects every uncached intrinsic leaf into one `BatchMeasurer` request before computing geometry. Requests contain value-only node type, text/image content, complete typed style, and constraints. Results are matched by integer request ID and rejected when missing, duplicated, unknown, or returned with a structured native error.

`MeasurementCache` keys results by content, style, and constraints and is safe for concurrent access. Hosts must invalidate or replace the cache when native font or asset availability changes. Protocol capability negotiation and bounded payload, mutation-count, and string limits are available in `runtime`.

Protocol v8 embeds the nested record implemented by `runtime.MarshalTypedStyles` and `UnmarshalTypedStyles` in every mutation. It serializes the portable style followed by complete iOS and Android overrides using declaration-ordered, fixed-width little-endian fields. The record has its own version, strict string and trailing-data validation, round-trip coverage, and a stable SHA-256 golden fixture. Both native readers validate and consume the bounded record; field-level native application is tracked separately.

UIKit and Android Views currently apply the portable appearance shell from that record: background and foreground RGBA colors, border width/color, corner radius, opacity, visibility, and disabled interaction. Platform-override merging, complete typography, transforms/shadows, and Go-computed layout geometry remain follow-up decoder work.

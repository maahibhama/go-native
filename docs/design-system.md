# Typed design foundation

Go Native v0.2 introduces the production design model and mutation protocol version 9.

## Public model

- `ui.Style` groups layout, appearance, text, and interaction properties.
- `ui.Theme` owns semantic colors, typography, spacing, radii, elevations, motion, icon sizes, control sizes, and component variants.
- `ui.Token[T]` resolves semantic values with deterministic fallbacks. Use `ColorToken`, `SpacingToken`, `TypographyToken`, and `StyleToken` to create typed tokens.
- `ui.Environment` carries theme, locale, layout direction, media query, lifecycle, and typed dependencies.
- `ui.BuildContext` is immutable and scoped by mounted component path.
- `runtime/layout.Pipeline` batches native intrinsic measurement, computes deterministic box geometry, and flattens it into integer-keyed frames.
- `runtime/headless.Renderer` applies mutation batches in memory for component and integration tests.

## Legacy API compatibility

`Styled` retains the complete typed style on `ui.Node`. Width, height, uniform padding, gap, alignment, font size, and bold weight are also projected into legacy `ui.Props`, while protocol v9 carries the full typed record and computed geometry to native renderers.

Existing modifiers such as `Width`, `Padding`, and `FontSize` populate both typed style and legacy props, so application source remains compatible. Custom native hosts are not wire-compatible with earlier protocols and must migrate their decoder as described in `docs/migrations/v0.2.md`.

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

Advanced layout fields are resolved by the Go-owned layout pipeline before each commit. Hosts receive the resulting logical-point frame in the same mutation record while native controls retain text input, focus, accessibility, selection, and scrolling behavior.

## Batched intrinsic measurement

`runtime/layout.Engine.LayoutMeasured` collects every uncached intrinsic leaf into one `BatchMeasurer` request before computing geometry. Requests contain value-only node type, text/image content, complete typed style, and constraints. Results are matched by integer request ID and rejected when missing, duplicated, unknown, or returned with a structured native error.

`MarshalMeasurementRequests`/`UnmarshalMeasurementRequests` and their result counterparts define the native adapter wire format. The bounded little-endian protocol carries a versioned batch header, integer request IDs, node type, constraints, content, typed style, measured size, and structured error text. A golden request fixture protects field ordering across Objective-C/JNI implementations.

`MeasurementCache` keys results by content, style, and constraints and is safe for concurrent access. Hosts must invalidate or replace the cache when native font or asset availability changes. Protocol capability negotiation and bounded payload, mutation-count, and string limits are available in `runtime`.

Protocol v9 embeds the nested record implemented by `runtime.MarshalTypedStyles` and `UnmarshalTypedStyles` in every mutation. It serializes the portable style followed by complete iOS and Android overrides using declaration-ordered, fixed-width little-endian fields, then an optional computed frame. The style record has its own version, strict string and trailing-data validation, round-trip coverage, and a stable SHA-256 golden fixture.

UIKit and Android Views apply the appearance and typography shell from that record: background and foreground RGBA colors, border width/color, corner radius, opacity, visibility, disabled interaction, font family/size/weight, line height, letter spacing, translation, scale, rotation, and native shadows/elevation. Each host resolves its own platform override after portable style and applies guarded computed geometry without replacing the full-screen root host.

## Focus and lifecycle

`FocusManager`, `FocusScope`, `FocusNode`, and `UseFocusNode` provide application-scoped focus identity, traversal, observation, programmatic requests, and deterministic mounted cleanup. Controls associate a node with `WithFocusNode`; native focus changes return only the integer `NodeID`, and programmatic requests are mirrored through the normal controlled `Focused` prop.

Native foreground, active, inactive, background, memory-pressure, and destroyed callbacks update the runtime environment. `Runtime.ObserveLifecycle` supports application services independently of rendering, while `UseLifecycle` rebuilds mounted UI through context. Both subscriptions clean up deterministically.

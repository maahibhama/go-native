# Typed design foundation

Go Native v0.2 introduces the production design model. The v0.3 text/input slice uses mutation protocol version 10.

## Public model

- `ui.Style` groups layout, appearance, text, and interaction properties.
- `ui.Theme` owns semantic colors, typography, spacing, radii, elevations, motion, icon sizes, control sizes, and component variants.
- `ui.Token[T]` resolves semantic values with deterministic fallbacks. Use `ColorToken`, `SpacingToken`, `TypographyToken`, and `StyleToken` to create typed tokens.
- `ui.Environment` carries theme, locale, layout direction, media query, lifecycle, and typed dependencies.
- `ui.BuildContext` is immutable and scoped by mounted component path.
- `runtime/layout.Pipeline` batches native intrinsic measurement, computes deterministic box geometry, and flattens it into integer-keyed frames.
- `runtime/headless.Renderer` applies mutation batches in memory for component and integration tests.

## Legacy API compatibility

`Styled` retains the complete typed style on `ui.Node`. Width, height, uniform padding, gap, alignment, font size, and bold weight are also projected into legacy `ui.Props`, while protocol v10 carries the full typed record, computed geometry, and text/input contracts to native renderers.

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

The v0.3 composition layer adds `Stack`, flexible and fixed spacers, `Center`, `AspectRatio`, `KeyboardAvoidingView`, and horizontal or vertical dividers. These components build ordinary portable nodes rather than introducing platform-specific containers. The fluent API exposes min/max constraints, margins, axis-specific padding, absolute insets, overflow, borders, shadows, transforms, visibility, detailed typography, and hit slop through the existing typed style record.

## Batched intrinsic measurement

`runtime/layout.Engine.LayoutMeasured` collects every uncached intrinsic leaf into one `BatchMeasurer` request before computing geometry. Requests contain value-only node type, text/image content, complete typed style, and constraints. Results are matched by integer request ID and rejected when missing, duplicated, unknown, or returned with a structured native error.

`MarshalMeasurementRequests`/`UnmarshalMeasurementRequests` and their result counterparts define the native adapter wire format. Measurement protocol v2 carries a bounded little-endian batch header, integer request IDs, node type, constraints, content, geometry-affecting text/input properties, typed style, measured size, and structured error text. A golden request fixture protects field ordering across Objective-C/JNI implementations. The cache key includes wrapping, truncation, rich spans, placeholder, keyboard kind, secure/multiline state, and maximum length so configuration changes cannot reuse stale geometry.

`MeasurementCache` keys results by content, style, and constraints and is safe for concurrent access. Hosts must invalidate or replace the cache when native font or asset availability changes. Protocol capability negotiation and bounded payload, mutation-count, and string limits are available in `runtime`.

Protocol v10 embeds the nested record implemented by `runtime.MarshalTypedStyles` and `UnmarshalTypedStyles` in every mutation. It serializes the portable style followed by complete iOS and Android overrides using declaration-ordered, fixed-width little-endian fields, then an optional computed frame. The style record has its own version, strict string and trailing-data validation, round-trip coverage, and a stable SHA-256 golden fixture.

UIKit and Android Views apply the appearance and typography shell from that record: background and foreground RGBA colors, border width/color, corner radius, opacity, visibility, disabled interaction, font family/size/weight, line height, letter spacing, translation, scale, rotation, and native shadows/elevation. Each host resolves its own platform override after portable style and applies guarded computed geometry without replacing the full-screen root host.

## Focus and lifecycle

`FocusManager`, `FocusScope`, `FocusNode`, and `UseFocusNode` provide application-scoped focus identity, traversal, observation, programmatic requests, and deterministic mounted cleanup. Controls associate a node with `WithFocusNode`; native focus changes return only the integer `NodeID`, and programmatic requests are mirrored through the normal controlled `Focused` prop.

Native foreground, active, inactive, background, memory-pressure, and destroyed callbacks update the runtime environment. `Runtime.ObserveLifecycle` supports application services independently of rendering, while `UseLifecycle` rebuilds mounted UI through context. Both subscriptions clean up deterministically.

## Text and input

Text presentation is portable through `TextWrapping`, `TextOverflow`, `LineLimit`, `Selectable`, and the existing dynamic-type and typed-font modifiers. `RichText` accepts immutable `Span` values with optional font, color, underline, italic, and application-defined link identifiers. Rich spans cross the bridge in a bounded versioned value payload; link callbacks remain Go-owned and native code stores only their integer handler ID.

`TextInput` is controlled: its `Text` value is reapplied by committed renders. `UncontrolledTextInput` initializes the native value once and then leaves editing ownership native-side. Both expose placeholder, secure/multiline/read-only modes, text/email/phone/URL/integer/decimal/search keyboards, return actions, capitalization, autocorrection, maximum length, validation metadata, selection, focus nodes, and submit callbacks. Selection offsets use UTF-16 code units on both platforms; `UncontrolledTextSelection()` leaves cursor ownership native-side.

`InputFormatter` is a deterministic Go event-path transform. Formatters compose in declaration order before `onChange`. Controlled inputs should commit the formatted value to state so the next render authoritatively mirrors it; uncontrolled inputs receive the formatted callback value but retain native editing ownership.

Mutation protocol v10 appends the text/input record before the typed-style record. Strings and rich-span bytes share the 1 MiB per-field bound. Submit, link, selection, and value callbacks use stable integer IDs and are replaced in place for surviving nodes, then released on removal or shutdown.

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

The Go layout engine supports intrinsic leaf measurement, logical-point and percentage dimensions, minimum and maximum constraints, padding, margins, gaps, row and column flow, cross-axis alignment, and absolute overlays. Flex growth/shrink, wrapping, aspect ratio enforcement, grids, RTL mirroring, and native measurement batching remain subsequent v0.2 work.

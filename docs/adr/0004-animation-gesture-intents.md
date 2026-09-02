# ADR 0004: Portable animation and gesture intents

## Status

Accepted. Runtime/wire integration is complete; native execution is the next phase.

## Context

UIKit and Android expose substantially different recognizer and animation APIs.
Passing those objects into application code would make components platform-specific,
while encoding callbacks or durations directly in the existing compact props would
couple the first public API draft to the mutation protocol prematurely.

## Decision

The `ui` package defines protocol-neutral `GestureIntent` and `AnimationIntent`
values. `WithGesture` and `WithAnimation` decorate a component, and `IntentsOf`
returns detached metadata. Gesture handlers remain Go-owned, following the same
ownership rule as button handlers.

The first gesture vocabulary is tap, long press, directional swipe, and drag. The
first animation vocabulary covers opacity, scale, position, and layout with linear,
ease, and spring timing. Distances use logical points and timing uses `time.Duration`.
Leading and trailing directions are used instead of left and right so native
renderers can respect interface direction.

The API records whether an animation remains acceptable under reduced-motion
accessibility settings. Renderers must otherwise suppress or replace nonessential
motion when that setting is active.

## Consequences

This establishes typed application-facing contracts without changing the current
node or wire format. The decorations build the same node today and therefore do not
produce native behavior yet. The renderer integration phase must carry the metadata
onto built nodes, register gesture handlers as integer IDs, define protocol fields,
and implement platform recognizers and animations with cleanup on deletion.

Protocol v7 appends a length-prefixed interaction payload after every v6 field.
The payload retains arbitrary ordered gesture and animation counts. Gesture callback
IDs remain stable by node and gesture index, and callbacks carry translation and
velocity as four logical-point `float32` values. Keeping the payload in a Go string
makes `Props` comparable, so reconciliation remains allocation-light and deterministic.

Animations carry explicit scalar `From`/`To` targets for opacity and scale and
logical-point `FromX`/`FromY`/`ToX`/`ToY` targets for position. Opacity is clamped
to 0...1, scale is clamped to non-negative values, and non-finite values become
zero before crossing the bridge. Layout animation uses the surrounding native
layout pass rather than inventing numeric geometry.

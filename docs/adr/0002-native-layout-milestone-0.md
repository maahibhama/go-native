# ADR 0002: Translate Milestone 0 layout to native containers

- Status: Accepted for Milestone 0; revisit before complex layout
- Date: 2026-09-02

## Decision

Translate `Column` and `Row` to `UIStackView` and map gap, padding, alignment, width, and height to UIKit constraints and stack properties.

## Rationale

This preserves genuine native layout behavior and keeps the POC small. It avoids implementing and validating a cross-platform layout engine before proving the bridge and update loop.

## Consequences

Subtle measurement can differ between platforms, so this does not yet guarantee pixel-identical layout. Before Android or advanced primitives, benchmark and prototype representative nested layouts. If predictable parity cannot be achieved through a documented common model translated to native constraints, move measurement into Go while leaving native controls responsible for rendering and interaction.

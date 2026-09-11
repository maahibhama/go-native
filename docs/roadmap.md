# Implementation roadmap

An item is complete only when it has tests, relevant documentation, successful iOS and Android builds, and device verification where native behavior is involved.

## Production foundation (v0.2 complete)

- [x] Typed `Style`, `Theme`, semantic `Token[T]`, and platform-style contracts
- [x] Environment, media-query, lifecycle, typed dependencies, and `BuildContext`
- [x] Context-aware runtime applications with environment-driven rendering
- [x] Protocol-v7 compatibility projection for existing component modifiers
- [x] Deterministic headless mutation renderer
- [x] Initial Go-owned constraint layout engine and intrinsic measurement interface
- [x] Mounted component paths and component-scoped state, reducer, ref, memo, callback, effect, lifecycle, and media hooks
- [x] Flex growth/shrink/basis, wrapping, aspect ratio, adaptive grid, RTL, and responsive breakpoint layout
- [x] Batched native intrinsic measurement and typed mutation protocol
  - [x] Value-only measurement batch API, result validation, and content/style/constraint cache
  - [x] Protocol capability negotiation and bounded decoder allocations
  - [x] Versioned, fixed-width typed-style record with portable/iOS/Android overrides and golden fixture
  - [x] Protocol v9 mutation integration with bounded UIKit/Android/template/example readers
  - [x] Portable native background/foreground, border, radius, opacity, visibility, and disabled styling
  - [x] Portable native typography, transforms, and shadows/elevation
  - [x] Native iOS/Android platform override group resolution
  - [x] Bounded native measurement request/result wire protocol and golden fixture
  - [x] UIKit/JNI measurement adapters and computed geometry transmission
- [x] Focus tree and lifecycle callbacks wired from both native hosts

## Component system (v0.3 complete)

- [x] Layout composition API: `View`, `Row`, `Column`, `Stack`, `Spacer`, `Center`, `AspectRatio`, `SafeArea`, `KeyboardAvoidingView`, `ScrollView`, and `Divider`
- [x] Fluent min/max constraints, margins, non-uniform padding, absolute positioning, overflow, border, shadow, transform, visibility, typography, and hit-slop modifiers
- [x] Text and rich-span system: selectable text, links, truncation, wrapping, line limits, dynamic type, and custom fonts
  - [x] Portable typed props, rich-span codec, link callbacks, reconciliation, and protocol round trip
- [x] Complete input system: secure, multiline, search, numeric, keyboard/return configuration, validation, formatters, selection, and focus traversal
  - [x] Portable controlled/uncontrolled API, formatters, selection/submit handlers, stable identity, and cleanup
- [x] Pressable, buttons, links, hit slop, interaction states, and native feedback
- [x] Selection controls, pickers, date/time controls, and progress variants
- [x] Media, icons, SVG/vector views, caching, video, and audio controls
- [x] Feedback overlays, menus, dialogs, sheets, and modals
- [x] Structural components: bars, scaffold, form, card, badge, avatar, skeleton, and empty/error states
- [x] Accessibility, dark mode, RTL, responsive, native contract, screenshot, and end-to-end qualification matrix

## Foundation

- [x] Typed virtual tree for View, Column, Row, Text, and Button
- [x] Stable structural and explicit identity
- [x] Create, delete, update, insert, remove, and move reconciliation
- [x] Versioned binary mutation batches
- [x] Thread-safe state and integer event registry
- [x] UIKit and Android Views counter loops
- [x] Minimal build/run/doctor CLI
- [x] Portable performance benchmarks
- [x] Runtime shutdown, iOS action cleanup, and Android JNI renderer teardown
- [x] Cross-platform native mutation and end-to-end event timing hooks

## Public API

- [x] SafeArea and TextFunc
- [x] Accessibility role, hint, focus, and scalable text
- [x] TextInput, Image, ScrollView, Switch, and ProgressIndicator
- [ ] Modal and navigation contracts

Every primitive requires typed props, reconciliation and protocol tests, both native mappings, accessibility behavior, and an example.

## Scale and interaction

- [ ] Keyed List backed by UICollectionView and RecyclerView
- [ ] Native gesture and animation intent contracts
- [x] Foreground/background lifecycle
- [ ] Cancellation and structured renderer errors

## Tooling and distribution

- [x] Minimal, no-overwrite `gonative init <name>` Go application scaffold
- [x] Conventional Android Gradle/AndroidX project, vendored wrapper, and multi-ABI packaging ([ ] offline dependency cache)
- [x] Explicit iOS physical-device compile and signing workflow
- [ ] Component/event logging and UI tree inspector
- [ ] Source reload investigation
- [x] Native timing JSONL collection harness ([ ] controlled interaction automation)

## Future platforms

- [ ] macOS, Windows, and Linux feasibility and renderers

Desktop work begins only after mobile contracts and lifecycle are stable.

# Implementation roadmap

An item is complete only when it has tests, relevant documentation, successful iOS and Android builds, and device verification where native behavior is involved.

## Production foundation (v0.2 in progress)

- [x] Typed `Style`, `Theme`, semantic `Token[T]`, and platform-style contracts
- [x] Environment, media-query, lifecycle, typed dependencies, and `BuildContext`
- [x] Context-aware runtime applications with environment-driven rendering
- [x] Protocol-v7 compatibility projection for existing component modifiers
- [x] Deterministic headless mutation renderer
- [x] Initial Go-owned constraint layout engine and intrinsic measurement interface
- [x] Mounted component paths and component-scoped state, reducer, ref, memo, callback, effect, lifecycle, and media hooks
- [ ] Flex growth/shrink, wrapping, aspect ratio, grid, and RTL layout
- [ ] Batched native intrinsic measurement and typed mutation protocol
- [ ] Focus tree and lifecycle callbacks wired from both native hosts

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
- [ ] Foreground/background lifecycle
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

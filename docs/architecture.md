# Architecture

Go Native's first milestone proves a narrow claim: a Go-owned declarative tree can drive genuine UIKit controls through one batched native boundary, and a native button event can update a label through state, rebuilding, and reconciliation.

## Milestone 0 data flow

```text
Go App → ui.Node tree → runtime.Reconcile → MutationBatch (binary)
       → one cgo/C call → Objective-C renderer → UIKit main queue
       → NodeID-to-UIView registry → UILabel / UIButton / UIStackView

UIButton → HandlerID → exported C-compatible Go entrypoint
         → EventRegistry → Go callback → State update → scheduled render
```

The `ui` package contains only portable primitives and typed properties. `runtime` owns stable structural identity, handler lifetime, reconciliation, scheduling, and serialization. `platform/ios` owns UIKit objects, the main-thread hop, native layout application, and view destruction. The example bridge is the only package that imports C.

The underlying Go and Apple primitives reviewed for this choice are recorded in [the iOS bridge research notes](research/ios-bridge-primitives.md).

## Minimum contracts

- `ui.Node`: stable integer identity, primitive type, typed `Props`, ordered children.
- `runtime.Mutation`: create, delete, update, insert, remove, or move with IDs and complete properties.
- `runtime.MutationBatch`: one ordered set serialized using a versioned little-endian binary protocol.
- `runtime.Renderer`: `Apply(MutationBatch) error`; it accepts a whole batch and is responsible for UI-thread delivery.
- `runtime.EventRegistry`: native code retains only `HandlerID`; no Go pointer crosses the boundary. It supports both action events and UTF-8 value events used by controlled text inputs.

The renderer receives complete props on create/update. This slightly enlarges updates but makes native application deterministic and avoids optional-field ambiguity. The protocol is deliberately small, versioned, and non-JSON.

## Identity and reconciliation

For unkeyed trees, the runtime preserves identity by matching type and structural position. `ui.WithID` marks explicit identity that survives reordering and also keeps event handlers attached to the logical node. Reconciliation updates equal-ID/equal-type nodes in place, recursively creates/deletes subtrees, and computes moves against the evolving child order so every emitted index is immediately applicable by native renderers. A counter click produces one `MutationUpdate` for the label.

## Threading and ownership

State is safe to read or update from goroutines. Rendering is serialized and repeated state changes are coalesced while a render is pending. The iOS renderer copies the batch bytes synchronously, then applies the copy on the UIKit main queue. Native callbacks enter Go on the UIKit thread, resolve an integer handler, and schedule rendering; UIKit is never mutated directly from Go application goroutines.

Go owns nodes, state, and callbacks. Objective-C owns UIKit views and action targets. Delete mutations remove native views and node-keyed action targets; removed virtual subtrees release handler registry entries. When the view-controller owner is destroyed, its bridge stops the Go runtime before releasing the remaining native registries. The iOS renderer does not retain Go pointers.

Android uses the same bytes without a second protocol: JNI copies the payload into a Java byte array, Java schedules it with `runOnUiThread`, and a `LongSparseArray<View>` retains native view identity. Go-created threads obtain a thread-local `JNIEnv` from the cached `JavaVM`, attach only when necessary, and detach before exit. Activity destruction stops the runtime, clears the view registry, and deletes the JNI global renderer reference under a mutex. The Java renderer retains no Go pointers.

## Layout scope

Milestone 0 translates `Row` and `Column` to `UIStackView`, including padding, gap, alignment, width, and height. This is intentionally not CSS. The longer-term layout choice and its tradeoffs are recorded in [ADR 0002](adr/0002-native-layout-milestone-0.md).

## Deferred work

Keyed-list APIs, native virtualized lists, navigation, animation, gestures, and hot reload remain deferred. Controlled `TextInput` is mapped to `UITextField` and `EditText`; edits cross the bridge as copied UTF-8 strings associated with stable handler IDs. Boolean `Switch` values use a typed boolean callback, while determinate `ProgressIndicator` values are clamped to 0...1 and map to native progress views. Accessibility labels, hints, roles, focus intent, and user-scaled text map to native UIKit and Android accessibility APIs through typed props.

`Image(source)` accepts a bundle/resource name, not a URL or filesystem path. UIKit resolves it with `imageNamed:`; Android resolves a `drawable` and then `mipmap` resource in the application package. Network fetching and caching stay outside the renderer contract. `ScrollView` owns exactly one child and maps to `UIScrollView`, Android `ScrollView`, or `HorizontalScrollView`; multiple items belong inside a `Row` or `Column`.

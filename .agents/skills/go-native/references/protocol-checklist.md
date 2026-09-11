# Binary protocol checklist

Use this checklist for changes to node or mutation enums, `ui.Props`, handler fields, accessibility, images, scrolling, gestures, animation, mutation ordering, or protocol metadata.

## Outer mutation batch

The canonical writer/Go test reader is `runtime/protocol.go`. Preserve:

- little-endian encoding;
- header order: `uint16 version` (10), `uint32 mutation count`, `uint64 sequence`;
- mutation enum and node enum numeric values;
- per-mutation field order and exact widths:
  1. `uint8 type`, `uint8 nodeType`, `uint64 nodeID`, `uint64 parentID`, `int32 index`, `int32 fromIndex`
  2. `float32 width`, `float32 height`, `float32 padding`, `float32 gap`, `uint8 alignment`, `uint8 bold`, `float32 fontSize`
  3. `uint64 onPress`, `uint64 onChange`, `uint64 onToggle`, `uint8 checked`, `float32 progress`
  4. length-prefixed strings: `Text`, `AccessLabel`, `AccessHint`
  5. `uint8 accessRole`, `uint8 focused`, `uint8 scalesText`
  6. length-prefixed string: `ImageSource`, `uint8 imageMode`, `uint8 horizontal`
  7. length-prefixed bytes: `Interactions`
  8. `uint8 textWrap`, `uint8 textOverflow`, `uint32 maxLines`, `uint8 selectable`
  9. length-prefixed rich-text bytes, `uint64 onLink`
  10. length-prefixed placeholder string
  11. `uint8 inputMode`, `uint8 inputKind`, `uint8 returnKey`, `uint8 capitalization`, `uint8 autoCorrect`
  12. `uint8 secure`, `uint8 multiline`, `uint8 readOnly`, `uint8 validationState`
  13. length-prefixed error string
  14. `int32 selectionStart`, `int32 selectionEnd`, `int32 maxLength`
  15. `uint64 onSubmit`, `uint64 onSelection`
  16. length-prefixed typed-style bytes
  17. `uint8 hasFrame`, `float32 x, y, width, height`
- complete props on create/update.

Mirror the decoder in:

- `platform/ios/GoNativeRenderer.m` (`GNApply`/`GNStyle`);
- `platform/android/src/dev/gonative/counter/MainActivity.java` (`applyOnUiThread`/`style`);
- renderer sources embedded in `cmd/gonative/templates.go`;
- checked-in generated renderers in `examples/my-project/`.

The outer mutation batch version is currently `10`. The native measurement protocol version is currently `2` (`runtime/layout/measurement_protocol.go`). Native implementations compare versions as literals. Bump every copy when the byte layout becomes incompatible.

## Interaction payload

`runtime/interactions.go` serializes gesture count/items followed by animation count/items into `Props.Interactions`. Both native renderers parse this nested payload independently. Keep gesture handler ID position, duration units, scalar ordering, curve/property enum values, and reduce-motion semantics aligned.

## Lifecycle and delivery

- Go serializes one ordered batch for one coarse boundary call.
- iOS must copy bytes before asynchronous main-queue application.
- Android must clone bytes before `runOnUiThread`.
- Native callbacks carry integer handler IDs and copied values only.
- Delete/replacement/stop must clean up native targets and Go registry entries.

## Verification

Add or update Go round-trip tests, runtime behavior tests, and template-generation assertions as applicable. Then run Go race tests, vet, formatting checks, and builds for both affected native platforms. A passing Go round trip alone does not prove the Objective-C and Java decoders are synchronized.

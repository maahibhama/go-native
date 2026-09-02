# Binary protocol checklist

Use this checklist for changes to node or mutation enums, `ui.Props`, handler fields, accessibility, images, scrolling, gestures, animation, mutation ordering, or protocol metadata.

## Outer mutation batch

The canonical writer/Go test reader is `runtime/protocol.go`. Preserve:

- little-endian encoding;
- header order: `uint16 version`, `uint32 mutation count`, `uint64 sequence`;
- mutation enum and node enum numeric values;
- per-mutation field order and exact widths;
- length-prefixed UTF-8 strings and opaque interaction bytes;
- complete props on create/update.

Mirror the decoder in:

- `platform/ios/GoNativeRenderer.m` (`GNApply`/`GNStyle`);
- `platform/android/src/dev/gonative/counter/MainActivity.java` (`applyOnUiThread`/`style`);
- renderer sources embedded in `cmd/gonative/templates.go`;
- checked-in generated renderers in `examples/my-project/`.

The version is currently `7`; native implementations compare it as a literal. Bump every copy when the byte layout becomes incompatible.

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

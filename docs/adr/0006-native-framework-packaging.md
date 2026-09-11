# ADR 0006: Package native renderers as framework dependencies

## Status

Accepted

## Context

Generated projects currently contain copies of the iOS and Android renderers. This
makes application projects large, couples JNI symbols to application package names,
and requires the framework renderer, generator templates, and checked-in example to
be updated together for every native change.

## Decision

Go Native will distribute its native implementation as `GoNativeKit.xcframework`
for iOS and `gonative-runtime.aar` for Android. Generated applications will retain
only minimal launchers, platform configuration, assets, and their app-specific
compiled Go library.

The native libraries communicate with that Go library through a framework-owned,
versioned C ABI using fixed-width scalars, byte buffers, and integer IDs. Android
will use package-independent JNI registration. Mutation and measurement protocols
remain separately versioned.

## Consequences

- Native renderer code has one source of truth per platform.
- Generated applications become smaller and easier to understand.
- Native framework artifacts require release, compatibility, and dependency tooling.
- App and framework versions must negotiate ABI and protocol support at startup.
- Native modules and permissions can be integrated through generated metadata rather
  than copied renderer changes.

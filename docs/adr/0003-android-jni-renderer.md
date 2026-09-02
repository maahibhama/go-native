# ADR 0003: Android uses a small JNI surface and native Views

- Status: Accepted for the Android counter POC
- Date: 2026-09-02

## Decision

Build the Go counter bridge as an arm64 Android `c-shared` library. Use a small C JNI shim to copy each existing binary `MutationBatch` into one Java byte array call. Java owns a `NodeID → View` registry and maps Text, Button, Row, and Column to `TextView`, `Button`, and `LinearLayout`.

Native-to-Go events pass only a `HandlerID` integer. The Java renderer applies every decoded batch through `Activity.runOnUiThread`.

## Rationale

This reuses the exact renderer contract proven on iOS, keeps JNI traffic coarse-grained, and lets Android retain responsibility for controls, drawing, touch processing, accessibility, and layout primitives. Android's JNI guidance recommends minimizing marshalling and call frequency, caching the `JavaVM`, and treating `JNIEnv` as thread-local.

The first build uses installed SDK command-line tools directly because Gradle is not present in the development environment. A conventional Gradle application/plugin should replace this packaging script before distribution; that change does not affect the bridge contract.

## Risks and mitigations

- Go render work may execute on an unattached native thread. The shim calls `GetEnv`, attaches only when needed, and detaches threads it attached.
- Java UI objects are not thread-safe. All mutation application is posted to the UI thread.
- JNI global references can leak across complex Activity lifecycles. The POC replaces the renderer reference on each start; explicit stop/release APIs are required before production lifecycle support.
- JNI exceptions do not unwind through native frames. The shim checks, reports, and clears exceptions after the batch call; production code should propagate structured renderer errors.
- The POC APK contains only arm64. Multi-ABI packaging and release signing are deferred.

## Sources

- [Android JNI tips](https://developer.android.com/ndk/guides/jni-tips)
- [Android processes and threads](https://developer.android.com/guide/components/processes-and-threads)

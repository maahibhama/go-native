# Fast Reload architecture

Go does not support replacing compiled functions or package code inside a running
iOS or Android process. A reliable first reload workflow should therefore rebuild
the native application and restart it, while preserving fast incremental steps:

1. watch project-owned Go source outside the application process;
2. debounce changes and invoke the existing platform build command;
3. reinstall/relaunch only after a successful build;
4. report compiler failures without terminating the last installed build;
5. optionally restore explicitly serializable application state after relaunch.

The iOS Simulator and Android emulator can both use this model. Physical iOS
devices add signing and installation latency. Android can replace the debug APK
when its persistent signing key matches, but application process state is lost.

Dynamic plugins are not a portable alternative: Go's plugin mechanism is not
supported on these mobile targets and would conflict with platform code-signing
rules. Interpreting Go or introducing a JavaScript/Lua runtime would undermine the
project's all-Go native-rendering goal and is out of scope.

The implemented `gonative dev <ios|android>` workflow hashes project Go sources,
module files, assets, and `gonative.yaml`, while excluding generated build output.
It keeps the last working application installed when a build fails and reuses the
native framework artifact after the initial build. Each development daemon assigns
a session ID to the generated Go bridge. `ui.UseReloadState` persists explicitly
selected JSON values under that session and restores them before the next initial
render. Arbitrary Go memory, effects, native handles, focus, and keyboard state are
intentionally not reflected or restored.

## Virtual Reload transport

The virtual-reload foundation uses `runtime/devtransport`, a debug-only,
authenticated, length-prefixed little-endian transport. It envelopes unchanged
mutation and measurement protocol payloads and assigns every frame to a monotonic
worker generation. `runtime.RemoteRenderer` sends mutation batches and reset-tree
commands; `layout.RemoteMeasurer` provides multiplexed native measurement RPC.
`Runtime.Reload` clears mounted hooks, effects, handlers, and geometry before a
fresh complete render while leaving its renderer connection alive.

`gonative dev start` provides an interactive terminal controller. Use `i` or
`a` to select and launch a platform, `r` to force a reload, `d` to run toolchain
diagnostics, and `q` to stop the controller. Source changes automatically reload
the most recently selected platform.

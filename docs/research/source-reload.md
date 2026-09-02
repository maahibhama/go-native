# Source reload feasibility

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

Before implementation, measure rebuild, install, and relaunch latency separately.
The watcher should live in `gonative`, remain opt-in, ignore generated build output,
cancel superseded builds, and never execute source received through the diagnostics
inspector. A later state-restoration contract must be versioned and application-
controlled rather than reflecting arbitrary Go memory.


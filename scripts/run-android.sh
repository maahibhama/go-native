#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
SDK=${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/Library/Android/sdk}}
ADB="$SDK/platform-tools/adb"
APK=$($ROOT/scripts/build-android.sh)
adb_target() {
    if [ -n "${GONATIVE_ANDROID_SERIAL:-}" ]; then
        "$ADB" -s "$GONATIVE_ANDROID_SERIAL" "$@"
    else
        "$ADB" "$@"
    fi
}
if ! adb_target get-state >/dev/null 2>&1; then
    echo "No running Android device or emulator. Start one, then rerun this command." >&2
    echo "If several are connected, set GONATIVE_ANDROID_SERIAL to the desired adb serial." >&2
    exit 1
fi
adb_target install -r "$APK"
adb_target shell am force-stop dev.gonative.counter
adb_target shell am start -n dev.gonative.counter/.MainActivity

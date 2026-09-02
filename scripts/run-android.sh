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
if install_output=$(adb_target install -r "$APK" 2>&1); then
    echo "$install_output"
elif echo "$install_output" | grep -q 'INSTALL_FAILED_UPDATE_INCOMPATIBLE'; then
    if [ "${GONATIVE_ANDROID_REINSTALL:-0}" = "1" ]; then
        echo "Removing the incompatible development install (its local app data will be deleted)." >&2
        adb_target uninstall dev.gonative.counter
        adb_target install "$APK"
    else
        echo "$install_output" >&2
        echo "The installed development app was signed by an older ephemeral key." >&2
        echo "To remove it once and reinstall, run:" >&2
        echo "  GONATIVE_ANDROID_REINSTALL=1 gonative run android" >&2
        echo "Warning: reinstall recovery deletes the counter app's local data." >&2
        exit 1
    fi
else
    echo "$install_output" >&2
    exit 1
fi
adb_target shell am force-stop dev.gonative.counter
adb_target shell am start -n dev.gonative.counter/.MainActivity

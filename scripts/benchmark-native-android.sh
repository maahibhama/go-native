#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
DURATION=${GONATIVE_BENCHMARK_SECONDS:-15}
SDK=${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/Library/Android/sdk}}
ADB="$SDK/platform-tools/adb"
adb_target() {
    if [ -n "${GONATIVE_ANDROID_SERIAL:-}" ]; then
        "$ADB" -s "$GONATIVE_ANDROID_SERIAL" "$@"
    else
        "$ADB" "$@"
    fi
}

echo "Launching the instrumented Android counter. Tap Increment during the next $DURATION seconds." >&2
GONATIVE_BENCHMARK=1 "$ROOT/scripts/run-android.sh" >/dev/null
adb_target logcat -c
adb_target shell am force-stop dev.gonative.counter
adb_target shell am start -n dev.gonative.counter/.MainActivity >/dev/null
LOG="$ROOT/build/android/native-timing.log"
adb_target logcat -v raw > "$LOG" &
log_pid=$!
trap 'kill "$log_pid" 2>/dev/null || true' EXIT INT TERM
sleep "$DURATION"
kill "$log_pid" 2>/dev/null || true
wait "$log_pid" 2>/dev/null || true
trap - EXIT INT TERM
grep 'GONATIVE_TIMING ' "$LOG" | sed 's/^.*GONATIVE_TIMING //' || true

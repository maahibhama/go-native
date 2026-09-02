#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
DURATION=${GONATIVE_BENCHMARK_SECONDS:-15}
DEVICE=${GONATIVE_SIMULATOR:-"iPhone 17 Pro"}
UDID=$(xcrun simctl list devices available -j | /usr/bin/python3 -c 'import json,sys; name=sys.argv[1]; devices=json.load(sys.stdin)["devices"]; matches=[d for group in devices.values() for d in group if d["name"]==name and d.get("isAvailable")]; print(matches[0]["udid"] if matches else "")' "$DEVICE")
if [ -z "$UDID" ]; then echo "No available simulator named: $DEVICE" >&2; exit 1; fi

echo "Launching the instrumented iOS counter. Tap Increment during the next $DURATION seconds." >&2
GONATIVE_BENCHMARK=1 "$ROOT/scripts/run-ios.sh" >/dev/null
LOG="$ROOT/build/ios-simulator/native-timing.log"
xcrun simctl spawn "$UDID" log stream --style compact --predicate 'eventMessage CONTAINS "GONATIVE_TIMING"' > "$LOG" &
log_pid=$!
trap 'kill "$log_pid" 2>/dev/null || true' EXIT INT TERM
sleep "$DURATION"
kill "$log_pid" 2>/dev/null || true
wait "$log_pid" 2>/dev/null || true
trap - EXIT INT TERM
grep 'GONATIVE_TIMING ' "$LOG" | sed 's/^.*GONATIVE_TIMING //' || true

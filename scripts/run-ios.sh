#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
APP=$($ROOT/scripts/build-ios.sh)
DEVICE=${GONATIVE_SIMULATOR:-"iPhone 17 Pro"}
UDID=$(xcrun simctl list devices available -j | /usr/bin/python3 -c 'import json,sys; name=sys.argv[1]; devices=json.load(sys.stdin)["devices"]; matches=[d for group in devices.values() for d in group if d["name"]==name and d.get("isAvailable")]; print(matches[0]["udid"] if matches else "")' "$DEVICE")
if [ -z "$UDID" ]; then echo "No available simulator named: $DEVICE" >&2; exit 1; fi
xcrun simctl boot "$UDID" 2>/dev/null || true
open -a Simulator
xcrun simctl install "$UDID" "$APP"
xcrun simctl launch "$UDID" dev.gonative.counter

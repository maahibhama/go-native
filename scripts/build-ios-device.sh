#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BUILD="$ROOT/build/ios-device"
APP="$BUILD/GoNativeCounter.app"
IDENTITY=${GONATIVE_IOS_SIGNING_IDENTITY:-}
PROFILE=${GONATIVE_IOS_PROVISIONING_PROFILE:-}

if [ -z "$IDENTITY" ]; then
    echo "Missing GONATIVE_IOS_SIGNING_IDENTITY (for example: Apple Development: Name (TEAMID))" >&2
    exit 1
fi
if [ -z "$PROFILE" ]; then
    echo "Missing GONATIVE_IOS_PROVISIONING_PROFILE (path to a .mobileprovision file)" >&2
    exit 1
fi
if [ ! -f "$PROFILE" ]; then
    echo "iOS provisioning profile does not exist: $PROFILE" >&2
    exit 1
fi
if ! security find-identity -v -p codesigning | grep -F "$IDENTITY" >/dev/null; then
    echo "iOS signing identity is not available in the keychain: $IDENTITY" >&2
    exit 1
fi

SDK=$(xcrun --sdk iphoneos --show-sdk-path)
rm -rf "$BUILD"
mkdir -p "$APP"

cd "$ROOT"
LDFLAGS=
if [ "${GONATIVE_BENCHMARK:-0}" = "1" ]; then LDFLAGS='-X main.benchmarkOutput=1'; fi
CGO_ENABLED=1 GOOS=ios GOARCH=arm64 \
CC="$(xcrun -f clang) -target arm64-apple-ios15.0 -isysroot $SDK" \
GOCACHE=${GOCACHE:-/tmp/go-native-gocache} \
GOPATH=${GOPATH:-/tmp/go-native-gopath} \
go build -ldflags "$LDFLAGS" -buildmode=c-archive -o "$BUILD/counter.a" ./examples/counter/bridge

xcrun --sdk iphoneos clang -target arm64-apple-ios15.0 -isysroot "$SDK" \
    -fobjc-arc -framework UIKit -framework Foundation -framework CoreGraphics \
    -I"$BUILD" -I"$ROOT/platform/ios" \
    "$ROOT/platform/ios/main.m" "$ROOT/platform/ios/GoNativeRenderer.m" \
    "$BUILD/counter.a" -o "$APP/GoNativeCounter"
cp "$ROOT/platform/ios/Info.plist" "$APP/Info.plist"
cp "$PROFILE" "$APP/embedded.mobileprovision"

security cms -D -i "$PROFILE" > "$BUILD/profile.plist"
/usr/libexec/PlistBuddy -x -c 'Print :Entitlements' "$BUILD/profile.plist" > "$BUILD/entitlements.plist"
codesign --force --sign "$IDENTITY" --entitlements "$BUILD/entitlements.plist" "$APP"
codesign --verify --deep --strict "$APP"
echo "$APP"

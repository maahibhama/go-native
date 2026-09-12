#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BUILD="$ROOT/build/ios-simulator"
APP="$BUILD/GoNativeCounter.app"
SDK=$(xcrun --sdk iphonesimulator --show-sdk-path)
mkdir -p "$BUILD" "$APP"
cd "$ROOT"
"$ROOT/scripts/build-ios-framework.sh" >/dev/null
LDFLAGS=
if [ "${GONATIVE_BENCHMARK:-0}" = "1" ]; then LDFLAGS='-X main.benchmarkOutput=1'; fi
CGO_ENABLED=1 GOOS=ios GOARCH=arm64 CC="$(xcrun -f clang) -target arm64-apple-ios15.0-simulator -isysroot $SDK" go build -ldflags "$LDFLAGS" -buildmode=c-archive -o "$BUILD/frameworktest.a" ./internal/frameworktest/bridge
xcrun --sdk iphonesimulator clang -target arm64-apple-ios15.0-simulator -isysroot "$SDK" -fobjc-arc \
    -framework UIKit -framework Foundation -framework CoreGraphics \
    -I"$BUILD" -I"$ROOT/build/native/GoNativeKit.xcframework/ios-arm64-simulator/Headers" \
    "$ROOT/platform/ios/main.m" \
    "$ROOT/build/native/GoNativeKit.xcframework/ios-arm64-simulator/libGoNativeKit.a" \
    "$BUILD/frameworktest.a" -o "$APP/GoNativeCounter"
cp "$ROOT/platform/ios/Info.plist" "$APP/Info.plist"
codesign --force --sign - "$APP"
echo "$APP"

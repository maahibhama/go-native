#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BUILD=${GONATIVE_ARTIFACT_DIR:-"$ROOT/build/native"}
WORK="$BUILD/ios-framework"
OUTPUT="$BUILD/GoNativeKit.xcframework"
SOURCES="GNProtocolReader.m GNViewRegistry.m GNRuntimeHost.m GNMeasurementHost.m GNControls.m GoNativeRenderer.m"

rm -rf "$WORK" "$OUTPUT"
mkdir -p "$WORK/device/Headers" "$WORK/simulator/Headers"

for header in GoNativeKit.h GoNativeRenderer.h GNControls.h GNMeasurementHost.h GNProtocolReader.h GNViewRegistry.h module.modulemap; do
    cp "$ROOT/platform/ios/$header" "$WORK/device/Headers/$header"
    cp "$ROOT/platform/ios/$header" "$WORK/simulator/Headers/$header"
done

build_slice() {
    name=$1
    sdk=$2
    target=$3
    slice="$WORK/$name"
    sdk_path=$(xcrun --sdk "$sdk" --show-sdk-path)
    objects=
    for source in $SOURCES; do
        object="$slice/${source%.m}.o"
        xcrun --sdk "$sdk" clang -target "$target" -isysroot "$sdk_path" \
            -fobjc-arc -fmodules -fmodule-name=GoNativeKit \
            -fmodules-cache-path="$WORK/module-cache" -I"$ROOT/platform/ios" \
            -c "$ROOT/platform/ios/$source" -o "$object"
        objects="$objects $object"
    done
    # App ABI symbols intentionally remain unresolved until the application Go archive is linked.
    xcrun ar rcs "$slice/libGoNativeKit.a" $objects
}

build_slice device iphoneos arm64-apple-ios15.0
build_slice simulator iphonesimulator arm64-apple-ios15.0-simulator

xcodebuild -create-xcframework \
    -library "$WORK/device/libGoNativeKit.a" -headers "$WORK/device/Headers" \
    -library "$WORK/simulator/libGoNativeKit.a" -headers "$WORK/simulator/Headers" \
    -output "$OUTPUT"

echo "$OUTPUT"

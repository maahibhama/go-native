#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
SDK=${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/Library/Android/sdk}}
NDK_VERSION=${GONATIVE_NDK_VERSION:-28.2.13676358}
NDK="$SDK/ndk/$NDK_VERSION"
if [ ! -d "$NDK" ]; then
    for d in "$SDK/ndk/"*; do
        if [ -d "$d" ]; then NDK="$d"; break; fi
    done
fi
HOST_TAG=${GONATIVE_NDK_HOST_TAG:-darwin-x86_64}
TOOLCHAIN="$NDK/toolchains/llvm/prebuilt/$HOST_TAG"
BUILD="$ROOT/build/android"
LIB_BUILD="$BUILD/lib.next.$$"
ABIS=${GONATIVE_ANDROID_ABIS:-"arm64-v8a,x86_64"}

mkdir -p "$BUILD" "$LIB_BUILD"
trap 'rm -rf "$LIB_BUILD"' EXIT INT TERM

old_ifs=$IFS
IFS=,
for abi in $ABIS; do
    case "$abi" in
        arm64-v8a) goarch=arm64; compiler=aarch64-linux-android23-clang ;;
        x86_64) goarch=amd64; compiler=x86_64-linux-android23-clang ;;
        *) continue ;;
    esac
    if [ ! -x "$TOOLCHAIN/bin/$compiler" ]; then
        echo "Missing Android NDK compiler: $TOOLCHAIN/bin/$compiler" >&2
        exit 1
    fi
    mkdir -p "$LIB_BUILD/$abi"
    (
        cd "$ROOT"
        CGO_ENABLED=1 GOOS=android GOARCH="$goarch" \
        CC="$TOOLCHAIN/bin/$compiler" \
        CGO_CFLAGS="--sysroot=$TOOLCHAIN/sysroot -I$TOOLCHAIN/sysroot/usr/include" \
        go build -buildmode=c-shared -o "$LIB_BUILD/$abi/libgonative.so" ./android/bridge
    )
    rm -f "$LIB_BUILD/$abi/libgonative.h"
done
IFS=$old_ifs

rm -rf "$BUILD/lib"
mv "$LIB_BUILD" "$BUILD/lib"
trap - EXIT INT TERM

#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
SDK=${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/Library/Android/sdk}}
NDK_VERSION=${GONATIVE_NDK_VERSION:-28.2.13676358}
NDK="$SDK/ndk/$NDK_VERSION"
HOST_TAG=${GONATIVE_NDK_HOST_TAG:-darwin-x86_64}
TOOLCHAIN="$NDK/toolchains/llvm/prebuilt/$HOST_TAG"
BUILD="$ROOT/build/android"
LIB_BUILD="$BUILD/lib.next.$$"
ABIS=${GONATIVE_ANDROID_ABIS:-"arm64-v8a,x86_64"}
LDFLAGS=
if [ "${GONATIVE_BENCHMARK:-0}" = "1" ]; then LDFLAGS='-X main.benchmarkOutput=1'; fi

mkdir -p "$BUILD"
mkdir -p "$LIB_BUILD"
trap 'rm -rf "$LIB_BUILD"' EXIT INT TERM
old_ifs=$IFS
IFS=,
for abi in $ABIS; do
    case "$abi" in
        arm64-v8a) goarch=arm64; compiler=aarch64-linux-android23-clang ;;
        x86_64) goarch=amd64; compiler=x86_64-linux-android23-clang ;;
        *) echo "Unsupported Android ABI: $abi" >&2; exit 1 ;;
    esac
    if [ ! -x "$TOOLCHAIN/bin/$compiler" ]; then echo "Missing Android compiler: $TOOLCHAIN/bin/$compiler" >&2; exit 1; fi
    mkdir -p "$LIB_BUILD/$abi"
    CGO_ENABLED=1 GOOS=android GOARCH="$goarch" \
    CC="$TOOLCHAIN/bin/$compiler" \
    CGO_CFLAGS="--sysroot=$TOOLCHAIN/sysroot -I$TOOLCHAIN/sysroot/usr/include" \
    GOCACHE=${GOCACHE:-/tmp/go-native-gocache} \
    GOPATH=${GOPATH:-/tmp/go-native-gopath} \
    go build -ldflags "$LDFLAGS" -buildmode=c-shared -o "$LIB_BUILD/$abi/libgonative.so" ./examples/counter/androidbridge
    rm -f "$LIB_BUILD/$abi/libgonative.h"
done
IFS=$old_ifs
rm -rf "$BUILD/lib"
mv "$LIB_BUILD" "$BUILD/lib"
trap - EXIT INT TERM

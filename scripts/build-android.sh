#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
SDK=${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/Library/Android/sdk}}
BUILD_TOOLS=${GONATIVE_BUILD_TOOLS:-36.0.0}
PLATFORM=${GONATIVE_ANDROID_PLATFORM:-35}
NDK_VERSION=${GONATIVE_NDK_VERSION:-28.2.13676358}
TOOLS="$SDK/build-tools/$BUILD_TOOLS"
ANDROID_JAR="$SDK/platforms/android-$PLATFORM/android.jar"
NDK="$SDK/ndk/$NDK_VERSION"
HOST_TAG=darwin-x86_64
TOOLCHAIN="$NDK/toolchains/llvm/prebuilt/$HOST_TAG"
BUILD="$ROOT/build/android"
CLASSES="$BUILD/classes"
APK="$BUILD/GoNativeCounter.apk"
STATE="$ROOT/.gonative"
KEYSTORE="$STATE/android-debug.keystore"

for path in "$TOOLS/aapt2" "$TOOLS/d8" "$TOOLS/apksigner" "$ANDROID_JAR" "$TOOLCHAIN/bin/aarch64-linux-android23-clang"; do
    if [ ! -e "$path" ]; then echo "Missing Android tool: $path" >&2; exit 1; fi
done

rm -rf "$BUILD"
mkdir -p "$CLASSES" "$BUILD/lib/arm64-v8a" "$BUILD/compiled-res" "$STATE"

cd "$ROOT"
CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
CC="$TOOLCHAIN/bin/aarch64-linux-android23-clang" \
CGO_CFLAGS="--sysroot=$TOOLCHAIN/sysroot -I$TOOLCHAIN/sysroot/usr/include" \
GOCACHE=${GOCACHE:-/tmp/go-native-gocache} \
GOPATH=${GOPATH:-/tmp/go-native-gopath} \
go build -buildmode=c-shared -o "$BUILD/lib/arm64-v8a/libgonative.so" ./examples/counter/androidbridge
rm -f "$BUILD/lib/arm64-v8a/libgonative.h"

find "$ROOT/platform/android/src" -name '*.java' -print > "$BUILD/java-sources.txt"
javac -source 8 -target 8 -bootclasspath "$ANDROID_JAR" -d "$CLASSES" @"$BUILD/java-sources.txt"
"$TOOLS/d8" --lib "$ANDROID_JAR" --min-api 23 --output "$BUILD" $(find "$CLASSES" -name '*.class' -print)
"$TOOLS/aapt2" compile --dir "$ROOT/platform/android/res" -o "$BUILD/compiled-res"
"$TOOLS/aapt2" link -o "$BUILD/unsigned.apk" -I "$ANDROID_JAR" --manifest "$ROOT/platform/android/AndroidManifest.xml" --min-sdk-version 23 --target-sdk-version "$PLATFORM" --java "$BUILD/generated" "$BUILD/compiled-res"/*.flat
cp "$BUILD/unsigned.apk" "$APK"
(cd "$BUILD" && zip -q "$APK" classes.dex && zip -qr "$APK" lib)

if [ ! -f "$KEYSTORE" ]; then
    keytool -genkeypair -keystore "$KEYSTORE" -storepass android -keypass android -alias androiddebugkey -dname "CN=Android Debug,O=Go Native,C=US" -keyalg RSA -keysize 2048 -validity 10000 >/dev/null 2>&1
fi
"$TOOLS/zipalign" -f 4 "$APK" "$BUILD/aligned.apk"
"$TOOLS/apksigner" sign --ks "$KEYSTORE" --ks-pass pass:android --key-pass pass:android --out "$APK" "$BUILD/aligned.apk"
"$TOOLS/apksigner" verify "$APK"
echo "$APK"

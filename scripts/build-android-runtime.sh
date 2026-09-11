#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BUILD=${GONATIVE_ARTIFACT_DIR:-"$ROOT/build/native"}
GRADLE=${GONATIVE_GRADLE:-"$ROOT/platform/android/gradlew"}

mkdir -p "$BUILD"
GONATIVE_MAVEN_REPO=${GONATIVE_MAVEN_REPO:-"$BUILD/maven"}
export GONATIVE_MAVEN_REPO
"$GRADLE" -p "$ROOT/platform/android" :runtime:assembleDebug :runtime:publishDebugPublicationToGoNativeRepository --no-daemon
cp "$ROOT/platform/android/runtime/build/outputs/aar/runtime-debug.aar" "$BUILD/gonative-runtime.aar"
echo "$BUILD/gonative-runtime.aar"
echo "$GONATIVE_MAVEN_REPO/dev/gonative/gonative-runtime"

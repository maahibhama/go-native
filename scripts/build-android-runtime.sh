#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BUILD=${GONATIVE_ARTIFACT_DIR:-"$ROOT/build/native"}
GRADLE=${GONATIVE_GRADLE:-"$ROOT/platform/android/gradlew"}

mkdir -p "$BUILD"
"$GRADLE" -p "$ROOT/platform/android" :runtime:assembleDebug --no-daemon
cp "$ROOT/platform/android/runtime/build/outputs/aar/runtime-debug.aar" "$BUILD/gonative-runtime.aar"
echo "$BUILD/gonative-runtime.aar"

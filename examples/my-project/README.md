# my-project

A Go Native application. The UI is declared in app.go and renders genuine platform-native controls on iOS and Android.

After signing in, the dashboard includes a production-layout API example covering mounted hooks, responsive breakpoints, flex metadata, adaptive grid columns, aspect ratio, and logical LTR/RTL layout.

## Project Map
- app.go: Declarative UI tree written in Go.
- ios/: Native iOS host project (UIKit, Xcode-compatible).
- android/: Native Android host project (Android Views + Gradle, Android Studio-compatible).

## Development Commands

```bash
# Check local toolchain
gonative doctor

# Build & run on iOS Simulator
gonative build ios
gonative run ios

# Build & run on Android
gonative build android
gonative run android
```

## IDE Usage
- **Xcode**: Open `ios/my-project.xcodeproj` in Xcode and click **Run** (Cmd+R).
- **Android Studio**: Open the `android/` directory in Android Studio and click **Run** (Shift+F10).

# Go Native (`go-native`)

A declarative Go UI framework that compiles and renders genuine platform-native controls on iOS (UIKit) and Android (Android Views) with zero WebView, JavaScript, React Native bridge JSON, Flutter, or Skia canvas overhead.

---

## Key Highlights

- **Direct Native Controls**: Renders UIKit (`UILabel`, `UIButton`, `UIStackView`, `UITextField`, `UISwitch`, `UIProgressView`, `UIImageView`, `UIScrollView`) on iOS and Android Views (`TextView`, `Button`, `LinearLayout`, `EditText`, `Switch`, `ProgressBar`, `ImageView`, `ScrollView`) on Android.
- **Binary Mutation Batch**: Virtual UI reconciliation emits a compact, little-endian binary batch (`runtime.MutationBatch`) sent across cgo/JNI in **one coarse call per render pass**.
- **Zero Go Pointers Across Boundaries**: Native event listeners retain only 64-bit integer `HandlerID`s. No Go pointers or objects cross the foreign function interface.
- **IDE-First Native Workflows**: `gonative init` generates genuine **Xcode** (`.xcodeproj`) and **Android Studio** projects ready for one-click build and debug.
- **Declarative Gestures & Spring Animations**: Declarative touch intents (Tap, LongPress, Swipe, Pan) and spring/cubic Bézier animations compiled directly into native platform animators (`UIView.animate` / `ValueAnimator`).
- **Live Diagnostics & Inspector**: Built-in loopback HTTP inspector (`GET /v1/tree`, `GET /v1/logs`) for inspecting virtual trees and performance metrics in real time.

---

## Supported Primitives

| Category | Go Native Primitive | iOS (UIKit) | Android (Views) |
|---|---|---|---|
| **Layout** | `ui.View`, `ui.SafeArea` | `UIView`, `GNSafeAreaView` | `LinearLayout` (`fitsSystemWindows`) |
| **Flex Containers** | `ui.Column`, `ui.Row` | `UIStackView` (vertical / horizontal) | `LinearLayout` (vertical / horizontal) |
| **Typography** | `ui.Text` | `UILabel` (Dynamic Type support) | `TextView` (SP font scaling) |
| **Actions** | `ui.Button` | `UIButton` (Target-Action) | `Button` (`OnClickListener`) |
| **Inputs** | `ui.TextInput` | `UITextField` | `EditText` (`TextWatcher`) |
| **Toggles** | `ui.Switch` | `UISwitch` | `Switch` (`OnCheckedChangeListener`) |
| **Indicators** | `ui.ProgressIndicator` | `UIProgressView` | `ProgressBar` (horizontal style) |
| **Media** | `ui.Image` | `UIImageView` (asset bundles) | `ImageView` (drawables / mipmaps) |
| **Scrolling** | `ui.ScrollView` | `UIScrollView` | `ScrollView` / `HorizontalScrollView` |

---

## Requirements

- **Host**: macOS on Apple Silicon (or Intel)
- **Go**: Go 1.24 or newer
- **iOS**: Xcode 15+ with iOS Simulator runtime (`xcode-select -p`)
- **Android**: Android SDK (Platform 35, Build-Tools 36.0.0, NDK 28.2+, JDK 17)

Verify your local toolchain at any time with:
```bash
go run ./cmd/gonative doctor
```

---

## Quick Start & CLI

Install the `gonative` CLI globally:
```bash
go install ./cmd/gonative
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 1. Initialize a new project
```bash
gonative init my-app
cd my-app
```

This scaffolds a complete, standalone native application:
```text
my-app/
├── app.go                  # Pure Go declarative UI tree
├── go.mod                  # Module dependencies
├── ios/                    # Native UIKit host + Xcode project
│   ├── my-app.xcodeproj/   # Open directly in Xcode (Cmd+R)
│   ├── main.m
│   ├── GoNativeRenderer.h/.m
│   └── bridge/
└── android/                # Native Android Views host + Gradle project
    ├── build.gradle        # Open directory in Android Studio (Shift+F10)
    ├── gradlew
    ├── build-libs.sh       # NDK shared library compiler
    ├── bridge/
    └── app/
        ├── build.gradle
        └── src/main/java/dev/gonative/my_app/MainActivity.java
```

### 2. Run via CLI

```bash
# Build & run on iOS Simulator
gonative run ios

# Build & run on Android device / emulator
gonative run android

# Build for physical iOS device (requires signing identity & profile)
GONATIVE_IOS_SIGNING_IDENTITY="Apple Development: Name (TEAMID)" \
GONATIVE_IOS_PROVISIONING_PROFILE=/path/to/profile.mobileprovision \
gonative build ios-device
```

### 3. Open & Run in IDEs

- **Xcode**: Double click or run `open ios/my-app.xcodeproj`. Select your target simulator and press **Cmd + R**. Xcode automatically invokes Go build phases to compile the bridge archive.
- **Android Studio**: Open Android Studio, select **Open**, and choose the `android/` directory. Let Gradle sync and press **Shift + F10**.

---

## Example UI Code

```go
package app

import (
	"fmt"
	"github.com/go-native/go-native/ui"
)

func App() ui.Component {
	count := ui.UseState(0)

	return ui.SafeArea(
		ui.Column(
			ui.Text("Go Native").FontSize(32).Bold(),
			ui.Text(fmt.Sprintf("Current count: %d", count.Get())).FontSize(18),
			ui.Row(
				ui.Button("Decrement").OnClick(func() {
					count.Set(count.Get() - 1)
				}),
				ui.Button("Increment").OnClick(func() {
					count.Set(count.Get() + 1)
				}),
			).Gap(16),
		).Padding(24).Gap(16).Align(ui.AlignCenter),
	)
}
```

---

## Testing & Benchmarks

```bash
# Run all unit tests with data race detector
go test -race ./...

# Run static analysis and code formatting verification
go vet ./...
gofmt -l .

# Run Go reconciliation microbenchmarks
make benchmark

# Run Native Sampling Benchmarks (end-to-end bridge roundtrip)
gonative benchmark native ios
gonative benchmark native android
```

See [docs/performance.md](docs/performance.md) for benchmark baselines and [docs/architecture.md](docs/architecture.md) for deep technical architecture details.

---

## Repository Map

- [`ui/`](./ui): Declarative primitives, typed `Props`, `State[T]`, gesture/animation intents, and presentation models.
- [`runtime/`](./runtime): Virtual tree reconciler, identity stabilization, little-endian binary protocol serializer, thread-safe event registry, and timing telemetry.
- [`runtime/inspector/`](./runtime/inspector): Loopback HTTP diagnostic server (`GET /v1/tree`, `GET /v1/logs`).
- [`platform/ios/`](./platform/ios): Objective-C UIKit host and renderer.
- [`platform/android/`](./platform/android): Java Android Views host, GapDrawable, and Gradle build harness.
- [`cmd/gonative/`](./cmd/gonative): Developer CLI (`init`, `doctor`, `build`, `run`, `benchmark native`).
- [`examples/counter/`](./examples/counter): Interactive counter demo application and cgo/JNI bridge entrypoints.
- [`docs/`](./docs): Architectural decision records (ADRs), roadmap, performance baselines, and diagnostics guides.

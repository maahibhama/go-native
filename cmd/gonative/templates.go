package main

import (
	"fmt"
	"strings"
)

func sanitizePackageName(name string) string {
	var b strings.Builder
	for _, ch := range strings.ToLower(name) {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			b.WriteRune(ch)
		} else if ch == '_' || ch == '-' {
			b.WriteRune('_')
		}
	}
	res := b.String()
	if len(res) == 0 || (res[0] >= '0' && res[0] <= '9') {
		res = "app_" + res
	}
	return res
}

func sanitizeBundleIdentifierPart(name string) string {
	return strings.ReplaceAll(sanitizePackageName(name), "_", "-")
}

func generatePbxproj(name, pkg string) string {
	return fmt.Sprintf(`// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 54;
	objects = {

/* Begin PBXBuildFile section */
		100000000000000000000014 /* AppDelegate.m in Sources */ = {isa = PBXBuildFile; fileRef = 100000000000000000000013 /* AppDelegate.m */; };
		100000000000000000000011 /* main.m in Sources */ = {isa = PBXBuildFile; fileRef = 100000000000000000000010 /* main.m */; };
/* End PBXBuildFile section */

/* Begin PBXFileReference section */
		100000000000000000000003 /* %s.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "%s.app"; sourceTree = BUILT_PRODUCTS_DIR; };
		100000000000000000000012 /* AppDelegate.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = AppDelegate.h; sourceTree = "<group>"; };
		100000000000000000000013 /* AppDelegate.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = AppDelegate.m; sourceTree = "<group>"; };
		100000000000000000000010 /* main.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = main.m; sourceTree = "<group>"; };
		100000000000000000000015 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXFrameworksBuildPhase section */
		100000000000000000000006 /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXFrameworksBuildPhase section */

/* Begin PBXGroup section */
		100000000000000000000004 = {
			isa = PBXGroup;
			children = (
				100000000000000000000012 /* AppDelegate.h */,
				100000000000000000000013 /* AppDelegate.m */,
				100000000000000000000010 /* main.m */,
				100000000000000000000015 /* Info.plist */,
				100000000000000000000003 /* %s.app */,
			);
			sourceTree = "<group>";
		};
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
		100000000000000000000002 /* %s */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 10000000000000000000000A /* Build configuration list for PBXNativeTarget "%s" */;
			buildPhases = (
				100000000000000000000008 /* Build Go Native Bridge */,
				100000000000000000000005 /* Sources */,
				100000000000000000000006 /* Frameworks */,
				100000000000000000000007 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = "%s";
			productName = "%s";
			productReference = 100000000000000000000003 /* %s.app */;
			productType = "com.apple.product-type.application";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		100000000000000000000001 /* Project object */ = {
			isa = PBXProject;
			attributes = {
				BuildIndependentTargetsInParallel = 1;
				LastUpgradeCheck = 1500;
			};
			buildConfigurationList = 100000000000000000000009 /* Build configuration list for PBXProject "%s" */;
			compatibilityVersion = "Xcode 14.0";
			developmentRegion = en;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 100000000000000000000004;
			productRefGroup = 100000000000000000000004;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				100000000000000000000002 /* %s */,
			);
		};
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
		100000000000000000000007 /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXResourcesBuildPhase section */

/* Begin PBXShellScriptBuildPhase section */
		100000000000000000000008 /* Build Go Native Bridge */ = {
			isa = PBXShellScriptBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			inputFileListPaths = (
			);
			inputPaths = (
			);
			name = "Build Go Native Bridge";
			outputFileListPaths = (
			);
			outputPaths = (
				"$(BUILT_PRODUCTS_DIR)/libcounter.a",
				"$(BUILT_PRODUCTS_DIR)/libGoNativeKit.a",
			);
			runOnlyForDeploymentPostprocessing = 0;
			shellPath = /bin/sh;
			shellScript = "export PATH=\"$PATH:/opt/homebrew/bin:/usr/local/bin:$HOME/go/bin\"\nexport GOCACHE=\"${SRCROOT}/../build/gocache\"\nFRAMEWORK_ROOT=$(cd \"${SRCROOT}/..\" && go list -m -f '{{.Dir}}' github.com/go-native/go-native)\n\"$FRAMEWORK_ROOT/scripts/build-ios-framework.sh\" >/dev/null\n\nif [ \"$PLATFORM_NAME\" = \"iphonesimulator\" ]; then\n    SDK_PATH=$(xcrun --sdk iphonesimulator --show-sdk-path)\n    SLICE=ios-arm64-simulator\n    export GOOS=ios\n    export GOARCH=arm64\n    export CGO_ENABLED=1\n    export CC=\"clang -target arm64-apple-ios15.0-simulator -isysroot $SDK_PATH\"\nelse\n    SDK_PATH=$(xcrun --sdk iphoneos --show-sdk-path)\n    SLICE=ios-arm64\n    export GOOS=ios\n    export GOARCH=arm64\n    export CGO_ENABLED=1\n    export CC=\"clang -target arm64-apple-ios15.0 -isysroot $SDK_PATH\"\nfi\n\nmkdir -p \"${BUILT_PRODUCTS_DIR}/GoNativeKitHeaders\"\ncp \"$FRAMEWORK_ROOT/build/native/GoNativeKit.xcframework/$SLICE/libGoNativeKit.a\" \"${BUILT_PRODUCTS_DIR}/libGoNativeKit.a\"\ncp -R \"$FRAMEWORK_ROOT/build/native/GoNativeKit.xcframework/$SLICE/Headers/.\" \"${BUILT_PRODUCTS_DIR}/GoNativeKitHeaders/\"\ncd \"${SRCROOT}/..\"\ngo build -buildmode=c-archive -o \"${BUILT_PRODUCTS_DIR}/libcounter.a\" ./ios/bridge\n";
		};
/* End PBXShellScriptBuildPhase section */

/* Begin PBXSourcesBuildPhase section */
		100000000000000000000005 /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				100000000000000000000014 /* AppDelegate.m in Sources */,
				100000000000000000000011 /* main.m in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXSourcesBuildPhase section */

/* Begin XCBuildConfiguration section */
		10000000000000000000000B /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_ANALYZER_NONNULL = YES;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++20";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CODE_SIGN_STYLE = Automatic;
				DEAD_CODE_STRIPPING = YES;
				DEBUG_INFORMATION_FORMAT = dwarf;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				ENABLE_TESTABILITY = YES;
				GCC_DYNAMIC_NO_PIC = NO;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_OPTIMIZATION_LEVEL = 0;
				GCC_PREPROCESSOR_DEFINITIONS = (
					"DEBUG=1",
					"$(inherited)",
				);
				IPHONEOS_DEPLOYMENT_TARGET = 15.0;
				MTL_ENABLE_DEBUG_INFO = INCLUDE_SOURCE;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = iphoneos;
			};
			name = Debug;
		};
		10000000000000000000000C /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_ANALYZER_NONNULL = YES;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++20";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CODE_SIGN_STYLE = Automatic;
				DEAD_CODE_STRIPPING = YES;
				DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
				ENABLE_NS_ASSERTIONS = NO;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				GCC_NO_COMMON_BLOCKS = YES;
				IPHONEOS_DEPLOYMENT_TARGET = 15.0;
				MTL_ENABLE_DEBUG_INFO = NO;
				SDKROOT = iphoneos;
				VALIDATE_PRODUCT = YES;
			};
			name = Release;
		};
		10000000000000000000000D /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ARCHS = arm64;
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_STYLE = Automatic;
				CURRENT_PROJECT_VERSION = 1;
				GENERATE_INFOPLIST_FILE = NO;
				HEADER_SEARCH_PATHS = (
					"$(BUILT_PRODUCTS_DIR)",
					"$(BUILT_PRODUCTS_DIR)/GoNativeKitHeaders",
				);
				INFOPLIST_FILE = Info.plist;
				IPHONEOS_DEPLOYMENT_TARGET = 15.0;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				LIBRARY_SEARCH_PATHS = (
					"$(inherited)",
					"$(BUILT_PRODUCTS_DIR)",
				);
				MARKETING_VERSION = 0.1;
				OTHER_LDFLAGS = (
					"-lcounter",
					"-lGoNativeKit",
					"-framework",
					"UIKit",
					"-framework",
					"Foundation",
					"-framework",
					"CoreGraphics",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "dev.gonative.%s";
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_EMIT_LOC_STRINGS = YES;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Debug;
		};
		10000000000000000000000E /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ARCHS = arm64;
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_STYLE = Automatic;
				CURRENT_PROJECT_VERSION = 1;
				GENERATE_INFOPLIST_FILE = NO;
				HEADER_SEARCH_PATHS = (
					"$(BUILT_PRODUCTS_DIR)",
					"$(BUILT_PRODUCTS_DIR)/GoNativeKitHeaders",
				);
				INFOPLIST_FILE = Info.plist;
				IPHONEOS_DEPLOYMENT_TARGET = 15.0;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				LIBRARY_SEARCH_PATHS = (
					"$(inherited)",
					"$(BUILT_PRODUCTS_DIR)",
				);
				MARKETING_VERSION = 0.1;
				OTHER_LDFLAGS = (
					"-lcounter",
					"-lGoNativeKit",
					"-framework",
					"UIKit",
					"-framework",
					"Foundation",
					"-framework",
					"CoreGraphics",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "dev.gonative.%s";
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_EMIT_LOC_STRINGS = YES;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Release;
		};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		100000000000000000000009 /* Build configuration list for PBXProject "%s" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				10000000000000000000000B /* Debug */,
				10000000000000000000000C /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		10000000000000000000000A /* Build configuration list for PBXNativeTarget "%s" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				10000000000000000000000D /* Debug */,
				10000000000000000000000E /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
/* End XCConfigurationList section */

	};
	rootObject = 100000000000000000000001 /* Project object */;
}
`, name, name, name, name, name, name, name, name, name, name, pkg, pkg, name, name)
}

func generateXcscheme(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Scheme
   LastUpgradeVersion = "1500"
   version = "1.7">
   <BuildAction
      parallelizeBuildables = "YES"
      buildImplicitDependencies = "YES">
      <BuildActionEntries>
         <BuildActionEntry
            buildForTesting = "YES"
            buildForRunning = "YES"
            buildForProfiling = "YES"
            buildForArchiving = "YES"
            buildForAnalyzing = "YES">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "100000000000000000000002"
               BuildableName = "%s.app"
               BlueprintName = "%s"
               ReferencedContainer = "container:%s.xcodeproj">
            </BuildableReference>
         </BuildActionEntry>
      </BuildActionEntries>
   </BuildAction>
   <TestAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      shouldUseLaunchSchemeArgsEnv = "YES">
      <Testables>
      </Testables>
   </TestAction>
   <LaunchAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      launchStyle = "0"
      useCustomWorkingDirectory = "NO"
      ignoresPersistentStateOnLaunch = "NO"
      debugDocumentVersioning = "YES"
      debugServiceExtension = "internal"
      allowLocationSimulation = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "100000000000000000000002"
            BuildableName = "%s.app"
            BlueprintName = "%s"
            ReferencedContainer = "container:%s.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
   </LaunchAction>
   <ProfileAction
      buildConfiguration = "Release"
      shouldUseLaunchSchemeArgsEnv = "YES"
      savedToolIdentifier = ""
      useCustomWorkingDirectory = "NO"
      debugDocumentVersioning = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "100000000000000000000002"
            BuildableName = "%s.app"
            BlueprintName = "%s"
            ReferencedContainer = "container:%s.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
   </ProfileAction>
   <AnalyzeAction
      buildConfiguration = "Debug">
   </AnalyzeAction>
   <ArchiveAction
      buildConfiguration = "Release"
      revealArchiveInOrganizer = "YES">
   </ArchiveAction>
</Scheme>
`, name, name, name, name, name, name, name, name, name)
}

func getProjectTemplates(name string) map[string]string {
	pkg := sanitizePackageName(name)
	iosPkg := sanitizeBundleIdentifierPart(name)

	templates := map[string]string{
		"go.mod": fmt.Sprintf("module %s\n\ngo 1.24\n\nrequire github.com/go-native/go-native v0.0.0\n", name),
		"gonative.yaml": `framework:
  version: 0.1.0
  protocol: 10
  measurementProtocol: 2

ios:
  deploymentTarget: "15.0"
  package: GoNativeKit

android:
  minSdk: 23
  targetSdk: 35
  dependency: dev.gonative:gonative-runtime:0.1.0
`,
		"app.go": `// Package app contains the application's declarative native UI.
package app

import "github.com/go-native/go-native/ui"

// App builds the root UI component.
func App() ui.Component {
	return ui.SafeArea(
		ui.Column(
			ui.Text("Hello from Go Native").FontSize(28).Bold(),
			ui.Text("Edit app.go to get started."),
		).Padding(20).Gap(12).Align(ui.AlignCenter),
	)
}
`,
		"assets/.gitkeep": "",
		".gitignore":      "build/\n.gonative/\n.gradle/\n*.app\n*.apk\n*.idsig\nDerivedData/\nandroid/app/libs/*.aar\n",
		"README.md": fmt.Sprintf(`# %s

A Go Native application. The UI is declared in app.go and renders genuine platform-native controls on iOS and Android.

## Project Map
- app.go: Declarative UI tree written in Go.
- ios/: Native iOS host project (UIKit, Xcode-compatible).
- android/: Native Android host project (Android Views + Gradle, Android Studio-compatible).

## Development Commands

`+"```bash"+`
# Check local toolchain
gonative doctor

# Build & run on iOS Simulator
gonative build ios
gonative run ios

# Build & run on Android
gonative build android
gonative run android
`+"```"+`

## IDE Usage
- **Xcode**: Open `+"`ios/%s.xcodeproj`"+` in Xcode and click **Run** (Cmd+R).
- **Android Studio**: Open the `+"`android/`"+` directory in Android Studio and click **Run** (Shift+F10).
`, name, name),

		// iOS Host & Xcode Project
		fmt.Sprintf("ios/%s.xcodeproj/project.pbxproj", name):                          generatePbxproj(name, iosPkg),
		fmt.Sprintf("ios/%s.xcodeproj/xcshareddata/xcschemes/%s.xcscheme", name, name): generateXcscheme(name),
		"ios/AppDelegate.h": `#import <UIKit/UIKit.h>

@interface AppDelegate : UIResponder <UIApplicationDelegate>

@property (strong, nonatomic) UIWindow *window;

@end
`,
		"ios/Package.swift": `// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "GoNativeAppDependencies",
    platforms: [.iOS(.v15)],
    products: [
        .library(name: "GoNativeAppDependencies", targets: ["GoNativeKit"]),
    ],
    targets: [
        .binaryTarget(name: "GoNativeKit", path: ".gonative/GoNativeKit.xcframework"),
    ]
)
`,
		"ios/AppDelegate.m": `#import "AppDelegate.h"
#import "GoNativeRenderer.h"

@implementation AppDelegate

- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    self.window.rootViewController = [GNRootViewController new];
    [self.window makeKeyAndVisible];
    return YES;
}

@end
`,
		"ios/main.m": `#import <UIKit/UIKit.h>
#import "AppDelegate.h"

int main(int argc, char * argv[]) {
    @autoreleasepool {
        return UIApplicationMain(argc, argv, nil, NSStringFromClass([AppDelegate class]));
    }
}
`,
		"ios/Info.plist": fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleDisplayName</key>
    <string>%s</string>
    <key>CFBundleExecutable</key>
    <string>%s</string>
    <key>CFBundleIdentifier</key>
    <string>dev.gonative.%s</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>%s</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>0.1</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSRequiresIPhoneOS</key>
    <true/>
    <key>UILaunchScreen</key>
    <dict/>
    <key>UIUserInterfaceStyle</key>
    <string>Light</string>
    <key>UISupportedInterfaceOrientations</key>
    <array>
        <string>UIInterfaceOrientationPortrait</string>
    </array>
</dict>
</plist>
`, name, name, iosPkg, name),

		"ios/bridge/main.go": fmt.Sprintf(`package main

/*
#include <stdint.h>
#include <stdlib.h>
void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
int32_t GNMeasureNativeBatch(const uint8_t *bytes, int32_t length, uint8_t **results, int32_t *resultLength);
void GNFreeNativeBuffer(void *buffer);
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/runtime/layout"
	"github.com/go-native/go-native/ui"
	"%s"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

var benchmarkOutput string
var goNativeReloadSession string

type iosRenderer struct{}

type iosNativeMeasurer struct{}

func (iosNativeMeasurer) MeasureBatch(ctx context.Context, requests []layout.MeasurementRequest) ([]layout.MeasurementResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := layout.MarshalMeasurementRequests(requests)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("empty native measurement request")
	}
	var output *C.uint8_t
	var outputLength C.int32_t
	status := C.GNMeasureNativeBatch((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int32_t(len(data)), &output, &outputLength)
	if output != nil {
		defer C.GNFreeNativeBuffer(unsafe.Pointer(output))
	}
	if status != 0 {
		return nil, fmt.Errorf("UIKit measurement failed with status %%d", int32(status))
	}
	if output == nil || outputLength <= 0 {
		return nil, errors.New("UIKit measurement returned an empty response")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return layout.UnmarshalMeasurementResults(C.GoBytes(unsafe.Pointer(output), C.int(outputLength)))
}

var _ layout.BatchMeasurer = iosNativeMeasurer{}

func (iosRenderer) Apply(batch gnruntime.MutationBatch) error {
	data, err := batch.MarshalBinary()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty mutation batch")
	}
	C.GNApplyMutationBatch((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
	return nil
}

var appRuntime *gnruntime.Runtime

//export GoNativeStart
func GoNativeStart() {
	configureReloadState()
	appRuntime = gnruntime.New(app.App, iosRenderer{})
	appRuntime.SetLayoutProvider(&layout.Pipeline{Measurer: iosNativeMeasurer{}, Cache: layout.NewMeasurementCache()})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
}

func configureReloadState() {
	if goNativeReloadSession == "" {
		return
	}
	directory, err := os.UserCacheDir()
	if err != nil {
		directory = os.TempDir()
	}
	ui.ConfigureReloadState(ui.ReloadStateOptions{Path: filepath.Join(directory, "go-native", "reload-state.json"), SessionID: goNativeReloadSession, OnWarning: func(message string) { fmt.Println("Go Native Fast Reload:", message) }})
}

//export GoNativeSetViewport
func GoNativeSetViewport(width, height, scale C.float) {
	if appRuntime == nil || width <= 0 || height <= 0 {
		return
	}
	current := appRuntime.Environment().MediaQuery
	if current.Viewport.Width == float32(width) && current.Viewport.Height == float32(height) && current.Scale == float32(scale) {
		return
	}
	appRuntime.UpdateEnvironment(func(environment ui.Environment) ui.Environment {
		environment.MediaQuery.Viewport = ui.Size{Width: float32(width), Height: float32(height)}
		environment.MediaQuery.Scale = float32(scale)
		return environment
	})
}

//export GoNativeDispatchEvent
func GoNativeDispatchEvent(handler C.uint64_t) {
	if appRuntime != nil {
		appRuntime.Dispatch(ui.HandlerID(handler))
	}
}

//export GoNativeDispatchValueEvent
func GoNativeDispatchValueEvent(handler C.uint64_t, value *C.char) {
	if appRuntime != nil {
		appRuntime.DispatchValue(ui.HandlerID(handler), C.GoString(value))
	}
}

//export GoNativeDispatchSelection
func GoNativeDispatchSelection(handler C.uint64_t, start, end C.int32_t) {
	if appRuntime != nil {
		appRuntime.DispatchSelection(ui.HandlerID(handler), int32(start), int32(end))
	}
}

//export GoNativeDispatchBoolEvent
func GoNativeDispatchBoolEvent(handler C.uint64_t, value C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchBool(ui.HandlerID(handler), value != 0)
	}
}

//export GoNativeDispatchGestureEvent
func GoNativeDispatchGestureEvent(handler C.uint64_t, translationX, translationY, velocityX, velocityY C.float) {
	if appRuntime != nil {
		appRuntime.DispatchGesture(ui.HandlerID(handler), ui.GestureEvent{
			TranslationX: float32(translationX), TranslationY: float32(translationY),
			VelocityX: float32(velocityX), VelocityY: float32(velocityY),
		})
	}
}

//export GoNativeDispatchFocus
func GoNativeDispatchFocus(nodeID C.uint64_t, focused C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchFocus(ui.NodeID(nodeID), focused != 0)
	}
}

//export GoNativeStop
func GoNativeStop() {
	if appRuntime != nil {
		appRuntime.Stop()
		appRuntime = nil
	}
}

//export GoNativeSetLifecycle
func GoNativeSetLifecycle(state C.uint8_t) {
	if appRuntime != nil && state <= C.uint8_t(ui.LifecycleDestroyed) {
		appRuntime.SetLifecycle(ui.LifecycleState(state))
	}
}

//export GoNativeReportBatchApplied
func GoNativeReportBatchApplied(sequence C.uint64_t, nativeNanos C.uint64_t) {
	if appRuntime != nil {
		appRuntime.RecordNativeApply(uint64(sequence), time.Duration(nativeNanos))
		emitTimingSample(uint64(sequence))
	}
}

func emitTimingSample(sequence uint64) {
	if benchmarkOutput != "1" {
		return
	}
	for _, sample := range appRuntime.TimingSamples() {
		if sample.Sequence == sequence {
			fmt.Printf("GONATIVE_TIMING {\"sequence\":%%d,\"mutations\":%%d,\"native_apply_ns\":%%d,\"bridge_to_apply_ns\":%%d,\"event_to_apply_ns\":%%d}\n", sample.Sequence, sample.MutationCount, sample.NativeApply.Nanoseconds(), sample.BridgeToApply.Nanoseconds(), sample.EventToApply.Nanoseconds())
			return
		}
	}
}

func main() {}
`, name),

		// Android Host & Android Studio Project
		"android/build.gradle": `plugins {
    id "com.android.application" version "8.9.1" apply false
}
`,

		"android/settings.gradle": fmt.Sprintf(`pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}
dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        maven { url = uri("${rootDir}/.gonative/m2") }
        google()
        mavenCentral()
    }
}
rootProject.name = '%s'
include ':app'
`, name),

		"android/gradle.properties": "android.useAndroidX=true\n",

		"android/app/build.gradle": fmt.Sprintf(`plugins {
    id "com.android.application"
}

def selectedAbis = (System.getenv("GONATIVE_ANDROID_ABIS") ?: "arm64-v8a,x86_64").split(",")
def rootDir = rootProject.projectDir.parentFile

android {
    namespace "dev.gonative.%s"
    compileSdk 35

    defaultConfig {
        applicationId "dev.gonative.%s"
        minSdk 23
        targetSdk 35
        versionCode 1
        versionName "0.1"
        ndk {
            abiFilters(*selectedAbis)
        }
    }

    sourceSets {
        main {
            jniLibs.srcDirs = ["../../build/android/lib"]
        }
    }

    buildTypes {
        debug {
            debuggable true
        }
        release {
            minifyEnabled false
        }
    }

    compileOptions {
        sourceCompatibility JavaVersion.VERSION_1_8
        targetCompatibility JavaVersion.VERSION_1_8
    }
}

dependencies {
    implementation "dev.gonative:gonative-runtime:0.1.0"
}

tasks.register("prepareGoNativeLibraries", Exec) {
    workingDir rootProject.projectDir.parentFile
    commandLine "sh", "${rootProject.projectDir}/build-libs.sh"
}

tasks.named("preBuild").configure {
    dependsOn("prepareGoNativeLibraries")
}
`, pkg, pkg),

		"android/build-libs.sh": `#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
SDK=${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/Library/Android/sdk}}
NDK_VERSION=${GONATIVE_NDK_VERSION:-28.2.13676358}
NDK="$SDK/ndk/$NDK_VERSION"
if [ ! -d "$NDK" ]; then
    for d in "$SDK/ndk/"* "$SDK/ndk-bundle"; do
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
        if [ -f "$BUILD/lib/$abi/libgonative.so" ]; then
            echo "Using existing pre-built $BUILD/lib/$abi/libgonative.so"
            continue
        fi
        echo "Missing Android NDK compiler: $TOOLCHAIN/bin/$compiler" >&2
        echo "Please set ANDROID_NDK_ROOT or install NDK via Android Studio SDK Manager." >&2
        exit 1
    fi
    mkdir -p "$LIB_BUILD/$abi"
    (
        cd "$ROOT"
        reload_ldflags=
        if [ -n "${GONATIVE_RELOAD_SESSION:-}" ]; then reload_ldflags="-X main.goNativeReloadSession=$GONATIVE_RELOAD_SESSION"; fi
        CGO_ENABLED=1 GOOS=android GOARCH="$goarch" \
        CC="$TOOLCHAIN/bin/$compiler" \
        CGO_CFLAGS="--sysroot=$TOOLCHAIN/sysroot -I$TOOLCHAIN/sysroot/usr/include" \
        go build -ldflags "$reload_ldflags" -buildmode=c-shared -o "$LIB_BUILD/$abi/libgonative.so" ./android/bridge
    )
    rm -f "$LIB_BUILD/$abi/libgonative.h"
done
IFS=$old_ifs

rm -rf "$BUILD/lib"
mv "$LIB_BUILD" "$BUILD/lib"
trap - EXIT INT TERM
`,

		"android/app/src/main/AndroidManifest.xml": `<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android" android:versionCode="1" android:versionName="0.1">
    <application android:theme="@style/AppTheme" android:label="@string/app_name" android:allowBackup="false" android:supportsRtl="true">
        <activity android:name=".MainActivity" android:screenOrientation="portrait" android:exported="true">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>
`,

		"android/app/src/main/res/values/styles.xml": `<?xml version="1.0" encoding="utf-8"?>
<resources>
    <style name="AppTheme" parent="android:style/Theme.Material.Light.NoActionBar">
        <item name="android:fontFamily">sans</item>
        <item name="android:colorAccent">#0066CC</item>
        <item name="android:windowBackground">#FFFFFF</item>
        <item name="android:textColor">#000000</item>
        <item name="android:textColorPrimary">#000000</item>
        <item name="android:editTextColor">#000000</item>
        <item name="android:windowLightStatusBar">true</item>
        <item name="android:navigationBarColor">#FFFFFF</item>
        <item name="android:statusBarColor">#FFFFFF</item>
    </style>
</resources>
`,

		"android/app/src/main/res/values/strings.xml": fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</string>
</resources>
`, name),

		fmt.Sprintf("android/app/src/main/java/dev/gonative/%s/MainActivity.java", pkg): fmt.Sprintf(`package dev.gonative.%s;

import dev.gonative.runtime.GoNativeActivity;

/** Application launcher. Rendering is supplied by gonative-runtime. */
public final class MainActivity extends GoNativeActivity {}
`, pkg),

		"android/bridge/main.go": fmt.Sprintf(`//go:build android

package main

/*
#include <stdint.h>
#include <stdlib.h>
void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length);
int32_t GNAndroidMeasureBatch(const uint8_t *bytes, int32_t length, uint8_t **result);
void GNAndroidFreeBuffer(uint8_t *bytes);
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unsafe"

	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/runtime/layout"
	"github.com/go-native/go-native/ui"
	"%s"
	"os"
	"path/filepath"
)

var benchmarkOutput string
var goNativeReloadSession string

type androidRenderer struct{}

func (androidRenderer) Apply(batch gnruntime.MutationBatch) error {
	data, err := batch.MarshalBinary()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty mutation batch")
	}
	C.GNAndroidApplyMutationBatch((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
	return nil
}

func (androidRenderer) MeasureBatch(ctx context.Context, requests []layout.MeasurementRequest) ([]layout.MeasurementResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	payload, err := layout.MarshalMeasurementRequests(requests)
	if err != nil {
		return nil, err
	}
	if len(payload) == 0 {
		return nil, errors.New("android measurement: empty request batch")
	}
	var result *C.uint8_t
	length := C.GNAndroidMeasureBatch((*C.uint8_t)(unsafe.Pointer(&payload[0])), C.int32_t(len(payload)), &result)
	if length <= 0 || result == nil {
		return nil, fmt.Errorf("android measurement: native adapter unavailable (%%d)", int32(length))
	}
	defer C.GNAndroidFreeBuffer(result)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return layout.UnmarshalMeasurementResults(C.GoBytes(unsafe.Pointer(result), C.int(length)))
}

var appRuntime *gnruntime.Runtime

//export GoNativeAndroidStart
func GoNativeAndroidStart() {
	configureReloadState()
	renderer := androidRenderer{}
	appRuntime = gnruntime.New(app.App, renderer)
	appRuntime.SetLayoutProvider(&layout.Pipeline{Measurer: renderer, Cache: layout.NewMeasurementCache()})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
}

func configureReloadState() {
	if goNativeReloadSession == "" {
		return
	}
	directory, err := os.UserCacheDir()
	if err != nil {
		directory = os.TempDir()
	}
	ui.ConfigureReloadState(ui.ReloadStateOptions{Path: filepath.Join(directory, "go-native", "reload-state.json"), SessionID: goNativeReloadSession, OnWarning: func(message string) { fmt.Println("Go Native Fast Reload:", message) }})
}

//export GoNativeAndroidUpdateViewport
func GoNativeAndroidUpdateViewport(width, height, scale C.float) {
	if appRuntime == nil || width <= 0 || height <= 0 || scale <= 0 {
		return
	}
	current := appRuntime.Environment().MediaQuery
	if current.Viewport.Width == float32(width) && current.Viewport.Height == float32(height) && current.Scale == float32(scale) {
		return
	}
	appRuntime.UpdateEnvironment(func(environment ui.Environment) ui.Environment {
		environment.MediaQuery.Viewport = ui.Size{Width: float32(width), Height: float32(height)}
		environment.MediaQuery.Scale = float32(scale)
		if width > height {
			environment.MediaQuery.Orientation = ui.OrientationLandscape
		} else {
			environment.MediaQuery.Orientation = ui.OrientationPortrait
		}
		return environment
	})
}

//export GoNativeAndroidSetLifecycle
func GoNativeAndroidSetLifecycle(state C.uint8_t) {
	if appRuntime != nil {
		appRuntime.SetLifecycle(ui.LifecycleState(state))
	}
}

//export GoNativeAndroidDispatchFocus
func GoNativeAndroidDispatchFocus(nodeID C.uint64_t, focused C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchFocus(ui.NodeID(nodeID), focused != 0)
	}
}

//export GoNativeAndroidDispatchEvent
func GoNativeAndroidDispatchEvent(handler C.uint64_t) {
	if appRuntime != nil {
		appRuntime.Dispatch(ui.HandlerID(handler))
	}
}

//export GoNativeAndroidDispatchValueEvent
func GoNativeAndroidDispatchValueEvent(handler C.uint64_t, value *C.char) {
	if appRuntime != nil {
		appRuntime.DispatchValue(ui.HandlerID(handler), C.GoString(value))
	}
}

//export GoNativeAndroidDispatchBoolEvent
func GoNativeAndroidDispatchBoolEvent(handler C.uint64_t, value C.uint8_t) {
	if appRuntime != nil {
		appRuntime.DispatchBool(ui.HandlerID(handler), value != 0)
	}
}

//export GoNativeAndroidDispatchGestureEvent
func GoNativeAndroidDispatchGestureEvent(handler C.uint64_t, translationX, translationY, velocityX, velocityY C.float) {
	if appRuntime != nil {
		appRuntime.DispatchGesture(ui.HandlerID(handler), ui.GestureEvent{
			TranslationX: float32(translationX), TranslationY: float32(translationY),
			VelocityX: float32(velocityX), VelocityY: float32(velocityY),
		})
	}
}

//export GoNativeAndroidDispatchSelectionEvent
func GoNativeAndroidDispatchSelectionEvent(handler C.uint64_t, start, end C.int32_t) {
	if appRuntime != nil {
		appRuntime.DispatchSelection(ui.HandlerID(handler), int32(start), int32(end))
	}
}

//export GoNativeAndroidStop
func GoNativeAndroidStop() {
	if appRuntime != nil {
		appRuntime.Stop()
		appRuntime = nil
	}
}

//export GoNativeAndroidReportBatchApplied
func GoNativeAndroidReportBatchApplied(sequence C.uint64_t, nativeNanos C.uint64_t) {
	if appRuntime != nil {
		appRuntime.RecordNativeApply(uint64(sequence), time.Duration(nativeNanos))
		emitTimingSample(uint64(sequence))
	}
}

func emitTimingSample(sequence uint64) {
	if benchmarkOutput != "1" {
		return
	}
	for _, sample := range appRuntime.TimingSamples() {
		if sample.Sequence == sequence {
			fmt.Printf("GONATIVE_TIMING {\"sequence\":%%d,\"mutations\":%%d,\"native_apply_ns\":%%d,\"bridge_to_apply_ns\":%%d,\"event_to_apply_ns\":%%d}\n", sample.Sequence, sample.MutationCount, sample.NativeApply.Nanoseconds(), sample.BridgeToApply.Nanoseconds(), sample.EventToApply.Nanoseconds())
			return
		}
	}
}

func main() {}
`, name),

		"android/bridge/jni.c": `//go:build android

#include <jni.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>

extern void GoNativeAndroidStart(void);
extern void GoNativeAndroidDispatchEvent(uint64_t handler);
extern void GoNativeAndroidDispatchValueEvent(uint64_t handler, const char *value);
extern void GoNativeAndroidDispatchBoolEvent(uint64_t handler, uint8_t value);
extern void GoNativeAndroidDispatchGestureEvent(uint64_t handler, float translationX, float translationY, float velocityX, float velocityY);
extern void GoNativeAndroidDispatchSelectionEvent(uint64_t handler, int32_t start, int32_t end);
extern void GoNativeAndroidStop(void);
extern void GoNativeAndroidReportBatchApplied(uint64_t sequence, uint64_t nativeNanos);
extern void GoNativeAndroidSetLifecycle(uint8_t state);
extern void GoNativeAndroidDispatchFocus(uint64_t nodeID, uint8_t focused);
extern void GoNativeAndroidUpdateViewport(float width, float height, float scale);

static JavaVM *gn_vm;
static jobject gn_renderer;
static jmethodID gn_apply;
static jmethodID gn_measure;
static pthread_mutex_t gn_renderer_mu = PTHREAD_MUTEX_INITIALIZER;

JNIEXPORT jint JNICALL JNI_OnLoad(JavaVM *vm, void *reserved) {
    (void)reserved;
    gn_vm = vm;
    return JNI_VERSION_1_6;
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeStart(JNIEnv *env, jobject renderer) {
    pthread_mutex_lock(&gn_renderer_mu);
    if (gn_renderer) {
        (*env)->DeleteGlobalRef(env, gn_renderer);
    }
    gn_renderer = (*env)->NewGlobalRef(env, renderer);
    jclass cls = (*env)->GetObjectClass(env, renderer);
    gn_apply = (*env)->GetMethodID(env, cls, "applyMutationBatch", "([B)V");
    gn_measure = (*env)->GetMethodID(env, cls, "measureNativeBatch", "([B)[B");
    (*env)->DeleteLocalRef(env, cls);
    pthread_mutex_unlock(&gn_renderer_mu);
    GoNativeAndroidStart();
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeSetLifecycle(JNIEnv *env, jobject renderer, jint state) {
    (void)env;
    (void)renderer;
    if (state >= 0 && state <= 6) GoNativeAndroidSetLifecycle((uint8_t)state);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeDispatchFocus(JNIEnv *env, jobject renderer, jlong nodeID, jboolean focused) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchFocus((uint64_t)nodeID, focused ? 1 : 0);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeUpdateViewport(JNIEnv *env, jobject renderer, jfloat width, jfloat height, jfloat scale) {
    (void)env;
    (void)renderer;
    GoNativeAndroidUpdateViewport((float)width, (float)height, (float)scale);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeDispatchEvent(JNIEnv *env, jobject renderer, jlong handler) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchEvent((uint64_t)handler);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeDispatchValueEvent(JNIEnv *env, jobject renderer, jlong handler, jstring value) {
    (void)renderer;
    const char *utf8 = value ? (*env)->GetStringUTFChars(env, value, NULL) : "";
    if (utf8) { GoNativeAndroidDispatchValueEvent((uint64_t)handler, utf8); }
    if (value && utf8) { (*env)->ReleaseStringUTFChars(env, value, utf8); }
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeDispatchBoolEvent(JNIEnv *env, jobject renderer, jlong handler, jboolean value) {
    (void)env; (void)renderer; GoNativeAndroidDispatchBoolEvent((uint64_t)handler, value ? 1 : 0);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeDispatchGestureEvent(JNIEnv *env, jobject renderer, jlong handler, jfloat translationX, jfloat translationY, jfloat velocityX, jfloat velocityY) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchGestureEvent((uint64_t)handler, (float)translationX, (float)translationY, (float)velocityX, (float)velocityY);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeDispatchSelectionEvent(JNIEnv *env, jobject renderer, jlong handler, jint start, jint end) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchSelectionEvent((uint64_t)handler, (int32_t)start, (int32_t)end);
}

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeStop(JNIEnv *env, jobject renderer) {
    (void)renderer;
    GoNativeAndroidStop();
    pthread_mutex_lock(&gn_renderer_mu);
    if (gn_renderer) {
        (*env)->DeleteGlobalRef(env, gn_renderer);
        gn_renderer = NULL;
    }
    gn_apply = NULL;
    gn_measure = NULL;
    pthread_mutex_unlock(&gn_renderer_mu);
}

int32_t GNAndroidMeasureBatch(const uint8_t *bytes, int32_t length, uint8_t **result) {
    if (result) *result = NULL;
    if (!gn_vm || !result || !bytes || length <= 0) return -1;
    JNIEnv *env = NULL;
    int attached = 0;
    jint status = (*gn_vm)->GetEnv(gn_vm, (void **)&env, JNI_VERSION_1_6);
    if (status == JNI_EDETACHED) {
        if ((*gn_vm)->AttachCurrentThread(gn_vm, &env, NULL) != JNI_OK) return -2;
        attached = 1;
    } else if (status != JNI_OK) return -2;
    int32_t result_length = -3;
    pthread_mutex_lock(&gn_renderer_mu);
    if (gn_renderer && gn_measure) {
        jbyteArray request = (*env)->NewByteArray(env, length);
        if (request) {
            (*env)->SetByteArrayRegion(env, request, 0, length, (const jbyte *)bytes);
            jbyteArray response = (jbyteArray)(*env)->CallObjectMethod(env, gn_renderer, gn_measure, request);
            if (!(*env)->ExceptionCheck(env) && response) {
                jsize response_length = (*env)->GetArrayLength(env, response);
                if (response_length > 0 && response_length <= 16777216) {
                    uint8_t *copy = (uint8_t *)malloc((size_t)response_length);
                    if (copy) {
                        (*env)->GetByteArrayRegion(env, response, 0, response_length, (jbyte *)copy);
                        *result = copy;
                        result_length = (int32_t)response_length;
                    } else result_length = -4;
                }
                (*env)->DeleteLocalRef(env, response);
            }
            (*env)->DeleteLocalRef(env, request);
        }
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
            result_length = -5;
        }
    }
    pthread_mutex_unlock(&gn_renderer_mu);
    if (attached) (*gn_vm)->DetachCurrentThread(gn_vm);
    return result_length;
}

void GNAndroidFreeBuffer(uint8_t *bytes) { free(bytes); }

JNIEXPORT void JNICALL
Java_dev_gonative_runtime_GoNativeActivity_nativeReportBatchApplied(JNIEnv *env, jobject renderer, jlong sequence, jlong nativeNanos) {
    (void)env;
    (void)renderer;
    GoNativeAndroidReportBatchApplied((uint64_t)sequence, (uint64_t)nativeNanos);
}

void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length) {
    if (!gn_vm || length <= 0) {
        return;
    }
    JNIEnv *env = NULL;
    int attached = 0;
    jint status = (*gn_vm)->GetEnv(gn_vm, (void **)&env, JNI_VERSION_1_6);
    if (status == JNI_EDETACHED) {
        if ((*gn_vm)->AttachCurrentThread(gn_vm, &env, NULL) != JNI_OK) {
            return;
        }
        attached = 1;
    } else if (status != JNI_OK) {
        return;
    }
    pthread_mutex_lock(&gn_renderer_mu);
    if (!gn_renderer || !gn_apply) {
        pthread_mutex_unlock(&gn_renderer_mu);
        if (attached) (*gn_vm)->DetachCurrentThread(gn_vm);
        return;
    }
    jbyteArray payload = (*env)->NewByteArray(env, length);
    if (payload) {
        (*env)->SetByteArrayRegion(env, payload, 0, length, (const jbyte *)bytes);
        (*env)->CallVoidMethod(env, gn_renderer, gn_apply, payload);
        (*env)->DeleteLocalRef(env, payload);
    }
    if ((*env)->ExceptionCheck(env)) {
        (*env)->ExceptionDescribe(env);
        (*env)->ExceptionClear(env);
    }
    pthread_mutex_unlock(&gn_renderer_mu);
    if (attached) {
        (*gn_vm)->DetachCurrentThread(gn_vm);
    }
}
`,

		"android/bridge/stub.go": `//go:build !android

package main

func main() {}
`,
	}

	return templates
}

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

func sanitizeJniName(pkg string) string {
	return strings.ReplaceAll(pkg, "_", "_1")
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
		100000000000000000000011 /* main.m in Sources */ = {isa = PBXBuildFile; fileRef = 100000000000000000000010 /* main.m */; };
		100000000000000000000014 /* GoNativeRenderer.m in Sources */ = {isa = PBXBuildFile; fileRef = 100000000000000000000013 /* GoNativeRenderer.m */; };
/* End PBXBuildFile section */

/* Begin PBXFileReference section */
		100000000000000000000003 /* %s.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "%s.app"; sourceTree = BUILT_PRODUCTS_DIR; };
		100000000000000000000010 /* main.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = main.m; sourceTree = "<group>"; };
		100000000000000000000012 /* GoNativeRenderer.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = GoNativeRenderer.h; sourceTree = "<group>"; };
		100000000000000000000013 /* GoNativeRenderer.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = GoNativeRenderer.m; sourceTree = "<group>"; };
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
				100000000000000000000010 /* main.m */,
				100000000000000000000012 /* GoNativeRenderer.h */,
				100000000000000000000013 /* GoNativeRenderer.m */,
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
			);
			runOnlyForDeploymentPostprocessing = 0;
			shellPath = /bin/sh;
			shellScript = "export PATH=\"$PATH:/opt/homebrew/bin:/usr/local/bin:$HOME/go/bin\"\nexport GOCACHE=\"${SRCROOT}/../build/gocache\"\n\nif [ \"$PLATFORM_NAME\" = \"iphonesimulator\" ]; then\n    SDK_PATH=$(xcrun --sdk iphonesimulator --show-sdk-path)\n    export GOOS=ios\n    export GOARCH=arm64\n    export CGO_ENABLED=1\n    export CC=\"clang -target arm64-apple-ios15.0-simulator -isysroot $SDK_PATH\"\nelse\n    SDK_PATH=$(xcrun --sdk iphoneos --show-sdk-path)\n    export GOOS=ios\n    export GOARCH=arm64\n    export CGO_ENABLED=1\n    export CC=\"clang -target arm64-apple-ios15.0 -isysroot $SDK_PATH\"\nfi\n\nmkdir -p \"${BUILT_PRODUCTS_DIR}\"\ncd \"${SRCROOT}/..\"\ngo build -buildmode=c-archive -o \"${BUILT_PRODUCTS_DIR}/libcounter.a\" ./ios/bridge\n";
		};
/* End PBXShellScriptBuildPhase section */

/* Begin PBXSourcesBuildPhase section */
		100000000000000000000005 /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				100000000000000000000011 /* main.m in Sources */,
				100000000000000000000014 /* GoNativeRenderer.m in Sources */,
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
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_STYLE = Automatic;
				CURRENT_PROJECT_VERSION = 1;
				GENERATE_INFOPLIST_FILE = NO;
				HEADER_SEARCH_PATHS = (
					"$(BUILT_PRODUCTS_DIR)",
					"$(SRCROOT)",
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
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_STYLE = Automatic;
				CURRENT_PROJECT_VERSION = 1;
				GENERATE_INFOPLIST_FILE = NO;
				HEADER_SEARCH_PATHS = (
					"$(BUILT_PRODUCTS_DIR)",
					"$(SRCROOT)",
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
	jniPkg := sanitizeJniName(pkg)

	return map[string]string{
		"go.mod": fmt.Sprintf("module %s\n\ngo 1.24\n\nrequire github.com/go-native/go-native v0.0.0\n", name),
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
		".gitignore": "build/\n.gonative/\n.gradle/\n*.app\n*.apk\n*.idsig\nDerivedData/\n",
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
		fmt.Sprintf("ios/%s.xcodeproj/project.pbxproj", name):                          generatePbxproj(name, pkg),
		fmt.Sprintf("ios/%s.xcodeproj/xcshareddata/xcschemes/%s.xcscheme", name, name): generateXcscheme(name),
		"ios/GoNativeRenderer.h": `#import <UIKit/UIKit.h>

@interface GNRootViewController : UIViewController
@end

void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
int32_t GNMeasureNativeBatch(const uint8_t *bytes, int32_t length, uint8_t **results, int32_t *resultLength);
void GNFreeNativeBuffer(void *buffer);
`,
		"ios/GoNativeRenderer.m": `#import "GoNativeRenderer.h"
#import "counter.h"
#include <time.h>
#include <math.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };
typedef NS_ENUM(uint8_t, GNNode) { GNView=1, GNText, GNButton, GNRow, GNColumn, GNSafeArea, GNTextInput, GNSwitch, GNProgressIndicator, GNImage, GNScrollView };

@interface GNSafeAreaView : UIView
@property(nonatomic) uint8_t gnAlignment;
@end
@implementation GNSafeAreaView
@end

@interface GNAction : NSObject
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint64_t nodeID;
- (void)invoke;
- (void)change:(UITextField *)sender;
- (void)toggle:(UISwitch *)sender;
- (void)focus:(UITextField *)sender;
- (void)blur:(UITextField *)sender;
@end
@implementation GNAction
- (void)invoke { GoNativeDispatchEvent(self.handler); }
- (void)change:(UITextField *)sender { GoNativeDispatchValueEvent(self.handler, (char *)sender.text.UTF8String); }
- (void)toggle:(UISwitch *)sender { GoNativeDispatchBoolEvent(self.handler, sender.isOn ? 1 : 0); }
- (void)focus:(UITextField *)sender { (void)sender; GoNativeDispatchFocus(self.nodeID, 1); }
- (void)blur:(UITextField *)sender { (void)sender; GoNativeDispatchFocus(self.nodeID, 0); }
@end

@interface GNGestureAction : NSObject
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint8_t kind;
- (void)recognize:(UIGestureRecognizer *)recognizer;
@end
@implementation GNGestureAction
- (void)recognize:(UIGestureRecognizer *)recognizer {
    if ((self.kind == 2 && recognizer.state != UIGestureRecognizerStateBegan) ||
        (self.kind == 4 && recognizer.state != UIGestureRecognizerStateEnded)) return;
    if (self.kind != 2 && self.kind != 4 && recognizer.state != UIGestureRecognizerStateRecognized) return;
    CGPoint translation=CGPointZero, velocity=CGPointZero;
    if ([recognizer isKindOfClass:UIPanGestureRecognizer.class]) {
        UIPanGestureRecognizer *pan=(UIPanGestureRecognizer *)recognizer;
        translation=[pan translationInView:pan.view]; velocity=[pan velocityInView:pan.view];
    }
    GoNativeDispatchGestureEvent(self.handler,(float)translation.x,(float)translation.y,(float)velocity.x,(float)velocity.y);
}
@end

static NSMutableDictionary<NSNumber *,UIView *> *GNViews;
static NSMutableDictionary<NSNumber *,GNAction *> *GNActions;
static NSMutableDictionary<NSNumber *,NSArray<GNGestureAction *> *> *GNGestureActions;
static NSMutableDictionary<NSNumber *,NSValue *> *GNComputedFrames;
static NSMutableDictionary<NSNumber *,NSArray<NSLayoutConstraint *> *> *GNFrameConstraints;
static __weak GNRootViewController *GNRoot;

typedef struct { const uint8_t *p; const uint8_t *end; } GNReader;
static uint8_t u8(GNReader *r){return r->p<r->end?*r->p++:0;}
static uint16_t u16(GNReader*r){uint16_t v=0;if(r->p+2<=r->end){memcpy(&v,r->p,2);r->p+=2;}return v;}
static uint32_t u32(GNReader*r){uint32_t v=0;if(r->p+4<=r->end){memcpy(&v,r->p,4);r->p+=4;}return v;}
static uint64_t u64(GNReader*r){uint64_t v=0;if(r->p+8<=r->end){memcpy(&v,r->p,8);r->p+=8;}return v;}
static int32_t i32(GNReader*r){return (int32_t)u32(r);}
static float f32(GNReader*r){float v=0;if(r->p+4<=r->end){memcpy(&v,r->p,4);r->p+=4;}return v;}
static NSString *str(GNReader*r){uint32_t n=u32(r);if(r->p+n>r->end){r->p=r->end;return @"";}NSString*s=[[NSString alloc]initWithBytes:r->p length:n encoding:NSUTF8StringEncoding]?:@"";r->p+=n;return s;}
static uint64_t GNNowNanos(void){struct timespec t;clock_gettime(CLOCK_MONOTONIC_RAW,&t);return (uint64_t)t.tv_sec*1000000000ull+(uint64_t)t.tv_nsec;}

static UIViewAnimationOptions GNAnimationOptions(uint8_t curve) {
    switch(curve){case 1:return UIViewAnimationOptionCurveEaseIn;case 2:return UIViewAnimationOptionCurveEaseOut;case 3:return UIViewAnimationOptionCurveLinear;default:return UIViewAnimationOptionCurveEaseInOut;}
}

static void GNConfigureInteractions(uint64_t nodeID, UIView *view, NSData *payload, BOOL animate) {
    for (UIGestureRecognizer *recognizer in [view.gestureRecognizers copy]) [view removeGestureRecognizer:recognizer];
    NSMutableArray<GNGestureAction *> *targets=[NSMutableArray array];
    GNReader r={(const uint8_t *)payload.bytes,(const uint8_t *)payload.bytes+payload.length};
    uint32_t gestureCount=u32(&r);
    for(uint32_t i=0;i<gestureCount && r.p<r.end;i++){
        uint8_t kind=u8(&r),direction=u8(&r);uint64_t minimumPress=u64(&r);float minimumTravel=f32(&r);uint64_t handler=u64(&r);
        GNGestureAction *target=[GNGestureAction new];target.handler=handler;target.kind=kind;
        UIGestureRecognizer *recognizer=nil;
        if(kind==1) recognizer=[[UITapGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];
        else if(kind==2){UILongPressGestureRecognizer *longPress=[[UILongPressGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];longPress.minimumPressDuration=(NSTimeInterval)minimumPress/1e9;recognizer=longPress;}
        else if(kind==3){UISwipeGestureRecognizer *swipe=[[UISwipeGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];swipe.direction=direction==1?UISwipeGestureRecognizerDirectionUp:direction==2?UISwipeGestureRecognizerDirectionDown:direction==3?UISwipeGestureRecognizerDirectionLeft:direction==4?UISwipeGestureRecognizerDirectionRight:(UISwipeGestureRecognizerDirectionUp|UISwipeGestureRecognizerDirectionDown|UISwipeGestureRecognizerDirectionLeft|UISwipeGestureRecognizerDirectionRight);recognizer=swipe;}
        else if(kind==4){UIPanGestureRecognizer *pan=[[UIPanGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];pan.minimumNumberOfTouches=1;recognizer=pan;}
        (void)minimumTravel;
        if(recognizer){[view addGestureRecognizer:recognizer];[targets addObject:target];}
    }
    GNGestureActions[@(nodeID)]=targets;
    uint32_t animationCount=u32(&r);
    for(uint32_t i=0;i<animationCount && r.p<r.end;i++){
        uint8_t property=u8(&r);int64_t duration=(int64_t)u64(&r),delay=(int64_t)u64(&r);uint8_t curve=u8(&r);float damping=f32(&r),velocity=f32(&r);BOOL reduceMotionOK=u8(&r);
        float from=f32(&r),to=f32(&r),fromX=f32(&r),fromY=f32(&r),toX=f32(&r),toY=f32(&r);
        void (^initial)(void)=^{if(property==1)view.alpha=from;else if(property==2)view.transform=CGAffineTransformMakeScale(from,from);else if(property==3)view.transform=CGAffineTransformMakeTranslation(fromX,fromY);};
        void (^final)(void)=^{if(property==1)view.alpha=to;else if(property==2)view.transform=CGAffineTransformMakeScale(to,to);else if(property==3)view.transform=CGAffineTransformMakeTranslation(toX,toY);else if(property==4)[view.superview layoutIfNeeded];};
        if(!animate) continue;
        initial();
        BOOL reduce=UIAccessibilityIsReduceMotionEnabled() && !reduceMotionOK;
        if(reduce || duration<=0){final();continue;}
        NSTimeInterval seconds=(NSTimeInterval)duration/1e9, wait=(NSTimeInterval)delay/1e9;
        if(curve==4)[UIView animateWithDuration:seconds delay:wait usingSpringWithDamping:MAX(0.01,MIN(1,damping)) initialSpringVelocity:velocity options:UIViewAnimationOptionBeginFromCurrentState animations:final completion:nil];
        else [UIView animateWithDuration:seconds delay:wait options:(GNAnimationOptions(curve)|UIViewAnimationOptionBeginFromCurrentState) animations:final completion:nil];
    }
}

static void GNStyle(uint64_t nodeID, UIView *view, GNNode kind, NSString *text, float width, float height, float padding, float gap, uint8_t alignment, float fontSize, BOOL bold, uint64_t handler, uint64_t changeHandler, uint64_t toggleHandler, BOOL checked, float progress, NSString *accessibility, NSString *hint, uint8_t role, BOOL focused, BOOL scalesText, NSString *imageSource, uint8_t imageMode, BOOL horizontal, NSData *interactions, BOOL animate) {
    if ([view isKindOfClass:UILabel.class]) { UILabel*l=(UILabel*)view;l.text=text;UIFont*base=bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:17]:[UIFont systemFontOfSize:fontSize>0?fontSize:17];l.font=scalesText?[[UIFontMetrics defaultMetrics] scaledFontForFont:base]:base;l.adjustsFontForContentSizeCategory=scalesText; }
    if ([view isKindOfClass:UIButton.class]) {
        UIButton*b=(UIButton*)view;
        UIFont*btnFont=bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:16]:[UIFont systemFontOfSize:fontSize>0?fontSize:16];
        if (@available(iOS 15.0, *)) {
            UIButtonConfiguration *cfg = b.configuration ?: [UIButtonConfiguration filledButtonConfiguration];
            cfg.attributedTitle = [[NSAttributedString alloc] initWithString:text?:@"" attributes:@{NSFontAttributeName:btnFont}];
            cfg.baseBackgroundColor = [UIColor systemBlueColor];
            cfg.baseForegroundColor = [UIColor whiteColor];
            cfg.cornerStyle = UIButtonConfigurationCornerStyleMedium;
            cfg.contentInsets = NSDirectionalEdgeInsetsMake(12,20,12,20);
            b.configuration = cfg;
        } else {
            [b setTitle:text forState:UIControlStateNormal];
            b.titleLabel.font = btnFont;
            b.backgroundColor = [UIColor systemBlueColor];
            [b setTitleColor:UIColor.whiteColor forState:UIControlStateNormal];
            b.layer.cornerRadius = 8.0;
            b.clipsToBounds = YES;
        }
        b.titleLabel.numberOfLines = 1;
        b.titleLabel.lineBreakMode = NSLineBreakByTruncatingTail;
        NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&handler){a=[GNAction new];GNActions[actionKey]=a;[b addTarget:a action:@selector(invoke) forControlEvents:UIControlEventTouchUpInside];}a.handler=handler;
    }
    if ([view isKindOfClass:UITextField.class]) {
        UITextField*f=(UITextField*)view;
        if(![f.text isEqualToString:text]) f.text=text;
        if(hint.length) f.placeholder = hint;
        f.font = bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:16]:[UIFont systemFontOfSize:fontSize>0?fontSize:16];
        f.borderStyle = UITextBorderStyleRoundedRect;
        f.autocapitalizationType = UITextAutocapitalizationTypeNone;
        f.autocorrectionType = UITextAutocorrectionTypeNo;
        NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a){a=[GNAction new];a.nodeID=nodeID;GNActions[actionKey]=a;[f addTarget:a action:@selector(change:) forControlEvents:UIControlEventEditingChanged];[f addTarget:a action:@selector(focus:) forControlEvents:UIControlEventEditingDidBegin];[f addTarget:a action:@selector(blur:) forControlEvents:UIControlEventEditingDidEnd];}a.handler=changeHandler;
    }
    if ([view isKindOfClass:UISwitch.class]) { UISwitch*s=(UISwitch*)view;[s setOn:checked animated:NO];NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&toggleHandler){a=[GNAction new];GNActions[actionKey]=a;[s addTarget:a action:@selector(toggle:) forControlEvents:UIControlEventValueChanged];}a.handler=toggleHandler; }
    if ([view isKindOfClass:UIProgressView.class]) { [(UIProgressView*)view setProgress:progress animated:NO]; }
    if ([view isKindOfClass:UIImageView.class]) {
        UIImageView*i=(UIImageView*)view;
        UIImage *img = [UIImage imageNamed:imageSource];
        if(!img && [imageSource isEqualToString:@"app_logo"]) {
            img = [UIImage systemImageNamed:@"lock.shield.fill"];
            i.tintColor = UIColor.systemBlueColor;
        } else if(!img && [imageSource isEqualToString:@"avatar"]) {
            img = [UIImage systemImageNamed:@"person.crop.circle.fill"];
            i.tintColor = UIColor.systemBlueColor;
        } else if(!img && imageSource.length) {
            img = [UIImage systemImageNamed:imageSource];
        }
        i.image=img;
        i.contentMode=imageMode==1?UIViewContentModeScaleAspectFill:imageMode==2?UIViewContentModeCenter:UIViewContentModeScaleAspectFit;
        i.clipsToBounds=YES;
    }
    if ([view isKindOfClass:UIStackView.class]) { UIStackView*s=(UIStackView*)view;s.spacing=gap;s.layoutMarginsRelativeArrangement=YES;s.directionalLayoutMargins=NSDirectionalEdgeInsetsMake(padding,padding,padding,padding);s.alignment=alignment==1?UIStackViewAlignmentCenter:alignment==2?UIStackViewAlignmentTrailing:UIStackViewAlignmentLeading; }
    if ([view isKindOfClass:GNSafeAreaView.class]) { ((GNSafeAreaView *)view).gnAlignment=alignment; }
    if(width>0)[view.widthAnchor constraintEqualToConstant:width].active=YES;if(height>0)[view.heightAnchor constraintEqualToConstant:height].active=YES;
    view.isAccessibilityElement=(kind==GNText||kind==GNButton||kind==GNTextInput||kind==GNSwitch||kind==GNProgressIndicator||role!=0);view.accessibilityLabel=accessibility.length?accessibility:text;view.accessibilityHint=hint;UIAccessibilityTraits traits=UIAccessibilityTraitNone;if(role==2||kind==GNButton)traits|=UIAccessibilityTraitButton;if(role==3)traits|=UIAccessibilityTraitHeader;if(role==4)traits|=UIAccessibilityTraitImage;view.accessibilityTraits=traits;if(focused){if(view.canBecomeFirstResponder&&![view isFirstResponder])[view becomeFirstResponder];UIAccessibilityPostNotification(UIAccessibilityScreenChangedNotification,view);}else if([view isFirstResponder]){[view resignFirstResponder];}
    GNConfigureInteractions(nodeID,view,interactions,animate);
}

static UIView *GNMake(GNNode kind){UIView*v;if(kind==GNText){UILabel*l=[UILabel new];l.numberOfLines=0;l.textColor=UIColor.labelColor;v=l;}else if(kind==GNButton){UIButton*b=[UIButton buttonWithType:UIButtonTypeSystem];v=b;}else if(kind==GNTextInput){UITextField*f=[UITextField new];f.borderStyle=UITextBorderStyleRoundedRect;v=f;}else if(kind==GNSwitch){v=[UISwitch new];}else if(kind==GNProgressIndicator){v=[[UIProgressView alloc]initWithProgressViewStyle:UIProgressViewStyleDefault];}else if(kind==GNImage){v=[UIImageView new];}else if(kind==GNScrollView){v=[UIScrollView new];}else if(kind==GNRow||kind==GNColumn){UIStackView*s=[UIStackView new];s.axis=kind==GNRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;v=s;}else if(kind==GNSafeArea){v=[GNSafeAreaView new];v.backgroundColor=UIColor.systemBackgroundColor;}else{v=[UIView new];v.backgroundColor=UIColor.systemBackgroundColor;}v.translatesAutoresizingMaskIntoConstraints=NO;return v;}

static UIColor *GNColor(const uint8_t *p){return [UIColor colorWithRed:p[0]/255.0 green:p[1]/255.0 blue:p[2]/255.0 alpha:p[3]/255.0];}
static NSUInteger GNStyleSize(const uint8_t *style,const uint8_t *end){if(style+185>end)return 0;uint32_t fontLength=0;memcpy(&fontLength,style+181,4);NSUInteger size=220+(NSUInteger)fontLength;return style+size<=end?size:0;}
static BOOL GNHasTypedValues(const uint8_t *p,NSUInteger length,const uint8_t *end){if(p+length>end)return NO;for(NSUInteger i=0;i<length;i++)if(p[i]!=0)return YES;return NO;}
static void GNApplyTypedStyle(UIView *view, NSData *payload){
    if(!view||payload.length<187)return;const uint8_t*record=payload.bytes,*end=record+payload.length;uint16_t version=0;memcpy(&version,record,2);if(version!=1)return;
    const uint8_t*portable=record+2;NSUInteger portableSize=GNStyleSize(portable,end);if(!portableSize)return;const uint8_t*ios=portable+portableSize;NSUInteger iosSize=GNStyleSize(ios,end);if(!iosSize)return;
    const uint8_t*appearance=GNHasTypedValues(ios+112,69,end)?ios+112:portable+112;
    uint32_t iosFontLength=0;memcpy(&iosFontLength,ios+181,4);const uint8_t*text=GNHasTypedValues(ios+181,22+(NSUInteger)iosFontLength,end)?ios:portable;
    const uint8_t*interaction=GNHasTypedValues(ios+203+(NSUInteger)iosFontLength,17,end)?ios:portable;
    float borderWidth=0,cornerRadius=0,opacity=0;memcpy(&borderWidth,appearance+8,4);memcpy(&cornerRadius,appearance+16,4);memcpy(&opacity,appearance+44,4);
    if(appearance[3]>0)view.backgroundColor=GNColor(appearance);
    if([view isKindOfClass:UILabel.class]&&appearance[7]>0)((UILabel*)view).textColor=GNColor(appearance+4);
    if([view isKindOfClass:UIButton.class]&&appearance[7]>0)[((UIButton*)view) setTitleColor:GNColor(appearance+4) forState:UIControlStateNormal];
    if(borderWidth>0){view.layer.borderWidth=borderWidth;view.layer.borderColor=GNColor(appearance+12).CGColor;}
    if(cornerRadius>0){view.layer.cornerRadius=cornerRadius;view.clipsToBounds=YES;}
    if(opacity>0)view.alpha=MIN(1,opacity);view.hidden=appearance[68]!=0;
    float tx=0,ty=0,sx=0,sy=0,rotation=0;memcpy(&tx,appearance+48,4);memcpy(&ty,appearance+52,4);memcpy(&sx,appearance+56,4);memcpy(&sy,appearance+60,4);memcpy(&rotation,appearance+64,4);view.transform=CGAffineTransformRotate(CGAffineTransformScale(CGAffineTransformMakeTranslation(tx,ty),sx==0?1:sx,sy==0?1:sy),rotation*(CGFloat)M_PI/180.0);
    float shadowX=0,shadowY=0,shadowBlur=0,shadowOpacity=0;memcpy(&shadowX,appearance+24,4);memcpy(&shadowY,appearance+28,4);memcpy(&shadowBlur,appearance+32,4);memcpy(&shadowOpacity,appearance+40,4);if(shadowOpacity>0){view.layer.shadowColor=GNColor(appearance+20).CGColor;view.layer.shadowOffset=CGSizeMake(shadowX,shadowY);view.layer.shadowRadius=shadowBlur;view.layer.shadowOpacity=shadowOpacity;}
    uint32_t fontLength=0;memcpy(&fontLength,text+181,4);NSUInteger disabledOffset=203;uint32_t interactionFontLength=0;memcpy(&interactionFontLength,interaction+181,4);disabledOffset+=(NSUInteger)interactionFontLength;if(interaction+disabledOffset>=end)return;NSUInteger fontOffset=185+(NSUInteger)fontLength;float fontSize=0,lineHeight=0,letterSpacing=0;uint16_t fontWeight=0;memcpy(&fontSize,text+fontOffset,4);memcpy(&fontWeight,text+fontOffset+4,2);memcpy(&lineHeight,text+fontOffset+6,4);memcpy(&letterSpacing,text+fontOffset+10,4);NSString*family=[[NSString alloc]initWithBytes:text+185 length:fontLength encoding:NSUTF8StringEncoding]?:@"";UIFont*font=family.length?[UIFont fontWithName:family size:fontSize]:nil;if(!font&&fontSize>0)font=[UIFont systemFontOfSize:fontSize weight:fontWeight>=600?UIFontWeightBold:UIFontWeightRegular];if(font){if([view isKindOfClass:UILabel.class])((UILabel*)view).font=font;else if([view isKindOfClass:UIButton.class])((UIButton*)view).titleLabel.font=font;else if([view isKindOfClass:UITextField.class])((UITextField*)view).font=font;}if([view isKindOfClass:UILabel.class]&&(lineHeight>0||letterSpacing!=0)){UILabel*l=(UILabel*)view;NSMutableParagraphStyle*paragraph=[NSMutableParagraphStyle new];if(lineHeight>0){paragraph.minimumLineHeight=lineHeight;paragraph.maximumLineHeight=lineHeight;}l.attributedText=[[NSAttributedString alloc]initWithString:l.text?:@"" attributes:@{NSKernAttributeName:@(letterSpacing),NSParagraphStyleAttributeName:paragraph}];}view.userInteractionEnabled=interaction[disabledOffset]==0;
}

static void GNConstrainSafeAreaChild(GNSafeAreaView *parent, UIView *view) {
    UILayoutGuide *guide=parent.safeAreaLayoutGuide;
    NSMutableArray<NSLayoutConstraint *> *constraints=[NSMutableArray arrayWithArray:@[[view.leadingAnchor constraintGreaterThanOrEqualToAnchor:guide.leadingAnchor],[view.trailingAnchor constraintLessThanOrEqualToAnchor:guide.trailingAnchor],[view.topAnchor constraintGreaterThanOrEqualToAnchor:guide.topAnchor],[view.bottomAnchor constraintLessThanOrEqualToAnchor:guide.bottomAnchor]]];
    if(parent.gnAlignment==1){[constraints addObject:[view.centerXAnchor constraintEqualToAnchor:guide.centerXAnchor]];[constraints addObject:[view.centerYAnchor constraintEqualToAnchor:guide.centerYAnchor]];}
    else if(parent.gnAlignment==2){[constraints addObject:[view.trailingAnchor constraintEqualToAnchor:guide.trailingAnchor]];[constraints addObject:[view.bottomAnchor constraintEqualToAnchor:guide.bottomAnchor]];}
    else{[constraints addObject:[view.leadingAnchor constraintEqualToAnchor:guide.leadingAnchor]];[constraints addObject:[view.topAnchor constraintEqualToAnchor:guide.topAnchor]];}
    [NSLayoutConstraint activateConstraints:constraints];
}

static void GNApplyComputedFrame(uint64_t nodeID,UIView *view) {
    if(!view||!view.superview||view.superview==GNRoot.view)return;
    NSValue *value=GNComputedFrames[@(nodeID)];if(!value)return;CGRect frame=value.CGRectValue;
    NSArray<NSLayoutConstraint *> *old=GNFrameConstraints[@(nodeID)];if(old.count)[NSLayoutConstraint deactivateConstraints:old];
    NSMutableArray<NSLayoutConstraint *> *legacySize=[NSMutableArray array];for(NSLayoutConstraint *constraint in view.constraints)if(constraint.firstItem==view&&(constraint.firstAttribute==NSLayoutAttributeWidth||constraint.firstAttribute==NSLayoutAttributeHeight))[legacySize addObject:constraint];if(legacySize.count)[NSLayoutConstraint deactivateConstraints:legacySize];
    NSMutableArray<NSLayoutConstraint *> *constraints=[NSMutableArray array];
    if(![view.superview isKindOfClass:UIStackView.class]){[constraints addObject:[view.leadingAnchor constraintEqualToAnchor:view.superview.leadingAnchor constant:frame.origin.x]];[constraints addObject:[view.topAnchor constraintEqualToAnchor:view.superview.topAnchor constant:frame.origin.y]];}
    if(frame.size.width>=0)[constraints addObject:[view.widthAnchor constraintEqualToConstant:frame.size.width]];
    if(frame.size.height>=0)[constraints addObject:[view.heightAnchor constraintEqualToConstant:frame.size.height]];
    [NSLayoutConstraint activateConstraints:constraints];GNFrameConstraints[@(nodeID)]=constraints;
}

static void GNApply(NSData *data){uint64_t started=GNNowNanos();GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};if(u16(&r)!=9)return;uint32_t count=u32(&r);uint64_t sequence=u64(&r);for(uint32_t op=0;op<count;op++){GNMutation mutation=(GNMutation)u8(&r);GNNode kind=(GNNode)u8(&r);uint64_t nodeID=u64(&r),parentID=u64(&r);int32_t index=i32(&r),from=i32(&r);float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);uint8_t alignment=u8(&r);BOOL bold=u8(&r);float fontSize=f32(&r);uint64_t handler=u64(&r),changeHandler=u64(&r),toggleHandler=u64(&r);BOOL checked=u8(&r);float progress=f32(&r);NSString*text=str(&r);NSString*accessibility=str(&r);NSString*hint=str(&r);uint8_t role=u8(&r);BOOL focused=u8(&r);BOOL scalesText=u8(&r);NSString*imageSource=str(&r);uint8_t imageMode=u8(&r);BOOL horizontal=u8(&r);uint32_t interactionLength=u32(&r);NSData *interactions;if(r.p+interactionLength<=r.end){interactions=[NSData dataWithBytes:r.p length:interactionLength];r.p+=interactionLength;}else{r.p=r.end;interactions=[NSData data];}uint32_t styleLength=u32(&r);if(styleLength>1048576||r.p+styleLength>r.end)return;NSData *typedStyle=[NSData dataWithBytes:r.p length:styleLength];r.p+=styleLength;if(r.p+17>r.end)return;BOOL hasFrame=u8(&r);float frameX=f32(&r),frameY=f32(&r),frameWidth=f32(&r),frameHeight=f32(&r);NSNumber*key=@(nodeID);if(hasFrame)GNComputedFrames[key]=[NSValue valueWithCGRect:CGRectMake(frameX,frameY,MAX(0,frameWidth),MAX(0,frameHeight))];else{[GNComputedFrames removeObjectForKey:key];NSArray<NSLayoutConstraint*>*old=GNFrameConstraints[key];if(old.count)[NSLayoutConstraint deactivateConstraints:old];[GNFrameConstraints removeObjectForKey:key];}UIView*view=GNViews[key];
    if(mutation==GNCreate){view=GNMake(kind);GNViews[key]=view;GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,NO);GNApplyTypedStyle(view,typedStyle);if(!GNRoot.view.subviews.count){view.backgroundColor=UIColor.systemBackgroundColor;[GNRoot.view addSubview:view];[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.bottomAnchor]]];}}
    else if(mutation==GNUpdate){GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,YES);GNApplyTypedStyle(view,typedStyle);GNApplyComputedFrame(nodeID,view);}
    else if(mutation==GNInsert){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}else{[parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];GNSafeAreaView*safe=[parent isKindOfClass:GNSafeAreaView.class]?(GNSafeAreaView*)parent:nil;if([parent isKindOfClass:UIScrollView.class]&&!hasFrame){UIScrollView*s=(UIScrollView*)parent;[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:s.contentLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:s.contentLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:s.contentLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:s.contentLayoutGuide.bottomAnchor],[view.widthAnchor constraintEqualToAnchor:s.frameLayoutGuide.widthAnchor]]];}else if(safe&&!hasFrame){GNConstrainSafeAreaChild(safe,view);}GNApplyComputedFrame(nodeID,view);}}
    else if(mutation==GNRemove){[view removeFromSuperview];}
    else if(mutation==GNMove){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s removeArrangedSubview:view];[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}}
    else if(mutation==GNDelete){NSArray<NSLayoutConstraint *>*frameConstraints=GNFrameConstraints[key];if(frameConstraints.count)[NSLayoutConstraint deactivateConstraints:frameConstraints];[view removeFromSuperview];[GNViews removeObjectForKey:key];[GNActions removeObjectForKey:key];[GNGestureActions removeObjectForKey:key];[GNComputedFrames removeObjectForKey:key];[GNFrameConstraints removeObjectForKey:key];}
    (void)from;
}GoNativeReportBatchApplied(sequence,GNNowNanos()-started);}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});}

static void GNAppendU16(NSMutableData *data,uint16_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendU32(NSMutableData *data,uint32_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendU64(NSMutableData *data,uint64_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendF32(NSMutableData *data,float value){[data appendBytes:&value length:sizeof(value)];}
static BOOL GNReadStringChecked(GNReader *reader,NSString **value){if(reader->p+4>reader->end)return NO;uint32_t length=u32(reader);if(length>1048576||reader->p+length>reader->end)return NO;*value=[[NSString alloc]initWithBytes:reader->p length:length encoding:NSUTF8StringEncoding]?:@"";reader->p+=length;return YES;}
static UIFont *GNMeasurementFont(NSData *typedStyle,CGFloat fallback){
    const uint8_t *record=typedStyle.bytes,*end=record+typedStyle.length;if(typedStyle.length<2)return [UIFont systemFontOfSize:fallback];uint16_t version=0;memcpy(&version,record,2);if(version!=1)return [UIFont systemFontOfSize:fallback];
    const uint8_t *style=record+2;NSUInteger styleSize=GNStyleSize(style,end);if(!styleSize)return [UIFont systemFontOfSize:fallback];uint32_t familyLength=0;memcpy(&familyLength,style+181,4);NSUInteger fontOffset=185+(NSUInteger)familyLength;if(style+fontOffset+6>end)return [UIFont systemFontOfSize:fallback];
    float fontSize=0;uint16_t weight=0;memcpy(&fontSize,style+fontOffset,4);memcpy(&weight,style+fontOffset+4,2);NSString *family=[[NSString alloc]initWithBytes:style+185 length:familyLength encoding:NSUTF8StringEncoding]?:@"";CGFloat size=fontSize>0?fontSize:fallback;UIFont *font=family.length?[UIFont fontWithName:family size:size]:nil;return font?:[UIFont systemFontOfSize:size weight:weight>=600?UIFontWeightBold:UIFontWeightRegular];
}
static CGSize GNMeasureControl(GNNode kind,NSString *text,NSString *imageSource,NSData *typedStyle,CGSize constraint){
    UIView *view=GNMake(kind);UIFont *font=GNMeasurementFont(typedStyle,kind==GNButton||kind==GNTextInput?16:17);
    if([view isKindOfClass:UILabel.class]){UILabel *label=(UILabel *)view;label.text=text;label.font=font;label.numberOfLines=0;}
    else if([view isKindOfClass:UIButton.class]){CGRect title=[(text?:@"") boundingRectWithSize:CGSizeMake(CGFLOAT_MAX,CGFLOAT_MAX) options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading attributes:@{NSFontAttributeName:font} context:nil];return CGSizeMake(MAX(44,ceil(title.size.width)+40),MAX(44,ceil(title.size.height)+24));}
    else if([view isKindOfClass:UITextField.class]){CGRect content=[(text.length?text:@"M") boundingRectWithSize:CGSizeMake(CGFLOAT_MAX,CGFLOAT_MAX) options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading attributes:@{NSFontAttributeName:font} context:nil];return CGSizeMake(MAX(240,ceil(content.size.width)+32),MAX(44,ceil(content.size.height)+20));}
    else if([view isKindOfClass:UIImageView.class]){UIImage *image=[UIImage imageNamed:imageSource]?:[UIImage systemImageNamed:imageSource];((UIImageView *)view).image=image;}
    CGSize size=[view sizeThatFits:constraint];if(size.width<=0||size.height<=0)size=view.intrinsicContentSize;if(size.width==UIViewNoIntrinsicMetric)size.width=0;if(size.height==UIViewNoIntrinsicMetric)size.height=0;return size;
}
int32_t GNMeasureNativeBatch(const uint8_t *bytes,int32_t length,uint8_t **results,int32_t *resultLength){
    if(!results||!resultLength){return 1;}*results=NULL;*resultLength=0;if(!bytes||length<6||length>16777216)return 2;NSData *input=[NSData dataWithBytes:bytes length:(NSUInteger)length];__block NSMutableData *output=nil;__block int32_t status=0;
    void (^measure)(void)=^{GNReader reader={(const uint8_t *)input.bytes,(const uint8_t *)input.bytes+input.length};uint16_t version=u16(&reader);uint32_t count=u32(&reader);if(version!=1||count>100000){status=3;return;}output=[NSMutableData data];GNAppendU16(output,1);GNAppendU32(output,count);
        for(uint32_t index=0;index<count;index++){if(reader.p+25>reader.end){status=4;return;}uint64_t requestID=u64(&reader);GNNode kind=(GNNode)u8(&reader);float minWidth=f32(&reader),maxWidth=f32(&reader),minHeight=f32(&reader),maxHeight=f32(&reader);NSString *text=@"",*image=@"";if(!GNReadStringChecked(&reader,&text)||!GNReadStringChecked(&reader,&image)||reader.p+4>reader.end){status=4;return;}uint32_t styleLength=u32(&reader);if(styleLength>1048576||reader.p+styleLength>reader.end){status=4;return;}NSData *style=[NSData dataWithBytes:reader.p length:styleLength];reader.p+=styleLength;CGFloat width=isfinite(maxWidth)&&maxWidth>0?maxWidth:CGFLOAT_MAX,height=isfinite(maxHeight)&&maxHeight>0?maxHeight:CGFLOAT_MAX;CGSize size=GNMeasureControl(kind,text,image,style,CGSizeMake(width,height));size.width=MAX(minWidth,MIN(size.width,width));size.height=MAX(minHeight,MIN(size.height,height));GNAppendU64(output,requestID);GNAppendF32(output,(float)size.width);GNAppendF32(output,(float)size.height);GNAppendU32(output,0);}
        if(reader.p!=reader.end||output.length>16777216){status=5;output=nil;}
    };if(NSThread.isMainThread)measure();else dispatch_sync(dispatch_get_main_queue(),measure);if(status!=0||!output)return status?:6;void *buffer=malloc(output.length);if(!buffer)return 7;memcpy(buffer,output.bytes,output.length);*results=buffer;*resultLength=(int32_t)output.length;return 0;
}
void GNFreeNativeBuffer(void *buffer){free(buffer);}

@implementation GNRootViewController
- (void)viewDidLoad {[super viewDidLoad];self.view.backgroundColor=UIColor.systemBackgroundColor;GNRoot=self;GNViews=[NSMutableDictionary dictionary];GNActions=[NSMutableDictionary dictionary];GNGestureActions=[NSMutableDictionary dictionary];GNComputedFrames=[NSMutableDictionary dictionary];GNFrameConstraints=[NSMutableDictionary dictionary];GoNativeStart();GoNativeSetLifecycle(1);[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnDidBecomeActive) name:UIApplicationDidBecomeActiveNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnWillResignActive) name:UIApplicationWillResignActiveNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnDidEnterBackground) name:UIApplicationDidEnterBackgroundNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnWillEnterForeground) name:UIApplicationWillEnterForegroundNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnMemoryWarning) name:UIApplicationDidReceiveMemoryWarningNotification object:nil];}
- (void)gnDidBecomeActive { GoNativeSetLifecycle(2); }
- (void)gnWillResignActive { GoNativeSetLifecycle(3); }
- (void)gnDidEnterBackground { GoNativeSetLifecycle(4); }
- (void)gnWillEnterForeground { GoNativeSetLifecycle(1); }
- (void)gnMemoryWarning { GoNativeSetLifecycle(5); }
- (void)viewDidLayoutSubviews {[super viewDidLayoutSubviews];CGRect viewport=self.view.safeAreaLayoutGuide.layoutFrame;GoNativeSetViewport((float)viewport.size.width,(float)viewport.size.height,(float)UIScreen.mainScreen.scale);}
- (void)dealloc { [[NSNotificationCenter defaultCenter]removeObserver:self];GoNativeSetLifecycle(6);GoNativeStop(); }
@end
`,
		"ios/main.m": `#import <UIKit/UIKit.h>
#import "GoNativeRenderer.h"

@interface GNAppDelegate : UIResponder <UIApplicationDelegate>
@property (strong, nonatomic) UIWindow *window;
@end

@implementation GNAppDelegate
- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    self.window.rootViewController = [GNRootViewController new];
    [self.window makeKeyAndVisible];
    return YES;
}
@end

int main(int argc, char * argv[]) {
    @autoreleasepool {
        return UIApplicationMain(argc, argv, nil, NSStringFromClass([GNAppDelegate class]));
    }
}
`,
		"ios/Info.plist": fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
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
`, name, pkg, name),

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
	"time"
	"unsafe"
)

var benchmarkOutput string

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
	appRuntime = gnruntime.New(app.App, iosRenderer{})
	appRuntime.SetLayoutProvider(&layout.Pipeline{Measurer: iosNativeMeasurer{}, Cache: layout.NewMeasurementCache()})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
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
    implementation "androidx.recyclerview:recyclerview:1.4.0"
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
`,

		"android/app/src/main/AndroidManifest.xml": fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android" package="dev.gonative.%s" android:versionCode="1" android:versionName="0.1">
    <application android:theme="@style/AppTheme" android:label="%s" android:allowBackup="false" android:supportsRtl="true">
        <activity android:name=".MainActivity" android:screenOrientation="portrait" android:exported="true">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>
`, pkg, name),

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

		fmt.Sprintf("android/app/src/main/java/dev/gonative/%s/GapDrawable.java", pkg): fmt.Sprintf(`package dev.gonative.%s;

import android.graphics.Canvas;
import android.graphics.ColorFilter;
import android.graphics.PixelFormat;
import android.graphics.drawable.Drawable;
import android.widget.LinearLayout;

@SuppressWarnings("deprecation")
final class GapDrawable extends Drawable {
    private final int size;
    private final int orientation;

    GapDrawable(int size, int orientation) {
        this.size = size;
        this.orientation = orientation;
    }

    @Override public void draw(Canvas canvas) {}
    @Override public void setAlpha(int alpha) {}
    @Override public void setColorFilter(ColorFilter filter) {}
    @Override public int getOpacity() { return PixelFormat.TRANSPARENT; }
    @Override public int getIntrinsicWidth() { return orientation == LinearLayout.HORIZONTAL ? size : 0; }
    @Override public int getIntrinsicHeight() { return orientation == LinearLayout.VERTICAL ? size : 0; }
}
`, pkg),

		fmt.Sprintf("android/app/src/main/java/dev/gonative/%s/MainActivity.java", pkg): fmt.Sprintf(`package dev.gonative.%s;

import android.app.Activity;
import android.graphics.Typeface;
import android.os.Bundle;
import android.os.Looper;
import android.util.LongSparseArray;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.view.accessibility.AccessibilityEvent;
import android.view.accessibility.AccessibilityNodeInfo;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.EditText;
import android.widget.Switch;
import android.widget.ProgressBar;
import android.widget.CompoundButton;
import android.widget.ImageView;
import android.widget.ScrollView;
import android.widget.HorizontalScrollView;
import android.text.Editable;
import android.text.TextWatcher;
import android.animation.Animator;
import android.animation.AnimatorSet;
import android.animation.ObjectAnimator;
import android.animation.ValueAnimator;
import android.view.MotionEvent;
import android.view.VelocityTracker;
import android.view.animation.AccelerateDecelerateInterpolator;
import android.view.animation.AccelerateInterpolator;
import android.view.animation.DecelerateInterpolator;
import android.view.animation.LinearInterpolator;
import android.view.animation.OvershootInterpolator;

import java.util.ArrayList;
import java.util.List;
import java.util.WeakHashMap;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicReference;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;

public final class MainActivity extends Activity {
    private static final int CREATE = 1;
    private static final int DELETE = 2;
    private static final int UPDATE = 3;
    private static final int INSERT = 4;
    private static final int REMOVE = 5;
    private static final int MOVE = 6;
    private static final int TEXT = 2;
    private static final int BUTTON = 3;
    private static final int ROW = 4;
    private static final int COLUMN = 5;
    private static final int SAFE_AREA = 6;
    private static final int TEXT_INPUT = 7;
    private static final int SWITCH = 8;
    private static final int PROGRESS_INDICATOR = 9;
    private static final int IMAGE = 10;
    private static final int SCROLL_VIEW = 11;

    static { System.loadLibrary("gonative"); }

    private final LongSparseArray<View> views = new LongSparseArray<>();
    private final LongSparseArray<GestureBinding> gestureBindings = new LongSparseArray<>();
    private final WeakHashMap<EditText, TextWatcher> textWatchers = new WeakHashMap<>();
    private long rootNodeID;
    private float lastViewportWidth, lastViewportHeight, lastViewportScale;
    private final View.OnLayoutChangeListener viewportListener = new View.OnLayoutChangeListener() {
        @Override public void onLayoutChange(View view, int left, int top, int right, int bottom,
                                             int oldLeft, int oldTop, int oldRight, int oldBottom) {
            reportViewport(view);
        }
    };

    private native void nativeStart();
    private native void nativeDispatchEvent(long handler);
    private native void nativeDispatchValueEvent(long handler, String value);
    private native void nativeDispatchBoolEvent(long handler, boolean value);
    private native void nativeDispatchGestureEvent(long handler, float translationX, float translationY, float velocityX, float velocityY);
    private native void nativeStop();
    private native void nativeSetLifecycle(int state);
    private native void nativeDispatchFocus(long nodeID, boolean focused);
    private native void nativeUpdateViewport(float width, float height, float scale);
    private native void nativeReportBatchApplied(long sequence, long nativeNanos);

    @Override protected void onCreate(Bundle state) {
        super.onCreate(state);
        getWindow().getDecorView().setBackgroundColor(android.graphics.Color.WHITE);
        nativeStart();
        nativeSetLifecycle(0);
        View content = findViewById(android.R.id.content);
        content.addOnLayoutChangeListener(viewportListener);
        content.post(new Runnable() { @Override public void run() { reportViewport(findViewById(android.R.id.content)); } });
    }

    @Override protected void onStart() { super.onStart(); nativeSetLifecycle(1); }
    @Override protected void onResume() { super.onResume(); nativeSetLifecycle(2); }
    @Override protected void onPause() { nativeSetLifecycle(3); super.onPause(); }
    @Override protected void onStop() { nativeSetLifecycle(4); super.onStop(); }
    @Override public void onLowMemory() { nativeSetLifecycle(5); super.onLowMemory(); }

    @Override protected void onDestroy() {
        nativeSetLifecycle(6);
        View content = findViewById(android.R.id.content);
        if (content != null) content.removeOnLayoutChangeListener(viewportListener);
        nativeStop();
        for (int i = 0; i < gestureBindings.size(); i++) gestureBindings.valueAt(i).dispose();
        gestureBindings.clear();
        textWatchers.clear();
        for (int i = 0; i < views.size(); i++) views.valueAt(i).animate().cancel();
        views.clear();
        super.onDestroy();
    }

    private void reportViewport(View content) {
        if (content == null || content.getWidth() <= 0 || content.getHeight() <= 0) return;
        float scale = getResources().getDisplayMetrics().density;
        float width = content.getWidth() / scale, height = content.getHeight() / scale;
        if (width == lastViewportWidth && height == lastViewportHeight && scale == lastViewportScale) return;
        lastViewportWidth = width; lastViewportHeight = height; lastViewportScale = scale;
        nativeUpdateViewport(width, height, scale);
    }

    @SuppressWarnings("unused")
    public byte[] measureNativeBatch(final byte[] payload) {
        if (payload == null || payload.length == 0 || payload.length > 16777216) return null;
        if (Looper.myLooper() == Looper.getMainLooper()) return measureNativeBatchOnUiThread(payload.clone());
        final AtomicReference<byte[]> result = new AtomicReference<>();
        final CountDownLatch ready = new CountDownLatch(1);
        runOnUiThread(new Runnable() {
            @Override public void run() {
                try { result.set(measureNativeBatchOnUiThread(payload.clone())); }
                finally { ready.countDown(); }
            }
        });
        try {
            return ready.await(10, TimeUnit.SECONDS) ? result.get() : null;
        } catch (InterruptedException interrupted) {
            Thread.currentThread().interrupt();
            return null;
        }
    }

    private byte[] measureNativeBatchOnUiThread(byte[] payload) {
        try {
            ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
            if (in.remaining() < 6 || Short.toUnsignedInt(in.getShort()) != 1) return null;
            int count = in.getInt();
            if (count < 0 || count > 100000) return null;
            ArrayList<NativeMeasurement> measured = new ArrayList<>(count);
            for (int i = 0; i < count; i++) {
                if (in.remaining() < 25) return null;
                long id = in.getLong();
                int kind = Byte.toUnsignedInt(in.get());
                float minWidth = in.getFloat(), maxWidth = in.getFloat(), minHeight = in.getFloat(), maxHeight = in.getFloat();
                String text = readRequiredString(in);
                String imageSource = readRequiredString(in);
                if (in.remaining() < 4) return null;
                int styleLength = in.getInt();
                if (styleLength < 0 || styleLength > 1048576 || styleLength > in.remaining()) return null;
                byte[] typedStyle = new byte[styleLength];
                in.get(typedStyle);
                measured.add(measureNativeControl(id, kind, text, imageSource, typedStyle, minWidth, maxWidth, minHeight, maxHeight));
            }
            if (in.hasRemaining()) return null;
            int capacity = 6;
            for (NativeMeasurement item : measured) capacity += 20 + item.error.getBytes(StandardCharsets.UTF_8).length;
            ByteBuffer out = ByteBuffer.allocate(capacity).order(ByteOrder.LITTLE_ENDIAN);
            out.putShort((short) 1).putInt(measured.size());
            for (NativeMeasurement item : measured) {
                byte[] error = item.error.getBytes(StandardCharsets.UTF_8);
                out.putLong(item.id).putFloat(item.width).putFloat(item.height).putInt(error.length).put(error);
            }
            return out.array();
        } catch (Throwable error) {
            android.util.Log.e("GoNative", "Native measurement batch failed", error);
            return null;
        }
    }

    private NativeMeasurement measureNativeControl(long id, int kind, String text, String imageSource, byte[] typedStyle,
                                                    float minWidth, float maxWidth, float minHeight, float maxHeight) {
        try {
            View view = makeView(kind, false);
            if (view instanceof TextView) {
                TextView label = (TextView) view;
                label.setText(text);
                label.setIncludeFontPadding(false);
                if (view instanceof Button) {
                    Button button = (Button) view;
                    button.setAllCaps(false);
                    button.setMinWidth(0);
                    button.setMinimumWidth(0);
                    button.setMinHeight(dp(44));
                    button.setMinimumHeight(dp(44));
                    button.setPadding(dp(16), 0, dp(16), 0);
                }
                if (view instanceof EditText) {
                    EditText field = (EditText) view;
                    field.setSingleLine(true);
                    field.setMinWidth(dp(240));
                    field.setMinimumWidth(dp(240));
                    field.setMinHeight(dp(44));
                    field.setMinimumHeight(dp(44));
                    field.setPadding(dp(12), dp(8), dp(12), dp(8));
                }
            }
            if (view instanceof ImageView && imageSource != null && !imageSource.isEmpty()) {
                int resource = getResources().getIdentifier(imageSource, "drawable", getPackageName());
                if (resource == 0) resource = getResources().getIdentifier(imageSource, "mipmap", getPackageName());
                if (resource != 0) ((ImageView) view).setImageResource(resource);
            }
            applyTypedStyle(view, typedStyle);
            int widthSpec = nativeMeasureSpec(maxWidth);
            int heightSpec = nativeMeasureSpec(maxHeight);
            view.measure(widthSpec, heightSpec);
            float density = getResources().getDisplayMetrics().density;
            float width = Math.max(finiteNonNegative(minWidth), view.getMeasuredWidth() / density);
            float height = Math.max(finiteNonNegative(minHeight), view.getMeasuredHeight() / density);
            // Normalized control shells are configured before measurement. An
            // explicit maximum remains authoritative for compact variants.
            if (isFinite(maxWidth) && maxWidth >= 0) width = Math.min(width, maxWidth);
            if (isFinite(maxHeight) && maxHeight >= 0) height = Math.min(height, maxHeight);
            return new NativeMeasurement(id, width, height, "");
        } catch (Throwable error) {
            return new NativeMeasurement(id, 0, 0, error.getClass().getSimpleName() + ": " + String.valueOf(error.getMessage()));
        }
    }

    private int nativeMeasureSpec(float maximum) {
        if (!isFinite(maximum) || maximum <= 0) return View.MeasureSpec.makeMeasureSpec(0, View.MeasureSpec.UNSPECIFIED);
        return View.MeasureSpec.makeMeasureSpec(dp(maximum), View.MeasureSpec.AT_MOST);
    }

    private static boolean isFinite(float value) { return !Float.isNaN(value) && !Float.isInfinite(value); }
    private static float finiteNonNegative(float value) { return isFinite(value) ? Math.max(0, value) : 0; }

    private static final class NativeMeasurement {
        final long id; final float width, height; final String error;
        NativeMeasurement(long id, float width, float height, String error) { this.id = id; this.width = width; this.height = height; this.error = error == null ? "" : error; }
    }

    @SuppressWarnings("unused")
    public void applyMutationBatch(byte[] payload) {
        final byte[] owned = payload.clone();
        runOnUiThread(new Runnable() {
            @Override public void run() { applyOnUiThread(owned); }
        });
    }

    private void applyOnUiThread(byte[] payload) {
        if (payload == null) return;
        try {
            long started = System.nanoTime();
            ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
            if (in.remaining() < 14) return;
            if (Short.toUnsignedInt(in.getShort()) != 9) return;
            int count = in.getInt();
            long sequence = in.getLong();
            for (int operation = 0; operation < count && in.hasRemaining(); operation++) {
                int mutation = Byte.toUnsignedInt(in.get());
                int kind = Byte.toUnsignedInt(in.get());
                long nodeID = in.getLong();
                long parentID = in.getLong();
                int index = in.getInt();
                int fromIndex = in.getInt();
                float width = in.getFloat();
                float height = in.getFloat();
                float padding = in.getFloat();
                float gap = in.getFloat();
                int alignment = Byte.toUnsignedInt(in.get());
                boolean bold = in.get() != 0;
                float fontSize = in.getFloat();
                long handler = in.getLong();
                long changeHandler = in.getLong();
                long toggleHandler = in.getLong();
                boolean checked = in.get() != 0;
                float progress = in.getFloat();
                String text = readString(in);
                String accessibility = readString(in);
                String hint = readString(in);
                int role = Byte.toUnsignedInt(in.get());
                boolean focused = in.get() != 0;
                boolean scalesText = in.get() != 0;
                String imageSource = readString(in);
                int imageMode = Byte.toUnsignedInt(in.get());
                boolean horizontal = in.get() != 0;
                int interactionLength = in.remaining() >= 4 ? in.getInt() : 0;
                byte[] interactions = new byte[Math.max(0, Math.min(interactionLength, in.remaining()))];
                if (interactions.length > 0) in.get(interactions);
                int styleLength = in.remaining() >= 4 ? in.getInt() : -1;
                if (styleLength < 0 || styleLength > 1048576 || styleLength > in.remaining()) return;
                byte[] typedStyle = new byte[styleLength];
                in.get(typedStyle);
                if (in.remaining() < 17) return;
                boolean hasFrame = in.get() != 0;
                float frameX = in.getFloat(), frameY = in.getFloat(), frameWidth = in.getFloat(), frameHeight = in.getFloat();
                View view = views.get(nodeID);

                if (mutation == CREATE) {
                    view = makeView(kind, horizontal);
                    view.setTag(nodeID);
                    final long focusNodeID = nodeID;
                    view.setOnFocusChangeListener(new View.OnFocusChangeListener() {
                        @Override public void onFocusChange(View changed, boolean hasFocus) { nativeDispatchFocus(focusNodeID, hasFocus); }
                    });
                    views.put(nodeID, view);
                    if (rootNodeID == 0) rootNodeID = nodeID;
                    style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                    applyTypedStyle(view, typedStyle);
                    applyComputedFrame(nodeID, view, hasFrame, frameX, frameY, frameWidth, frameHeight);
                    applyInteractions(nodeID, view, interactions);
                    if (views.size() == 1) {
                        view.setLayoutParams(new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT));
                        view.setBackgroundColor(android.graphics.Color.WHITE);
                        setContentView(view);
                    }
                } else if (mutation == UPDATE) {
                    if (view != null) {
                        style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                        applyTypedStyle(view, typedStyle);
                        applyComputedFrame(nodeID, view, hasFrame, frameX, frameY, frameWidth, frameHeight);
                        applyInteractions(nodeID, view, interactions);
                    }
                } else if (mutation == INSERT) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        if (parent instanceof ScrollView || parent instanceof HorizontalScrollView) {
                            parent.removeAllViews();
                            ViewGroup.LayoutParams lp = new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
                            parent.addView(view, lp);
                        } else {
                            parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                        }
                    }
                } else if (mutation == REMOVE) {
                    detach(view);
                } else if (mutation == MOVE) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        if (parent instanceof ScrollView || parent instanceof HorizontalScrollView) {
                            parent.removeAllViews();
                            ViewGroup.LayoutParams lp = new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
                            parent.addView(view, lp);
                        } else {
                            parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                        }
                    }
                } else if (mutation == DELETE) {
                    detach(view);
                    if (view instanceof EditText) textWatchers.remove((EditText) view);
                    GestureBinding binding = gestureBindings.get(nodeID);
                    if (binding != null) binding.dispose();
                    gestureBindings.remove(nodeID);
                    if (view != null) view.animate().cancel();
                    views.remove(nodeID);
                    if (nodeID == rootNodeID) rootNodeID = 0;
                }
                // Retained in the cross-platform protocol for renderer diagnostics.
                if (fromIndex == Integer.MIN_VALUE) throw new AssertionError();
            }
            nativeReportBatchApplied(sequence, System.nanoTime() - started);
        } catch (Throwable t) {
            android.util.Log.e("GoNative", "Error applying mutation batch", t);
        }
    }

    private void applyComputedFrame(long nodeID, View view, boolean hasFrame, float x, float y, float width, float height) {
        if (!hasFrame || view == null || nodeID == rootNodeID) return;
        if (!isFinite(x) || !isFinite(y) || !isFinite(width) || !isFinite(height) || width < 0 || height < 0) return;
        int measuredWidth = dp(Math.min(width, 1000000f));
        int measuredHeight = dp(Math.min(height, 1000000f));
        ViewGroup.LayoutParams params = view.getLayoutParams();
        if (params == null) params = new ViewGroup.LayoutParams(measuredWidth, measuredHeight);
        params.width = measuredWidth;
        params.height = measuredHeight;
        view.setLayoutParams(params);
        final float targetX = Math.max(-1000000f, Math.min(1000000f, x));
        final float targetY = Math.max(-1000000f, Math.min(1000000f, y));
        view.post(new Runnable() {
            @Override public void run() { view.setX(dp(targetX)); view.setY(dp(targetY)); }
        });
    }

    private View makeView(int kind, boolean horizontal) {
        if (kind == TEXT) {
            TextView tv = new TextView(this);
            tv.setTextColor(android.graphics.Color.BLACK);
            return tv;
        }
        if (kind == BUTTON) return new Button(this);
        if (kind == TEXT_INPUT) return new EditText(this);
        if (kind == SWITCH) return new Switch(this);
        if (kind == PROGRESS_INDICATOR) { ProgressBar bar = new ProgressBar(this, null, android.R.attr.progressBarStyleHorizontal); bar.setMax(10000); return bar; }
        if (kind == IMAGE) return new ImageView(this);
        if (kind == SCROLL_VIEW) {
            if (horizontal) {
                HorizontalScrollView hsv = new HorizontalScrollView(this);
                hsv.setFillViewport(true);
                return hsv;
            } else {
                ScrollView sv = new ScrollView(this);
                sv.setFillViewport(true);
                return sv;
            }
        }
        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(kind == ROW ? LinearLayout.HORIZONTAL : LinearLayout.VERTICAL);
        if (kind == SAFE_AREA) layout.setFitsSystemWindows(true);
        return layout;
    }

    private void applyTypedStyle(View view, byte[] payload) {
        if (view == null || payload == null || payload.length < 187) return;
        ByteBuffer style = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
        if (Short.toUnsignedInt(style.getShort(0)) != 1) return;
        int portable = 2, ios = portable + typedStyleSize(style, portable, payload.length), androidStyle = ios + typedStyleSize(style, ios, payload.length);
        if (androidStyle <= ios || androidStyle >= payload.length || typedStyleSize(style, androidStyle, payload.length) == 0) return;
        int appearanceBase = hasTypedValues(payload, androidStyle + 112, 69) ? androidStyle + 112 : portable + 112;
        int androidFontLength = style.getInt(androidStyle + 181);
        int textBase = hasTypedValues(payload, androidStyle + 181, 22 + Math.max(0, androidFontLength)) ? androidStyle : portable;
        int fontLength = style.getInt(textBase + 181), interactionBase = hasTypedValues(payload, androidStyle + 203 + Math.max(0, androidFontLength), 17) ? androidStyle : portable;
        int background = rgba(style, appearanceBase), foreground = rgba(style, appearanceBase + 4);
        float borderWidth = style.getFloat(appearanceBase + 8), cornerRadius = style.getFloat(appearanceBase + 16), opacity = style.getFloat(appearanceBase + 44);
        int borderColor = rgba(style, appearanceBase + 12), visibility = Byte.toUnsignedInt(style.get(appearanceBase + 68));
        int disabledOffset = interactionBase + 203 + style.getInt(interactionBase + 181);
        if (fontLength < 0 || disabledOffset >= payload.length) return;
        int fontOffset = textBase + 185 + fontLength;
        float fontSize = style.getFloat(fontOffset), lineHeight = style.getFloat(fontOffset + 6), letterSpacing = style.getFloat(fontOffset + 10);
        int fontWeight = Short.toUnsignedInt(style.getShort(fontOffset + 4));
        if (view instanceof TextView) {
            TextView text = (TextView) view;
            String family = new String(payload, textBase + 185, fontLength, java.nio.charset.StandardCharsets.UTF_8);
            Typeface face = family.isEmpty() ? Typeface.DEFAULT : Typeface.create(family, Typeface.NORMAL);
            text.setTypeface(face, fontWeight >= 600 ? Typeface.BOLD : Typeface.NORMAL);
            if (fontSize > 0) text.setTextSize(TypedValue.COMPLEX_UNIT_DIP, fontSize);
            if (lineHeight > 0 && android.os.Build.VERSION.SDK_INT >= 28) text.setLineHeight(dp(lineHeight));
            if (letterSpacing != 0 && fontSize > 0) text.setLetterSpacing(letterSpacing / fontSize);
        }
        if (android.graphics.Color.alpha(background) > 0 || borderWidth > 0) {
            android.graphics.drawable.GradientDrawable drawable = new android.graphics.drawable.GradientDrawable();
            drawable.setColor(background);
            if (cornerRadius > 0) drawable.setCornerRadius(dp(cornerRadius));
            if (borderWidth > 0) drawable.setStroke(dp(borderWidth), borderColor);
            view.setBackground(drawable);
        }
        if (view instanceof TextView && android.graphics.Color.alpha(foreground) > 0) ((TextView) view).setTextColor(foreground);
        if (opacity > 0) view.setAlpha(Math.min(1f, opacity));
        float translateX = style.getFloat(appearanceBase + 48), translateY = style.getFloat(appearanceBase + 52), scaleX = style.getFloat(appearanceBase + 56), scaleY = style.getFloat(appearanceBase + 60), rotation = style.getFloat(appearanceBase + 64);
        view.setTranslationX(dp(translateX)); view.setTranslationY(dp(translateY));
        view.setScaleX(scaleX == 0 ? 1 : scaleX); view.setScaleY(scaleY == 0 ? 1 : scaleY); view.setRotation(rotation);
        float shadowBlur = style.getFloat(appearanceBase + 32), shadowOpacity = style.getFloat(appearanceBase + 40);
        if (shadowBlur > 0 && shadowOpacity > 0) view.setElevation(dp(shadowBlur));
        view.setVisibility(visibility == 2 ? View.GONE : visibility == 1 ? View.INVISIBLE : View.VISIBLE);
        view.setEnabled(style.get(disabledOffset) == 0);
    }

    private int rgba(ByteBuffer style, int offset) {
        return android.graphics.Color.argb(Byte.toUnsignedInt(style.get(offset + 3)), Byte.toUnsignedInt(style.get(offset)), Byte.toUnsignedInt(style.get(offset + 1)), Byte.toUnsignedInt(style.get(offset + 2)));
    }

    private int typedStyleSize(ByteBuffer style, int base, int limit) {
        if (base < 0 || base + 185 > limit) return 0;
        int fontLength = style.getInt(base + 181);
        int size = 220 + fontLength;
        return fontLength < 0 || base + size > limit ? 0 : size;
    }

    private boolean hasTypedValues(byte[] payload, int offset, int length) {
        if (offset < 0 || length < 0 || offset + length > payload.length) return false;
        for (int i = offset; i < offset + length; i++) if (payload[i] != 0) return true;
        return false;
    }

    private void style(View view, int kind, String text, float width, float height, float padding,
                       float gap, int alignment, float fontSize, boolean bold, long handler, long changeHandler, long toggleHandler, boolean checked, float progress,
                       String accessibility, String hint, int role, boolean focused, boolean scalesText, String imageSource, int imageMode) {
        if (kind == TEXT && view instanceof TextView) {
            TextView textView = (TextView) view;
            textView.setText(text);
            if (fontSize > 0) textView.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            textView.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            textView.setIncludeFontPadding(false);
            textView.setTextColor(android.graphics.Color.parseColor("#111111"));
        }
        if (kind == BUTTON && view instanceof Button) {
            Button btn = (Button) view;
            btn.setText(text);
            if (fontSize > 0) btn.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            btn.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            btn.setAllCaps(false);
            btn.setIncludeFontPadding(false);
            btn.setGravity(Gravity.CENTER);
            btn.setMinHeight(0);
            btn.setMinimumHeight(0);
            btn.setElevation(0);
            android.graphics.drawable.GradientDrawable btnBg = new android.graphics.drawable.GradientDrawable();
            btnBg.setColor(android.graphics.Color.parseColor("#007AFF"));
            btnBg.setCornerRadius(dp(8));
            btn.setBackground(btnBg);
            btn.setTextColor(android.graphics.Color.WHITE);
            btn.setPadding(dp(16), 0, dp(16), 0);
            final long eventHandler = handler;
            if (eventHandler != 0) {
                view.setOnClickListener(new View.OnClickListener() {
                    @Override public void onClick(View clicked) { nativeDispatchEvent(eventHandler); }
                });
            } else {
                view.setOnClickListener(null);
            }
        }
        if (view instanceof EditText) {
            EditText field = (EditText) view;
            field.setSingleLine(true);
            field.setGravity(Gravity.CENTER_VERTICAL);
            field.setIncludeFontPadding(false);
            field.setMinHeight(dp(44));
            if (hint != null && !hint.isEmpty()) field.setHint(hint);
            field.setTextColor(android.graphics.Color.BLACK);
            field.setHintTextColor(android.graphics.Color.parseColor("#8E8E93"));
            android.graphics.drawable.GradientDrawable fieldBg = new android.graphics.drawable.GradientDrawable();
            fieldBg.setColor(android.graphics.Color.parseColor("#FAFAFC"));
            fieldBg.setCornerRadius(dp(8));
            fieldBg.setStroke(dp(1), android.graphics.Color.parseColor("#D1D1D6"));
            field.setBackground(fieldBg);
            int padX = dp(12), padY = dp(8);
            field.setPadding(padX, padY, padX, padY);
            TextWatcher existing = textWatchers.get(field);
            if (existing != null) field.removeTextChangedListener(existing);
            if (text != null && !field.getText().toString().equals(text)) { field.setText(text); field.setSelection(field.length()); }
            field.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize > 0 ? fontSize : 16);
            field.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            final long eventHandler = changeHandler;
            TextWatcher watcher = new TextWatcher() {
                public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                public void onTextChanged(CharSequence s, int start, int before, int count) { if (eventHandler != 0) nativeDispatchValueEvent(eventHandler, s.toString()); }
                public void afterTextChanged(Editable s) {}
            };
            field.addTextChangedListener(watcher); textWatchers.put(field, watcher);
        }
        if (view instanceof Switch) {
            Switch toggle = (Switch) view;
            toggle.setOnCheckedChangeListener(null); toggle.setChecked(checked);
            final long eventHandler = toggleHandler;
            if (eventHandler != 0) {
                toggle.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
                    @Override public void onCheckedChanged(CompoundButton button, boolean value) { nativeDispatchBoolEvent(eventHandler, value); }
                });
            }
        }
        if (view instanceof ProgressBar) { ((ProgressBar) view).setProgress(Math.round(progress * 10000)); }
        if (view instanceof ImageView) {
            ImageView image = (ImageView) view;
            if (imageSource != null && !imageSource.isEmpty()) {
                int resource = getResources().getIdentifier(imageSource, "drawable", getPackageName());
                if (resource == 0) resource = getResources().getIdentifier(imageSource, "mipmap", getPackageName());
                if (resource != 0) {
                    image.setImageResource(resource);
                } else if ("app_logo".equals(imageSource)) {
                    android.graphics.drawable.GradientDrawable gd = new android.graphics.drawable.GradientDrawable();
                    gd.setColor(android.graphics.Color.parseColor("#007AFF"));
                    gd.setCornerRadius(dp(16));
                    image.setBackground(gd);
                    image.setImageResource(android.R.drawable.ic_lock_lock);
                    image.setColorFilter(android.graphics.Color.WHITE);
                    int pad = dp(12);
                    image.setPadding(pad, pad, pad, pad);
                } else if ("avatar".equals(imageSource)) {
                    android.graphics.drawable.GradientDrawable gd = new android.graphics.drawable.GradientDrawable();
                    gd.setColor(android.graphics.Color.parseColor("#007AFF"));
                    gd.setShape(android.graphics.drawable.GradientDrawable.OVAL);
                    image.setBackground(gd);
                    image.setImageResource(android.R.drawable.ic_menu_myplaces);
                    image.setColorFilter(android.graphics.Color.WHITE);
                    int pad = dp(10);
                    image.setPadding(pad, pad, pad, pad);
                } else {
                    image.setImageDrawable(null);
                }
            } else {
                image.setImageDrawable(null);
            }
            image.setScaleType(imageMode == 1 ? ImageView.ScaleType.CENTER_CROP : imageMode == 2 ? ImageView.ScaleType.CENTER : ImageView.ScaleType.FIT_CENTER);
        }
        int paddingPx = dp(padding);
        if (view instanceof EditText) {
            view.setPadding(dp(12) + paddingPx, dp(8) + paddingPx, dp(12) + paddingPx, dp(8) + paddingPx);
        } else {
            view.setPadding(paddingPx, paddingPx, paddingPx, paddingPx);
        }
        view.setContentDescription(accessibility.isEmpty() ? text : accessibility);
        view.setAccessibilityDelegate(new View.AccessibilityDelegate() {
            @Override public void onInitializeAccessibilityNodeInfo(View host, AccessibilityNodeInfo info) {
                super.onInitializeAccessibilityNodeInfo(host, info);
                if (role == 1 || role == 3) info.setClassName(TextView.class.getName());
                else if (role == 2) info.setClassName(Button.class.getName());
                else if (role == 4) info.setClassName("android.widget.ImageView");
                if (android.os.Build.VERSION.SDK_INT >= 26 && !hint.isEmpty()) info.setHintText(hint);
                if (android.os.Build.VERSION.SDK_INT >= 28 && role == 3) info.setHeading(true);
            }
        });
        if (focused) {
            view.requestFocus();
            view.sendAccessibilityEvent(AccessibilityEvent.TYPE_VIEW_FOCUSED);
        } else if (view.hasFocus()) {
            view.clearFocus();
        }
        ViewGroup.LayoutParams current = view.getLayoutParams();
        int requestedWidth = width > 0 ? dp(width) : ViewGroup.LayoutParams.WRAP_CONTENT;
        int requestedHeight = height > 0 ? dp(height) : ViewGroup.LayoutParams.WRAP_CONTENT;
        if (current == null) current = new LinearLayout.LayoutParams(requestedWidth, requestedHeight);
        current.width = requestedWidth;
        current.height = requestedHeight;
        view.setLayoutParams(current);
        if (view instanceof LinearLayout) {
            LinearLayout container = (LinearLayout) view;
            container.setGravity(alignment == 1 ? Gravity.CENTER : alignment == 2 ? Gravity.END : Gravity.START);
            container.setShowDividers(gap > 0 ? LinearLayout.SHOW_DIVIDER_MIDDLE : LinearLayout.SHOW_DIVIDER_NONE);
            if (gap > 0) container.setDividerDrawable(new GapDrawable(dp(gap), container.getOrientation()));
        }
    }

    private int dp(float value) { return Math.round(value * getResources().getDisplayMetrics().density); }

    private void applyInteractions(long nodeID, View view, byte[] payload) {
        if (view == null) return;
        GestureBinding old = gestureBindings.get(nodeID);
        if (old != null) old.dispose();
        gestureBindings.remove(nodeID);
        if (payload == null || payload.length < 4) return;
        ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
        ArrayList<GestureSpec> gestures = new ArrayList<>();
        int gestureCount = in.getInt();
        for (int i = 0; i < gestureCount && in.remaining() >= 26; i++) {
            GestureSpec spec = new GestureSpec();
            spec.kind = Byte.toUnsignedInt(in.get());
            spec.direction = Byte.toUnsignedInt(in.get());
            spec.minimumPressNanos = in.getLong();
            spec.minimumTravel = in.getFloat();
            spec.handler = in.getLong();
            gestures.add(spec);
        }
        if (!gestures.isEmpty()) {
            GestureBinding binding = new GestureBinding(view, gestures);
            gestureBindings.put(nodeID, binding);
            view.setOnTouchListener(binding);
        } else view.setOnTouchListener(null);

        if (in.remaining() < 4) return;
        int animationCount = in.getInt();
        ArrayList<Animator> animations = new ArrayList<>();
        for (int i = 0; i < animationCount && in.remaining() >= 42; i++) {
            int property = Byte.toUnsignedInt(in.get());
            long durationNanos = in.getLong();
            long delayNanos = in.getLong();
            int curve = Byte.toUnsignedInt(in.get());
            float damping = in.getFloat();
            float velocity = in.getFloat();
            boolean reduceMotionOK = in.get() != 0;
            float from = in.getFloat(), to = in.getFloat(), fromX = in.getFloat(), fromY = in.getFloat(), toX = in.getFloat(), toY = in.getFloat();
            Animator animator = makeAnimator(view, property, from, to, fromX, fromY, toX, toY);
            if (animator == null) continue;
            animator.setDuration(Math.max(0, durationNanos / 1000000L));
            animator.setStartDelay(Math.max(0, delayNanos / 1000000L));
            if (curve == 1) animator.setInterpolator(new AccelerateInterpolator());
            else if (curve == 2) animator.setInterpolator(new DecelerateInterpolator());
            else if (curve == 3) animator.setInterpolator(new LinearInterpolator());
            else if (curve == 4) animator.setInterpolator(new OvershootInterpolator(Math.max(.1f, (1f - damping) * 2f + Math.abs(velocity) * .1f)));
            else animator.setInterpolator(new AccelerateDecelerateInterpolator());
            if (!reduceMotionOK && android.os.Build.VERSION.SDK_INT >= 26 && !ValueAnimator.areAnimatorsEnabled()) animator.setDuration(0);
            animations.add(animator);
        }
        if (!animations.isEmpty()) { AnimatorSet set = new AnimatorSet(); set.playSequentially(animations); set.start(); }
    }

    private Animator makeAnimator(final View view, int property, float from, float to, float fromX, float fromY, float toX, float toY) {
        if (property == 1) return ObjectAnimator.ofFloat(view, View.ALPHA, from, to);
        if (property == 2) { AnimatorSet set = new AnimatorSet(); set.playTogether(ObjectAnimator.ofFloat(view, View.SCALE_X, from, to), ObjectAnimator.ofFloat(view, View.SCALE_Y, from, to)); return set; }
        if (property == 3) { AnimatorSet set = new AnimatorSet(); set.playTogether(ObjectAnimator.ofFloat(view, View.TRANSLATION_X, dp(fromX), dp(toX)), ObjectAnimator.ofFloat(view, View.TRANSLATION_Y, dp(fromY), dp(toY))); return set; }
        if (property == 4) {
            ValueAnimator animator = ValueAnimator.ofFloat(0, 1);
            animator.addUpdateListener(new ValueAnimator.AnimatorUpdateListener() {
                @Override public void onAnimationUpdate(ValueAnimator animation) { view.requestLayout(); }
            });
            return animator;
        }
        return null;
    }

    private static final class GestureSpec { int kind, direction; long minimumPressNanos, handler; float minimumTravel; }

    private final class GestureBinding implements View.OnTouchListener {
        final View view; final List<GestureSpec> specs; final ArrayList<Runnable> pending = new ArrayList<>();
        VelocityTracker velocity; float downX, downY; long downTime; boolean moved;
        GestureBinding(View view, List<GestureSpec> specs) { this.view = view; this.specs = specs; }
        void dispose() { for (Runnable r : pending) view.removeCallbacks(r); pending.clear(); if (velocity != null) velocity.recycle(); velocity = null; view.setOnTouchListener(null); }
        void emit(GestureSpec s, float x, float y, float vx, float vy) { if (s.handler != 0) nativeDispatchGestureEvent(s.handler, x / getResources().getDisplayMetrics().density, y / getResources().getDisplayMetrics().density, vx / getResources().getDisplayMetrics().density, vy / getResources().getDisplayMetrics().density); }
        @Override public boolean onTouch(View ignored, MotionEvent event) {
            if (event.getActionMasked() == MotionEvent.ACTION_DOWN) {
                downX=event.getX(); downY=event.getY(); downTime=System.nanoTime(); moved=false;
                velocity=VelocityTracker.obtain(); velocity.addMovement(event);
                for (final GestureSpec s : specs) if (s.kind == 2) {
                    Runnable r = new Runnable() {
                        @Override public void run() { if (!moved) emit(s, 0, 0, 0, 0); }
                    };
                    pending.add(r);
                    view.postDelayed(r, Math.max(0, s.minimumPressNanos / 1000000L));
                }
                return true;
            }
            if (velocity != null) velocity.addMovement(event);
            float dx=event.getX()-downX, dy=event.getY()-downY;
            if (event.getActionMasked() == MotionEvent.ACTION_MOVE) {
                for (GestureSpec s:specs) if (Math.hypot(dx,dy) >= dp(s.minimumTravel)) { moved=true; if(s.kind==4) { velocity.computeCurrentVelocity(1000); emit(s,dx,dy,velocity.getXVelocity(),velocity.getYVelocity()); } }
                if(moved) { for(Runnable r:pending)view.removeCallbacks(r); pending.clear(); }
            } else if (event.getActionMasked() == MotionEvent.ACTION_UP) {
                for(Runnable r:pending)view.removeCallbacks(r); pending.clear(); velocity.computeCurrentVelocity(1000); float vx=velocity.getXVelocity(),vy=velocity.getYVelocity();
                for(GestureSpec s:specs) {
                    float distance=(float)Math.hypot(dx,dy);
                    if(s.kind==1 && distance < dp(Math.max(8,s.minimumTravel))) emit(s,dx,dy,vx,vy);
                    if(s.kind==3 && distance >= dp(s.minimumTravel) && directionMatches(s.direction,dx,dy)) emit(s,dx,dy,vx,vy);
                    if(s.kind==4) emit(s,dx,dy,vx,vy);
                }
                if (view instanceof Button && !moved) view.performClick();
                velocity.recycle(); velocity=null;
            } else if(event.getActionMasked()==MotionEvent.ACTION_CANCEL) disposePending();
            return true;
        }
        void disposePending(){for(Runnable r:pending)view.removeCallbacks(r);pending.clear();if(velocity!=null){velocity.recycle();velocity=null;}}
        boolean directionMatches(int direction,float dx,float dy){ if(direction==0)return true;if(direction==1)return dy<0&&Math.abs(dy)>=Math.abs(dx);if(direction==2)return dy>0&&Math.abs(dy)>=Math.abs(dx);boolean rtl=view.getLayoutDirection()==View.LAYOUT_DIRECTION_RTL;if(direction==3)return rtl?dx>0:dx<0;if(direction==4)return rtl?dx<0:dx>0;return false; }
    }

    private static void detach(View view) {
        if (view != null && view.getParent() instanceof ViewGroup) ((ViewGroup) view.getParent()).removeView(view);
    }

    private static String readString(ByteBuffer in) {
        if (in.remaining() < 4) return "";
        int length = in.getInt();
        if (length <= 0 || in.remaining() < length) return "";
        byte[] bytes = new byte[length];
        in.get(bytes);
        return new String(bytes, StandardCharsets.UTF_8);
    }

    private static String readRequiredString(ByteBuffer in) {
        if (in.remaining() < 4) throw new IllegalArgumentException("missing string length");
        int length = in.getInt();
        if (length < 0 || length > 1048576 || length > in.remaining()) throw new IllegalArgumentException("invalid string length");
        byte[] bytes = new byte[length];
        in.get(bytes);
        return new String(bytes, StandardCharsets.UTF_8);
    }
}
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
)

var benchmarkOutput string

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
	renderer := androidRenderer{}
	appRuntime = gnruntime.New(app.App, renderer)
	appRuntime.SetLayoutProvider(&layout.Pipeline{Measurer: renderer, Cache: layout.NewMeasurementCache()})
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
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

		"android/bridge/jni.c": fmt.Sprintf(`//go:build android

#include <jni.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>

extern void GoNativeAndroidStart(void);
extern void GoNativeAndroidDispatchEvent(uint64_t handler);
extern void GoNativeAndroidDispatchValueEvent(uint64_t handler, const char *value);
extern void GoNativeAndroidDispatchBoolEvent(uint64_t handler, uint8_t value);
extern void GoNativeAndroidDispatchGestureEvent(uint64_t handler, float translationX, float translationY, float velocityX, float velocityY);
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
Java_dev_gonative_%s_MainActivity_nativeStart(JNIEnv *env, jobject renderer) {
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
Java_dev_gonative_%s_MainActivity_nativeSetLifecycle(JNIEnv *env, jobject renderer, jint state) {
    (void)env;
    (void)renderer;
    if (state >= 0 && state <= 6) GoNativeAndroidSetLifecycle((uint8_t)state);
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeDispatchFocus(JNIEnv *env, jobject renderer, jlong nodeID, jboolean focused) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchFocus((uint64_t)nodeID, focused ? 1 : 0);
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeUpdateViewport(JNIEnv *env, jobject renderer, jfloat width, jfloat height, jfloat scale) {
    (void)env;
    (void)renderer;
    GoNativeAndroidUpdateViewport((float)width, (float)height, (float)scale);
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeDispatchEvent(JNIEnv *env, jobject renderer, jlong handler) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchEvent((uint64_t)handler);
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeDispatchValueEvent(JNIEnv *env, jobject renderer, jlong handler, jstring value) {
    (void)renderer;
    const char *utf8 = value ? (*env)->GetStringUTFChars(env, value, NULL) : "";
    if (utf8) { GoNativeAndroidDispatchValueEvent((uint64_t)handler, utf8); }
    if (value && utf8) { (*env)->ReleaseStringUTFChars(env, value, utf8); }
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeDispatchBoolEvent(JNIEnv *env, jobject renderer, jlong handler, jboolean value) {
    (void)env; (void)renderer; GoNativeAndroidDispatchBoolEvent((uint64_t)handler, value ? 1 : 0);
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeDispatchGestureEvent(JNIEnv *env, jobject renderer, jlong handler, jfloat translationX, jfloat translationY, jfloat velocityX, jfloat velocityY) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchGestureEvent((uint64_t)handler, (float)translationX, (float)translationY, (float)velocityX, (float)velocityY);
}

JNIEXPORT void JNICALL
Java_dev_gonative_%s_MainActivity_nativeStop(JNIEnv *env, jobject renderer) {
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
Java_dev_gonative_%s_MainActivity_nativeReportBatchApplied(JNIEnv *env, jobject renderer, jlong sequence, jlong nativeNanos) {
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
`, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg),

		"android/bridge/stub.go": `//go:build !android

package main

func main() {}
`,
	}
}

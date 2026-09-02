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
		fmt.Sprintf("ios/%s.xcodeproj/project.pbxproj", name): generatePbxproj(name, pkg),
		fmt.Sprintf("ios/%s.xcodeproj/xcshareddata/xcschemes/%s.xcscheme", name, name): generateXcscheme(name),
		"ios/GoNativeRenderer.h": `#import <UIKit/UIKit.h>

@interface GNRootViewController : UIViewController
@end

void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
`,
		"ios/GoNativeRenderer.m": `#import "GoNativeRenderer.h"
#import "counter.h"
#include <time.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };
typedef NS_ENUM(uint8_t, GNNode) { GNView=1, GNText, GNButton, GNRow, GNColumn, GNSafeArea, GNTextInput, GNSwitch, GNProgressIndicator, GNImage, GNScrollView };

@interface GNSafeAreaView : UIView
@end
@implementation GNSafeAreaView
@end

@interface GNAction : NSObject
@property(nonatomic) uint64_t handler;
- (void)invoke;
- (void)change:(UITextField *)sender;
- (void)toggle:(UISwitch *)sender;
@end
@implementation GNAction
- (void)invoke { GoNativeDispatchEvent(self.handler); }
- (void)change:(UITextField *)sender { GoNativeDispatchValueEvent(self.handler, (char *)sender.text.UTF8String); }
- (void)toggle:(UISwitch *)sender { GoNativeDispatchBoolEvent(self.handler, sender.isOn ? 1 : 0); }
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
    if ([view isKindOfClass:UIButton.class]) { [(UIButton*)view setTitle:text forState:UIControlStateNormal]; NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&handler){a=[GNAction new];GNActions[actionKey]=a;[(UIButton*)view addTarget:a action:@selector(invoke) forControlEvents:UIControlEventTouchUpInside];}a.handler=handler; }
    if ([view isKindOfClass:UITextField.class]) { UITextField*f=(UITextField*)view;if(![f.text isEqualToString:text])f.text=text;NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&changeHandler){a=[GNAction new];GNActions[actionKey]=a;[f addTarget:a action:@selector(change:) forControlEvents:UIControlEventEditingChanged];}a.handler=changeHandler; }
    if ([view isKindOfClass:UISwitch.class]) { UISwitch*s=(UISwitch*)view;[s setOn:checked animated:NO];NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&toggleHandler){a=[GNAction new];GNActions[actionKey]=a;[s addTarget:a action:@selector(toggle:) forControlEvents:UIControlEventValueChanged];}a.handler=toggleHandler; }
    if ([view isKindOfClass:UIProgressView.class]) { [(UIProgressView*)view setProgress:progress animated:NO]; }
    if ([view isKindOfClass:UIImageView.class]) { UIImageView*i=(UIImageView*)view;i.image=[UIImage imageNamed:imageSource];i.contentMode=imageMode==1?UIViewContentModeScaleAspectFill:imageMode==2?UIViewContentModeCenter:UIViewContentModeScaleAspectFit;i.clipsToBounds=YES; }
    if ([view isKindOfClass:UIStackView.class]) { UIStackView*s=(UIStackView*)view;s.spacing=gap;s.layoutMarginsRelativeArrangement=YES;s.directionalLayoutMargins=NSDirectionalEdgeInsetsMake(padding,padding,padding,padding);s.alignment=alignment==1?UIStackViewAlignmentCenter:alignment==2?UIStackViewAlignmentTrailing:UIStackViewAlignmentLeading; }
    if(width>0)[view.widthAnchor constraintEqualToConstant:width].active=YES;if(height>0)[view.heightAnchor constraintEqualToConstant:height].active=YES;
    view.isAccessibilityElement=(kind==GNText||kind==GNButton||kind==GNTextInput||kind==GNSwitch||kind==GNProgressIndicator||role!=0);view.accessibilityLabel=accessibility.length?accessibility:text;view.accessibilityHint=hint;UIAccessibilityTraits traits=UIAccessibilityTraitNone;if(role==2||kind==GNButton)traits|=UIAccessibilityTraitButton;if(role==3)traits|=UIAccessibilityTraitHeader;if(role==4)traits|=UIAccessibilityTraitImage;view.accessibilityTraits=traits;if(focused)UIAccessibilityPostNotification(UIAccessibilityScreenChangedNotification,view);
    GNConfigureInteractions(nodeID,view,interactions,animate);
}

static UIView *GNMake(GNNode kind){UIView*v;if(kind==GNText){UILabel*l=[UILabel new];l.numberOfLines=0;v=l;}else if(kind==GNButton){UIButton*b=[UIButton buttonWithType:UIButtonTypeSystem];v=b;}else if(kind==GNTextInput){UITextField*f=[UITextField new];f.borderStyle=UITextBorderStyleRoundedRect;v=f;}else if(kind==GNSwitch){v=[UISwitch new];}else if(kind==GNProgressIndicator){v=[[UIProgressView alloc]initWithProgressViewStyle:UIProgressViewStyleDefault];}else if(kind==GNImage){v=[UIImageView new];}else if(kind==GNScrollView){v=[UIScrollView new];}else if(kind==GNRow||kind==GNColumn){UIStackView*s=[UIStackView new];s.axis=kind==GNRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;v=s;}else if(kind==GNSafeArea){v=[GNSafeAreaView new];}else{v=[UIView new];}v.translatesAutoresizingMaskIntoConstraints=NO;return v;}

static void GNApply(NSData *data){uint64_t started=GNNowNanos();GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};if(u16(&r)!=7)return;uint32_t count=u32(&r);uint64_t sequence=u64(&r);for(uint32_t op=0;op<count;op++){GNMutation mutation=(GNMutation)u8(&r);GNNode kind=(GNNode)u8(&r);uint64_t nodeID=u64(&r),parentID=u64(&r);int32_t index=i32(&r),from=i32(&r);float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);uint8_t alignment=u8(&r);BOOL bold=u8(&r);float fontSize=f32(&r);uint64_t handler=u64(&r),changeHandler=u64(&r),toggleHandler=u64(&r);BOOL checked=u8(&r);float progress=f32(&r);NSString*text=str(&r);NSString*accessibility=str(&r);NSString*hint=str(&r);uint8_t role=u8(&r);BOOL focused=u8(&r);BOOL scalesText=u8(&r);NSString*imageSource=str(&r);uint8_t imageMode=u8(&r);BOOL horizontal=u8(&r);uint32_t interactionLength=u32(&r);NSData *interactions;if(r.p+interactionLength<=r.end){interactions=[NSData dataWithBytes:r.p length:interactionLength];r.p+=interactionLength;}else{r.p=r.end;interactions=[NSData data];}NSNumber*key=@(nodeID);UIView*view=GNViews[key];
    if(mutation==GNCreate){view=GNMake(kind);GNViews[key]=view;GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,NO);if(!GNRoot.view.subviews.count){[GNRoot.view addSubview:view];[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor]]];}}
    else if(mutation==GNUpdate){GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,YES);}
    else if(mutation==GNInsert){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}else{[parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];UILayoutGuide*guide=[parent isKindOfClass:GNSafeAreaView.class]?parent.safeAreaLayoutGuide:nil;if([parent isKindOfClass:UIScrollView.class]){UIScrollView*s=(UIScrollView*)parent;[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:s.contentLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:s.contentLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:s.contentLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:s.contentLayoutGuide.bottomAnchor]]];}else if(guide){[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:guide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:guide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:guide.topAnchor],[view.bottomAnchor constraintLessThanOrEqualToAnchor:guide.bottomAnchor]]];}}}
    else if(mutation==GNRemove){[view removeFromSuperview];}
    else if(mutation==GNMove){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s removeArrangedSubview:view];[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}}
    else if(mutation==GNDelete){[view removeFromSuperview];[GNViews removeObjectForKey:key];[GNActions removeObjectForKey:key];[GNGestureActions removeObjectForKey:key];}
    (void)from;
}GoNativeReportBatchApplied(sequence,GNNowNanos()-started);}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});}

@implementation GNRootViewController
- (void)viewDidLoad {[super viewDidLoad];self.view.backgroundColor=UIColor.systemBackgroundColor;GNRoot=self;GNViews=[NSMutableDictionary dictionary];GNActions=[NSMutableDictionary dictionary];GNGestureActions=[NSMutableDictionary dictionary];GoNativeStart();}
- (void)dealloc { GoNativeStop(); }
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
*/
import "C"

import (
	"errors"
	"fmt"
	"%s"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
	"time"
	"unsafe"
)

var benchmarkOutput string

type iosRenderer struct{}

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
	if err := appRuntime.Start(); err != nil {
		panic(err)
	}
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

//export GoNativeStop
func GoNativeStop() {
	if appRuntime != nil {
		appRuntime.Stop()
		appRuntime = nil
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

import android.animation.Animator;
import android.animation.AnimatorSet;
import android.animation.ObjectAnimator;
import android.animation.ValueAnimator;
import android.app.Activity;
import android.graphics.Typeface;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.LongSparseArray;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.MotionEvent;
import android.view.VelocityTracker;
import android.view.View;
import android.view.ViewGroup;
import android.view.accessibility.AccessibilityEvent;
import android.view.accessibility.AccessibilityNodeInfo;
import android.view.animation.AccelerateDecelerateInterpolator;
import android.view.animation.AccelerateInterpolator;
import android.view.animation.DecelerateInterpolator;
import android.view.animation.LinearInterpolator;
import android.view.animation.OvershootInterpolator;
import android.widget.Button;
import android.widget.CompoundButton;
import android.widget.EditText;
import android.widget.HorizontalScrollView;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.ProgressBar;
import android.widget.ScrollView;
import android.widget.Switch;
import android.widget.TextView;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

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

    private native void nativeStart();
    private native void nativeDispatchEvent(long handler);
    private native void nativeDispatchValueEvent(long handler, String value);
    private native void nativeDispatchBoolEvent(long handler, boolean value);
    private native void nativeDispatchGestureEvent(long handler, float translationX, float translationY, float velocityX, float velocityY);
    private native void nativeStop();
    private native void nativeReportBatchApplied(long sequence, long nativeNanos);

    @Override protected void onCreate(Bundle state) {
        super.onCreate(state);
        nativeStart();
    }

    @Override protected void onDestroy() {
        nativeStop();
        for (int i = 0; i < gestureBindings.size(); i++) gestureBindings.valueAt(i).dispose();
        gestureBindings.clear();
        for (int i = 0; i < views.size(); i++) views.valueAt(i).animate().cancel();
        views.clear();
        super.onDestroy();
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
            if (Short.toUnsignedInt(in.getShort()) != 7) return;
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
                View view = views.get(nodeID);

                if (mutation == CREATE) {
                    view = makeView(kind, horizontal);
                    view.setTag(nodeID);
                    views.put(nodeID, view);
                    style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                    applyInteractions(nodeID, view, interactions);
                    if (views.size() == 1) setContentView(view);
                } else if (mutation == UPDATE) {
                    if (view != null) {
                        style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                        applyInteractions(nodeID, view, interactions);
                    }
                } else if (mutation == INSERT) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                    }
                } else if (mutation == REMOVE) {
                    detach(view);
                } else if (mutation == MOVE) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                    }
                } else if (mutation == DELETE) {
                    detach(view);
                    GestureBinding binding = gestureBindings.get(nodeID);
                    if (binding != null) binding.dispose();
                    gestureBindings.remove(nodeID);
                    if (view != null) view.animate().cancel();
                    views.remove(nodeID);
                }
                if (fromIndex == Integer.MIN_VALUE) throw new AssertionError();
            }
            nativeReportBatchApplied(sequence, System.nanoTime() - started);
        } catch (Throwable t) {
            android.util.Log.e("GoNative", "Error applying mutation batch", t);
        }
    }

    private View makeView(int kind, boolean horizontal) {
        if (kind == TEXT) return new TextView(this);
        if (kind == BUTTON) return new Button(this);
        if (kind == TEXT_INPUT) return new EditText(this);
        if (kind == SWITCH) return new Switch(this);
        if (kind == PROGRESS_INDICATOR) { ProgressBar bar = new ProgressBar(this, null, android.R.attr.progressBarStyleHorizontal); bar.setMax(10000); return bar; }
        if (kind == IMAGE) return new ImageView(this);
        if (kind == SCROLL_VIEW) return horizontal ? new HorizontalScrollView(this) : new ScrollView(this);
        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(kind == ROW ? LinearLayout.HORIZONTAL : LinearLayout.VERTICAL);
        if (kind == SAFE_AREA) layout.setFitsSystemWindows(true);
        return layout;
    }

    private void style(View view, int kind, String text, float width, float height, float padding,
                       float gap, int alignment, float fontSize, boolean bold, long handler, long changeHandler, long toggleHandler, boolean checked, float progress,
                       String accessibility, String hint, int role, boolean focused, boolean scalesText, String imageSource, int imageMode) {
        if (kind == TEXT && view instanceof TextView) {
            TextView textView = (TextView) view;
            textView.setText(text);
            if (fontSize > 0) textView.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            textView.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
        }
        if (kind == BUTTON && view instanceof Button) {
            Button btn = (Button) view;
            btn.setText(text);
            if (fontSize > 0) btn.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            btn.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
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
            Object existing = field.getTag(android.R.id.custom);
            if (existing instanceof TextWatcher) field.removeTextChangedListener((TextWatcher) existing);
            if (text != null && !field.getText().toString().equals(text)) { field.setText(text); field.setSelection(field.length()); }
            if (fontSize > 0) field.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            field.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            final long eventHandler = changeHandler;
            TextWatcher watcher = new TextWatcher() {
                public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                public void onTextChanged(CharSequence s, int start, int before, int count) { if (eventHandler != 0) nativeDispatchValueEvent(eventHandler, s.toString()); }
                public void afterTextChanged(Editable s) {}
            };
            field.addTextChangedListener(watcher); field.setTag(android.R.id.custom, watcher);
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
                if (resource != 0) image.setImageResource(resource);
                else image.setImageDrawable(null);
            } else {
                image.setImageDrawable(null);
            }
            image.setScaleType(imageMode == 1 ? ImageView.ScaleType.CENTER_CROP : imageMode == 2 ? ImageView.ScaleType.CENTER : ImageView.ScaleType.FIT_CENTER);
        }
        int paddingPx = dp(padding);
        view.setPadding(paddingPx, paddingPx, paddingPx, paddingPx);
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
}
`, pkg),

		"android/bridge/main.go": fmt.Sprintf(`//go:build android

package main

/*
#include <stdint.h>
void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length);
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"%s"
	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
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

var appRuntime *gnruntime.Runtime

//export GoNativeAndroidStart
func GoNativeAndroidStart() {
	appRuntime = gnruntime.New(app.App, androidRenderer{})
	if err := appRuntime.Start(); err != nil {
		panic(err)
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

static JavaVM *gn_vm;
static jobject gn_renderer;
static jmethodID gn_apply;
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
    (*env)->DeleteLocalRef(env, cls);
    pthread_mutex_unlock(&gn_renderer_mu);
    GoNativeAndroidStart();
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
    pthread_mutex_unlock(&gn_renderer_mu);
}

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
`, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg, jniPkg),

		"android/bridge/stub.go": `//go:build !android

package main

func main() {}
`,
	}
}

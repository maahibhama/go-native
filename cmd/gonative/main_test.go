package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type recordedRun struct {
	name, dir string
	args, env []string
}
type fakeRunner struct {
	calls []recordedRun
	err   error
}

func (f *fakeRunner) Run(name string, args []string, dir string, env []string, stdout, stderr io.Writer) error {
	f.calls = append(f.calls, recordedRun{name: name, args: args, dir: dir, env: env})
	return f.err
}

func TestHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"help"}, &fakeRunner{}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "gonative build") {
		t.Fatalf("unexpected help: %s", out.String())
	}
}

func TestInitCreatesCompleteNativeScaffold(t *testing.T) {
	parent := t.TempDir()
	var out bytes.Buffer
	if err := initProject(parent, "hello-native", &out); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(parent, "hello-native")

	expectedFiles := []string{
		"go.mod",
		"gonative.yaml",
		"app.go",
		"README.md",
		".gitignore",
		"assets/.gitkeep",
		// iOS & Xcode
		"ios/hello-native.xcodeproj/project.pbxproj",
		"ios/hello-native.xcodeproj/xcshareddata/xcschemes/hello-native.xcscheme",
		"ios/AppDelegate.h",
		"ios/AppDelegate.m",
		"ios/Package.swift",
		"ios/main.m",
		"ios/Info.plist",
		"ios/bridge/main.go",
		// Android & Android Studio
		"android/build.gradle",
		"android/settings.gradle",
		"android/gradle.properties",
		"android/build-libs.sh",
		"android/app/build.gradle",
		"android/app/src/main/AndroidManifest.xml",
		"android/app/src/main/res/values/strings.xml",
		"android/app/src/main/res/values/styles.xml",
		"android/app/src/main/java/dev/gonative/hello_native/MainActivity.java",
		"android/bridge/main.go",
		"android/bridge/jni.c",
		"android/bridge/stub.go",
	}

	for _, rel := range expectedFiles {
		fullPath := filepath.Join(destination, rel)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("missing scaffold file %s: %v", rel, err)
		}
	}

	app, err := os.ReadFile(filepath.Join(destination, "app.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(app), "func App() ui.Component") {
		t.Fatalf("unexpected app.go:\n%s", app)
	}

	manifest, err := os.ReadFile(filepath.Join(destination, "android", "app", "src", "main", "AndroidManifest.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(manifest), "package=") || !strings.Contains(string(manifest), "android:label=\"@string/app_name\"") {
		t.Fatalf("unexpected Android manifest:\n%s", manifest)
	}

	stringsXML, err := os.ReadFile(filepath.Join(destination, "android", "app", "src", "main", "res", "values", "strings.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(stringsXML), `<string name="app_name">hello-native</string>`) {
		t.Fatalf("unexpected strings.xml:\n%s", stringsXML)
	}

	infoPlist, err := os.ReadFile(filepath.Join(destination, "ios", "Info.plist"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(infoPlist), "<key>CFBundleDisplayName</key>\n    <string>hello-native</string>") {
		t.Fatalf("unexpected Info.plist:\n%s", infoPlist)
	}

	pbx, err := os.ReadFile(filepath.Join(destination, "ios", "hello-native.xcodeproj", "project.pbxproj"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pbx), "PBXNativeTarget") {
		t.Fatalf("unexpected pbxproj:\n%s", pbx)
	}
	if !strings.Contains(string(pbx), "AppDelegate.m in Sources") {
		t.Fatalf("pbxproj missing AppDelegate.m in Sources:\n%s", pbx)
	}

	appDelegateH, err := os.ReadFile(filepath.Join(destination, "ios", "AppDelegate.h"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(appDelegateH), "@interface AppDelegate : UIResponder <UIApplicationDelegate>") {
		t.Fatalf("unexpected AppDelegate.h:\n%s", appDelegateH)
	}

	appDelegateM, err := os.ReadFile(filepath.Join(destination, "ios", "AppDelegate.m"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(appDelegateM), "@implementation AppDelegate") || !strings.Contains(string(appDelegateM), "GNRootViewController") {
		t.Fatalf("unexpected AppDelegate.m:\n%s", appDelegateM)
	}

	packageSwift, err := os.ReadFile(filepath.Join(destination, "ios", "Package.swift"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(packageSwift), `.binaryTarget(name: "GoNativeKit"`) {
		t.Fatalf("unexpected Package.swift:\n%s", packageSwift)
	}

	androidBuild, err := os.ReadFile(filepath.Join(destination, "android", "app", "build.gradle"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(androidBuild), `dev.gonative:gonative-runtime:0.1.0`) {
		t.Fatalf("Android app does not use Maven runtime dependency:\n%s", androidBuild)
	}

	mainM, err := os.ReadFile(filepath.Join(destination, "ios", "main.m"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(mainM), "NSStringFromClass([AppDelegate class])") {
		t.Fatalf("unexpected main.m:\n%s", mainM)
	}

	jni, err := os.ReadFile(filepath.Join(destination, "android", "bridge", "jni.c"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(jni), "Java_dev_gonative_runtime_GoNativeActivity_nativeStart") {
		t.Fatalf("unexpected jni.c:\n%s", jni)
	}
	launcher, err := os.ReadFile(filepath.Join(destination, "android", "app", "src", "main", "java", "dev", "gonative", "hello_native", "MainActivity.java"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(launcher), "extends GoNativeActivity") {
		t.Fatalf("Android launcher contains framework implementation:\n%s", launcher)
	}

	iosBridge, err := os.ReadFile(filepath.Join(destination, "ios", "bridge", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	androidBridge, err := os.ReadFile(filepath.Join(destination, "android", "bridge", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	for path, bridge := range map[string][]byte{"ios bridge": iosBridge, "android bridge": androidBridge} {
		if !strings.Contains(string(bridge), "goNativeReloadSession") || !strings.Contains(string(bridge), "ui.ConfigureReloadState") {
			t.Fatalf("%s missing Fast Reload setup", path)
		}
	}

	if !strings.Contains(out.String(), "Created ") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestInitDispatchWorksOutsideFrameworkRepository(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	parent := t.TempDir()
	if err = os.Chdir(parent); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })
	if err = run([]string{"init", "standalone"}, &fakeRunner{}, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(parent, "standalone", "app.go")); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(parent, "standalone", "ios", "main.m")); err != nil {
		t.Fatal(err)
	}
	module, err := os.ReadFile(filepath.Join(parent, "standalone", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(module), "replace github.com/go-native/go-native =>") {
		t.Fatalf("standalone project did not retain its local framework dependency:\n%s", module)
	}
	if _, err = os.Stat(filepath.Join(parent, "standalone", "android", "gradlew")); err != nil {
		t.Fatalf("standalone project is missing Gradle wrapper: %v", err)
	}
}

func TestFrameworkRootEnvironmentOverride(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GONATIVE_FRAMEWORK_ROOT", root)
	got, err := findFrameworkRoot()
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("root=%q want=%q", got, root)
	}
}

func TestInitRejectsUnsafeNames(t *testing.T) {
	for _, name := range []string{"", ".", "../escape", "nested/app", "two words", "9app"} {
		if err := initProject(t.TempDir(), name, io.Discard); err == nil {
			t.Errorf("name %q accepted", name)
		}
	}
}

func TestInitDoesNotOverwrite(t *testing.T) {
	parent := t.TempDir()
	destination := filepath.Join(parent, "existing")
	if err := os.Mkdir(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(destination, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := initProject(parent, "existing", io.Discard); err == nil {
		t.Fatal("expected existing destination error")
	}
	got, err := os.ReadFile(marker)
	if err != nil || string(got) != "keep" {
		t.Fatalf("existing content changed: %q, %v", got, err)
	}
}

func TestBuildIOSDispatchesScript(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{}
	var out bytes.Buffer
	if err = run([]string{"build", "ios"}, runner, &out, &out); err != nil {
		t.Fatal(err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("calls=%d", len(runner.calls))
	}
	want := filepath.Join(root, "scripts", "build-ios.sh")
	if runner.calls[0].name != want {
		t.Fatalf("script=%s want=%s", runner.calls[0].name, want)
	}
}

func TestBuildIOSDeviceDispatchesScript(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{}
	if err = run([]string{"build", "ios-device"}, runner, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "scripts", "build-ios-device.sh")
	if len(runner.calls) != 1 || runner.calls[0].name != want {
		t.Fatalf("calls=%v want script=%s", runner.calls, want)
	}
}

func TestRunIOSDeviceIsRejected(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"run", "ios-device"}, &fakeRunner{}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "unsupported platform") {
		t.Fatalf("error=%v", err)
	}
}

func TestNativeBenchmarkDispatchesScript(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{}
	if err = run([]string{"benchmark", "native", "android"}, runner, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "scripts", "benchmark-native-android.sh")
	if len(runner.calls) != 1 || runner.calls[0].name != want {
		t.Fatalf("calls=%v want script=%s", runner.calls, want)
	}
}

func TestNativeBenchmarkRejectsUnsupportedPlatform(t *testing.T) {
	err := run([]string{"benchmark", "native", "web"}, &fakeRunner{}, io.Discard, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "unsupported native benchmark platform") {
		t.Fatalf("error=%v", err)
	}
}

func TestAndroidGradleProjectReferencesSharedNativeLibraries(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	buildFile, err := os.ReadFile(filepath.Join(root, "platform", "android", "app", "build.gradle"))
	if err != nil {
		t.Fatal(err)
	}
	contents := string(buildFile)
	for _, want := range []string{"prepareGoNativeLibraries", "build/android/lib", "arm64-v8a,x86_64", "../AndroidManifest.xml", "implementation project(\":runtime\")"} {
		if !strings.Contains(contents, want) {
			t.Errorf("Gradle configuration missing %q", want)
		}
	}
}

func TestUnsupportedPlatform(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"build", "web"}, &fakeRunner{}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "unsupported platform") {
		t.Fatalf("error=%v", err)
	}
}

func TestDefaultEnvPreservesValue(t *testing.T) {
	env := defaultEnv([]string{"GOCACHE=/custom"}, "GOCACHE", "/tmp/default")
	if len(env) != 1 || env[0] != "GOCACHE=/custom" {
		t.Fatalf("env=%v", env)
	}
}

func TestFindProjectRootFromChild(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(filepath.Join(root, "ui")); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(original)
	got, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("got=%s want=%s", got, root)
	}
}

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
		"app.go",
		"README.md",
		".gitignore",
		// iOS & Xcode
		"ios/hello-native.xcodeproj/project.pbxproj",
		"ios/hello-native.xcodeproj/xcshareddata/xcschemes/hello-native.xcscheme",
		"ios/main.m",
		"ios/GoNativeRenderer.h",
		"ios/GoNativeRenderer.m",
		"ios/Info.plist",
		"ios/bridge/main.go",
		// Android & Android Studio
		"android/build.gradle",
		"android/settings.gradle",
		"android/gradle.properties",
		"android/build-libs.sh",
		"android/app/build.gradle",
		"android/app/src/main/AndroidManifest.xml",
		"android/app/src/main/res/values/styles.xml",
		"android/app/src/main/java/dev/gonative/hello_native/MainActivity.java",
		"android/app/src/main/java/dev/gonative/hello_native/GapDrawable.java",
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
	if !strings.Contains(string(manifest), "package=\"dev.gonative.hello_native\"") {
		t.Fatalf("unexpected Android manifest:\n%s", manifest)
	}

	pbx, err := os.ReadFile(filepath.Join(destination, "ios", "hello-native.xcodeproj", "project.pbxproj"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pbx), "PBXNativeTarget") {
		t.Fatalf("unexpected pbxproj:\n%s", pbx)
	}

	jni, err := os.ReadFile(filepath.Join(destination, "android", "bridge", "jni.c"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(jni), "Java_dev_gonative_hello_1native_MainActivity_nativeStart") {
		t.Fatalf("unexpected jni.c:\n%s", jni)
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
	for _, want := range []string{"prepareGoNativeLibraries", "build/android/lib", "arm64-v8a,x86_64", "../AndroidManifest.xml", "androidx.recyclerview:recyclerview"} {
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

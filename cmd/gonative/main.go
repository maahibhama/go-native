// Command gonative provides the developer workflow for Go Native applications.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const usage = `Go Native CLI

Usage:
  gonative init <name>
  gonative build <ios|ios-device|android>
  gonative run <ios|android>
  gonative benchmark native <ios|android>
  gonative doctor
  gonative help

Environment:
  GONATIVE_SIMULATOR                  iOS Simulator device name
  GONATIVE_IOS_SIGNING_IDENTITY       Apple code-signing identity
  GONATIVE_IOS_PROVISIONING_PROFILE   Path to a .mobileprovision file
  GONATIVE_ANDROID_ABIS               Comma-separated Android ABIs
  GONATIVE_ANDROID_SERIAL             adb device serial
  ANDROID_SDK_ROOT                    Android SDK location
`

type commandRunner interface {
	Run(name string, args []string, dir string, env []string, stdout, stderr io.Writer) error
}

type processRunner struct{}

func (processRunner) Run(name string, args []string, dir string, env []string, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func main() {
	if err := run(os.Args[1:], processRunner{}, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "gonative:", err)
		os.Exit(1)
	}
}

func run(args []string, runner commandRunner, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, usage)
		return nil
	}
	if args[0] == "init" {
		if len(args) != 2 {
			return fmt.Errorf("init requires one project name\n\n%s", usage)
		}
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		return initProject(dir, args[1], stdout)
	}
	if args[0] == "benchmark" {
		if len(args) != 3 || args[1] != "native" {
			return fmt.Errorf("benchmark requires native and one platform\n\n%s", usage)
		}
		root, err := findProjectRoot()
		if err != nil {
			return err
		}
		return nativeBenchmarkCommand(root, args[2], runner, stdout, stderr)
	}
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	switch args[0] {
	case "doctor":
		return doctor(root, stdout)
	case "build", "run":
		if len(args) != 2 {
			return fmt.Errorf("%s requires one platform\n\n%s", args[0], usage)
		}
		return platformCommand(root, args[0], args[1], runner, stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func nativeBenchmarkCommand(root, platform string, runner commandRunner, stdout, stderr io.Writer) error {
	if platform != "ios" && platform != "android" {
		return fmt.Errorf("unsupported native benchmark platform %q", platform)
	}
	script := filepath.Join(root, "scripts", "benchmark-native-"+platform+".sh")
	env := append([]string{}, os.Environ()...)
	env = defaultEnv(env, "GOCACHE", filepath.Join(os.TempDir(), "go-native-gocache"))
	env = defaultEnv(env, "GOPATH", filepath.Join(os.TempDir(), "go-native-gopath"))
	if err := runner.Run(script, nil, root, env, stdout, stderr); err != nil {
		return fmt.Errorf("native %s benchmark failed: %w", platform, err)
	}
	return nil
}

var projectNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`)

func initProject(parent, name string, out io.Writer) error {
	if !projectNamePattern.MatchString(name) {
		return fmt.Errorf("invalid project name %q; use a letter followed by letters, numbers, hyphens, or underscores", name)
	}
	destination := filepath.Join(parent, name)
	if _, err := os.Stat(destination); err == nil {
		return fmt.Errorf("destination %q already exists", destination)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect destination: %w", err)
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	templates := getProjectTemplates(name)
	frameworkRoot, frameworkErr := findFrameworkRoot()
	if frameworkErr == nil {
		if rel, err := filepath.Rel(destination, frameworkRoot); err == nil {
			templates["go.mod"] = fmt.Sprintf("module %s\n\ngo 1.24\n\nrequire github.com/go-native/go-native v0.0.0\n\nreplace github.com/go-native/go-native => %s\n", name, rel)
		}
	}

	for relPath, contents := range templates {
		targetFile := filepath.Join(destination, relPath)
		if err := os.MkdirAll(filepath.Dir(targetFile), 0o755); err != nil {
			return fmt.Errorf("mkdir for %s: %w", relPath, err)
		}
		if err := os.WriteFile(targetFile, []byte(contents), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", relPath, err)
		}
	}

	// Copy Gradle wrapper from framework if available
	if frameworkErr == nil {
		srcGradlew := filepath.Join(frameworkRoot, "platform", "android", "gradlew")
		if _, err := os.Stat(srcGradlew); err == nil {
			_ = copyFile(srcGradlew, filepath.Join(destination, "android", "gradlew"))
			_ = os.Chmod(filepath.Join(destination, "android", "gradlew"), 0o755)
			_ = copyFile(filepath.Join(frameworkRoot, "platform", "android", "gradlew.bat"), filepath.Join(destination, "android", "gradlew.bat"))
			_ = copyDir(filepath.Join(frameworkRoot, "platform", "android", "gradle"), filepath.Join(destination, "android", "gradle"))
		}
	}

	fmt.Fprintf(out, "Created %s with native ios/ and android/ project directories.\n\nNext:\n  cd %s\n  gonative doctor\n  gonative run ios\n  gonative run android\n", destination, name)
	return nil
}

func platformCommand(root, action, platform string, runner commandRunner, stdout, stderr io.Writer) error {
	if platform != "ios" && platform != "android" && !(action == "build" && platform == "ios-device") {
		return fmt.Errorf("unsupported platform %q for %s", platform, action)
	}
	script := filepath.Join(root, "scripts", action+"-"+platform+".sh")
	if _, err := os.Stat(script); err == nil {
		env := append([]string{}, os.Environ()...)
		env = defaultEnv(env, "GOCACHE", filepath.Join(os.TempDir(), "go-native-gocache"))
		env = defaultEnv(env, "GOPATH", filepath.Join(os.TempDir(), "go-native-gopath"))
		if err := runner.Run(script, nil, root, env, stdout, stderr); err != nil {
			return fmt.Errorf("%s %s failed: %w", action, platform, err)
		}
		return nil
	}

	// Standalone project execution
	return runStandalonePlatformCommand(root, action, platform, runner, stdout, stderr)
}

func runStandalonePlatformCommand(root, action, platform string, runner commandRunner, stdout, stderr io.Writer) error {
	env := append([]string{}, os.Environ()...)
	env = defaultEnv(env, "GOCACHE", filepath.Join(os.TempDir(), "go-native-gocache"))
	env = defaultEnv(env, "GOPATH", filepath.Join(os.TempDir(), "go-native-gopath"))

	appName := filepath.Base(root)
	pkg := sanitizePackageName(appName)

	switch {
	case action == "build" && platform == "ios":
		sdkOut, err := exec.Command("xcrun", "--sdk", "iphonesimulator", "--show-sdk-path").Output()
		if err != nil {
			return fmt.Errorf("lookup iOS simulator SDK: %w", err)
		}
		sdk := strings.TrimSpace(string(sdkOut))
		buildDir := filepath.Join(root, "build", "ios-simulator")
		appBundle := filepath.Join(buildDir, appName+".app")
		_ = os.MkdirAll(appBundle, 0o755)

		cgoCC := fmt.Sprintf("clang -target arm64-apple-ios15.0-simulator -isysroot %s", sdk)
		cgoEnv := append(env, "CGO_ENABLED=1", "GOOS=ios", "GOARCH=arm64", "CC="+cgoCC)
		if err := runner.Run("go", []string{"build", "-buildmode=c-archive", "-o", filepath.Join(buildDir, "counter.a"), "./ios/bridge"}, root, cgoEnv, stdout, stderr); err != nil {
			return fmt.Errorf("compile ios bridge: %w", err)
		}

		clangArgs := []string{
			"--sdk", "iphonesimulator", "clang",
			"-target", "arm64-apple-ios15.0-simulator",
			"-isysroot", sdk,
			"-fobjc-arc",
			"-framework", "UIKit",
			"-framework", "Foundation",
			"-framework", "CoreGraphics",
			"-I" + buildDir,
			"-I" + filepath.Join(root, "ios"),
			filepath.Join(root, "ios", "main.m"),
			filepath.Join(root, "ios", "GoNativeRenderer.m"),
			filepath.Join(buildDir, "counter.a"),
			"-o", filepath.Join(appBundle, appName),
		}
		if err := runner.Run("xcrun", clangArgs, root, env, stdout, stderr); err != nil {
			return fmt.Errorf("link ios simulator binary: %w", err)
		}
		_ = copyFile(filepath.Join(root, "ios", "Info.plist"), filepath.Join(appBundle, "Info.plist"))
		_ = runner.Run("codesign", []string{"--force", "--sign", "-", appBundle}, root, env, io.Discard, io.Discard)
		fmt.Fprintln(stdout, appBundle)
		return nil

	case action == "run" && platform == "ios":
		if err := runStandalonePlatformCommand(root, "build", "ios", runner, stdout, stderr); err != nil {
			return err
		}
		simName := os.Getenv("GONATIVE_SIMULATOR")
		if simName == "" {
			simName = "booted"
		}
		appBundle := filepath.Join(root, "build", "ios-simulator", appName+".app")
		_ = runner.Run("xcrun", []string{"simctl", "boot", simName}, root, env, io.Discard, io.Discard)
		if err := runner.Run("xcrun", []string{"simctl", "install", simName, appBundle}, root, env, stdout, stderr); err != nil {
			return fmt.Errorf("simctl install: %w", err)
		}
		bundleID := fmt.Sprintf("dev.gonative.%s", pkg)
		if err := runner.Run("xcrun", []string{"simctl", "launch", simName, bundleID}, root, env, stdout, stderr); err != nil {
			return fmt.Errorf("simctl launch: %w", err)
		}
		return nil

	case action == "build" && platform == "android":
		if err := buildStandaloneAndroidLibs(root, env, runner, stdout, stderr); err != nil {
			return err
		}
		gradlew := filepath.Join(root, "android", "gradlew")
		if _, err := os.Stat(gradlew); err == nil {
			return runner.Run(gradlew, []string{"-p", "android", "assembleDebug", "--no-daemon"}, root, env, stdout, stderr)
		}
		// Fallback to framework script if available
		if frameworkRoot, err := findFrameworkRoot(); err == nil {
			script := filepath.Join(frameworkRoot, "scripts", "build-android.sh")
			if _, err := os.Stat(script); err == nil {
				return runner.Run(script, nil, frameworkRoot, env, stdout, stderr)
			}
		}
		return fmt.Errorf("android build requires gradlew in %s/android", root)

	case action == "run" && platform == "android":
		if err := runStandalonePlatformCommand(root, "build", "android", runner, stdout, stderr); err != nil {
			return err
		}
		apk := filepath.Join(root, "android", "app", "build", "outputs", "apk", "debug", "app-debug.apk")
		if _, err := os.Stat(apk); err != nil {
			apk = filepath.Join(root, "build", "android", "GoNativeCounter.apk")
		}
		if _, err := os.Stat(apk); err != nil {
			return fmt.Errorf("built APK not found: %w", err)
		}
		if err := runner.Run("adb", []string{"install", "-r", apk}, root, env, stdout, stderr); err != nil {
			return fmt.Errorf("adb install: %w", err)
		}
		activity := fmt.Sprintf("dev.gonative.%s/.MainActivity", pkg)
		if err := runner.Run("adb", []string{"shell", "am", "start", "-n", activity}, root, env, stdout, stderr); err != nil {
			return fmt.Errorf("adb am start: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unsupported platform command %s %s", action, platform)
}

func findNDKToolchain(sdk, preferredVersion, hostTag string) (string, error) {
	if preferredVersion != "" {
		p := filepath.Join(sdk, "ndk", preferredVersion, "toolchains", "llvm", "prebuilt", hostTag)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	ndkParent := filepath.Join(sdk, "ndk")
	if entries, err := os.ReadDir(ndkParent); err == nil {
		for i := len(entries) - 1; i >= 0; i-- {
			if entries[i].IsDir() {
				p := filepath.Join(ndkParent, entries[i].Name(), "toolchains", "llvm", "prebuilt", hostTag)
				if _, err := os.Stat(p); err == nil {
					return p, nil
				}
			}
		}
	}
	bundlePath := filepath.Join(sdk, "ndk-bundle", "toolchains", "llvm", "prebuilt", hostTag)
	if _, err := os.Stat(bundlePath); err == nil {
		return bundlePath, nil
	}
	return "", fmt.Errorf("Android NDK LLVM toolchain not found under %s", sdk)
}

func buildStandaloneAndroidLibs(root string, env []string, runner commandRunner, stdout, stderr io.Writer) error {
	sdk := os.Getenv("ANDROID_SDK_ROOT")
	if sdk == "" {
		sdk = os.Getenv("ANDROID_HOME")
	}
	if sdk == "" {
		home, _ := os.UserHomeDir()
		sdk = filepath.Join(home, "Library", "Android", "sdk")
	}
	ndkVersion := os.Getenv("GONATIVE_NDK_VERSION")
	if ndkVersion == "" {
		ndkVersion = "28.2.13676358"
	}
	hostTag := os.Getenv("GONATIVE_NDK_HOST_TAG")
	if hostTag == "" {
		hostTag = "darwin-x86_64"
	}
	toolchain, err := findNDKToolchain(sdk, ndkVersion, hostTag)
	if err != nil {
		return err
	}
	abis := os.Getenv("GONATIVE_ANDROID_ABIS")
	if abis == "" {
		abis = "arm64-v8a,x86_64"
	}

	buildDir := filepath.Join(root, "build", "android", "lib")
	compiledCount := 0
	for _, abi := range strings.Split(abis, ",") {
		abi = strings.TrimSpace(abi)
		if abi == "" {
			continue
		}
		var goarch, compiler string
		switch abi {
		case "arm64-v8a":
			goarch = "arm64"
			compiler = "aarch64-linux-android23-clang"
		case "x86_64":
			goarch = "amd64"
			compiler = "x86_64-linux-android23-clang"
		default:
			continue
		}
		compilerPath := filepath.Join(toolchain, "bin", compiler)
		if _, err := os.Stat(compilerPath); err != nil {
			continue
		}
		outDir := filepath.Join(buildDir, abi)
		_ = os.MkdirAll(outDir, 0o755)

		cgoEnv := append([]string{}, env...)
		cgoEnv = defaultEnv(cgoEnv, "CGO_ENABLED", "1")
		cgoEnv = defaultEnv(cgoEnv, "GOOS", "android")
		cgoEnv = defaultEnv(cgoEnv, "GOARCH", goarch)
		cgoEnv = defaultEnv(cgoEnv, "CC", compilerPath)
		cgoEnv = defaultEnv(cgoEnv, "CGO_CFLAGS", fmt.Sprintf("--sysroot=%s/sysroot -I%s/sysroot/usr/include", toolchain, toolchain))

		outFile := filepath.Join(outDir, "libgonative.so")
		if err := runner.Run("go", []string{"build", "-buildmode=c-shared", "-o", outFile, "./android/bridge"}, root, cgoEnv, stdout, stderr); err != nil {
			return fmt.Errorf("compile android native lib for %s: %w", abi, err)
		}
		_ = os.Remove(filepath.Join(outDir, "libgonative.h"))
		compiledCount++
	}
	if compiledCount == 0 {
		return fmt.Errorf("no supported ABI compilers found in toolchain %s", toolchain)
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func defaultEnv(env []string, key, value string) []string {
	prefix := key + "="
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			return env
		}
	}
	return append(env, prefix+value)
}

func findFrameworkRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		module, readErr := os.ReadFile(filepath.Join(dir, "go.mod"))
		if readErr == nil && strings.Contains(string(module), "module github.com/go-native/go-native") {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("framework root not found")
		}
		dir = parent
	}
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not inside a go-native project or repository")
		}
		dir = parent
	}
}

type check struct{ label, command, path string }

func doctor(root string, out io.Writer) error {
	fmt.Fprintf(out, "Go Native doctor (%s/%s)\n", runtime.GOOS, runtime.GOARCH)
	checks := []check{
		{label: "Go", command: "go"},
		{label: "Xcode build tools", command: "xcodebuild"},
		{label: "Java compiler", command: "javac"},
	}
	iosScript := filepath.Join(root, "scripts", "build-ios.sh")
	if _, err := os.Stat(iosScript); err == nil {
		checks = append(checks, check{label: "iOS build script", path: iosScript})
	} else if _, err := os.Stat(filepath.Join(root, "ios")); err == nil {
		checks = append(checks, check{label: "iOS project folder", path: filepath.Join(root, "ios")})
	}

	androidScript := filepath.Join(root, "scripts", "build-android.sh")
	if _, err := os.Stat(androidScript); err == nil {
		checks = append(checks, check{label: "Android build script", path: androidScript})
	} else if _, err := os.Stat(filepath.Join(root, "android")); err == nil {
		checks = append(checks, check{label: "Android project folder", path: filepath.Join(root, "android")})
	}

	sdk := os.Getenv("ANDROID_SDK_ROOT")
	if sdk == "" {
		sdk = os.Getenv("ANDROID_HOME")
	}
	if sdk == "" {
		home, _ := os.UserHomeDir()
		sdk = filepath.Join(home, "Library", "Android", "sdk")
	}
	checks = append(checks, check{label: "Android SDK", path: sdk})
	missing := 0
	for _, item := range checks {
		var value string
		var err error
		if item.command != "" {
			value, err = exec.LookPath(item.command)
		} else {
			value = item.path
			_, err = os.Stat(item.path)
		}
		if err != nil {
			missing++
			fmt.Fprintf(out, "  missing  %s\n", item.label)
		} else {
			fmt.Fprintf(out, "  ok       %-20s %s\n", item.label, value)
		}
	}
	if missing > 0 {
		return fmt.Errorf("%d required tool or path missing", missing)
	}
	return nil
}

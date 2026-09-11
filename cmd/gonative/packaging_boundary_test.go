package main

import (
	"os"
	"path/filepath"
	"testing"
)

// These sources have already moved behind the framework packaging boundary.
// During the transition, the generated project still carries the legacy renderer,
// but newly extracted framework modules must not be copied into applications.
func TestExtractedIOSModulesRemainFrameworkOwned(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	modules := []string{
		"GNProtocolReader.h",
		"GNProtocolReader.m",
		"GNViewRegistry.h",
		"GNViewRegistry.m",
		"GNRuntimeHost.m",
		"GNMeasurementHost.h",
		"GNMeasurementHost.m",
		"GNControls.h",
		"GNControls.m",
		"GoNativeRenderer.h",
		"GoNativeRenderer.m",
	}
	templates := getProjectTemplates("ownership-check")

	for _, name := range modules {
		t.Run(name, func(t *testing.T) {
			frameworkPath := filepath.Join(root, "platform", "ios", name)
			if _, err := os.Stat(frameworkPath); err != nil {
				t.Fatalf("framework-owned source missing: %v", err)
			}

			generatedPath := filepath.ToSlash(filepath.Join("ios", name))
			if _, copied := templates[generatedPath]; copied {
				t.Fatalf("framework-owned source %q must not be copied by gonative init", generatedPath)
			}

			fixturePath := filepath.Join(root, "examples", "showcase-app", "ios", name)
			if _, err := os.Stat(fixturePath); err == nil {
				t.Fatalf("framework-owned source was copied into generated fixture: %s", fixturePath)
			} else if !os.IsNotExist(err) {
				t.Fatalf("inspect generated fixture: %v", err)
			}
		})
	}
}

func TestExtractedAndroidModulesRemainFrameworkOwned(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	modules := []string{
		"ProtocolReader.java",
		"ViewRegistry.java",
		"NativeMeasurer.java",
		"ControlFactory.java",
		"GestureBinding.java",
		"EventDispatcher.java",
		"GapDrawable.java",
	}

	templates := getProjectTemplates("ownership-check")

	for _, name := range modules {
		t.Run(name, func(t *testing.T) {
			frameworkPath := filepath.Join(root, "platform", "android", "runtime", "src", "main", "java", "dev", "gonative", "runtime", name)
			if _, err := os.Stat(frameworkPath); err != nil {
				t.Fatalf("framework-owned Android source missing: %v", err)
			}

			for generatedPath := range templates {
				if filepath.Base(generatedPath) == name {
					t.Fatalf("framework-owned Android module %q must not be copied by gonative init: %s", name, generatedPath)
				}
			}

			fixturePath := filepath.Join(root, "examples", "showcase-app", "android", "app", "src", "main", "java", "dev", "gonative", "runtime", name)
			if _, err := os.Stat(fixturePath); err == nil {
				t.Fatalf("framework-owned Android source was copied into generated fixture: %s", fixturePath)
			} else if !os.IsNotExist(err) {
				t.Fatalf("inspect generated fixture: %v", err)
			}
		})
	}
}

package android_test

import (
	"os"
	"strings"
	"testing"
)

func TestRootViewKeepsCommittedBackgroundStyle(t *testing.T) {
	source, err := os.ReadFile("runtime/src/main/java/dev/gonative/runtime/GoNativeActivity.java")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(source), "view.setBackgroundColor(android.graphics.Color.WHITE)") {
		t.Fatal("root mounting must not overwrite the typed background applied by ControlFactory")
	}
}

func TestRootViewRemainsMatchParentAfterUpdates(t *testing.T) {
	source, err := os.ReadFile("runtime/src/main/java/dev/gonative/runtime/GoNativeActivity.java")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	if !strings.Contains(text, "if (nodeID == registry.getRootNodeID()) fillRoot(view);") {
		t.Fatal("root updates must restore MATCH_PARENT after legacy styling")
	}
	if !strings.Contains(text, "ViewGroup.LayoutParams.MATCH_PARENT") {
		t.Fatal("root mounting must fill the Android content viewport")
	}
}

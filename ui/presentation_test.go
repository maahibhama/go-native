package ui

import "testing"

func TestNavigationStackBuildsVisibleRoute(t *testing.T) {
	stack := NavigationStack(
		Route{Key: "home", Title: "Home", Content: Text("Home")},
		Route{Key: "details", Title: "Details", Content: Text("Details")},
	)
	if got := stack.Build().Props.Text; got != "Details" {
		t.Fatalf("visible route = %q, want Details", got)
	}
	intent, ok := NavigationOf(stack)
	if !ok || len(intent.Routes) != 2 || intent.Routes[0].Key != "home" {
		t.Fatalf("unexpected navigation intent: %#v, %v", intent, ok)
	}
}

func TestNavigationIntentRoutesAreCopied(t *testing.T) {
	stack := NavigationStack(Route{Key: "home", Content: Text("Home")})
	intent, _ := NavigationOf(stack)
	intent.Routes[0].Key = "changed"
	second, _ := NavigationOf(stack)
	if second.Routes[0].Key != "home" {
		t.Fatal("navigation metadata exposed its route slice")
	}
}

func TestNavigationStackSkipsNilContentFallback(t *testing.T) {
	stack := NavigationStack(Route{Key: "home", Content: Text("Home")}, Route{Key: "pending"})
	if got := stack.Build().Props.Text; got != "Home" {
		t.Fatalf("fallback route = %q, want Home", got)
	}
	if node := NavigationStack().Build(); node != nil {
		t.Fatalf("empty navigation stack built %#v", node)
	}
}

func TestModalPreservesBaseAndMetadata(t *testing.T) {
	dismissed := false
	component := PresentModal(Text("Base"), ModalIntent{
		Key: "settings", Content: Text("Settings"), Style: ModalSheet,
		Dismissible: true, OnDismiss: func() { dismissed = true },
	})
	if got := component.Build().Props.Text; got != "Base" {
		t.Fatalf("base = %q, want Base", got)
	}
	intent, presented := ModalOf(component)
	if !presented || intent.Key != "settings" || intent.Style != ModalSheet {
		t.Fatalf("unexpected modal intent: %#v, %v", intent, presented)
	}
	intent.OnDismiss()
	if !dismissed {
		t.Fatal("dismiss callback was not retained")
	}
}

func TestModalWithoutContentIsNotPresented(t *testing.T) {
	component := PresentModal(nil, ModalIntent{Key: "absent"})
	if component.Build() != nil {
		t.Fatal("nil base built a node")
	}
	if _, presented := ModalOf(component); presented {
		t.Fatal("modal without content reported as presented")
	}
	if _, ok := NavigationOf(Text("plain")); ok {
		t.Fatal("plain component reported navigation metadata")
	}
}

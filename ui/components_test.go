package ui

import "testing"

func TestTextFuncEvaluatesOnBuildDescription(t *testing.T) {
	value := "first"
	build := func() *Node { return TextFunc(func() string { return value }).Build() }
	if got := build().Props.Text; got != "first" {
		t.Fatalf("text=%q", got)
	}
	value = "second"
	if got := build().Props.Text; got != "second" {
		t.Fatalf("text=%q", got)
	}
}

func TestSafeAreaNode(t *testing.T) {
	n := SafeArea(Text("inside")).Build()
	if n.Type != NodeSafeArea || len(n.Children) != 1 || n.Children[0].Type != NodeText {
		t.Fatalf("unexpected safe area: %#v", n)
	}
}

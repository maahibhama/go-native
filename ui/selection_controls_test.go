package ui

import (
	"testing"
)

func TestCheckbox(t *testing.T) {
	changed := false
	var gotVal bool
	onChange := func(val bool) {
		changed = true
		gotVal = val
	}

	cb := Checkbox(false, onChange)
	node := cb.Build()

	if node.Type != NodeSwitch {
		t.Fatalf("expected NodeSwitch, got %v", node.Type)
	}
	if node.Props.AccessRole != RoleButton {
		t.Fatalf("expected RoleButton, got %v", node.Props.AccessRole)
	}
	if node.Props.Checked {
		t.Fatal("expected Checked to be false")
	}
	if node.Toggle == nil {
		t.Fatal("expected non-nil node.Toggle")
	}

	// Verify touch target constraints
	if node.Style.Layout.MinHeight.Value != 44 || node.Platform.Android.Layout.MinHeight.Value != 48 {
		t.Fatalf("expected accessible touch targets on Checkbox")
	}

	// Verify toggle handler invocation
	node.Toggle(true)
	if !changed || !gotVal {
		t.Fatalf("node.Toggle did not set gotVal=true: changed=%v, gotVal=%v", changed, gotVal)
	}

	// Verify Press toggles checked state
	changed = false
	node.Press()
	if !changed || !gotVal {
		t.Fatalf("node.Press did not toggle false to true: changed=%v, gotVal=%v", changed, gotVal)
	}

	// Test checked=true toggling to false
	changed = false
	cbChecked := Checkbox(true, onChange).Build()
	cbChecked.Press()
	if !changed || gotVal {
		t.Fatalf("node.Press did not toggle true to false: changed=%v, gotVal=%v", changed, gotVal)
	}

	// Test disabled Checkbox
	disabledChanged := false
	disabledCb := Checkbox(false, func(bool) { disabledChanged = true }).Disabled(true).Build()
	disabledCb.Press()
	if disabledChanged {
		t.Fatal("disabled Checkbox should not invoke onChange on Press")
	}
}

func TestRadio(t *testing.T) {
	selected := false
	onSelect := func() { selected = true }

	// Unselected radio
	radio := Radio(false, onSelect)
	node := radio.Build()

	if node.Type != NodeButton {
		t.Fatalf("expected NodeButton, got %v", node.Type)
	}
	if node.Props.AccessRole != RoleButton {
		t.Fatalf("expected RoleButton, got %v", node.Props.AccessRole)
	}
	if node.Props.Checked {
		t.Fatal("expected Checked to be false")
	}
	if node.Press == nil {
		t.Fatal("expected non-nil node.Press")
	}

	node.Press()
	if !selected {
		t.Fatal("pressing unselected radio should invoke onSelect")
	}

	// Already selected radio should not fire onSelect again
	selectedAgain := false
	selectedRadio := Radio(true, func() { selectedAgain = true }).Build()
	selectedRadio.Press()
	if selectedAgain {
		t.Fatal("pressing already selected radio should not invoke onSelect")
	}

	// Test Toggle handler
	toggleFired := false
	toggleRadio := Radio(false, func() { toggleFired = true }).Build()
	if toggleRadio.Toggle != nil {
		toggleRadio.Toggle(true)
		if !toggleFired {
			t.Fatal("Toggle(true) should invoke onSelect")
		}
		toggleFired = false
		toggleRadio.Toggle(false)
		if toggleFired {
			t.Fatal("Toggle(false) should not invoke onSelect")
		}
	}
}

func TestRadioGroup(t *testing.T) {
	var selectedFruit string
	options := []RadioOption[string]{
		{Value: "apple", Label: "Apple"},
		{Value: "banana", Label: "Banana"},
		{Value: "cherry", Label: "Cherry", Disabled: true},
	}

	group := RadioGroup("banana", func(val string) {
		selectedFruit = val
	}, options)

	node := group.Build()
	if node.Type != NodeColumn {
		t.Fatalf("expected NodeColumn for RadioGroup, got %v", node.Type)
	}
	if len(node.Children) != 3 {
		t.Fatalf("expected 3 options, got %d", len(node.Children))
	}

	// Check option 0: "apple" (unselected)
	opt0Row := node.Children[0]
	if len(opt0Row.Children) < 2 {
		t.Fatalf("expected radio + label row, got %#v", opt0Row.Children)
	}
	radio0 := opt0Row.Children[0]
	if radio0.Props.Checked {
		t.Fatal("apple should not be checked")
	}

	// Check option 1: "banana" (selected)
	opt1Row := node.Children[1]
	radio1 := opt1Row.Children[0]
	if !radio1.Props.Checked {
		t.Fatal("banana should be checked")
	}

	// Check option 2: "cherry" (disabled)
	opt2Row := node.Children[2]
	radio2 := opt2Row.Children[0]
	if !radio2.Style.Interaction.Disabled {
		t.Fatal("cherry radio should be disabled")
	}

	// Selecting "apple" via radio0 Press
	radio0.Press()
	if selectedFruit != "apple" {
		t.Fatalf("expected selectedFruit to be 'apple', got %q", selectedFruit)
	}

	// Selecting disabled option should not change selectedFruit
	radio2.Press()
	if selectedFruit != "apple" {
		t.Fatal("clicking disabled option should not invoke onSelect")
	}
}

func TestSlider(t *testing.T) {
	var sliderVal float32
	onChange := func(val float32) {
		sliderVal = val
	}

	slider := Slider(0.35, onChange).
		SliderMin(0.0).
		SliderMax(100.0).
		SliderStep(5.0)

	node := slider.Build()
	if node.Type != NodeProgressIndicator {
		t.Fatalf("expected NodeProgressIndicator, got %v", node.Type)
	}
	if node.Props.Progress != 0.35 {
		t.Fatalf("expected Progress 0.35, got %v", node.Props.Progress)
	}
	if node.SliderMin != 0.0 || node.SliderMax != 100.0 || node.SliderStep != 5.0 {
		t.Fatalf("slider min/max/step mismatch: got %v/%v/%v", node.SliderMin, node.SliderMax, node.SliderStep)
	}
	if node.SliderChange == nil {
		t.Fatal("expected non-nil node.SliderChange")
	}

	node.SliderChange(50.0)
	if sliderVal != 50.0 {
		t.Fatalf("expected sliderVal 50.0, got %v", sliderVal)
	}

	// Test string change conversion
	if node.Change != nil {
		node.Change("75.5")
		if sliderVal != 75.5 {
			t.Fatalf("expected sliderVal 75.5 from Change, got %v", sliderVal)
		}
	}
}

func TestSegmentedControl(t *testing.T) {
	var selectedTab string
	items := []SegmentItem[string]{
		{Value: "day", Label: "Day"},
		{Value: "week", Label: "Week"},
		{Value: "month", Label: "Month", Disabled: true},
	}

	seg := SegmentedControl("week", func(val string) {
		selectedTab = val
	}, items)

	node := seg.Build()
	if node.Type != NodeRow {
		t.Fatalf("expected NodeRow container, got %v", node.Type)
	}
	if len(node.Children) != 3 {
		t.Fatalf("expected 3 segment buttons, got %d", len(node.Children))
	}

	// Selected segment ("week") should be bold with non-zero background
	selectedBtn := node.Children[1]
	if !selectedBtn.Props.Bold {
		t.Fatal("selected segment should be bold")
	}
	if selectedBtn.Style.Appearance.Background.IsTransparent() {
		t.Fatal("selected segment should have solid background")
	}

	// Unselected segment ("day") should have transparent background
	unselectedBtn := node.Children[0]
	if !unselectedBtn.Style.Appearance.Background.IsTransparent() {
		t.Fatal("unselected segment should have transparent background")
	}

	// Disabled segment ("month")
	disabledBtn := node.Children[2]
	if !disabledBtn.Style.Interaction.Disabled {
		t.Fatal("expected month segment to be disabled")
	}

	// Clicking unselected segment triggers onSelect
	unselectedBtn.Press()
	if selectedTab != "day" {
		t.Fatalf("expected selectedTab to be 'day', got %q", selectedTab)
	}

	// Clicking disabled segment does not trigger onSelect
	disabledBtn.Press()
	if selectedTab != "day" {
		t.Fatal("clicking disabled segment should not invoke onSelect")
	}
}

func TestProgressIndicators(t *testing.T) {
	// Linear determinate
	linear := LinearProgress(0.7).Build()
	if linear.Type != NodeProgressIndicator {
		t.Fatalf("expected NodeProgressIndicator, got %v", linear.Type)
	}
	if linear.Props.Progress != 0.7 {
		t.Fatalf("expected Progress 0.7, got %v", linear.Props.Progress)
	}
	if linear.Indeterminate() {
		t.Fatal("0.7 should not be indeterminate")
	}
	if linear.IsCircular {
		t.Fatal("LinearProgress should not be circular")
	}

	// Linear indeterminate (negative value)
	linearIndet := LinearProgress(-1.0).Build()
	if !linearIndet.Indeterminate() {
		t.Fatal("LinearProgress(-1.0) should be indeterminate")
	}
	if linearIndet.Props.Progress >= 0 {
		t.Fatalf("expected negative progress for indeterminate, got %v", linearIndet.Props.Progress)
	}

	// Clamped values for LinearProgress
	linearOver := LinearProgress(1.5).Build()
	if linearOver.Props.Progress != 1.0 {
		t.Fatalf("expected Progress clamped to 1.0, got %v", linearOver.Props.Progress)
	}

	// Circular determinate
	circ := CircularProgress(0.4).Build()
	if circ.Type != NodeProgressIndicator {
		t.Fatalf("expected NodeProgressIndicator, got %v", circ.Type)
	}
	if circ.Props.Progress != 0.4 {
		t.Fatalf("expected Progress 0.4, got %v", circ.Props.Progress)
	}
	if circ.Indeterminate() {
		t.Fatal("0.4 should not be indeterminate")
	}
	if !circ.IsCircular {
		t.Fatal("CircularProgress should have IsCircular=true")
	}

	// Circular indeterminate (negative value)
	circIndet := CircularProgress(-0.5).Build()
	if !circIndet.Indeterminate() {
		t.Fatal("CircularProgress(-0.5) should be indeterminate")
	}
	if !circIndet.IsCircular {
		t.Fatal("CircularProgress should have IsCircular=true")
	}
}

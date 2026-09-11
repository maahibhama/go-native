package ui

import (
	"testing"
	"time"
)

func TestScaffold(t *testing.T) {
	appBar := AppBar("Test Title")
	body := Text("Body Content")
	bottomBar := View(Text("Bottom"))
	fabClicked := false
	fab := Button("+", func() { fabClicked = true })
	bg := RGB(240, 240, 240)

	scaffold := Scaffold(ScaffoldProps{
		AppBar:               appBar,
		Body:                 body,
		BottomBar:            bottomBar,
		FloatingActionButton: fab,
		BackgroundColor:      bg,
	})

	// Verify ScaffoldProps accessor
	props := scaffold.(*scaffoldComponent).Props()
	if props.AppBar != appBar || props.Body != body || props.BottomBar != bottomBar {
		t.Fatal("unexpected ScaffoldProps values")
	}

	node := scaffold.Build()
	if node == nil {
		t.Fatal("expected non-nil node from Scaffold.Build()")
	}

	if node.Style.Appearance.Background != bg {
		t.Fatalf("expected background %v, got %v", bg, node.Style.Appearance.Background)
	}

	// Since FAB is provided, root is a Stack containing mainColumn and FAB overlay
	if len(node.Children) != 2 {
		t.Fatalf("expected 2 children in Scaffold root (Stack), got %d", len(node.Children))
	}

	mainCol := node.Children[0]
	if len(mainCol.Children) != 3 {
		t.Fatalf("expected 3 children in main column (AppBar, Body, BottomBar), got %d", len(mainCol.Children))
	}

	// FAB overlay is positioned absolute
	fabOverlay := node.Children[1]
	if fabOverlay.Style.Layout.Position != PositionAbsolute {
		t.Fatalf("expected FAB overlay to be PositionAbsolute, got %v", fabOverlay.Style.Layout.Position)
	}

	// Verify fab press callback
	fabNode := fabOverlay.Children[0]
	if fabNode.Press != nil {
		fabNode.Press()
		if !fabClicked {
			t.Fatal("expected FAB click handler to be invoked")
		}
	}

	// Test Scaffold with nil FAB
	scaffoldNoFAB := Scaffold(ScaffoldProps{
		Body: Text("Simple"),
	}).Build()
	if scaffoldNoFAB == nil || len(scaffoldNoFAB.Children) == 0 {
		t.Fatal("expected valid tree for Scaffold without FAB")
	}
}

func TestAppBar(t *testing.T) {
	actionCalled := false
	actionBtn := Button("Save", func() { actionCalled = true })
	leadingIcon := Text("<")

	bar := AppBarWithProps(AppBarProps{
		Title:     "Settings",
		Leading:   leadingIcon,
		Actions:   []Component{actionBtn},
		Elevation: 4.0,
	})

	appBarComp := bar.(*appBarComponent)
	if appBarComp.Props().Title != "Settings" {
		t.Fatalf("expected title 'Settings', got %q", appBarComp.Props().Title)
	}

	node := bar.Build()
	if node == nil {
		t.Fatal("expected non-nil node for AppBar")
	}

	if node.Type != NodeRow {
		t.Fatalf("expected AppBar root to be NodeRow, got %v", node.Type)
	}

	if node.Style.Appearance.Shadow.Blur != 4.0 {
		t.Fatalf("expected shadow blur 4.0, got %v", node.Style.Appearance.Shadow.Blur)
	}

	if actionBtn.node.Press != nil {
		actionBtn.node.Press()
		if !actionCalled {
			t.Fatal("expected action button click handler to be invoked")
		}
	}

	// Short helper AppBar(title, actions...)
	shortBar := AppBar("Short Title", actionBtn).Build()
	if shortBar == nil {
		t.Fatal("expected non-nil node from AppBar short constructor")
	}
}

func TestForm(t *testing.T) {
	submitted := false
	form := Form(
		TextInput("user", nil),
		TextInput("pass", nil),
	).OnSubmit(func() {
		submitted = true
	})

	if form.IsValid() != true {
		t.Fatal("expected empty form to be valid")
	}

	form.SetError("pass", "Password too short")
	if form.IsValid() {
		t.Fatal("expected form with error to be invalid")
	}
	if form.Error("pass") != "Password too short" {
		t.Fatalf("expected error 'Password too short', got %q", form.Error("pass"))
	}

	errors := form.Errors()
	if len(errors) != 1 || errors["pass"] != "Password too short" {
		t.Fatalf("unexpected errors map: %#v", errors)
	}

	form.ClearError("pass")
	if !form.IsValid() {
		t.Fatal("expected form to be valid after ClearError")
	}

	form.SetError("field1", "err1")
	form.SetError("field2", "err2")
	form.ClearErrors()
	if !form.IsValid() || len(form.Errors()) != 0 {
		t.Fatal("expected form to have no errors after ClearErrors")
	}

	form.Submit()
	if !submitted {
		t.Fatal("expected submit handler to be called by form.Submit()")
	}

	// Verify node has Submit callback attached
	node := form.Build()
	if node == nil || node.Submit == nil {
		t.Fatal("expected built node to have Submit callback")
	}

	// Verify FormValidationKey injection into context dependencies
	env := DefaultEnvironment()
	ctx := NewBuildContext(env)
	_ = form.BuildContext(ctx)
	resolved, ok := Resolve(env.Dependencies, FormValidationKey)
	if !ok || resolved == nil {
		t.Fatal("expected FormValidationState to be resolvable from dependencies")
	}

	// Test FormWithProps
	customValidation := NewFormValidationState()
	form2 := FormWithProps(FormProps{
		OnSubmit:   func() {},
		Validation: customValidation,
	}, Text("inside"))
	if form2.ValidationState() != customValidation {
		t.Fatal("expected FormWithProps to retain custom ValidationState")
	}
}

func TestCard(t *testing.T) {
	c := Card(Text("Hello Card"))
	node := c.Build()
	if node == nil {
		t.Fatal("expected non-nil node for Card")
	}

	// Default is CardElevated
	if node.Style.Appearance.Shadow.Blur == 0 {
		t.Fatal("expected default Card to have elevation shadow")
	}

	outlined := Card(Text("Outlined")).CardOutlined().Build()
	if outlined.Style.Appearance.Border.Width != 1 {
		t.Fatalf("expected border width 1 for CardOutlined, got %v", outlined.Style.Appearance.Border.Width)
	}
	if outlined.Style.Appearance.Shadow.Blur != 0 {
		t.Fatal("expected no shadow for CardOutlined")
	}

	filled := Card(Text("Filled")).CardFilled().Build()
	if filled.Style.Appearance.Border.Width != 0 {
		t.Fatal("expected no border for CardFilled")
	}
	if filled.Style.Appearance.Shadow.Blur != 0 {
		t.Fatal("expected no shadow for CardFilled")
	}
	if filled.Style.Appearance.Background != RGB(245, 245, 245) {
		t.Fatalf("expected filled background RGB(245, 245, 245), got %v", filled.Style.Appearance.Background)
	}
}

func TestBadgeAndDotBadge(t *testing.T) {
	// Badge with positive count
	content := Text("Inbox")
	badge := Badge(content, 5)
	bComp := badge.(*badgeComponent)
	if bComp.Count() != 5 || bComp.Content() != content {
		t.Fatal("unexpected badgeComponent properties")
	}

	node := badge.Build()
	if node == nil || len(node.Children) != 2 {
		t.Fatalf("expected Stack with 2 children for Badge, got %#v", node)
	}

	// Badge with count > 99
	badge99 := Badge(content, 150).Build()
	if badge99 == nil {
		t.Fatal("expected non-nil node for Badge > 99")
	}

	// Badge with count <= 0 renders only content
	badgeZero := Badge(content, 0).Build()
	if badgeZero == nil || badgeZero.Type != NodeText || badgeZero.Props.Text != "Inbox" {
		t.Fatalf("expected only content rendered for Badge count 0, got %#v", badgeZero)
	}

	// DotBadge visible
	dotBadge := DotBadge(content, true)
	dComp := dotBadge.(*dotBadgeComponent)
	if !dComp.Visible() || dComp.Content() != content {
		t.Fatal("unexpected dotBadgeComponent properties")
	}

	dotNode := dotBadge.Build()
	if dotNode == nil || len(dotNode.Children) != 2 {
		t.Fatalf("expected Stack with 2 children for visible DotBadge, got %#v", dotNode)
	}

	// DotBadge invisible renders only content
	dotBadgeHidden := DotBadge(content, false).Build()
	if dotBadgeHidden == nil || dotBadgeHidden.Type != NodeText {
		t.Fatalf("expected only content rendered for hidden DotBadge, got %#v", dotBadgeHidden)
	}

	// DotBadge nil content
	dotOnly := DotBadge(nil, true).Build()
	if dotOnly == nil {
		t.Fatal("expected non-nil node for DotBadge with nil content")
	}
}

func TestAvatar(t *testing.T) {
	// Image avatar
	imgAvatar := Avatar("https://example.com/avatar.png", "JD", 50)
	aComp := imgAvatar.(*avatarComponent)
	if aComp.Source() != "https://example.com/avatar.png" || aComp.Initials() != "JD" || aComp.Size() != 50 {
		t.Fatal("unexpected avatarComponent properties")
	}

	imgNode := imgAvatar.Build()
	if imgNode == nil || imgNode.Type != NodeImage {
		t.Fatalf("expected NodeImage for Avatar with source, got %#v", imgNode)
	}
	if imgNode.Style.Layout.Width != Points(50) || imgNode.Style.Appearance.CornerRadius != 25 {
		t.Fatalf("expected size 50 and corner radius 25, got width=%v, radius=%v", imgNode.Style.Layout.Width, imgNode.Style.Appearance.CornerRadius)
	}

	// Fallback avatar with initials
	initAvatar := Avatar("", "AB", 40)
	initNode := initAvatar.Build()
	if initNode == nil || initNode.Type != NodeView {
		t.Fatalf("expected NodeView container for initials avatar fallback, got %#v", initNode)
	}
	if initNode.Style.Layout.Width != Points(40) || initNode.Style.Appearance.CornerRadius != 20 {
		t.Fatalf("expected size 40 and corner radius 20, got width=%v, radius=%v", initNode.Style.Layout.Width, initNode.Style.Appearance.CornerRadius)
	}

	// Default size when size <= 0
	defAvatar := Avatar("", "CD", 0).Build()
	if defAvatar.Style.Layout.Width != Points(40) {
		t.Fatalf("expected default size 40, got %v", defAvatar.Style.Layout.Width)
	}
}

func TestSkeleton(t *testing.T) {
	skel := Skeleton(120, 24)
	node := skel.Build()
	if node == nil {
		t.Fatal("expected non-nil node for Skeleton")
	}
	if node.Style.Layout.Width != Points(120) || node.Style.Layout.Height != Points(24) {
		t.Fatalf("expected 120x24 size, got %v x %v", node.Style.Layout.Width, node.Style.Layout.Height)
	}
	if node.Style.Appearance.CornerRadius != 4 {
		t.Fatalf("expected corner radius 4, got %v", node.Style.Appearance.CornerRadius)
	}
	if len(node.Intents.Animations) != 1 {
		t.Fatalf("expected 1 animation intent, got %d", len(node.Intents.Animations))
	}
	anim := node.Intents.Animations[0]
	if anim.Property != AnimateOpacity || anim.From != 0.4 || anim.To != 1.0 {
		t.Fatalf("expected opacity pulse animation from 0.4 to 1.0, got %#v", anim)
	}
	if anim.Duration != 1000*time.Millisecond {
		t.Fatalf("expected 1s duration, got %v", anim.Duration)
	}
}

func TestEmptyStateAndErrorState(t *testing.T) {
	actionCalled := false
	actionBtn := Button("Add", func() { actionCalled = true })
	empty := EmptyState("No Items", "Create your first item to get started", Text("icon"), actionBtn)
	eComp := empty.(*emptyStateComponent)
	if eComp.Title() != "No Items" || eComp.Description() != "Create your first item to get started" {
		t.Fatal("unexpected EmptyState properties")
	}

	emptyNode := empty.Build()
	if emptyNode == nil {
		t.Fatal("expected non-nil node for EmptyState")
	}
	if actionBtn.node.Press != nil {
		actionBtn.node.Press()
		if !actionCalled {
			t.Fatal("expected action callback to be called")
		}
	}

	retryCalled := false
	errState := ErrorState("Failed to load", "Network timeout occurred", func() { retryCalled = true })
	errComp := errState.(*errorStateComponent)
	if errComp.Title() != "Failed to load" || errComp.Message() != "Network timeout occurred" {
		t.Fatal("unexpected ErrorState properties")
	}
	if errComp.OnRetry() != nil {
		errComp.OnRetry()()
		if !retryCalled {
			t.Fatal("expected retry callback to be called")
		}
	}

	errNode := errState.Build()
	if errNode == nil {
		t.Fatal("expected non-nil node for ErrorState")
	}

	// ErrorState without retry
	errNoRetry := ErrorState("Notice", "Simple notice", nil).Build()
	if errNoRetry == nil {
		t.Fatal("expected non-nil node for ErrorState without retry")
	}
}

package ui

import "testing"

func TestDialog(t *testing.T) {
	confirmed := false
	cancelled := false
	dismissed := false

	dialog := Dialog(DialogProps{
		Title:        "Delete Item",
		Message:      "Are you sure you want to permanently delete this item?",
		Content:      Text("Warning details"),
		ConfirmLabel: "Delete",
		OnConfirm:    func() { confirmed = true },
		CancelLabel:  "Cancel",
		OnCancel:     func() { cancelled = true },
		OnDismiss:    func() { dismissed = true },
	})

	dComp := dialog.(*dialogComponent)
	props := dComp.Props()
	if props.Title != "Delete Item" || props.ConfirmLabel != "Delete" || props.CancelLabel != "Cancel" {
		t.Fatal("unexpected DialogProps")
	}

	if props.OnConfirm != nil {
		props.OnConfirm()
		if !confirmed {
			t.Fatal("expected confirmed to be true")
		}
	}
	if props.OnCancel != nil {
		props.OnCancel()
		if !cancelled {
			t.Fatal("expected cancelled to be true")
		}
	}

	node := dialog.Build()
	if node == nil {
		t.Fatal("expected non-nil node for Dialog.Build()")
	}

	// Verify Dialog satisfies ModalComponent
	intent, ok := ModalOf(dialog)
	if !ok {
		t.Fatal("expected Dialog to provide ModalIntent")
	}
	if !intent.Dismissible {
		t.Fatal("expected Dialog with OnDismiss to be Dismissible")
	}
	if intent.OnDismiss != nil {
		intent.OnDismiss()
		if !dismissed {
			t.Fatal("expected dismissed callback to be invoked")
		}
	}

	// Verify overlay backdrop
	if node.Style.Layout.Position != PositionAbsolute {
		t.Fatalf("expected dialog overlay to be PositionAbsolute, got %v", node.Style.Layout.Position)
	}

	// Dialog without dismiss
	noDismissDialog := Dialog(DialogProps{Title: "Locked"}).Build()
	if noDismissDialog == nil {
		t.Fatal("expected non-nil node for simple Dialog")
	}
}

func TestSheet(t *testing.T) {
	dismissed := false
	sheet := Sheet(SheetProps{
		Title:     "Options",
		Content:   Text("Sheet Content"),
		OnDismiss: func() { dismissed = true },
	})

	sComp := sheet.(*sheetComponent)
	if sComp.Props().Title != "Options" {
		t.Fatalf("expected title 'Options', got %q", sComp.Props().Title)
	}

	node := sheet.Build()
	if node == nil {
		t.Fatal("expected non-nil node for Sheet.Build()")
	}

	// Verify Sheet satisfies ModalComponent with ModalSheet style
	intent, ok := ModalOf(sheet)
	if !ok {
		t.Fatal("expected Sheet to provide ModalIntent")
	}
	if intent.Style != ModalSheet {
		t.Fatalf("expected ModalSheet style, got %v", intent.Style)
	}
	if !intent.Dismissible {
		t.Fatal("expected Sheet to be Dismissible")
	}
	if intent.OnDismiss != nil {
		intent.OnDismiss()
		if !dismissed {
			t.Fatal("expected sheet OnDismiss to be invoked")
		}
	}
}

func TestModal(t *testing.T) {
	dismissed := false
	content := Text("Modal Body")

	// Visible modal
	visibleModal := Modal(true, func() { dismissed = true }, content)
	mComp := visibleModal.(*modalOverlayComponent)
	if !mComp.Visible() || mComp.Content() != content {
		t.Fatal("unexpected modalOverlayComponent properties")
	}

	node := visibleModal.Build()
	if node == nil {
		t.Fatal("expected non-nil node for visible Modal")
	}

	intent, ok := ModalOf(visibleModal)
	if !ok {
		t.Fatal("expected visible Modal to have ModalIntent")
	}
	if !intent.Dismissible {
		t.Fatal("expected Modal to be Dismissible")
	}
	if intent.OnDismiss != nil {
		intent.OnDismiss()
		if !dismissed {
			t.Fatal("expected Modal onDismiss to be called")
		}
	}

	// Hidden modal
	hiddenModal := Modal(false, func() {}, content)
	hiddenNode := hiddenModal.Build()
	if hiddenNode != nil {
		t.Fatalf("expected nil node for hidden Modal, got %#v", hiddenNode)
	}

	_, hiddenOk := ModalOf(hiddenModal)
	if hiddenOk {
		t.Fatal("expected hidden Modal to return false for ModalOf")
	}
}

func TestSnackbar(t *testing.T) {
	actionCalled := false
	snack := Snackbar("Item deleted", "Undo", func() {
		actionCalled = true
	})

	sComp := snack.(*snackbarComponent)
	if sComp.Message() != "Item deleted" || sComp.ActionLabel() != "Undo" {
		t.Fatal("unexpected snackbarComponent properties")
	}
	if sComp.OnAction() != nil {
		sComp.OnAction()()
		if !actionCalled {
			t.Fatal("expected snackbar action callback to be called")
		}
	}

	node := snack.Build()
	if node == nil {
		t.Fatal("expected non-nil node for Snackbar")
	}

	if node.Type != NodeRow {
		t.Fatalf("expected Snackbar to be NodeRow, got %v", node.Type)
	}

	// Should contain message text and action button
	if len(node.Children) < 2 {
		t.Fatalf("expected at least 2 children (text, button), got %d", len(node.Children))
	}

	// Snackbar without action
	simpleSnack := Snackbar("Saved successfully", "", nil).Build()
	if simpleSnack == nil {
		t.Fatal("expected non-nil node for simple Snackbar")
	}
}

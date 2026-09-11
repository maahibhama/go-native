package ui

import (
	"strings"
	"testing"
)

func TestTextPresentationProps(t *testing.T) {
	node := Text("A long title").TextWrapping(TextNoWrap).TextOverflow(TextOverflowEllipsisTail).LineLimit(2).Selectable(true).Build()
	if node.Props.TextWrap != TextNoWrap || node.Props.TextOverflow != TextOverflowEllipsisTail || node.Props.MaxLines != 2 || !node.Props.Selectable {
		t.Fatalf("text props = %#v", node.Props)
	}
}

func TestRichTextProducesPortableSpanPayload(t *testing.T) {
	node := RichText([]Span{{Text: "Read ", FontWeight: 600}, {Text: "docs", Link: "docs://start", Underline: true, Color: RGB(10, 20, 30), HasColor: true}}, func(string) {}).Build()
	if node.Props.Text != "Read docs" || len(node.Props.RichText) == 0 || node.Link == nil {
		t.Fatalf("rich text node = %#v", node)
	}
	decoded, err := DecodeRichTextPayload(node.Props.RichText)
	if err != nil || len(decoded) != 2 || decoded[1].Link != "docs://start" || decoded[1].Color != RGB(10, 20, 30) {
		t.Fatalf("decoded spans = %#v, %v", decoded, err)
	}
	if _, err := DecodeRichTextPayload(node.Props.RichText + "x"); err == nil {
		t.Fatal("trailing rich text data accepted")
	}
	if got := RichText([]Span{{Text: strings.Repeat("x", maxRichTextBytes+1)}}, nil).Build().Props.RichText; got != "" {
		t.Fatal("oversized rich text payload was retained")
	}
}

func TestTextInputConfigurationAndDefaults(t *testing.T) {
	selection := TextSelection{Start: 1, End: 3}
	node := TextInput("hello", func(string) {}).Placeholder("Email").InputType(InputEmail).
		ReturnAction(ReturnKeyNext, func() {}).Capitalization(CapitalizeWords).
		AutoCorrect(AutoCorrectDisabled).SecureEntry(true).MultilineInput(true).
		ReadOnly(true).Validation(ValidationInvalid, "Required").SelectionRange(selection, func(TextSelection) {}).
		MaximumLength(12).Build()
	if node.Props.Placeholder != "Email" || node.Props.InputKind != InputEmail || node.Props.ReturnKey != ReturnKeyNext || node.Props.Capitalization != CapitalizeWords || node.Props.AutoCorrect != AutoCorrectDisabled {
		t.Fatalf("input enums = %#v", node.Props)
	}
	if !node.Props.Secure || !node.Props.Multiline || !node.Props.ReadOnly || node.Props.Validation != ValidationInvalid || node.Props.ErrorText != "Required" || node.Props.SelectionStart != 1 || node.Props.SelectionEnd != 3 || node.Props.MaxLength != 12 {
		t.Fatalf("input props = %#v", node.Props)
	}
	defaults := TextInput("", nil).Build().Props
	if defaults.SelectionStart != -1 || defaults.SelectionEnd != -1 || defaults.InputMode != InputControlled {
		t.Fatalf("controlled defaults = %#v", defaults)
	}
	if UncontrolledTextInput("draft", nil).Build().Props.InputMode != InputUncontrolled {
		t.Fatal("uncontrolled input mode not retained")
	}
}

func TestInputFormattersComposeInOrder(t *testing.T) {
	trim := InputFormatterFunc(func(_, proposed string) string { return strings.TrimSpace(proposed) })
	upper := InputFormatterFunc(func(_, proposed string) string { return strings.ToUpper(proposed) })
	if got := ApplyInputFormatters("", "  hello ", []InputFormatter{trim, upper}); got != "HELLO" {
		t.Fatalf("formatted = %q", got)
	}
}

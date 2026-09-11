package ui

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// TextWrapMode controls whether a native text view may wrap between lines.
type TextWrapMode uint8

const (
	TextWrapDefault TextWrapMode = iota
	TextWrap
	TextNoWrap
)

// TextOverflow controls clipping and ellipsis placement when text is constrained.
type TextOverflow uint8

const (
	TextOverflowClip TextOverflow = iota
	TextOverflowEllipsisHead
	TextOverflowEllipsisMiddle
	TextOverflowEllipsisTail
)

// Span is an immutable description of one rich-text run. Link is delivered to
// the RichText onLink callback and is never interpreted by the runtime.
type Span struct {
	Text       string
	Link       string
	FontSize   float32
	FontWeight uint16
	Underline  bool
	Italic     bool
	Color      Color
	HasColor   bool
}

func PlainSpan(text string) Span      { return Span{Text: text} }
func LinkSpan(text, link string) Span { return Span{Text: text, Link: link, Underline: true} }

// RichText creates a native attributed-text view from bounded, value-only spans.
func RichText(spans []Span, onLink func(string)) *element {
	var plain bytes.Buffer
	for _, span := range spans {
		plain.WriteString(span.Text)
	}
	e := newElement(NodeText, Props{Text: plain.String(), RichText: marshalRichText(spans)})
	e.node.Link = onLink
	return e
}

const maxRichTextBytes = 1 << 20
const MaxRichTextSpans = 16_384

func marshalRichText(spans []Span) string {
	if len(spans) == 0 {
		return ""
	}
	if len(spans) > MaxRichTextSpans {
		return ""
	}
	var out bytes.Buffer
	_ = binary.Write(&out, binary.LittleEndian, uint16(1))
	_ = binary.Write(&out, binary.LittleEndian, uint32(len(spans)))
	for _, span := range spans {
		if !utf8.ValidString(span.Text) || !utf8.ValidString(span.Link) || len(span.Text) > maxRichTextBytes || len(span.Link) > maxRichTextBytes {
			return ""
		}
		_ = binary.Write(&out, binary.LittleEndian, uint32(len(span.Text)))
		out.WriteString(span.Text)
		_ = binary.Write(&out, binary.LittleEndian, uint32(len(span.Link)))
		out.WriteString(span.Link)
		_ = binary.Write(&out, binary.LittleEndian, span.FontSize)
		_ = binary.Write(&out, binary.LittleEndian, span.FontWeight)
		var flags byte
		if span.Underline {
			flags |= 1
		}
		if span.Italic {
			flags |= 2
		}
		out.WriteByte(flags)
		if span.HasColor {
			out.WriteByte(1)
		} else {
			out.WriteByte(0)
		}
		out.Write([]byte{span.Color.R, span.Color.G, span.Color.B, span.Color.A})
		if out.Len() > maxRichTextBytes {
			return ""
		}
	}
	return out.String()
}

// DecodeRichTextPayload is provided for custom native hosts and headless tools.
// It rejects unsupported versions, oversized fields, counts, and trailing bytes.
func DecodeRichTextPayload(payload string) ([]Span, error) {
	if len(payload) == 0 {
		return nil, nil
	}
	if len(payload) > maxRichTextBytes {
		return nil, errors.New("rich text payload exceeds maximum length")
	}
	r := bytes.NewReader([]byte(payload))
	var version uint16
	var count uint32
	if binary.Read(r, binary.LittleEndian, &version) != nil || version != 1 {
		return nil, errors.New("unsupported rich text payload")
	}
	if binary.Read(r, binary.LittleEndian, &count) != nil || count > MaxRichTextSpans {
		return nil, errors.New("invalid rich text span count")
	}
	spans := make([]Span, 0, count)
	for range count {
		text, err := readRichString(r)
		if err != nil {
			return nil, err
		}
		link, err := readRichString(r)
		if err != nil {
			return nil, err
		}
		span := Span{Text: text, Link: link}
		if err := binary.Read(r, binary.LittleEndian, &span.FontSize); err != nil {
			return nil, err
		}
		if err := binary.Read(r, binary.LittleEndian, &span.FontWeight); err != nil {
			return nil, err
		}
		flags, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		hasColor, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		rgba := make([]byte, 4)
		if _, err := io.ReadFull(r, rgba); err != nil {
			return nil, err
		}
		span.Underline, span.Italic, span.HasColor = flags&1 != 0, flags&2 != 0, hasColor != 0
		span.Color = RGBA(rgba[0], rgba[1], rgba[2], rgba[3])
		spans = append(spans, span)
	}
	if r.Len() != 0 {
		return nil, errors.New("trailing rich text data")
	}
	return spans, nil
}

func readRichString(r *bytes.Reader) (string, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	if length > maxRichTextBytes || uint64(length) > uint64(r.Len()) {
		return "", fmt.Errorf("invalid rich text string length %d", length)
	}
	value := make([]byte, length)
	if _, err := io.ReadFull(r, value); err != nil {
		return "", err
	}
	if !utf8.Valid(value) {
		return "", errors.New("invalid UTF-8 in rich text")
	}
	return string(value), nil
}

type InputValueMode uint8

const (
	InputControlled InputValueMode = iota
	InputUncontrolled
)

type InputKind uint8

const (
	InputText InputKind = iota
	InputEmail
	InputPhone
	InputURL
	InputInteger
	InputDecimal
	InputSearch
)

type ReturnKey uint8

const (
	ReturnKeyDefault ReturnKey = iota
	ReturnKeyDone
	ReturnKeyGo
	ReturnKeyNext
	ReturnKeySearch
	ReturnKeySend
)

type TextCapitalization uint8

const (
	CapitalizeNone TextCapitalization = iota
	CapitalizeSentences
	CapitalizeWords
	CapitalizeCharacters
)

type AutoCorrectMode uint8

const (
	AutoCorrectDefault AutoCorrectMode = iota
	AutoCorrectEnabled
	AutoCorrectDisabled
)

type ValidationState uint8

const (
	ValidationNone ValidationState = iota
	ValidationValid
	ValidationInvalid
)

// TextSelection uses UTF-16 code-unit offsets, matching UIKit and Android text APIs.
// {-1,-1} means the selection is not controlled by Go.
type TextSelection struct{ Start, End int32 }

// UncontrolledTextSelection returns the sentinel that leaves selection native-owned.
func UncontrolledTextSelection() TextSelection { return TextSelection{Start: -1, End: -1} }

// InputFormatter transforms a proposed edit before it reaches onChange.
// Implementations must be deterministic and safe to call from the runtime event path.
type InputFormatter interface {
	Format(previous, proposed string) string
}
type InputFormatterFunc func(previous, proposed string) string

func (f InputFormatterFunc) Format(previous, proposed string) string { return f(previous, proposed) }

func ApplyInputFormatters(previous, proposed string, formatters []InputFormatter) string {
	value := proposed
	for _, formatter := range formatters {
		if formatter != nil {
			value = formatter.Format(previous, value)
		}
	}
	return value
}

// UncontrolledTextInput creates a native-owned input initialized once with initialValue.
func UncontrolledTextInput(initialValue string, onChange func(string)) *element {
	e := TextInput(initialValue, onChange)
	e.node.Props.InputMode = InputUncontrolled
	return e
}

func (e *element) TextWrapping(mode TextWrapMode) *element { e.node.Props.TextWrap = mode; return e }
func (e *element) TextOverflow(value TextOverflow) *element {
	e.node.Props.TextOverflow = value
	return e
}
func (e *element) LineLimit(value uint32) *element    { e.node.Props.MaxLines = value; return e }
func (e *element) Selectable(value bool) *element     { e.node.Props.Selectable = value; return e }
func (e *element) InputType(value InputKind) *element { e.node.Props.InputKind = value; return e }
func (e *element) ReturnAction(value ReturnKey, onSubmit func()) *element {
	e.node.Props.ReturnKey, e.node.Submit = value, onSubmit
	return e
}
func (e *element) Capitalization(value TextCapitalization) *element {
	e.node.Props.Capitalization = value
	return e
}
func (e *element) AutoCorrect(value AutoCorrectMode) *element {
	e.node.Props.AutoCorrect = value
	return e
}
func (e *element) SecureEntry(value bool) *element    { e.node.Props.Secure = value; return e }
func (e *element) MultilineInput(value bool) *element { e.node.Props.Multiline = value; return e }
func (e *element) ReadOnly(value bool) *element       { e.node.Props.ReadOnly = value; return e }
func (e *element) Validation(state ValidationState, message string) *element {
	e.node.Props.Validation, e.node.Props.ErrorText = state, message
	return e
}
func (e *element) SelectionRange(value TextSelection, onChange func(TextSelection)) *element {
	e.node.Props.SelectionStart, e.node.Props.SelectionEnd, e.node.Selection = value.Start, value.End, onChange
	return e
}
func (e *element) MaximumLength(value int32) *element {
	if value < 0 {
		value = 0
	}
	e.node.Props.MaxLength = value
	return e
}
func (e *element) FormatWith(formatters ...InputFormatter) *element {
	e.node.InputFormatters = append(e.node.InputFormatters, formatters...)
	return e
}

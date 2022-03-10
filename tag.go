package html5tag

import (
	"fmt"
	"html"
	"io"
	"strings"
)

// The LabelDrawingMode describes how to draw a label when it is drawn.
// Various CSS frameworks expect it a certain way. Many are not very forgiving when
// you don't do it the way they expect.
type LabelDrawingMode int

const (
	// LabelDefault means the mode is defined elsewhere, like in a config setting
	LabelDefault LabelDrawingMode = iota
	// LabelBefore indicates the label is in front of the control.
	// Example: <label>MyLabel</label><input ... />
	LabelBefore
	// LabelAfter indicates the label is after the control.
	// Example: <input ... /><label>MyLabel</label>
	LabelAfter
	// LabelWrapBefore indicates the label is before the control's tag, and wraps the control tag.
	// Example: <label>MyLabel<input ... /></label>
	LabelWrapBefore
	// LabelWrapAfter indicates the label is after the control's tag, and wraps the control tag.
	// Example: <label><input ... />MyLabel</label>
	LabelWrapAfter
)

// VoidTag represents a void tag, which is a tag that does not need a matching closing tag.
type VoidTag struct {
	Tag  string
	Attr Attributes
}

// Render returns the rendered version of the tag.
func (t VoidTag) Render() string {
	return RenderVoidTag(t.Tag, t.Attr)
}

// RenderVoidTag renders a void tag using the given tag name and attributes.
func RenderVoidTag(tag string, attr Attributes) (s string) {
	b := strings.Builder{}
	_, err := WriteVoidTag(&b, tag, attr)
	if err != nil {
		panic(err)
	}
	return b.String()
}

// WriteVoidTag writes a void tag to the io.Writer.
func WriteVoidTag(w io.Writer, tag string, attr Attributes) (n int, err error) {
	return writeTag(w, tag, attr, nil, true, false, false)
}

// RenderTag renders a standard html tag with a closing tag.
//
// innerHtml is html, and must already be escaped if needed.
//
// The tag will be surrounded with newlines to force general formatting consistency.
// This will cause the tag to be rendered with a space between it and its neighbors if the tag is
// not a block tag.
//
// In the few situations where you would want to
// get rid of this space, call RenderTagNoSpace()
func RenderTag(tag string, attr Attributes, innerHtml string) string {
	b := strings.Builder{}
	var wto io.WriterTo
	if innerHtml != "" {
		wto = strings.NewReader(innerHtml)
	}

	_, err := WriteTag(&b, tag, attr, wto)
	if err != nil {
		panic(err)
	}
	return b.String()
}

// RenderTagFormatted renders the tag, pretty prints the innerHtml and sorts the attributes.
//
// Do not use this for tags where changing the innerHtml will change the appearance.
func RenderTagFormatted(tag string, attr Attributes, innerHtml string) string {
	b := strings.Builder{}
	var wto io.WriterTo
	if innerHtml != "" {
		wto = strings.NewReader(innerHtml)
	}
	_, err := WriteTagFormatted(&b, tag, attr, wto)
	if err != nil {
		panic(err)
	}
	return b.String()
}

// WriteTag writes the tag to the io.Writer.
func WriteTag(w io.Writer, tag string, attr Attributes, innerHtml io.WriterTo) (n int, err error) {
	return writeTag(w, tag, attr, innerHtml, false, false, false)
}

// WriteTagFormatted writes the tag to the io.Writer, pretty prints the innerHtml and sorts the attributes.
func WriteTagFormatted(w io.Writer, tag string, attr Attributes, innerHtml io.WriterTo) (n int, err error) {
	return writeTag(w, tag, attr, innerHtml, false, false, true)
}

// RenderTagNoSpace is similar to RenderTag, but should be used in situations where the tag is an
// inline tag that you want to visually be right next to its neighbors with no space.
func RenderTagNoSpace(tag string, attr Attributes, innerHtml string) string {
	b := strings.Builder{}
	var wto io.WriterTo
	if innerHtml != "" {
		wto = strings.NewReader(innerHtml)
	}
	_, err := WriteTagNoSpace(&b, tag, attr, wto)
	if err != nil {
		panic(err)
	}
	return b.String()
}

// WriteTagNoSpace writes the tag to the io.Writer, and does not add any spaces between the tag and the innerHtml.
func WriteTagNoSpace(w io.Writer, tag string, attr Attributes, innerHtml io.WriterTo) (n int, err error) {
	return writeTag(w, tag, attr, innerHtml, false, true, false)
}

// RenderTagNoSpaceFormatted will render without formatting the innerHtml, but WILL sort the attributes.
func RenderTagNoSpaceFormatted(tag string, attr Attributes, innerHtml string) string {
	b := strings.Builder{}
	var wto io.WriterTo
	if innerHtml != "" {
		wto = strings.NewReader(innerHtml)
	}

	_, err := WriteTagNoSpaceFormatted(&b, tag, attr, wto)
	if err != nil {
		panic(err)
	}
	return b.String()
}

// WriteTagNoSpaceFormatted writes to tag without formatting the innerHtml, but WILL sort the attributes.
func WriteTagNoSpaceFormatted(w io.Writer, tag string, attr Attributes, innerHtml io.WriterTo) (n int, err error) {
	return writeTag(w, tag, attr, innerHtml, false, true, true)
}

// writeString is a version of io.WriteString that accumulates the total written from previous writes.
func writeString(w io.Writer, s string, n int) (n2 int, err error) {
	n2, err = io.WriteString(w, s)
	n2 += n
	return
}

// writeTag is the main formatter of tags.
func writeTag(w io.Writer, tag string, attr Attributes, innerHtml io.WriterTo, isVoid bool, noSpace bool, format bool) (n int, err error) {
	var n3 int64

	if n, err = writeString(w, "<", n); err != nil {
		return
	}
	if n, err = writeString(w, tag, n); err != nil {
		return
	}
	if len(attr) != 0 {
		if n, err = writeString(w, " ", n); err != nil {
			return
		}

		if format {
			n3, err = attr.WriteSortedTo(w)
			n += int(n3)
			if err != nil {
				return
			}
		} else {
			n3, err = attr.WriteTo(w)
			n += int(n3)
			if err != nil {
				return
			}
		}
	}
	if n, err = writeString(w, ">", n); err != nil {
		return
	}

	if isVoid {
		return
	}

	if innerHtml != nil {
		builder := strings.Builder{}
		innerW := w
		var innerN int

		if format {
			innerW = &builder
		}
		if !noSpace {
			// required for consistency, will force a space between itself and its neighbors in certain situations
			if innerN, err = writeString(innerW, "\n", innerN); err != nil {
				return
			}
		}
		n3, err = innerHtml.WriteTo(innerW)
		innerN += int(n3)
		if err != nil {
			if !format {
				n += innerN
			}
			return
		}
		if !noSpace {
			if innerN, err = writeString(innerW, "\n", innerN); err != nil {
				if !format {
					n += innerN
				}
				return
			}
		}
		if format {
			s := builder.String()
			if !noSpace {
				s = Indent(s)
			}
			if n, err = writeString(w, s, n); err != nil {
				return
			}
		} else {
			n += innerN
		}
	}
	if n, err = writeString(w, "</", n); err != nil {
		return
	}
	if n, err = writeString(w, tag, n); err != nil {
		return
	}
	n, err = writeString(w, ">", n)
	return
}

// RenderLabel is a utility function to render a label, together with its text.
// Various CSS frameworks require labels to be rendered a certain way.
func RenderLabel(labelAttributes Attributes, label string, ctrlHtml string, mode LabelDrawingMode) string {
	b := strings.Builder{}

	var wto io.WriterTo
	if ctrlHtml != "" {
		wto = strings.NewReader(ctrlHtml)
	}
	_, err := WriteLabel(&b, labelAttributes, label, wto, mode)
	if err != nil {
		panic(err)
	}
	return b.String()
}

type writerItems []io.WriterTo

// WriteTo implements the io.WriterTo interface.
func (i writerItems) WriteTo(w io.Writer) (n int64, err error) {
	for _, item := range i {
		n2, err2 := item.WriteTo(w)
		n += n2
		if err2 != nil {
			return n, err2
		}
	}
	return
}

func makeWritersTo(items ...io.WriterTo) io.WriterTo {
	b := writerItems(items)
	return b
}

// WriteLabel is a utility function to render a label, together with its text.
// Various CSS frameworks require labels to be rendered a certain way.
func WriteLabel(w io.Writer, labelAttributes Attributes, label string, ctrlHtml io.WriterTo, mode LabelDrawingMode) (n int, err error) {
	var n64 int64
	var n2 int
	label = html.EscapeString(label)
	switch mode {
	case LabelBefore:
		if n, err = WriteTagNoSpace(w, "label", labelAttributes, strings.NewReader(label)); err != nil {
			return
		}
		if n, err = writeString(w, " ", n); err != nil {
			return
		}
		n64, err = ctrlHtml.WriteTo(w)
		n += int(n64)
		return
	case LabelAfter:
		n64, err = ctrlHtml.WriteTo(w)
		n += int(n64)
		if err != nil {
			return
		}
		if n, err = writeString(w, " ", n); err != nil {
			return
		}
		n2, err = WriteTagNoSpace(w, "label", labelAttributes, strings.NewReader(label))
		n += n2
		return
	case LabelWrapBefore:
		return WriteTag(w, "label", labelAttributes, makeWritersTo(strings.NewReader(label+" "), ctrlHtml))
	case LabelWrapAfter:
		return WriteTag(w, "label", labelAttributes, makeWritersTo(ctrlHtml, strings.NewReader(" "+label)))
	}
	panic("Unknown label mode")
}

// RenderImage renders an image tag with the given source, alt and attribute values.
// Panics on error.
func RenderImage(src string, alt string, attributes Attributes) string {
	b := strings.Builder{}
	_, err := WriteImage(&b, src, alt, attributes)
	if err != nil {
		panic(err)
	}
	return b.String()
}

// WriteImage writes an image tag.
func WriteImage(w io.Writer, src string, alt string, attributes Attributes) (n int, err error) {
	a := attributes.Copy().Set("src", src).Set("alt", alt)
	return WriteVoidTag(w, "img", a)
}

// Indent will add space to the front of every line in the string. Since indent is used to format code for reading
// while we are in development mode, we do not need it to be particularly efficient.
// It will not do this for textarea tags, since that would change the text in the tag.
func Indent(s string) string {
	var out string
	var taOffset int
	for {
		taOffset = strings.Index(s, "<textarea")
		if taOffset == -1 {
			out += indent(s)
			return out
		}
		out += indent(s[:taOffset])
		s = s[taOffset:]
		taOffset = strings.Index(s, "</textarea>")
		if taOffset == -1 {
			// This is an error in the html, so just return the original
			return s
		}
		out += s[:taOffset+11] // skip textarea close tag
		s = s[taOffset+11:]
	}
}

// indents the string unsafely, in that it does not check for allowable tags to indent
func indent(s string) string {
	var newLines []string
	a := strings.Split(s, "\n")
	for _, l := range a {
		if l != "" {
			l = "  " + l
		}
		newLines = append(newLines, l)
	}
	return strings.Join(newLines, "\n")
}

// Comment turns the given text into an HTML comment and returns the rendered comment
func Comment(s string) string {
	return fmt.Sprintf("<!-- %s -->", s)
}

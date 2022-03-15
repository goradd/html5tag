package html5tag

import (
	"html"
)

var voidTags = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

// A TagBuilder creates a tag using a builder pattern, starting out with the
// tag name and slowly adding parts to it, describing it, until you are ready to print
// out the entire html tag. The zero value is usable.
type TagBuilder struct {
	tag        string
	attributes Attributes
	innerHtml  string
	isVoid     bool
}

// NewTagBuilder starts a tag build, though you can use a tag builder from its zero value too.
func NewTagBuilder() *TagBuilder {
	return &TagBuilder{}
}

// Tag sets the tag value
func (b *TagBuilder) Tag(tag string) *TagBuilder {
	b.tag = tag
	b.isVoid, _ = voidTags[tag]
	return b
}

// Set sets the attribute to the given value
func (b *TagBuilder) Set(attribute string, value string) *TagBuilder {
	if b.attributes == nil {
		b.attributes = NewAttributes()
	}
	b.attributes.Set(attribute, value)
	return b
}

// ID sets the id attribute
func (b *TagBuilder) ID(id string) *TagBuilder {
	b.Set("id", id)
	return b
}

// Class sets the class attribute to the value given.
// If you prefix the value with "+ " the given value will be appended to the end of the current class list.
// If you prefix the value with "- " the given value will be removed from the class list.
// Otherwise, the current class value is replaced.
// The given class can be multiple classes separated by a space.
func (b *TagBuilder) Class(class string) *TagBuilder {
	if b.attributes == nil {
		b.attributes = NewAttributes()
	}
	b.attributes.SetClass(class)
	return b
}

// Link is a shortcut that will set the tag to "a" and the "href" to the given destination.
// This is not the same as an actual "link" tag, which points to resources from the header.
func (b *TagBuilder) Link(href string) *TagBuilder {
	b.tag = "a"
	b.Set("href", href)
	return b
}

// IsVoid will make the builder output a void tag instead of one with inner html.
func (b *TagBuilder) IsVoid() *TagBuilder {
	b.isVoid = true
	return b
}

// InnerHtml sets the inner html of the tag.
//
// Remember this is HTML, and will not be escaped.
func (b *TagBuilder) InnerHtml(html string) *TagBuilder {
	b.innerHtml = html
	return b
}

// InnerText sets the inner part of the tag to the given text. The text will be escaped.
func (b *TagBuilder) InnerText(text string) *TagBuilder {
	b.innerHtml = html.EscapeString(text)
	return b
}

// String ends the builder and returns the html.
func (b *TagBuilder) String() string {
	if b.tag == "" {
		panic("You cannot output the tag builder with no tag")
	}
	if b.isVoid {
		return RenderVoidTag(b.tag, b.attributes)
	}
	return RenderTag(b.tag, b.attributes, b.innerHtml)
}

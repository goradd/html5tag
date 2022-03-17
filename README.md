[![Go Reference](https://pkg.go.dev/badge/github.com/goradd/html5tag.svg)](https://pkg.go.dev/github.com/goradd/html5tag)
![Build Status](https://img.shields.io/github/workflow/status/goradd/got/Go)
[![Go Report Card](https://goreportcard.com/badge/github.com/goradd/html5tag)](https://goreportcard.com/report/github.com/goradd/html5tag)
[![codecov](https://codecov.io/gh/goradd/html5tag/branch/main/graph/badge.svg?token=L8KC75KWWR)](https://codecov.io/gh/goradd/html5tag)

# html5tag

The html5tag package contains utilities to generate html 5 tags. 
Choose between string versions of the
functions for easy tag creation, or io.Writer versions for speed.

html5tag also has a tag builder for convenience and can perform math operations
on numeric style values.

html5tag does some checks to make sure tags are well-formed. For example,
when adding data-* attributes, it will make sure the key used for the
attribute does not violate html syntax rules.

html5tag has options to pretty-print tags and the content of tags so they appear formatted
in a document. However, in certain contexts, like in inline text, or in a textarea tag, adding
extra returns and spaces changes the look of the output. In these situations, use the functions
that do not add spaces to the inner HTML of a tag.

Some examples:

```go
package main

import . "github.com/goradd/html5tag"

main() {
	
	// Render an input tag, inside a div tag, inside a body tag using different tag building mechanisms

	a := NewAttributes().
	SetID("myText").
	Set("text", "Hi there").
	Set("placeholder", "Write here").
	SetClass("class1 class2").
	SetStyle("width":"20em")
	
	inputTag := RenderVoidTag("input", a)
	divTag := RenderTag("div", Attriubtes{"id":"inputWrapper"}, inputTag)
	
	bodyTag := NewTagBuilder().
		Tag("body").
		ID("bodyId").
		InnerHtml(divTag).
		String()
	
	fmt.Print(bodyTag)
}
```

For complete documentation, start at the documentation for `RenderTag()` and `WriteTag()` and drill down from there.

package html5tag

import (
	"fmt"
	"testing"
)

func TestRandomString(t *testing.T) {
	s := RandomString(40)
	fmt.Printf(s + " ")

	if len(s) != 40 {
		t.Error("Wrong size")
	}
}

func ExampleTextToHtml() {
	s := TextToHtml("This is a & test.\n\nA paragraph\nwith a forced break.")
	fmt.Println(s)
	// Output: This is a &amp; test.<p>A paragraph<br />with a forced break.
}

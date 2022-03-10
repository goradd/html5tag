package html5tag

import (
	"fmt"
	"testing"
)

func ExampleToDataAttr() {
	s, _ := ToDataAttr("thisIsMyTest")
	fmt.Println(s)
	// Output: this-is-my-test

}

func ExampleToDataJqKey() {
	s, _ := ToDataKey("this-is-my-test")
	fmt.Println(s)
	// Output: thisIsMyTest

}

func TestToDataAttr(t *testing.T) {

	cases := []struct {
		in, expected string
		err          bool
	}{
		{"ThisThat", "", true},
		{"thisANDthat", "", true},
		{"That", "", true},
		{"", "", false},
		{"this", "this", false},
		{"thisAndThat", "this-and-that", false},
		{"this and that", "", true},
	}

	for _, c := range cases {
		result, err := ToDataAttr(c.in)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q): %v", c.in, err)
				continue
			}
		}

		if c.err && err == nil { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q)", c.in)
			continue // no sense in checking other things, since we were expecting an error
		}

		if result != c.expected {
			t.Errorf("ToDataAttr failed on (%q) expected (%q) got (%q)", c.in, c.expected, result)
		}
	}

}

func TestToDataKey(t *testing.T) {

	cases := []struct {
		in, expected string
		err          bool
	}{
		{"ThisThat", "", true},
		{"thisANDthat", "", true},
		{"That", "", true},
		{"", "", false},
		{"this", "this", false},
		{"this-and-that", "thisAndThat", false},
		{"this and that", "", true},
		{"a-b-c", "", true},
	}

	for _, c := range cases {
		result, err := ToDataKey(c.in)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q): %v", c.in, err)
				continue
			}
		}

		if c.err && err == nil { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q)", c.in)
			continue // no sense in checking other things, since we were expecting an error
		}

		if result != c.expected {
			t.Errorf("ToDataKey failed on (%q) expected (%q) got (%q)", c.in, c.expected, result)
		}
	}

}

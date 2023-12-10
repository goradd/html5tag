package html5tag

import (
	"fmt"
	"strconv"
	"testing"
)

func TestBasicAttributes(t *testing.T) {
	cases := []struct {
		attr, val    string
		changed, err bool
	}{
		{"a", "t", true, false},
		{"b", "t", true, false},
		{"id", "9", true, false},
		{"a", "t", false, false},
		{"id", "t y", true, true},
		{"class", "t y", true, false},
		{"bad name", "t y", true, true},
		{"style", "a b", false, true},
	}

	a := NewAttributes()
	for _, c := range cases {
		changed, err := a.SetChanged(c.attr, c.val)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q, %q): %v", c.attr, c.val, err)
				continue
			}
		}

		if c.err { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q, %q)", c.attr, c.val)
			continue // no sense in checking other things, since we were expecting an error
		}

		got := a.Get(c.attr)
		if got != c.val {
			t.Errorf("Basic set/get (%q, %q)", c.attr, c.val)
		}

		if changed != c.changed {
			t.Errorf("Basic changed test (%q, %q)", c.attr, c.val)
		}
	}

	a.Set("a", FalseValue)
	if a.Has("a") {
		t.Error("Missing 'a' attribute")
	}
	if a.Has("style") {
		t.Error("Should not have style attribute")
	}
}

func TestBasicStyles(t *testing.T) {
	cases := []struct {
		attr         string
		val          string
		changed, err bool
	}{
		{"color", "blue", true, false},
		{"color", "red", true, false},
		{"color", "red", false, false},
		{"height", "4px", true, false},
	}

	var changed bool
	var err error

	a := NewAttributes()
	for _, c := range cases {
		changed, err = a.SetStyleChanged(c.attr, c.val)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q, %q): %v", c.attr, c.val, err)
				continue
			}
		}

		if c.err { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q, %q)", c.attr, c.val)
			continue // no sense in checking other things, since we were expecting an error
		}

		got := a.GetStyle(c.attr)
		if got != c.val {
			t.Errorf("BasicStyle set/get (%q, %q), got %q", c.attr, c.val, got)
		}
		if changed != c.changed {
			t.Errorf("BasicStyle changed test (%q, %q)", c.attr, c.val)
		}
	}

	a = NewAttributes()
	if changed, err = a.SetChanged("style", "height:4px;width:6px"); !changed || err != nil {
		t.Error("Problem setting style")
	}

	a.SetStyle("width", "+ 2")

	if changed, err = a.SetChanged("style", "width:8px; height:4px"); changed || err != nil {
		t.Errorf("Problem setting same style again: (%v, %v)", changed, err)
	}

	if a.GetStyle("width") != "8px" {
		t.Error("Problem with setting style")
	}

	a = NewAttributes()
	var b []bool
	changed, _ = a.SetStyleChanged("height", "9")
	b = append(b, changed)
	changed, _ = a.SetStyleChanged("height", "10")
	b = append(b, changed)
	changed, _ = a.SetStyleChanged("height", "10")
	b = append(b, changed)
	changed, _ = a.SetStyleChanged("width", "4")
	b = append(b, changed)
	out := fmt.Sprint(b)
	if out != "[true true false true]" {
		t.Error("Changing styles failed. Got: " + out)
	}

}

func TestClass(t *testing.T) {
	cases := []struct {
		val     string
		got     string
		changed bool
	}{
		{"c1", "c1", true},
		{"c2", "c2", true},
		{"c1 c2", "c1 c2", true},
		{"+ c3", "c1 c2 c3", true},
		{"+ c3", "c1 c2 c3", false},
		{"- c1", "c2 c3", true},
	}

	a := NewAttributes()
	for _, c := range cases {
		changed := a.SetClassChanged(c.val)
		got := a.Class()
		if got != c.got {
			t.Errorf("Class set (%q), expected (%q), got (%q)", c.val, c.got, got)
		}

		if changed != c.changed {
			t.Errorf("Class changed test (%q)", c.val)
		}
	}

}

func TestDataAttributes(t *testing.T) {
	cases := []struct {
		attr, val    string
		changed, err bool
	}{
		{"data-a", "t", true, false},
		{"data-b", "t", true, false},
		{"data-b", "t", false, false},
		{"data-id", "9", true, false},
		{"data-$a", "t", false, true},
		{"data-bad name", "t y", true, true},
	}

	a := NewAttributes()
	for _, c := range cases {
		changed, err := a.SetChanged(c.attr, c.val)
		if err != nil {
			if c.err { // expected an error
				continue
			} else {
				t.Errorf("Unexpected error on (%q, %q): %v", c.attr, c.val, err)
				continue
			}
		}

		if c.err { // expected an error, but didn't get one
			t.Errorf("Expected error on (%q, %q)", c.attr, c.val)
			continue // no sense in checking other things, since we were expecting an error
		}

		got := a.DataAttribute(c.attr[5:])
		if got != c.val {
			t.Errorf("Data Attribute set/get (%q, %q)", c.attr, c.val)
		}

		if changed != c.changed {
			t.Errorf("Data Attribute changed test (%q, %q)", c.attr, c.val)
		}

		if !a.HasDataAttribute(c.attr[5:]) {
			t.Errorf("Has data attribute (%q)", c.attr)
		}
	}

	a.RemoveDataAttribute("data-id")

	if a.HasDataAttribute("data-id") {
		t.Error("Removed data attribute (data-id)")
	}

}

func TestOutput(t *testing.T) {
	var s string
	a := NewAttributes()
	a.Set("class", "a")
	a.Set("id", "b")

	s = a.String()

	if !(s == `class="a" id="b"` || s == `id="b" class="a"`) {
		t.Error("No style test")
	}

	a.RemoveAttribute("class")
	a.RemoveAttribute("id")
	a.SetStyle("height", "4")

	s = a.String()

	if !(s == `style="height:4px"`) {
		t.Error("With style test: " + s)
	}

	// Test some escapes here
	a = NewAttributes()
	a.Set("placeholder", "<& I Am Here >")
	expected := "placeholder=\"&lt;&amp; I Am Here &gt;\""
	if s = a.String(); s != expected {
		t.Errorf("Not escaping. Expected (%q) got (%q)", expected, s)
	}

	a = Attributes{"ok": "", "id": "3"}
	if `id="3" ok` != a.SortedString() {
		t.Error("Sorted string failed")
	}
}

func TestOverride(t *testing.T) {
	a := NewAttributes()
	a.Set("class", "a")
	a.Set("id", "b")
	a.Set("style", "height:4px; width:3px")

	m := a.Override(map[string]string{"id": "c", "style": "height:7px"})

	if m.Get("id") != "c" {
		t.Errorf("Error overriding id. Wanted 'c', got %s", m.Get("id"))
	}

	if m.GetStyle("height") != "7px" {
		t.Errorf("Error overriding height style. Wanted 7px, got %s", m.GetStyle("height"))
	}

	if m.GetStyle("width") != "" {
		t.Error("Error overriding style")
	}

}

// Examples
func ExampleAttributes_Set() {
	a := Attributes{}
	a = a.Set("class", "a").Set("id", "b")
	fmt.Println(a.SortedString())
	//Output: id="b" class="a"
}

func ExampleAttributes_SetClass() {
	a := NewAttributes()
	a.SetClass("this")
	a.SetClass("+ that")
	s := a.Class()
	fmt.Println(s)
	a.SetClass("")
	fmt.Println(a.Has("class"))
	// Output: this that
	// false
}

func ExampleAttributes_SetStyle() {
	a := NewAttributes()
	a.SetStyle("height", "4em")
	a.SetStyle("width", "8")
	a.SetStyle("width", "- 2")
	fmt.Println(a.GetStyle("height"))
	fmt.Println(a.GetStyle("width"))
	// Output:
	// 4em
	// 6px
}

func ExampleAttributes_SetID() {
	a := Attributes{}
	a = a.SetID("a")
	fmt.Println(a.ID())
	a = a.SetID("")
	fmt.Println(a.Has("id"))
	//Output: a
	// false
}

func ExampleAttributes_Override() {
	a := NewAttributes().SetClass("this").SetStyle("height", "4em")
	b := NewAttributes().Set("class", "that").SetStyle("width", "6")

	a = a.Override(b)
	fmt.Println(a.SortedString())
	//Output: class="that" style="width:6px"
}

func ExampleAttributes_Merge() {
	a := NewAttributes().SetClass("this").SetStyle("height", "4em")
	b := NewAttributes().Set("class", "that").SetStyle("width", "6")

	a = a.Override(b)
	fmt.Println(a.SortedString())
	// Output: class="that" style="width:6px"
}

func ExampleAttributes_AddClass() {
	a := NewAttributes()
	a.AddClass("this")
	a.AddClass("that")
	a.AddClass("")
	fmt.Println(a)
	//Output: class="this that"
}

func ExampleAttributes_HasClass() {
	a := NewAttributes()
	if !a.HasClass("that") {
		fmt.Println("Not found")
	}
	a.SetClass("this that other")
	if a.HasClass("that") {
		fmt.Println("found")
	}
	// Output: Not found
	// found
}

func ExampleAttributes_HasStyle() {
	a := NewAttributes()
	var b []bool
	var found bool
	found = a.HasStyle("height")
	b = append(b, found)
	a.SetStyle("height", strconv.Itoa(10))
	found = a.HasStyle("height")
	b = append(b, found)
	fmt.Println(b)
	// Output: [false true]
}

func ExampleAttributes_RemoveStyle() {
	a := NewAttributes()
	a.SetStyle("height", "10")
	a.SetStyle("width", strconv.Itoa(5))
	a.RemoveStyle("height")
	fmt.Println(a)
	// Output: style="width:5px"
}

func ExampleAttributes_RemoveClass() {
	a := Attributes{"class": "this that"}
	changed := a.RemoveClass("this")
	fmt.Println(changed)
	fmt.Println(a.String())
	changed = a.RemoveClass("other")
	fmt.Println(changed)
	fmt.Println(a.String())
	// Output: true
	// class="that"
	// false
	// class="that"
}

func ExampleAttributes_RemoveClassesWithPrefix() {
	a := Attributes{"class": "col-2 that"}
	a.RemoveClassesWithPrefix("col-")
	fmt.Println(a.String())
	// Output: class="that"
}

func ExampleAttributes_HasClassWithPrefix() {
	a := Attributes{"class": "col-2 that"}
	found := a.HasClassWithPrefix("col-")
	fmt.Println(found)
	// Output: true
}

func ExampleAttributes_AddValues() {
	a := Attributes{"abc": "123"}
	a.AddValues("abc", "456")
	fmt.Println(a.String())
	// Output: abc="123 456"
}

func ExampleAttributes_SetData() {
	a := Attributes{"abc": "123"}
	a.SetData("myVal", "456")
	fmt.Println(a.SortedString())
	// Output: abc="123" data-my-val="456"
}

func ExampleAttributes_SetStyles() {
	a := Attributes{"style": "color:blue"}
	s := Style{"color": "yellow"}
	a.SetStyles(s)
	fmt.Println(a.String())
	// Output: style="color:yellow"
}

func ExampleAttributes_SetStylesTo() {
	a := Attributes{"style": "color:blue"}
	a.SetStylesTo("color:red")
	fmt.Println(a.String())
	// Output: style="color:red"
}

func ExampleAttributes_SetDisabled() {
	a := Attributes{"style": "color:blue"}
	a.SetDisabled(true)
	fmt.Println(a.SortedString())
	a.SetDisabled(false)
	fmt.Println(a.SortedString())
	// Output: style="color:blue" disabled
	// style="color:blue"
}

func ExampleAttributes_SetDisplay() {
	a := Attributes{"style": "color:blue"}
	a.SetDisplay("none")
	fmt.Println(a.SortedString())
	// Output: style="color:blue;display:none"
}

func ExampleAttributes_IsDisplayed() {
	a := Attributes{"style": "color:blue"}
	a.SetDisplay("none")
	fmt.Println(a.IsDisplayed())
	// Output: false
}

func ExampleValueString() {
	a := Attributes{}
	a.Set("a", ValueString(1))
	a.Set("b", ValueString(float32(2.2)))
	a.Set("c", ValueString("test"))
	a.Set("d", ValueString(true))
	a.Set("e", ValueString(false))
	fmt.Println(a.SortedString())
	// Output: a="1" b="2.2" c="test" d
}

func TestMergeString(t *testing.T) {
	a := NewAttributes()
	a.MergeString(`class="here"`)
	c := a.Class()
	if c != "here" {
		t.Error("Attribute string failed")
	}

	a.MergeString(`class="there" m="g"`)
	c = a.Class()
	if c != "here there" {
		t.Error("Attribute string failed")
	}
	d := a.Get("m")
	if d != "g" {
		t.Error("Attribute string failed")
	}
	a.Merge(nil)
	if a.Len() != 2 {
		t.Error("Nil merge failed")
	}

	a.Merge(Attributes{"style": "color:white"})
	if !a.Has("style") {
		t.Error("Style merge failed")
	}

	a.Merge(Attributes{"style": "color:black"})
	if !a.HasStyle("color") {
		t.Error("Color style merge failed")
	}
	a.Merge(map[string]string{"style": "color:yellow"})
	if a.GetStyle("color") != "yellow" {
		t.Error("Color style override failed")
	}
}

func TestNilAttributes(t *testing.T) {
	var a Attributes
	if a.Len() != 0 {
		t.Error("Nil length failed")
	}
	if a.Has("b") {
		t.Error("Nil Has failed")
	}
	if a.String() != "" {
		t.Error("Nil String failed")
	}
	a.Range(func(k string, v string) bool {
		t.Error("Should not range")
		return false
	})
	if a.ID() != "" {
		t.Error("ID should be empty")
	}
}

func ExampleAttributes_Len() {
	a := Attributes{"id": "45", "class": "aclass"}
	fmt.Print(a.Len())
	//Output: 2
}

func ExampleAttributes_Range() {
	a := Attributes{"y": "7", "x": "10", "id": "1", "class": "2", "z": "4"}
	a.Range(func(k string, v string) bool {
		if k == "z" {
			return false
		}
		fmt.Println(k, "=", v)
		return true
	})
	// Output: id = 1
	// class = 2
	// x = 10
	// y = 7
}

func TestAttributes_RemoveClass(t *testing.T) {
	tests := []struct {
		name        string
		a           Attributes
		removeClass string
		changed     bool
		finalClass  string
	}{
		{"remove one", Attributes{"id": "1", "class": "this"}, "this", true, ""},
		{"remove from multiple", Attributes{"id": "1", "class": "this that"}, "this", true, "that"},
		{"remove from none", Attributes{"id": "1"}, "this", false, ""},
		{"remove not existing", Attributes{"id": "1", "class": "this that"}, "other", false, "this that"},
		{"remove multiple", Attributes{"id": "1", "class": "this that other"}, "this other", true, "that"},
		{"remove multiple one not existing", Attributes{"id": "1", "class": "this that other"}, "nothere other", true, "this that"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed := tt.a.RemoveClass(tt.removeClass)
			if tt.changed != changed {
				t.Errorf("Error on changed. Wanted %v, got %v", tt.changed, changed)
			}
			if tt.finalClass != tt.a.Class() {
				t.Errorf("Error on class. Wanted %v, got %v", tt.finalClass, tt.a.Class())
			}
		})
	}
}

func ExampleAttributes_IsDisabled() {
	a := Attributes{"disabled": ""}
	fmt.Print(a.IsDisabled())
	// Output: true
}

func BenchmarkSortAttr(b *testing.B) {
	a := Attributes{"a": "b", "id": "c", "width": "14", "d": "e"}

	for i := 0; i < b.N; i++ {
		_ = a.String()
	}
}
func BenchmarkSortedKeys(b *testing.B) {
	a := Attributes{"a": "b", "id": "c", "width": "14", "d": "e"}

	for i := 0; i < b.N; i++ {
		a.sortedKeys()
	}
}

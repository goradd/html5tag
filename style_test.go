package html5tag

import (
	"fmt"
	"testing"
)

func ExampleStyle_Copy() {
	s := Style{"color": "green", "size": "9"}
	s2 := s.Copy()

	fmt.Print(s2)
	//Output: color:green;size:9
}

func ExampleStyle_Len() {
	s := Style{"color": "green", "size": "9"}
	fmt.Print(s.Len())
	//Output: 2
}

func ExampleStyle_SetString() {
	s := NewStyle()
	_, _ = s.SetString("height: 9em; width: 100%; position:absolute")
	fmt.Print(s)
	//Output: height:9em;position:absolute;width:100%
}

func ExampleStyle_Set_a() {
	s := NewStyle()
	s.Set("height", "9")
	fmt.Print(s)
	//Output: height:9px
}

func ExampleStyle_Set_b() {
	s := NewStyle()
	_, _ = s.SetString("height:9px")
	s.Set("height", "+ 10")
	fmt.Print(s)
	//Output: height:19px
}

func ExampleStyle_Get() {
	s := NewStyle()
	_, _ = s.SetString("height: 9em; width: 100%; position:absolute")
	fmt.Print(s.Get("width"))
	//Output: 100%
}

func ExampleStyle_Remove() {
	s := NewStyle()
	_, _ = s.SetString("height: 9em; width: 100%; position:absolute")
	s.Remove("position")
	fmt.Print(s)
	//Output: height:9em;width:100%
}

func ExampleStyle_RemoveAll() {
	s := NewStyle()
	_, _ = s.SetString("height: 9em; width: 100%; position:absolute")
	s.RemoveAll()
	fmt.Print(s)
	//Output:
}

func ExampleStyle_Has() {
	s := NewStyle()
	_, _ = s.SetString("height: 9em; width: 100%; position:absolute")
	fmt.Print(s.Has("width"), s.Has("display"))
	//Output:true false
}

func TestStyleSet(t *testing.T) {
	s := NewStyle()

	changed, err := s.SetChanged("height", "4")
	if !changed {
		t.Error("Expected a change")
	}
	if err != nil {
		t.Error(err)
	}

	s.RemoveAll()
	if s.Has("height") {
		t.Error("Expected no height")
	}

	s.Set("height", "4")
	changed, err = s.SetString("height: 3; width: 5")
	if !changed {
		t.Error("Expected a change")
	}
	if err != nil {
		t.Error(err)
	}

	if !changed {
		t.Error("Expected a change")
	}
	if "5px" != s.Get("width") {
		t.Errorf("Wrong width. Expected 5px, got %s.", s.Get("width"))
	}
	if "3px" != s.Get("height") {
		t.Errorf("Wrong height. Expected 5px, got %s.", s.Get("height"))
	}

	// test error
	changed, err = s.SetString("height of: 3; width: 4")
	if changed {
		t.Error("Expected no change")
	}
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestStyleLengths(t *testing.T) {
	s := NewStyle()

	s.Set("height", "4px")
	changed, err := s.SetChanged("height", "4")

	if changed {
		t.Error("Expected no change")
	}
	if err != nil {
		t.Error(err)
	}

	changed, err = s.SetChanged("height", "4em")
	if !changed {
		t.Error("Expected change")
	}
	if err != nil {
		t.Error(err)
	}

	changed, err = s.SetChanged("width", "0")
	if !changed {
		t.Error("Expected change")
	}
	if err != nil {
		t.Error(err)
	}

	if w := s.Get("width"); w != "0" {
		t.Error("Expected a 0")
	}

	changed, err = s.SetChanged("width", "1")
	if !changed {
		t.Error("Expected change")
	}
	if err != nil {
		t.Error(err)
	}

	if w := s.Get("width"); w != "1px" {
		t.Error("Expected a 1px")
	}

	// test a non-length numeric
	changed, err = s.SetChanged("volume", "4")
	if !changed {
		t.Error("Expected change")
	}
	if err != nil {
		t.Error(err)
	}

	if w := s.Get("volume"); w != "4" {
		t.Error("Expected a 4")
	}
}

func TestStyleMath(t *testing.T) {
	s := NewStyle()

	s.Set("height", "4em")
	s.Set("height", "* 2")
	if h := s.Get("height"); h != "8em" {
		t.Error("Expected 8em, got " + h)
	}

	s.Set("height", "2em 9px")
	s.Set("height", "/ 2")
	if h := s.Get("height"); h != "1em 4.5px" {
		t.Error("Expected 1em 4.5px, got " + h)
	}

	s.Set("width", "7.6in")
	s.Set("width", "+ 2")
	if h := s.Get("width"); h != "9.6in" {
		t.Error("Expected 9.6in, got " + h)
	}

	s.Set("width", "1.6in")
	s.Set("width", "- 2")
	if h := s.Get("width"); h != "-0.4in" { // this test in particular can produce a rounding error if not handled carefully
		t.Error("Expected -0.4in, got " + h)
	}
}

func TestStyle(t *testing.T) {
	s := NewStyle()

	s.Set("position", "9")

	if a := s.Get("position"); a != "9px" {
		t.Error("Style test failed: " + a)
	}
}

func TestNilStyle(t *testing.T) {
	var s Style

	if s.Len() != 0 {
		t.Error("Nil style should have zero length")
	}

	if s.Has("a") {
		t.Error("Nil style should be empty")
	}
}

func TestStyle_mathOp(t *testing.T) {
	c := Style{"height": "10", "margin": "", "width": "20en"}

	type args struct {
		attribute string
		op        string
		val       string
	}
	tests := []struct {
		name        string
		s           Style
		args        args
		wantChanged bool
		wantErr     bool
		wantString  string
	}{
		{"Test empty", c.Copy(), args{"margin", "+", "1"}, true, false, "height:10;margin:1;width:20en"},
		{"Test float error", c.Copy(), args{"margin", "+", "1a"}, false, true, "height:10;margin:;width:20en"},
		{"Test mul no unit", c.Copy(), args{"height", "*", "2"}, true, false, "height:20;margin:;width:20en"},
		{"Test div w/ unit", c.Copy(), args{"width", "/", "2"}, true, false, "height:10;margin:;width:10en"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChanged, err := tt.s.mathOp(tt.args.attribute, tt.args.op, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("mathOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChanged != tt.wantChanged {
				t.Errorf("mathOp() gotChanged = %v, want %v", gotChanged, tt.wantChanged)
			}
			if tt.wantString != fmt.Sprint(tt.s) {
				t.Errorf("mathOp() got = %v, want %v", fmt.Sprint(tt.s), tt.wantString)
			}
		})
	}
}

func TestStyleString(t *testing.T) {
	tests := []struct {
		name string
		i    interface{}
		want string
	}{
		{"int", int(5), "5px"},
		{"float", float32(5.1), "5.1px"},
		{"double", float64(5.2), "5.2px"},
		{"string", "9em", "9em"},
		{"string2", "9", "9"},
		{"Stringer", NewStyle(), ""},
		{"default", []string{"a", "b"}, "[a b]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StyleString(tt.i); got != tt.want {
				t.Errorf("StyleString() = %v, want %v", got, tt.want)
			}
		})
	}
}

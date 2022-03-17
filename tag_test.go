package html5tag

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type errBuffer struct {
	cap int
	buf *bytes.Buffer
}

func (b *errBuffer) Write(p []byte) (n int, err error) {
	if len(p)+b.buf.Len() > b.cap {
		l := b.cap - b.buf.Len()
		n, _ = b.Write(p[:l])
		return n, fmt.Errorf("out of memory")
	}
	return b.buf.Write(p)
}

func newErrBuf(cap int) *errBuffer {
	return &errBuffer{buf: &bytes.Buffer{}, cap: cap}
}

func ExampleVoidTag_Render() {
	v := VoidTag{"br", Attributes{"id": "hi"}}
	fmt.Println(v.Render())
	//Output: <br id="hi">
}

func ExampleRenderTagNoSpace() {
	fmt.Println(RenderTagNoSpace("div", Attributes{"id": "me"}, "Here I am"))
	// Output: <div id="me">Here I am</div>
}

func ExampleRenderVoidTag() {
	fmt.Println(RenderVoidTag("img", Attributes{"src": "thisFile"}))
	// Output: <img src="thisFile">
}

func ExampleRenderLabel() {
	s1 := RenderLabel(nil, "Title", "<input>", LabelBefore)
	s2 := RenderLabel(nil, "Title", "<input>", LabelAfter)
	s3 := RenderLabel(nil, "Title", "<input>", LabelWrapBefore)
	s4 := RenderLabel(nil, "Title", "<input>", LabelWrapAfter)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)
	// Output: <label>Title</label> <input>
	// <input> <label>Title</label>
	// <label>
	// Title <input>
	// </label>
	// <label>
	// <input> Title
	// </label>
}

func TestRenderTagNoSpace(t *testing.T) {
	type args struct {
		tag       string
		attr      Attributes
		innerHtml string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test empty", args{"p", nil, ""}, `<p></p>`},
		{"Test empty with attributes", args{"p", Attributes{"height": "10"}, ""}, `<p height="10"></p>`},
		{"Test text", args{"p", nil, "I am here"}, `<p>I am here</p>`},
		{"Test html", args{"p", nil, "<p>I am here</p>"}, `<p><p>I am here</p></p>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderTagNoSpace(tt.args.tag, tt.args.attr, tt.args.innerHtml); got != tt.want {
				t.Errorf("RenderTagNoSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleComment() {
	s := Comment("This is a test")
	fmt.Print(s)
	//Output: <!-- This is a test -->
}

func BenchmarkWriteVoidTag(b *testing.B) {
	buf := bytes.Buffer{}
	s := "tag"
	var n int
	a := Attributes{"a": "b"}
	for i := 0; i < b.N; i++ {
		n2, _ := WriteVoidTag(&buf, s, a)
		n += n2
	}
}

func BenchmarkRenderVoidTag(b *testing.B) {
	s := "tag"
	var s2 string
	a := Attributes{"a": "b"}
	for i := 0; i < b.N; i++ {
		s2 = RenderVoidTag(s, a)
		_ = s2
	}
}

func BenchmarkWriteTag(b *testing.B) {
	buf := bytes.Buffer{}
	s := "tag"
	var n int
	a := Attributes{"a": "b"}
	w2 := strings.NewReader("abc")
	for i := 0; i < b.N; i++ {
		n2, _ := WriteTag(&buf, s, a, w2)
		n += n2
	}
}

func BenchmarkWriterTag(b *testing.B) {
	buf := bytes.Buffer{}
	s := "tag"
	var n int
	a := Attributes{"a": "b"}
	w2 := strings.NewReader("abc" + s + "cd")
	for i := 0; i < b.N; i++ {
		n2, _ := WriteTag(&buf, s, a, w2)
		n += n2
	}
}

func BenchmarkWriterTag2(b *testing.B) {
	buf := bytes.Buffer{}
	s := "tag"
	var n int
	a := Attributes{"a": "b"}
	w2 := makeWritersTo(strings.NewReader("abc"), strings.NewReader(s), strings.NewReader("cd"))
	for i := 0; i < b.N; i++ {
		n2, _ := WriteTag(&buf, s, a, w2)
		n += n2
	}
}

func BenchmarkWriterTag3(b *testing.B) {
	buf := bytes.Buffer{}
	s := "tag"
	var n int
	a := Attributes{"a": "b"}
	w2 := makeWritersTo(strings.NewReader("abc"+s), strings.NewReader("cd"))
	for i := 0; i < b.N; i++ {
		n2, _ := WriteTag(&buf, s, a, w2)
		n += n2
	}
}

func BenchmarkRenderTag(b *testing.B) {
	s := "tag"
	inner := "abc"
	var s2 string
	a := Attributes{"a": "b"}
	for i := 0; i < b.N; i++ {
		s2 = RenderTag(s, a, inner)
		_ = s2
	}
}

func Test_writeTag(t *testing.T) {
	type args struct {
		tag       string
		attr      Attributes
		innerHtml io.WriterTo
		isVoid    bool
		noSpace   bool
		format    bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"simple void tag", args{"a", nil, nil, true, false, false}, "<a>", false},
		{"void tag with attribute", args{"a", Attributes{"b": "c"}, nil, true, false, false}, `<a b="c">`, false},
		{"no space", args{"a", Attributes{"b": "c"}, strings.NewReader("d"), false, true, false}, `<a b="c">d</a>`, false},
		{"space", args{"a", Attributes{"b": "c"}, strings.NewReader("d"), false, false, false}, `<a b="c">` + "\n" + `d` + "\n" + `</a>`, false},
		{"format", args{"a", Attributes{"b": "c"}, strings.NewReader("d"), false, false, true}, `<a b="c">` + "\n" + `  d` + "\n" + `</a>`, false},
		{"format no space", args{"a", Attributes{"b": "c"}, strings.NewReader("d"), false, true, true}, `<a b="c">d</a>`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			gotN, err := writeTag(w, tt.args.tag, tt.args.attr, tt.args.innerHtml, tt.args.isVoid, tt.args.noSpace, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeTag() gotW = %v, want %v", gotW, tt.wantW)
			}
			if gotN != len(tt.wantW) {
				t.Errorf("writeTag() gotN = %v, want %v", gotN, len(tt.wantW))
			}
		})
	}
}

func Test_writeTagErr(t *testing.T) {
	type args struct {
		tag       string
		attr      Attributes
		innerHtml io.WriterTo
		isVoid    bool
		noSpace   bool
		format    bool
	}
	tests := []struct {
		name string
		args args
		n    int
	}{
		{"void tag 0", args{"a", nil, nil, true, false, false}, 0},
		{"void tag 1", args{"a", nil, nil, true, false, false}, 1},
		{"void tag 2", args{"a", nil, nil, true, false, false}, 2},
		{"void tag 3", args{"ab", nil, nil, true, false, false}, 3},
		{"void tag attr 2", args{"a", Attributes{"b": "c"}, nil, true, false, false}, 2},
		{"void tag attr 3", args{"a", Attributes{"b": "c"}, nil, true, false, false}, 3},
		{"tag 3", args{"a", nil, strings.NewReader("abc"), false, false, false}, 3},
		{"tag 4", args{"a", nil, strings.NewReader("abc"), false, false, false}, 4},
		{"tag 5", args{"a", nil, strings.NewReader("abc"), false, false, false}, 5},
		{"tag 5.2", args{"a", nil, strings.NewReader("b"), false, false, false}, 5},
		{"tag 7", args{"a", nil, strings.NewReader("b"), false, false, false}, 7},
		{"tag 8", args{"a", nil, strings.NewReader("b"), false, false, false}, 8},
		{"tag 9", args{"a", nil, strings.NewReader("b"), false, false, false}, 9},
		{"tag attr 3", args{"a", Attributes{"b": "c"}, nil, false, false, false}, 3},
		{"tag attr 5", args{"a", Attributes{"b": "c"}, nil, false, false, false}, 5},
		{"tag attr 7", args{"a", Attributes{"b": "c"}, nil, false, false, false}, 7},
		{"tag attr formatted 5", args{"a", Attributes{"b": "c"}, nil, false, false, true}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newErrBuf(tt.n)
			gotN, err := writeTag(w, tt.args.tag, tt.args.attr, tt.args.innerHtml, tt.args.isVoid, tt.args.noSpace, tt.args.format)
			if err == nil {
				t.Errorf("writeTagErr() want err, got no error")
			}
			if gotN != tt.n {
				t.Errorf("writeTag() gotN = %v, want %v", gotN, tt.n)
			}
		})
	}
}

func TestIndent1(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"simple", `a`, `  a`},
		{"wrapped newlines", "\na\n", "\n  a\n"},
		{"inside newlines", "a\nb\nc", "  a\n  b\n  c"},
		{"inside space and newlines", "a\n  b\nc", "  a\n    b\n  c"},
		{"textarea", `<textarea height="10">a
  b
    c</textarea>`, `<textarea height="10">a
  b
    c</textarea>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Indent(tt.s); got != tt.want {
				t.Errorf("Indent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderTagFormatted(t *testing.T) {
	type args struct {
		tag       string
		attr      Attributes
		innerHtml string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with innerHtml", args{"a", Attributes{"b": "c"}, "d"}, `<a b="c">` + "\n" + `  d` + "\n" + `</a>`},
		{"without innerHtml", args{"a", Attributes{"b": "c"}, ""}, `<a b="c"></a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderTagFormatted(tt.args.tag, tt.args.attr, tt.args.innerHtml); got != tt.want {
				t.Errorf("RenderTagFormatted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderTagNoSpaceFormatted(t *testing.T) {
	type args struct {
		tag       string
		attr      Attributes
		innerHtml string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with innerHtml", args{"a", Attributes{"b": "c"}, "d"}, `<a b="c">d</a>`},
		{"without innerHtml", args{"a", Attributes{"b": "c"}, ""}, `<a b="c"></a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderTagNoSpaceFormatted(tt.args.tag, tt.args.attr, tt.args.innerHtml); got != tt.want {
				t.Errorf("RenderTagNoSpaceFormatted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderImage(t *testing.T) {
	s := RenderImage("http://a/b.img", "alt", nil)
	if s[:4] != "<img" {
		t.Errorf("TestRenderImage tag not rendered")
	}
}

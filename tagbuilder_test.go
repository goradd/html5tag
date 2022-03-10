package html5tag

import "fmt"

func ExampleTagBuilder_Tag() {
	fmt.Println(NewTagBuilder().Tag("div"))
	// Output: <div></div>
}

func ExampleTagBuilder_Set() {
	fmt.Println(NewTagBuilder().Tag("div").Set("me", "you"))
	// Output: <div me="you"></div>
}

func ExampleTagBuilder_ID() {
	fmt.Println(NewTagBuilder().Tag("div").ID("bob"))
	// Output: <div id="bob"></div>
}

func ExampleTagBuilder_Class() {
	fmt.Println(NewTagBuilder().Tag("div").Class("bob sam"))
	// Output: <div class="bob sam"></div>
}

func ExampleTagBuilder_Link() {
	fmt.Println(NewTagBuilder().Link("http://example.com"))
	// Output: <a href="http://example.com"></a>
}

func ExampleTagBuilder_IsVoid() {
	fmt.Println(NewTagBuilder().Tag("img").IsVoid())
	// Output: <img>
}

func ExampleTagBuilder_InnerHtml() {
	fmt.Println(NewTagBuilder().Tag("div").InnerHtml("<p>A big deal</p>"))
	// Output:
	// <div>
	// <p>A big deal</p>
	// </div>
}

func ExampleTagBuilder_InnerText() {
	fmt.Println(NewTagBuilder().Tag("div").InnerText("<p>A big deal</p>"))
	// Output:
	// <div>
	// &lt;p&gt;A big deal&lt;/p&gt;
	// </div>
}

func ExampleTagBuilder_String() {
	s := NewTagBuilder().Tag("div").InnerHtml("<p>A big deal</p>").String()
	fmt.Println(s)
	// Output:
	// <div>
	// <p>A big deal</p>
	// </div>
}

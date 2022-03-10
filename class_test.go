package html5tag

import (
	"fmt"
	"strconv"
	"testing"
)

func ExampleMergeWords() {
	classes := MergeWords("myClass1 myClass2", "myClass1 myClass3")
	fmt.Println(classes)
	// Output: myClass1 myClass2 myClass3
}

func ExampleRemoveWords() {
	classes := RemoveWords("myClass1 myClass2", "myClass1 myClass3")
	fmt.Println(classes)
	// Output: myClass2
}

func ExampleHasWord() {
	found := HasWord("myClass31 myClass2", "myClass3")
	fmt.Println(strconv.FormatBool(found))
	// Output: false
}

func ExampleRemoveClassesWithPrefix() {
	classes := RemoveClassesWithPrefix("col-6 col-brk col4-other", "col-")
	fmt.Println(classes)
	// Output: col4-other
}

func TestMergeWords1(t *testing.T) {
	tests := []struct {
		name           string
		originalValues string
		newValues      string
		want           string
	}{
		{"same", "myClass1", "myClass1", "myClass1"},
		{"empty1", "", "myClass1", "myClass1"},
		{"empty2", "myClass1", "", "myClass1"},
		{"remove spaces", " myClass1  myClass2", "", "myClass1 myClass2"},
		{"no shuffle", "myClass1 myClass2", "myClass2 myClass1", "myClass1 myClass2"},
		{"append", "myClass1 myClass2", "myClass3", "myClass1 myClass2 myClass3"},
		{"append1", "myClass1 myClass2", "myClass3 myClass1", "myClass1 myClass2 myClass3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeWords(tt.originalValues, tt.newValues); got != tt.want {
				t.Errorf("MergeWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

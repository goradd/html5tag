package html5tag

import (
	"strings"
)

// Utilities to manage class strings

// MergeWords is a utility function that appends the given space separated words to the end
// of the given string, if the words are not already in the string. This is primarily used for
// adding classes to a class attribute, but other attributes work as well, like
// aria-labelledby and aria-describedby attributes.
//
// MergeWords returns the new string, which will have no duplicates.
//
// Since the order of a class list in html makes a difference, you should take care in the
// order of the classes you add if this matters in your situation.
func MergeWords(originalValues string, newValues string) string {
	var found bool

	wordArray := strings.Fields(originalValues)
	newWordArray := strings.Fields(newValues)
	for _, s := range newWordArray {
		found = false
		for _, s2 := range wordArray {
			if s2 == s {
				found = true
			}
		}
		if !found {
			wordArray = append(wordArray, s)
		}
	}
	return strings.Join(wordArray, " ")
}

// HasWord searches haystack for the given needle.
func HasWord(haystack string, needle string) (found bool) {
	classArray := strings.Fields(haystack)
	for _, s := range classArray {
		if s == needle {
			found = true
			break
		}
	}
	return
}

// RemoveWords removes a value from the list of space-separated values given.
// You can give it more than one value to remove by
// separating the values with spaces in the removeValue string. This is particularly useful
// for removing a class from a class list in a class attribute.
func RemoveWords(originalValues string, removeValues string) string {
	classes := strings.Fields(originalValues)
	removeClasses := strings.Fields(removeValues)
	ret := ""
	var found bool

	for _, s := range classes {
		found = false
		for _, s2 := range removeClasses {
			if s2 == s {
				found = true
			}
		}
		if !found {
			ret = ret + s + " "
		}
	}

	ret = strings.TrimSpace(ret)

	return ret
}

// RemoveClassesWithPrefix will remove all classes from the class string with the given prefix.
//
// Many CSS frameworks use families of classes, which are built up from a base family name. For example,
// Bootstrap uses 'col-lg-6' to represent a table that is 6 units wide on large screens and Foundation
// uses 'large-6' to do the same thing. This utility removes classes that start with a particular prefix
// to remove whatever sizing class was specified.
// Returns the resulting class list.
func RemoveClassesWithPrefix(class string, prefix string) string {
	classes := strings.Fields(class)
	ret := ""

	for _, s := range classes {
		if !strings.HasPrefix(s, prefix) {
			ret = ret + s + " "
		}
	}

	ret = strings.TrimSpace(ret)

	return ret
}

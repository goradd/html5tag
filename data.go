package html5tag

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ToDataAttr is a helper function to convert a name from camelCase to kabob-case for data attributes in particular.
//
// data-* html attributes have special conversion rules. Attribute names should always be lower case. Dashes in the
// name get converted to camel case javascript variable names.
// For example, if you want to pass the value with key name "testVar" to javascript by printing it in
// the html, you would use this function to help convert it to "data-test-var", after which you can retrieve
// in javascript by calling ".data('testVar')". on the object.
// This will also test for the existence of a camel case string it cannot handle
func ToDataAttr(s string) (string, error) {
	if matched, _ := regexp.MatchString("^[^a-z]|[A-Z][A-Z]|\\W", s); matched {
		err := fmt.Errorf("%s is not an acceptable camelCase name", s)
		return s, err
	}
	re, err := regexp.Compile("[A-Z]")
	if err == nil {
		s = re.ReplaceAllStringFunc(s, func(s2 string) string { return "-" + strings.ToLower(s2) })
	}

	return strings.TrimSpace(strings.TrimPrefix(s, "-")), err
}

// ToDataKey is a helper function to convert a name from kabob-case to camelCase.
//
// data-* html attributes have special conversion rules. Key names should always be lower case. Dashes in the
// name get converted to camel case javascript variable names.
// For example, if you want to pass the value with key name "testVar" to javascript by printing it in
//the html, you would use this function to help convert it to "data-test-var", after which you can retrieve
//in javascript by calling ".dataset.testVar" on the object.
func ToDataKey(s string) (string, error) {
	if matched, _ := regexp.MatchString("[A-Z]|[^a-z0-9-]", s); matched {
		err := errors.New("this is not an acceptable kabob-case name")
		return s, err
	}

	pieces := strings.Split(s, "-")
	var ret string
	for i, p := range pieces {
		if len(p) == 1 {
			err := errors.New("individual kabob words must be at least 2 characters long")
			return s, err
		}
		if i != 0 {
			p = strings.Title(p)
		}
		ret += p
	}
	return ret, nil
}

package html5tag

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const numericMatch = `-?[\d]*(\.[\d]+)?`

var numericReplacer, _ = regexp.Compile(numericMatch)
var numericMatcher, _ = regexp.Compile("^" + numericMatch + "$")

// keys for style attributes that take a number that is not a length
var nonLengthNumerics = map[string]bool{
	"volume":            true,
	"speech-rate":       true,
	"orphans":           true,
	"widows":            true,
	"pitch-range":       true,
	"font-weight":       true,
	"z-index":           true,
	"counter-increment": true,
	"counter-reset":     true,
}

// Style makes it easy to add and manipulate individual properties in a generated style sheet.
//
// Its main use is for generating a style attribute in an HTML tag.
// It implements the String interface to get the style properties as an HTML embeddable string.
type Style map[string]string

// NewStyle initializes an empty Style object.
func NewStyle() Style {
	return make(map[string]string)
}

// Copy copies the given style. It also turns a map[string]string into a Style.
func (s Style) Copy() Style {
	s2 := NewStyle()
	s2.Merge(s)
	return s2
}

// Merge merges the styles from one style to another. Conflicts will overwrite the current style.
func (s Style) Merge(m Style) {
	for k, v := range m {
		s[k] = v
	}
}

// Len returns the number of properties in the style.
func (s Style) Len() int {
	if s == nil {
		return 0
	}
	return len(s)
}

// Has returns true if the given property is in the style.
func (s Style) Has(property string) bool {
	if s == nil {
		return false
	}
	_, ok := s[property]
	return ok
}

// Get returns the property.
func (s Style) Get(property string) string {
	return s[property]
}

// Remove removes the property.
func (s Style) Remove(property string) {
	delete(s, property)
}

// SetString receives a style encoded "style" attribute into the Style structure (e.g. "width: 4px; border: 1px solid black")
func (s Style) SetString(text string) (changed bool, err error) {
	s.RemoveAll()
	a := strings.Split(text, ";") // break apart into pairs
	changed = false
	err = nil
	for _, value := range a {
		b := strings.Split(value, ":")
		if len(b) != 2 {
			err = errors.New("Css must be a name/value pair separated by a colon. '" + string(text) + "' was given.")
			return
		}
		newChange, newErr := s.SetChanged(strings.TrimSpace(b[0]), strings.TrimSpace(b[1]))
		if newErr != nil {
			err = newErr
			return
		}
		changed = changed || newChange
	}
	return
}

// SetChanged sets the given property to the given value.
//
// If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value
// For example, Set ("height", "* 2") will double the height value without changing the unit specifier
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.
func (s Style) SetChanged(property string, value string) (changed bool, err error) {
	if strings.Contains(property, " ") {
		err = errors.New("attribute names cannot contain spaces")
		return
	}

	if strings.HasPrefix(value, "+ ") ||
		strings.HasPrefix(value, "- ") || // the space here distinguishes between a math operation and a negative value
		strings.HasPrefix(value, "* ") ||
		strings.HasPrefix(value, "/ ") {

		return s.mathOp(property, value[0:1], value[2:])
	}

	if value == "0" {
		changed = s.set(property, value)
		return
	}

	isNumeric := numericMatcher.MatchString(value)
	if isNumeric {
		if !nonLengthNumerics[property] {
			value = value + "px"
		}
		changed = s.set(property, value)
		return
	}

	changed = s.set(property, value)
	return
}

// Set is like SetChanged, but returns the Style for chaining.
func (s Style) Set(property string, value string) Style {
	_, err := s.SetChanged(property, value)
	if err != nil {
		panic(err)
	}
	return s
}

// opReplacer is used in the regular expression replacement function below
func opReplacer(op string, v float64) func(string) string {
	return func(cur string) string {
		if cur == "" {
			return ""
		} // bug workaround
		//fmt.Println(cur)
		f, err := strconv.ParseFloat(cur, 0)
		if err != nil {
			panic("The number detector is broken on " + cur) // this is coming directly from the regular expression match
		}
		var newVal float64
		switch op {
		case "+":
			newVal = f + v
		case "-":
			newVal = f - v
		case "*":
			newVal = f * v
		case "/":
			newVal = f / v
		default:
			panic("Unexpected operation")
		}

		// floating point operations sometimes are not accurate. This is an attempt to correct epsilons.
		return fmt.Sprint(roundFloat(newVal, 6))

	}
}

// mathOp applies the given math operation and value to all the numeric values found in the given property.
// Bug(r) If the operation is working on a zero, and the result is not a zero, we may get a raw number with no unit. Not a big deal, but result will use default unit of browser, which is not always px
func (s Style) mathOp(property string, op string, val string) (changed bool, err error) {
	cur := s.Get(property)
	if cur == "" {
		cur = "0"
	}

	f, err := strconv.ParseFloat(val, 0)
	if err != nil {
		return
	}
	newStr := numericReplacer.ReplaceAllStringFunc(cur, opReplacer(op, f))
	changed = s.set(property, newStr)
	return
}

// RemoveAll resets the style to contain no styles
func (s Style) RemoveAll() {
	for k := range s {
		delete(s, k)
	}
}

// String returns the string version of the style attribute, suitable for inclusion in an HTML style tag
func (s Style) String() string {
	return s.encode()
}

// set is a raw set and return true if changed
func (s Style) set(k string, v string) bool {
	oldVal, existed := s[k]
	s[k] = v
	return !existed || oldVal != v
}

// roundFloat takes out rounding errors when doing length math
func roundFloat(f float64, digits int) float64 {
	f = f * math.Pow10(digits)
	if math.Abs(f) < 0.5 {
		return 0
	}
	v := int(f + math.Copysign(0.5, f))
	f = float64(v) / math.Pow10(digits)
	return f
}

// encode will output a text version of the style, suitable for inclusion in an HTML "style" attribute.
// it will sort the keys so that they are presented in a consistent and testable way.
func (s Style) encode() (text string) {
	var keys []string
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			text += ";"
		}
		text += k + ":" + s.Get(k)
	}
	return text
}

// StyleString converts an interface type that is being used to set a style value to a string that can be fed into
// the SetStyle* functions
func StyleString(i interface{}) string {
	var sValue string
	switch v := i.(type) {
	case int:
		sValue = fmt.Sprintf("%dpx", v)
	case float32:
		sValue = fmt.Sprintf("%gpx", v)
	case float64:
		sValue = fmt.Sprintf("%gpx", v)
	case string:
		sValue = v
	case fmt.Stringer:
		sValue = v.String()
	default:
		sValue = fmt.Sprint(v)
	}
	return sValue
}

// MergeStyleStrings merges the styles found in the two style strings.
// s2 wins conflicts.
func MergeStyleStrings(s1, s2 string) string {
	style1 := NewStyle()
	_, _ = style1.SetString(s1)
	style2 := NewStyle()
	_, _ = style2.SetString(s2)
	style1.Merge(style2)
	return style1.String()
}

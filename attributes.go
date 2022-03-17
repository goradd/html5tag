package html5tag

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FalseValue is use by Set to set a boolean attribute to false. The Has() function will return true, but
// the value will not appear in the attribute list when converted to a string.
const FalseValue = "**GORADD-FALSE**"

// Attributer is a general purpose interface for objects that return attributes based on information given.
type Attributer interface {
	Attributes(...interface{}) Attributes
}

// Attributes is an HTML attribute manager.
//
// Use Set to set specific attribute values,
// and then convert it to a string to get the attributes embeddable in an HTML tag.
//
// To create new attributes, the easiest is to do this:
//   a := Attributes{"id":"theId", "class":"myClass"}
type Attributes map[string]string

// NewAttributes creates a new Attributes collection.
func NewAttributes() Attributes {
	return make(map[string]string)
}

// Copy returns a copy of the attributes.
func (a Attributes) Copy() Attributes {
	a2 := Attributes{}
	return a2.Merge(a)
}

// Len returns the number of attributes.
func (a Attributes) Len() int {
	if a == nil {
		return 0
	}
	return len(a)
}

// Has returns true if the Attributes has the named attribute.
func (a Attributes) Has(attr string) bool {
	if a == nil {
		return false
	}
	_, ok := a[attr]
	return ok
}

// Get returns the named attribute.
func (a Attributes) Get(attr string) string {
	return a[attr]
}

// Remove deletes the given attribute.
func (a Attributes) Remove(attr string) {
	delete(a, attr)
}

// SetChanged sets the value of an attribute and returns changed if something in the attribute
// structure changed.
//
// It looks for special attributes like "class" and "style" to do some error checking
// on them. Returns err if the given attribute name or value is not valid.
//
// Use SetDataChanged when setting data attributes for additional validity checks.
func (a Attributes) SetChanged(name string, v string) (changed bool, err error) {
	if strings.Contains(name, " ") {
		err = errors.New("attribute names cannot contain spaces")
		return
	}

	if v == FalseValue {
		changed = a.RemoveAttribute(name)
		return
	}

	if name == "style" {
		styles := NewStyle()
		_, err = styles.SetString(v)
		if err != nil {
			return
		}

		oldStyles := a.StyleMap()

		if !reflect.DeepEqual(oldStyles, styles) { // since maps are not ordered, we must use a special equality test. We can't just compare strings for equality here.
			changed = true
			a["style"] = styles.String()
		}
		return
	}
	if name == "id" {
		return a.SetIDChanged(v)
	}
	if name == "class" {
		changed = a.SetClassChanged(v)
		return
	}
	if strings.HasPrefix(name, "data-") {
		return a.SetDataChanged(name[5:], v)
	}
	changed = a.set(name, v)
	return
}

// set is a raw set and return true if changed
func (a Attributes) set(k string, v string) bool {
	oldVal, existed := a[k]
	a[k] = v
	return !existed || oldVal != v
}

// Set sets a particular attribute and returns Attributes so that it can be chained.
//
// It looks for special attributes like "class", "style" and "data" to do some error checking
// on them. Use SetData to set data attributes.
//
// Pass v an empty string to create a boolean TRUE attribute, or to FalseValue to set the attribute
// such that you know it has been set, but will not print in the final html string.
func (a Attributes) Set(name string, v string) Attributes {
	_, err := a.SetChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// RemoveAttribute removes the named attribute.
// Returns true if the attribute existed.
func (a Attributes) RemoveAttribute(name string) bool {
	if a == nil {
		return false
	}
	if a.Has(name) {
		a.Remove(name)
		return true
	}
	return false
}

// This is a helper to sort the attribute keys so that special attributes
// are returned in a consistent order
var attrSpecialSort = map[string]int{
	"id":    1,
	"class": 2,
	"style": 3,
	// keep name and value together
	"name":  4,
	"value": 5,
	// keep src and alt together
	"src": 6,
	"alt": 7,
	// keep width and height together
	"width":  8,
	"height": 9,
}

func (a Attributes) sortedKeys() []string {
	keys := make([]string, len(a), len(a))
	idx := 0
	for k := range a {
		keys[idx] = k
		idx++
	}
	sort.Slice(keys, func(i1, i2 int) bool {
		k1 := keys[i1]
		k2 := keys[i2]
		v1, ok1 := attrSpecialSort[k1]
		v2, ok2 := attrSpecialSort[k2]
		if ok1 {
			if ok2 {
				return v1 < v2
			}
			return true
		} else if ok2 { // and !ok1
			return false
		} else { // !ok1 && !ok2
			return k1 < k2
		}
	})
	return keys
}

// String returns the attributes escaped and encoded, ready to be placed in an HTML tag
func (a Attributes) String() string {
	if a == nil {
		return ""
	}
	b := strings.Builder{}
	_, _ = a.WriteTo(&b)
	return b.String()
}

// SortedString returns the attributes escaped and encoded, ready to be placed in an HTML tag
// For consistency, it will use attrSpecialSort to order the keys.
func (a Attributes) SortedString() string {
	if a == nil {
		return ""
	}
	b := strings.Builder{}
	_, err := a.WriteSortedTo(&b)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func writeKV(w io.Writer, k, v string) (n int, err error) {
	if v == "" {
		if n, err = writeString(w, k, n); err != nil {
			return
		}
	} else {
		v = html.EscapeString(v)
		if n, err = writeString(w, k, n); err != nil {
			return
		}
		if n, err = writeString(w, `="`, n); err != nil {
			return
		}
		if n, err = writeString(w, v, n); err != nil {
			return
		}
		if n, err = writeString(w, `"`, n); err != nil {
			return
		}
	}
	return
}

// WriteSortedTo writes the attributes escaped, encoded and with sorted keys.
func (a Attributes) WriteSortedTo(w io.Writer) (n int64, err error) {
	if a == nil {
		return
	}
	var n1 int

	sk := a.sortedKeys()
	lastKey := len(sk) - 1
	for i, k := range sk {
		v := a[k]
		n1, err = writeKV(w, k, v)
		n += int64(n1)
		if err != nil {
			return
		}
		if i < lastKey {
			n1, err = io.WriteString(w, " ")
			n += int64(n1)
			if err != nil {
				return
			}
		}
	}
	return
}

// WriteTo writes the attributes escaped and encoded as fast as possible.
func (a Attributes) WriteTo(w io.Writer) (n int64, err error) {
	if a == nil {
		return
	}
	var n1 int
	i := 1
	length := len(a)
	for k, v := range a {
		n1, err = writeKV(w, k, v)
		n += int64(n1)
		if err != nil {
			return
		}
		if i < length {
			n1, err = io.WriteString(w, " ")
			n += int64(n1)
			if err != nil {
				return
			}

		}
		i++
	}
	return
}

// Range will call f for each item in the attributes.
//
// Keys will be ranged over such that repeating the range will produce the same ordering of keys.
// Return true from the range function to continue iterating, or false to stop.
func (a Attributes) Range(f func(key string, value string) bool) {
	if a == nil {
		return
	}
	for _, k := range a.sortedKeys() {
		if !f(k, a[k]) {
			break
		}
	}
}

// Override will replace attributes with the attributes in overrides.
// Conflicts are won by the given overrides.
func (a Attributes) Override(overrides Attributes) Attributes {
	if overrides == nil {
		return a
	}
	for k, v := range overrides {
		a[k] = v
	}
	return a
}

// Merge merges the given attributes into the current attributes. Conflicts are generally won by the passed in Attributes.
// However, styles are merged, so that if both the passed in map and the current map have a styles attribute, the
// actual style properties will get merged together. Style conflicts are won by the passed in map.
// The class attribute will merge so that the final classes will be a union of the two.
//
// See Override for a merge that does not merge the styles or classes.
func (a Attributes) Merge(aIn Attributes) Attributes {
	if aIn == nil {
		return a
	}
	for k, v := range aIn {
		if k == "style" {
			if v2, ok := a[k]; ok {
				v = MergeStyleStrings(v2, v)
			}
		} else if k == "class" {
			if v2, ok := a[k]; ok {
				v = MergeWords(v2, v)
			}
		}
		a[k] = v
	}
	return a
}

// OverrideString merges an attribute string into the attributes. Conflicts are won by the string.
//
// It takes an attribute string of the form
//   a="b" c="d"
func (a Attributes) OverrideString(s string) Attributes {
	if s == "" {
		return a
	}
	a2 := getAttributesFromTemplate(s)
	a.Override(a2)
	return a
}

// MergeString merges an attribute string into the attributes.
// Conflicts are won by the string, but styles and classes merge.
//
// It takes an attribute string of the form
//   a="b" c="d"
func (a Attributes) MergeString(s string) Attributes {
	if s == "" {
		return a
	}
	a2 := getAttributesFromTemplate(s)
	a.Merge(a2)
	return a
}

// SetIDChanged sets the id to the given value and returns true if something changed.
// In other words, if you set the id to the same value that it currently is, it will return false.
// It will return an error if you attempt to set the id to an illegal value.
func (a Attributes) SetIDChanged(i string) (changed bool, err error) {
	if i == "" { // empty attribute is not allowed, so it is the same as removal
		changed = a.RemoveAttribute("id")
		return
	}

	if strings.ContainsAny(i, " ") {
		err = errors.New("id attributes cannot contain spaces")
		return
	}

	changed = a.set("id", i)
	return
}

// SetID sets the id attribute to the given value
func (a Attributes) SetID(i string) Attributes {
	_, err := a.SetIDChanged(i)
	if err != nil {
		panic(err)
	}
	return a
}

// ID returns the value of the id attribute.
func (a Attributes) ID() string {
	if a == nil {
		return ""
	}
	return a.Get("id")
}

// SetClassChanged sets the class attribute to the value given.
//
// If you prefix the value with "+ " the given value will be appended to the end of the current class list.
// If you prefix the value with "- " the given value will be removed from a class list.
// Otherwise, the current class value is replaced.
// Returns whether something actually changed or not.
// value can be multiple classes separated by a space
func (a Attributes) SetClassChanged(value string) bool {
	if value == "" { // empty attribute is not allowed, so it is the same as removal
		return a.RemoveAttribute("class")
	}

	if strings.HasPrefix(value, "+ ") {
		return a.AddClassChanged(value[2:])
	} else if strings.HasPrefix(value, "- ") {
		return a.RemoveClass(value[2:])
	}

	changed := a.set("class", value)
	return changed
}

// SetClass will set the class to the given value, and return the attributes so that you can chain calls.
func (a Attributes) SetClass(v string) Attributes {
	a.SetClassChanged(v)
	return a
}

// RemoveClass removes the named class from the list of classes in the class attribute.
//
// Returns true if the attribute changed.
func (a Attributes) RemoveClass(v string) bool {
	if a.Has("class") {
		oldClass := a.Get("class")
		newClass := RemoveWords(oldClass, v)
		if oldClass != newClass {
			a.set("class", newClass)
			return true
		}
		return false
	}
	return false
}

// RemoveClassesWithPrefix removes classes with the given prefix.
//
// Many CSS frameworks use families of classes, which are built up from a base family name. For example,
// Bootstrap uses 'col-lg-6' to represent a table that is 6 units wide on large screens and Foundation
// uses 'large-6' to do the same thing. This utility removes classes that start with a particular prefix
// to remove whatever sizing class was specified.
// Returns true if the list actually changed.
func (a Attributes) RemoveClassesWithPrefix(v string) bool {
	if a.Has("class") {
		oldClass := a.Get("class")
		newClass := RemoveClassesWithPrefix(oldClass, v)
		if oldClass != newClass {
			a.set("class", newClass)
			return true
		}
		return false
	}
	return false
}

// AddValuesChanged adds the given space separated values to the end of the values in the
// given attribute, removing duplicates and returning true if the attribute was changed at all.
// An example of a place to use this is the aria-labelledby attribute, which can take multiple
// space-separated id numbers.
func (a Attributes) AddValuesChanged(attrKey string, values string) bool {
	if values == "" {
		return false // nothing to add
	}
	if a.Has(attrKey) {
		attrValue := a.Get(attrKey)
		newValues := MergeWords(attrValue, values)
		if newValues != attrValue {
			a.set(attrKey, newValues)
			return true
		}
		return false
	}

	a.set(attrKey, values)
	return true
}

// AddValues adds space separated values to the end of an attribute value.
// If a value is not present, the value will be added to the end of the value list.
// If a value is present, it will not be added, and the position of the current value in the list will not change.
func (a Attributes) AddValues(attr string, values string) Attributes {
	a.AddValuesChanged(attr, values)
	return a
}

// AddClassChanged is similar to AddClass, but will return true if the class changed at all.
func (a Attributes) AddClassChanged(v string) bool {
	return a.AddValuesChanged("class", v)
}

// AddClass adds a class or classes. Multiple classes can be separated by spaces.
// If a class is not present, the class will be added to the end of the class list.
// If a class is present, it will not be added, and the position of the current class in the list will not change.
func (a Attributes) AddClass(v string) Attributes {
	a.AddClassChanged(v)
	return a
}

// Class returns the value of the class attribute.
func (a Attributes) Class() string {
	return a.Get("class")
}

// HasAttributeValue returns true if the given value exists in the space-separated attribute value.
func (a Attributes) HasAttributeValue(attr string, value string) bool {
	var curValue string
	if curValue = a.Get(attr); curValue == "" {
		return false
	}
	f := strings.Fields(curValue)
	for _, s := range f {
		if s == value {
			return true
		}
	}
	return false
}

// HasClass returns true if the given class is in the class list in the class attribute.
func (a Attributes) HasClass(c string) bool {
	return a.HasAttributeValue("class", c)
}

// SetDataChanged sets the given value as an HTML "data-*" attribute.
// The named value will be retrievable in javascript by using
//
//	$obj.dataset.valname;
//
// Note: Data name cases are handled specially. data-* attribute names are supposed to be lower kebab case. Javascript
// converts dashed notation to camelCase when converting html attributes into object properties.
// In other words, we give it a camelCase name here, it shows up in the html as
// a kebab-case name, and then you retrieve it using javascript as camelCase again.
//
// For example, if your html looks like this:
//
//	<div id='test1' data-test-case="my test"></div>
//
// You would get that value in javascript by doing:
//	g$('test1').data('testCase');
//
// Conversion to special html data-* name formatting is handled here automatically. So if you SetData('testCase') here,
// you can get it using .dataset.testCase in javascript
func (a Attributes) SetDataChanged(name string, v string) (changed bool, err error) {
	// validate the name
	if strings.ContainsAny(name, " !$") {
		err = errors.New("data attribute names cannot contain spaces or $ or ! chars")
		return
	}
	suffix, err := ToDataAttr(name)
	if err == nil {
		name = "data-" + suffix
		changed = a.set(name, v)
	}
	return
}

// SetData sets the given data attribute. Data attribute keys must be in camelCase notation and
// cannot be hyphenated. The key will get converted to kebab-case for output in html. When referring to
// the attribute in javascript, javascript will convert it back into camelCase.
func (a Attributes) SetData(name string, v string) Attributes {
	_, err := a.SetDataChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// DataAttribute gets the data attribute value that was set previously. The key should be in camelCase.
func (a Attributes) DataAttribute(key string) string {
	if a == nil {
		return ""
	}
	suffix, _ := ToDataAttr(key)
	key = "data-" + suffix
	return a.Get(key)
}

// RemoveDataAttribute removes the named data attribute. The key should be in camelCase.
// Returns true if the data attribute existed.
func (a Attributes) RemoveDataAttribute(key string) bool {
	if a == nil {
		return false
	}
	suffix, _ := ToDataAttr(key)
	key = "data-" + suffix
	return a.RemoveAttribute(key)
}

// HasDataAttribute returns true if the data attribute is set. The key should be in camelCase.
func (a Attributes) HasDataAttribute(key string) bool {
	if a == nil {
		return false
	}
	suffix, _ := ToDataAttr(key)
	key = "data-" + suffix
	return a.Has(key)
}

// StyleString returns the css style string, or a blank string if there is none.
func (a Attributes) StyleString() string {
	return a.Get("style")
}

// StyleMap returns a special Style structure which lets you refer to the styles as a string map.
func (a Attributes) StyleMap() Style {
	s := NewStyle()
	_, _ = s.SetString(a.StyleString())
	return s
}

// SetStyleChanged sets the given style to the given value. If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value.
// For example, SetStyle ("height", "* 2") will double the height value without changing the unit specifier.
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.
func (a Attributes) SetStyleChanged(name string, v string) (changed bool, err error) {
	s := a.StyleMap()
	changed, err = s.SetChanged(name, v)
	if err == nil {
		a.set("style", s.String())
	}
	return
}

// SetStyle sets the given property in the style attribute
func (a Attributes) SetStyle(name string, v string) Attributes {
	_, err := a.SetStyleChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// SetStyles merges the given styles with the current styles. The given style wins on collision.
func (a Attributes) SetStyles(s Style) Attributes {
	styles := a.StyleMap()
	styles.Merge(s)
	a.set("style", styles.String())
	return a
}

// SetStylesTo sets the styles using a traditional css style string with colon and semicolon separators
func (a Attributes) SetStylesTo(s string) Attributes {
	styles := a.StyleMap()
	if _, err := styles.SetString(s); err != nil {
		return a
	}
	a.set("style", styles.String())
	return a
}

// GetStyle gives you the value of a single style attribute value. If you want all the attributes as a style string, use
// StyleString().
func (a Attributes) GetStyle(name string) string {
	if a == nil {
		return ""
	}
	s := a.StyleMap()
	return s.Get(name)
}

// HasStyle returns true if the given style is set to any value, and false if not.
func (a Attributes) HasStyle(name string) bool {
	if a == nil {
		return false
	}
	s := a.StyleMap()
	return s.Has(name)
}

// RemoveStyle removes the style from the style list. Returns true if there was a change.
func (a Attributes) RemoveStyle(name string) (changed bool) {
	if a == nil {
		return false
	}
	s := a.StyleMap()
	if s.Has(name) {
		changed = true
		s.Remove(name)
		a.set("style", s.String())
	}
	return changed
}

// SetDisabled sets the "disabled" attribute to the given value.
func (a Attributes) SetDisabled(d bool) Attributes {
	if d {
		a.Set("disabled", "")
	} else {
		a.RemoveAttribute("disabled")
	}
	return a
}

// IsDisabled returns true if the "disabled" attribute is set to true.
func (a Attributes) IsDisabled() bool {
	if a == nil {
		return false
	}
	return a.Has("disabled")
}

// SetDisplay sets the "display" attribute to the given value.
func (a Attributes) SetDisplay(d string) Attributes {
	a.SetStyle("display", d)
	return a
}

// IsDisplayed returns true if the "display" attribute is not set, or if it is set, if it is not set to "none".
func (a Attributes) IsDisplayed() bool {
	if a == nil {
		return true
	}
	return a.GetStyle("display") != "none"
}

// ValueString is a helper function to convert an interface type to a string that is appropriate for the value
// in the Set function.
func ValueString(i interface{}) string {
	switch v := i.(type) {
	case fmt.Stringer:
		return v.String()
	case bool:
		if v {
			return "" // boolean true
		}
		return FalseValue // Our special value to indicate to NOT print the attribute at all
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	}
	return fmt.Sprint(i)
}

// getAttributesFromTemplate returns Attributes extracted from a string in the form
// of name="value"
func getAttributesFromTemplate(s string) Attributes {
	pairs := templateMatcher.FindAllString(s, -1)
	if len(pairs) == 0 {
		return nil
	}
	a := NewAttributes()
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		val := kv[1][1 : len(kv[1])-1] // remove quotes
		a.Set(kv[0], val)
	}
	return a
}

/*
type AttributeCreator map[string]string

func (c AttributeCreator) Create() Attributes {
	return Attributes(c)
}
*/
var templateMatcher *regexp.Regexp

func init() {
	gob.Register(Attributes{})
	templateMatcher = regexp.MustCompile(`\w+=".*?"`)
}

/*
Package html5tag includes functions for manipulating html 5 formatted tags.
It includes specific functions for manipulating attributes inside of tags, including various
special attributes like styles, classes, and data-* attributes.

Many of the routines return a boolean to indicate whether the data actually changed. This can be used to prevent
needlessly redrawing html after setting values that had no effect on the attribute list.

You can choose to build tags using strings for convenience, or io.Writer for speed.
*/
package html5tag

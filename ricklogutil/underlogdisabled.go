//go:build !underlog
// +build !underlog

package ricklogutil

var Underlog = func(...interface{}) {}

// My IDE and associated tools seem to insist on adding
// //go:build !underlog, even though I am using Go 1.12, and
// //go:build doesn't yet exist.  I guess the Go command line
// tools I have installed on my computer are 1.25, and that's
// what the IDE and the linter and whatever are going by.

//go:build !underlog

package abelogutil

var UnderEnabled = false

func UnderPrintf(format string, arguments ...any) {
	// Do nothing.
}

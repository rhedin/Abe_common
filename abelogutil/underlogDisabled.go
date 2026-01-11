//go:build !underlog

package abelogutil

const UnderEnabled = false

func UnderPrintf(format string, arguments ...any) {
	// Do nothing.
}

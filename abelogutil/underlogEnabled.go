//go:build underlog

package abelogutil

import (
	"fmt"
	"runtime"
)

var UnderEnabled = true

func UnderPrintf(format string, arguments ...any) {
	fmt.Printf(format, arguments...)

	pc := make([]uintptr, 5)    // Slice to hold up to 5 program counters
	n := runtime.Callers(2, pc) // Skip 2 frames
	if n == 0 {
		fmt.Printf("The call to runtime.Callers returned 0.\n")
		return
	}
	pc = pc[:n] // Trim the slice to actual number of PCs

	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	fmt.Printf("- %s:%d %s\n", frame.File, frame.Line, frame.Function)
}

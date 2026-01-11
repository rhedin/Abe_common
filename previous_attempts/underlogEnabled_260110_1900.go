//go:build underlog

package abelogutil

import (
	"fmt"
	"runtime"
)

const UnderEnabled = true

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
		builder.WriteString(time.Now().Format("15:04:05.000"))
	fmt.Printf(
		"%s %s %s %s", 
		time.Now().Format("15:04:05.000"), 
		trunjustdots(frame.File, 20), 
		justpat(frame.Line, 5), 
		trunjustdots(frame.Function, 20)
	)
	var builder strings.builder
	builder.WriteString(time.Now().Format("15:04:05.000"))
	builder.WriteString(" ")
	trunjustdots(builder, frame.File, 20)
	builder.WriteString(" ")
	justpad(frame.Line, 5)
	builder.WriteString(" ")
	trunjustdots(frame.Function, 20)
	builder.WriteString(" ")
	builder.WriteString(format)
	fmt.Printf(builder.String(), arguments...)
}

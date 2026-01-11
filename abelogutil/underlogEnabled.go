//go:build underlog

package abelogutil

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const UnderEnabled = true

// This function appends a string to the passed-in builder.String.
// It makes sure the string is exactly n characters wide,
// justified to the right, padded with spaces on the left,
// and preceded by three dots if necessary.
//
// dog -> ^^^^^^^^^^^^^^^^^dog (hats are spaces)
// supercalifragilisticexpialidocious -> ...listicexpialidocious
func trunjustdots(builder *strings.Builder, input string, length int) {
	if len(input) > length {
		builder.WriteString("...")
		builder.WriteString(input[len(input)-length:])
	} else if len(input) > length {
		builder.WriteString("   ")
		builder.WriteString(input)
	} else {
		builder.WriteString("   ")
		builder.WriteString(strings.Repeat(" ", length-len(input)))
		builder.WriteString(input)
	}
}

// This function handles numbers.  We want numbers to be right-
// justified in a field, but if the field is too small, it expands.
func justpad(builder *strings.Builder, inputnum int, length int) {
	var inputstr = strconv.Itoa(inputnum)
	if len(inputstr) >= length {
		builder.WriteString(inputstr)
	} else {
		builder.WriteString(strings.Repeat(" ", length-len(inputstr)))
		builder.WriteString(inputstr)
	}
}

func UnderPrintf(format string, arguments ...any) {
	pc := make([]uintptr, 5)    // Slice to hold up to 5 program counters
	n := runtime.Callers(2, pc) // Skip 2 frames
	if n == 0 {
		fmt.Printf("The call to runtime.Callers returned 0.\n")
		return
	}
	pc = pc[:n] // Trim the slice to actual number of PCs
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	var builder strings.Builder
	builder.WriteString(time.Now().Format("15:04:05.000"))
	builder.WriteString(" ")
	trunjustdots(&builder, frame.File, 20)
	builder.WriteString(" ")
	justpad(&builder, frame.Line, 4)
	builder.WriteString(" ")
	trunjustdots(&builder, frame.Function, 20)
	builder.WriteString(" ")
	builder.WriteString(format)
	fmt.Printf(builder.String(), arguments...)
}

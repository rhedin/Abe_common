//go:build !underlog

package abelogutil

const UnderEnabled = false

func UnderPrintf(format string, arguments ...any) {
	// Do nothing.
}

func FormatQuery(requestString string) string {
	return ""
}

// Here's what Grok said about this.
// Go's dead-code elimination (which can remove the whole if false { ... } block)
// happens after type checking. The symbol must be resolvable first.
// So yes, the entire block in rest.go is removed during compilation.
// if abelog.UnderEnabled {
//     . . . .
//     formattedQuery := abelog.FormatQuery(string(requestDump))
//     . . . .
// }
// But the compiler wants to understand it. (first phase)  Before it removes it. (second phase)

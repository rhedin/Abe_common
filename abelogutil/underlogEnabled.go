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
	} else if len(input) == length {
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
	whoCalledUs, _ := frames.Next()
	whoCalledThem, more := frames.Next()
	if !more {
		fmt.Printf("We do not have information about the caller of our caller.\n")
		fmt.Printf("Leaving UnderPrintf early.\n")
	}
	var builder strings.Builder
	builder.WriteString(time.Now().Format("15:04:05.000"))
	builder.WriteString(" ")
	trunjustdots(&builder, whoCalledUs.File, 40)
	builder.WriteString(" ")
	justpad(&builder, whoCalledUs.Line, 4)
	builder.WriteString(" ")
	trunjustdots(&builder, whoCalledUs.Function, 40)
	builder.WriteString(" From:")
	builder.WriteString(" ")
	trunjustdots(&builder, whoCalledThem.File, 40)
	builder.WriteString(" ")
	justpad(&builder, whoCalledThem.Line, 4)
	builder.WriteString(" ")
	trunjustdots(&builder, whoCalledThem.Function, 40)
	builder.WriteString(" ")
	builder.WriteString(format)
	fmt.Printf(builder.String(), arguments...)
}

// We put in something like this:
// POST /db/v1/graphql/main HTTP/1.1
// Host: localhost:9090
// Accept: */*
// Accept-Encoding: gzip, deflate, br, zstd
// Accept-Language: en-US,en;q=0.9
// Connection: keep-alive
// Content-Length: 183
// Content-Type: application/json
// Cookie: _ga=GA1.1.899125396.1755970306; _ga_7W0YET4Q10=GS2.1.s1756308492$o3$g1$t1756308741$j60$l0$h0; ~sid=S-4133a929aa6ced4888449bf22ad503255ae6228f33fdecd6ce29f6aa22dfc634; ~aid=A-2cf749e015bbdc15fee5ceb41847f8a9ba7199d5df7804494220f82475bb07a1
// Origin: https://localhost:9090
// Referer: https://localhost:9090/
// Sec-Ch-Ua: "Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"
// Sec-Ch-Ua-Mobile: ?0
// Sec-Ch-Ua-Platform: "macOS"
// Sec-Fetch-Dest: empty
// Sec-Fetch-Mode: cors
// Sec-Fetch-Site: same-origin
// User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36
//
// {"operationName":"","variables":{"node":{"key":"1768665757868","kind":"general","message":"Hello."}},"query":"\nmutation($node : NodeTemplate) {\n  general(storeNode : $node) { }\n}"}
//
// And we return just the last line, reformatted like this:
//
//	{
//		"operationName":"",
//		"variables":{
//			"node":{
//				"key":"1768665757868",
//				"kind":"general",
//				"message":"Hello."
//			}
//		},
//		query:mutation($node:NodeTemplate){
//			general(storeNode:$node){}
//		}
//	}
func FormatQuery(requestString string) string {

	// Pick off the last line.

	var indexOfLastNewline = strings.LastIndex(requestString, "\n")
	if indexOfLastNewline == -1 {
		return "I can't find the last line."
	}
	var lastLine = requestString[indexOfLastNewline+1:]
	fmt.Printf("lastLine = -->%s<--\n", lastLine)

	if !strings.HasPrefix(lastLine, `{"operationName":`) {
		fmt.Printf("underlogEnabled.go FormatQuery: I can't find the operationName field.\n")
	}

	//       That's because there is no operationName field.  Should there be one?
	//
	// 21:18:32.065 ...edin/work/260102/Abe_eliasdb/api/rest.go  157 ...liasdb/api.RegisterRestEndpoints.func1.1 From: ...v/versions/1.25.4/src/net/http/server.go 2322              net/http.HandlerFunc.ServeHTTP Running a handler function for handlerURL = /db/v1/graphql/ and handlerInst = 0x104a32d30 with http.Request:
	// POST /db/v1/graphql/main HTTP/1.1
	// Host: localhost:9090
	// Accept: */*
	// Content-Length: 49
	// Content-Type: application/json
	// User-Agent: curl/8.7.1
	//
	// { "query": "{ Person(key: \"hans\") { name } }" }
	// formatted the graphql query:
	// I can't find the operationName field.
	//
	//       I asked Grok about it.  You don't have to name your operations.  The operationName
	//       field doesn't have to be there.
	//
	// No, you should not be concerned that your curl requests (or Apollo Sandbox queries)
	// lack an operationName field.
	//
	// Why operationName is optional in GraphQL
	//
	// The operationName field in a GraphQL request payload is optional in the spec and
	// in nearly all implementations (including graphql-go, gqlgen, Apollo Server, and
	// EliasDB's custom runtime).
	//
	// It is only required in two rare cases:
	//
	//   1. When the request contains multiple operations (e.g. multiple named queries/mutations
	//      in the same string) — the server needs to know which one to execute.
	//   2. When the client wants to explicitly name the operation for logging, caching,
	//      or debugging purposes.
	//
	// In your case (and in 99% of real-world usage):
	//
	//   1. Requests have exactly one operation.
	//   2. The query/mutation is anonymous (no query MyQuery { ... } or mutation MyMutation { ... } name).
	//
	// 	The server knows exactly which operation to run — no ambiguity.
	//
	//       So I changed the code so that lack of an operationName field won't make
	//       formatting bail out.

	// Compress the "query":"..."} part of the line.

	var queryIndex = strings.Index(lastLine, "\"query\":")
	if queryIndex == -1 {
		return "I can't find the query field."
	}
	var leftPortion = lastLine[:queryIndex]
	var rightPortion = lastLine[queryIndex:]
	var replacer = strings.NewReplacer(
		`\"`, `"`,
		`"`, ``,
		`\n`, ``,
		` `, ``,
	)
	rightPortion = replacer.Replace(rightPortion)
	lastLine = leftPortion + rightPortion
	fmt.Printf("Revised lastLine = -->%s<--\n", lastLine)
	// We might discover that removing all the quote marks from the query portion is
	// too root and branch.  Change the code when we discover problems.
	//
	// I did discover that removing all the quote marks was too "root and branch".
	// See [Note 1] below this function.

	// Format for easy reading.

	var indent = 0
	var inQuote = false
	var builder strings.Builder
	builder.Grow(500)

	var carriageReturn = func() {
		var spaces = strings.Repeat(" ", 4*indent)
		builder.WriteByte('\n')
		builder.WriteString(spaces)
	}

	for i := 0; i < len(lastLine); i++ {
		switch lastLine[i] {
		case '{':
			if inQuote {
				builder.WriteByte('{')
			} else if i+1 < len(lastLine) && lastLine[i+1] == '}' {
				// If an open brace is immediately followed by a close brace, we don't want to consume a line.
				builder.WriteString("{}")
				i++
			} else {
				builder.WriteByte('{')
				indent++
				carriageReturn()
			}
		case '}':
			if inQuote {
				builder.WriteByte('}')
			} else {
				if indent == 0 {
					carriageReturn()
					builder.WriteString("} Tried to indent a negative amount. ")
				} else {
					indent--
					carriageReturn()
					builder.WriteByte('}')
				}
			}
		case ',':
			if inQuote {
				builder.WriteByte(',')
			} else {
				builder.WriteByte(',')
				carriageReturn()
			}
		case '"':
			if inQuote {
				inQuote = false
				builder.WriteByte('"')
			} else {
				inQuote = true
				builder.WriteByte('"')
			}
		case '\\':
			if inQuote {
				if i+1 < len(lastLine) && lastLine[i+1] == '"' {
					builder.WriteString("\\\"")
					i++
				} else {
					builder.WriteByte('\\')
				}
			} else {
				builder.WriteByte('\\')
			}
		default:
			builder.WriteByte(lastLine[i])
		}
	}

	return builder.String()
}

/***

Here's what I found in my log file.

[cors] 2026/03/07 08:33:23 Handler: Actual request
[cors] 2026/03/07 08:33:23   Actual request no headers added: missing origin
lastLine = -->{ "query": "{ Person(key: \"hans\") { name } }" }<--
underlogEnabled.go FormatQuery: I can't find the operationName field.
Revised lastLine = -->{ query:{Person(key:\hans\){name}}}<--
08:33:23.708 ...edin/work/260102/Abe_eliasdb/api/rest.go  157 ...liasdb/api.RegisterRestEndpoints.func1.1 From: ...v/versions/1.25.4/src/net/http/server.go 2322              net/http.HandlerFunc.ServeHTTP Running a handler function for handlerURL = /db/v1/graphql/ and handlerInst = 0x100fa6dc0 with http.Request:
POST /db/v1/graphql/main HTTP/1.1
Host: localhost:9090
The original text was   Accept: asteriskslashasterisk   .  I changed it because it messed up our block comment.
Content-Length: 49
Content-Type: application/json
User-Agent: curl/8.7.1

{ "query": "{ Person(key: \"hans\") { name } }" }
formatted the graphql query:
{
     query:{
        Person(key:\hans\){
            name
        }
    }
}

I would rather convert \"s to "s.  I think.

Something Gai said:

If multiple patterns match at the same starting position (e.g., \" and "),
the longer match or the one appearing earlier in your NewReplacer arguments
takes precedence.

***/

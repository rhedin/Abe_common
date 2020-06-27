/*
 * Public Domain Software
 *
 * I (Matthias Ladkau) am the author of the source code in this file.
 * I have placed the source code in this file in the public domain.
 *
 * For further information see: http://creativecommons.org/publicdomain/zero/1.0/
 */

package parser

import (
	"fmt"
	"testing"
)

func TestStatementParsing(t *testing.T) {

	// Comment parsing without statements

	input := `a := 1
	b := 2; c:= 3`
	expectedOutput := `
statements
  :=
    identifier: a
    number: 1
  :=
    identifier: b
    number: 2
  :=
    identifier: c
    number: 3
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}
}

func TestIdentifierParsing(t *testing.T) {

	input := `a := 1
	a.foo := 2
	a.b.c.foo := a.b
	`
	expectedOutput := `
statements
  :=
    identifier: a
    number: 1
  :=
    identifier: a
      identifier: foo
    number: 2
  :=
    identifier: a
      identifier: b
        identifier: c
          identifier: foo
    identifier: a
      identifier: b
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	input = `a := b[1 + 1]
	a[4].foo["aaa"] := c[i]
	`
	expectedOutput = `
statements
  :=
    identifier: a
    identifier: b
      compaccess
        plus
          number: 1
          number: 1
  :=
    identifier: a
      compaccess
        number: 4
      identifier: foo
        compaccess
          string: 'aaa'
    identifier: c
      compaccess
        identifier: i
`[1:]
	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}
}

func TestCommentParsing(t *testing.T) {

	// Comment parsing without statements

	input := `/* This is  a comment*/ a := 1 + 1 # foo bar`
	expectedOutput := `
:=
  identifier: a #  This is  a comment
  plus
    number: 1
    number: 1 #  foo bar
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	input = `/* foo */ 1 # foo bar`
	expectedOutput = `
number: 1 #  foo   foo bar
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}
}

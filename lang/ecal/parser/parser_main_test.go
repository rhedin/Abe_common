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

	input = `a := 1
a().foo := x2.foo()
a.b.c().foo := a()
	`
	expectedOutput = `
statements
  :=
    identifier: a
    number: 1
  :=
    identifier: a
      funccall
      identifier: foo
    identifier: x2
      identifier: foo
        funccall
  :=
    identifier: a
      identifier: b
        identifier: c
          funccall
          identifier: foo
    identifier: a
      funccall
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	input = `a(1+2).foo := x2.foo(foo)
a.b.c(x()).foo := a(1,a(),3, x, y) + 1
	`
	expectedOutput = `
statements
  :=
    identifier: a
      funccall
        plus
          number: 1
          number: 2
      identifier: foo
    identifier: x2
      identifier: foo
        funccall
          identifier: foo
  :=
    identifier: a
      identifier: b
        identifier: c
          funccall
            identifier: x
              funccall
          identifier: foo
    plus
      identifier: a
        funccall
          number: 1
          identifier: a
            funccall
          number: 3
          identifier: x
          identifier: y
      number: 1
`[1:]

	if res, err := UnitTestParseWithPPResult("mytest", input, ""); err != nil || fmt.Sprint(res) != expectedOutput {
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

func TestSimpleExpressionParsing(t *testing.T) {

	// Test error output

	input := `"bl\*a"conversion`
	if _, err := UnitTestParse("mytest", input); err.Error() !=
		"Parse error in mytest: Lexical error (invalid syntax while parsing string) (Line:1 Pos:1)" {
		t.Error(err)
		return
	}

	// Test incomplete expression

	input = `a *`
	if _, err := UnitTestParse("mytest", input); err.Error() !=
		"Parse error in mytest: Unexpected end" {
		t.Error(err)
		return
	}

	input = `not ==`
	if _, err := UnitTestParse("mytest", input); err.Error() !=
		"Parse error in mytest: Term cannot start an expression (==) (Line:1 Pos:5)" {
		t.Error(err)
		return
	}

	input = `(==)`
	if _, err := UnitTestParse("mytest", input); err.Error() !=
		"Parse error in mytest: Term cannot start an expression (==) (Line:1 Pos:2)" {
		t.Error(err)
		return
	}

	input = "5 ( 5"
	if _, err := UnitTestParse("mytest", input); err.Error() !=
		"Parse error in mytest: Term can only start an expression (() (Line:1 Pos:3)" {
		t.Error(err)
		return
	}

	input = "5 + \""
	if _, err := UnitTestParse("mytest", input); err.Error() !=
		"Parse error in mytest: Lexical error (Unexpected end while reading string value (unclosed quotes)) (Line:1 Pos:5)" {
		t.Error(err)
		return
	}

	// Test prefix operator

	input = ` + a - -5`
	expectedOutput := `
minus
  plus
    identifier: a
  minus
    number: 5
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

}

func TestArithmeticParsing(t *testing.T) {
	input := "a + b * 5 /2"
	expectedOutput := `
plus
  identifier: a
  div
    times
      identifier: b
      number: 5
    number: 2
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	// Test brackets

	input = "a + 1 * (5 + 6)"
	expectedOutput = `
plus
  identifier: a
  times
    number: 1
    plus
      number: 5
      number: 6
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	// Test needless brackets

	input = "(a + 1) * (5 / (6 - 2))"
	expectedOutput = `
times
  plus
    identifier: a
    number: 1
  div
    number: 5
    minus
      number: 6
      number: 2
`[1:]

	// Pretty printer should get rid of the needless brackets

	res, err := UnitTestParseWithPPResult("mytest", input, "(a + 1) * 5 / (6 - 2)")
	if err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}
}

func TestLogicParsing(t *testing.T) {
	input := "not (a + 1) * 5 and tRue == false or not 1 - 5 != test"
	expectedOutput := `
or
  and
    not
      times
        plus
          identifier: a
          number: 1
        number: 5
    ==
      true
      false
  not
    !=
      minus
        number: 1
        number: 5
      identifier: test
`[1:]

	res, err := UnitTestParseWithPPResult("mytest", input, "not (a + 1) * 5 and true == false or not 1 - 5 != test")

	if err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	input = "a > b or a <= p or b hasSuffix 'test' or c hasPrefix 'test' and x < 4 or x >= 10"
	expectedOutput = `
or
  or
    or
      or
        >
          identifier: a
          identifier: b
        <=
          identifier: a
          identifier: p
      hassuffix
        identifier: b
        string: 'test'
    and
      hasprefix
        identifier: c
        string: 'test'
      <
        identifier: x
        number: 4
  >=
    identifier: x
    number: 10
`[1:]

	res, err = UnitTestParseWithPPResult("mytest", input, `a > b or a <= p or b hassuffix "test" or c hasprefix "test" and x < 4 or x >= 10`)

	if err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}

	input = "(a in null or c notin d) and false like 9 or x // 6 > 2 % 1"
	expectedOutput = `
or
  and
    or
      in
        identifier: a
        null
      notin
        identifier: c
        identifier: d
    like
      false
      number: 9
  >
    divint
      identifier: x
      number: 6
    modint
      number: 2
      number: 1
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}
}

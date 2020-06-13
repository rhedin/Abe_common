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

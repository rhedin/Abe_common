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

func TestImportParsing(t *testing.T) {

	input := `import "foo/bar.ecal" as foobar
	i := foobar`
	expectedOutput := `
statements
  import
    string: 'foo/bar.ecal'
    identifier: foobar
  :=
    identifier: i
    identifier: foobar
`[1:]

	if res, err := UnitTestParse("mytest", input); err != nil || fmt.Sprint(res) != expectedOutput {
		t.Error("Unexpected parser output:\n", res, "expected was:\n", expectedOutput, "Error:", err)
		return
	}
}

func TestFunctionCalling(t *testing.T) {

	input := `import "foo/bar.ecal" as foobar
	foobar.test()`
	expectedOutput := `
statements
  import
    string: 'foo/bar.ecal'
    identifier: foobar
  identifier: foobar
    identifier: test
      funccall
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

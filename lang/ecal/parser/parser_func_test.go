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

/*
TODO:

func TestFunctionCalling(t *testing.T) {

	input := `import "foo/bar.ecal" as foobar
	foobar.test()`
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
*/

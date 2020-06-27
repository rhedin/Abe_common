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
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	"devt.de/krotik/common/errorutil"
	"devt.de/krotik/common/stringutil"
)

/*
Map of AST nodes corresponding to lexer tokens
*/
var prettyPrinterMap map[string]*template.Template

/*
Map of nodes where the precedence might have changed because of parentheses
*/
var bracketPrecedenceMap map[string]bool

func init() {
	prettyPrinterMap = map[string]*template.Template{

		NodeSTRING: template.Must(template.New(NodeSTRING).Parse("{{.qval}}")),
		NodeNUMBER: template.Must(template.New(NodeNUMBER).Parse("{{.val}}")),
		// NodeIDENTIFIER - Special case (handled in code)

		// Constructed tokens

		// NodeSTATEMENTS - Special case (handled in code)
		// NodeFUNCCALL - Special case (handled in code)
		NodeCOMPACCESS + "_1": template.Must(template.New(NodeCOMPACCESS).Parse("[{{.c1}}]")),
		// TokenLIST - Special case (handled in code)
		// TokenMAP - Special case (handled in code)
		// TokenPARAMS - Special case (handled in code)

		/*


			NodeSTATEMENTS = "statements" // List of statements

			// Assignment statement

			NodeASSIGN = ":="
		*/

		// Assignment statement

		NodeASSIGN + "_2": template.Must(template.New(NodeMINUS).Parse("{{.c1}} := {{.c2}}")),

		// Import statement

		NodeIMPORT + "_2": template.Must(template.New(NodeMINUS).Parse("import {{.c1}} as {{.c2}}")),

		// Arithmetic operators

		NodePLUS + "_1":   template.Must(template.New(NodePLUS).Parse("+{{.c1}}")),
		NodePLUS + "_2":   template.Must(template.New(NodePLUS).Parse("{{.c1}} + {{.c2}}")),
		NodeMINUS + "_1":  template.Must(template.New(NodeMINUS).Parse("-{{.c1}}")),
		NodeMINUS + "_2":  template.Must(template.New(NodeMINUS).Parse("{{.c1}} - {{.c2}}")),
		NodeTIMES + "_2":  template.Must(template.New(NodeTIMES).Parse("{{.c1}} * {{.c2}}")),
		NodeDIV + "_2":    template.Must(template.New(NodeDIV).Parse("{{.c1}} / {{.c2}}")),
		NodeMODINT + "_2": template.Must(template.New(NodeMODINT).Parse("{{.c1}} % {{.c2}}")),
		NodeDIVINT + "_2": template.Must(template.New(NodeDIVINT).Parse("{{.c1}} // {{.c2}}")),

		// Function definition

		NodeFUNC + "_3":   template.Must(template.New(NodeDIVINT).Parse("func {{.c1}}{{.c2}} {\n{{.c3}}}")),
		NodeRETURN:        template.Must(template.New(NodeDIVINT).Parse("return")),
		NodeRETURN + "_1": template.Must(template.New(NodeDIVINT).Parse("return {{.c1}}")),

		// Boolean operators

		NodeOR + "_2":  template.Must(template.New(NodeGEQ).Parse("{{.c1}} or {{.c2}}")),
		NodeAND + "_2": template.Must(template.New(NodeLEQ).Parse("{{.c1}} and {{.c2}}")),
		NodeNOT + "_1": template.Must(template.New(NodeNOT).Parse("not {{.c1}}")),

		// Condition operators

		NodeLIKE + "_2":      template.Must(template.New(NodeGEQ).Parse("{{.c1}} like {{.c2}}")),
		NodeIN + "_2":        template.Must(template.New(NodeLEQ).Parse("{{.c1}} in {{.c2}}")),
		NodeHASPREFIX + "_2": template.Must(template.New(NodeLEQ).Parse("{{.c1}} hasprefix {{.c2}}")),
		NodeHASSUFFIX + "_2": template.Must(template.New(NodeLEQ).Parse("{{.c1}} hassuffix {{.c2}}")),
		NodeNOTIN + "_2":     template.Must(template.New(NodeLEQ).Parse("{{.c1}} notin {{.c2}}")),

		NodeGEQ + "_2": template.Must(template.New(NodeGEQ).Parse("{{.c1}} >= {{.c2}}")),
		NodeLEQ + "_2": template.Must(template.New(NodeLEQ).Parse("{{.c1}} <= {{.c2}}")),
		NodeNEQ + "_2": template.Must(template.New(NodeNEQ).Parse("{{.c1}} != {{.c2}}")),
		NodeEQ + "_2":  template.Must(template.New(NodeEQ).Parse("{{.c1}} == {{.c2}}")),
		NodeGT + "_2":  template.Must(template.New(NodeGT).Parse("{{.c1}} > {{.c2}}")),
		NodeLT + "_2":  template.Must(template.New(NodeLT).Parse("{{.c1}} < {{.c2}}")),

		// Separators

		NodeKVP + "_2":    template.Must(template.New(NodeLT).Parse("{{.c1}} : {{.c2}}")),
		NodePRESET + "_2": template.Must(template.New(NodeLT).Parse("{{.c1}}={{.c2}}")),

		// Constants

		NodeTRUE:  template.Must(template.New(NodeTRUE).Parse("true")),
		NodeFALSE: template.Must(template.New(NodeFALSE).Parse("false")),
		NodeNULL:  template.Must(template.New(NodeNULL).Parse("null")),
	}

	bracketPrecedenceMap = map[string]bool{
		NodePLUS:  true,
		NodeMINUS: true,
		NodeAND:   true,
		NodeOR:    true,
	}
}

/*
PrettyPrint produces pretty printed code from a given AST.
*/
func PrettyPrint(ast *ASTNode) (string, error) {
	var visit func(ast *ASTNode, level int) (string, error)

	ppMetaData := func(ast *ASTNode, ppString string) string {
		ret := ppString

		// Add meta data

		if len(ast.Meta) > 0 {
			for _, meta := range ast.Meta {
				if meta.Type() == MetaDataPreComment {
					ret = fmt.Sprintf("/*%v*/ %v", meta.Value(), ret)
				} else if meta.Type() == MetaDataPostComment {
					ret = fmt.Sprintf("%v #%v", ret, meta.Value())
				}
			}
		}

		return ret
	}

	visit = func(ast *ASTNode, level int) (string, error) {
		var buf bytes.Buffer
		var numChildren = len(ast.Children)

		tempKey := ast.Name
		tempParam := make(map[string]string)

		// First pretty print children

		if numChildren > 0 {
			for i, child := range ast.Children {
				res, err := visit(child, level+1)
				if err != nil {
					return "", err
				}

				if _, ok := bracketPrecedenceMap[child.Name]; ok && ast.binding > child.binding {

					// Put the expression in brackets iff (if and only if) the binding would
					// normally order things differently

					res = fmt.Sprintf("(%v)", res)
				}

				tempParam[fmt.Sprint("c", i+1)] = res
			}

			tempKey += fmt.Sprint("_", len(tempParam))
		}

		// Handle special cases - children in tempParam have been resolved

		if ast.Name == NodeSTATEMENTS {

			// For statements just concat all children

			for i := 0; i < numChildren; i++ {
				buf.WriteString(stringutil.GenerateRollingString(" ", level*4))
				buf.WriteString(tempParam[fmt.Sprint("c", i+1)])
				buf.WriteString("\n")
			}

			return ppMetaData(ast, buf.String()), nil

		} else if ast.Name == NodeFUNCCALL {

			// For statements just concat all children

			for i := 0; i < numChildren; i++ {
				buf.WriteString(tempParam[fmt.Sprint("c", i+1)])
				if i < numChildren-1 {
					buf.WriteString(", ")
				}
			}

			return ppMetaData(ast, buf.String()), nil

		} else if ast.Name == NodeIDENTIFIER {

			buf.WriteString(ast.Token.Val)

			for i := 0; i < numChildren; i++ {
				if ast.Children[i].Name == NodeIDENTIFIER {
					buf.WriteString(".")
					buf.WriteString(tempParam[fmt.Sprint("c", i+1)])
				} else if ast.Children[i].Name == NodeFUNCCALL {
					buf.WriteString("(")
					buf.WriteString(tempParam[fmt.Sprint("c", i+1)])
					buf.WriteString(")")
				} else if ast.Children[i].Name == NodeCOMPACCESS {
					buf.WriteString(tempParam[fmt.Sprint("c", i+1)])
				}
			}

			return ppMetaData(ast, buf.String()), nil
		} else if ast.Name == NodeLIST {

			buf.WriteString("[")
			i := 1
			for ; i < numChildren; i++ {
				buf.WriteString(tempParam[fmt.Sprint("c", i)])
				buf.WriteString(", ")
			}
			buf.WriteString(tempParam[fmt.Sprint("c", i)])
			buf.WriteString("]")

			return ppMetaData(ast, buf.String()), nil

		} else if ast.Name == NodeMAP {

			buf.WriteString("{")
			i := 1
			for ; i < numChildren; i++ {
				buf.WriteString(tempParam[fmt.Sprint("c", i)])
				buf.WriteString(", ")
			}
			buf.WriteString(tempParam[fmt.Sprint("c", i)])
			buf.WriteString("}")

			return ppMetaData(ast, buf.String()), nil
		} else if ast.Name == NodePARAMS {

			buf.WriteString("(")
			i := 1
			for ; i < numChildren; i++ {
				buf.WriteString(tempParam[fmt.Sprint("c", i)])
				buf.WriteString(", ")
			}
			buf.WriteString(tempParam[fmt.Sprint("c", i)])
			buf.WriteString(")")

			return ppMetaData(ast, buf.String()), nil
		}

		if ast.Token != nil {

			// Adding node value to template parameters

			tempParam["val"] = ast.Token.Val
			tempParam["qval"] = strconv.Quote(ast.Token.Val)
		}

		// Retrieve the template

		temp, ok := prettyPrinterMap[tempKey]
		if !ok {
			return "", fmt.Errorf("Could not find template for %v (tempkey: %v)",
				ast.Name, tempKey)
		}

		// Use the children as parameters for template

		errorutil.AssertOk(temp.Execute(&buf, tempParam))

		return ppMetaData(ast, buf.String()), nil
	}

	return visit(ast, 0)
}

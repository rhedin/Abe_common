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

	"devt.de/krotik/common/datautil"
	"devt.de/krotik/common/stringutil"
)

// AST Nodes
// =========

/*
ASTNode models a node in the AST
*/
type ASTNode struct {
	Name     string     // Name of the node
	Token    *LexToken  // Lexer token of this ASTNode
	Children []*ASTNode // Child nodes
	Runtime  Runtime    // Runtime component for this ASTNode

	binding        int                                                             // Binding power of this node
	nullDenotation func(p *parser, self *ASTNode) (*ASTNode, error)                // Configure token as beginning node
	leftDenotation func(p *parser, self *ASTNode, left *ASTNode) (*ASTNode, error) // Configure token as left node
}

/*
Create a new instance of this ASTNode which is connected to a concrete lexer token.
*/
func (n *ASTNode) instance(p *parser, t *LexToken) *ASTNode {

	ret := &ASTNode{n.Name, t, make([]*ASTNode, 0, 2), nil, n.binding, n.nullDenotation, n.leftDenotation}

	if p.rp != nil {
		ret.Runtime = p.rp.Runtime(ret)
	}

	return ret
}

/*
String returns a string representation of this token.
*/
func (n *ASTNode) String() string {
	var buf bytes.Buffer
	n.levelString(0, &buf)
	return buf.String()
}

/*
levelString function to recursively print the tree.
*/
func (n *ASTNode) levelString(indent int, buf *bytes.Buffer) {

	// Print current level

	buf.WriteString(stringutil.GenerateRollingString(" ", indent*2))

	if n.Name == NodeCOMMENT {
		buf.WriteString(fmt.Sprintf("%v: %20v", n.Name, n.Token.Val))
	} else if n.Name == NodeSTRING {
		buf.WriteString(fmt.Sprintf("%v: '%v'", n.Name, n.Token.Val))
	} else if n.Name == NodeNUMBER {
		buf.WriteString(fmt.Sprintf("%v: %v", n.Name, n.Token.Val))
	} else if n.Name == NodeIDENTIFIER {
		buf.WriteString(fmt.Sprintf("%v: %v", n.Name, n.Token.Val))
	} else {
		buf.WriteString(n.Name)
	}

	buf.WriteString("\n")

	// Print children

	for _, child := range n.Children {
		child.levelString(indent+1, buf)
	}
}

// Look ahead buffer
// =================

/*
ASTNode models a node in the AST
*/
type LABuffer struct {
	tokens chan LexToken
	buffer *datautil.RingBuffer
}

/*
Create a new instance of this ASTNode which is connected to a concrete lexer token.
*/
func NewLABuffer(c chan LexToken, size int) *LABuffer {

	if size < 1 {
		size = 1
	}

	ret := &LABuffer{c, datautil.NewRingBuffer(size)}

	v, more := <-ret.tokens
	ret.buffer.Add(v)

	for ret.buffer.Size() < size && more && v.ID != TokenEOF {
		v, more = <-ret.tokens
		ret.buffer.Add(v)
	}

	return ret
}

/*
Next returns the next item.
*/
func (b *LABuffer) Next() (LexToken, bool) {

	ret := b.buffer.Poll()

	if v, more := <-b.tokens; more {
		b.buffer.Add(v)
	}

	if ret == nil {
		return LexToken{ID: TokenEOF}, false
	}

	return ret.(LexToken), true
}

/*
Peek looks inside the buffer starting with 0 as the next item.
*/
func (b *LABuffer) Peek(pos int) (LexToken, bool) {

	if pos >= b.buffer.Size() {
		return LexToken{ID: TokenEOF}, false
	}

	return b.buffer.Get(pos).(LexToken), true
}

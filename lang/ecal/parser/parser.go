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
)

/*
Map of AST nodes corresponding to lexer tokens
*/
var astNodeMap map[LexTokenID]*ASTNode

func init() {
	astNodeMap = map[LexTokenID]*ASTNode{
		TokenEOF: {NodeEOF, nil, nil, nil, 0, ndTerm, nil},

		// Value tokens

		TokenCOMMENT:    {NodeCOMMENT, nil, nil, nil, 0, ndTerm, nil},
		TokenSTRING:     {NodeSTRING, nil, nil, nil, 0, ndTerm, nil},
		TokenNUMBER:     {NodeNUMBER, nil, nil, nil, 0, ndTerm, nil},
		TokenIDENTIFIER: {NodeIDENTIFIER, nil, nil, nil, 0, ndTerm, nil},

		// Constructed tokens

		TokenSTATEMENTS: {NodeSTATEMENTS, nil, nil, nil, 0, nil, nil},
		TokenSEMICOLON:  {"", nil, nil, nil, 0, nil, nil},
		/*
			TokenLIST:       {NodeLIST, nil, nil, nil, 0, nil, nil},
			TokenMAP:        {NodeMAP, nil, nil, nil, 0, nil, nil},
			TokenGUARD:      {NodeGUARD, nil, nil, nil, 0, nil, nil},
		*/

		// Grouping symbols

		TokenLPAREN: {"", nil, nil, nil, 150, ndInner, nil},
		TokenRPAREN: {"", nil, nil, nil, 0, nil, nil},

		// Separators

		TokenCOMMA: {"", nil, nil, nil, 0, nil, nil},

		// Assignment statement

		TokenASSIGN: {NodeASSIGN, nil, nil, nil, 10, nil, ldInfix},

		// Simple arithmetic expressions

		TokenPLUS:   {NodePLUS, nil, nil, nil, 110, ndPrefix, ldInfix},
		TokenMINUS:  {NodeMINUS, nil, nil, nil, 110, ndPrefix, ldInfix},
		TokenTIMES:  {NodeTIMES, nil, nil, nil, 120, nil, ldInfix},
		TokenDIV:    {NodeDIV, nil, nil, nil, 120, nil, ldInfix},
		TokenDIVINT: {NodeDIVINT, nil, nil, nil, 120, nil, ldInfix},
		TokenMODINT: {NodeMODINT, nil, nil, nil, 120, nil, ldInfix},

		// Boolean operators

		TokenOR:  {NodeOR, nil, nil, nil, 30, nil, ldInfix},
		TokenAND: {NodeAND, nil, nil, nil, 40, nil, ldInfix},
		TokenNOT: {NodeNOT, nil, nil, nil, 20, ndPrefix, nil},

		// Condition operators

		TokenLIKE:      {NodeLIKE, nil, nil, nil, 60, nil, ldInfix},
		TokenIN:        {NodeIN, nil, nil, nil, 60, nil, ldInfix},
		TokenHASPREFIX: {NodeHASPREFIX, nil, nil, nil, 60, nil, ldInfix},
		TokenHASSUFFIX: {NodeHASSUFFIX, nil, nil, nil, 60, nil, ldInfix},
		TokenNOTIN:     {NodeNOTIN, nil, nil, nil, 60, nil, ldInfix},

		TokenGEQ: {NodeGEQ, nil, nil, nil, 60, nil, ldInfix},
		TokenLEQ: {NodeLEQ, nil, nil, nil, 60, nil, ldInfix},
		TokenNEQ: {NodeNEQ, nil, nil, nil, 60, nil, ldInfix},
		TokenEQ:  {NodeEQ, nil, nil, nil, 60, nil, ldInfix},
		TokenGT:  {NodeGT, nil, nil, nil, 60, nil, ldInfix},
		TokenLT:  {NodeLT, nil, nil, nil, 60, nil, ldInfix},

		// Constants

		TokenFALSE: {NodeFALSE, nil, nil, nil, 0, ndTerm, nil},
		TokenTRUE:  {NodeTRUE, nil, nil, nil, 0, ndTerm, nil},
		TokenNULL:  {NodeNULL, nil, nil, nil, 0, ndTerm, nil},
	}
}

// Parser
// ======

/*
Parser data structure
*/
type parser struct {
	name   string          // Name to identify the input
	node   *ASTNode        // Current ast node
	tokens *LABuffer       // Buffer which is connected to the channel which contains lex tokens
	rp     RuntimeProvider // Runtime provider which creates runtime components
}

/*
Parse parses a given input string and returns an AST.
*/
func Parse(name string, input string) (*ASTNode, error) {
	return ParseWithRuntime(name, input, nil)
}

/*
ParseWithRuntime parses a given input string and returns an AST decorated with
runtime components.
*/
func ParseWithRuntime(name string, input string, rp RuntimeProvider) (*ASTNode, error) {

	// Create a new parser with a look-ahead buffer of 3

	p := &parser{name, nil, NewLABuffer(Lex(name, input), 3), rp}

	// Read and set initial AST node

	node, err := p.next()

	if err != nil {
		return nil, err
	}

	p.node = node

	n, err := p.run(0)

	if err == nil && hasMoreStatements(p, n) {

		st := astNodeMap[TokenSTATEMENTS].instance(p, nil)
		st.Children = append(st.Children, n)

		for err == nil && hasMoreStatements(p, n) {

			// Skip semicolons

			if p.node.Token.ID == TokenSEMICOLON {
				skipToken(p, TokenSEMICOLON)
			}

			n, err = p.run(0)
			st.Children = append(st.Children, n)
		}

		n = st
	}

	if err == nil && p.node != nil && p.node.Token.ID != TokenEOF {
		token := *p.node.Token
		err = p.newParserError(ErrUnexpectedEnd, fmt.Sprintf("extra token id:%v (%v)",
			token.ID, token), token)
	}

	return n, err
}

/*
run models the main parser function.
*/
func (p *parser) run(rightBinding int) (*ASTNode, error) {
	var err error

	n := p.node

	p.node, err = p.next()
	if err != nil {
		return nil, err
	}

	// Start with the null denotation of this statement / expression

	if n.nullDenotation == nil {
		return nil, p.newParserError(ErrImpossibleNullDenotation,
			n.Token.String(), *n.Token)
	}

	left, err := n.nullDenotation(p, n)
	if err != nil {
		return nil, err
	}

	// Collect left denotations as long as the left binding power is greater
	// than the initial right one

	for rightBinding < p.node.binding {
		var nleft *ASTNode

		n = p.node

		if n.leftDenotation == nil {

			if left.Token.Lline < n.Token.Lline {

				// If the impossible left denotation is on a new line
				// we might be parsing a new statement

				return left, nil
			}

			return nil, p.newParserError(ErrImpossibleLeftDenotation,
				n.Token.String(), *n.Token)
		}

		p.node, err = p.next()

		if err != nil {
			return nil, err
		}

		// Get the next left denotation

		nleft, err = n.leftDenotation(p, n, left)

		left = nleft

		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

/*
next retrieves the next lexer token.
*/
func (p *parser) next() (*ASTNode, error) {

	token, more := p.tokens.Next()

	if !more {

		// Unexpected end of input - the associated token is an empty error token

		return nil, p.newParserError(ErrUnexpectedEnd, "", token)

	} else if token.ID == TokenError {

		// There was a lexer error wrap it in a parser error

		return nil, p.newParserError(ErrLexicalError, token.Val, token)

	} else if node, ok := astNodeMap[token.ID]; ok {

		// We got a normal AST component

		return node.instance(p, &token), nil
	}

	return nil, p.newParserError(ErrUnknownToken, fmt.Sprintf("id:%v (%v)", token.ID, token), token)
}

// Standard null denotation functions
// ==================================

/*
ndTerm is used for terminals.
*/
func ndTerm(p *parser, self *ASTNode) (*ASTNode, error) {
	return self, nil
}

/*
ndInner returns the inner expression of an enclosed block and discard the
block token. This method is used for brackets.
*/
func ndInner(p *parser, self *ASTNode) (*ASTNode, error) {

	// Get the inner expression

	exp, err := p.run(0)
	if err != nil {
		return nil, err
	}

	// We return here the inner expression - discarding the bracket tokens

	return exp, skipToken(p, TokenRPAREN)
}

/*
ndPrefix is used for prefix operators.
*/
func ndPrefix(p *parser, self *ASTNode) (*ASTNode, error) {

	// Make sure a prefix will only prefix the next item

	val, err := p.run(self.binding + 20)
	if err != nil {
		return nil, err
	}

	self.Children = append(self.Children, val)

	return self, nil
}

// Standard left denotation functions
// ==================================

/*
ldInfix is used for infix operators.
*/
func ldInfix(p *parser, self *ASTNode, left *ASTNode) (*ASTNode, error) {

	right, err := p.run(self.binding)
	if err != nil {
		return nil, err
	}

	self.Children = append(self.Children, left)
	self.Children = append(self.Children, right)

	return self, nil
}

// Helper functions
// ================

/*
hasMoreStatements returns true if there are more statements to parse.
*/
func hasMoreStatements(p *parser, currentNode *ASTNode) bool {
	nextNode := p.node

	if nextNode == nil || nextNode.Token.ID == TokenEOF {
		return false
	} else if nextNode.Token.ID == TokenSEMICOLON {
		return true
	}

	return currentNode != nil && currentNode.Token.Lline < nextNode.Token.Lline
}

/*
skipToken skips over a given token.
*/
func skipToken(p *parser, ids ...LexTokenID) error {
	var err error

	canSkip := func(id LexTokenID) bool {
		for _, i := range ids {
			if i == id {
				return true
			}
		}
		return false
	}

	if !canSkip(p.node.Token.ID) {
		if p.node.Token.ID == TokenEOF {
			return p.newParserError(ErrUnexpectedEnd, "", *p.node.Token)
		}
		return p.newParserError(ErrUnexpectedToken, p.node.Token.Val, *p.node.Token)
	}

	// This should never return an error unless we skip over EOF or complex tokens
	// like values

	p.node, err = p.next()

	return err
}

/*
acceptChild accepts the current token as a child.
*/
func acceptChild(p *parser, self *ASTNode, id LexTokenID) error {
	var err error

	current := p.node

	p.node, err = p.next()
	if err != nil {
		return err
	}

	if current.Token.ID == id {
		self.Children = append(self.Children, current)
		return nil
	}

	return p.newParserError(ErrUnexpectedToken, current.Token.Val, *current.Token)
}

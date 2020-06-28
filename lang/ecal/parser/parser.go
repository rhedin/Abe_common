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
Map of AST nodes corresponding to lexer tokens. The map determines how a given
sequence of lexer tokens are organized into an AST.
*/
var astNodeMap map[LexTokenID]*ASTNode

func init() {
	astNodeMap = map[LexTokenID]*ASTNode{
		TokenEOF: {NodeEOF, nil, nil, nil, nil, 0, ndTerm, nil},

		// Value tokens

		TokenSTRING:     {NodeSTRING, nil, nil, nil, nil, 0, ndTerm, nil},
		TokenNUMBER:     {NodeNUMBER, nil, nil, nil, nil, 0, ndTerm, nil},
		TokenIDENTIFIER: {NodeIDENTIFIER, nil, nil, nil, nil, 0, ndIdentifier, nil},

		// Constructed tokens

		TokenSTATEMENTS: {NodeSTATEMENTS, nil, nil, nil, nil, 0, nil, nil},
		TokenFUNCCALL:   {NodeFUNCCALL, nil, nil, nil, nil, 0, nil, nil},
		TokenCOMPACCESS: {NodeCOMPACCESS, nil, nil, nil, nil, 0, nil, nil},
		TokenLIST:       {NodeLIST, nil, nil, nil, nil, 0, nil, nil},
		TokenMAP:        {NodeMAP, nil, nil, nil, nil, 0, nil, nil},
		TokenPARAMS:     {NodePARAMS, nil, nil, nil, nil, 0, nil, nil},
		/*
			TokenGUARD:      {NodeGUARD, nil, nil, nil, 0, nil, nil},
		*/

		// Condition operators

		TokenGEQ: {NodeGEQ, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenLEQ: {NodeLEQ, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenNEQ: {NodeNEQ, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenEQ:  {NodeEQ, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenGT:  {NodeGT, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenLT:  {NodeLT, nil, nil, nil, nil, 60, nil, ldInfix},

		// Grouping symbols

		TokenLPAREN: {"", nil, nil, nil, nil, 150, ndInner, nil},
		TokenRPAREN: {"", nil, nil, nil, nil, 0, nil, nil},
		TokenLBRACK: {"", nil, nil, nil, nil, 150, ndList, nil},
		TokenRBRACK: {"", nil, nil, nil, nil, 0, nil, nil},
		TokenLBRACE: {"", nil, nil, nil, nil, 150, ndMap, nil},
		TokenRBRACE: {"", nil, nil, nil, nil, 0, nil, nil},

		// Separators

		TokenDOT:       {"", nil, nil, nil, nil, 0, nil, nil},
		TokenCOMMA:     {"", nil, nil, nil, nil, 0, nil, nil},
		TokenSEMICOLON: {"", nil, nil, nil, nil, 0, nil, nil},

		// Grouping

		TokenCOLON: {NodeKVP, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenEQUAL: {NodePRESET, nil, nil, nil, nil, 60, nil, ldInfix},

		// Arithmetic operators

		TokenPLUS:   {NodePLUS, nil, nil, nil, nil, 110, ndPrefix, ldInfix},
		TokenMINUS:  {NodeMINUS, nil, nil, nil, nil, 110, ndPrefix, ldInfix},
		TokenTIMES:  {NodeTIMES, nil, nil, nil, nil, 120, nil, ldInfix},
		TokenDIV:    {NodeDIV, nil, nil, nil, nil, 120, nil, ldInfix},
		TokenDIVINT: {NodeDIVINT, nil, nil, nil, nil, 120, nil, ldInfix},
		TokenMODINT: {NodeMODINT, nil, nil, nil, nil, 120, nil, ldInfix},

		// Assignment statement

		TokenASSIGN: {NodeASSIGN, nil, nil, nil, nil, 10, nil, ldInfix},

		// Import statement

		TokenIMPORT: {NodeIMPORT, nil, nil, nil, nil, 0, ndImport, nil},
		TokenAS:     {"", nil, nil, nil, nil, 0, ndImport, nil},

		// Sink definition

		TokenSINK:       {NodeSINK, nil, nil, nil, nil, 0, ndSkink, nil},
		TokenKINDMATCH:  {NodeKINDMATCH, nil, nil, nil, nil, 150, ndPrefix, nil},
		TokenSCOPEMATCH: {NodeSCOPEMATCH, nil, nil, nil, nil, 150, ndPrefix, nil},
		TokenSTATEMATCH: {NodeSTATEMATCH, nil, nil, nil, nil, 150, ndPrefix, nil},
		TokenPRIORITY:   {NodePRIORITY, nil, nil, nil, nil, 150, ndPrefix, nil},
		TokenSUPPRESSES: {NodeSUPPRESSES, nil, nil, nil, nil, 150, ndPrefix, nil},

		// Function definition

		TokenFUNC:   {NodeFUNC, nil, nil, nil, nil, 0, ndFunc, nil},
		TokenRETURN: {NodeRETURN, nil, nil, nil, nil, 0, ndReturn, nil},

		// Boolean operators

		TokenAND: {NodeAND, nil, nil, nil, nil, 40, nil, ldInfix},
		TokenOR:  {NodeOR, nil, nil, nil, nil, 30, nil, ldInfix},
		TokenNOT: {NodeNOT, nil, nil, nil, nil, 20, ndPrefix, nil},

		// Condition operators

		TokenLIKE:      {NodeLIKE, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenIN:        {NodeIN, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenHASPREFIX: {NodeHASPREFIX, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenHASSUFFIX: {NodeHASSUFFIX, nil, nil, nil, nil, 60, nil, ldInfix},
		TokenNOTIN:     {NodeNOTIN, nil, nil, nil, nil, 60, nil, ldInfix},

		// Constant terminals

		TokenFALSE: {NodeFALSE, nil, nil, nil, nil, 0, ndTerm, nil},
		TokenTRUE:  {NodeTRUE, nil, nil, nil, nil, 0, ndTerm, nil},
		TokenNULL:  {NodeNULL, nil, nil, nil, nil, 0, ndTerm, nil},

		/*
			// Conditional statements

			TokenIF
			TokenELIF
			TokenELSE

			// Loop statements

			TokenFOR
			TokenBREAK
			TokenCONTINUE
		*/
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
	var preComments []MetaData
	var postComments []MetaData

	token, more := p.tokens.Next()

	// Skip over pre comment token

	for more && token.ID == TokenPRECOMMENT {
		preComments = append(preComments, NewLexTokenInstance(token))
		token, more = p.tokens.Next()
	}

	// Skip over post comment token

	for more && token.ID == TokenPOSTCOMMENT {
		postComments = append(postComments, NewLexTokenInstance(token))
		token, more = p.tokens.Next()
	}

	if !more {

		// Unexpected end of input - the associated token is an empty error token

		return nil, p.newParserError(ErrUnexpectedEnd, "", token)

	} else if token.ID == TokenError {

		// There was a lexer error wrap it in a parser error

		return nil, p.newParserError(ErrLexicalError, token.Val, token)

	} else if node, ok := astNodeMap[token.ID]; ok {

		// We got a normal AST component

		ret := node.instance(p, &token)

		ret.Meta = append(ret.Meta, preComments...) // Attach pre comments to the next AST node
		if len(postComments) > 0 && p.node != nil {
			p.node.Meta = append(p.node.Meta, postComments...) // Attach post comments to the previous AST node
		}

		return ret, nil
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

// Null denotation functions for specific expressions
// ==================================================

/*
ndImport is used to parse imports.
*/
func ndImport(p *parser, self *ASTNode) (*ASTNode, error) {

	// Must specify a file path

	err := acceptChild(p, self, TokenSTRING)

	if err == nil {

		// Must specify AS

		if err = skipToken(p, TokenAS); err == nil {

			// Must specify an identifier

			err = acceptChild(p, self, TokenIDENTIFIER)
		}
	}

	return self, err
}

/*
ndSink is used to parse sinks.
*/
func ndSkink(p *parser, self *ASTNode) (*ASTNode, error) {
	var ret *ASTNode

	// Must specify a name

	err := acceptChild(p, self, TokenIDENTIFIER)

	if err == nil {

		// Parse the rest of the parameters as children until we reach the body

		for p.node.Token.ID != TokenEOF && p.node.Token.ID != TokenLBRACE {
			exp, err := p.run(150)
			if err != nil {
				return nil, err
			}

			self.Children = append(self.Children, exp)

			// Skip commas

			if p.node.Token.ID == TokenCOMMA {
				skipToken(p, TokenCOMMA)
			}
		}

		// Parse the body

		ret, err = parseInnerStatements(p, self)
	}

	return ret, err
}

/*
ndFunc is used to parse function definitions.
*/
func ndFunc(p *parser, self *ASTNode) (*ASTNode, error) {

	// Must specify a function name

	err := acceptChild(p, self, TokenIDENTIFIER)

	// Read in parameters

	if err == nil {
		err = skipToken(p, TokenLPAREN)

		params := astNodeMap[TokenPARAMS].instance(p, nil)
		self.Children = append(self.Children, params)

		for err == nil && p.node.Token.ID != TokenRPAREN {

			// Parse all the expressions inside

			exp, err := p.run(0)
			if err == nil {
				params.Children = append(params.Children, exp)

				if p.node.Token.ID == TokenCOMMA {
					err = skipToken(p, TokenCOMMA)
				}
			}
		}

		if err == nil {
			err = skipToken(p, TokenRPAREN)
		}
	}

	if err == nil {

		// Parse the body

		self, err = parseInnerStatements(p, self)
	}

	return self, err
}

/*
ndReturn is used to parse return statements.
*/
func ndReturn(p *parser, self *ASTNode) (*ASTNode, error) {
	var err error

	if self.Token.Lline == p.node.Token.Lline {
		var val *ASTNode

		// Consume the next expression only if it is on the same line

		val, err = p.run(0)

		if err == nil {
			self.Children = append(self.Children, val)
		}
	}

	return self, err
}

/*
ndIdentifier is to parse identifiers and function calls.
*/
func ndIdentifier(p *parser, self *ASTNode) (*ASTNode, error) {
	var parseMore, parseSegment, parseFuncCall, parseCompositionAccess func(parent *ASTNode) error

	parseMore = func(current *ASTNode) error {
		var err error

		if p.node.Token.ID == TokenDOT {
			err = parseSegment(current)
		} else if p.node.Token.ID == TokenLPAREN {
			err = parseFuncCall(current)
		} else if p.node.Token.ID == TokenLBRACK {
			err = parseCompositionAccess(current)
		}

		return err
	}

	parseSegment = func(current *ASTNode) error {
		var err error
		var next *ASTNode

		if err = skipToken(p, TokenDOT); err == nil {
			next = p.node
			if err = acceptChild(p, current, TokenIDENTIFIER); err == nil {
				err = parseMore(next)
			}
		}

		return err
	}

	parseFuncCall = func(current *ASTNode) error {
		err := skipToken(p, TokenLPAREN)

		fc := astNodeMap[TokenFUNCCALL].instance(p, nil)
		current.Children = append(current.Children, fc)

		// Read in parameters

		for err == nil && p.node.Token.ID != TokenRPAREN {

			// Parse all the expressions inside the directives

			exp, err := p.run(0)
			if err == nil {
				fc.Children = append(fc.Children, exp)

				if p.node.Token.ID == TokenCOMMA {
					err = skipToken(p, TokenCOMMA)
				}
			}
		}

		if err == nil {
			err = skipToken(p, TokenRPAREN)
			if err == nil {
				err = parseMore(current)
			}
		}

		return err
	}

	parseCompositionAccess = func(current *ASTNode) error {
		err := skipToken(p, TokenLBRACK)

		ca := astNodeMap[TokenCOMPACCESS].instance(p, nil)
		current.Children = append(current.Children, ca)

		// Parse all the expressions inside the directives

		exp, err := p.run(0)
		if err == nil {
			ca.Children = append(ca.Children, exp)

			if err = skipToken(p, TokenRBRACK); err == nil {
				err = parseMore(current)
			}
		}

		return err
	}

	return self, parseMore(self)
}

/*
ndList is used to collect elements of a list.
*/
func ndList(p *parser, self *ASTNode) (*ASTNode, error) {

	// Create a list token

	st := astNodeMap[TokenLIST].instance(p, self.Token)

	// Get the inner expression

	for p.node.Token.ID != TokenRBRACK {

		// Parse all the expressions inside

		exp, err := p.run(0)
		if err != nil {
			return nil, err
		}

		st.Children = append(st.Children, exp)

		if p.node.Token.ID == TokenCOMMA {
			skipToken(p, TokenCOMMA)
		}
	}

	// Must have a closing bracket

	return st, skipToken(p, TokenRBRACK)
}

/*
ndMap is used to collect elements of a map.
*/
func ndMap(p *parser, self *ASTNode) (*ASTNode, error) {

	// Create a map token

	st := astNodeMap[TokenMAP].instance(p, self.Token)

	// Get the inner expression

	for p.node.Token.ID != TokenRBRACE {

		// Parse all the expressions inside

		exp, err := p.run(0)
		if err != nil {
			return nil, err
		}

		st.Children = append(st.Children, exp)

		if p.node.Token.ID == TokenCOMMA {
			if err := skipToken(p, TokenCOMMA); err != nil {
				return nil, err
			}
		}
	}

	// Must have a closing brace

	return st, skipToken(p, TokenRBRACE)
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

/*
parseInnerStatements collects the inner statements of a block statement. It
is assumed that a block statement starts with a left brace '{' and ends with
a right brace '}'.
*/
func parseInnerStatements(p *parser, self *ASTNode) (*ASTNode, error) {

	// Must start with an opening brace

	if err := skipToken(p, TokenLBRACE); err != nil {
		return nil, err
	}

	// Always create a statements node

	st := astNodeMap[TokenSTATEMENTS].instance(p, nil)
	self.Children = append(self.Children, st)

	// Check if there are actually children

	if p.node != nil && p.node.Token.ID != TokenRBRACE {

		n, err := p.run(0)

		if p.node != nil && p.node.Token.ID != TokenEOF {

			st.Children = append(st.Children, n)

			for hasMoreStatements(p, n) {

				if p.node.Token.ID == TokenSEMICOLON {
					skipToken(p, TokenSEMICOLON)
				} else if p.node.Token.ID == TokenRBRACE {
					break
				}

				n, err = p.run(0)
				st.Children = append(st.Children, n)
			}
		}

		if err != nil {
			return nil, err
		}
	}

	// Must end with a closing brace

	return self, skipToken(p, TokenRBRACE)
}

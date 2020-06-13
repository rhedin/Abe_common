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
	"testing"
)

func TestLABuffer(t *testing.T) {

	buf := NewLABuffer(Lex("test", "1 2 3 4 5 6 7 8 9"), 3)

	if token, ok := buf.Next(); token.Val != "1" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "2" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Check Peek

	if token, ok := buf.Peek(0); token.Val != "3" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(1); token.Val != "4" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(2); token.Val != "5" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(3); token.ID != TokenEOF || ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Continue

	if token, ok := buf.Next(); token.Val != "3" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "4" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "5" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "6" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "7" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "8" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Check Peek

	if token, ok := buf.Peek(0); token.Val != "9" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(1); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(2); token.ID != TokenEOF || ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Continue

	if token, ok := buf.Next(); token.Val != "9" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Check Peek

	if token, ok := buf.Peek(0); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(1); token.ID != TokenEOF || ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Continue

	if token, ok := buf.Next(); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// New Buffer

	buf = NewLABuffer(Lex("test", "1 2 3"), 3)

	if token, ok := buf.Next(); token.Val != "1" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "2" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// Check Peek

	if token, ok := buf.Peek(0); token.Val != "3" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(1); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(2); token.ID != TokenEOF || ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.Val != "3" || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	// New Buffer - test edge case

	buf = NewLABuffer(Lex("test", ""), 0)

	if token, ok := buf.Peek(0); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.ID != TokenEOF || !ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Peek(0); token.ID != TokenEOF || ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}

	if token, ok := buf.Next(); token.ID != TokenEOF || ok {
		t.Error("Unexpected result: ", token, ok)
		return
	}
}

package lex

import (
	"strings"
)

// AtomNode handles atoms (symbols, identifiers, operators).
// This is typically the last node in the chain and handles anything
// that wasn't matched by previous nodes.
type AtomNode struct {
	BaseNode
}

// NewAtomNode creates a new AtomNode.
func NewAtomNode() *AtomNode {
	return &AtomNode{}
}

// Handle parses an atom token.
func (a *AtomNode) Handle(r Reader) (Token, bool) {
	if r.EOF() {
		return nil, false
	}

	ch := r.Current()

	// Skip if it's whitespace or a delimiter that shouldn't be an atom
	if isWhitespace(ch) {
		return nil, false
	}

	// Reserved characters that are not valid atom starters
	if ch == '(' || ch == ')' || ch == '"' {
		return nil, false
	}

	file, pos, line, linePos := r.Snapshot()

	// Parse the atom
	var sb strings.Builder

	for !r.EOF() {
		ch := r.Current()

		// Stop at whitespace or delimiters
		if isWhitespace(ch) || isDelimiter(ch) {
			break
		}

		sb.WriteByte(ch)
		r.Advance(1)
	}

	raw := sb.String()
	if raw == "" {
		return nil, false
	}

	return NewToken(file, pos, line, linePos, raw, Atom), true
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDelimiter(ch byte) bool {
	return ch == '(' || ch == ')' || ch == '"' || ch == '\'' || ch == '`' || ch == ','
}

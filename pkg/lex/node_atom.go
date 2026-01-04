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
func (a *AtomNode) Handle(input string, pos int, line int, linePos int, file string) (Token, int, bool) {
	if pos >= len(input) {
		return nil, 0, false
	}

	ch := input[pos]

	// Skip if it's whitespace or a delimiter that shouldn't be an atom
	if isWhitespace(ch) {
		return nil, 0, false
	}

	// Reserved characters that are not valid atom starters
	if ch == '(' || ch == ')' || ch == '"' {
		return nil, 0, false
	}

	// Parse the atom
	var sb strings.Builder
	i := pos

	for i < len(input) {
		ch := input[i]

		// Stop at whitespace or delimiters
		if isWhitespace(ch) || isDelimiter(ch) {
			break
		}

		sb.WriteByte(ch)
		i++
	}

	raw := sb.String()
	if raw == "" {
		return nil, 0, false
	}

	return NewToken(file, pos, line, linePos, raw, Atom), i - pos, true
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDelimiter(ch byte) bool {
	return ch == '(' || ch == ')' || ch == '"' || ch == '\'' || ch == '`' || ch == ','
}

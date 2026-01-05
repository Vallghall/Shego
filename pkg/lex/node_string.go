package lex

import (
	"errors"
	"strings"
)

// ErrUnclosedString is returned when a string literal is not properly closed.
var ErrUnclosedString = errors.New("unclosed string literal")

// StringNode handles quoted string tokens: "..."
type StringNode struct {
	BaseNode
}

// NewStringNode creates a new StringNode.
func NewStringNode() *StringNode {
	return &StringNode{}
}

// Handle checks if the current character starts a string literal and parses it.
func (s *StringNode) Handle(r Reader) (Token, bool) {
	if r.EOF() {
		return nil, false
	}

	if r.Current() != '"' {
		return s.PassToNext(r)
	}

	file, pos, line, linePos := r.Snapshot()

	// Parse the string literal
	var sb strings.Builder
	sb.WriteByte('"')
	r.Advance(1) // Skip opening quote

	escaped := false

	for !r.EOF() {
		ch := r.Current()

		if escaped {
			// Handle escape sequences
			switch ch {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			default:
				// Unknown escape, keep as-is
				sb.WriteByte('\\')
				sb.WriteByte(ch)
			}
			escaped = false
			r.Advance(1)
			continue
		}

		if ch == '\\' {
			escaped = true
			r.Advance(1)
			continue
		}

		if ch == '"' {
			sb.WriteByte('"')
			r.Advance(1)
			// Successfully closed the string
			return NewToken(file, pos, line, linePos, sb.String(), String), true
		}

		// Don't allow unescaped newlines in strings
		if ch == '\n' {
			return nil, false // Will be handled as error by lexer
		}

		sb.WriteByte(ch)
		r.Advance(1)
	}

	// Reached end of input without closing quote
	return nil, false
}

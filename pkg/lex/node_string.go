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
func (s *StringNode) Handle(input string, pos int, line int, linePos int, file string) (Token, int, bool) {
	if pos >= len(input) {
		return nil, 0, false
	}

	if input[pos] != '"' {
		return s.PassToNext(input, pos, line, linePos, file)
	}

	// Parse the string literal
	var sb strings.Builder
	sb.WriteByte('"')
	i := pos + 1
	escaped := false

	for i < len(input) {
		ch := input[i]

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
			i++
			continue
		}

		if ch == '\\' {
			escaped = true
			i++
			continue
		}

		if ch == '"' {
			sb.WriteByte('"')
			i++
			// Successfully closed the string
			return NewToken(file, pos, line, linePos, sb.String(), String), i - pos, true
		}

		// Don't allow unescaped newlines in strings
		if ch == '\n' {
			return nil, 0, false // Will be handled as error by lexer
		}

		sb.WriteByte(ch)
		i++
	}

	// Reached end of input without closing quote
	return nil, 0, false
}

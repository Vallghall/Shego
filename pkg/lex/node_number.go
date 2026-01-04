package lex

import (
	"strings"
)

// NumberNode handles numeric tokens, including signed numbers.
// Special handling for +/- : if followed by a digit, parse as signed number.
// If followed by whitespace or delimiter, pass to next node (Atom).
type NumberNode struct {
	BaseNode
}

// NewNumberNode creates a new NumberNode.
func NewNumberNode() *NumberNode {
	return &NumberNode{}
}

// Handle checks if the current position starts a number and parses it.
func (n *NumberNode) Handle(input string, pos int, line int, linePos int, file string) (Token, int, bool) {
	if pos >= len(input) {
		return nil, 0, false
	}

	ch := input[pos]

	// Check for sign prefix
	hasSign := ch == '+' || ch == '-'
	if hasSign {
		// Look ahead to see if this is a signed number or a standalone operator
		if pos+1 >= len(input) {
			// End of input, pass to atom handler
			return n.PassToNext(input, pos, line, linePos, file)
		}

		nextCh := input[pos+1]
		if !isDigit(nextCh) && nextCh != '.' {
			// Not followed by digit or decimal point, pass to atom handler
			return n.PassToNext(input, pos, line, linePos, file)
		}
	} else if !isDigit(ch) && ch != '.' {
		// Doesn't start with digit, sign, or decimal point
		return n.PassToNext(input, pos, line, linePos, file)
	}

	// Parse the number
	var sb strings.Builder
	i := pos

	// Handle optional sign
	if hasSign {
		sb.WriteByte(input[i])
		i++
	}

	// Parse integer part
	hasIntPart := false
	for i < len(input) && isDigit(input[i]) {
		sb.WriteByte(input[i])
		hasIntPart = true
		i++
	}

	// Parse optional decimal part
	hasDecimal := false
	if i < len(input) && input[i] == '.' {
		// Look ahead to ensure there's at least one digit after the dot
		// or we already have an integer part
		if i+1 < len(input) && isDigit(input[i+1]) {
			sb.WriteByte('.')
			i++
			hasDecimal = true
			for i < len(input) && isDigit(input[i]) {
				sb.WriteByte(input[i])
				i++
			}
		} else if hasIntPart {
			// Allow trailing dot like "42."
			sb.WriteByte('.')
			i++
			hasDecimal = true
		}
	}

	// Must have at least some digits
	if !hasIntPart && !hasDecimal {
		return n.PassToNext(input, pos, line, linePos, file)
	}

	// Check that we actually parsed something beyond just a sign
	raw := sb.String()
	if raw == "+" || raw == "-" || raw == "." {
		return n.PassToNext(input, pos, line, linePos, file)
	}

	return NewToken(file, pos, line, linePos, raw, Number), i - pos, true
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

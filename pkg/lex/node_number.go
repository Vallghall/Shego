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
func (n *NumberNode) Handle(r Reader) (Token, bool) {
	if r.EOF() {
		return nil, false
	}

	ch := r.Current()

	// Check for sign prefix
	hasSign := ch == '+' || ch == '-'
	if hasSign {
		// Look ahead to see if this is a signed number or a standalone operator
		nextCh, ok := r.Peek(1)
		if !ok {
			// End of input, pass to atom handler
			return n.PassToNext(r)
		}

		if !isDigit(nextCh) && nextCh != '.' {
			// Not followed by digit or decimal point, pass to atom handler
			return n.PassToNext(r)
		}
	} else if !isDigit(ch) && ch != '.' {
		// Doesn't start with digit, sign, or decimal point
		return n.PassToNext(r)
	}

	file, pos, line, linePos := r.Snapshot()

	// Parse the number
	var sb strings.Builder
	consumed := 0

	// Handle optional sign
	if hasSign {
		sb.WriteByte(r.Current())
		r.Advance(1)
		consumed++
	}

	// Parse integer part
	hasIntPart := false
	for !r.EOF() && isDigit(r.Current()) {
		sb.WriteByte(r.Current())
		r.Advance(1)
		consumed++
		hasIntPart = true
	}

	// Parse optional decimal part
	hasDecimal := false
	if !r.EOF() && r.Current() == '.' {
		// Look ahead to ensure there's at least one digit after the dot
		// or we already have an integer part
		nextCh, hasNext := r.Peek(1)
		if hasNext && isDigit(nextCh) {
			sb.WriteByte('.')
			r.Advance(1)
			consumed++
			hasDecimal = true
			for !r.EOF() && isDigit(r.Current()) {
				sb.WriteByte(r.Current())
				r.Advance(1)
				consumed++
			}
		} else if hasIntPart {
			// Allow trailing dot like "42."
			sb.WriteByte('.')
			r.Advance(1)
			consumed++
			hasDecimal = true
		}
	}

	// Must have at least some digits
	if !hasIntPart && !hasDecimal {
		return n.PassToNext(r)
	}

	// Check that we actually parsed something beyond just a sign
	raw := sb.String()
	if raw == "+" || raw == "-" || raw == "." {
		return n.PassToNext(r)
	}

	return NewToken(file, pos, line, linePos, raw, Number), true
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

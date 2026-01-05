package ast

import (
	"github.com/Vallghall/schego/pkg/lex"
)

// ListParseNode handles parenthesized S-expressions (non-special forms).
type ListParseNode struct {
	BaseParseNode
	// chain is the full parsing chain for recursive parsing
	chain ParseNode
}

// NewListParseNode creates a new ListParseNode.
func NewListParseNode() *ListParseNode {
	return &ListParseNode{}
}

// SetChain sets the parsing chain (called after chain is fully built).
func (l *ListParseNode) SetChain(chain ParseNode) {
	l.chain = chain
}

// Handle parses a parenthesized list and its contents.
func (l *ListParseNode) Handle(r TokenReader) (Node, bool) {
	if r.EOF() {
		return nil, false
	}

	tok := r.Current()
	if tok == nil || tok.Kind() != lex.ParenOpen {
		return l.PassToNext(r)
	}

	openParen := tok
	r.Advance() // consume '('

	// Parse list elements
	var elements []Node
	for !r.EOF() {
		current := r.Current()
		if current == nil {
			break
		}

		// Check for closing paren
		if current.Kind() == lex.ParenClose {
			r.Advance() // consume ')'

			// If has elements, treat as a call
			if len(elements) > 0 {
				return NewCallNode(elements[0], elements[1:]), true
			}

			// Empty list
			return NewListNode(openParen, elements), true
		}

		// Parse the next element using the full chain
		elem, handled := l.chain.Handle(r)
		if !handled {
			// Parsing error - unexpected token
			return nil, false
		}
		elements = append(elements, elem)
	}

	// Unclosed parenthesis
	return nil, false
}

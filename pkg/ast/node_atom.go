package ast

import (
	"github.com/Vallghall/schego/pkg/lex"
)

// AtomParseNode handles atomic tokens: numbers, strings, and symbols.
type AtomParseNode struct {
	BaseParseNode
}

// NewAtomParseNode creates a new AtomParseNode.
func NewAtomParseNode() *AtomParseNode {
	return &AtomParseNode{}
}

// Handle parses atomic tokens into their corresponding AST nodes.
func (a *AtomParseNode) Handle(r TokenReader) (Node, bool) {
	if r.EOF() {
		return nil, false
	}

	tok := r.Current()
	if tok == nil {
		return nil, false
	}

	switch tok.Kind() {
	case lex.Number:
		r.Advance()
		return NewNumberNode(tok), true

	case lex.String:
		r.Advance()
		return NewStringNode(tok), true

	case lex.Atom:
		r.Advance()
		return NewSymbolNode(tok), true

	default:
		return a.PassToNext(r)
	}
}

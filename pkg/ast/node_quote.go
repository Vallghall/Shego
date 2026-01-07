package ast

import (
	"github.com/Vallghall/schego/pkg/lex"
)

// QuotePrefixNode handles prefix quote tokens: ', `, ,, and ,@
// These are transformed into QuoteNode, QuasiquoteNode, UnquoteNode, and UnquoteSplicingNode respectively.
type QuotePrefixNode struct {
	BaseParseNode
	// chain is the full parsing chain for recursive parsing
	chain ParseNode
}

// NewQuotePrefixNode creates a new QuotePrefixNode.
func NewQuotePrefixNode() *QuotePrefixNode {
	return &QuotePrefixNode{}
}

// SetChain sets the parsing chain (called after chain is fully built).
func (q *QuotePrefixNode) SetChain(chain ParseNode) {
	q.chain = chain
}

// Handle checks if the current token is a quote prefix and parses accordingly.
func (q *QuotePrefixNode) Handle(r TokenReader) (Node, bool) {
	if r.EOF() {
		return nil, false
	}

	tok := r.Current()
	if tok == nil {
		return nil, false
	}

	switch tok.Kind() {
	case lex.Quote:
		r.Advance() // consume '
		return q.parseQuote(r, tok)
	case lex.Backtick:
		r.Advance() // consume `
		return q.parseQuasiquote(r, tok)
	case lex.Comma:
		r.Advance() // consume ,
		return q.parseUnquote(r, tok)
	case lex.CommaAt:
		r.Advance() // consume ,@
		return q.parseUnquoteSplicing(r, tok)
	default:
		return q.PassToNext(r)
	}
}

// parseQuote parses the expression following ' and wraps it in a QuoteNode.
func (q *QuotePrefixNode) parseQuote(r TokenReader, tok lex.Token) (Node, bool) {
	node, handled := q.chain.Handle(r)
	if !handled || node == nil {
		return nil, false
	}
	return NewQuoteNode(tok, node), true
}

// parseQuasiquote parses the expression following ` and wraps it in a QuasiquoteNode.
func (q *QuotePrefixNode) parseQuasiquote(r TokenReader, tok lex.Token) (Node, bool) {
	node, handled := q.chain.Handle(r)
	if !handled || node == nil {
		return nil, false
	}
	return NewQuasiquoteNode(tok, node), true
}

// parseUnquote parses the expression following , and wraps it in an UnquoteNode.
func (q *QuotePrefixNode) parseUnquote(r TokenReader, tok lex.Token) (Node, bool) {
	node, handled := q.chain.Handle(r)
	if !handled || node == nil {
		return nil, false
	}
	return NewUnquoteNode(tok, node), true
}

// parseUnquoteSplicing parses the expression following ,@ and wraps it in an UnquoteSplicingNode.
func (q *QuotePrefixNode) parseUnquoteSplicing(r TokenReader, tok lex.Token) (Node, bool) {
	node, handled := q.chain.Handle(r)
	if !handled || node == nil {
		return nil, false
	}
	return NewUnquoteSplicingNode(tok, node), true
}

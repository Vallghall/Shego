package ast

import (
	"github.com/Vallghall/schego/pkg/lex"
)

// TokenReader provides methods for reading tokens during parsing.
type TokenReader interface {
	// EOF returns true if all tokens have been consumed.
	EOF() bool
	// Current returns the current token, or nil if at EOF.
	Current() lex.Token
	// Peek returns the token at offset positions ahead.
	Peek(offset int) lex.Token
	// Advance moves to the next token.
	Advance()
	// Position returns the current token index.
	Position() int
}

// tokenReader is the concrete implementation of TokenReader.
type tokenReader struct {
	tokens []lex.Token
	pos    int
}

// NewTokenReader creates a new TokenReader for the given tokens.
func NewTokenReader(tokens []lex.Token) TokenReader {
	return &tokenReader{
		tokens: tokens,
		pos:    0,
	}
}

func (r *tokenReader) EOF() bool {
	return r.pos >= len(r.tokens)
}

func (r *tokenReader) Current() lex.Token {
	if r.EOF() {
		return nil
	}
	return r.tokens[r.pos]
}

func (r *tokenReader) Peek(offset int) lex.Token {
	idx := r.pos + offset
	if idx < 0 || idx >= len(r.tokens) {
		return nil
	}
	return r.tokens[idx]
}

func (r *tokenReader) Advance() {
	if !r.EOF() {
		r.pos++
	}
}

func (r *tokenReader) Position() int {
	return r.pos
}

// ParseNode represents a handler in the Chain of Responsibility for parsing.
type ParseNode interface {
	// Handle attempts to parse starting from the current token.
	// Returns:
	//   - node: the parsed AST node if successful, nil otherwise
	//   - handled: true if this handler processed the input
	Handle(r TokenReader) (node Node, handled bool)

	// SetNext sets the next node in the chain.
	SetNext(ParseNode)

	// Next returns the next node in the chain.
	Next() ParseNode
}

// BaseParseNode provides common chain traversal logic.
type BaseParseNode struct {
	next ParseNode
}

// SetNext sets the next node in the chain.
func (b *BaseParseNode) SetNext(n ParseNode) {
	b.next = n
}

// Next returns the next node in the chain.
func (b *BaseParseNode) Next() ParseNode {
	return b.next
}

// PassToNext delegates handling to the next node in the chain.
func (b *BaseParseNode) PassToNext(r TokenReader) (Node, bool) {
	if b.next != nil {
		return b.next.Handle(r)
	}
	return nil, false
}

// BuildParseChain constructs a chain from the given nodes in order.
func BuildParseChain(nodes ...ParseNode) ParseNode {
	if len(nodes) == 0 {
		return nil
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].SetNext(nodes[i+1])
	}
	return nodes[0]
}

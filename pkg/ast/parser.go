package ast

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/lex"
)

// Parser parses tokens into an AST.
type Parser interface {
	// Parse parses the given tokens and returns a slice of AST nodes.
	Parse(tokens []lex.Token) ([]Node, error)
}

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message string
	Token   lex.Token
}

func (e *ParseError) Error() string {
	if e.Token != nil {
		file := e.Token.File()
		if file != "" {
			return fmt.Sprintf("%s:%d:%d: %s", file, e.Token.Line(), e.Token.LinePosition(), e.Message)
		}
		return fmt.Sprintf("%d:%d: %s", e.Token.Line(), e.Token.LinePosition(), e.Message)
	}
	return e.Message
}

// parser is the concrete implementation of Parser.
type parser struct {
	chain ParseNode
}

// New creates a new Parser with the default parsing chain.
func New() Parser {
	p := &parser{}

	// Create the nodes
	specialForms := NewSpecialFormsNode()
	listNode := NewListParseNode()
	atomNode := NewAtomParseNode()

	// Build the chain: SpecialForms -> List -> Atom
	p.chain = BuildParseChain(
		specialForms,
		listNode,
		atomNode,
	)

	// Set the chain references for recursive parsing
	specialForms.SetChain(p.chain)
	listNode.SetChain(p.chain)

	return p
}

// Parse parses the given tokens and returns a slice of AST nodes.
func (p *parser) Parse(tokens []lex.Token) ([]Node, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	r := NewTokenReader(tokens)
	var nodes []Node

	for !r.EOF() {
		node, handled := p.chain.Handle(r)
		if !handled {
			tok := r.Current()
			if tok != nil {
				return nil, &ParseError{
					Message: fmt.Sprintf("unexpected token: %s", tok.Raw()),
					Token:   tok,
				}
			}
			return nil, &ParseError{
				Message: "unexpected end of input",
			}
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

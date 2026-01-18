package ast

import (
	"github.com/Vallghall/schego/pkg/lex"
)

// specialKeywords maps keyword names to their handlers
var specialKeywords = map[string]bool{
	"if":               true,
	"begin":            true,
	"define":           true,
	"set!":             true,
	"lambda":           true,
	"let":              true,
	"cond":             true,
	"quote":            true,
	"quasiquote":       true,
	"unquote":          true,
	"unquote-splicing": true,
}

// SpecialFormsNode handles special form expressions.
// It checks if the current list starts with a special keyword and parses accordingly.
type SpecialFormsNode struct {
	BaseParseNode
	// chain is the full parsing chain for recursive parsing
	chain ParseNode
}

// NewSpecialFormsNode creates a new SpecialFormsNode.
func NewSpecialFormsNode() *SpecialFormsNode {
	return &SpecialFormsNode{}
}

// SetChain sets the parsing chain (called after chain is fully built).
func (s *SpecialFormsNode) SetChain(chain ParseNode) {
	s.chain = chain
}

// Handle checks if this is a special form and parses it, otherwise passes to next.
func (s *SpecialFormsNode) Handle(r TokenReader) (Node, bool) {
	if r.EOF() {
		return nil, false
	}

	tok := r.Current()
	if tok == nil || tok.Kind() != lex.ParenOpen {
		return s.PassToNext(r)
	}

	// Peek at the next token to see if it's a special keyword
	nextTok := r.Peek(1)
	if nextTok == nil || nextTok.Kind() != lex.Atom {
		return s.PassToNext(r)
	}

	keyword := nextTok.Raw()
	if !specialKeywords[keyword] {
		return s.PassToNext(r)
	}

	// This is a special form - parse it
	openParen := tok
	r.Advance() // consume '('

	// Parse all elements in the list
	var elements []Node
	for !r.EOF() {
		current := r.Current()
		if current == nil {
			break
		}

		if current.Kind() == lex.ParenClose {
			r.Advance() // consume ')'
			return s.transform(openParen, keyword, elements)
		}

		elem, handled := s.chain.Handle(r)
		if !handled {
			return nil, false
		}
		elements = append(elements, elem)
	}

	// Unclosed parenthesis
	return nil, false
}

// transform converts parsed elements into the appropriate special form node.
func (s *SpecialFormsNode) transform(openParen lex.Token, keyword string, elements []Node) (Node, bool) {
	// First element should be the keyword symbol - skip it
	if len(elements) == 0 {
		return nil, false
	}
	elements = elements[1:] // Remove the keyword from elements

	switch keyword {
	case "if":
		return s.parseIf(openParen, elements)
	case "begin":
		return s.parseBegin(openParen, elements)
	case "define":
		return s.parseDefine(openParen, elements)
	case "set!":
		return s.parseSet(openParen, elements)
	case "lambda":
		return s.parseLambda(openParen, elements)
	case "let":
		return s.parseLet(openParen, elements)
	case "cond":
		return s.parseCond(openParen, elements)
	case "quote":
		return s.parseQuote(openParen, elements)
	case "quasiquote":
		return s.parseQuasiquote(openParen, elements)
	case "unquote":
		return s.parseUnquote(openParen, elements)
	case "unquote-splicing":
		return s.parseUnquoteSplicing(openParen, elements)
	default:
		return nil, false
	}
}

// extractListElements extracts elements from either a ListNode or CallNode.
// This is needed because nested lists get transformed to CallNodes during parsing.
func extractListElements(n Node) ([]Node, bool) {
	switch node := n.(type) {
	case *ListNode:
		return node.Elements, true
	case *CallNode:
		// Reconstruct elements: operator + args
		elements := make([]Node, 0, 1+len(node.Args))
		elements = append(elements, node.Operator)
		elements = append(elements, node.Args...)
		return elements, true
	default:
		return nil, false
	}
}

// parseIf handles (if condition then [else]) expressions.
func (s *SpecialFormsNode) parseIf(tok lex.Token, elements []Node) (Node, bool) {
	// elements: [condition, then] or [condition, then, else]
	if len(elements) < 2 || len(elements) > 3 {
		return nil, false
	}

	condition := elements[0]
	then := elements[1]
	var elseExpr Node
	if len(elements) == 3 {
		elseExpr = elements[2]
	}

	return NewIfNode(tok, condition, then, elseExpr), true
}

// parseBegin handles (begin expr ...) expressions.
func (s *SpecialFormsNode) parseBegin(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) < 1 {
		return nil, false
	}
	return NewBeginNode(tok, elements), true
}

// parseDefine handles (define name value) and (define (name args) body) expressions.
func (s *SpecialFormsNode) parseDefine(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) < 2 {
		return nil, false
	}

	first := elements[0]

	// Check for (define name value) form
	if sym, ok := first.(*SymbolNode); ok {
		if len(elements) != 2 {
			return nil, false
		}
		return NewDefineNode(tok, sym, elements[1]), true
	}

	// Check for (define (name args...) body...) form
	nameElements, ok := extractListElements(first)
	if !ok || len(nameElements) == 0 {
		return nil, false
	}

	funcName, ok := nameElements[0].(*SymbolNode)
	if !ok {
		return nil, false
	}

	// Initialize params as empty slice (not nil) to distinguish from simple define
	params := make([]*SymbolNode, 0, len(nameElements)-1)
	for _, elem := range nameElements[1:] {
		param, ok := elem.(*SymbolNode)
		if !ok {
			return nil, false
		}
		params = append(params, param)
	}

	body := elements[1:]
	if len(body) == 0 {
		return nil, false
	}

	return NewDefineFuncNode(tok, funcName, params, body), true
}

// parseSet handles (set! name value) expressions.
func (s *SpecialFormsNode) parseSet(tok lex.Token, elements []Node) (Node, bool) {
	// elements: [name, value]
	if len(elements) != 2 {
		return nil, false
	}

	name, ok := elements[0].(*SymbolNode)
	if !ok {
		return nil, false
	}

	return NewSetNode(tok, name, elements[1]), true
}

// parseLambda handles (lambda (params) body) expressions.
func (s *SpecialFormsNode) parseLambda(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) < 2 {
		return nil, false
	}

	paramElements, ok := extractListElements(elements[0])
	if !ok {
		return nil, false
	}

	var params []*SymbolNode
	for _, elem := range paramElements {
		param, ok := elem.(*SymbolNode)
		if !ok {
			return nil, false
		}
		params = append(params, param)
	}

	body := elements[1:]
	return NewLambdaNode(tok, params, body), true
}

// parseLet handles (let ((name value) ...) body) expressions.
func (s *SpecialFormsNode) parseLet(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) < 2 {
		return nil, false
	}

	bindingsElements, ok := extractListElements(elements[0])
	if !ok {
		return nil, false
	}

	var bindings []Binding
	for _, elem := range bindingsElements {
		bindingElements, ok := extractListElements(elem)
		if !ok || len(bindingElements) != 2 {
			return nil, false
		}

		name, ok := bindingElements[0].(*SymbolNode)
		if !ok {
			return nil, false
		}

		bindings = append(bindings, Binding{
			Name:  name,
			Value: bindingElements[1],
		})
	}

	body := elements[1:]
	return NewLetNode(tok, bindings, body), true
}

// parseCond handles (cond (test expr) ... [(else expr)]) expressions.
func (s *SpecialFormsNode) parseCond(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) < 1 {
		return nil, false
	}

	var clauses []CondClause
	for _, elem := range elements {
		clauseElements, ok := extractListElements(elem)
		if !ok || len(clauseElements) < 2 {
			return nil, false
		}

		if sym, ok := clauseElements[0].(*SymbolNode); ok && sym.Name == "else" {
			clauses = append(clauses, CondClause{
				Condition: nil,
				Body:      clauseElements[1:],
				IsElse:    true,
			})
		} else {
			clauses = append(clauses, CondClause{
				Condition: clauseElements[0],
				Body:      clauseElements[1:],
				IsElse:    false,
			})
		}
	}

	return NewCondNode(tok, clauses), true
}

// parseQuote handles (quote expr) expressions.
func (s *SpecialFormsNode) parseQuote(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) != 1 {
		return nil, false
	}
	return NewQuoteNode(tok, elements[0]), true
}

// parseQuasiquote handles (quasiquote expr) expressions.
func (s *SpecialFormsNode) parseQuasiquote(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) != 1 {
		return nil, false
	}
	return NewQuasiquoteNode(tok, elements[0]), true
}

// parseUnquote handles (unquote expr) expressions.
func (s *SpecialFormsNode) parseUnquote(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) != 1 {
		return nil, false
	}
	return NewUnquoteNode(tok, elements[0]), true
}

// parseUnquoteSplicing handles (unquote-splicing expr) expressions.
func (s *SpecialFormsNode) parseUnquoteSplicing(tok lex.Token, elements []Node) (Node, bool) {
	if len(elements) != 1 {
		return nil, false
	}
	return NewUnquoteSplicingNode(tok, elements[0]), true
}

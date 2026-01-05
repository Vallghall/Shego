package ast

import (
	"github.com/Vallghall/schego/pkg/lex"
)

// Node is the interface implemented by all AST nodes.
type Node interface {
	node() // marker method to ensure type safety
	// Token returns the primary token associated with this node for error reporting.
	Token() lex.Token
}

// baseNode provides common functionality for AST nodes.
type baseNode struct {
	token lex.Token
}

func (b *baseNode) Token() lex.Token { return b.token }

// =============================================================================
// Atom types (self-evaluating)
// =============================================================================

// NumberNode represents a numeric literal.
type NumberNode struct {
	baseNode
	Value string
}

func (*NumberNode) node() {}

// NewNumberNode creates a new NumberNode from a token.
func NewNumberNode(tok lex.Token) *NumberNode {
	return &NumberNode{
		baseNode: baseNode{token: tok},
		Value:    tok.Raw(),
	}
}

// StringNode represents a string literal.
type StringNode struct {
	baseNode
	Value string
}

func (*StringNode) node() {}

// NewStringNode creates a new StringNode from a token.
func NewStringNode(tok lex.Token) *StringNode {
	// Strip the surrounding quotes from raw value
	raw := tok.Raw()
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	return &StringNode{
		baseNode: baseNode{token: tok},
		Value:    raw,
	}
}

// SymbolNode represents an identifier/symbol.
type SymbolNode struct {
	baseNode
	Name string
}

func (*SymbolNode) node() {}

// NewSymbolNode creates a new SymbolNode from a token.
func NewSymbolNode(tok lex.Token) *SymbolNode {
	return &SymbolNode{
		baseNode: baseNode{token: tok},
		Name:     tok.Raw(),
	}
}

// =============================================================================
// List and call expressions
// =============================================================================

// ListNode represents a raw list (before special form transformation).
type ListNode struct {
	baseNode
	Elements []Node
}

func (*ListNode) node() {}

// NewListNode creates a new ListNode.
func NewListNode(openParen lex.Token, elements []Node) *ListNode {
	return &ListNode{
		baseNode: baseNode{token: openParen},
		Elements: elements,
	}
}

// CallNode represents a function/procedure call.
type CallNode struct {
	baseNode
	Operator Node
	Args     []Node
}

func (*CallNode) node() {}

// NewCallNode creates a new CallNode.
func NewCallNode(operator Node, args []Node) *CallNode {
	return &CallNode{
		baseNode: baseNode{token: operator.Token()},
		Operator: operator,
		Args:     args,
	}
}

// =============================================================================
// Special forms
// =============================================================================

// IfNode represents an (if condition then else) expression.
type IfNode struct {
	baseNode
	Condition Node
	Then      Node
	Else      Node // may be nil
}

func (*IfNode) node() {}

// NewIfNode creates a new IfNode.
func NewIfNode(tok lex.Token, condition, then, elseExpr Node) *IfNode {
	return &IfNode{
		baseNode:  baseNode{token: tok},
		Condition: condition,
		Then:      then,
		Else:      elseExpr,
	}
}

// BeginNode represents a (begin expr ...) expression.
type BeginNode struct {
	baseNode
	Body []Node
}

func (*BeginNode) node() {}

// NewBeginNode creates a new BeginNode.
func NewBeginNode(tok lex.Token, body []Node) *BeginNode {
	return &BeginNode{
		baseNode: baseNode{token: tok},
		Body:     body,
	}
}

// DefineNode represents a (define name value) or (define (name args) body) expression.
type DefineNode struct {
	baseNode
	Name   *SymbolNode
	Params []*SymbolNode // nil for simple define, populated for function define
	Body   []Node        // single element for simple define, multiple for function
}

func (*DefineNode) node() {}

// NewDefineNode creates a simple (define name value) node.
func NewDefineNode(tok lex.Token, name *SymbolNode, value Node) *DefineNode {
	return &DefineNode{
		baseNode: baseNode{token: tok},
		Name:     name,
		Params:   nil,
		Body:     []Node{value},
	}
}

// NewDefineFuncNode creates a (define (name args) body) node.
func NewDefineFuncNode(tok lex.Token, name *SymbolNode, params []*SymbolNode, body []Node) *DefineNode {
	return &DefineNode{
		baseNode: baseNode{token: tok},
		Name:     name,
		Params:   params,
		Body:     body,
	}
}

// LambdaNode represents a (lambda (params) body) expression.
type LambdaNode struct {
	baseNode
	Params []*SymbolNode
	Body   []Node
}

func (*LambdaNode) node() {}

// NewLambdaNode creates a new LambdaNode.
func NewLambdaNode(tok lex.Token, params []*SymbolNode, body []Node) *LambdaNode {
	return &LambdaNode{
		baseNode: baseNode{token: tok},
		Params:   params,
		Body:     body,
	}
}

// Binding represents a single (name value) binding in let expressions.
type Binding struct {
	Name  *SymbolNode
	Value Node
}

// LetNode represents a (let ((name value) ...) body) expression.
type LetNode struct {
	baseNode
	Bindings []Binding
	Body     []Node
}

func (*LetNode) node() {}

// NewLetNode creates a new LetNode.
func NewLetNode(tok lex.Token, bindings []Binding, body []Node) *LetNode {
	return &LetNode{
		baseNode: baseNode{token: tok},
		Bindings: bindings,
		Body:     body,
	}
}

// CondClause represents a single clause in a cond expression.
type CondClause struct {
	Condition Node // nil for else clause
	Body      []Node
	IsElse    bool
}

// CondNode represents a (cond (test expr) ... (else expr)) expression.
type CondNode struct {
	baseNode
	Clauses []CondClause
}

func (*CondNode) node() {}

// NewCondNode creates a new CondNode.
func NewCondNode(tok lex.Token, clauses []CondClause) *CondNode {
	return &CondNode{
		baseNode: baseNode{token: tok},
		Clauses:  clauses,
	}
}

// QuoteNode represents a (quote expr) or 'expr expression.
type QuoteNode struct {
	baseNode
	Value Node
}

func (*QuoteNode) node() {}

// NewQuoteNode creates a new QuoteNode.
func NewQuoteNode(tok lex.Token, value Node) *QuoteNode {
	return &QuoteNode{
		baseNode: baseNode{token: tok},
		Value:    value,
	}
}

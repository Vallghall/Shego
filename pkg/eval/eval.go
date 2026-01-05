package eval

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/atom"
	"github.com/Vallghall/schego/pkg/mem"
)

// Evaluator evaluates AST nodes and manages interpreter state.
type Evaluator interface {
	// State returns the interpreter state.
	State() *State

	// Eval evaluates a sequence of AST nodes and returns the last result.
	Eval(nodes []ast.Node) (mem.Object, error)

	// EvalOne evaluates a single AST node.
	EvalOne(node ast.Node) (mem.Object, error)
}

// evaluator is the concrete implementation of Evaluator.
type evaluator struct {
	state *State
}

// New creates a new Evaluator with empty state.
func New() Evaluator {
	return &evaluator{
		state: NewState(),
	}
}

// NewWithState creates a new Evaluator with the given state.
func NewWithState(state *State) Evaluator {
	return &evaluator{
		state: state,
	}
}

// State returns the interpreter state.
func (e *evaluator) State() *State {
	return e.state
}

// Eval evaluates a sequence of AST nodes and returns the last result.
func (e *evaluator) Eval(nodes []ast.Node) (mem.Object, error) {
	var result mem.Object = mem.Void

	for _, node := range nodes {
		var err error
		result, err = e.EvalOne(node)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// EvalOne evaluates a single AST node.
func (e *evaluator) EvalOne(node ast.Node) (mem.Object, error) {
	switch n := node.(type) {
	case *ast.NumberNode:
		return e.evalNumber(n)
	case *ast.StringNode:
		return e.evalString(n)
	case *ast.SymbolNode:
		return e.evalSymbol(n)
	case *ast.ListNode:
		return e.evalList(n)
	case *ast.CallNode:
		return e.evalCall(n)
	case *ast.IfNode:
		return e.evalIf(n)
	case *ast.BeginNode:
		return e.evalBegin(n)
	case *ast.DefineNode:
		return e.evalDefine(n)
	case *ast.LambdaNode:
		return e.evalLambda(n)
	case *ast.LetNode:
		return e.evalLet(n)
	case *ast.CondNode:
		return e.evalCond(n)
	case *ast.QuoteNode:
		return e.evalQuote(n)
	default:
		return nil, fmt.Errorf("unknown AST node type: %T", node)
	}
}

// =============================================================================
// Self-Evaluating Forms
// =============================================================================

func (e *evaluator) evalNumber(n *ast.NumberNode) (mem.Object, error) {
	// Parse the number string
	var value float64
	_, err := fmt.Sscanf(n.Value, "%f", &value)
	if err != nil {
		return nil, mem.NewRuntimeError(fmt.Sprintf("invalid number: %s", n.Value))
	}
	return mem.NewNumber(value), nil
}

func (e *evaluator) evalString(n *ast.StringNode) (mem.Object, error) {
	return mem.NewString(n.Value), nil
}

func (e *evaluator) evalSymbol(n *ast.SymbolNode) (mem.Object, error) {
	// Handle special symbols
	switch n.Name {
	case "#t":
		return mem.True, nil
	case "#f":
		return mem.False, nil
	}

	// Look up the symbol in the current context
	id := e.state.Intern(n.Name)
	value, err := e.state.LookupOrError(id)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (e *evaluator) evalList(n *ast.ListNode) (mem.Object, error) {
	// Empty list evaluates to nil
	if len(n.Elements) == 0 {
		return mem.Nil, nil
	}
	// Non-empty list is an error (should have been parsed as CallNode)
	return nil, mem.NewRuntimeError("unexpected list in evaluation")
}

// =============================================================================
// Procedure Call
// =============================================================================

func (e *evaluator) evalCall(n *ast.CallNode) (mem.Object, error) {
	// Evaluate the operator
	opValue, err := e.EvalOne(n.Operator)
	if err != nil {
		return nil, err
	}

	// Check if callable
	proc, err := mem.AsProcedure(opValue)
	if err != nil {
		return nil, mem.NewNotCallableError(opValue)
	}

	// Evaluate arguments
	args := make([]mem.Object, len(n.Args))
	for i, arg := range n.Args {
		args[i], err = e.EvalOne(arg)
		if err != nil {
			return nil, err
		}
	}

	// Call the procedure
	return e.apply(proc, args)
}

func (e *evaluator) apply(proc *mem.Procedure, args []mem.Object) (mem.Object, error) {
	// Check arity
	if proc.Arity() >= 0 && len(args) != proc.Arity() {
		return nil, mem.NewArityError(proc.Name(), proc.Arity(), len(args))
	}

	switch proc.Kind() {
	case mem.ProcPrimitive:
		return e.applyPrimitive(proc, args)
	case mem.ProcLambda:
		return e.applyLambda(proc, args)
	default:
		return nil, mem.NewRuntimeError("unknown procedure kind")
	}
}

func (e *evaluator) applyPrimitive(proc *mem.Procedure, args []mem.Object) (mem.Object, error) {
	// The native function is stored as mem.PrimitiveFunc
	fn, ok := proc.Native().(mem.PrimitiveFunc)
	if !ok {
		return nil, mem.NewRuntimeError("invalid primitive function")
	}
	return fn(args)
}

func (e *evaluator) applyLambda(proc *mem.Procedure, args []mem.Object) (mem.Object, error) {
	// Create new context extending the closure's captured context
	ctx, err := mem.ExtendWith(proc.Context(), proc.Params(), args)
	if err != nil {
		return nil, err
	}

	// Evaluate body in the new context
	var result mem.Object = mem.Void
	err = e.state.WithContext(ctx, func() error {
		body := proc.Body()
		for _, expr := range body {
			astNode, ok := expr.(ast.Node)
			if !ok {
				return mem.NewRuntimeError("invalid body expression")
			}
			var evalErr error
			result, evalErr = e.EvalOne(astNode)
			if evalErr != nil {
				return evalErr
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// =============================================================================
// Special Forms
// =============================================================================

func (e *evaluator) evalIf(n *ast.IfNode) (mem.Object, error) {
	// Evaluate condition
	cond, err := e.EvalOne(n.Condition)
	if err != nil {
		return nil, err
	}

	// In Scheme, only #f is falsy
	if cond.IsTruthy() {
		return e.EvalOne(n.Then)
	}

	// Evaluate else branch if present
	if n.Else != nil {
		return e.EvalOne(n.Else)
	}

	return mem.Void, nil
}

func (e *evaluator) evalBegin(n *ast.BeginNode) (mem.Object, error) {
	var result mem.Object = mem.Void

	for _, expr := range n.Body {
		var err error
		result, err = e.EvalOne(expr)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (e *evaluator) evalDefine(n *ast.DefineNode) (mem.Object, error) {
	nameID := e.state.Intern(n.Name.Name)

	if n.Params == nil {
		// Simple define: (define name value)
		value, err := e.EvalOne(n.Body[0])
		if err != nil {
			return nil, err
		}
		if err := e.state.Define(nameID, value); err != nil {
			return nil, err
		}
	} else {
		// Function define: (define (name params...) body...)
		params := make([]atom.ID, len(n.Params))
		for i, p := range n.Params {
			params[i] = e.state.Intern(p.Name)
		}

		// Convert body to []any for storage
		body := make([]any, len(n.Body))
		for i, b := range n.Body {
			body[i] = b
		}

		proc := mem.NewLambda(n.Name.Name, params, body, e.state.Current())
		if err := e.state.Define(nameID, proc); err != nil {
			return nil, err
		}
	}

	return mem.Void, nil
}

func (e *evaluator) evalLambda(n *ast.LambdaNode) (mem.Object, error) {
	// Capture parameters
	params := make([]atom.ID, len(n.Params))
	for i, p := range n.Params {
		params[i] = e.state.Intern(p.Name)
	}

	// Convert body to []any for storage
	body := make([]any, len(n.Body))
	for i, b := range n.Body {
		body[i] = b
	}

	// Create lambda capturing current context
	return mem.NewLambda("", params, body, e.state.Current()), nil
}

func (e *evaluator) evalLet(n *ast.LetNode) (mem.Object, error) {
	// Evaluate all binding values in current context
	names := make([]atom.ID, len(n.Bindings))
	values := make([]mem.Object, len(n.Bindings))

	for i, binding := range n.Bindings {
		names[i] = e.state.Intern(binding.Name.Name)
		var err error
		values[i], err = e.EvalOne(binding.Value)
		if err != nil {
			return nil, err
		}
	}

	// Create new context with bindings
	ctx, err := e.state.ExtendCurrentWith(names, values)
	if err != nil {
		return nil, err
	}

	// Evaluate body in new context
	var result mem.Object = mem.Void
	err = e.state.WithContext(ctx, func() error {
		for _, expr := range n.Body {
			var evalErr error
			result, evalErr = e.EvalOne(expr)
			if evalErr != nil {
				return evalErr
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *evaluator) evalCond(n *ast.CondNode) (mem.Object, error) {
	for _, clause := range n.Clauses {
		if clause.IsElse {
			// Else clause - always execute
			return e.evalCondBody(clause.Body)
		}

		// Evaluate condition
		cond, err := e.EvalOne(clause.Condition)
		if err != nil {
			return nil, err
		}

		if cond.IsTruthy() {
			return e.evalCondBody(clause.Body)
		}
	}

	// No clause matched
	return mem.Void, nil
}

func (e *evaluator) evalCondBody(body []ast.Node) (mem.Object, error) {
	var result mem.Object = mem.Void
	for _, expr := range body {
		var err error
		result, err = e.EvalOne(expr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (e *evaluator) evalQuote(n *ast.QuoteNode) (mem.Object, error) {
	// Convert AST to Scheme object without evaluating
	return e.astToObject(n.Value)
}

// astToObject converts an AST node to a Scheme object (for quote).
func (e *evaluator) astToObject(node ast.Node) (mem.Object, error) {
	switch n := node.(type) {
	case *ast.NumberNode:
		var value float64
		fmt.Sscanf(n.Value, "%f", &value)
		return mem.NewNumber(value), nil
	case *ast.StringNode:
		return mem.NewString(n.Value), nil
	case *ast.SymbolNode:
		id := e.state.Intern(n.Name)
		return mem.NewSymbol(id, e.state.Pool()), nil
	case *ast.ListNode:
		if len(n.Elements) == 0 {
			return mem.Nil, nil
		}
		elements := make([]mem.Object, len(n.Elements))
		for i, elem := range n.Elements {
			var err error
			elements[i], err = e.astToObject(elem)
			if err != nil {
				return nil, err
			}
		}
		return mem.SliceToList(elements), nil
	case *ast.CallNode:
		// In quote, a call node is just a list
		elements := make([]mem.Object, 1+len(n.Args))
		var err error
		elements[0], err = e.astToObject(n.Operator)
		if err != nil {
			return nil, err
		}
		for i, arg := range n.Args {
			elements[i+1], err = e.astToObject(arg)
			if err != nil {
				return nil, err
			}
		}
		return mem.SliceToList(elements), nil
	default:
		return nil, mem.NewRuntimeError(fmt.Sprintf("cannot quote: %T", node))
	}
}



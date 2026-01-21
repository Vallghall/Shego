package eval

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/atom"
	"github.com/Vallghall/schego/pkg/core"
	"github.com/Vallghall/schego/pkg/core/builtin"
	"github.com/Vallghall/schego/pkg/lex"
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

// New creates a new Evaluator with builtins loaded.
func New() (Evaluator, error) {
	state := NewState()
	// Load core builtins automatically
	if err := state.LoadBuiltins(core.All); err != nil {
		return nil, fmt.Errorf("failed to load builtins: %w", err)
	}

	// Create evaluator
	eval := &evaluator{
		state: state,
	}

	// Register eval, apply, and load primitives (these need evaluator access)
	if err := eval.registerEvalPrimitives(); err != nil {
		return nil, fmt.Errorf("failed to register eval primitives: %w", err)
	}

	// Load functional.scm (map, filter, fold) before eval/apply are available
	// We need to load it directly using Go code
	if err := eval.loadFunctionalLibrary(); err != nil {
		return nil, fmt.Errorf("failed to load functional library: %w", err)
	}

	return eval, nil
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
	case *ast.SetNode:
		return e.evalSet(n)
	case *ast.LambdaNode:
		return e.evalLambda(n)
	case *ast.LetNode:
		return e.evalLet(n)
	case *ast.CondNode:
		return e.evalCond(n)
	case *ast.QuoteNode:
		return e.evalQuote(n)
	case *ast.QuasiquoteNode:
		return e.evalQuasiquote(n)
	case *ast.UnquoteNode:
		return nil, mem.NewRuntimeError("unquote outside of quasiquote")
	case *ast.UnquoteSplicingNode:
		return nil, mem.NewRuntimeError("unquote-splicing outside of quasiquote")
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

func (e *evaluator) evalSet(n *ast.SetNode) (mem.Object, error) {
	nameID := e.state.Intern(n.Name.Name)

	// Evaluate the new value
	value, err := e.EvalOne(n.Value)
	if err != nil {
		return nil, err
	}

	// Set the binding (this will error if not defined)
	if err := e.state.Set(nameID, value); err != nil {
		return nil, err
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

func (e *evaluator) evalQuasiquote(n *ast.QuasiquoteNode) (mem.Object, error) {
	// Convert AST to Scheme object, but evaluate unquote expressions
	return e.quasiquoteToObject(n.Value)
}

// quasiquoteToObject converts an AST node to a Scheme object for quasiquote.
// Unlike astToObject, this evaluates unquote expressions.
func (e *evaluator) quasiquoteToObject(node ast.Node) (mem.Object, error) {
	switch n := node.(type) {
	case *ast.UnquoteNode:
		// Evaluate the unquoted expression
		return e.EvalOne(n.Value)
	case *ast.UnquoteSplicingNode:
		// unquote-splicing at top level is an error
		return nil, mem.NewRuntimeError("unquote-splicing in non-list context")
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
		return e.quasiquoteList(n.Elements)
	case *ast.CallNode:
		// In quasiquote, a call node is just a list
		// Combine operator and args into a single list
		allElements := make([]ast.Node, 1+len(n.Args))
		allElements[0] = n.Operator
		copy(allElements[1:], n.Args)
		return e.quasiquoteList(allElements)
	case *ast.QuasiquoteNode:
		// Nested quasiquote - just quote it
		return e.astToObject(node)
	case *ast.QuoteNode:
		// Quote inside quasiquote - just quote it
		return e.astToObject(node)
	default:
		return nil, mem.NewRuntimeError(fmt.Sprintf("cannot quasiquote: %T", node))
	}
}

// quasiquoteList handles list elements, supporting unquote-splicing.
func (e *evaluator) quasiquoteList(elements []ast.Node) (mem.Object, error) {
	var result []mem.Object

	for _, elem := range elements {
		// Check if this element is unquote-splicing
		if splice, ok := elem.(*ast.UnquoteSplicingNode); ok {
			// Evaluate the spliced expression
			value, err := e.EvalOne(splice.Value)
			if err != nil {
				return nil, err
			}

			// Convert to slice and append all elements
			items, err := mem.ListToSlice(value)
			if err != nil {
				return nil, mem.NewRuntimeError("unquote-splicing requires a list")
			}
			result = append(result, items...)
		} else {
			// Regular element - process normally
			obj, err := e.quasiquoteToObject(elem)
			if err != nil {
				return nil, err
			}
			result = append(result, obj)
		}
	}

	return mem.SliceToList(result), nil
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
	case *ast.QuoteNode:
		// Nested quote - return as (quote value)
		quoteSym := mem.NewSymbol(e.state.Intern("quote"), e.state.Pool())
		value, err := e.astToObject(n.Value)
		if err != nil {
			return nil, err
		}
		return mem.SliceToList([]mem.Object{quoteSym, value}), nil
	case *ast.QuasiquoteNode:
		// Quasiquote in quote - return as (quasiquote value)
		qqSym := mem.NewSymbol(e.state.Intern("quasiquote"), e.state.Pool())
		value, err := e.astToObject(n.Value)
		if err != nil {
			return nil, err
		}
		return mem.SliceToList([]mem.Object{qqSym, value}), nil
	case *ast.UnquoteNode:
		// Unquote in quote - return as (unquote value)
		uqSym := mem.NewSymbol(e.state.Intern("unquote"), e.state.Pool())
		value, err := e.astToObject(n.Value)
		if err != nil {
			return nil, err
		}
		return mem.SliceToList([]mem.Object{uqSym, value}), nil
	case *ast.UnquoteSplicingNode:
		// Unquote-splicing in quote - return as (unquote-splicing value)
		uqsSym := mem.NewSymbol(e.state.Intern("unquote-splicing"), e.state.Pool())
		value, err := e.astToObject(n.Value)
		if err != nil {
			return nil, err
		}
		return mem.SliceToList([]mem.Object{uqsSym, value}), nil
	default:
		return nil, mem.NewRuntimeError(fmt.Sprintf("cannot quote: %T", node))
	}
}

// objectToAST converts a Scheme object to an AST node (for eval).
func (e *evaluator) objectToAST(obj mem.Object) (ast.Node, error) {
	switch obj.Type() {
	case mem.TypeNumber:
		num := obj.(*mem.Number)
		// Create a token for the number
		tok := lex.NewToken("", 0, 1, 1, fmt.Sprintf("%g", num.Value()), lex.Number)
		return ast.NewNumberNode(tok), nil

	case mem.TypeString:
		str := obj.(*mem.String)
		// Create a token for the string (with quotes)
		tok := lex.NewToken("", 0, 1, 1, fmt.Sprintf("%q", str.Value()), lex.String)
		return ast.NewStringNode(tok), nil

	case mem.TypeSymbol:
		sym := obj.(*mem.Symbol)
		name := sym.Name()
		// Create a token for the symbol
		tok := lex.NewToken("", 0, 1, 1, name, lex.Atom)
		return ast.NewSymbolNode(tok), nil

	case mem.TypeNil:
		// Empty list -> ListNode with no elements
		tok := lex.NewToken("", 0, 1, 1, "()", lex.ParenOpen)
		return ast.NewListNode(tok, nil), nil

	case mem.TypePair:
		// Convert list to AST
		elements, err := mem.ListToSlice(obj)
		if err != nil {
			return nil, mem.NewRuntimeError(fmt.Sprintf("eval: cannot convert to list: %v", err))
		}

		if len(elements) == 0 {
			tok := lex.NewToken("", 0, 1, 1, "()", lex.ParenOpen)
			return ast.NewListNode(tok, nil), nil
		}

		// Convert elements to AST nodes
		astElements := make([]ast.Node, len(elements))
		for i, elem := range elements {
			astElem, err := e.objectToAST(elem)
			if err != nil {
				return nil, err
			}
			astElements[i] = astElem
		}

		// If first element is a symbol, treat as CallNode
		if firstSym, ok := astElements[0].(*ast.SymbolNode); ok {
			return ast.NewCallNode(firstSym, astElements[1:]), nil
		}

		// Otherwise, treat as ListNode
		tok := lex.NewToken("", 0, 1, 1, "(", lex.ParenOpen)
		return ast.NewListNode(tok, astElements), nil

	default:
		return nil, mem.NewRuntimeError(fmt.Sprintf("eval: cannot convert %s to AST", obj.Type()))
	}
}

// registerEvalPrimitives registers eval, apply, and load primitives with evaluator access.
func (e *evaluator) registerEvalPrimitives() error {
	// Register eval primitive
	evalDefs := builtin.EvalDefinitions(func(obj mem.Object) (mem.Object, error) {
		// Convert object to AST
		astNode, err := e.objectToAST(obj)
		if err != nil {
			return nil, err
		}
		// Evaluate the AST node
		return e.EvalOne(astNode)
	})
	if err := e.state.LoadDefinitions(evalDefs); err != nil {
		return fmt.Errorf("failed to register eval: %w", err)
	}

	// Register apply primitive
	applyDefs := builtin.ApplyDefinitions(func(proc *mem.Procedure, args []mem.Object) (mem.Object, error) {
		return e.apply(proc, args)
	})
	if err := e.state.LoadDefinitions(applyDefs); err != nil {
		return fmt.Errorf("failed to register apply: %w", err)
	}

	// Register load primitive
	loadDefs := builtin.LoadDefinitions(func(source string) (mem.Object, error) {
		// Tokenize
		lexer := lex.New()
		tokens, err := lexer.Tokenize(source)
		if err != nil {
			return nil, fmt.Errorf("load: lexer error: %w", err)
		}

		// Parse
		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		if err != nil {
			return nil, fmt.Errorf("load: parser error: %w", err)
		}

		// Evaluate
		return e.Eval(nodes)
	})
	if err := e.state.LoadDefinitions(loadDefs); err != nil {
		return fmt.Errorf("failed to register load: %w", err)
	}

	return nil
}

// loadFunctionalLibrary loads the functional.scm file containing map, filter, fold.
func (e *evaluator) loadFunctionalLibrary() error {
	// Try multiple paths to find functional.scm
	paths := []string{
		filepath.Join("core", "builtin", "functional.scm"),
		"functional.scm",
		filepath.Join("..", "core", "builtin", "functional.scm"),
		filepath.Join(".", "core", "builtin", "functional.scm"),
	}

	var content []byte
	var err error
	var found bool

	for _, path := range paths {
		content, err = os.ReadFile(path)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		// If not found, it's not critical - map, filter, fold just won't be available
		// Return nil to allow evaluator to work without them
		return nil
	}

	// Tokenize, parse, and evaluate
	lexer := lex.New()
	tokens, err := lexer.Tokenize(string(content))
	if err != nil {
		return fmt.Errorf("failed to tokenize functional.scm: %w", err)
	}

	parser := ast.New()
	nodes, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Errorf("failed to parse functional.scm: %w", err)
	}

	_, err = e.Eval(nodes)
	if err != nil {
		return fmt.Errorf("failed to evaluate functional.scm: %w", err)
	}

	return nil
}

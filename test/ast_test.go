package test

import (
	"testing"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/lex"
)

func TestParser(t *testing.T) {
	t.Run("Atoms", func(t *testing.T) {
		testAtoms(t)
	})

	t.Run("SimpleLists", func(t *testing.T) {
		testSimpleLists(t)
	})

	t.Run("NestedLists", func(t *testing.T) {
		testNestedLists(t)
	})

	t.Run("IfExpression", func(t *testing.T) {
		testIfExpression(t)
	})

	t.Run("BeginExpression", func(t *testing.T) {
		testBeginExpression(t)
	})

	t.Run("DefineExpression", func(t *testing.T) {
		testDefineExpression(t)
	})

	t.Run("LambdaExpression", func(t *testing.T) {
		testLambdaExpression(t)
	})

	t.Run("LetExpression", func(t *testing.T) {
		testLetExpression(t)
	})

	t.Run("CondExpression", func(t *testing.T) {
		testCondExpression(t)
	})

	t.Run("QuoteExpression", func(t *testing.T) {
		testQuoteExpression(t)
	})

	t.Run("FunctionCalls", func(t *testing.T) {
		testFunctionCalls(t)
	})
}

func parseString(t *testing.T, input string) []ast.Node {
	t.Helper()
	lexer := lex.New()
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		t.Fatalf("lexer error: %v", err)
	}

	parser := ast.New()
	nodes, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	return nodes
}

func testAtoms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		nodeType string
		value    string
	}{
		{"number", "42", "*ast.NumberNode", "42"},
		{"negative number", "-3", "*ast.NumberNode", "-3"},
		{"decimal", "3.14", "*ast.NumberNode", "3.14"},
		{"string", `"hello"`, "*ast.StringNode", "hello"},
		{"symbol", "foo", "*ast.SymbolNode", "foo"},
		{"plus operator", "+", "*ast.SymbolNode", "+"},
		{"minus operator", "-", "*ast.SymbolNode", "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := parseString(t, tt.input)
			if len(nodes) != 1 {
				t.Fatalf("expected 1 node, got %d", len(nodes))
			}

			nodeType := getTypeName(nodes[0])
			if nodeType != tt.nodeType {
				t.Errorf("expected type %s, got %s", tt.nodeType, nodeType)
			}

			switch n := nodes[0].(type) {
			case *ast.NumberNode:
				if n.Value != tt.value {
					t.Errorf("expected value %q, got %q", tt.value, n.Value)
				}
			case *ast.StringNode:
				if n.Value != tt.value {
					t.Errorf("expected value %q, got %q", tt.value, n.Value)
				}
			case *ast.SymbolNode:
				if n.Name != tt.value {
					t.Errorf("expected name %q, got %q", tt.value, n.Name)
				}
			}
		})
	}
}

func testSimpleLists(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		nodes := parseString(t, "()")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		_, ok := nodes[0].(*ast.ListNode)
		if !ok {
			t.Errorf("expected ListNode, got %T", nodes[0])
		}
	})

	t.Run("list with atoms", func(t *testing.T) {
		nodes := parseString(t, "(1 2 3)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		call, ok := nodes[0].(*ast.CallNode)
		if !ok {
			t.Fatalf("expected CallNode, got %T", nodes[0])
		}
		if len(call.Args) != 2 {
			t.Errorf("expected 2 args, got %d", len(call.Args))
		}
	})
}

func testNestedLists(t *testing.T) {
	t.Run("nested list", func(t *testing.T) {
		nodes := parseString(t, "((1 2) 3)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		call, ok := nodes[0].(*ast.CallNode)
		if !ok {
			t.Fatalf("expected CallNode, got %T", nodes[0])
		}
		_, ok = call.Operator.(*ast.CallNode)
		if !ok {
			t.Errorf("expected nested CallNode as operator, got %T", call.Operator)
		}
	})

	t.Run("deeply nested", func(t *testing.T) {
		nodes := parseString(t, "(((a)))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		// (((a))) -> CallNode { Operator: ((a)) }
		// ((a)) -> CallNode { Operator: (a) }
		// (a) -> CallNode { Operator: a (SymbolNode) }
		call1, ok := nodes[0].(*ast.CallNode)
		if !ok {
			t.Fatalf("expected CallNode, got %T", nodes[0])
		}
		call2, ok := call1.Operator.(*ast.CallNode)
		if !ok {
			t.Fatalf("expected nested CallNode, got %T", call1.Operator)
		}
		call3, ok := call2.Operator.(*ast.CallNode)
		if !ok {
			t.Fatalf("expected innermost CallNode, got %T", call2.Operator)
		}
		_, ok = call3.Operator.(*ast.SymbolNode)
		if !ok {
			t.Errorf("expected SymbolNode as operator of innermost call, got %T", call3.Operator)
		}
	})
}

func testIfExpression(t *testing.T) {
	t.Run("if with else", func(t *testing.T) {
		nodes := parseString(t, "(if #t 1 2)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		ifNode, ok := nodes[0].(*ast.IfNode)
		if !ok {
			t.Fatalf("expected IfNode, got %T", nodes[0])
		}
		if ifNode.Condition == nil {
			t.Error("condition is nil")
		}
		if ifNode.Then == nil {
			t.Error("then is nil")
		}
		if ifNode.Else == nil {
			t.Error("else is nil")
		}
	})

	t.Run("if without else", func(t *testing.T) {
		nodes := parseString(t, "(if #t 1)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		ifNode, ok := nodes[0].(*ast.IfNode)
		if !ok {
			t.Fatalf("expected IfNode, got %T", nodes[0])
		}
		if ifNode.Else != nil {
			t.Error("else should be nil")
		}
	})

	t.Run("nested if", func(t *testing.T) {
		nodes := parseString(t, "(if (if #t #f) 1 2)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		ifNode, ok := nodes[0].(*ast.IfNode)
		if !ok {
			t.Fatalf("expected IfNode, got %T", nodes[0])
		}
		_, ok = ifNode.Condition.(*ast.IfNode)
		if !ok {
			t.Errorf("expected nested IfNode as condition, got %T", ifNode.Condition)
		}
	})
}

func testBeginExpression(t *testing.T) {
	t.Run("begin with single expression", func(t *testing.T) {
		nodes := parseString(t, "(begin 1)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		beginNode, ok := nodes[0].(*ast.BeginNode)
		if !ok {
			t.Fatalf("expected BeginNode, got %T", nodes[0])
		}
		if len(beginNode.Body) != 1 {
			t.Errorf("expected 1 body expression, got %d", len(beginNode.Body))
		}
	})

	t.Run("begin with multiple expressions", func(t *testing.T) {
		nodes := parseString(t, "(begin 1 2 3)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		beginNode, ok := nodes[0].(*ast.BeginNode)
		if !ok {
			t.Fatalf("expected BeginNode, got %T", nodes[0])
		}
		if len(beginNode.Body) != 3 {
			t.Errorf("expected 3 body expressions, got %d", len(beginNode.Body))
		}
	})
}

func testDefineExpression(t *testing.T) {
	t.Run("simple define", func(t *testing.T) {
		nodes := parseString(t, "(define x 42)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		defNode, ok := nodes[0].(*ast.DefineNode)
		if !ok {
			t.Fatalf("expected DefineNode, got %T", nodes[0])
		}
		if defNode.Name.Name != "x" {
			t.Errorf("expected name 'x', got %q", defNode.Name.Name)
		}
		if defNode.Params != nil {
			t.Error("expected nil params for simple define")
		}
	})

	t.Run("function define", func(t *testing.T) {
		nodes := parseString(t, "(define (square x) (* x x))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		defNode, ok := nodes[0].(*ast.DefineNode)
		if !ok {
			t.Fatalf("expected DefineNode, got %T", nodes[0])
		}
		if defNode.Name.Name != "square" {
			t.Errorf("expected name 'square', got %q", defNode.Name.Name)
		}
		if len(defNode.Params) != 1 {
			t.Errorf("expected 1 param, got %d", len(defNode.Params))
		}
		if defNode.Params[0].Name != "x" {
			t.Errorf("expected param 'x', got %q", defNode.Params[0].Name)
		}
	})

	t.Run("function define with multiple params", func(t *testing.T) {
		nodes := parseString(t, "(define (add a b) (+ a b))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		defNode, ok := nodes[0].(*ast.DefineNode)
		if !ok {
			t.Fatalf("expected DefineNode, got %T", nodes[0])
		}
		if len(defNode.Params) != 2 {
			t.Errorf("expected 2 params, got %d", len(defNode.Params))
		}
	})
}

func testLambdaExpression(t *testing.T) {
	t.Run("lambda with no params", func(t *testing.T) {
		nodes := parseString(t, "(lambda () 42)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		lambdaNode, ok := nodes[0].(*ast.LambdaNode)
		if !ok {
			t.Fatalf("expected LambdaNode, got %T", nodes[0])
		}
		if len(lambdaNode.Params) != 0 {
			t.Errorf("expected 0 params, got %d", len(lambdaNode.Params))
		}
	})

	t.Run("lambda with params", func(t *testing.T) {
		nodes := parseString(t, "(lambda (x y) (+ x y))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		lambdaNode, ok := nodes[0].(*ast.LambdaNode)
		if !ok {
			t.Fatalf("expected LambdaNode, got %T", nodes[0])
		}
		if len(lambdaNode.Params) != 2 {
			t.Errorf("expected 2 params, got %d", len(lambdaNode.Params))
		}
		if len(lambdaNode.Body) != 1 {
			t.Errorf("expected 1 body expression, got %d", len(lambdaNode.Body))
		}
	})

	t.Run("lambda with multiple body expressions", func(t *testing.T) {
		nodes := parseString(t, "(lambda (x) (display x) x)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		lambdaNode, ok := nodes[0].(*ast.LambdaNode)
		if !ok {
			t.Fatalf("expected LambdaNode, got %T", nodes[0])
		}
		if len(lambdaNode.Body) != 2 {
			t.Errorf("expected 2 body expressions, got %d", len(lambdaNode.Body))
		}
	})
}

func testLetExpression(t *testing.T) {
	t.Run("let with single binding", func(t *testing.T) {
		nodes := parseString(t, "(let ((x 1)) x)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		letNode, ok := nodes[0].(*ast.LetNode)
		if !ok {
			t.Fatalf("expected LetNode, got %T", nodes[0])
		}
		if len(letNode.Bindings) != 1 {
			t.Errorf("expected 1 binding, got %d", len(letNode.Bindings))
		}
		if letNode.Bindings[0].Name.Name != "x" {
			t.Errorf("expected binding name 'x', got %q", letNode.Bindings[0].Name.Name)
		}
	})

	t.Run("let with multiple bindings", func(t *testing.T) {
		nodes := parseString(t, "(let ((x 1) (y 2)) (+ x y))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		letNode, ok := nodes[0].(*ast.LetNode)
		if !ok {
			t.Fatalf("expected LetNode, got %T", nodes[0])
		}
		if len(letNode.Bindings) != 2 {
			t.Errorf("expected 2 bindings, got %d", len(letNode.Bindings))
		}
	})
}

func testCondExpression(t *testing.T) {
	t.Run("cond with clauses", func(t *testing.T) {
		nodes := parseString(t, "(cond ((= x 1) \"one\") ((= x 2) \"two\"))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		condNode, ok := nodes[0].(*ast.CondNode)
		if !ok {
			t.Fatalf("expected CondNode, got %T", nodes[0])
		}
		if len(condNode.Clauses) != 2 {
			t.Errorf("expected 2 clauses, got %d", len(condNode.Clauses))
		}
	})

	t.Run("cond with else", func(t *testing.T) {
		nodes := parseString(t, "(cond ((= x 1) \"one\") (else \"other\"))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		condNode, ok := nodes[0].(*ast.CondNode)
		if !ok {
			t.Fatalf("expected CondNode, got %T", nodes[0])
		}
		if len(condNode.Clauses) != 2 {
			t.Errorf("expected 2 clauses, got %d", len(condNode.Clauses))
		}
		if !condNode.Clauses[1].IsElse {
			t.Error("expected second clause to be else")
		}
	})
}

func testQuoteExpression(t *testing.T) {
	t.Run("quote symbol", func(t *testing.T) {
		nodes := parseString(t, "(quote foo)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		quoteNode, ok := nodes[0].(*ast.QuoteNode)
		if !ok {
			t.Fatalf("expected QuoteNode, got %T", nodes[0])
		}
		sym, ok := quoteNode.Value.(*ast.SymbolNode)
		if !ok {
			t.Fatalf("expected SymbolNode in quote, got %T", quoteNode.Value)
		}
		if sym.Name != "foo" {
			t.Errorf("expected 'foo', got %q", sym.Name)
		}
	})

	t.Run("quote list", func(t *testing.T) {
		nodes := parseString(t, "(quote (1 2 3))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		quoteNode, ok := nodes[0].(*ast.QuoteNode)
		if !ok {
			t.Fatalf("expected QuoteNode, got %T", nodes[0])
		}
		// Inside quote, lists are not transformed to CallNode
		_, ok = quoteNode.Value.(*ast.CallNode)
		if !ok {
			// It could also be a ListNode depending on implementation
			_, ok = quoteNode.Value.(*ast.ListNode)
		}
		if !ok {
			t.Errorf("expected CallNode or ListNode in quote, got %T", quoteNode.Value)
		}
	})
}

func testFunctionCalls(t *testing.T) {
	t.Run("simple call", func(t *testing.T) {
		nodes := parseString(t, "(+ 1 2)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		call, ok := nodes[0].(*ast.CallNode)
		if !ok {
			t.Fatalf("expected CallNode, got %T", nodes[0])
		}
		op, ok := call.Operator.(*ast.SymbolNode)
		if !ok {
			t.Fatalf("expected SymbolNode operator, got %T", call.Operator)
		}
		if op.Name != "+" {
			t.Errorf("expected operator '+', got %q", op.Name)
		}
		if len(call.Args) != 2 {
			t.Errorf("expected 2 args, got %d", len(call.Args))
		}
	})

	t.Run("nested calls", func(t *testing.T) {
		nodes := parseString(t, "(* (+ 1 2) (- 4 3))")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		call, ok := nodes[0].(*ast.CallNode)
		if !ok {
			t.Fatalf("expected CallNode, got %T", nodes[0])
		}
		if len(call.Args) != 2 {
			t.Errorf("expected 2 args, got %d", len(call.Args))
		}
		_, ok = call.Args[0].(*ast.CallNode)
		if !ok {
			t.Errorf("expected CallNode as first arg, got %T", call.Args[0])
		}
	})

	t.Run("call with no args", func(t *testing.T) {
		nodes := parseString(t, "(newline)")
		if len(nodes) != 1 {
			t.Fatalf("expected 1 node, got %d", len(nodes))
		}
		call, ok := nodes[0].(*ast.CallNode)
		if !ok {
			t.Fatalf("expected CallNode, got %T", nodes[0])
		}
		if len(call.Args) != 0 {
			t.Errorf("expected 0 args, got %d", len(call.Args))
		}
	})
}

func getTypeName(n ast.Node) string {
	switch n.(type) {
	case *ast.NumberNode:
		return "*ast.NumberNode"
	case *ast.StringNode:
		return "*ast.StringNode"
	case *ast.SymbolNode:
		return "*ast.SymbolNode"
	case *ast.ListNode:
		return "*ast.ListNode"
	case *ast.CallNode:
		return "*ast.CallNode"
	case *ast.IfNode:
		return "*ast.IfNode"
	case *ast.BeginNode:
		return "*ast.BeginNode"
	case *ast.DefineNode:
		return "*ast.DefineNode"
	case *ast.LambdaNode:
		return "*ast.LambdaNode"
	case *ast.LetNode:
		return "*ast.LetNode"
	case *ast.CondNode:
		return "*ast.CondNode"
	case *ast.QuoteNode:
		return "*ast.QuoteNode"
	default:
		return "unknown"
	}
}

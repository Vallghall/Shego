package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/eval"
	"github.com/Vallghall/schego/pkg/lex"
	"github.com/Vallghall/schego/pkg/mem"
)

// evalProgram tokenizes, parses, and evaluates a Scheme program.
func evalProgram(t *testing.T, program string) (mem.Object, *eval.State) {
	t.Helper()

	// Tokenize
	lexer := lex.New()
	tokens, err := lexer.Tokenize(program)
	require.NoError(t, err, "lexer error")

	// Parse
	parser := ast.New()
	nodes, err := parser.Parse(tokens)
	require.NoError(t, err, "parser error")

	// Create evaluator (builtins are loaded automatically)
	evaluator, err := eval.New()
	require.NoError(t, err, "evaluator creation error")

	// Evaluate
	result, err := evaluator.Eval(nodes)
	require.NoError(t, err, "eval error")

	return result, evaluator.State()
}

func expectNumber(t *testing.T, obj mem.Object, expected float64) {
	t.Helper()
	n, err := mem.AsNumber(obj)
	require.NoError(t, err, "expected number, got %T", obj)
	require.Equal(t, expected, n.Value())
}

func TestEvalArithmetic(t *testing.T) {
	t.Run("Addition", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(+ 1 2)", 3},
			{"(+ 1 2 3)", 6},
			{"(+ 1 2 3 4 5)", 15},
			{"(+)", 0},
			{"(+ 10)", 10},
			{"(+ -1 1)", 0},
			{"(+ 1.5 2.5)", 4},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("Subtraction", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(- 5 3)", 2},
			{"(- 10 3 2)", 5},
			{"(- 5)", -5},
			{"(- 0)", 0},
			{"(- 10 5 3 1)", 1},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("Multiplication", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(* 2 3)", 6},
			{"(* 2 3 4)", 24},
			{"(*)", 1},
			{"(* 5)", 5},
			{"(* 2 -3)", -6},
			{"(* 1.5 2)", 3},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("Division", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(/ 6 2)", 3},
			{"(/ 24 2 3)", 4},
			{"(/ 2)", 0.5},
			{"(/ 10 2 2)", 2.5},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("Modulo", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(modulo 10 3)", 1},
			{"(modulo 10 5)", 0},
			{"(modulo 7 2)", 1},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("Abs", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(abs 5)", 5},
			{"(abs -5)", 5},
			{"(abs 0)", 0},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("MinMax", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(min 1 2 3)", 1},
			{"(min 5)", 5},
			{"(min 3 1 4 1 5)", 1},
			{"(max 1 2 3)", 3},
			{"(max 5)", 5},
			{"(max 3 1 4 1 5)", 5},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})

	t.Run("NestedArithmetic", func(t *testing.T) {
		tests := []struct {
			program  string
			expected float64
		}{
			{"(+ 1 (* 2 3))", 7},
			{"(* (+ 1 2) (+ 3 4))", 21},
			{"(- (+ 10 5) (* 2 3))", 9},
			{"(/ (+ 10 20) (- 10 5))", 6},
			{"(+ (* 2 3) (* 4 5))", 26},
		}

		for _, tt := range tests {
			t.Run(tt.program, func(t *testing.T) {
				result, _ := evalProgram(t, tt.program)
				expectNumber(t, result, tt.expected)
			})
		}
	})
}

func TestEvalDefine(t *testing.T) {
	t.Run("SimpleDefine", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 42)
			x
		`)
		expectNumber(t, result, 42)
	})

	t.Run("DefineWithExpression", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x (+ 1 2))
			x
		`)
		expectNumber(t, result, 3)
	})

	t.Run("MultipleDefines", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 10)
			(define y 20)
			(+ x y)
		`)
		expectNumber(t, result, 30)
	})

	t.Run("FunctionDefine", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (square x) (* x x))
			(square 5)
		`)
		expectNumber(t, result, 25)
	})

	t.Run("FunctionWithMultipleParams", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (add a b) (+ a b))
			(add 3 4)
		`)
		expectNumber(t, result, 7)
	})

	t.Run("RecursiveFunction", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (factorial n)
				(if (= n 0)
					1
					(* n (factorial (- n 1)))))
			(factorial 5)
		`)
		expectNumber(t, result, 120)
	})
}

func TestEvalLambda(t *testing.T) {
	t.Run("SimpleLambda", func(t *testing.T) {
		result, _ := evalProgram(t, `
			((lambda (x) (* x x)) 4)
		`)
		expectNumber(t, result, 16)
	})

	t.Run("LambdaWithMultipleParams", func(t *testing.T) {
		result, _ := evalProgram(t, `
			((lambda (a b) (+ a b)) 3 7)
		`)
		expectNumber(t, result, 10)
	})

	t.Run("LambdaAsValue", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define square (lambda (x) (* x x)))
			(square 6)
		`)
		expectNumber(t, result, 36)
	})

	t.Run("Closure", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (make-adder n)
				(lambda (x) (+ x n)))
			(define add5 (make-adder 5))
			(add5 10)
		`)
		expectNumber(t, result, 15)
	})
}

func TestEvalIf(t *testing.T) {
	t.Run("IfTrue", func(t *testing.T) {
		result, _ := evalProgram(t, "(if #t 1 2)")
		expectNumber(t, result, 1)
	})

	t.Run("IfFalse", func(t *testing.T) {
		result, _ := evalProgram(t, "(if #f 1 2)")
		expectNumber(t, result, 2)
	})

	t.Run("IfWithoutElse", func(t *testing.T) {
		result, _ := evalProgram(t, "(if #t 42)")
		expectNumber(t, result, 42)
	})

	t.Run("IfWithExpression", func(t *testing.T) {
		result, _ := evalProgram(t, "(if (> 5 3) (* 2 3) (+ 1 1))")
		expectNumber(t, result, 6)
	})
}

func TestEvalLet(t *testing.T) {
	t.Run("SimpleLet", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(let ((x 10))
				x)
		`)
		expectNumber(t, result, 10)
	})

	t.Run("LetWithMultipleBindings", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(let ((x 10) (y 20))
				(+ x y))
		`)
		expectNumber(t, result, 30)
	})

	t.Run("NestedLet", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(let ((x 10))
				(let ((y 20))
					(+ x y)))
		`)
		expectNumber(t, result, 30)
	})

	t.Run("LetShadowing", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(let ((x 10))
				(let ((x 20))
					x))
		`)
		expectNumber(t, result, 20)
	})
}

func TestEvalBegin(t *testing.T) {
	t.Run("BeginReturnsLast", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(begin
				1
				2
				3)
		`)
		expectNumber(t, result, 3)
	})

	t.Run("BeginWithSideEffects", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 0)
			(define y 0)
			(begin
				(set! x 10)
				(set! y 20)
				(+ x y))
		`)
		expectNumber(t, result, 30)
	})
}

func TestSchemePrograms(t *testing.T) {
	t.Run("SumOfSquares", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (square x) (* x x))
			(define (sum-of-squares a b)
				(+ (square a) (square b)))
			(sum-of-squares 3 4)
		`)
		expectNumber(t, result, 25) // 9 + 16
	})

	t.Run("Fibonacci", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (fib n)
				(if (< n 2)
					n
					(+ (fib (- n 1)) (fib (- n 2)))))
			(fib 10)
		`)
		expectNumber(t, result, 55)
	})

	t.Run("GCD", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (gcd a b)
				(if (= b 0)
					a
					(gcd b (modulo a b))))
			(gcd 48 18)
		`)
		expectNumber(t, result, 6)
	})

	t.Run("Power", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (power base exp)
				(if (= exp 0)
					1
					(* base (power base (- exp 1)))))
			(power 2 10)
		`)
		expectNumber(t, result, 1024)
	})

	t.Run("Average", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define (average a b) (/ (+ a b) 2))
			(average 10 20)
		`)
		expectNumber(t, result, 15)
	})
}

func TestEvalSet(t *testing.T) {
	t.Run("SetSimple", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 10)
			(set! x 20)
			x
		`)
		expectNumber(t, result, 20)
	})

	t.Run("SetMultipleTimes", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 1)
			(set! x 2)
			(set! x 3)
			(set! x 4)
			x
		`)
		expectNumber(t, result, 4)
	})

	t.Run("SetWithExpression", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 10)
			(set! x (+ x 5))
			x
		`)
		expectNumber(t, result, 15)
	})

	t.Run("SetInNestedScope", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 10)
			(let ((y 5))
				(set! x (+ x y)))
			x
		`)
		expectNumber(t, result, 15)
	})

	t.Run("SetInLambda", func(t *testing.T) {
		// Test that set! works when called from within a function
		result, _ := evalProgram(t, `
			(define counter 0)
			(define (inc) (set! counter (+ counter 1)))
			(inc)
			(inc)
			(inc)
			counter
		`)
		expectNumber(t, result, 3)
	})

	t.Run("SetUndefinedError", func(t *testing.T) {
		lexer := lex.New()
		tokens, err := lexer.Tokenize("(set! undefined-var 42)")
		require.NoError(t, err)

		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)

		evaluator, err := eval.New()
		require.NoError(t, err)

		_, err = evaluator.Eval(nodes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unbound variable")
	})
}

func TestEvalRedefinitionError(t *testing.T) {
	t.Run("RedefineSameScope", func(t *testing.T) {
		lexer := lex.New()
		tokens, err := lexer.Tokenize(`
			(define x 10)
			(define x 20)
		`)
		require.NoError(t, err)

		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)

		evaluator, err := eval.New()
		require.NoError(t, err)

		_, err = evaluator.Eval(nodes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot redefine variable")
	})

	t.Run("DefineShadowInLet", func(t *testing.T) {
		// Shadowing in a new scope should be allowed
		result, _ := evalProgram(t, `
			(define x 10)
			(let ((x 20))
				x)
		`)
		expectNumber(t, result, 20)
	})

	t.Run("OriginalValuePreserved", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 10)
			(let ((x 20))
				x)
			x
		`)
		expectNumber(t, result, 10)
	})

	t.Run("RedefineFunctionError", func(t *testing.T) {
		lexer := lex.New()
		tokens, err := lexer.Tokenize(`
			(define (f x) x)
			(define (f x) (* x 2))
		`)
		require.NoError(t, err)

		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)

		evaluator, err := eval.New()
		require.NoError(t, err)

		_, err = evaluator.Eval(nodes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot redefine variable")
	})
}

func TestEvalQuote(t *testing.T) {
	t.Run("QuoteSymbol", func(t *testing.T) {
		result, state := evalProgram(t, "'foo")
		sym, err := mem.AsSymbol(result)
		require.NoError(t, err)
		name, ok := state.Pool().Resolve(sym.ID())
		require.True(t, ok)
		require.Equal(t, "foo", name)
	})

	t.Run("QuoteNumber", func(t *testing.T) {
		result, _ := evalProgram(t, "'42")
		expectNumber(t, result, 42)
	})

	t.Run("QuoteString", func(t *testing.T) {
		result, _ := evalProgram(t, `'"hello"`)
		str, err := mem.AsString(result)
		require.NoError(t, err)
		require.Equal(t, "hello", str.Value())
	})

	t.Run("QuoteList", func(t *testing.T) {
		result, _ := evalProgram(t, "'(1 2 3)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		expectNumber(t, items[0], 1)
		expectNumber(t, items[1], 2)
		expectNumber(t, items[2], 3)
	})

	t.Run("QuoteEmptyList", func(t *testing.T) {
		result, _ := evalProgram(t, "'()")
		require.Equal(t, mem.Nil, result)
	})

	t.Run("QuoteLongForm", func(t *testing.T) {
		result, state := evalProgram(t, "(quote foo)")
		sym, err := mem.AsSymbol(result)
		require.NoError(t, err)
		name, ok := state.Pool().Resolve(sym.ID())
		require.True(t, ok)
		require.Equal(t, "foo", name)
	})

	t.Run("QuoteNestedList", func(t *testing.T) {
		result, _ := evalProgram(t, "'((1 2) (3 4))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
		// Each item should be a list
		inner1, err := mem.ListToSlice(items[0])
		require.NoError(t, err)
		require.Len(t, inner1, 2)
		inner2, err := mem.ListToSlice(items[1])
		require.NoError(t, err)
		require.Len(t, inner2, 2)
	})

	t.Run("QuotePreservesSymbols", func(t *testing.T) {
		result, state := evalProgram(t, "'(a b c)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		// All should be symbols
		for i, item := range items {
			sym, err := mem.AsSymbol(item)
			require.NoError(t, err, "item %d should be a symbol", i)
			_ = sym
		}
		sym0, _ := mem.AsSymbol(items[0])
		name, ok := state.Pool().Resolve(sym0.ID())
		require.True(t, ok)
		require.Equal(t, "a", name)
	})
}

func TestEvalQuasiquote(t *testing.T) {
	t.Run("QuasiquoteSimple", func(t *testing.T) {
		result, _ := evalProgram(t, "`(1 2 3)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		expectNumber(t, items[0], 1)
		expectNumber(t, items[1], 2)
		expectNumber(t, items[2], 3)
	})

	t.Run("QuasiquoteWithUnquote", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 42)
			` + "`(a ,x b)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		// First and third are symbols, second is number
		expectNumber(t, items[1], 42)
	})

	t.Run("QuasiquoteWithUnquoteExpression", func(t *testing.T) {
		result, _ := evalProgram(t, "`(result is ,(+ 1 2))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		expectNumber(t, items[2], 3)
	})

	t.Run("QuasiquoteLongForm", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 5)
			(quasiquote (a (unquote x) c))
		`)
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		expectNumber(t, items[1], 5)
	})

	t.Run("QuasiquoteMultipleUnquotes", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 1)
			(define y 2)
			` + "`(,x ,y)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
		expectNumber(t, items[0], 1)
		expectNumber(t, items[1], 2)
	})
}

func TestEvalUnquoteSplicing(t *testing.T) {
	t.Run("UnquoteSplicingBasic", func(t *testing.T) {
		result, _ := evalProgram(t, "`(a ,@(list 1 2 3) b)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 5) // a, 1, 2, 3, b
		expectNumber(t, items[1], 1)
		expectNumber(t, items[2], 2)
		expectNumber(t, items[3], 3)
	})

	t.Run("UnquoteSplicingWithVariable", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define nums (list 1 2 3))
			` + "`(start ,@nums end)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 5) // start, 1, 2, 3, end
		expectNumber(t, items[1], 1)
		expectNumber(t, items[2], 2)
		expectNumber(t, items[3], 3)
	})

	t.Run("UnquoteSplicingEmptyList", func(t *testing.T) {
		result, _ := evalProgram(t, "`(a ,@(list) b)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2) // a, b - empty list disappears
	})

	t.Run("UnquoteSplicingMultiple", func(t *testing.T) {
		result, _ := evalProgram(t, "`(,@(list 1 2) ,@(list 3 4))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 4) // 1, 2, 3, 4
		expectNumber(t, items[0], 1)
		expectNumber(t, items[1], 2)
		expectNumber(t, items[2], 3)
		expectNumber(t, items[3], 4)
	})

	t.Run("UnquoteSplicingLongForm", func(t *testing.T) {
		result, _ := evalProgram(t, "(quasiquote (a (unquote-splicing (list 1 2)) b))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 4) // a, 1, 2, b
	})

	t.Run("UnquoteSplicingMixedWithUnquote", func(t *testing.T) {
		result, _ := evalProgram(t, `
			(define x 42)
			(define xs (list 1 2))
			` + "`(,x ,@xs)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3) // 42, 1, 2
		expectNumber(t, items[0], 42)
		expectNumber(t, items[1], 1)
		expectNumber(t, items[2], 2)
	})
}

func TestEvalNestedQuoting(t *testing.T) {
	t.Run("DoubleQuote", func(t *testing.T) {
		result, _ := evalProgram(t, "''foo")
		// Result should be (quote foo) as a list
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
		// First element should be symbol 'quote'
		sym, err := mem.AsSymbol(items[0])
		require.NoError(t, err)
		_ = sym
	})

	t.Run("QuoteQuasiquote", func(t *testing.T) {
		result, _ := evalProgram(t, "'`foo")
		// Result should be (quasiquote foo) as a list
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
	})

	t.Run("QuoteUnquote", func(t *testing.T) {
		result, _ := evalProgram(t, "',x")
		// Result should be (unquote x) as a list
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
	})

	t.Run("QuoteUnquoteSplicing", func(t *testing.T) {
		result, _ := evalProgram(t, "',@xs")
		// Result should be (unquote-splicing xs) as a list
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
	})

	t.Run("QuotedListWithQuotes", func(t *testing.T) {
		result, _ := evalProgram(t, "'('a 'b)")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
		// Each item should be a (quote x) list
		for _, item := range items {
			require.True(t, mem.IsList(item))
			inner, err := mem.ListToSlice(item)
			require.NoError(t, err)
			require.Len(t, inner, 2)
		}
	})
}

func TestEvalQuotingWithListOps(t *testing.T) {
	t.Run("CarOfQuotedList", func(t *testing.T) {
		result, state := evalProgram(t, "(car '(a b c))")
		sym, err := mem.AsSymbol(result)
		require.NoError(t, err)
		name, ok := state.Pool().Resolve(sym.ID())
		require.True(t, ok)
		require.Equal(t, "a", name)
	})

	t.Run("CdrOfQuotedList", func(t *testing.T) {
		result, _ := evalProgram(t, "(cdr '(a b c))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 2)
	})

	t.Run("LengthOfQuotedList", func(t *testing.T) {
		result, _ := evalProgram(t, "(length '(a b c d e))")
		expectNumber(t, result, 5)
	})

	t.Run("AppendQuotedLists", func(t *testing.T) {
		result, _ := evalProgram(t, "(append '(1 2) '(3 4))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 4)
	})

	t.Run("ReverseQuotedList", func(t *testing.T) {
		result, _ := evalProgram(t, "(reverse '(1 2 3))")
		require.True(t, mem.IsList(result))
		items, err := mem.ListToSlice(result)
		require.NoError(t, err)
		require.Len(t, items, 3)
		expectNumber(t, items[0], 3)
		expectNumber(t, items[1], 2)
		expectNumber(t, items[2], 1)
	})

	t.Run("NullOfQuotedEmptyList", func(t *testing.T) {
		result, _ := evalProgram(t, "(null? '())")
		require.Equal(t, mem.True, result)
	})

	t.Run("NullOfQuotedNonEmptyList", func(t *testing.T) {
		result, _ := evalProgram(t, "(null? '(a))")
		require.Equal(t, mem.False, result)
	})

	t.Run("PairOfQuotedList", func(t *testing.T) {
		result, _ := evalProgram(t, "(pair? '(a b))")
		require.Equal(t, mem.True, result)
	})

	t.Run("ListOfQuotedList", func(t *testing.T) {
		result, _ := evalProgram(t, "(list? '(a b c))")
		require.Equal(t, mem.True, result)
	})
}

func TestEvalQuotingErrors(t *testing.T) {
	t.Run("UnquoteOutsideQuasiquote", func(t *testing.T) {
		lexer := lex.New()
		tokens, err := lexer.Tokenize(",x")
		require.NoError(t, err)

		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)

		evaluator, err := eval.New()
		require.NoError(t, err)

		_, err = evaluator.Eval(nodes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unquote outside of quasiquote")
	})

	t.Run("UnquoteSplicingOutsideQuasiquote", func(t *testing.T) {
		lexer := lex.New()
		tokens, err := lexer.Tokenize(",@xs")
		require.NoError(t, err)

		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)

		evaluator, err := eval.New()
		require.NoError(t, err)

		_, err = evaluator.Eval(nodes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unquote-splicing outside of quasiquote")
	})

	t.Run("UnquoteSplicingNonList", func(t *testing.T) {
		lexer := lex.New()
		tokens, err := lexer.Tokenize("`(a ,@5 b)")
		require.NoError(t, err)

		parser := ast.New()
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)

		evaluator, err := eval.New()
		require.NoError(t, err)

		_, err = evaluator.Eval(nodes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unquote-splicing requires a list")
	})
}

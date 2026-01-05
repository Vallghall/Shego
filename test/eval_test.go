package test

import (
	"testing"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/builtin"
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
	if err != nil {
		t.Fatalf("lexer error: %v", err)
	}

	// Parse
	parser := ast.New()
	nodes, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	// Create evaluator with builtins
	evaluator := eval.New()
	if err := evaluator.State().LoadBuiltins(builtin.All); err != nil {
		t.Fatalf("builtin load error: %v", err)
	}

	// Evaluate
	result, err := evaluator.Eval(nodes)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	return result, evaluator.State()
}

func expectNumber(t *testing.T, obj mem.Object, expected float64) {
	t.Helper()
	n, err := mem.AsNumber(obj)
	if err != nil {
		t.Fatalf("expected number, got %T", obj)
	}
	if n.Value() != expected {
		t.Errorf("expected %v, got %v", expected, n.Value())
	}
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
			(begin
				(define x 10)
				(define y 20)
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


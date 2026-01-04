package test

import (
	"testing"

	"github.com/Vallghall/schego/pkg/lex"
)

func TestLexer(t *testing.T) {
	t.Run("ParenTokens", func(t *testing.T) {
		testParenTokens(t)
	})

	t.Run("StringTokens", func(t *testing.T) {
		testStringTokens(t)
	})

	t.Run("NumberTokens", func(t *testing.T) {
		testNumberTokens(t)
	})

	t.Run("AtomTokens", func(t *testing.T) {
		testAtomTokens(t)
	})

	t.Run("MixedExpressions", func(t *testing.T) {
		testMixedExpressions(t)
	})

	t.Run("ErrorCases", func(t *testing.T) {
		testErrorCases(t)
	})

	t.Run("PositionTracking", func(t *testing.T) {
		testPositionTracking(t)
	})
}

func testParenTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			kind lex.TKind
			raw  string
		}
	}{
		{
			name:  "single open paren",
			input: "(",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
			},
		},
		{
			name:  "single close paren",
			input: ")",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "empty parens",
			input: "()",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "nested parens",
			input: "((()))",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.ParenOpen, "("},
				{lex.ParenOpen, "("},
				{lex.ParenClose, ")"},
				{lex.ParenClose, ")"},
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "parens with spaces",
			input: "( ( ) )",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.ParenOpen, "("},
				{lex.ParenClose, ")"},
				{lex.ParenClose, ")"},
			},
		},
	}

	lexer := lex.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, exp := range tt.expected {
				if tokens[i].Kind() != exp.kind {
					t.Errorf("token %d: expected kind %v, got %v", i, exp.kind, tokens[i].Kind())
				}
				if tokens[i].Raw() != exp.raw {
					t.Errorf("token %d: expected raw %q, got %q", i, exp.raw, tokens[i].Raw())
				}
			}
		})
	}
}

func testStringTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			kind lex.TKind
			raw  string
		}
	}{
		{
			name:  "simple string",
			input: `"hello"`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.String, `"hello"`},
			},
		},
		{
			name:  "string with spaces",
			input: `"hello world"`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.String, `"hello world"`},
			},
		},
		{
			name:  "escaped quote",
			input: `"say \"hi\""`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.String, `"say "hi""`},
			},
		},
		{
			name:  "escaped backslash",
			input: `"path\\to\\file"`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.String, `"path\to\file"`},
			},
		},
		{
			name:  "empty string",
			input: `""`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.String, `""`},
			},
		},
		{
			name:  "multiple strings",
			input: `"a" "b"`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.String, `"a"`},
				{lex.String, `"b"`},
			},
		},
	}

	lexer := lex.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, exp := range tt.expected {
				if tokens[i].Kind() != exp.kind {
					t.Errorf("token %d: expected kind %v, got %v", i, exp.kind, tokens[i].Kind())
				}
				if tokens[i].Raw() != exp.raw {
					t.Errorf("token %d: expected raw %q, got %q", i, exp.raw, tokens[i].Raw())
				}
			}
		})
	}
}

func testNumberTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			kind lex.TKind
			raw  string
		}
	}{
		{
			name:  "integer",
			input: "42",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "42"},
			},
		},
		{
			name:  "negative integer",
			input: "-3",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "-3"},
			},
		},
		{
			name:  "positive integer",
			input: "+5",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "+5"},
			},
		},
		{
			name:  "decimal",
			input: "3.14",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "3.14"},
			},
		},
		{
			name:  "negative decimal",
			input: "-2.5",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "-2.5"},
			},
		},
		{
			name:  "zero",
			input: "0",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "0"},
			},
		},
		{
			name:  "multiple numbers",
			input: "1 2 3",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, "1"},
				{lex.Number, "2"},
				{lex.Number, "3"},
			},
		},
		{
			name:  "leading decimal",
			input: ".5",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Number, ".5"},
			},
		},
	}

	lexer := lex.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, exp := range tt.expected {
				if tokens[i].Kind() != exp.kind {
					t.Errorf("token %d: expected kind %v, got %v", i, exp.kind, tokens[i].Kind())
				}
				if tokens[i].Raw() != exp.raw {
					t.Errorf("token %d: expected raw %q, got %q", i, exp.raw, tokens[i].Raw())
				}
			}
		})
	}
}

func testAtomTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			kind lex.TKind
			raw  string
		}
	}{
		{
			name:  "simple identifier",
			input: "foo",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "foo"},
			},
		},
		{
			name:  "define keyword",
			input: "define",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "define"},
			},
		},
		{
			name:  "lambda keyword",
			input: "lambda",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "lambda"},
			},
		},
		{
			name:  "plus operator",
			input: "+",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "+"},
			},
		},
		{
			name:  "minus operator",
			input: "-",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "-"},
			},
		},
		{
			name:  "multiply operator",
			input: "*",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "*"},
			},
		},
		{
			name:  "divide operator",
			input: "/",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "/"},
			},
		},
		{
			name:  "comparison operators",
			input: "< > <= >= =",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "<"},
				{lex.Atom, ">"},
				{lex.Atom, "<="},
				{lex.Atom, ">="},
				{lex.Atom, "="},
			},
		},
		{
			name:  "hyphenated identifier",
			input: "my-function",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "my-function"},
			},
		},
		{
			name:  "question mark predicate",
			input: "null?",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "null?"},
			},
		},
		{
			name:  "bang mutator",
			input: "set!",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.Atom, "set!"},
			},
		},
	}

	lexer := lex.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, exp := range tt.expected {
				if tokens[i].Kind() != exp.kind {
					t.Errorf("token %d: expected kind %v, got %v", i, exp.kind, tokens[i].Kind())
				}
				if tokens[i].Raw() != exp.raw {
					t.Errorf("token %d: expected raw %q, got %q", i, exp.raw, tokens[i].Raw())
				}
			}
		})
	}
}

func testMixedExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			kind lex.TKind
			raw  string
		}
	}{
		{
			name:  "simple addition",
			input: "(+ 1 2)",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.Atom, "+"},
				{lex.Number, "1"},
				{lex.Number, "2"},
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "define with string",
			input: `(define x "hello")`,
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.Atom, "define"},
				{lex.Atom, "x"},
				{lex.String, `"hello"`},
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "nested expression",
			input: "(* (+ 1 2) (- 4 3))",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.Atom, "*"},
				{lex.ParenOpen, "("},
				{lex.Atom, "+"},
				{lex.Number, "1"},
				{lex.Number, "2"},
				{lex.ParenClose, ")"},
				{lex.ParenOpen, "("},
				{lex.Atom, "-"},
				{lex.Number, "4"},
				{lex.Number, "3"},
				{lex.ParenClose, ")"},
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "lambda expression",
			input: "(lambda (x) (* x x))",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.Atom, "lambda"},
				{lex.ParenOpen, "("},
				{lex.Atom, "x"},
				{lex.ParenClose, ")"},
				{lex.ParenOpen, "("},
				{lex.Atom, "*"},
				{lex.Atom, "x"},
				{lex.Atom, "x"},
				{lex.ParenClose, ")"},
				{lex.ParenClose, ")"},
			},
		},
		{
			name:  "signed numbers in expression",
			input: "(+ -3 +5)",
			expected: []struct {
				kind lex.TKind
				raw  string
			}{
				{lex.ParenOpen, "("},
				{lex.Atom, "+"},
				{lex.Number, "-3"},
				{lex.Number, "+5"},
				{lex.ParenClose, ")"},
			},
		},
	}

	lexer := lex.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, exp := range tt.expected {
				if tokens[i].Kind() != exp.kind {
					t.Errorf("token %d: expected kind %v, got %v", i, exp.kind, tokens[i].Kind())
				}
				if tokens[i].Raw() != exp.raw {
					t.Errorf("token %d: expected raw %q, got %q", i, exp.raw, tokens[i].Raw())
				}
			}
		})
	}
}

func testErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unclosed string",
			input: `"hello`,
		},
		{
			name:  "unclosed string with escape",
			input: `"hello\"`,
		},
	}

	lexer := lex.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := lexer.Tokenize(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q, got none", tt.input)
			}
		})
	}
}

func testPositionTracking(t *testing.T) {
	lexer := lex.New(lex.WithFile("test.scm"))

	input := "(+ 1\n   2)"
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []struct {
		raw     string
		line    int
		linePos int
	}{
		{"(", 1, 1},
		{"+", 1, 2},
		{"1", 1, 4},
		{"2", 2, 4},
		{")", 2, 5},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		tok := tokens[i]
		if tok.Raw() != exp.raw {
			t.Errorf("token %d: expected raw %q, got %q", i, exp.raw, tok.Raw())
		}
		if tok.Line() != exp.line {
			t.Errorf("token %d (%q): expected line %d, got %d", i, exp.raw, exp.line, tok.Line())
		}
		if tok.LinePosition() != exp.linePos {
			t.Errorf("token %d (%q): expected linePos %d, got %d", i, exp.raw, exp.linePos, tok.LinePosition())
		}
		if tok.File() != "test.scm" {
			t.Errorf("token %d: expected file %q, got %q", i, "test.scm", tok.File())
		}
	}
}

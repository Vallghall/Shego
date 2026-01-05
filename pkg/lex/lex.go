package lex

import (
	"fmt"
)

// Lexer tokenizes Scheme source code.
type Lexer interface {
	// Tokenize parses the input string and returns a slice of tokens.
	Tokenize(input string) ([]Token, error)
}

// LexError represents an error that occurred during lexing.
type LexError struct {
	Message  string
	Position int
	Line     int
	LinePos  int
	File     string
}

func (e *LexError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("%s:%d:%d: %s", e.File, e.Line, e.LinePos, e.Message)
	}
	return fmt.Sprintf("%d:%d: %s", e.Line, e.LinePos, e.Message)
}

// lexer is the concrete implementation of Lexer using Chain of Responsibility.
type lexer struct {
	chain Node
	file  string
}

// Option is a functional option for configuring the lexer.
type Option func(*lexer)

// WithFile sets the filename for error messages and token metadata.
func WithFile(file string) Option {
	return func(l *lexer) {
		l.file = file
	}
}

// New creates a new Lexer with the default chain of responsibility.
// Chain order: Paren -> String -> Number -> Atom
func New(opts ...Option) Lexer {
	l := &lexer{
		chain: BuildChain(
			NewParenNode(),
			NewStringNode(),
			NewNumberNode(),
			NewAtomNode(),
		),
		file: "",
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Tokenize parses the input string and returns a slice of tokens.
func (l *lexer) Tokenize(input string) ([]Token, error) {
	var tokens []Token
	r := NewReader(input, l.file)

	for !r.EOF() {
		// Skip whitespace
		r.SkipWhitespace()
		if r.EOF() {
			break
		}

		// Save position for error reporting
		file, pos, line, linePos := r.Snapshot()
		ch := r.Current()

		// Try to parse a token using the chain
		token, handled := l.chain.Handle(r)

		if !handled || token == nil {
			// Check for specific error conditions
			if ch == '"' {
				return nil, &LexError{
					Message:  "unclosed string literal",
					Position: pos,
					Line:     line,
					LinePos:  linePos,
					File:     file,
				}
			}

			return nil, &LexError{
				Message:  fmt.Sprintf("unexpected character: %q", ch),
				Position: pos,
				Line:     line,
				LinePos:  linePos,
				File:     file,
			}
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

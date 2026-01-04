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

	pos := 0
	line := 1
	linePos := 1

	for pos < len(input) {
		ch := input[pos]

		// Skip whitespace while tracking position
		if ch == ' ' || ch == '\t' {
			pos++
			linePos++
			continue
		}

		if ch == '\n' {
			pos++
			line++
			linePos = 1
			continue
		}

		if ch == '\r' {
			pos++
			// Handle \r\n as single newline
			if pos < len(input) && input[pos] == '\n' {
				pos++
			}
			line++
			linePos = 1
			continue
		}

		// Try to parse a token using the chain
		token, consumed, handled := l.chain.Handle(input, pos, line, linePos, l.file)

		if !handled || token == nil {
			// Check for specific error conditions
			if ch == '"' {
				return nil, &LexError{
					Message:  "unclosed string literal",
					Position: pos,
					Line:     line,
					LinePos:  linePos,
					File:     l.file,
				}
			}

			return nil, &LexError{
				Message:  fmt.Sprintf("unexpected character: %q", ch),
				Position: pos,
				Line:     line,
				LinePos:  linePos,
				File:     l.file,
			}
		}

		tokens = append(tokens, token)

		// Update position tracking
		for i := 0; i < consumed; i++ {
			if input[pos+i] == '\n' {
				line++
				linePos = 1
			} else {
				linePos++
			}
		}
		pos += consumed
	}

	return tokens, nil
}

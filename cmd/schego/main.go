package main

import (
	"fmt"
	"os"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/eval"
	"github.com/Vallghall/schego/pkg/lex"
	"github.com/Vallghall/schego/pkg/mem"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: schego <file.scm>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read source file
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Tokenize
	lexer := lex.New()
	tokens, err := lexer.Tokenize(string(source))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lexer error: %v\n", err)
		os.Exit(1)
	}

	// Parse
	parser := ast.New()
	nodes, err := parser.Parse(tokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parser error: %v\n", err)
		os.Exit(1)
	}

	// Create evaluator (builtins are loaded automatically)
	evaluator, err := eval.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating evaluator: %v\n", err)
		os.Exit(1)
	}

	// Evaluate
	result, err := evaluator.Eval(nodes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}

	// Print result if not void
	if result.Type() != mem.TypeVoid {
		fmt.Println(result)
	}
}

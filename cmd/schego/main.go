package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/eval"
	"github.com/Vallghall/schego/pkg/lex"
	"github.com/Vallghall/schego/pkg/mem"
)

var (
	// Flag variables
	evalExpr string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "schego [file.scm]",
	Short: "Schego - A Scheme interpreter written in Go",
	Long: `Schego is a Scheme interpreter written in Go that supports
core Scheme features including lambdas, closures, and tail recursion.

Examples:
  schego program.scm          Run a Scheme file
  schego -e "(+ 1 2 3)"       Evaluate an expression and exit`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSchego,
}

func init() {
	// Bind flags
	rootCmd.Flags().StringVarP(&evalExpr, "eval", "e", "", "evaluate expression and exit")

	// Bind to viper for env support
	viper.BindPFlag("eval", rootCmd.Flags().Lookup("eval"))

	// Environment variable support (SCHEGO_EVAL)
	viper.SetEnvPrefix("SCHEGO")
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runSchego(cmd *cobra.Command, args []string) error {
	// Get eval expression from flag or env
	expr := viper.GetString("eval")

	// Determine source code
	var source string

	if expr != "" {
		// -e flag: evaluate expression
		source = expr
	} else if len(args) == 1 {
		// File argument: read and evaluate file
		content, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		source = string(content)
	} else {
		// No input provided
		return fmt.Errorf("no input provided. Use -e <expr> or provide a file")
	}

	// Execute the source
	result, err := execute(source)
	if err != nil {
		return err
	}

	// Print result if not void
	if result.Type() != mem.TypeVoid {
		fmt.Println(result)
	}

	return nil
}

// execute tokenizes, parses, and evaluates Scheme source code.
func execute(source string) (mem.Object, error) {
	// Tokenize
	lexer := lex.New()
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		return nil, fmt.Errorf("lexer error: %w", err)
	}

	// Parse
	parser := ast.New()
	nodes, err := parser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	// Create evaluator (builtins are loaded automatically)
	evaluator, err := eval.New()
	if err != nil {
		return nil, fmt.Errorf("error creating evaluator: %w", err)
	}

	// Evaluate
	result, err := evaluator.Eval(nodes)
	if err != nil {
		return nil, fmt.Errorf("runtime error: %w", err)
	}

	return result, nil
}

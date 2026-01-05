package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Vallghall/schego/pkg/ast"
	"github.com/Vallghall/schego/pkg/builtin"
	"github.com/Vallghall/schego/pkg/eval"
	"github.com/Vallghall/schego/pkg/lex"
	"github.com/Vallghall/schego/pkg/mem"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create evaluator with builtins
	evaluator := eval.New()
	if err := evaluator.State().LoadBuiltins(builtin.All); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading builtins: %v\n", err)
		os.Exit(1)
	}

	// Create lexer and parser
	lexer := lex.New()
	parser := ast.New()

	// Command history
	var history []string
	historyIndex := -1

	// Create a text view for output
	outputView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	outputView.SetBorder(true).SetTitle(" Schego REPL ")

	fmt.Fprintln(outputView, "[yellow]Welcome to Schego Scheme REPL[white]")
	fmt.Fprintln(outputView, "[dim]Type (exit) or (quit) to exit[white]")
	fmt.Fprintln(outputView, "")

	// Create an input field
	inputField := tview.NewInputField().
		SetLabel("[green]λ>[white] ").
		SetLabelColor(tcell.ColorGreen).
		SetFieldWidth(0)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			input := strings.TrimSpace(inputField.GetText())
			if input == "" {
				return
			}

			// Add to history
			history = append(history, input)
			historyIndex = len(history)

			// Echo input
			fmt.Fprintf(outputView, "[green]λ>[white] %s\n", input)

			// Check for exit commands
			if input == "(exit)" || input == "(quit)" {
				app.Stop()
				return
			}

			// Evaluate
			result, err := evalInput(lexer, parser, evaluator, input)
			if err != nil {
				fmt.Fprintf(outputView, "[red]Error:[white] %v\n", err)
			} else if result.Type() != mem.TypeVoid {
				fmt.Fprintf(outputView, "[cyan]=> %s[white]\n", result.String())
			}

			// Clear input
			inputField.SetText("")
		}
	})

	// Handle up/down arrows for history
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if historyIndex > 0 {
				historyIndex--
				inputField.SetText(history[historyIndex])
			}
			return nil
		case tcell.KeyDown:
			if historyIndex < len(history)-1 {
				historyIndex++
				inputField.SetText(history[historyIndex])
			} else {
				historyIndex = len(history)
				inputField.SetText("")
			}
			return nil
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		}
		return event
	})

	inputField.SetBorder(true)

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(outputView, 0, 1, false).
		AddItem(inputField, 3, 0, true)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// evalInput tokenizes, parses, and evaluates a single input line.
func evalInput(lexer lex.Lexer, parser ast.Parser, evaluator eval.Evaluator, input string) (mem.Object, error) {
	// Tokenize
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		return nil, fmt.Errorf("syntax error: %v", err)
	}

	// Parse
	nodes, err := parser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	// Evaluate
	result, err := evaluator.Eval(nodes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

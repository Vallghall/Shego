package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create a text view for output
	outputView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	outputView.SetBorder(true).SetTitle(" Schego REPL ")

	fmt.Fprintln(outputView, "Welcome to Schego Scheme REPL")
	fmt.Fprintln(outputView, "Type (exit) to quit")
	fmt.Fprintln(outputView, "")

	// Create an input field
	inputField := tview.NewInputField().
		SetLabel("λ> ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// TODO: Implement expression evaluation
				fmt.Fprintf(outputView, "=> (not yet implemented)\n")
			}
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

package builtin

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/mem"
)

// IODefinitions returns all I/O operation definitions.
func IODefinitions() []mem.Definition {
	return []mem.Definition{
		{Name: "display", Value: mem.NewPrimitive("display", 1, display)},
		{Name: "write", Value: mem.NewPrimitive("write", 1, write)},
		{Name: "newline", Value: mem.NewPrimitive("newline", 0, newline)},
	}
}

// display implements (display obj) - prints object to stdout.
// For strings, prints without quotes. For other types, prints string representation.
func display(args []mem.Object) (mem.Object, error) {
	obj := args[0]

	// For strings, print the raw value without quotes
	if s, ok := obj.(*mem.String); ok {
		fmt.Print(s.Value())
	} else {
		fmt.Print(obj.String())
	}

	return mem.Void, nil
}

// write implements (write obj) - prints object to stdout with quotes for strings.
// Preserves the Scheme readable representation.
func write(args []mem.Object) (mem.Object, error) {
	obj := args[0]
	fmt.Print(obj.String())
	return mem.Void, nil
}

// newline implements (newline) - prints a newline to stdout.
func newline(args []mem.Object) (mem.Object, error) {
	fmt.Println()
	return mem.Void, nil
}


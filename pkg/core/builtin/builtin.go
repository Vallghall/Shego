package builtin

import (
	"github.com/Vallghall/schego/pkg/mem"
)

// Definitions returns all builtin operation definitions.
// This aggregates arithmetic, comparison, conversion, and I/O operations.
func Definitions() []mem.Definition {
	var defs []mem.Definition

	// Arithmetic operations
	defs = append(defs, ArythDefinitions()...)

	// Comparison operations
	defs = append(defs, CmpDefinitions()...)

	// Conversion operations
	defs = append(defs, ConvDefinitions()...)

	// I/O operations
	defs = append(defs, IODefinitions()...)

	return defs
}

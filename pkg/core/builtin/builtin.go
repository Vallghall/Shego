package builtin

import (
	"github.com/Vallghall/schego/pkg/mem"
)

// Definitions returns all builtin operation definitions.
// This aggregates arithmetic, comparison, and conversion operations.
func Definitions() []mem.Definition {
	var defs []mem.Definition

	// Arithmetic operations
	defs = append(defs, ArythDefinitions()...)

	// Comparison operations
	defs = append(defs, CmpDefinitions()...)

	// Conversion operations
	defs = append(defs, ConvDefinitions()...)

	return defs
}

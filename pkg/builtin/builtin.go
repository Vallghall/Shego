package builtin

import (
	"github.com/Vallghall/schego/pkg/builtin/aryth"
	"github.com/Vallghall/schego/pkg/builtin/cmp"
	"github.com/Vallghall/schego/pkg/mem"
)

// All returns all built-in definitions.
// This is the main entry point for loading builtins into the evaluator.
func All() []mem.Definition {
	var defs []mem.Definition

	// Arithmetic operations
	defs = append(defs, aryth.Definitions()...)

	// Comparison operations
	defs = append(defs, cmp.Definitions()...)

	// Add more builtin categories here:
	// defs = append(defs, list.Definitions()...)
	// etc.

	return defs
}

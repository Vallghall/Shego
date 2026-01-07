package core

import (
	"github.com/Vallghall/schego/pkg/core/builtin"
	schmath "github.com/Vallghall/schego/pkg/core/math"
	"github.com/Vallghall/schego/pkg/mem"
)

// All returns all built-in definitions from all core subpackages.
// This is the main entry point for loading builtins into the evaluator.
func All() []mem.Definition {
	var defs []mem.Definition

	// Builtin operations (arithmetic, comparison, conversion)
	defs = append(defs, builtin.Definitions()...)

	// Math operations
	defs = append(defs, schmath.Definitions()...)

	// Add more core categories here:
	// defs = append(defs, list.Definitions()...)
	// etc.

	return defs
}

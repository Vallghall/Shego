package builtin

import (
	"github.com/Vallghall/schego/pkg/mem"
)

// EvalDefinitions returns eval primitive definition.
// The eval function is created with evaluator access via closure.
func EvalDefinitions(evalFn func(mem.Object) (mem.Object, error)) []mem.Definition {
	return []mem.Definition{
		{Name: "eval", Value: mem.NewPrimitive("eval", 1, func(args []mem.Object) (mem.Object, error) {
			return evalFn(args[0])
		})},
	}
}

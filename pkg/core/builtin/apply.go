package builtin

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/mem"
)

// ApplyDefinitions returns apply primitive definition.
// The apply function is created with evaluator access via closure.
func ApplyDefinitions(applyFn func(*mem.Procedure, []mem.Object) (mem.Object, error)) []mem.Definition {
	return []mem.Definition{
		{Name: "apply", Value: mem.NewVariadicPrimitive("apply", func(args []mem.Object) (mem.Object, error) {
			if len(args) < 2 {
				return nil, mem.NewArityAtLeastError("apply", 2, len(args))
			}

			// First arg must be a procedure
			proc, err := mem.AsProcedure(args[0])
			if err != nil {
				return nil, mem.NewNotCallableError(args[0])
			}

			// Last arg must be a list
			lastArg := args[len(args)-1]
			if !mem.IsList(lastArg) {
				return nil, mem.NewTypeError(mem.TypePair, lastArg.Type(), "last argument must be a list")
			}

			// Convert list to slice
			listArgs, err := mem.ListToSlice(lastArg)
			if err != nil {
				return nil, fmt.Errorf("apply: %w", err)
			}

			// Prepend all variadic args (except last) to the list
			variadicArgs := args[1 : len(args)-1] // all args except procedure and list
			combinedArgs := make([]mem.Object, len(variadicArgs)+len(listArgs))
			copy(combinedArgs, variadicArgs)                    // variadic args first
			copy(combinedArgs[len(variadicArgs):], listArgs)     // list args appended

			// Call the procedure
			return applyFn(proc, combinedArgs)
		})},
	}
}

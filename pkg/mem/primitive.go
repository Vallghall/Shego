package mem

// PrimitiveFunc is the signature for built-in procedures.
// Primitives receive evaluated arguments and return a result.
// They don't have direct access to evaluation state - if state access
// is needed, it should be provided through closures or specific object types.
type PrimitiveFunc func(args []Object) (Object, error)

// Definition represents a name-value pair for bulk loading of bindings.
// Used to define built-in procedures and constants.
type Definition struct {
	Name  string
	Value Object
}


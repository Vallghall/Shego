package builtin

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/mem"
)

// ListDefinitions returns all list operation definitions.
func ListDefinitions() []mem.Definition {
	return []mem.Definition{
		// Basic pair operations
		{Name: "cons", Value: mem.NewPrimitive("cons", 2, cons)},
		{Name: "car", Value: mem.NewPrimitive("car", 1, car)},
		{Name: "cdr", Value: mem.NewPrimitive("cdr", 1, cdr)},

		// Two-level compositions
		{Name: "caar", Value: mem.NewPrimitive("caar", 1, caar)},
		{Name: "cadr", Value: mem.NewPrimitive("cadr", 1, cadr)},
		{Name: "cdar", Value: mem.NewPrimitive("cdar", 1, cdar)},
		{Name: "cddr", Value: mem.NewPrimitive("cddr", 1, cddr)},

		// Three-level compositions
		{Name: "caaar", Value: mem.NewPrimitive("caaar", 1, caaar)},
		{Name: "caadr", Value: mem.NewPrimitive("caadr", 1, caadr)},
		{Name: "cadar", Value: mem.NewPrimitive("cadar", 1, cadar)},
		{Name: "caddr", Value: mem.NewPrimitive("caddr", 1, caddr)},
		{Name: "cdaar", Value: mem.NewPrimitive("cdaar", 1, cdaar)},
		{Name: "cdadr", Value: mem.NewPrimitive("cdadr", 1, cdadr)},
		{Name: "cddar", Value: mem.NewPrimitive("cddar", 1, cddar)},
		{Name: "cdddr", Value: mem.NewPrimitive("cdddr", 1, cdddr)},

		// List utilities
		{Name: "list", Value: mem.NewVariadicPrimitive("list", list)},
		{Name: "length", Value: mem.NewPrimitive("length", 1, length)},
		{Name: "append", Value: mem.NewVariadicPrimitive("append", appendLists)},
		{Name: "reverse", Value: mem.NewPrimitive("reverse", 1, reverse)},

		// Predicates
		{Name: "null?", Value: mem.NewPrimitive("null?", 1, isNull)},
		{Name: "pair?", Value: mem.NewPrimitive("pair?", 1, isPair)},
		{Name: "list?", Value: mem.NewPrimitive("list?", 1, isList)},
	}
}

// =============================================================================
// Basic Pair Operations
// =============================================================================

// cons implements (cons a b) - creates a new pair.
func cons(args []mem.Object) (mem.Object, error) {
	return mem.NewPair(args[0], args[1]), nil
}

// car implements (car pair) - returns the first element of a pair.
func car(args []mem.Object) (mem.Object, error) {
	pair, err := mem.AsPair(args[0])
	if err != nil {
		return nil, fmt.Errorf("car: expected pair, got %s", args[0].Type())
	}
	return pair.Car(), nil
}

// cdr implements (cdr pair) - returns the second element of a pair.
func cdr(args []mem.Object) (mem.Object, error) {
	pair, err := mem.AsPair(args[0])
	if err != nil {
		return nil, fmt.Errorf("cdr: expected pair, got %s", args[0].Type())
	}
	return pair.Cdr(), nil
}

// =============================================================================
// Two-level Compositions
// =============================================================================

// caar implements (caar x) = (car (car x))
func caar(args []mem.Object) (mem.Object, error) {
	result, err := car(args)
	if err != nil {
		return nil, fmt.Errorf("caar: %w", err)
	}
	return car([]mem.Object{result})
}

// cadr implements (cadr x) = (car (cdr x)) - second element
func cadr(args []mem.Object) (mem.Object, error) {
	result, err := cdr(args)
	if err != nil {
		return nil, fmt.Errorf("cadr: %w", err)
	}
	return car([]mem.Object{result})
}

// cdar implements (cdar x) = (cdr (car x))
func cdar(args []mem.Object) (mem.Object, error) {
	result, err := car(args)
	if err != nil {
		return nil, fmt.Errorf("cdar: %w", err)
	}
	return cdr([]mem.Object{result})
}

// cddr implements (cddr x) = (cdr (cdr x)) - rest after second
func cddr(args []mem.Object) (mem.Object, error) {
	result, err := cdr(args)
	if err != nil {
		return nil, fmt.Errorf("cddr: %w", err)
	}
	return cdr([]mem.Object{result})
}

// =============================================================================
// Three-level Compositions
// =============================================================================

// caaar implements (caaar x) = (car (car (car x)))
func caaar(args []mem.Object) (mem.Object, error) {
	result, err := caar(args)
	if err != nil {
		return nil, fmt.Errorf("caaar: %w", err)
	}
	return car([]mem.Object{result})
}

// caadr implements (caadr x) = (car (car (cdr x)))
func caadr(args []mem.Object) (mem.Object, error) {
	result, err := cadr(args)
	if err != nil {
		return nil, fmt.Errorf("caadr: %w", err)
	}
	return car([]mem.Object{result})
}

// cadar implements (cadar x) = (car (cdr (car x)))
func cadar(args []mem.Object) (mem.Object, error) {
	result, err := cdar(args)
	if err != nil {
		return nil, fmt.Errorf("cadar: %w", err)
	}
	return car([]mem.Object{result})
}

// caddr implements (caddr x) = (car (cdr (cdr x))) - third element
func caddr(args []mem.Object) (mem.Object, error) {
	result, err := cddr(args)
	if err != nil {
		return nil, fmt.Errorf("caddr: %w", err)
	}
	return car([]mem.Object{result})
}

// cdaar implements (cdaar x) = (cdr (car (car x)))
func cdaar(args []mem.Object) (mem.Object, error) {
	result, err := caar(args)
	if err != nil {
		return nil, fmt.Errorf("cdaar: %w", err)
	}
	return cdr([]mem.Object{result})
}

// cdadr implements (cdadr x) = (cdr (car (cdr x)))
func cdadr(args []mem.Object) (mem.Object, error) {
	result, err := cadr(args)
	if err != nil {
		return nil, fmt.Errorf("cdadr: %w", err)
	}
	return cdr([]mem.Object{result})
}

// cddar implements (cddar x) = (cdr (cdr (car x)))
func cddar(args []mem.Object) (mem.Object, error) {
	result, err := cdar(args)
	if err != nil {
		return nil, fmt.Errorf("cddar: %w", err)
	}
	return cdr([]mem.Object{result})
}

// cdddr implements (cdddr x) = (cdr (cdr (cdr x))) - rest after third
func cdddr(args []mem.Object) (mem.Object, error) {
	result, err := cddr(args)
	if err != nil {
		return nil, fmt.Errorf("cdddr: %w", err)
	}
	return cdr([]mem.Object{result})
}

// =============================================================================
// List Utilities
// =============================================================================

// list implements (list a b c ...) - creates a proper list from arguments.
func list(args []mem.Object) (mem.Object, error) {
	return mem.SliceToList(args), nil
}

// length implements (length lst) - returns the length of a proper list.
func length(args []mem.Object) (mem.Object, error) {
	count := 0
	current := args[0]

	for {
		switch current.Type() {
		case mem.TypeNil:
			return mem.NewNumber(float64(count)), nil
		case mem.TypePair:
			pair := current.(*mem.Pair)
			count++
			current = pair.Cdr()
		default:
			return nil, fmt.Errorf("length: not a proper list")
		}
	}
}

// appendLists implements (append lst1 lst2 ...) - concatenates lists.
func appendLists(args []mem.Object) (mem.Object, error) {
	if len(args) == 0 {
		return mem.Nil, nil
	}

	// Last argument is returned as-is (can be improper)
	if len(args) == 1 {
		return args[0], nil
	}

	// Collect all elements from all but last list
	var elements []mem.Object
	for i := 0; i < len(args)-1; i++ {
		lst := args[i]
		for {
			switch lst.Type() {
			case mem.TypeNil:
				goto nextList
			case mem.TypePair:
				pair := lst.(*mem.Pair)
				elements = append(elements, pair.Car())
				lst = pair.Cdr()
			default:
				return nil, fmt.Errorf("append: expected list, got %s", lst.Type())
			}
		}
	nextList:
	}

	// Build result by prepending elements to last argument
	result := args[len(args)-1]
	for i := len(elements) - 1; i >= 0; i-- {
		result = mem.NewPair(elements[i], result)
	}

	return result, nil
}

// reverse implements (reverse lst) - returns a reversed copy of the list.
func reverse(args []mem.Object) (mem.Object, error) {
	elements, err := mem.ListToSlice(args[0])
	if err != nil {
		return nil, fmt.Errorf("reverse: %w", err)
	}

	// Reverse the slice
	for i, j := 0, len(elements)-1; i < j; i, j = i+1, j-1 {
		elements[i], elements[j] = elements[j], elements[i]
	}

	return mem.SliceToList(elements), nil
}

// =============================================================================
// Predicates
// =============================================================================

// isNull implements (null? obj) - returns #t if obj is the empty list.
func isNull(args []mem.Object) (mem.Object, error) {
	return mem.Boolean(args[0].Type() == mem.TypeNil), nil
}

// isPair implements (pair? obj) - returns #t if obj is a pair.
func isPair(args []mem.Object) (mem.Object, error) {
	return mem.Boolean(args[0].Type() == mem.TypePair), nil
}

// isList implements (list? obj) - returns #t if obj is a proper list.
func isList(args []mem.Object) (mem.Object, error) {
	return mem.Boolean(mem.IsList(args[0])), nil
}


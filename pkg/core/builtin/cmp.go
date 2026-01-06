package builtin

import (
	"github.com/Vallghall/schego/pkg/mem"
)

// CmpDefinitions returns all comparison operation definitions.
func CmpDefinitions() []mem.Definition {
	return []mem.Definition{
		{Name: "=", Value: mem.NewVariadicPrimitive("=", numEqual)},
		{Name: "<", Value: mem.NewVariadicPrimitive("<", lessThan)},
		{Name: ">", Value: mem.NewVariadicPrimitive(">", greaterThan)},
		{Name: "<=", Value: mem.NewVariadicPrimitive("<=", lessOrEqual)},
		{Name: ">=", Value: mem.NewVariadicPrimitive(">=", greaterOrEqual)},
		{Name: "zero?", Value: mem.NewPrimitive("zero?", 1, isZero)},
		{Name: "positive?", Value: mem.NewPrimitive("positive?", 1, isPositive)},
		{Name: "negative?", Value: mem.NewPrimitive("negative?", 1, isNegative)},
		{Name: "eq?", Value: mem.NewPrimitive("eq?", 2, isEq)},
		{Name: "equal?", Value: mem.NewPrimitive("equal?", 2, isEqual)},
		{Name: "not", Value: mem.NewPrimitive("not", 1, not)},
	}
}

// numEqual implements (= n1 n2 ...) - returns #t if all numbers are equal.
func numEqual(args []mem.Object) (mem.Object, error) {
	if len(args) < 2 {
		return nil, mem.NewArityAtLeastError("=", 2, len(args))
	}

	first, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}

	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(), "")
		}
		if first.Value() != n.Value() {
			return mem.False, nil
		}
	}
	return mem.True, nil
}

// lessThan implements (< n1 n2 ...) - returns #t if strictly increasing.
func lessThan(args []mem.Object) (mem.Object, error) {
	if len(args) < 2 {
		return nil, mem.NewArityAtLeastError("<", 2, len(args))
	}

	prev, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}

	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(), "")
		}
		if prev.Value() >= n.Value() {
			return mem.False, nil
		}
		prev = n
	}
	return mem.True, nil
}

// greaterThan implements (> n1 n2 ...) - returns #t if strictly decreasing.
func greaterThan(args []mem.Object) (mem.Object, error) {
	if len(args) < 2 {
		return nil, mem.NewArityAtLeastError(">", 2, len(args))
	}

	prev, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}

	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(), "")
		}
		if prev.Value() <= n.Value() {
			return mem.False, nil
		}
		prev = n
	}
	return mem.True, nil
}

// lessOrEqual implements (<= n1 n2 ...) - returns #t if non-decreasing.
func lessOrEqual(args []mem.Object) (mem.Object, error) {
	if len(args) < 2 {
		return nil, mem.NewArityAtLeastError("<=", 2, len(args))
	}

	prev, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}

	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(), "")
		}
		if prev.Value() > n.Value() {
			return mem.False, nil
		}
		prev = n
	}
	return mem.True, nil
}

// greaterOrEqual implements (>= n1 n2 ...) - returns #t if non-increasing.
func greaterOrEqual(args []mem.Object) (mem.Object, error) {
	if len(args) < 2 {
		return nil, mem.NewArityAtLeastError(">=", 2, len(args))
	}

	prev, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}

	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(), "")
		}
		if prev.Value() < n.Value() {
			return mem.False, nil
		}
		prev = n
	}
	return mem.True, nil
}

// isZero implements (zero? n) - returns #t if n is zero.
func isZero(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}
	return mem.Boolean(n.Value() == 0), nil
}

// isPositive implements (positive? n) - returns #t if n > 0.
func isPositive(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}
	return mem.Boolean(n.Value() > 0), nil
}

// isNegative implements (negative? n) - returns #t if n < 0.
func isNegative(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(), "")
	}
	return mem.Boolean(n.Value() < 0), nil
}

// isEq implements (eq? a b) - identity comparison.
func isEq(args []mem.Object) (mem.Object, error) {
	a, b := args[0], args[1]

	// For symbols, compare IDs
	if symA, ok := a.(*mem.Symbol); ok {
		if symB, ok := b.(*mem.Symbol); ok {
			return mem.Boolean(symA.ID() == symB.ID()), nil
		}
		return mem.False, nil
	}

	// For other types, use pointer equality or value equality for primitives
	switch a.Type() {
	case mem.TypeNil, mem.TypeVoid, mem.TypeBoolean:
		return mem.Boolean(a == b), nil
	case mem.TypeNumber:
		if na, ok := a.(*mem.Number); ok {
			if nb, ok := b.(*mem.Number); ok {
				return mem.Boolean(na.Value() == nb.Value()), nil
			}
		}
	}

	return mem.Boolean(a == b), nil
}

// isEqual implements (equal? a b) - structural equality.
func isEqual(args []mem.Object) (mem.Object, error) {
	return mem.Boolean(args[0].Equal(args[1])), nil
}

// not implements (not x) - boolean negation.
func not(args []mem.Object) (mem.Object, error) {
	return mem.Boolean(!args[0].IsTruthy()), nil
}

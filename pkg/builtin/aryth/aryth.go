package aryth

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/mem"
)

// Definitions returns all arithmetic operation definitions.
func Definitions() []mem.Definition {
	return []mem.Definition{
		{Name: "+", Value: mem.NewVariadicPrimitive("+", add)},
		{Name: "-", Value: mem.NewVariadicPrimitive("-", subtract)},
		{Name: "*", Value: mem.NewVariadicPrimitive("*", multiply)},
		{Name: "/", Value: mem.NewVariadicPrimitive("/", divide)},
		{Name: "modulo", Value: mem.NewPrimitive("modulo", 2, modulo)},
		{Name: "abs", Value: mem.NewPrimitive("abs", 1, abs)},
		{Name: "min", Value: mem.NewVariadicPrimitive("min", min)},
		{Name: "max", Value: mem.NewVariadicPrimitive("max", max)},
	}
}

// add implements (+ n1 n2 ...) - returns sum of all arguments, 0 if none.
func add(args []mem.Object) (mem.Object, error) {
	result := 0.0
	for i, arg := range args {
		n, err := mem.AsNumber(arg)
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, arg.Type(),
				argPosition("+", i))
		}
		result += n.Value()
	}
	return mem.NewNumber(result), nil
}

// subtract implements (- n) or (- n1 n2 ...).
// With one argument: negation. With multiple: n1 - n2 - n3 - ...
func subtract(args []mem.Object) (mem.Object, error) {
	if len(args) == 0 {
		return nil, mem.NewArityAtLeastError("-", 1, 0)
	}

	first, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			argPosition("-", 0))
	}

	if len(args) == 1 {
		// Negation
		return mem.NewNumber(-first.Value()), nil
	}

	// Subtraction
	result := first.Value()
	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(),
				argPosition("-", i))
		}
		result -= n.Value()
	}
	return mem.NewNumber(result), nil
}

// multiply implements (* n1 n2 ...) - returns product of all arguments, 1 if none.
func multiply(args []mem.Object) (mem.Object, error) {
	result := 1.0
	for i, arg := range args {
		n, err := mem.AsNumber(arg)
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, arg.Type(),
				argPosition("*", i))
		}
		result *= n.Value()
	}
	return mem.NewNumber(result), nil
}

// divide implements (/ n) or (/ n1 n2 ...).
// With one argument: reciprocal. With multiple: n1 / n2 / n3 / ...
func divide(args []mem.Object) (mem.Object, error) {
	if len(args) == 0 {
		return nil, mem.NewArityAtLeastError("/", 1, 0)
	}

	first, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			argPosition("/", 0))
	}

	if len(args) == 1 {
		// Reciprocal
		if first.Value() == 0 {
			return nil, &mem.DivisionByZeroError{}
		}
		return mem.NewNumber(1.0 / first.Value()), nil
	}

	// Division
	result := first.Value()
	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(),
				argPosition("/", i))
		}
		if n.Value() == 0 {
			return nil, &mem.DivisionByZeroError{}
		}
		result /= n.Value()
	}
	return mem.NewNumber(result), nil
}

// modulo implements (modulo n1 n2) - returns remainder of integer division.
func modulo(args []mem.Object) (mem.Object, error) {
	a, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			argPosition("modulo", 0))
	}
	b, err := mem.AsNumber(args[1])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[1].Type(),
			argPosition("modulo", 1))
	}

	if b.Value() == 0 {
		return nil, &mem.DivisionByZeroError{}
	}

	// Use Go's modulo with proper sign handling for Scheme semantics
	aInt := int64(a.Value())
	bInt := int64(b.Value())
	result := aInt % bInt

	// Scheme's modulo always has the same sign as the divisor
	if (result < 0 && bInt > 0) || (result > 0 && bInt < 0) {
		result += bInt
	}

	return mem.NewNumber(float64(result)), nil
}

// abs implements (abs n) - returns absolute value.
func abs(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			argPosition("abs", 0))
	}

	if n.Value() < 0 {
		return mem.NewNumber(-n.Value()), nil
	}
	return n, nil
}

// min implements (min n1 n2 ...) - returns minimum value.
func min(args []mem.Object) (mem.Object, error) {
	if len(args) == 0 {
		return nil, mem.NewArityAtLeastError("min", 1, 0)
	}

	first, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			argPosition("min", 0))
	}

	result := first.Value()
	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(),
				argPosition("min", i))
		}
		if n.Value() < result {
			result = n.Value()
		}
	}
	return mem.NewNumber(result), nil
}

// max implements (max n1 n2 ...) - returns maximum value.
func max(args []mem.Object) (mem.Object, error) {
	if len(args) == 0 {
		return nil, mem.NewArityAtLeastError("max", 1, 0)
	}

	first, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			argPosition("max", 0))
	}

	result := first.Value()
	for i := 1; i < len(args); i++ {
		n, err := mem.AsNumber(args[i])
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeNumber, args[i].Type(),
				argPosition("max", i))
		}
		if n.Value() > result {
			result = n.Value()
		}
	}
	return mem.NewNumber(result), nil
}

// argPosition returns a descriptive string for error messages.
func argPosition(proc string, idx int) string {
	return fmt.Sprintf("argument %d of %s", idx+1, proc)
}

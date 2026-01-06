package math

import (
	"fmt"
	gomath "math"

	"github.com/Vallghall/schego/pkg/mem"
)

// Definitions returns all math operation definitions.
func Definitions() []mem.Definition {
	return []mem.Definition{
		// Rounding operations
		{Name: "floor", Value: mem.NewPrimitive("floor", 1, floor)},
		{Name: "ceiling", Value: mem.NewPrimitive("ceiling", 1, ceiling)},
		{Name: "round", Value: mem.NewPrimitive("round", 1, round)},
		{Name: "truncate", Value: mem.NewPrimitive("truncate", 1, truncate)},

		// Exponential and logarithmic
		{Name: "sqrt", Value: mem.NewPrimitive("sqrt", 1, sqrt)},
		{Name: "expt", Value: mem.NewPrimitive("expt", 2, expt)},
		{Name: "exp", Value: mem.NewPrimitive("exp", 1, exp)},
		{Name: "log", Value: mem.NewVariadicPrimitive("log", log)},

		// Trigonometric functions
		{Name: "sin", Value: mem.NewPrimitive("sin", 1, sin)},
		{Name: "cos", Value: mem.NewPrimitive("cos", 1, cos)},
		{Name: "tan", Value: mem.NewPrimitive("tan", 1, tan)},
		{Name: "asin", Value: mem.NewPrimitive("asin", 1, asin)},
		{Name: "acos", Value: mem.NewPrimitive("acos", 1, acos)},
		{Name: "atan", Value: mem.NewVariadicPrimitive("atan", atan)},

		// Constants (as zero-arity procedures for Scheme compatibility)
		{Name: "pi", Value: mem.NewNumber(gomath.Pi)},
		{Name: "e", Value: mem.NewNumber(gomath.E)},

		// Predicates
		{Name: "nan?", Value: mem.NewPrimitive("nan?", 1, isNaN)},
		{Name: "infinite?", Value: mem.NewPrimitive("infinite?", 1, isInfinite)},
		{Name: "finite?", Value: mem.NewPrimitive("finite?", 1, isFinite)},
		{Name: "integer?", Value: mem.NewPrimitive("integer?", 1, isInteger)},
	}
}

// floor implements (floor n) - returns largest integer <= n.
func floor(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of floor")
	}
	return mem.NewNumber(gomath.Floor(n.Value())), nil
}

// ceiling implements (ceiling n) - returns smallest integer >= n.
func ceiling(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of ceiling")
	}
	return mem.NewNumber(gomath.Ceil(n.Value())), nil
}

// round implements (round n) - returns nearest integer (banker's rounding).
func round(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of round")
	}
	return mem.NewNumber(gomath.Round(n.Value())), nil
}

// truncate implements (truncate n) - returns integer part (toward zero).
func truncate(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of truncate")
	}
	return mem.NewNumber(gomath.Trunc(n.Value())), nil
}

// sqrt implements (sqrt n) - returns square root.
func sqrt(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of sqrt")
	}

	if n.Value() < 0 {
		return nil, fmt.Errorf("sqrt: negative argument: %v", n.Value())
	}

	return mem.NewNumber(gomath.Sqrt(n.Value())), nil
}

// expt implements (expt base exponent) - returns base raised to exponent.
func expt(args []mem.Object) (mem.Object, error) {
	base, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of expt")
	}

	exp, err := mem.AsNumber(args[1])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[1].Type(),
			"argument 2 of expt")
	}

	return mem.NewNumber(gomath.Pow(base.Value(), exp.Value())), nil
}

// exp implements (exp n) - returns e raised to the power n.
func exp(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of exp")
	}
	return mem.NewNumber(gomath.Exp(n.Value())), nil
}

// log implements (log n) or (log n base) - returns natural log or log with specified base.
func log(args []mem.Object) (mem.Object, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("log: expected 1 or 2 arguments, got %d", len(args))
	}

	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of log")
	}

	if n.Value() <= 0 {
		return nil, fmt.Errorf("log: argument must be positive: %v", n.Value())
	}

	if len(args) == 1 {
		// Natural logarithm
		return mem.NewNumber(gomath.Log(n.Value())), nil
	}

	// Logarithm with specified base
	base, err := mem.AsNumber(args[1])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[1].Type(),
			"argument 2 of log")
	}

	if base.Value() <= 0 || base.Value() == 1 {
		return nil, fmt.Errorf("log: invalid base: %v", base.Value())
	}

	return mem.NewNumber(gomath.Log(n.Value()) / gomath.Log(base.Value())), nil
}

// sin implements (sin n) - returns sine of n (radians).
func sin(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of sin")
	}
	return mem.NewNumber(gomath.Sin(n.Value())), nil
}

// cos implements (cos n) - returns cosine of n (radians).
func cos(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of cos")
	}
	return mem.NewNumber(gomath.Cos(n.Value())), nil
}

// tan implements (tan n) - returns tangent of n (radians).
func tan(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of tan")
	}
	return mem.NewNumber(gomath.Tan(n.Value())), nil
}

// asin implements (asin n) - returns arcsine of n (result in radians).
func asin(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of asin")
	}

	if n.Value() < -1 || n.Value() > 1 {
		return nil, fmt.Errorf("asin: argument out of range [-1, 1]: %v", n.Value())
	}

	return mem.NewNumber(gomath.Asin(n.Value())), nil
}

// acos implements (acos n) - returns arccosine of n (result in radians).
func acos(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of acos")
	}

	if n.Value() < -1 || n.Value() > 1 {
		return nil, fmt.Errorf("acos: argument out of range [-1, 1]: %v", n.Value())
	}

	return mem.NewNumber(gomath.Acos(n.Value())), nil
}

// atan implements (atan n) or (atan y x) - returns arctangent.
// With one argument: returns arctangent of n.
// With two arguments: returns arctangent of y/x, using signs to determine quadrant.
func atan(args []mem.Object) (mem.Object, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("atan: expected 1 or 2 arguments, got %d", len(args))
	}

	y, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of atan")
	}

	if len(args) == 1 {
		return mem.NewNumber(gomath.Atan(y.Value())), nil
	}

	// Two-argument form (atan2)
	x, err := mem.AsNumber(args[1])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[1].Type(),
			"argument 2 of atan")
	}

	return mem.NewNumber(gomath.Atan2(y.Value(), x.Value())), nil
}

// isNaN implements (nan? n) - returns #t if n is NaN.
func isNaN(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of nan?")
	}
	return mem.Boolean(gomath.IsNaN(n.Value())), nil
}

// isInfinite implements (infinite? n) - returns #t if n is infinite.
func isInfinite(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of infinite?")
	}
	return mem.Boolean(gomath.IsInf(n.Value(), 0)), nil
}

// isFinite implements (finite? n) - returns #t if n is finite (not NaN or infinite).
func isFinite(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of finite?")
	}
	v := n.Value()
	return mem.Boolean(!gomath.IsNaN(v) && !gomath.IsInf(v, 0)), nil
}

// isInteger implements (integer? n) - returns #t if n is an integer.
func isInteger(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		// Non-numbers are not integers
		return mem.False, nil
	}
	return mem.Boolean(n.IsInteger()), nil
}

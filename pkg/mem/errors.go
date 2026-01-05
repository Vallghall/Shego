package mem

import (
	"fmt"

	"github.com/Vallghall/schego/pkg/atom"
)

// TypeError indicates an operation was attempted on the wrong type.
type TypeError struct {
	Expected ObjectType
	Got      ObjectType
	Message  string // Optional additional context
}

func (e *TypeError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("type error: expected %s, got %s: %s",
			e.Expected, e.Got, e.Message)
	}
	return fmt.Sprintf("type error: expected %s, got %s", e.Expected, e.Got)
}

// UnboundError indicates a variable was referenced but not defined.
type UnboundError struct {
	Name atom.ID
	Pool *atom.Pool // Optional: for resolving name to string
}

func (e *UnboundError) Error() string {
	if e.Pool != nil {
		if name, ok := e.Pool.Resolve(e.Name); ok {
			return fmt.Sprintf("unbound variable: %s", name)
		}
	}
	return fmt.Sprintf("unbound variable: #<atom:%d>", e.Name)
}

// BindingError indicates a problem creating or modifying a binding.
type BindingError struct {
	Name    atom.ID
	Pool    *atom.Pool
	Message string
}

func (e *BindingError) Error() string {
	var nameStr string
	if e.Pool != nil {
		if name, ok := e.Pool.Resolve(e.Name); ok {
			nameStr = name
		}
	}
	if nameStr == "" {
		nameStr = fmt.Sprintf("#<atom:%d>", e.Name)
	}
	return fmt.Sprintf("binding error for %s: %s", nameStr, e.Message)
}

// ArityError indicates wrong number of arguments.
type ArityError struct {
	Name     string // Procedure name (if known)
	Expected int    // Expected argument count (-1 for variadic minimum)
	Got      int    // Actual argument count
	AtLeast  bool   // True if Expected is a minimum (for variadic)
}

func (e *ArityError) Error() string {
	nameStr := "procedure"
	if e.Name != "" {
		nameStr = e.Name
	}

	if e.AtLeast {
		return fmt.Sprintf("%s: expected at least %d argument(s), got %d",
			nameStr, e.Expected, e.Got)
	}
	return fmt.Sprintf("%s: expected %d argument(s), got %d",
		nameStr, e.Expected, e.Got)
}

// NotCallableError indicates an attempt to call a non-procedure.
type NotCallableError struct {
	Value Object
}

func (e *NotCallableError) Error() string {
	return fmt.Sprintf("not callable: %s is not a procedure", e.Value.Type())
}

// DivisionByZeroError indicates division by zero.
type DivisionByZeroError struct{}

func (e *DivisionByZeroError) Error() string {
	return "division by zero"
}

// ImmutableError indicates an attempt to modify an immutable value.
type ImmutableError struct {
	Type    ObjectType
	Message string
}

func (e *ImmutableError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("immutable %s: %s", e.Type, e.Message)
	}
	return fmt.Sprintf("cannot modify immutable %s", e.Type)
}

// IndexError indicates an out-of-bounds access.
type IndexError struct {
	Index int
	Size  int
	Type  string // "list", "string", etc.
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("index %d out of bounds for %s of size %d",
		e.Index, e.Type, e.Size)
}

// SyntaxError indicates malformed input during evaluation.
type SyntaxError struct {
	Message string
	Context string // Optional: where the error occurred
}

func (e *SyntaxError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("syntax error in %s: %s", e.Context, e.Message)
	}
	return fmt.Sprintf("syntax error: %s", e.Message)
}

// RuntimeError is a general-purpose runtime error.
type RuntimeError struct {
	Message string
	Cause   error // Optional: underlying error
}

func (e *RuntimeError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("runtime error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("runtime error: %s", e.Message)
}

func (e *RuntimeError) Unwrap() error {
	return e.Cause
}

// =============================================================================
// Error Constructors
// =============================================================================

// NewTypeError creates a type error with an optional message.
func NewTypeError(expected, got ObjectType, msg string) *TypeError {
	return &TypeError{Expected: expected, Got: got, Message: msg}
}

// NewArityError creates an arity error.
func NewArityError(name string, expected, got int) *ArityError {
	return &ArityError{Name: name, Expected: expected, Got: got}
}

// NewArityAtLeastError creates an arity error for variadic procedures.
func NewArityAtLeastError(name string, minimum, got int) *ArityError {
	return &ArityError{Name: name, Expected: minimum, Got: got, AtLeast: true}
}

// NewUnboundError creates an unbound variable error.
func NewUnboundError(name atom.ID, pool *atom.Pool) *UnboundError {
	return &UnboundError{Name: name, Pool: pool}
}

// NewNotCallableError creates a not-callable error.
func NewNotCallableError(value Object) *NotCallableError {
	return &NotCallableError{Value: value}
}

// NewRuntimeError creates a general runtime error.
func NewRuntimeError(msg string) *RuntimeError {
	return &RuntimeError{Message: msg}
}

// WrapRuntimeError wraps an existing error as a runtime error.
func WrapRuntimeError(msg string, cause error) *RuntimeError {
	return &RuntimeError{Message: msg, Cause: cause}
}


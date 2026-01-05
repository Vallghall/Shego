package mem

import (
	"fmt"
	"strconv"

	"github.com/Vallghall/schego/pkg/atom"
)

// ObjectType represents the type of a Scheme object.
type ObjectType int

const (
	TypeNil ObjectType = iota
	TypeVoid
	TypeBoolean
	TypeNumber
	TypeString
	TypeSymbol
	TypePair
	TypeProcedure
)

// String returns the type name.
func (t ObjectType) String() string {
	switch t {
	case TypeNil:
		return "nil"
	case TypeVoid:
		return "void"
	case TypeBoolean:
		return "boolean"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeSymbol:
		return "symbol"
	case TypePair:
		return "pair"
	case TypeProcedure:
		return "procedure"
	default:
		return "unknown"
	}
}

// Object represents any Scheme value in memory.
// Objects encapsulate their type internally and provide type-agnostic access.
type Object interface {
	// Type returns the concrete type of this object.
	Type() ObjectType

	// IsCallable returns true if this object can be invoked as a procedure.
	IsCallable() bool

	// IsTruthy returns the boolean interpretation.
	// In Scheme, only #f is falsy; everything else is truthy.
	IsTruthy() bool

	// String returns a human-readable representation.
	String() string

	// Equal checks equality with another object following Scheme's equal? semantics.
	Equal(other Object) bool
}

// =============================================================================
// Singleton Objects
// =============================================================================

// Nil represents the empty list '()
var Nil Object = &nilObject{}

type nilObject struct{}

func (*nilObject) Type() ObjectType  { return TypeNil }
func (*nilObject) IsCallable() bool  { return false }
func (*nilObject) IsTruthy() bool    { return true } // nil is truthy in Scheme!
func (*nilObject) String() string    { return "()" }
func (*nilObject) Equal(o Object) bool { return o.Type() == TypeNil }

// Void represents no meaningful value (result of side-effect operations)
var Void Object = &voidObject{}

type voidObject struct{}

func (*voidObject) Type() ObjectType  { return TypeVoid }
func (*voidObject) IsCallable() bool  { return false }
func (*voidObject) IsTruthy() bool    { return true }
func (*voidObject) String() string    { return "#<void>" }
func (*voidObject) Equal(o Object) bool { return o.Type() == TypeVoid }

// True is the boolean true value
var True Object = &boolObject{value: true}

// False is the boolean false value (the only falsy value in Scheme)
var False Object = &boolObject{value: false}

// =============================================================================
// Boolean Object
// =============================================================================

type boolObject struct {
	value bool
}

func (*boolObject) Type() ObjectType { return TypeBoolean }
func (*boolObject) IsCallable() bool { return false }
func (b *boolObject) IsTruthy() bool { return b.value }
func (b *boolObject) String() string {
	if b.value {
		return "#t"
	}
	return "#f"
}
func (b *boolObject) Equal(o Object) bool {
	if other, ok := o.(*boolObject); ok {
		return b.value == other.value
	}
	return false
}

// Boolean returns the appropriate boolean object.
func Boolean(v bool) Object {
	if v {
		return True
	}
	return False
}

// =============================================================================
// Number Object
// =============================================================================

// Number represents a numeric value.
// Internally stored as float64 to support both integers and reals.
type Number struct {
	value float64
}

func NewNumber(v float64) *Number {
	return &Number{value: v}
}

func NewNumberFromInt(v int64) *Number {
	return &Number{value: float64(v)}
}

func (*Number) Type() ObjectType { return TypeNumber }
func (*Number) IsCallable() bool { return false }
func (*Number) IsTruthy() bool   { return true }

func (n *Number) String() string {
	// Format as integer if possible
	if n.value == float64(int64(n.value)) {
		return strconv.FormatInt(int64(n.value), 10)
	}
	return strconv.FormatFloat(n.value, 'g', -1, 64)
}

func (n *Number) Equal(o Object) bool {
	if other, ok := o.(*Number); ok {
		return n.value == other.value
	}
	return false
}

// Value returns the numeric value.
func (n *Number) Value() float64 {
	return n.value
}

// IsInteger returns true if this number has no fractional part.
func (n *Number) IsInteger() bool {
	return n.value == float64(int64(n.value))
}

// Int64 returns the integer value, truncating if necessary.
func (n *Number) Int64() int64 {
	return int64(n.value)
}

// =============================================================================
// String Object
// =============================================================================

// String represents a string value.
type String struct {
	value string
}

func NewString(v string) *String {
	return &String{value: v}
}

func (*String) Type() ObjectType { return TypeString }
func (*String) IsCallable() bool { return false }
func (*String) IsTruthy() bool   { return true }

func (s *String) String() string {
	return fmt.Sprintf("%q", s.value)
}

func (s *String) Equal(o Object) bool {
	if other, ok := o.(*String); ok {
		return s.value == other.value
	}
	return false
}

// Value returns the string value.
func (s *String) Value() string {
	return s.value
}

// =============================================================================
// Symbol Object
// =============================================================================

// Symbol represents an interned symbol.
type Symbol struct {
	id   atom.ID
	pool *atom.Pool
}

func NewSymbol(id atom.ID, pool *atom.Pool) *Symbol {
	return &Symbol{id: id, pool: pool}
}

func (*Symbol) Type() ObjectType { return TypeSymbol }
func (*Symbol) IsCallable() bool { return false }
func (*Symbol) IsTruthy() bool   { return true }

func (s *Symbol) String() string {
	if s.pool != nil {
		return s.pool.MustResolve(s.id)
	}
	return fmt.Sprintf("#<symbol:%d>", s.id)
}

func (s *Symbol) Equal(o Object) bool {
	if other, ok := o.(*Symbol); ok {
		return s.id == other.id
	}
	return false
}

// ID returns the atom ID.
func (s *Symbol) ID() atom.ID {
	return s.id
}

// Name returns the symbol name.
func (s *Symbol) Name() string {
	if s.pool != nil {
		return s.pool.MustResolve(s.id)
	}
	return ""
}

// =============================================================================
// Pair Object
// =============================================================================

// Pair represents a cons cell (car . cdr).
type Pair struct {
	car Object
	cdr Object
}

func NewPair(car, cdr Object) *Pair {
	return &Pair{car: car, cdr: cdr}
}

func (*Pair) Type() ObjectType { return TypePair }
func (*Pair) IsCallable() bool { return false }
func (*Pair) IsTruthy() bool   { return true }

func (p *Pair) String() string {
	return formatList(p)
}

func (p *Pair) Equal(o Object) bool {
	if other, ok := o.(*Pair); ok {
		return p.car.Equal(other.car) && p.cdr.Equal(other.cdr)
	}
	return false
}

// Car returns the first element.
func (p *Pair) Car() Object {
	return p.car
}

// Cdr returns the second element.
func (p *Pair) Cdr() Object {
	return p.cdr
}

// SetCar sets the first element (for set-car!).
func (p *Pair) SetCar(v Object) {
	p.car = v
}

// SetCdr sets the second element (for set-cdr!).
func (p *Pair) SetCdr(v Object) {
	p.cdr = v
}

// formatList formats a pair as a proper or improper list.
func formatList(p *Pair) string {
	var elements []string
	current := Object(p)

	for {
		if pair, ok := current.(*Pair); ok {
			elements = append(elements, pair.car.String())
			current = pair.cdr
		} else if current.Type() == TypeNil {
			// Proper list
			return "(" + joinStrings(elements, " ") + ")"
		} else {
			// Improper list
			return "(" + joinStrings(elements, " ") + " . " + current.String() + ")"
		}
	}
}

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for _, s := range ss[1:] {
		result += sep + s
	}
	return result
}

// =============================================================================
// Procedure Object
// =============================================================================

// ProcedureType distinguishes primitive vs user-defined procedures.
type ProcedureType int

const (
	ProcPrimitive ProcedureType = iota // Built-in Go function
	ProcLambda                         // User-defined lambda
)

// Procedure represents a callable object.
// It captures its definition context for lexical scoping.
type Procedure struct {
	kind    ProcedureType
	name    string   // Optional name for debugging
	arity   int      // Expected argument count, -1 for variadic
	context Context  // Lexical environment where defined (for lambdas)
	params  []atom.ID // Parameter names (for lambdas)
	body    []any    // Body expressions (for lambdas) - stored as AST nodes
	native  any      // Native function (for primitives)
}

func (*Procedure) Type() ObjectType { return TypeProcedure }
func (*Procedure) IsCallable() bool { return true }
func (*Procedure) IsTruthy() bool   { return true }

func (p *Procedure) String() string {
	if p.name != "" {
		return fmt.Sprintf("#<procedure:%s>", p.name)
	}
	return "#<procedure>"
}

func (p *Procedure) Equal(o Object) bool {
	// Procedures are only equal to themselves (identity)
	return p == o
}

// Kind returns whether this is a primitive or lambda.
func (p *Procedure) Kind() ProcedureType {
	return p.kind
}

// Name returns the procedure name.
func (p *Procedure) Name() string {
	return p.name
}

// Arity returns the expected argument count (-1 for variadic).
func (p *Procedure) Arity() int {
	return p.arity
}

// Context returns the captured lexical environment.
func (p *Procedure) Context() Context {
	return p.context
}

// Params returns the parameter names (for lambdas).
func (p *Procedure) Params() []atom.ID {
	return p.params
}

// Body returns the body expressions (for lambdas).
func (p *Procedure) Body() []any {
	return p.body
}

// Native returns the native function (for primitives).
func (p *Procedure) Native() any {
	return p.native
}

// NewPrimitive creates a primitive procedure.
func NewPrimitive(name string, arity int, fn PrimitiveFunc) *Procedure {
	return &Procedure{
		kind:   ProcPrimitive,
		name:   name,
		arity:  arity,
		native: fn,
	}
}

// NewLambda creates a user-defined procedure.
func NewLambda(name string, params []atom.ID, body []any, ctx Context) *Procedure {
	return &Procedure{
		kind:    ProcLambda,
		name:    name,
		arity:   len(params),
		context: ctx,
		params:  params,
		body:    body,
	}
}

// NewVariadicPrimitive creates a variadic primitive procedure.
func NewVariadicPrimitive(name string, fn PrimitiveFunc) *Procedure {
	return &Procedure{
		kind:   ProcPrimitive,
		name:   name,
		arity:  -1, // variadic
		native: fn,
	}
}

// =============================================================================
// Type Coercion Helpers
// =============================================================================

// AsNumber attempts to convert an object to a number.
// Returns an error if conversion is not possible.
func AsNumber(o Object) (*Number, error) {
	if n, ok := o.(*Number); ok {
		return n, nil
	}
	return nil, &TypeError{Expected: TypeNumber, Got: o.Type()}
}

// AsString attempts to convert an object to a string.
// Returns an error if conversion is not possible.
func AsString(o Object) (*String, error) {
	if s, ok := o.(*String); ok {
		return s, nil
	}
	return nil, &TypeError{Expected: TypeString, Got: o.Type()}
}

// AsSymbol attempts to convert an object to a symbol.
func AsSymbol(o Object) (*Symbol, error) {
	if s, ok := o.(*Symbol); ok {
		return s, nil
	}
	return nil, &TypeError{Expected: TypeSymbol, Got: o.Type()}
}

// AsPair attempts to convert an object to a pair.
func AsPair(o Object) (*Pair, error) {
	if p, ok := o.(*Pair); ok {
		return p, nil
	}
	return nil, &TypeError{Expected: TypePair, Got: o.Type()}
}

// AsProcedure attempts to convert an object to a procedure.
func AsProcedure(o Object) (*Procedure, error) {
	if p, ok := o.(*Procedure); ok {
		return p, nil
	}
	return nil, &TypeError{Expected: TypeProcedure, Got: o.Type()}
}

// IsList checks if an object is a proper list (nil-terminated).
func IsList(o Object) bool {
	current := o
	for {
		switch current.Type() {
		case TypeNil:
			return true
		case TypePair:
			current = current.(*Pair).cdr
		default:
			return false
		}
	}
}

// ListToSlice converts a proper list to a Go slice.
// Returns an error if not a proper list.
func ListToSlice(o Object) ([]Object, error) {
	var result []Object
	current := o

	for {
		switch current.Type() {
		case TypeNil:
			return result, nil
		case TypePair:
			pair := current.(*Pair)
			result = append(result, pair.car)
			current = pair.cdr
		default:
			return nil, fmt.Errorf("not a proper list: ends with %s", current.Type())
		}
	}
}

// SliceToList converts a Go slice to a proper list.
func SliceToList(elements []Object) Object {
	result := Nil
	for i := len(elements) - 1; i >= 0; i-- {
		result = NewPair(elements[i], result)
	}
	return result
}


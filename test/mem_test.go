package test

import (
	"testing"

	"github.com/Vallghall/schego/pkg/atom"
	"github.com/Vallghall/schego/pkg/mem"
)

func TestObjects(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		if mem.Nil.Type() != mem.TypeNil {
			t.Errorf("expected TypeNil, got %v", mem.Nil.Type())
		}
		if !mem.Nil.IsTruthy() {
			t.Error("nil should be truthy in Scheme")
		}
		if mem.Nil.String() != "()" {
			t.Errorf("expected '()', got %q", mem.Nil.String())
		}
	})

	t.Run("Void", func(t *testing.T) {
		if mem.Void.Type() != mem.TypeVoid {
			t.Errorf("expected TypeVoid, got %v", mem.Void.Type())
		}
		if mem.Void.String() != "#<void>" {
			t.Errorf("expected '#<void>', got %q", mem.Void.String())
		}
	})

	t.Run("Boolean", func(t *testing.T) {
		if !mem.True.IsTruthy() {
			t.Error("True should be truthy")
		}
		if mem.False.IsTruthy() {
			t.Error("False should not be truthy")
		}
		if mem.True.String() != "#t" {
			t.Errorf("expected '#t', got %q", mem.True.String())
		}
		if mem.False.String() != "#f" {
			t.Errorf("expected '#f', got %q", mem.False.String())
		}
		if !mem.True.Equal(mem.Boolean(true)) {
			t.Error("True should equal Boolean(true)")
		}
	})

	t.Run("Number", func(t *testing.T) {
		n := mem.NewNumber(42.0)
		if n.Type() != mem.TypeNumber {
			t.Errorf("expected TypeNumber, got %v", n.Type())
		}
		if n.Value() != 42.0 {
			t.Errorf("expected 42.0, got %v", n.Value())
		}
		if !n.IsInteger() {
			t.Error("42.0 should be integer")
		}
		if n.String() != "42" {
			t.Errorf("expected '42', got %q", n.String())
		}

		f := mem.NewNumber(3.14)
		if f.IsInteger() {
			t.Error("3.14 should not be integer")
		}

		n2 := mem.NewNumber(42.0)
		if !n.Equal(n2) {
			t.Error("equal numbers should be equal")
		}
	})

	t.Run("String", func(t *testing.T) {
		s := mem.NewString("hello")
		if s.Type() != mem.TypeString {
			t.Errorf("expected TypeString, got %v", s.Type())
		}
		if s.Value() != "hello" {
			t.Errorf("expected 'hello', got %q", s.Value())
		}
		if s.String() != `"hello"` {
			t.Errorf("expected '\"hello\"', got %q", s.String())
		}

		s2 := mem.NewString("hello")
		if !s.Equal(s2) {
			t.Error("equal strings should be equal")
		}
	})

	t.Run("Symbol", func(t *testing.T) {
		pool := atom.NewPool()
		id := pool.Intern("foo")
		sym := mem.NewSymbol(id, pool)

		if sym.Type() != mem.TypeSymbol {
			t.Errorf("expected TypeSymbol, got %v", sym.Type())
		}
		if sym.ID() != id {
			t.Errorf("expected ID %d, got %d", id, sym.ID())
		}
		if sym.Name() != "foo" {
			t.Errorf("expected 'foo', got %q", sym.Name())
		}
		if sym.String() != "foo" {
			t.Errorf("expected 'foo', got %q", sym.String())
		}

		sym2 := mem.NewSymbol(id, pool)
		if !sym.Equal(sym2) {
			t.Error("symbols with same ID should be equal")
		}
	})

	t.Run("Pair", func(t *testing.T) {
		p := mem.NewPair(mem.NewNumber(1), mem.NewNumber(2))
		if p.Type() != mem.TypePair {
			t.Errorf("expected TypePair, got %v", p.Type())
		}

		car, _ := mem.AsNumber(p.Car())
		cdr, _ := mem.AsNumber(p.Cdr())
		if car.Value() != 1 || cdr.Value() != 2 {
			t.Errorf("expected (1 . 2), got (%v . %v)", car.Value(), cdr.Value())
		}

		// Test improper list string
		if p.String() != "(1 . 2)" {
			t.Errorf("expected '(1 . 2)', got %q", p.String())
		}
	})

	t.Run("ProperList", func(t *testing.T) {
		list := mem.NewPair(
			mem.NewNumber(1),
			mem.NewPair(
				mem.NewNumber(2),
				mem.NewPair(
					mem.NewNumber(3),
					mem.Nil,
				),
			),
		)

		if !mem.IsList(list) {
			t.Error("should be a proper list")
		}
		if list.String() != "(1 2 3)" {
			t.Errorf("expected '(1 2 3)', got %q", list.String())
		}

		elements, err := mem.ListToSlice(list)
		if err != nil {
			t.Fatalf("ListToSlice error: %v", err)
		}
		if len(elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(elements))
		}
	})

	t.Run("SliceToList", func(t *testing.T) {
		slice := []mem.Object{
			mem.NewNumber(1),
			mem.NewNumber(2),
			mem.NewNumber(3),
		}
		list := mem.SliceToList(slice)

		if !mem.IsList(list) {
			t.Error("should be a proper list")
		}

		elements, _ := mem.ListToSlice(list)
		if len(elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(elements))
		}
	})

	t.Run("Procedure", func(t *testing.T) {
		prim := mem.NewPrimitive("+", 2, nil)
		if prim.Type() != mem.TypeProcedure {
			t.Errorf("expected TypeProcedure, got %v", prim.Type())
		}
		if !prim.IsCallable() {
			t.Error("procedure should be callable")
		}
		if prim.Kind() != mem.ProcPrimitive {
			t.Error("expected ProcPrimitive")
		}
		if prim.Arity() != 2 {
			t.Errorf("expected arity 2, got %d", prim.Arity())
		}

		variadic := mem.NewVariadicPrimitive("list", nil)
		if variadic.Arity() != -1 {
			t.Errorf("variadic should have arity -1, got %d", variadic.Arity())
		}
	})

	t.Run("TypeCoercion", func(t *testing.T) {
		n := mem.NewNumber(42)
		_, err := mem.AsNumber(n)
		if err != nil {
			t.Errorf("AsNumber should succeed: %v", err)
		}

		_, err = mem.AsString(n)
		if err == nil {
			t.Error("AsString on number should fail")
		}
	})
}

func TestContext(t *testing.T) {
	pool := atom.NewPool()
	x := pool.Intern("x")
	y := pool.Intern("y")
	z := pool.Intern("z")

	t.Run("Define and Lookup", func(t *testing.T) {
		ctx := mem.NewContext()

		err := ctx.Define(x, mem.NewNumber(42))
		if err != nil {
			t.Fatalf("Define error: %v", err)
		}

		val, ok := ctx.Lookup(x)
		if !ok {
			t.Fatal("Lookup failed")
		}
		n, _ := mem.AsNumber(val)
		if n.Value() != 42 {
			t.Errorf("expected 42, got %v", n.Value())
		}
	})

	t.Run("Lookup Missing", func(t *testing.T) {
		ctx := mem.NewContext()

		_, ok := ctx.Lookup(z)
		if ok {
			t.Error("should not find undefined variable")
		}
	})

	t.Run("Set Existing", func(t *testing.T) {
		ctx := mem.NewContext()
		ctx.Define(x, mem.NewNumber(1))

		err := ctx.Set(x, mem.NewNumber(2))
		if err != nil {
			t.Fatalf("Set error: %v", err)
		}

		val, _ := ctx.Lookup(x)
		n, _ := mem.AsNumber(val)
		if n.Value() != 2 {
			t.Errorf("expected 2, got %v", n.Value())
		}
	})

	t.Run("Set Missing", func(t *testing.T) {
		ctx := mem.NewContext()

		err := ctx.Set(x, mem.NewNumber(1))
		if err == nil {
			t.Error("Set on undefined should fail")
		}
	})

	t.Run("Nested Scopes", func(t *testing.T) {
		parent := mem.NewContext()
		parent.Define(x, mem.NewNumber(1))

		child := parent.Extend()
		child.Define(y, mem.NewNumber(2))

		// Child can see parent's binding
		valX, ok := child.Lookup(x)
		if !ok {
			t.Fatal("child should see parent's x")
		}
		nX, _ := mem.AsNumber(valX)
		if nX.Value() != 1 {
			t.Errorf("expected x=1, got %v", nX.Value())
		}

		// Child has its own binding
		valY, ok := child.Lookup(y)
		if !ok {
			t.Fatal("child should have y")
		}
		nY, _ := mem.AsNumber(valY)
		if nY.Value() != 2 {
			t.Errorf("expected y=2, got %v", nY.Value())
		}

		// Parent doesn't see child's binding
		_, ok = parent.Lookup(y)
		if ok {
			t.Error("parent should not see child's y")
		}
	})

	t.Run("Shadow Parent", func(t *testing.T) {
		parent := mem.NewContext()
		parent.Define(x, mem.NewNumber(1))

		child := parent.Extend()
		child.Define(x, mem.NewNumber(2)) // Shadow

		// Child sees its own binding
		valChild, _ := child.Lookup(x)
		nChild, _ := mem.AsNumber(valChild)
		if nChild.Value() != 2 {
			t.Errorf("child should see x=2, got %v", nChild.Value())
		}

		// Parent still sees original
		valParent, _ := parent.Lookup(x)
		nParent, _ := mem.AsNumber(valParent)
		if nParent.Value() != 1 {
			t.Errorf("parent should still see x=1, got %v", nParent.Value())
		}
	})

	t.Run("Set in Parent Scope", func(t *testing.T) {
		parent := mem.NewContext()
		parent.Define(x, mem.NewNumber(1))

		child := parent.Extend()

		// Set modifies parent's binding
		err := child.Set(x, mem.NewNumber(99))
		if err != nil {
			t.Fatalf("Set error: %v", err)
		}

		// Both see the change
		valParent, _ := parent.Lookup(x)
		nParent, _ := mem.AsNumber(valParent)
		if nParent.Value() != 99 {
			t.Errorf("parent should see x=99, got %v", nParent.Value())
		}

		valChild, _ := child.Lookup(x)
		nChild, _ := mem.AsNumber(valChild)
		if nChild.Value() != 99 {
			t.Errorf("child should see x=99, got %v", nChild.Value())
		}
	})

	t.Run("IsDefined", func(t *testing.T) {
		parent := mem.NewContext()
		parent.Define(x, mem.NewNumber(1))

		child := parent.Extend()
		child.Define(y, mem.NewNumber(2))

		if !child.IsDefined(y) {
			t.Error("child should have y defined")
		}
		if child.IsDefined(x) {
			t.Error("x is in parent, not child's immediate scope")
		}
	})

	t.Run("ExtendWith", func(t *testing.T) {
		parent := mem.NewContext()
		names := []atom.ID{x, y}
		values := []mem.Object{mem.NewNumber(1), mem.NewNumber(2)}

		child, err := mem.ExtendWith(parent, names, values)
		if err != nil {
			t.Fatalf("ExtendWith error: %v", err)
		}

		valX, _ := child.Lookup(x)
		valY, _ := child.Lookup(y)
		nX, _ := mem.AsNumber(valX)
		nY, _ := mem.AsNumber(valY)

		if nX.Value() != 1 || nY.Value() != 2 {
			t.Errorf("expected x=1, y=2, got x=%v, y=%v", nX.Value(), nY.Value())
		}
	})

	t.Run("ExtendWith Arity Mismatch", func(t *testing.T) {
		parent := mem.NewContext()
		names := []atom.ID{x, y}
		values := []mem.Object{mem.NewNumber(1)} // Only 1 value

		_, err := mem.ExtendWith(parent, names, values)
		if err == nil {
			t.Error("should fail with arity mismatch")
		}
	})
}

func TestErrors(t *testing.T) {
	t.Run("TypeError", func(t *testing.T) {
		err := mem.NewTypeError(mem.TypeNumber, mem.TypeString, "in addition")
		if err.Error() != "type error: expected number, got string: in addition" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("UnboundError", func(t *testing.T) {
		pool := atom.NewPool()
		id := pool.Intern("foo")
		err := mem.NewUnboundError(id, pool)
		if err.Error() != "unbound variable: foo" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("ArityError", func(t *testing.T) {
		err := mem.NewArityError("+", 2, 3)
		if err.Error() != "+: expected 2 argument(s), got 3" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("NotCallableError", func(t *testing.T) {
		err := mem.NewNotCallableError(mem.NewNumber(42))
		if err.Error() != "not callable: number is not a procedure" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}


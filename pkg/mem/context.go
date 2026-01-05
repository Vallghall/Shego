package mem

import (
	"github.com/Vallghall/schego/pkg/atom"
)

// Context represents a lexical environment for variable bindings.
// It maps symbolic atoms to objects and supports nested scopes.
type Context interface {
	// Define creates a new binding in this context.
	// If the name already exists in this immediate context, it is overwritten.
	// Returns an error if the binding cannot be created.
	Define(name atom.ID, value Object) error

	// Lookup finds a binding by searching this context and parent scopes.
	// Returns the bound value and true if found, or nil and false if not found.
	Lookup(name atom.ID) (Object, bool)

	// Set modifies an existing binding in the nearest enclosing scope.
	// Returns an error if the binding does not exist.
	Set(name atom.ID, value Object) error

	// Parent returns the enclosing context, or nil for top-level.
	Parent() Context

	// Extend creates a new child context with this context as parent.
	Extend() Context

	// IsDefined checks if a name is defined in this immediate context (not parents).
	IsDefined(name atom.ID) bool

	// Names returns all names defined in this immediate context.
	Names() []atom.ID
}

// =============================================================================
// MapContext Implementation
// =============================================================================

// MapContext is a Context implementation using a map for bindings.
type MapContext struct {
	bindings map[atom.ID]Object
	parent   Context
}

// NewContext creates a new top-level context with no parent.
func NewContext() *MapContext {
	return &MapContext{
		bindings: make(map[atom.ID]Object),
		parent:   nil,
	}
}

// NewChildContext creates a new context with the given parent.
func NewChildContext(parent Context) *MapContext {
	return &MapContext{
		bindings: make(map[atom.ID]Object),
		parent:   parent,
	}
}

// Define creates or overwrites a binding in this context.
func (c *MapContext) Define(name atom.ID, value Object) error {
	if value == nil {
		return &BindingError{
			Name:    name,
			Message: "cannot bind nil value",
		}
	}
	c.bindings[name] = value
	return nil
}

// Lookup searches for a binding in this context and parent scopes.
func (c *MapContext) Lookup(name atom.ID) (Object, bool) {
	// Check this context first
	if value, ok := c.bindings[name]; ok {
		return value, true
	}

	// Search parent scopes
	if c.parent != nil {
		return c.parent.Lookup(name)
	}

	return nil, false
}

// Set modifies an existing binding in the nearest enclosing scope.
func (c *MapContext) Set(name atom.ID, value Object) error {
	if value == nil {
		return &BindingError{
			Name:    name,
			Message: "cannot set to nil value",
		}
	}

	// Check this context first
	if _, ok := c.bindings[name]; ok {
		c.bindings[name] = value
		return nil
	}

	// Try parent scopes
	if c.parent != nil {
		return c.parent.Set(name, value)
	}

	return &UnboundError{Name: name}
}

// Parent returns the enclosing context.
func (c *MapContext) Parent() Context {
	return c.parent
}

// Extend creates a new child context.
func (c *MapContext) Extend() Context {
	return NewChildContext(c)
}

// IsDefined checks if a name is defined in this immediate context.
func (c *MapContext) IsDefined(name atom.ID) bool {
	_, ok := c.bindings[name]
	return ok
}

// Names returns all names defined in this immediate context.
func (c *MapContext) Names() []atom.ID {
	names := make([]atom.ID, 0, len(c.bindings))
	for name := range c.bindings {
		names = append(names, name)
	}
	return names
}

// Size returns the number of bindings in this immediate context.
func (c *MapContext) Size() int {
	return len(c.bindings)
}

// Depth returns the nesting depth of this context (0 for top-level).
func (c *MapContext) Depth() int {
	depth := 0
	current := c.parent
	for current != nil {
		depth++
		current = current.Parent()
	}
	return depth
}

// =============================================================================
// Helper Functions
// =============================================================================

// DefineMultiple defines multiple bindings at once.
// Useful for setting up initial environments with primitives.
func DefineMultiple(ctx Context, bindings map[atom.ID]Object) error {
	for name, value := range bindings {
		if err := ctx.Define(name, value); err != nil {
			return err
		}
	}
	return nil
}

// ExtendWith creates a child context with the given bindings.
// Useful for creating procedure application environments.
func ExtendWith(parent Context, names []atom.ID, values []Object) (Context, error) {
	if len(names) != len(values) {
		return nil, &ArityError{
			Expected: len(names),
			Got:      len(values),
		}
	}

	child := parent.Extend()
	for i, name := range names {
		if err := child.Define(name, values[i]); err != nil {
			return nil, err
		}
	}
	return child, nil
}

// LookupOrError looks up a binding and returns a descriptive error if not found.
func LookupOrError(ctx Context, name atom.ID, pool *atom.Pool) (Object, error) {
	value, ok := ctx.Lookup(name)
	if !ok {
		return nil, &UnboundError{Name: name, Pool: pool}
	}
	return value, nil
}

// CopyContext creates a shallow copy of a context's immediate bindings.
// The parent reference is preserved, not copied.
func CopyContext(ctx Context) Context {
	if mapCtx, ok := ctx.(*MapContext); ok {
		newCtx := NewChildContext(mapCtx.parent)
		for name, value := range mapCtx.bindings {
			newCtx.bindings[name] = value
		}
		return newCtx
	}
	// Fallback: create new context and copy bindings
	newCtx := NewChildContext(ctx.Parent())
	for _, name := range ctx.Names() {
		if value, ok := ctx.Lookup(name); ok {
			newCtx.Define(name, value)
		}
	}
	return newCtx
}

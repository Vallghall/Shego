package eval

import (
	"github.com/Vallghall/schego/pkg/atom"
	"github.com/Vallghall/schego/pkg/mem"
)

// State holds the global application state for the interpreter.
// It manages the atom pool for symbol interning and the global context
// for variable bindings.
type State struct {
	pool    *atom.Pool
	global  mem.Context
	current mem.Context // Current evaluation context (may differ during calls)
}

// NewState creates a new interpreter state with empty global context.
func NewState() *State {
	pool := atom.NewPool()
	global := mem.NewContext()
	return &State{
		pool:    pool,
		global:  global,
		current: global,
	}
}

// NewStateWithPool creates a new interpreter state using an existing atom pool.
func NewStateWithPool(pool *atom.Pool) *State {
	global := mem.NewContext()
	return &State{
		pool:    pool,
		global:  global,
		current: global,
	}
}

// Pool returns the atom pool.
func (s *State) Pool() *atom.Pool {
	return s.pool
}

// Global returns the global context.
func (s *State) Global() mem.Context {
	return s.global
}

// Current returns the current evaluation context.
func (s *State) Current() mem.Context {
	return s.current
}

// SetCurrent sets the current evaluation context.
// Used when entering/exiting procedure calls.
func (s *State) SetCurrent(ctx mem.Context) {
	s.current = ctx
}

// ResetCurrent resets the current context to global.
func (s *State) ResetCurrent() {
	s.current = s.global
}

// =============================================================================
// Symbol Management
// =============================================================================

// Intern interns a string as an atom and returns its ID.
func (s *State) Intern(name string) atom.ID {
	return s.pool.Intern(name)
}

// Resolve returns the string for an atom ID.
func (s *State) Resolve(id atom.ID) string {
	return s.pool.MustResolve(id)
}

// =============================================================================
// Binding Operations (operate on current context)
// =============================================================================

// Define creates a new binding in the current context.
func (s *State) Define(name atom.ID, value mem.Object) error {
	return s.current.Define(name, value)
}

// DefineByName creates a new binding using a string name.
func (s *State) DefineByName(name string, value mem.Object) error {
	return s.current.Define(s.pool.Intern(name), value)
}

// Lookup finds a binding in the current context (searching parent scopes).
func (s *State) Lookup(name atom.ID) (mem.Object, bool) {
	return s.current.Lookup(name)
}

// LookupByName finds a binding using a string name.
func (s *State) LookupByName(name string) (mem.Object, bool) {
	id, ok := s.pool.Lookup(name)
	if !ok {
		return nil, false
	}
	return s.current.Lookup(id)
}

// LookupOrError looks up a binding and returns an error if not found.
func (s *State) LookupOrError(name atom.ID) (mem.Object, error) {
	value, ok := s.current.Lookup(name)
	if !ok {
		return nil, mem.NewUnboundError(name, s.pool)
	}
	return value, nil
}

// Set modifies an existing binding in the current context.
func (s *State) Set(name atom.ID, value mem.Object) error {
	return s.current.Set(name, value)
}

// SetByName modifies an existing binding using a string name.
func (s *State) SetByName(name string, value mem.Object) error {
	id, ok := s.pool.Lookup(name)
	if !ok {
		return mem.NewUnboundError(0, nil)
	}
	return s.current.Set(id, value)
}

// =============================================================================
// Global Binding Operations (always operate on global context)
// =============================================================================

// DefineGlobal creates a binding in the global context.
func (s *State) DefineGlobal(name atom.ID, value mem.Object) error {
	return s.global.Define(name, value)
}

// DefineGlobalByName creates a binding in the global context using a string name.
func (s *State) DefineGlobalByName(name string, value mem.Object) error {
	return s.global.Define(s.pool.Intern(name), value)
}

// LookupGlobal finds a binding in the global context only.
func (s *State) LookupGlobal(name atom.ID) (mem.Object, bool) {
	return s.global.Lookup(name)
}

// =============================================================================
// Context Management
// =============================================================================

// ExtendCurrent creates a new context extending the current one.
func (s *State) ExtendCurrent() mem.Context {
	return s.current.Extend()
}

// ExtendCurrentWith creates a new context with parameter bindings.
func (s *State) ExtendCurrentWith(names []atom.ID, values []mem.Object) (mem.Context, error) {
	return mem.ExtendWith(s.current, names, values)
}

// WithContext executes a function with a temporary context.
// The original context is restored after the function returns.
func (s *State) WithContext(ctx mem.Context, fn func() error) error {
	saved := s.current
	s.current = ctx
	defer func() { s.current = saved }()
	return fn()
}

// =============================================================================
// Bulk Operations
// =============================================================================

// LoadDefinitions loads multiple definitions into the global context.
func (s *State) LoadDefinitions(defs []mem.Definition) error {
	for _, def := range defs {
		if err := s.DefineGlobalByName(def.Name, def.Value); err != nil {
			return err
		}
	}
	return nil
}

// LoadBuiltins loads built-in procedures into the global context.
// The provider function should return a slice of definitions.
func (s *State) LoadBuiltins(provider func() []mem.Definition) error {
	defs := provider()
	return s.LoadDefinitions(defs)
}


package atom

import (
	"sync"
)

// ID is a unique identifier for an interned atom.
type ID uint64

// Pool is a thread-safe atom interning pool.
// It stores atoms (strings) and assigns them unique IDs for fast access.
// Once an atom is added to the pool, it remains there for the program's lifetime.
type Pool struct {
	mu     sync.RWMutex
	toID   map[string]ID // string -> ID lookup
	toAtom []string      // ID -> string lookup (ID is the index)
	nextID ID
}

// NewPool creates a new empty atom pool.
func NewPool() *Pool {
	return &Pool{
		toID:   make(map[string]ID),
		toAtom: make([]string, 0, 256), // pre-allocate for common atoms
		nextID: 0,
	}
}

// Intern adds an atom to the pool if not present and returns its ID.
// If the atom already exists, returns the existing ID.
// This operation is thread-safe.
func (p *Pool) Intern(s string) ID {
	// Fast path: check if already interned (read lock)
	p.mu.RLock()
	if id, ok := p.toID[s]; ok {
		p.mu.RUnlock()
		return id
	}
	p.mu.RUnlock()

	// Slow path: need to add (write lock)
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if id, ok := p.toID[s]; ok {
		return id
	}

	id := p.nextID
	p.nextID++
	p.toID[s] = id
	p.toAtom = append(p.toAtom, s)

	return id
}

// Lookup returns the ID for an atom if it exists in the pool.
// Returns the ID and true if found, or 0 and false if not found.
func (p *Pool) Lookup(s string) (ID, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	id, ok := p.toID[s]
	return id, ok
}

// Resolve returns the atom string for a given ID.
// Returns the string and true if the ID is valid, or empty string and false otherwise.
func (p *Pool) Resolve(id ID) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if int(id) >= len(p.toAtom) {
		return "", false
	}
	return p.toAtom[id], true
}

// MustResolve returns the atom string for a given ID.
// Panics if the ID is not valid.
func (p *Pool) MustResolve(id ID) string {
	s, ok := p.Resolve(id)
	if !ok {
		panic("atom: invalid ID")
	}
	return s
}

// Size returns the number of atoms in the pool.
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.toAtom)
}

// Contains checks if an atom exists in the pool.
func (p *Pool) Contains(s string) bool {
	_, ok := p.Lookup(s)
	return ok
}

// All returns a copy of all interned atoms in ID order.
func (p *Pool) All() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]string, len(p.toAtom))
	copy(result, p.toAtom)
	return result
}

// Global is the default global atom pool.
// It can be used for convenience when a single pool is sufficient.
var Global = NewPool()

// Intern adds an atom to the global pool.
func Intern(s string) ID {
	return Global.Intern(s)
}

// Resolve returns the atom for an ID from the global pool.
func Resolve(id ID) (string, bool) {
	return Global.Resolve(id)
}

// MustResolve returns the atom for an ID from the global pool, panicking if invalid.
func MustResolve(id ID) string {
	return Global.MustResolve(id)
}

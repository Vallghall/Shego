package test

import (
	"sync"
	"testing"

	"github.com/Vallghall/schego/pkg/atom"
)

func TestPool(t *testing.T) {
	t.Run("BasicIntern", func(t *testing.T) {
		pool := atom.NewPool()

		id1 := pool.Intern("foo")
		id2 := pool.Intern("bar")
		id3 := pool.Intern("foo") // same as id1

		if id1 != id3 {
			t.Errorf("expected same ID for 'foo', got %d and %d", id1, id3)
		}
		if id1 == id2 {
			t.Errorf("expected different IDs for 'foo' and 'bar'")
		}
	})

	t.Run("Resolve", func(t *testing.T) {
		pool := atom.NewPool()

		id := pool.Intern("hello")
		s, ok := pool.Resolve(id)
		if !ok {
			t.Fatal("expected to resolve valid ID")
		}
		if s != "hello" {
			t.Errorf("expected 'hello', got %q", s)
		}
	})

	t.Run("ResolveInvalid", func(t *testing.T) {
		pool := atom.NewPool()

		_, ok := pool.Resolve(999)
		if ok {
			t.Error("expected false for invalid ID")
		}
	})

	t.Run("MustResolve", func(t *testing.T) {
		pool := atom.NewPool()

		id := pool.Intern("test")
		s := pool.MustResolve(id)
		if s != "test" {
			t.Errorf("expected 'test', got %q", s)
		}
	})

	t.Run("MustResolvePanic", func(t *testing.T) {
		pool := atom.NewPool()

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid ID")
			}
		}()

		pool.MustResolve(999)
	})

	t.Run("Lookup", func(t *testing.T) {
		pool := atom.NewPool()

		_, ok := pool.Lookup("missing")
		if ok {
			t.Error("expected false for non-existent atom")
		}

		pool.Intern("present")
		id, ok := pool.Lookup("present")
		if !ok {
			t.Error("expected true for existing atom")
		}
		if id != 0 {
			t.Errorf("expected ID 0, got %d", id)
		}
	})

	t.Run("Size", func(t *testing.T) {
		pool := atom.NewPool()

		if pool.Size() != 0 {
			t.Errorf("expected size 0, got %d", pool.Size())
		}

		pool.Intern("a")
		pool.Intern("b")
		pool.Intern("a") // duplicate

		if pool.Size() != 2 {
			t.Errorf("expected size 2, got %d", pool.Size())
		}
	})

	t.Run("Contains", func(t *testing.T) {
		pool := atom.NewPool()

		if pool.Contains("x") {
			t.Error("expected false for non-existent")
		}

		pool.Intern("x")
		if !pool.Contains("x") {
			t.Error("expected true for existing")
		}
	})

	t.Run("All", func(t *testing.T) {
		pool := atom.NewPool()

		pool.Intern("alpha")
		pool.Intern("beta")
		pool.Intern("gamma")

		all := pool.All()
		if len(all) != 3 {
			t.Fatalf("expected 3 atoms, got %d", len(all))
		}

		// IDs are assigned in order
		expected := []string{"alpha", "beta", "gamma"}
		for i, s := range expected {
			if all[i] != s {
				t.Errorf("atom %d: expected %q, got %q", i, s, all[i])
			}
		}
	})

	t.Run("IDOrdering", func(t *testing.T) {
		pool := atom.NewPool()

		id0 := pool.Intern("first")
		id1 := pool.Intern("second")
		id2 := pool.Intern("third")

		if id0 != 0 || id1 != 1 || id2 != 2 {
			t.Errorf("expected sequential IDs 0,1,2, got %d,%d,%d", id0, id1, id2)
		}
	})

	t.Run("EmptyString", func(t *testing.T) {
		pool := atom.NewPool()

		id := pool.Intern("")
		s, ok := pool.Resolve(id)
		if !ok || s != "" {
			t.Error("expected to handle empty string")
		}
	})
}

func TestPoolConcurrency(t *testing.T) {
	pool := atom.NewPool()
	atoms := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for _, s := range atoms {
				id := pool.Intern(s)
				resolved, ok := pool.Resolve(id)
				if !ok || resolved != s {
					t.Errorf("goroutine %d: failed to resolve %q", n, s)
				}
			}
		}(i)
	}
	wg.Wait()

	if pool.Size() != len(atoms) {
		t.Errorf("expected %d atoms, got %d", len(atoms), pool.Size())
	}
}

func TestGlobalPool(t *testing.T) {
	// Note: Global pool persists across tests, so we just verify it works
	id := atom.Intern("global_test_atom")
	s, ok := atom.Resolve(id)
	if !ok {
		t.Fatal("expected to resolve from global pool")
	}
	if s != "global_test_atom" {
		t.Errorf("expected 'global_test_atom', got %q", s)
	}

	s2 := atom.MustResolve(id)
	if s2 != "global_test_atom" {
		t.Errorf("MustResolve: expected 'global_test_atom', got %q", s2)
	}
}

func BenchmarkInternNew(b *testing.B) {
	pool := atom.NewPool()
	atoms := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		atoms[i] = string(rune('a' + (i % 26)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Intern(atoms[i])
	}
}

func BenchmarkInternExisting(b *testing.B) {
	pool := atom.NewPool()
	pool.Intern("existing")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Intern("existing")
	}
}

func BenchmarkResolve(b *testing.B) {
	pool := atom.NewPool()
	id := pool.Intern("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Resolve(id)
	}
}


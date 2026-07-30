package engine

import (
	"sync"

	"github.com/ioriimasu/jervis/internal/memory/contracts"
)

// Engine implements WorkingMemory with a sliding window (FIFO) policy.
type Engine struct {
	mu       sync.RWMutex
	entries  []contracts.Entry
	index    map[string]int
	capacity int
}

func New(capacity int) *Engine {
	if capacity <= 0 {
		capacity = 50 // Default capacity
	}
	return &Engine{
		entries:  make([]contracts.Entry, 0, capacity),
		index:    make(map[string]int),
		capacity: capacity,
	}
}

func (e *Engine) Add(entry contracts.Entry) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// If entry already exists, remove the old one first to maintain FIFO order for the new one
	if idx, exists := e.index[entry.ID()]; exists {
		e.entries = append(e.entries[:idx], e.entries[idx+1:]...)
		// We'll re-index after adding the new one or just update index now.
		// It's easier to just re-build index if we prune or move.
	} else if len(e.entries) >= e.capacity {
		// Prune the oldest (FIFO)
		oldest := e.entries[0]
		delete(e.index, oldest.ID())
		e.entries = e.entries[1:]
	}

	// Add to end
	e.entries = append(e.entries, entry)

	// Re-build index to ensure correctness
	e.index = make(map[string]int, len(e.entries))
	for i, ent := range e.entries {
		e.index[ent.ID()] = i
	}

	return nil
}

func (e *Engine) Get(id string) (contracts.Entry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	idx, ok := e.index[id]
	if !ok {
		return nil, false
	}
	return e.entries[idx], true
}

func (e *Engine) All() []contracts.Entry {
	e.mu.RLock()
	defer e.mu.RUnlock()

	all := make([]contracts.Entry, len(e.entries))
	copy(all, e.entries)
	return all
}

func (e *Engine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.entries = make([]contracts.Entry, 0, e.capacity)
	e.index = make(map[string]int)
}

func (e *Engine) Capacity() int {
	return e.capacity
}

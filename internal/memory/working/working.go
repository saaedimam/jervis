package working

import (
	"github.com/ioriimasu/jervis/internal/memory/contracts"
	"github.com/ioriimasu/jervis/internal/memory/working/engine"
)

// WorkingMemory provides the high-level facade for active context.
type WorkingMemory struct {
	contracts.WorkingMemory
}

// New constructs a new WorkingMemory instance with the specified capacity.
func New(capacity int) *WorkingMemory {
	return &WorkingMemory{
		WorkingMemory: engine.New(capacity),
	}
}

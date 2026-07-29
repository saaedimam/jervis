package timeline

import (
	"github.com/ioriimasu/jervis/internal/memory/contracts"
	storecontracts "github.com/ioriimasu/jervis/internal/memory/store/contracts"
	"github.com/ioriimasu/jervis/internal/memory/timeline/engine"
)

// Timeline provides the high-level facade for the event ledger.
type Timeline struct {
	contracts.Timeline
}

// New constructs a new Timeline instance with the specified persistent store.
func New(store storecontracts.Store) *Timeline {
	return &Timeline{
		Timeline: engine.New(store),
	}
}

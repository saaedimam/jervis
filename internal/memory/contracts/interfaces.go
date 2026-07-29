package contracts

import (
	"context"
	"time"
	runtimecontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

// Entry represents a single unit of information in working memory.
type Entry interface {
	// ID returns the unique identifier for this entry.
	ID() string
	
	// Content returns the payload of the entry.
	Content() any
	
	// Metadata returns entry-specific key-value pairs.
	Metadata() map[string]string
	
	// Timestamp returns when the entry was created.
	Timestamp() time.Time
}

// WorkingMemory defines the interface for active context storage.
type WorkingMemory interface {
	// Add inserts a new entry into working memory, potentially pruning old ones.
	Add(entry Entry) error
	
	// Get retrieves an entry by its ID.
	Get(id string) (Entry, bool)
	
	// All returns all entries in chronological order.
	All() []Entry
	
	// Clear removes all entries from working memory.
	Clear()
	
	// Capacity returns the maximum number of entries allowed.
	Capacity() int
}

// Timeline defines the interface for the immutable event ledger.
type Timeline interface {
	// Append adds an event to the ledger.
	Append(ctx context.Context, event runtimecontracts.Event) error
	
	// Query retrieves events based on filter criteria.
	Query(ctx context.Context, filter Filter) ([]runtimecontracts.Event, error)
}

// Filter defines criteria for querying the timeline.
type Filter struct {
	From     time.Time
	To       time.Time
	Type     string
	Limit    int
}

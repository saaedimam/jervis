package device

import (
	"sync"
)

const maxCacheEntries = 500

// IdempotencyManager tracks recently processed event and command IDs in memory
// to quickly filter duplicates before querying or writing to SQLite.
type IdempotencyManager struct {
	mu            sync.Mutex
	recentEvents  map[string]int64
	recentResults map[string]int64
}

// NewIdempotencyManager creates a thread-safe idempotency tracker.
func NewIdempotencyManager() *IdempotencyManager {
	return &IdempotencyManager{
		recentEvents:  make(map[string]int64),
		recentResults: make(map[string]int64),
	}
}

// MarkEventSeen records an event ID. Returns true if already seen (duplicate).
func (im *IdempotencyManager) MarkEventSeen(eventID string, timestamp int64) bool {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, exists := im.recentEvents[eventID]; exists {
		return true
	}

	if len(im.recentEvents) >= maxCacheEntries {
		// Evict an arbitrary entry
		for k := range im.recentEvents {
			delete(im.recentEvents, k)
			break
		}
	}

	im.recentEvents[eventID] = timestamp
	return false
}

// MarkResultSeen records a command result. Returns true if already seen (duplicate).
func (im *IdempotencyManager) MarkResultSeen(commandID string, timestamp int64) bool {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, exists := im.recentResults[commandID]; exists {
		return true
	}

	if len(im.recentResults) >= maxCacheEntries {
		for k := range im.recentResults {
			delete(im.recentResults, k)
			break
		}
	}

	im.recentResults[commandID] = timestamp
	return false
}

package contracts

import (
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// Session represents an active user session or workspace context.
type Session interface {
	// ID returns the unique session identifier.
	ID() types.SessionID

	// Metadata returns a defensive copy of session metadata.
	Metadata() map[string]string

	// SetMetadata sets a metadata key-value pair.
	SetMetadata(key, value string)

	// GetMetadata retrieves a metadata value.
	GetMetadata(key string) (string, bool)

	// State returns the current lifecycle state of the session.
	State() types.State
}

// Registry manages the collection of active sessions.
type Registry interface {
	// Register adds a session to the registry.
	Register(session Session) error

	// Unregister removes a session by its ID.
	Unregister(id types.SessionID) error

	// Get retrieves a session by its ID.
	Get(id types.SessionID) (Session, bool)

	// All returns a defensive copy slice of all registered sessions.
	All() []Session

	// Count returns the number of active sessions.
	Count() int

	// Clear removes all sessions from the registry.
	Clear()
}

// Manager is the high-level facade for session operations.
type Manager interface {
	// CreateSession initializes a new session with an optional ID.
	CreateSession(id string) (Session, error)

	// GetSession retrieves an active session.
	GetSession(id types.SessionID) (Session, bool)

	// CloseSession terminates an active session.
	CloseSession(id types.SessionID) error
}

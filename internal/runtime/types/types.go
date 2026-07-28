package types

import (
	"fmt"
	"time"
)

// RuntimeID represents a canonical runtime instance identifier.
type RuntimeID struct {
	val string
}

// NewRuntimeID constructs a validated RuntimeID.
func NewRuntimeID(val string) (RuntimeID, error) {
	if val == "" {
		return RuntimeID{}, fmt.Errorf("runtime ID cannot be empty")
	}
	return RuntimeID{val: val}, nil
}

// String returns the string representation of RuntimeID.
func (id RuntimeID) String() string {
	return id.val
}

// IsZero checks if the RuntimeID is uninitialized.
func (id RuntimeID) IsZero() bool {
	return id.val == ""
}

// SessionID represents a canonical session identifier.
type SessionID struct {
	val string
}

// NewSessionID constructs a validated SessionID.
func NewSessionID(val string) (SessionID, error) {
	if val == "" {
		return SessionID{}, fmt.Errorf("session ID cannot be empty")
	}
	return SessionID{val: val}, nil
}

// String returns the string representation of SessionID.
func (id SessionID) String() string {
	return id.val
}

// IsZero checks if the SessionID is uninitialized.
func (id SessionID) IsZero() bool {
	return id.val == ""
}

// EventID represents a canonical event identifier.
type EventID struct {
	val string
}

// NewEventID constructs a validated EventID.
func NewEventID(val string) (EventID, error) {
	if val == "" {
		return EventID{}, fmt.Errorf("event ID cannot be empty")
	}
	return EventID{val: val}, nil
}

// String returns the string representation of EventID.
func (id EventID) String() string {
	return id.val
}

// IsZero checks if the EventID is uninitialized.
func (id EventID) IsZero() bool {
	return id.val == ""
}

// Timestamp represents a canonical runtime domain timestamp wrapper.
type Timestamp struct {
	t time.Time
}

// NewTimestamp creates a new Timestamp wrapping the provided time.Time in UTC.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{t: t.UTC()}
}

// Now creates a new Timestamp representing the current UTC time.
func Now() Timestamp {
	return Timestamp{t: time.Now().UTC()}
}

// Time returns the underlying time.Time in UTC.
func (ts Timestamp) Time() time.Time {
	return ts.t
}

// UnixNano returns the timestamp in nanoseconds since Unix epoch.
func (ts Timestamp) UnixNano() int64 {
	return ts.t.UnixNano()
}

// String returns the ISO8601/RFC3339 string format of the timestamp.
func (ts Timestamp) String() string {
	if ts.t.IsZero() {
		return ""
	}
	return ts.t.Format(time.RFC3339Nano)
}

// IsZero checks if the Timestamp is uninitialized.
func (ts Timestamp) IsZero() bool {
	return ts.t.IsZero()
}

// State represents runtime lifecycle state.
type State string

const (
	// StateCreated indicates the component has been created.
	StateCreated State = "Created"

	// StateInitializing indicates the component is initializing.
	StateInitializing State = "Initializing"

	// StateRunning indicates the component is actively running.
	StateRunning State = "Running"

	// StateStopping indicates the component is in the process of shutting down.
	StateStopping State = "Stopping"

	// StateStopped indicates the component is cleanly stopped.
	StateStopped State = "Stopped"

	// StateFailed indicates the component has encountered an unrecoverable failure.
	StateFailed State = "Failed"
)

// String returns the string representation of State.
func (s State) String() string {
	return string(s)
}

// IsValid validates whether the State is one of the permitted canonical states.
func (s State) IsValid() bool {
	switch s {
	case StateCreated, StateInitializing, StateRunning, StateStopping, StateStopped, StateFailed:
		return true
	default:
		return false
	}
}

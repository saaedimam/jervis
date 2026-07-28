package events

import (
	"time"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/errors"
)

// Header is a map of string keys to string values for event metadata.
type Header map[string]string

// Priority defines the priority level of an event.
// Lower numbers indicate higher priority.
type Priority int8

const (
	// PriorityLow is the lowest priority.
	PriorityLow Priority = -1
	// PriorityNormal is the default priority.
	PriorityNormal Priority = 0
	// PriorityHigh is the highest priority.
	PriorityHigh Priority = 1
)

// Metadata holds additional metadata for an event.
type Metadata struct {
	Headers Header
	Source  string
	Tags    map[string]string
}

// Envelope wraps an event with additional delivery information.
type Envelope struct {
	Event       Event
	Attempts    int
	MaxAttempts int
}

// EventID is an alias for contracts.EventID to satisfy the requirement.
type EventID = contracts.EventID

// EventType is an alias for string to represent the event type.
type EventType string

// Event represents an immutable event that implements contracts.Event.
type Event struct {
	id            EventID
	eventType     EventType
	source        string
	timestamp     Timestamp
	correlationID string
	causationID   string
	priority      int
	payload       []byte
	metadata      Metadata
	version       string
}

// Timestamp represents a timestamp in nanoseconds since epoch.
// It is an alias for contracts.Timestamp.
type Timestamp = contracts.Timestamp

// NewBuilder returns a new event builder with default values.
func NewBuilder() *Builder {
	return &Builder{
		event: Event{
			priority: 0,
			version:  "1.0.0",
			metadata: Metadata{
				Headers: make(Header),
				Tags:    make(map[string]string),
			},
		},
	}
}

// Builder provides a fluent interface for constructing events.
type Builder struct {
	event Event
}

// WithID sets the event ID.
func (b *Builder) WithID(id EventID) *Builder {
	b.event.id = id
	return b
}

// WithType sets the event type.
func (b *Builder) WithType(t EventType) *Builder {
	b.event.eventType = t
	return b
}

// WithSource sets the event source.
func (b *Builder) WithSource(s string) *Builder {
	b.event.source = s
	return b
}

// WithTimestamp sets the event timestamp.
func (b *Builder) WithTimestamp(ts Timestamp) *Builder {
	b.event.timestamp = ts
	return b
}

// WithCorrelationID sets the correlation ID.
func (b *Builder) WithCorrelationID(cid string) *Builder {
	b.event.correlationID = cid
	return b
}

// WithCausationID sets the causation ID.
func (b *Builder) WithCausationID(cid string) *Builder {
	b.event.causationID = cid
	return b
}

// WithPriority sets the event priority.
func (b *Builder) WithPriority(p Priority) *Builder {
	b.event.priority = int(p)
	return b
}

// WithPayload sets the event payload.
// The payload is copied to ensure immutability.
func (b *Builder) WithPayload(p []byte) *Builder {
	if p == nil {
		b.event.payload = nil
	} else {
		b.event.payload = make([]byte, len(p))
		copy(b.event.payload, p)
	}
	return b
}

// WithMetadata sets the event metadata.
// A deep copy of headers and tags is made.
func (b *Builder) WithMetadata(m Metadata) *Builder {
	headers := make(Header)
	for k, v := range m.Headers {
		headers[k] = v
	}
	tags := make(map[string]string)
	for k, v := range m.Tags {
		tags[k] = v
	}
	b.event.metadata = Metadata{
		Headers: headers,
		Source:  m.Source,
		Tags:    tags,
	}
	return b
}

// WithVersion sets the event version.
func (b *Builder) WithVersion(v string) *Builder {
	b.event.version = v
	return b
}

// Build constructs and returns an immutable event.
// It fills in default values for missing fields.
func (b *Builder) Build() (Event, error) {
	e := b.event
	if e.id == "" {
		// Generate a simple ID based on current time nanoseconds.
		e.id = EventID("evt-" + string(Timestamp(time.Now().UnixNano())))
	}
	if e.timestamp == 0 {
		e.timestamp = Timestamp(time.Now().UnixNano())
	}
	if e.source == "" {
		e.source = "unknown"
	}
	if e.eventType == "" {
		return Event{}, errors.ErrInvalidEvent
	}
	return e, nil
}

// Validate checks if the event is valid.
func (e Event) Validate() error {
	if e.id == "" {
		return errors.ErrInvalidEvent
	}
	if e.eventType == "" {
		return errors.ErrInvalidEvent
	}
	if e.timestamp == 0 {
		return errors.ErrInvalidEvent
	}
	// Additional validation can be added here.
	return nil
}

// Payload returns a copy of the event payload to maintain immutability.
func (e Event) Payload() any {
	if e.payload == nil {
		return nil
	}
	cp := make([]byte, len(e.payload))
	copy(cp, e.payload)
	return cp
}

// Accessor methods to satisfy contracts.Event interface.

func (e Event) ID() EventID {
	return e.id
}

func (e Event) Type() string {
	return string(e.eventType)
}

func (e Event) Source() string {
	return e.source
}

func (e Event) Timestamp() Timestamp {
	return e.timestamp
}

func (e Event) CorrelationID() string {
	return e.correlationID
}

func (e Event) CausationID() string {
	return e.causationID
}

func (e Event) Priority() int {
	return e.priority
}

func (e Event) Metadata() map[string]string {
	// Return a copy of the headers map.
	headers := make(map[string]string, len(e.metadata.Headers))
	for k, v := range e.metadata.Headers {
		headers[k] = v
	}
	return headers
}

func (e Event) Version() string {
	return e.version
}

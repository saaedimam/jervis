package events

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/ioriimasu/jervis/internal/runtime/eventbus/errors"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// Priority represents event processing priority as a byte.
type Priority uint8

const (
	// PriorityLow is for low priority background events.
	PriorityLow Priority = iota
	// PriorityNormal is the standard default priority.
	PriorityNormal
	// PriorityHigh is for high priority urgent events.
	PriorityHigh
	// PriorityCritical is for critical system events.
	PriorityCritical
)

const (
	DefaultPriority = PriorityNormal
	DefaultVersion  = "1.0.0"
)

// EventID represents a canonical event instance identifier.
type EventID = types.EventID

// EventType represents a namespaced event type classification.
type EventType string

func (t EventType) String() string {
	return string(t)
}

// Header holds header metadata for an event envelope.
type Header struct {
	ID            types.EventID
	Type          EventType
	Source        string
	Timestamp     types.Timestamp
	CorrelationID string
	CausationID   string
	Priority      Priority
	Version       string
}

// Metadata represents key-value header metadata attributes.
type Metadata map[string]string

// Clone returns a deep defensive copy of the Metadata map.
func (m Metadata) Clone() Metadata {
	if m == nil {
		return make(Metadata)
	}
	cp := make(Metadata, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// Envelope represents an immutable event container implementing contracts.Event.
type Envelope struct {
	header   Header
	payload  any
	metadata Metadata
}

var _ contracts.Event = (*Envelope)(nil)

// ID returns the unique EventID.
func (e *Envelope) ID() types.EventID {
	return e.header.ID
}

// Type returns the namespaced event type string.
func (e *Envelope) Type() string {
	return e.header.Type.String()
}

// Source returns the originating subsystem or component path.
func (e *Envelope) Source() string {
	return e.header.Source
}

// Timestamp returns the UTC creation timestamp.
func (e *Envelope) Timestamp() types.Timestamp {
	return e.header.Timestamp
}

// CorrelationID returns the root workflow correlation identifier.
func (e *Envelope) CorrelationID() string {
	return e.header.CorrelationID
}

// CausationID returns the direct parent event identifier.
func (e *Envelope) CausationID() string {
	return e.header.CausationID
}

// Priority returns the handler execution priority byte.
func (e *Envelope) Priority() uint8 {
	return uint8(e.header.Priority)
}

// Payload returns the event payload object.
func (e *Envelope) Payload() any {
	return e.payload
}

// Metadata returns a defensive copy of key-value header metadata.
func (e *Envelope) Metadata() map[string]string {
	return e.metadata.Clone()
}

// Version returns the event schema semantic version.
func (e *Envelope) Version() string {
	return e.header.Version
}

// Header returns a copy of the event Header struct.
func (e *Envelope) Header() Header {
	return e.header
}

// Clone creates and returns a deep defensive copy of the Envelope.
func (e *Envelope) Clone() *Envelope {
	if e == nil {
		return nil
	}
	return &Envelope{
		header:   e.header,
		payload:  e.payload,
		metadata: e.metadata.Clone(),
	}
}

// ValidatePriority checks if priority falls within permitted bounds [PriorityLow, PriorityCritical].
func ValidatePriority(p Priority) error {
	if p > PriorityCritical {
		return fmt.Errorf("%w: priority %d must be between %d and %d", errs.ErrInvalidPriority, p, PriorityLow, PriorityCritical)
	}
	return nil
}

// ValidateEventType verifies that event type complies with lowercase dot-separated format (<layer>.<component>.<verb>).
func ValidateEventType(t EventType) error {
	str := t.String()
	if str == "" {
		return fmt.Errorf("%w: event type cannot be empty", errs.ErrValidationFailed)
	}
	for _, r := range str {
		if unicode.IsUpper(r) || unicode.IsSpace(r) {
			return fmt.Errorf("%w: event type %q must be lowercase and contain no spaces", errs.ErrValidationFailed, str)
		}
	}
	parts := strings.Split(str, ".")
	if len(parts) < 2 {
		return fmt.Errorf("%w: event type %q must contain at least one namespace separator '.'", errs.ErrValidationFailed, str)
	}
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("%w: event type %q contains empty namespace segment", errs.ErrValidationFailed, str)
		}
	}
	return nil
}

// ValidateEvent performs strict validation on an event envelope.
func ValidateEvent(e contracts.Event) error {
	if e == nil {
		return fmt.Errorf("%w: event cannot be nil", errs.ErrInvalidEvent)
	}
	if e.ID().IsZero() {
		return fmt.Errorf("%w: event ID cannot be empty", errs.ErrValidationFailed)
	}
	if err := ValidateEventType(EventType(e.Type())); err != nil {
		return err
	}
	if e.Source() == "" {
		return fmt.Errorf("%w: event source cannot be empty", errs.ErrValidationFailed)
	}
	if e.Timestamp().IsZero() {
		return fmt.Errorf("%w: event timestamp cannot be zero", errs.ErrValidationFailed)
	}
	if err := ValidatePriority(Priority(e.Priority())); err != nil {
		return err
	}
	if e.Payload() == nil {
		return fmt.Errorf("%w: event payload cannot be nil", errs.ErrValidationFailed)
	}
	if e.Version() == "" {
		return fmt.Errorf("%w: event version cannot be empty", errs.ErrValidationFailed)
	}
	return nil
}

// Builder constructs an immutable Envelope.
type Builder struct {
	header   Header
	payload  any
	metadata Metadata
}

// NewBuilder constructs an initialized Builder with default priority and version.
func NewBuilder() *Builder {
	return &Builder{
		header: Header{
			Priority: DefaultPriority,
			Version:  DefaultVersion,
		},
		metadata: make(Metadata),
	}
}

// SetID sets the event instance ID.
func (b *Builder) SetID(id types.EventID) *Builder {
	b.header.ID = id
	return b
}

// SetType sets the event type string.
func (b *Builder) SetType(t EventType) *Builder {
	b.header.Type = t
	return b
}

// SetSource sets the originating component path.
func (b *Builder) SetSource(source string) *Builder {
	b.header.Source = source
	return b
}

// SetTimestamp sets the event creation timestamp.
func (b *Builder) SetTimestamp(ts types.Timestamp) *Builder {
	b.header.Timestamp = ts
	return b
}

// SetCorrelationID sets the workflow correlation ID.
func (b *Builder) SetCorrelationID(id string) *Builder {
	b.header.CorrelationID = id
	return b
}

// SetCausationID sets the direct parent event causation ID.
func (b *Builder) SetCausationID(id string) *Builder {
	b.header.CausationID = id
	return b
}

// SetPriority sets the event processing priority.
func (b *Builder) SetPriority(p Priority) *Builder {
	b.header.Priority = p
	return b
}

// SetPayload sets the event payload object (must not be nil).
func (b *Builder) SetPayload(payload any) *Builder {
	b.payload = payload
	return b
}

// SetMetadata sets or updates a key-value header metadata pair.
func (b *Builder) SetMetadata(key, value string) *Builder {
	if b.metadata == nil {
		b.metadata = make(Metadata)
	}
	b.metadata[key] = value
	return b
}

// SetVersion sets the semantic version string.
func (b *Builder) SetVersion(ver string) *Builder {
	b.header.Version = ver
	return b
}

// Build validates and constructs an immutable *Envelope instance.
func (b *Builder) Build() (*Envelope, error) {
	if b.header.Timestamp.IsZero() {
		b.header.Timestamp = types.Now()
	}
	if b.header.CorrelationID == "" && !b.header.ID.IsZero() {
		b.header.CorrelationID = b.header.ID.String()
	}
	if b.header.CausationID == "" && !b.header.ID.IsZero() {
		b.header.CausationID = b.header.ID.String()
	}

	env := &Envelope{
		header:   b.header,
		payload:  b.payload,
		metadata: b.metadata.Clone(),
	}

	if err := ValidateEvent(env); err != nil {
		return nil, err
	}
	return env, nil
}

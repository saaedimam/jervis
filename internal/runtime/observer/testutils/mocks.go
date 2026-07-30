package testutils

import (
	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/observer/contracts"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

var (
	_ eventcontracts.Event = (*MockEvent)(nil)
	_ contracts.Observer   = (*MockObserver)(nil)
)

type MockEvent struct {
	IDVal            types.EventID
	TypeVal          string
	SourceVal        string
	TimestampVal     types.Timestamp
	CorrelationIDVal string
	CausationIDVal   string
	PriorityVal      uint8
	PayloadVal       any
	MetadataVal      map[string]string
	VersionVal       string
}

func (m *MockEvent) ID() types.EventID           { return m.IDVal }
func (m *MockEvent) Type() string                { return m.TypeVal }
func (m *MockEvent) Source() string              { return m.SourceVal }
func (m *MockEvent) Timestamp() types.Timestamp  { return m.TimestampVal }
func (m *MockEvent) CorrelationID() string       { return m.CorrelationIDVal }
func (m *MockEvent) CausationID() string         { return m.CausationIDVal }
func (m *MockEvent) Priority() uint8             { return m.PriorityVal }
func (m *MockEvent) Payload() any                { return m.PayloadVal }
func (m *MockEvent) Metadata() map[string]string { return m.MetadataVal }
func (m *MockEvent) Version() string             { return m.VersionVal }

type MockObserver struct {
	IDVal      string
	HandleFunc func(n contracts.Notification) error
}

func (m *MockObserver) ID() string { return m.IDVal }
func (m *MockObserver) Handle(n contracts.Notification) error {
	if m.HandleFunc != nil {
		return m.HandleFunc(n)
	}
	return nil
}

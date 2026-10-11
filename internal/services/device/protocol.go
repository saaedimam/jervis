package device

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	// SchemaVersion is the immutable wire protocol schema version.
	SchemaVersion = "phone-agent.v1"

	// Envelope types
	TypeEvent   = "event"
	TypeCommand = "command"
	TypeResult  = "result"
	TypeError   = "error"

	// Event types emitted by Android agent
	EventSMSReceived     = "sms.received"
	EventBatteryChanged  = "battery.changed"
	EventDeviceBoot      = "device.boot"
	EventDevicePing      = "device.ping"

	// Allowed Capability Actions (Frozen Allowlist)
	ActionSendSMS       = "send_sms"
	ActionCaptureAudio  = "capture_audio"
	ActionNotify        = "notify"
	ActionBatteryStatus = "battery_status"
	ActionClipboardGet  = "clipboard_get"
	ActionClipboardSet  = "clipboard_set"
	ActionDeviceStatus  = "device_status"

	// Result status
	StatusSuccess = "success"
	StatusError   = "error"
)

var (
	ErrInvalidSchema    = errors.New("device: invalid or unsupported schema version")
	ErrMissingID        = errors.New("device: missing id field")
	ErrMissingDevice    = errors.New("device: missing device field")
	ErrInvalidTimestamp = errors.New("device: invalid timestamp (must be positive)")
	ErrMissingType      = errors.New("device: missing type field")
	ErrMissingPayload   = errors.New("device: missing payload field")
	ErrDisallowedAction = errors.New("device: action not in capability allowlist")
	ErrMissingCommandID = errors.New("device: missing command_id in payload")
	ErrMissingAction    = errors.New("device: missing action in payload")
	ErrUnknownField     = errors.New("device: unknown field in payload")
)

// AllowedActions contains the capability allowlist for commands.
var AllowedActions = map[string]struct{}{
	ActionSendSMS:       {},
	ActionCaptureAudio:  {},
	ActionNotify:        {},
	ActionBatteryStatus: {},
	ActionClipboardGet:  {},
	ActionClipboardSet:  {},
	ActionDeviceStatus:  {},
}

// IsAllowedAction checks if an action is present in the capability allowlist.
func IsAllowedAction(action string) bool {
	_, ok := AllowedActions[action]
	return ok
}

// Envelope represents the canonical phone-agent.v1 wire envelope.
type Envelope struct {
	Schema    string          `json:"schema"`
	ID        string          `json:"id"`
	Device    string          `json:"device"`
	Timestamp int64           `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

// ValidateEnvelope verifies envelope integrity against phone-agent.v1 invariants.
func (e *Envelope) ValidateEnvelope() error {
	if e.Schema != SchemaVersion {
		return fmt.Errorf("%w: got %q, expected %q", ErrInvalidSchema, e.Schema, SchemaVersion)
	}
	if strings.TrimSpace(e.ID) == "" {
		return ErrMissingID
	}
	if strings.TrimSpace(e.Device) == "" {
		return ErrMissingDevice
	}
	if e.Timestamp <= 0 {
		return ErrInvalidTimestamp
	}
	if strings.TrimSpace(e.Type) == "" {
		return ErrMissingType
	}
	if len(e.Payload) == 0 || string(e.Payload) == "null" {
		return ErrMissingPayload
	}
	return nil
}

// CommandPayload defines the strict typed structure for commands dispatched to the device.
type CommandPayload struct {
	CommandID       string `json:"command_id"`
	Action          string `json:"action"`
	Recipient       string `json:"recipient,omitempty"`
	Text            string `json:"text,omitempty"`
	Title           string `json:"title,omitempty"`
	Message         string `json:"message,omitempty"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
}

// Validate verifies command action against allowlist and required parameters.
func (cp *CommandPayload) Validate() error {
	if strings.TrimSpace(cp.CommandID) == "" {
		return ErrMissingCommandID
	}
	if strings.TrimSpace(cp.Action) == "" {
		return ErrMissingAction
	}
	if !IsAllowedAction(cp.Action) {
		return fmt.Errorf("%w: %q", ErrDisallowedAction, cp.Action)
	}
	switch cp.Action {
	case ActionSendSMS:
		if strings.TrimSpace(cp.Recipient) == "" {
			return errors.New("device: send_sms requires non-empty recipient")
		}
		if strings.TrimSpace(cp.Text) == "" {
			return errors.New("device: send_sms requires non-empty text")
		}
	case ActionNotify:
		if strings.TrimSpace(cp.Message) == "" {
			return errors.New("device: notify requires non-empty message")
		}
	case ActionCaptureAudio:
		if cp.DurationSeconds < 1 || cp.DurationSeconds > 10 {
			return errors.New("device: capture_audio duration_seconds must be between 1 and 10")
		}
	}
	return nil
}

// ResultPayload defines the strict typed structure for execution results returned from the device.
type ResultPayload struct {
	CommandID string          `json:"command_id"`
	Action    string          `json:"action"`
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     *string         `json:"error"`
}

// Validate verifies result payload integrity.
func (rp *ResultPayload) Validate() error {
	if strings.TrimSpace(rp.CommandID) == "" {
		return ErrMissingCommandID
	}
	if strings.TrimSpace(rp.Action) == "" {
		return ErrMissingAction
	}
	if !IsAllowedAction(rp.Action) {
		return fmt.Errorf("%w: %q", ErrDisallowedAction, rp.Action)
	}
	if rp.Status != StatusSuccess && rp.Status != StatusError {
		return fmt.Errorf("device: invalid result status %q (must be success or error)", rp.Status)
	}
	return nil
}

// SMSReceivedPayload represents inbound SMS telemetry.
type SMSReceivedPayload struct {
	Sender       string `json:"sender"`
	Body         string `json:"body"`
	SMSTimestamp int64  `json:"sms_timestamp"`
}

// BatteryChangedPayload represents device battery events.
type BatteryChangedPayload struct {
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

// DeviceBootPayload represents device boot events.
type DeviceBootPayload struct {
	BootAction string `json:"boot_action"`
	BootTime   int64  `json:"boot_time"`
}

// DevicePingPayload represents device heartbeat pings.
type DevicePingPayload struct {
	Ping bool  `json:"ping"`
	Time int64 `json:"time"`
}

// DecodeStrict decodes JSON bytes into target struct with DisallowUnknownFields.
func DecodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("%w: %v", ErrUnknownField, err)
	}
	return nil
}

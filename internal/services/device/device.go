package device

import (
	"context"
	"time"
)

// DeviceStatus exposes server-side known device telemetry without leaking sensitive contents.
type DeviceStatus struct {
	DeviceID    string `json:"device_id"`
	Connected   bool   `json:"connected"`
	LastEvent   string `json:"last_event,omitempty"`
	LastCommand string `json:"last_command,omitempty"`
	LastResult  string `json:"last_result,omitempty"`
	QueueDepth  int    `json:"queue_depth"`
	LastSeen    int64  `json:"last_seen,omitempty"`
}

// Service defines the core interface for the device control plane.
type Service interface {
	// IngestEvent processes an incoming event from the phone.
	// Returns isDuplicate=true if the event was already persisted.
	IngestEvent(ctx context.Context, env *Envelope) (isDuplicate bool, err error)

	// QueueCommand validates and enqueues an approved command for delivery to the phone.
	QueueCommand(ctx context.Context, cmd *Envelope) error

	// WaitForCommand retrieves the next pending command for the phone (bounded long-poll).
	WaitForCommand(ctx context.Context, deviceID string, timeout time.Duration) (*Envelope, error)

	// IngestResult processes a capability execution result reported by the phone.
	// Returns isDuplicate=true if the result was already received.
	IngestResult(ctx context.Context, res *Envelope) (isDuplicate bool, err error)

	// GetDeviceStatus returns authenticated server-side status for a device.
	GetDeviceStatus(ctx context.Context, deviceID string) (*DeviceStatus, error)
}

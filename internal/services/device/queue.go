package device

import (
	"context"
	"errors"
	"sync"
	"time"
)

const MaxQueueSize = 100

var ErrQueueFull = errors.New("device: command queue reached maximum capacity (100)")

// CommandQueue manages bounded, per-device FIFO command queues with long-poll wait capabilities.
type CommandQueue struct {
	mu      sync.Mutex
	queues  map[string][]*Envelope
	waiters map[string][]chan struct{}
}

// NewCommandQueue initializes a new thread-safe bounded command queue manager.
func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		queues:  make(map[string][]*Envelope),
		waiters: make(map[string][]chan struct{}),
	}
}

// Enqueue adds a command envelope to the device's FIFO queue.
// Returns ErrQueueFull if the queue has reached MaxQueueSize (100).
func (q *CommandQueue) Enqueue(cmd *Envelope) error {
	if cmd == nil {
		return errors.New("device: cannot enqueue nil command")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	deviceID := cmd.Device
	items := q.queues[deviceID]
	if len(items) >= MaxQueueSize {
		return ErrQueueFull
	}

	q.queues[deviceID] = append(items, cmd)

	// Notify waiting long-poll listeners for this device
	if waiterList, exists := q.waiters[deviceID]; exists && len(waiterList) > 0 {
		// Wake the first waiter
		ch := waiterList[0]
		q.waiters[deviceID] = waiterList[1:]
		close(ch)
	}

	return nil
}

// Dequeue pops the next command for the device immediately without waiting.
func (q *CommandQueue) Dequeue(deviceID string) *Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()

	items := q.queues[deviceID]
	if len(items) == 0 {
		return nil
	}

	cmd := items[0]
	q.queues[deviceID] = items[1:]
	return cmd
}

// Depth returns the current number of queued commands for the device.
func (q *CommandQueue) Depth(deviceID string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queues[deviceID])
}

// WaitForCommand implements long-polling:
// If a command exists, it is returned immediately.
// Otherwise, it waits up to timeout or until ctx is cancelled.
// Returns (nil, nil) on normal timeout expiration (e.g. 25s without commands).
func (q *CommandQueue) WaitForCommand(ctx context.Context, deviceID string, timeout time.Duration) (*Envelope, error) {
	q.mu.Lock()

	// 1. Check if a command is already queued
	items := q.queues[deviceID]
	if len(items) > 0 {
		cmd := items[0]
		q.queues[deviceID] = items[1:]
		q.mu.Unlock()
		return cmd, nil
	}

	// 2. Register a one-shot notification channel
	notifyCh := make(chan struct{})
	q.waiters[deviceID] = append(q.waiters[deviceID], notifyCh)
	q.mu.Unlock()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-notifyCh:
		// Notified: pop the enqueued command
		q.mu.Lock()
		defer q.mu.Unlock()
		items := q.queues[deviceID]
		if len(items) > 0 {
			cmd := items[0]
			q.queues[deviceID] = items[1:]
			return cmd, nil
		}
		return nil, nil

	case <-timer.C:
		// Timeout reached; remove waiter from slice if not yet notified
		q.removeWaiter(deviceID, notifyCh)
		return nil, nil

	case <-ctx.Done():
		// Context cancelled (e.g. client disconnected or server shutting down)
		q.removeWaiter(deviceID, notifyCh)
		return nil, ctx.Err()
	}
}

func (q *CommandQueue) removeWaiter(deviceID string, ch chan struct{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	list := q.waiters[deviceID]
	for i, waiter := range list {
		if waiter == ch {
			q.waiters[deviceID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

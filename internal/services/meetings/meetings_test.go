package meetings_test

import (
	"context"
	"testing"
	"time"

	"github.com/saaedimam/jervis/internal/memory/store/sqlite"
	"github.com/saaedimam/jervis/internal/services/meetings"
)

func setupTestMeetings(t *testing.T) (meetings.Service, func()) {
	ctx := context.Background()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}

	if err := meetings.Initialize(ctx, store); err != nil {
		t.Fatalf("failed to initialize meetings schema: %v", err)
	}

	svc := meetings.New(store, nil)
	return svc, func() { _ = store.Close() }
}

func TestMeetingsService(t *testing.T) {
	svc, cleanup := setupTestMeetings(t)
	defer cleanup()

	ctx := context.Background()

	now := time.Now().UTC()
	start := now.Add(1 * time.Hour)
	end := now.Add(2 * time.Hour)

	// Test meeting creation
	m1, err := svc.CreateMeeting(ctx, "meet-1", "Design Review", "Review Phase 3 architecture", start, end, "Zoom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m1.Title != "Design Review" {
		t.Errorf("expected title 'Design Review', got '%s'", m1.Title)
	}

	// Test GetMeeting
	fetched, err := svc.GetMeeting(ctx, "meet-1")
	if err != nil {
		t.Fatalf("unexpected error getting meeting: %v", err)
	}
	if fetched.Location != "Zoom" {
		t.Errorf("expected location 'Zoom', got '%s'", fetched.Location)
	}

	// Test update status
	updated, err := svc.UpdateMeetingStatus(ctx, "meet-1", meetings.StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
	if updated.Status != meetings.StatusCompleted {
		t.Errorf("expected StatusCompleted, got %v", updated.Status)
	}

	// Test list meetings
	list, err := svc.ListMeetings(ctx, now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error listing meetings: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 meeting, got %d", len(list))
	}
}

package habits_test

import (
	"context"
	"testing"
	"time"

	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
	"github.com/ioriimasu/jervis/internal/services/habits"
)

func setupTestHabits(t *testing.T) (habits.Service, func()) {
	ctx := context.Background()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}

	if err := habits.Initialize(ctx, store); err != nil {
		t.Fatalf("failed to initialize habits schema: %v", err)
	}

	svc := habits.New(store, nil)
	return svc, func() { store.Close() }
}

func TestHabitsService(t *testing.T) {
	svc, cleanup := setupTestHabits(t)
	defer cleanup()

	ctx := context.Background()

	// Test habit creation
	h1, err := svc.CreateHabit(ctx, "habit-1", "Drink Water", "Drink 2L of water daily", habits.FrequencyDaily)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h1.Name != "Drink Water" {
		t.Errorf("expected name 'Drink Water', got '%s'", h1.Name)
	}

	// Test habit logging
	today := time.Now().UTC()
	err = svc.LogHabit(ctx, "habit-1", today, true)
	if err != nil {
		t.Fatalf("unexpected error logging habit: %v", err)
	}

	// Test retrieving logs
	logs, err := svc.GetHabitLogs(ctx, "habit-1", today.AddDate(0, 0, -1), today.AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("unexpected error getting logs: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
	if !logs[0].Completed {
		t.Error("expected habit to be logged as completed")
	}

	// Test updating status
	updated, err := svc.UpdateHabitStatus(ctx, "habit-1", habits.StatusArchived)
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
	if updated.Status != habits.StatusArchived {
		t.Errorf("expected StatusArchived, got %v", updated.Status)
	}

	// Test list habits
	_, _ = svc.CreateHabit(ctx, "habit-2", "Exercise", "30 mins exercise", habits.FrequencyDaily)
	list, err := svc.ListHabits(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing habits: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 habits, got %d", len(list))
	}
}

package model

import (
	"context"
	"testing"
	"time"
)

func TestIntervalSchedule(t *testing.T) {
	s, err := NewIntervalSchedule(1 * time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := NewIntervalSchedule(0); err == nil {
		t.Error("expected error for zero interval")
	}

	now := time.Now()
	next := s.NextRun(now)
	if !next.Equal(now.Add(1 * time.Minute)) {
		t.Errorf("expected %v, got %v", now.Add(1*time.Minute), next)
	}
}

func TestOnceSchedule(t *testing.T) {
	target := time.Now().Add(1 * time.Hour)
	s := NewOnceSchedule(target)

	if !s.NextRun(time.Now()).Equal(target) {
		t.Error("expected target time")
	}

	if !s.NextRun(target.Add(1 * time.Second)).IsZero() {
		t.Error("expected zero time after target has passed")
	}
}

func TestCronSchedule_Validation(t *testing.T) {
	tests := []struct {
		expr  string
		valid bool
	}{
		{"* * * * *", true},
		{"0 12 * * *", true},
		{"* * * *", false},
		{"60 * * * *", false},
		{"* 24 * * *", false},
		{"* * 0 * *", false},
		{"* * 32 * *", false},
		{"* * * 0 *", false},
		{"* * * 13 *", false},
		{"* * * * 7", false},
		{"A * * * *", false},
	}

	for _, tt := range tests {
		_, err := NewCronSchedule(tt.expr)
		if (err == nil) != tt.valid {
			t.Errorf("expr %q: expected valid=%v, got err=%v", tt.expr, tt.valid, err)
		}
	}
}

func TestJob_Model(t *testing.T) {
	j := NewJob("id", "name", nil, func(ctx context.Context) error { return nil })
	if j.ID() != "id" || j.Name() != "name" {
		t.Error("mismatched id/name")
	}
	if err := j.Handle(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestCronSchedule_NextRun(t *testing.T) {
	// 0 12 * * * -> Every day at 12:00
	s, _ := NewCronSchedule("0 12 * * *")

	ref := time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC)
	next := s.NextRun(ref)
	expected := time.Date(2026, 7, 29, 12, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}

	// If current time is 12:00, next run should be tomorrow 12:00
	ref2 := time.Date(2026, 7, 29, 12, 0, 0, 0, time.UTC)
	next2 := s.NextRun(ref2)
	expected2 := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	if !next2.Equal(expected2) {
		t.Errorf("expected %v, got %v", expected2, next2)
	}
}

func TestCronSchedule_ComplexMatches(t *testing.T) {
	// Specific day of month and month
	s, _ := NewCronSchedule("0 0 1 1 *") // Jan 1st at 00:00

	ref := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if !s.NextRun(ref).Equal(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Error("expected next year Jan 1st")
	}
}

func TestCronSchedule_DOW(t *testing.T) {
	// Every Sunday at 00:00 (dow=0)
	s, _ := NewCronSchedule("0 0 * * 0")

	ref := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC) // 2026-07-26 is Sunday
	next := s.NextRun(ref)
	expected := time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC) // Next Sunday
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

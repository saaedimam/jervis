package habits

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	storecontracts "github.com/saaedimam/jervis/internal/memory/store/contracts"
	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

var (
	ErrInvalidHabit   = errors.New("habits: invalid habit")
	ErrHabitNotFound  = errors.New("habits: habit not found")
	ErrDuplicateHabit = errors.New("habits: duplicate habit ID")
)

type HabitFrequency string

const (
	FrequencyDaily  HabitFrequency = "DAILY"
	FrequencyWeekly HabitFrequency = "WEEKLY"
)

type HabitStatus string

const (
	StatusActive   HabitStatus = "ACTIVE"
	StatusArchived HabitStatus = "ARCHIVED"
)

type Habit struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Frequency   HabitFrequency `json:"frequency"`
	Status      HabitStatus    `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type HabitLog struct {
	HabitID   string    `json:"habit_id"`
	Date      time.Time `json:"date"`
	Completed bool      `json:"completed"`
}

type Service interface {
	CreateHabit(ctx context.Context, id, name, description string, frequency HabitFrequency) (*Habit, error)
	GetHabit(ctx context.Context, id string) (*Habit, error)
	UpdateHabitStatus(ctx context.Context, id string, status HabitStatus) (*Habit, error)
	ListHabits(ctx context.Context) ([]*Habit, error)
	LogHabit(ctx context.Context, habitID string, date time.Time, completed bool) error
	GetHabitLogs(ctx context.Context, habitID string, from, to time.Time) ([]*HabitLog, error)
}

type service struct {
	store     storecontracts.Store
	publisher eventcontracts.Publisher
}

func New(store storecontracts.Store, publisher eventcontracts.Publisher) Service {
	return &service{
		store:     store,
		publisher: publisher,
	}
}

func (s *service) CreateHabit(ctx context.Context, id, name, description string, frequency HabitFrequency) (*Habit, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return nil, ErrInvalidHabit
	}

	now := time.Now().UTC()
	h := &Habit{
		ID:          id,
		Name:        name,
		Description: description,
		Frequency:   frequency,
		Status:      StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := s.store.Exec(ctx,
		"INSERT INTO habits (id, name, description, frequency, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		h.ID, h.Name, h.Description, h.Frequency, h.Status, h.CreatedAt, h.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert habit: %w", err)
	}

	s.publishEvent(ctx, "habits.habit.created", h)

	return h, nil
}

func (s *service) GetHabit(ctx context.Context, id string) (*Habit, error) {
	row := s.store.QueryRow(ctx,
		"SELECT id, name, description, frequency, status, created_at, updated_at FROM habits WHERE id = ?",
		id,
	)

	var h Habit
	err := row.Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.Status, &h.CreatedAt, &h.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrHabitNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan habit: %w", err)
	}

	return &h, nil
}

func (s *service) UpdateHabitStatus(ctx context.Context, id string, status HabitStatus) (*Habit, error) {
	now := time.Now().UTC()

	result, err := s.store.Exec(ctx,
		"UPDATE habits SET status = ?, updated_at = ? WHERE id = ?",
		status, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrHabitNotFound
	}

	h, err := s.GetHabit(ctx, id)
	if err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "habits.habit.updated", h)

	return h, nil
}

func (s *service) ListHabits(ctx context.Context) ([]*Habit, error) {
	rows, err := s.store.Query(ctx, "SELECT id, name, description, frequency, status, created_at, updated_at FROM habits ORDER BY created_at ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query habits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []*Habit
	for rows.Next() {
		var h Habit
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.Status, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan habit row: %w", err)
		}
		results = append(results, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habit rows: %w", err)
	}

	return results, nil
}

func (s *service) LogHabit(ctx context.Context, habitID string, date time.Time, completed bool) error {
	// Normalize date to YYYY-MM-DD
	dateStr := date.Format("2006-01-02")

	_, err := s.store.Exec(ctx,
		"INSERT INTO habit_logs (habit_id, logged_date, completed) VALUES (?, ?, ?) ON CONFLICT(habit_id, logged_date) DO UPDATE SET completed = excluded.completed",
		habitID, dateStr, completed,
	)
	if err != nil {
		return fmt.Errorf("failed to log habit: %w", err)
	}

	s.publishEvent(ctx, "habits.log.created", map[string]any{
		"habit_id":  habitID,
		"date":      dateStr,
		"completed": completed,
	})

	return nil
}

func (s *service) GetHabitLogs(ctx context.Context, habitID string, from, to time.Time) ([]*HabitLog, error) {
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	rows, err := s.store.Query(ctx,
		"SELECT habit_id, logged_date, completed FROM habit_logs WHERE habit_id = ? AND logged_date >= ? AND logged_date <= ? ORDER BY logged_date ASC",
		habitID, fromStr, toStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query habit logs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []*HabitLog
	for rows.Next() {
		var log HabitLog
		var dateStr string
		if err := rows.Scan(&log.HabitID, &dateStr, &log.Completed); err != nil {
			return nil, fmt.Errorf("failed to scan habit log row: %w", err)
		}
		log.Date, _ = time.Parse("2006-01-02", dateStr)
		results = append(results, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habit log rows: %w", err)
	}

	return results, nil
}

func (s *service) publishEvent(ctx context.Context, eventType string, payload any) {
	if s.publisher == nil {
		return
	}

	id, _ := types.NewEventID(fmt.Sprintf("%s-%d", eventType, time.Now().UnixNano()))
	event, err := events.NewBuilder().
		SetID(id).
		SetType(events.EventType(eventType)).
		SetSource("habits").
		SetPayload(payload).
		Build()

	if err == nil {
		_ = s.publisher.Publish(event)
	}
}

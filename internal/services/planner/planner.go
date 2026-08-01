package planner

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
	ErrInvalidTask   = errors.New("planner: invalid task")
	ErrTaskNotFound  = errors.New("planner: task not found")
	ErrDuplicateTask = errors.New("planner: duplicate task ID")
)

// TaskStatus represents the lifecycle state of a planned task.
type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

// Task represents a unit of planned work.
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Service defines the Planner domain service interface.
type Service interface {
	CreateTask(ctx context.Context, id, title, description string) (*Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	UpdateTaskStatus(ctx context.Context, id string, status TaskStatus) (*Task, error)
	ListTasks(ctx context.Context) ([]*Task, error)
}

type service struct {
	store     storecontracts.Store
	publisher eventcontracts.Publisher
}

// New constructs a new Planner domain service with persistence.
func New(store storecontracts.Store, publisher eventcontracts.Publisher) Service {
	return &service{
		store:     store,
		publisher: publisher,
	}
}

func (s *service) CreateTask(ctx context.Context, id, title, description string) (*Task, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		return nil, ErrInvalidTask
	}

	now := time.Now().UTC()
	t := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := s.store.Exec(ctx,
		"INSERT INTO planner_tasks (id, title, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		t.ID, t.Title, t.Description, t.Status, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		// Basic check for duplicate key; ideally we'd use a more driver-specific check if needed,
		// but for now we'll just return a general error or check if it's a conflict.
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	s.publishEvent(ctx, "planner.task.created", t)

	return t, nil
}

func (s *service) GetTask(ctx context.Context, id string) (*Task, error) {
	row := s.store.QueryRow(ctx,
		"SELECT id, title, description, status, created_at, updated_at FROM planner_tasks WHERE id = ?",
		id,
	)

	var t Task
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	return &t, nil
}

func (s *service) UpdateTaskStatus(ctx context.Context, id string, status TaskStatus) (*Task, error) {
	now := time.Now().UTC()

	result, err := s.store.Exec(ctx,
		"UPDATE planner_tasks SET status = ?, updated_at = ? WHERE id = ?",
		status, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrTaskNotFound
	}

	t, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "planner.task.updated", t)

	return t, nil
}

func (s *service) ListTasks(ctx context.Context) ([]*Task, error) {
	rows, err := s.store.Query(ctx, "SELECT id, title, description, status, created_at, updated_at FROM planner_tasks ORDER BY created_at ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		results = append(results, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}

	return results, nil
}

func (s *service) publishEvent(ctx context.Context, eventType string, task *Task) {
	if s.publisher == nil {
		return
	}

	id, _ := types.NewEventID(fmt.Sprintf("%s-%d", task.ID, time.Now().UnixNano()))
	event, err := events.NewBuilder().
		SetID(id).
		SetType(events.EventType(eventType)).
		SetSource("planner").
		SetPayload(task).
		Build()

	if err == nil {
		_ = s.publisher.Publish(event)
	}
}

package planner

import (
	"errors"
	"strings"
	"sync"
	"time"

	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

var (
	ErrInvalidTask  = errors.New("planner: invalid task")
	ErrTaskNotFound = errors.New("planner: task not found")
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
	CreateTask(id, title, description string) (*Task, error)
	GetTask(id string) (*Task, error)
	UpdateTaskStatus(id string, status TaskStatus) (*Task, error)
	ListTasks() []*Task
}

type service struct {
	mu        sync.RWMutex
	tasks     map[string]*Task
	order     []string
	publisher eventcontracts.Publisher
}

// New constructs a new Planner domain service.
func New(publisher eventcontracts.Publisher) Service {
	return &service{
		tasks:     make(map[string]*Task),
		order:     make([]string, 0),
		publisher: publisher,
	}
}

func (s *service) CreateTask(id, title, description string) (*Task, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		return nil, ErrInvalidTask
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; exists {
		return nil, ErrDuplicateTask
	}

	now := time.Now()
	t := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.tasks[id] = t
	s.order = append(s.order, id)

	return t, nil
}

func (s *service) GetTask(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	// Return defensive copy
	cp := *t
	return &cp, nil
}

func (s *service) UpdateTaskStatus(id string, status TaskStatus) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	t.Status = status
	t.UpdatedAt = time.Now()

	cp := *t
	return &cp, nil
}

func (s *service) ListTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Task, 0, len(s.order))
	for _, id := range s.order {
		if t, ok := s.tasks[id]; ok {
			cp := *t
			result = append(result, &cp)
		}
	}
	return result
}

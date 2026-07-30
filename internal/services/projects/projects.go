package projects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	storecontracts "github.com/ioriimasu/jervis/internal/memory/store/contracts"
	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

var (
	ErrInvalidProject   = errors.New("projects: invalid project")
	ErrProjectNotFound  = errors.New("projects: project not found")
	ErrDuplicateProject = errors.New("projects: duplicate project ID")
)

type ProjectStatus string

const (
	StatusActive    ProjectStatus = "ACTIVE"
	StatusCompleted ProjectStatus = "COMPLETED"
	StatusArchived  ProjectStatus = "ARCHIVED"
)

type Project struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Service interface {
	CreateProject(ctx context.Context, id, name, description string) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	UpdateProjectStatus(ctx context.Context, id string, status ProjectStatus) (*Project, error)
	ListProjects(ctx context.Context) ([]*Project, error)
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

func (s *service) CreateProject(ctx context.Context, id, name, description string) (*Project, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return nil, ErrInvalidProject
	}

	now := time.Now().UTC()
	p := &Project{
		ID:          id,
		Name:        name,
		Description: description,
		Status:      StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := s.store.Exec(ctx,
		"INSERT INTO projects (id, name, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		p.ID, p.Name, p.Description, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	s.publishEvent(ctx, "projects.project.created", p)

	return p, nil
}

func (s *service) GetProject(ctx context.Context, id string) (*Project, error) {
	row := s.store.QueryRow(ctx,
		"SELECT id, name, description, status, created_at, updated_at FROM projects WHERE id = ?",
		id,
	)

	var p Project
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return &p, nil
}

func (s *service) UpdateProjectStatus(ctx context.Context, id string, status ProjectStatus) (*Project, error) {
	now := time.Now().UTC()

	result, err := s.store.Exec(ctx,
		"UPDATE projects SET status = ?, updated_at = ? WHERE id = ?",
		status, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrProjectNotFound
	}

	p, err := s.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "projects.project.updated", p)

	return p, nil
}

func (s *service) ListProjects(ctx context.Context) ([]*Project, error) {
	rows, err := s.store.Query(ctx, "SELECT id, name, description, status, created_at, updated_at FROM projects ORDER BY created_at ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var results []*Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		results = append(results, &p)
	}

	return results, nil
}

func (s *service) publishEvent(ctx context.Context, eventType string, project *Project) {
	if s.publisher == nil {
		return
	}

	id, _ := types.NewEventID(fmt.Sprintf("%s-%d", project.ID, time.Now().UnixNano()))
	event, err := events.NewBuilder().
		SetID(id).
		SetType(events.EventType(eventType)).
		SetSource("projects").
		SetPayload(project).
		Build()

	if err == nil {
		_ = s.publisher.Publish(event)
	}
}

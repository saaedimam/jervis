package projects

import (
	"errors"
	"strings"
	"sync"
	"time"

	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

var (
	ErrInvalidProject  = errors.New("projects: invalid project")
	ErrProjectNotFound = errors.New("projects: project not found")
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
	CreateProject(id, name, description string) (*Project, error)
	GetProject(id string) (*Project, error)
	UpdateProjectStatus(id string, status ProjectStatus) (*Project, error)
	ListProjects() []*Project
}

type service struct {
	mu        sync.RWMutex
	projects  map[string]*Project
	order     []string
	publisher eventcontracts.Publisher
}

func New(publisher eventcontracts.Publisher) Service {
	return &service{
		projects:  make(map[string]*Project),
		order:     make([]string, 0),
		publisher: publisher,
	}
}

func (s *service) CreateProject(id, name, description string) (*Project, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return nil, ErrInvalidProject
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.projects[id]; exists {
		return nil, ErrDuplicateProject
	}

	now := time.Now()
	p := &Project{
		ID:          id,
		Name:        name,
		Description: description,
		Status:      StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.projects[id] = p
	s.order = append(s.order, id)

	return p, nil
}

func (s *service) GetProject(id string) (*Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, exists := s.projects[id]
	if !exists {
		return nil, ErrProjectNotFound
	}

	cp := *p
	return &cp, nil
}

func (s *service) UpdateProjectStatus(id string, status ProjectStatus) (*Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, exists := s.projects[id]
	if !exists {
		return nil, ErrProjectNotFound
	}

	p.Status = status
	p.UpdatedAt = time.Now()

	cp := *p
	return &cp, nil
}

func (s *service) ListProjects() []*Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Project, 0, len(s.order))
	for _, id := range s.order {
		if p, ok := s.projects[id]; ok {
			cp := *p
			result = append(result, &cp)
		}
	}
	return result
}

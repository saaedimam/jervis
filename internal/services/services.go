package services

import (
	"context"
	"fmt"

	storecontracts "github.com/saaedimam/jervis/internal/memory/store/contracts"
	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/services/automation"
	"github.com/saaedimam/jervis/internal/services/calendar"
	"github.com/saaedimam/jervis/internal/services/habits"
	"github.com/saaedimam/jervis/internal/services/meetings"
	"github.com/saaedimam/jervis/internal/services/notion"
	"github.com/saaedimam/jervis/internal/services/planner"
	"github.com/saaedimam/jervis/internal/services/projects"
)

// Container holds all domain services for the Jervis platform.
type Container struct {
	Planner    planner.Service
	Projects   projects.Service
	Habits     habits.Service
	Meetings   meetings.Service
	Notion     notion.Service
	Calendar   calendar.Service
	Automation automation.Service
}

// NewContainer initializes all domain services with the provided store and publisher.
func NewContainer(ctx context.Context, store storecontracts.Store, publisher eventcontracts.Publisher, notionToken string) (*Container, error) {
	// Initialize schemas for all services
	if err := planner.Initialize(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to init planner: %w", err)
	}
	if err := projects.Initialize(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to init projects: %w", err)
	}
	if err := habits.Initialize(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to init habits: %w", err)
	}
	if err := meetings.Initialize(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to init meetings: %w", err)
	}

	meetingsService := meetings.New(store, publisher)

	return &Container{
		Planner:    planner.New(store, publisher),
		Projects:   projects.New(store, publisher),
		Habits:     habits.New(store, publisher),
		Meetings:   meetingsService,
		Notion:     notion.New(notionToken, publisher),
		Calendar:   calendar.New(publisher, meetingsService),
		Automation: automation.NewService(publisher),
	}, nil
}

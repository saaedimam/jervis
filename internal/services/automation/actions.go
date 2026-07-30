package automation

import (
	"context"
	"fmt"

	"github.com/saaedimam/jervis/internal/aiprovider/contracts"
	"github.com/saaedimam/jervis/internal/aiprovider/registry"
	"github.com/saaedimam/jervis/internal/services/notion"
	"github.com/saaedimam/jervis/internal/services/planner"
	"github.com/saaedimam/jervis/internal/services/projects"
)

// SyncNotionAction synchronizes local state to Notion.
type SyncNotionAction struct {
	Notion   notion.Service
	Planner  planner.Service
	Projects projects.Service
}

func (a *SyncNotionAction) Execute(ctx context.Context, payload map[string]any) error {
	syncType, _ := payload["sync_type"].(string)
	targetID, _ := payload["target_id"].(string)
	name, _ := payload["name"].(string)
	file, _ := payload["file"].(string)

	switch syncType {
	case "tasks":
		return a.Notion.SyncTasks(ctx, targetID, a.Planner)
	case "projects":
		return a.Notion.SyncProjects(ctx, targetID, a.Projects)
	case "context":
		return a.Notion.SyncContext(ctx, name, file, targetID)
	default:
		return fmt.Errorf("unsupported sync type: %s", syncType)
	}
}

// CreateTaskAction creates a new task.
type CreateTaskAction struct {
	Planner planner.Service
}

func (a *CreateTaskAction) Execute(ctx context.Context, payload map[string]any) error {
	id, _ := payload["id"].(string)
	title, _ := payload["title"].(string)
	desc, _ := payload["description"].(string)

	if id == "" || title == "" {
		return fmt.Errorf("id and title are required for CreateTaskAction")
	}

	_, err := a.Planner.CreateTask(ctx, id, title, desc)
	return err
}

// ChatAction interacts with an AI provider.
type ChatAction struct {
	Registry *registry.Registry
}

func (a *ChatAction) Execute(ctx context.Context, payload map[string]any) error {
	provider, _ := payload["provider"].(string)
	model, _ := payload["model"].(string)
	prompt, _ := payload["prompt"].(string)

	if prompt == "" {
		return fmt.Errorf("prompt is required for ChatAction")
	}

	p, err := a.Registry.Get(provider)
	if err != nil {
		return err
	}

	_, err = p.Chat(ctx, model, []contracts.Message{
		{Role: contracts.RoleUser, Content: prompt},
	}, contracts.ChatOptions{})

	return err
}

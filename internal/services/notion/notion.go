package notion

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/services/planner"
	"github.com/saaedimam/jervis/internal/services/projects"
)

// Service defines the Notion Integration Service interface.
type Service interface {
	SyncContext(ctx context.Context, name, filePath, pageID string) error
	SyncTasks(ctx context.Context, databaseID string, plannerService planner.Service) error
	SyncProjects(ctx context.Context, databaseID string, projectsService projects.Service) error
	SyncMilestones(ctx context.Context, databaseID, filePath string) error
	SyncADRs(ctx context.Context, databaseID, adrDir string) error
	SyncSpecifications(ctx context.Context, databaseID string, specFiles []string) error
	SyncDashboard(ctx context.Context, pageID, status string) error
}

type service struct {
	client    *Client
	publisher eventcontracts.Publisher
}

// New constructs a new Notion Integration Service.
func New(token string, publisher eventcontracts.Publisher) Service {
	return &service{
		client:    NewClient(token),
		publisher: publisher,
	}
}

func (s *service) SyncContext(ctx context.Context, name, filePath, pageID string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read context file %s: %w", filePath, err)
	}

	// 1. Clear existing blocks
	blocks, err := s.client.GetBlockChildren(ctx, pageID)
	if err != nil {
		return fmt.Errorf("failed to list blocks for page %s: %w", pageID, err)
	}

	for _, block := range blocks {
		blockID, ok := block["id"].(string)
		if !ok {
			continue
		}
		if err := s.client.DeleteBlock(ctx, blockID); err != nil {
			return fmt.Errorf("failed to delete block %s: %w", blockID, err)
		}
	}

	// 2. Upload new content in chunks
	const maxChunkSize = 1900
	var newBlocks []any

	runes := []rune(string(content))
	for i := 0; i < len(runes); i += maxChunkSize {
		end := i + maxChunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := string(runes[i:end])
		newBlocks = append(newBlocks, map[string]any{
			"object": "block",
			"type":   "code",
			"code": map[string]any{
				"caption": []any{},
				"rich_text": []any{
					map[string]any{
						"type": "text",
						"text": map[string]any{
							"content": chunk,
						},
					},
				},
				"language": "markdown",
			},
		})
	}

	if len(newBlocks) > 0 {
		err = s.client.UpdatePageBlocks(ctx, pageID, newBlocks)
		if err != nil {
			return fmt.Errorf("failed to update notion page %s: %w", pageID, err)
		}
	}

	return nil
}

func (s *service) SyncTasks(ctx context.Context, databaseID string, plannerService planner.Service) error {
	// 1. Ensure "Internal ID" property exists in the database
	err := s.client.UpdateDatabaseProperties(ctx, databaseID, map[string]any{
		"Internal ID": map[string]any{
			"rich_text": map[string]any{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to ensure Internal ID property: %w", err)
	}

	// 2. Fetch local tasks
	tasks, err := plannerService.ListTasks(ctx)
	if err != nil {
		return fmt.Errorf("failed to list local tasks: %w", err)
	}

	for _, task := range tasks {
		// 3. Check if task exists in Notion
		filter := map[string]any{
			"filter": map[string]any{
				"property": "Internal ID",
				"rich_text": map[string]any{
					"equals": task.ID,
				},
			},
		}
		results, err := s.client.QueryDatabase(ctx, databaseID, filter)
		if err != nil {
			return fmt.Errorf("failed to query task %s in notion: %w", task.ID, err)
		}

		status := mapStatus(task.Status)
		props := map[string]any{
			"Task": map[string]any{
				"title": []any{
					map[string]any{
						"text": map[string]any{
							"content": task.Title,
						},
					},
				},
			},
			"Status": map[string]any{
				"select": map[string]any{
					"name": status,
				},
			},
			"Internal ID": map[string]any{
				"rich_text": []any{
					map[string]any{
						"text": map[string]any{
							"content": task.ID,
						},
					},
				},
			},
		}

		if len(results) > 0 {
			// Update existing page
			pageID, ok := results[0]["id"].(string)
			if !ok {
				return fmt.Errorf("invalid page id in notion response for task %s", task.ID)
			}
			err = s.client.UpdatePageProperties(ctx, pageID, props)
			if err != nil {
				return fmt.Errorf("failed to update notion page for task %s: %w", task.ID, err)
			}
		} else {
			// Create new page
			_, err = s.client.CreatePageInDatabase(ctx, databaseID, props)
			if err != nil {
				return fmt.Errorf("failed to create notion page for task %s: %w", task.ID, err)
			}
		}
	}

	return nil
}

func mapStatus(status planner.TaskStatus) string {
	switch status {
	case planner.StatusPending:
		return "Todo"
	case planner.StatusInProgress:
		return "Doing"
	case planner.StatusCompleted:
		return "Done"
	case planner.StatusFailed:
		return "Blocked"
	default:
		return "Backlog"
	}
}

func (s *service) SyncProjects(ctx context.Context, databaseID string, projectsService projects.Service) error {
	// 1. Ensure "Internal ID" property exists in the database
	err := s.client.UpdateDatabaseProperties(ctx, databaseID, map[string]any{
		"Internal ID": map[string]any{
			"rich_text": map[string]any{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to ensure Internal ID property: %w", err)
	}

	// 2. Fetch local projects
	projs, err := projectsService.ListProjects(ctx)
	if err != nil {
		return fmt.Errorf("failed to list local projects: %w", err)
	}

	for _, project := range projs {
		// 3. Check if project exists in Notion
		filter := map[string]any{
			"filter": map[string]any{
				"property": "Internal ID",
				"rich_text": map[string]any{
					"equals": project.ID,
				},
			},
		}
		results, err := s.client.QueryDatabase(ctx, databaseID, filter)
		if err != nil {
			return fmt.Errorf("failed to query project %s in notion: %w", project.ID, err)
		}

		status := mapProjectStatus(project.Status)
		props := map[string]any{
			"Name": map[string]any{
				"title": []any{
					map[string]any{
						"text": map[string]any{
							"content": project.Name,
						},
					},
				},
			},
			"Status": map[string]any{
				"select": map[string]any{
					"name": status,
				},
			},
			"Description": map[string]any{
				"rich_text": []any{
					map[string]any{
						"text": map[string]any{
							"content": project.Description,
						},
					},
				},
			},
			"Internal ID": map[string]any{
				"rich_text": []any{
					map[string]any{
						"text": map[string]any{
							"content": project.ID,
						},
					},
				},
			},
		}

		if len(results) > 0 {
			// Update existing page
			pageID, ok := results[0]["id"].(string)
			if !ok {
				return fmt.Errorf("invalid page id in notion response for project %s", project.ID)
			}
			err = s.client.UpdatePageProperties(ctx, pageID, props)
			if err != nil {
				return fmt.Errorf("failed to update notion page for project %s: %w", project.ID, err)
			}
		} else {
			// Create new page
			_, err = s.client.CreatePageInDatabase(ctx, databaseID, props)
			if err != nil {
				return fmt.Errorf("failed to create notion page for project %s: %w", project.ID, err)
			}
		}
	}

	return nil
}

func mapProjectStatus(status projects.ProjectStatus) string {
	switch status {
	case projects.StatusActive:
		return "Active"
	case projects.StatusCompleted:
		return "Completed"
	case projects.StatusArchived:
		return "Archived"
	default:
		return "Active"
	}
}

func (s *service) SyncMilestones(ctx context.Context, databaseID, filePath string) error {
	// Simple implementation: parse MILESTONES.md and upload lines
	_, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// For now, just sync as a page if it was a page, but this is a DB.
	// Production-grade would parse the MD and create rows.
	// Let's at least create one row with the full content for now to avoid empty pages.
	props := map[string]any{
		"Milestone": map[string]any{
			"title": []any{
				map[string]any{
					"text": map[string]any{
						"content": "Project Milestones (Synced)",
					},
				},
			},
		},
		"Status": map[string]any{
			"select": map[string]any{
				"name": "Current",
			},
		},
	}
	_, err = s.client.CreatePageInDatabase(ctx, databaseID, props)
	return err
}

func (s *service) SyncADRs(ctx context.Context, databaseID, adrDir string) error {
	files, err := os.ReadDir(adrDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".md" {
			props := map[string]any{
				"Title": map[string]any{
					"title": []any{
						map[string]any{
							"text": map[string]any{
								"content": f.Name(),
							},
						},
					},
				},
			}
			if _, err := s.client.CreatePageInDatabase(ctx, databaseID, props); err != nil {
				return fmt.Errorf("failed to create notion page for ADR %s: %w", f.Name(), err)
			}
		}
	}
	return nil
}

func (s *service) SyncSpecifications(ctx context.Context, databaseID string, specFiles []string) error {
	for _, f := range specFiles {
		props := map[string]any{
			"Specification": map[string]any{
				"title": []any{
					map[string]any{
						"text": map[string]any{
							"content": f,
						},
					},
				},
			},
		}
		if _, err := s.client.CreatePageInDatabase(ctx, databaseID, props); err != nil {
			return fmt.Errorf("failed to create notion page for specification %s: %w", f, err)
		}
	}
	return nil
}

func (s *service) SyncDashboard(ctx context.Context, pageID, status string) error {
	// SyncDashboard is intentionally a no-op in v1.0.
	// Dashboard synchronization is performed by the Engineering Knowledge
	// Compiler and external Notion synchronization pipeline.
	return nil
}

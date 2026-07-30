package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ioriimasu/jervis/internal/app"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server wraps the MCP server implementation.
type Server struct {
	mcpServer *server.MCPServer
	app       *app.App
}

// NewServer creates a new Jervis MCP server.
func NewServer(a *app.App) *Server {
	s := server.NewMCPServer(
		"Jervis OS",
		"1.0.0",
		server.WithLogging(),
	)

	srv := &Server{
		mcpServer: s,
		app:       a,
	}

	srv.registerTools()
	srv.registerResources()

	return srv
}

// Serve starts the MCP server on stdio.
func (s *Server) Serve() error {
	return server.ServeStdio(s.mcpServer)
}

func (s *Server) registerTools() {
	// Planner Tools
	s.mcpServer.AddTool(mcp.NewTool("create_task",
		mcp.WithDescription("Create a new task in the Jervis planner"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Unique task ID")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Task title")),
		mcp.WithString("description", mcp.Description("Detailed task description")),
	), s.handleCreateTask)

	s.mcpServer.AddTool(mcp.NewTool("list_tasks",
		mcp.WithDescription("List all tasks in the Jervis planner"),
	), s.handleListTasks)

	// Notion Tools
	s.mcpServer.AddTool(mcp.NewTool("notion_sync",
		mcp.WithDescription("Synchronize local state to Notion"),
		mcp.WithString("type", mcp.Required(), mcp.Description("Sync type: context, tasks, projects")),
		mcp.WithString("name", mcp.Description("Context name (for type=context)")),
		mcp.WithString("file", mcp.Description("Local file path (for type=context)")),
		mcp.WithString("id", mcp.Description("Notion ID (optional, defaults to env)")),
	), s.handleNotionSync)
}

func (s *Server) registerResources() {
	// Expose files in context/ as resources
	contextDir := "context"
	files, _ := os.ReadDir(contextDir)
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".md" {
			name := f.Name()
			uri := fmt.Sprintf("file://context/%s", name)
			s.mcpServer.AddResource(mcp.NewResource(uri, name,
				mcp.WithResourceDescription(fmt.Sprintf("Jervis context document: %s", name)),
				mcp.WithMIMEType("text/markdown"),
			), s.handleReadResource)
		}
	}
}

func (s *Server) handleCreateTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}
	id, _ := args["id"].(string)
	title, _ := args["title"].(string)
	desc, _ := args["description"].(string)

	task, err := s.app.Services.Planner.CreateTask(ctx, id, title, desc)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Task created: [%s] %s", task.ID, task.Title)), nil
}

func (s *Server) handleListTasks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tasks, err := s.app.Services.Planner.ListTasks(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var result string
	for _, t := range tasks {
		result += fmt.Sprintf("- [%s] %s (%s)\n", t.ID, t.Title, t.Status)
	}
	return mcp.NewToolResultText(result), nil
}

func (s *Server) handleNotionSync(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}
	syncType, _ := args["type"].(string)
	name, _ := args["name"].(string)
	file, _ := args["file"].(string)
	targetID, _ := args["id"].(string)

	var err error
	switch syncType {
	case "context":
		if targetID == "" {
			targetID = os.Getenv("MASTER_CONTEXT_ID")
		}
		err = s.app.Services.Notion.SyncContext(ctx, name, file, targetID)
	case "tasks":
		if targetID == "" {
			targetID = os.Getenv("TASKS_DB")
		}
		err = s.app.Services.Notion.SyncTasks(ctx, targetID, s.app.Services.Planner)
	case "projects":
		if targetID == "" {
			targetID = os.Getenv("PACKAGES_DB")
		}
		err = s.app.Services.Notion.SyncProjects(ctx, targetID, s.app.Services.Projects)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown sync type: %s", syncType)), nil
	}

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Sync successful!"), nil
}

func (s *Server) handleReadResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Validate URI prefix - only allow context/ files
	const prefix = "file://context/"
	if len(request.Params.URI) <= len(prefix) || request.Params.URI[:len(prefix)] != prefix {
		return nil, fmt.Errorf("invalid URI scheme or path")
	}

	// Extract and validate the path
	relPath := request.Params.URI[len(prefix):]

	// Clean the path to prevent traversal
	cleanPath := filepath.Clean(relPath)

	// Reject if path tries to escape context/ directory
	if cleanPath == ".." || len(cleanPath) >= 3 && cleanPath[:3] == ".." {
		return nil, fmt.Errorf("path traversal not allowed")
	}

	// Build full path and verify it stays under context/
	fullPath := filepath.Join("context", cleanPath)
	absBase, _ := filepath.Abs("context")
	absFull, _ := filepath.Abs(fullPath)

	// Ensure the resolved path is under context/
	if len(absFull) < len(absBase) || absFull[:len(absBase)] != absBase {
		return nil, fmt.Errorf("path escapes allowed directory")
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			Text:     string(content),
			MIMEType: "text/markdown",
		},
	}, nil
}

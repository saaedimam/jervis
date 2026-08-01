package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/saaedimam/jervis/internal/aiprovider/contracts"
	"github.com/saaedimam/jervis/internal/app"
	"github.com/saaedimam/jervis/internal/interfaces/mcp"
	"github.com/saaedimam/jervis/internal/interfaces/rest"
	"github.com/saaedimam/jervis/internal/runtime/buildinfo"
)

func main() {
	_ = godotenv.Load()
	if len(os.Args) < 2 {
		usage()
		return
	}

	command := os.Args[1]
	switch command {
	case "version", "--version":
		printVersion()
	case "planner":
		runPlanner(os.Args[2:])
	case "sync":
		runSync(os.Args[2:])
	case "calendar":
		runCalendar(os.Args[2:])
	case "chat":
		runChat(os.Args[2:])
	case "mcp":
		runMCP()
	case "api":
		runAPI(os.Args[2:])
	case "automation":
		runAutomation(os.Args[2:])
	case "daemon":
		runDaemon()
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Jervis Personal OS — AI-Native Engineering Platform")
	fmt.Println("\nUsage:")
	fmt.Println("  jervis version         Show build information")
	fmt.Println("  jervis planner         Manage planned tasks")
	fmt.Println("  jervis sync            Sync local state to Notion (use --type context|tasks|projects)")
	fmt.Println("  jervis calendar        Manage calendar integrations (sync/export)")
	fmt.Println("  jervis chat            Chat with AI providers")
	fmt.Println("  jervis mcp             Start MCP server")
	fmt.Println("  jervis api             Start REST API server")
	fmt.Println("  jervis automation      Manage automation workflows")
	fmt.Println("  jervis daemon          Start runtime background daemon")
	fmt.Println("\nRun 'jervis [command] --help' for more information on a command.")
}

func printVersion() {
	info := buildinfo.Get()
	fmt.Printf("Jervis OS v%s\n", info.SemVer())
	fmt.Printf("Commit: %s\n", info.GitCommit())
	fmt.Printf("Built:  %s\n", info.BuildDate())
}

func runPlanner(args []string) {
	plannerCmd := flag.NewFlagSet("planner", flag.ExitOnError)
	createCmd := plannerCmd.Bool("create", false, "Create a new task")
	listCmd := plannerCmd.Bool("list", false, "List all tasks")
	id := plannerCmd.String("id", "", "Task ID")
	title := plannerCmd.String("title", "", "Task Title")
	desc := plannerCmd.String("desc", "", "Task Description")

	_ = plannerCmd.Parse(args)

	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Error initializing Jervis: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	if *createCmd {
		if *id == "" || *title == "" {
			fmt.Println("Error: --id and --title are required for creation")
			os.Exit(1)
		}
		t, err := a.Services.Planner.CreateTask(ctx, *id, *title, *desc)
		if err != nil {
			fmt.Printf("Error creating task: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task created: [%s] %s\n", t.ID, t.Title)
	} else if *listCmd {
		tasks, err := a.Services.Planner.ListTasks(ctx)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ID\tStatus\tTitle")
		fmt.Println("--\t------\t-----")
		for _, t := range tasks {
			fmt.Printf("%s\t%s\t%s\n", t.ID, t.Status, t.Title)
		}
	} else {
		plannerCmd.Usage()
	}
}

func runDaemon() {
	fmt.Println("Starting Jervis Daemon...")
	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	fmt.Println("Jervis OS Runtime Ready.")
	for {
		time.Sleep(1 * time.Hour)
	}
}

func runSync(args []string) {
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncType := syncCmd.String("type", "context", "Sync type (context, tasks, projects)")
	name := syncCmd.String("name", "MASTER_CONTEXT", "Context name (for type=context)")
	file := syncCmd.String("file", "context/MASTER_CONTEXT.md", "Local file path (for type=context)")
	id := syncCmd.String("id", "", "Notion ID (Page ID for context, Database ID for tasks/projects)")

	_ = syncCmd.Parse(args)

	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")
	cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.GoogleKey = os.Getenv("GOOGLE_API_KEY")
	cfg.OllamaBaseURL = os.Getenv("OLLAMA_BASE_URL")

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	var syncErr error
	targetID := *id
	switch *syncType {
	case "context":
		if targetID == "" {
			targetID = os.Getenv("MASTER_CONTEXT_ID")
		}
		if targetID == "" {
			fmt.Println("Error: --id or MASTER_CONTEXT_ID env var is required")
			os.Exit(1)
		}
		fmt.Printf("Syncing context %s to Notion page %s...\n", *file, targetID)
		syncErr = a.Services.Notion.SyncContext(ctx, *name, *file, targetID)
	case "tasks":
		if targetID == "" {
			targetID = os.Getenv("TASKS_DB")
		}
		if targetID == "" {
			fmt.Println("Error: --id or TASKS_DB env var is required")
			os.Exit(1)
		}
		fmt.Printf("Syncing tasks to Notion database %s...\n", targetID)
		syncErr = a.Services.Notion.SyncTasks(ctx, targetID, a.Services.Planner)
	case "projects":
		if targetID == "" {
			targetID = os.Getenv("PACKAGES_DB")
		}
		if targetID == "" {
			fmt.Println("Error: --id or PACKAGES_DB env var is required")
			os.Exit(1)
		}
		fmt.Printf("Syncing projects to Notion database %s...\n", targetID)
		syncErr = a.Services.Notion.SyncProjects(ctx, targetID, a.Services.Projects)
	case "all":
		fmt.Println("Starting full synchronization...")

		// 1. Context
		fmt.Print("Syncing MASTER_CONTEXT... ")
		if err := a.Services.Notion.SyncContext(ctx, "MASTER_CONTEXT", "context/MASTER_CONTEXT.md", os.Getenv("MASTER_CONTEXT_ID")); err != nil {
			fmt.Printf("FAIL: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		// 2. Tasks
		fmt.Print("Syncing Tasks... ")
		if err := a.Services.Notion.SyncTasks(ctx, os.Getenv("TASKS_DB"), a.Services.Planner); err != nil {
			fmt.Printf("FAIL: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		// 3. Projects
		fmt.Print("Syncing Projects... ")
		if err := a.Services.Notion.SyncProjects(ctx, os.Getenv("PACKAGES_DB"), a.Services.Projects); err != nil {
			fmt.Printf("FAIL: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		// 4. Milestones
		fmt.Print("Syncing Milestones... ")
		if err := a.Services.Notion.SyncMilestones(ctx, os.Getenv("MILESTONES_DB"), "context/MILESTONES.md"); err != nil {
			fmt.Printf("FAIL: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		// 5. ADRs
		fmt.Print("Syncing ADRs... ")
		if err := a.Services.Notion.SyncADRs(ctx, os.Getenv("ADRS_DB"), "context"); err != nil {
			fmt.Printf("FAIL: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		// 6. Specifications
		fmt.Print("Syncing Specifications... ")
		specs := []string{"ARCHITECTURE.md", "01_ENGINEERING_PRINCIPLES.md", "06_SECURITY_MODEL.md"}
		if err := a.Services.Notion.SyncSpecifications(ctx, os.Getenv("SPECIFICATIONS_DB"), specs); err != nil {
			fmt.Printf("FAIL: %v\n", err)
		} else {
			fmt.Println("OK")
		}

		fmt.Println("Full synchronization complete.")
		return
	default:
		fmt.Printf("Error: Unknown sync type: %s\n", *syncType)
		os.Exit(1)
	}

	if syncErr != nil {
		fmt.Printf("Sync failed: %v\n", syncErr)
		os.Exit(1)
	}
	fmt.Println("Sync successful!")
}

func runCalendar(args []string) {
	calCmd := flag.NewFlagSet("calendar", flag.ExitOnError)
	importCmd := calCmd.Bool("import", false, "Import events from iCal URL")
	exportCmd := calCmd.Bool("export", false, "Export local meetings to iCal")
	url := calCmd.String("url", "", "iCal feed URL (for import)")

	_ = calCmd.Parse(args)

	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	if *importCmd {
		if *url == "" {
			fmt.Println("Error: --url is required for import")
			os.Exit(1)
		}
		fmt.Printf("Importing calendar from %s...\n", *url)
		err = a.Services.Calendar.ImportICal(ctx, *url)
		if err != nil {
			fmt.Printf("Import failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Import successful!")
	} else if *exportCmd {
		fmt.Println("Exporting local meetings to iCal...")
		data, err := a.Services.Calendar.ExportICal(ctx)
		if err != nil {
			fmt.Printf("Export failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(data)
	} else {
		calCmd.Usage()
	}
}

func runChat(args []string) {
	chatCmd := flag.NewFlagSet("chat", flag.ExitOnError)
	provider := chatCmd.String("provider", "ollama", "AI provider (openai, anthropic, google, ollama)")
	model := chatCmd.String("model", "llama3", "Model name")
	system := chatCmd.String("system", "You are Jervis, a helpful AI assistant.", "System prompt")
	prompt := chatCmd.String("prompt", "", "User prompt")

	_ = chatCmd.Parse(args)

	if *prompt == "" {
		fmt.Println("Error: --prompt is required")
		os.Exit(1)
	}

	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")
	cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.GoogleKey = os.Getenv("GOOGLE_API_KEY")
	cfg.OllamaBaseURL = os.Getenv("OLLAMA_BASE_URL")
	cfg.DefaultAIProv = *provider

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	p, err := a.AIProviders.Get(*provider)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	messages := []contracts.Message{
		{Role: contracts.RoleSystem, Content: *system},
		{Role: contracts.RoleUser, Content: *prompt},
	}

	opts := contracts.ChatOptions{
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	fmt.Printf("[%s/%s] thinking...\n", *provider, *model)
	resp, err := p.Chat(ctx, *model, messages, opts)
	if err != nil {
		fmt.Printf("Chat failed: %v\n", err)
		os.Exit(1)
	}

	// Check if response has choices to avoid panic
	if resp == nil || len(resp.Choices) == 0 {
		fmt.Println("Chat returned no response choices")
		os.Exit(1)
	}

	fmt.Println("\n" + resp.Choices[0].Message.Content)
	fmt.Printf("\n(Tokens: P=%d, C=%d, T=%d)\n", resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
}

func runMCP() {
	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")
	cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.GoogleKey = os.Getenv("GOOGLE_API_KEY")
	cfg.OllamaBaseURL = os.Getenv("OLLAMA_BASE_URL")

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	s := mcp.NewServer(a)
	if err := s.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server failed: %v\n", err)
		os.Exit(1)
	}
}

func runAPI(args []string) {
	apiCmd := flag.NewFlagSet("api", flag.ExitOnError)
	port := apiCmd.Int("port", 8080, "API server port")
	_ = apiCmd.Parse(args)

	// -----------------------------------------------------------
	// Configuration validation – the REST API *must* have an auth key.
	// The server package itself never reads the environment.
	// -----------------------------------------------------------
	authKey := strings.TrimSpace(os.Getenv("JERVIS_API_AUTH_KEY"))
	if authKey == "" {
		fmt.Fprintln(os.Stderr, "Error: JERVIS_API_AUTH_KEY environment variable is required for the REST API")
		os.Exit(1)
	}

	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")
	cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.GoogleKey = os.Getenv("GOOGLE_API_KEY")
	cfg.OllamaBaseURL = os.Getenv("OLLAMA_BASE_URL")

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	// Explicitly pass the validated key to the server constructor.
	s := rest.NewServerWithAuth(a, *port, authKey)
	if err := s.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "API server failed: %v\n", err)
		os.Exit(1)
	}
}

func runAutomation(args []string) {
	autoCmd := flag.NewFlagSet("automation", flag.ExitOnError)
	list := autoCmd.Bool("list", false, "List all workflows")
	execute := autoCmd.String("execute", "", "Execute a workflow by ID")
	_ = autoCmd.Parse(args)

	ctx := context.Background()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = "jervis.db"
	cfg.NotionToken = os.Getenv("NOTION_TOKEN")
	cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.GoogleKey = os.Getenv("GOOGLE_API_KEY")
	cfg.OllamaBaseURL = os.Getenv("OLLAMA_BASE_URL")

	a, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = a.Close() }()

	if *list {
		workflows := a.Services.Automation.Registry().List()
		fmt.Println("Registered Workflows:")
		for _, w := range workflows {
			fmt.Printf("- [%s] %s\n", w.ID, w.Name)
		}
		return
	}

	if *execute != "" {
		wf, exists := a.Services.Automation.Registry().Get(*execute)
		if !exists {
			fmt.Printf("Error: Workflow %s not found\n", *execute)
			os.Exit(1)
		}
		// In a real scenario, we'd use the service to execute it properly
		// For now, let's just confirm it exists.
		fmt.Printf("Executing workflow: %s\n", wf.Name)
		// ... execution logic ...
		return
	}

	autoCmd.Usage()
}

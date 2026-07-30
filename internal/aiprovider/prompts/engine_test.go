package prompts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEngine(t *testing.T) {
	// Create a temporary directory for templates
	dir, err := os.MkdirTemp("", "prompts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	engine := NewEngine(dir)

	// Test Hydrate not found
	if _, err := engine.Hydrate("missing", nil); err == nil {
		t.Error("expected error for missing template")
	}

	// Test ListTemplates empty
	templates, err := engine.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Errorf("expected 0 templates, got %d", len(templates))
	}

	// Write a test template
	tmplContent := `Hello {{.Name}}!`
	if err := os.WriteFile(filepath.Join(dir, "greeting.md"), []byte(tmplContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Also write a non-md file and a directory
	if err := os.WriteFile(filepath.Join(dir, "ignore.txt"), []byte("ignore"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub.md"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Test ListTemplates again
	templates, err = engine.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 1 || templates[0] != "greeting" {
		t.Errorf("expected [greeting], got %v", templates)
	}

	// Test Hydrate success
	data := struct{ Name string }{Name: "World"}
	res, err := engine.Hydrate("greeting", data)
	if err != nil {
		t.Fatal(err)
	}
	if res != "Hello World!" {
		t.Errorf("expected 'Hello World!', got %s", res)
	}

	// Write an invalid template
	invalidContent := `Hello {{.Name` // syntax error
	if err := os.WriteFile(filepath.Join(dir, "invalid.md"), []byte(invalidContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test Hydrate parse error
	if _, err := engine.Hydrate("invalid", nil); err == nil {
		t.Error("expected parse error")
	}

	// Test Hydrate execution error
	// For execution error, missing field with strict mode? In Go templates, missing fields might just be empty or error depending on options.
	// We can provoke an error by calling a non-existent method.
	execErrContent := `Hello {{.MissingMethod}}`
	if err := os.WriteFile(filepath.Join(dir, "execerr.md"), []byte(execErrContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Hydrate("execerr", struct{}{}); err == nil {
		t.Error("expected execution error")
	}
}

func TestEngine_ListError(t *testing.T) {
	// Point to a non-existent directory
	engine := NewEngine("/does/not/exist/definitely")
	if _, err := engine.ListTemplates(); err == nil {
		t.Error("expected error listing non-existent directory")
	}
}

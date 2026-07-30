package prompts

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Engine handles loading and hydrating prompt templates.
type Engine struct {
	templateDir string
}

// NewEngine creates a new prompt engine.
func NewEngine(templateDir string) *Engine {
	return &Engine{
		templateDir: templateDir,
	}
}

// Hydrate loads a template by name and hydrants it with the provided data.
func (e *Engine) Hydrate(name string, data any) (string, error) {
	path := filepath.Join(e.templateDir, name+".md")
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to hydrate template %s: %w", name, err)
	}

	return buf.String(), nil
}

// ListTemplates returns a list of available template names.
func (e *Engine) ListTemplates() ([]string, error) {
	files, err := os.ReadDir(e.templateDir)
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".md" {
			name := f.Name()[:len(f.Name())-3]
			templates = append(templates, name)
		}
	}
	return templates, nil
}

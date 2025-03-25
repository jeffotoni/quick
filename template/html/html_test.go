package html

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTestFiles prepares a temporary directory with base and layout templates
func createTestFiles(t *testing.T, baseDir string) {
	t.Helper()

	err := os.MkdirAll(filepath.Join(baseDir, "layouts"), 0755)
	if err != nil {
		t.Fatalf("failed to create layout directory: %v", err)
	}

	// Create index.html as the base template
	indexHTML := `
		<h1>{{ .Title }}</h1>
		<p>{{ .Message }}</p>
	`
	err = os.WriteFile(filepath.Join(baseDir, "index.html"), []byte(indexHTML), 0644)
	if err != nil {
		t.Fatalf("failed to write index.html: %v", err)
	}

	// Create layouts/main.html with a yield placeholder
	mainLayout := `
		<!DOCTYPE html>
		<html>
		<head><title>Layout</title></head>
		<body>{{ .yield }}</body>
		</html>
	`
	err = os.WriteFile(filepath.Join(baseDir, "layouts", "main.html"), []byte(mainLayout), 0644)
	if err != nil {
		t.Fatalf("failed to write main.html: %v", err)
	}
}

// TestRenderNoLayout ensures template rendering works without layout
func TestRenderNoLayout(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir)

	engine := New(dir, ".html")
	if err := engine.Load(); err != nil {
		t.Fatalf("failed to load templates: %v", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Title":   "Hello",
		"Message": "World",
	}

	if err := engine.Render(&buf, "index.html", data); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Hello") || !strings.Contains(output, "World") {
		t.Errorf("output missing expected content: %s", output)
	}
}

// TestRenderWithLayout verifies rendering with a layout wrapper
func TestRenderWithLayout(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir)

	engine := New(dir, ".html")
	if err := engine.Load(); err != nil {
		t.Fatalf("template loading failed: %v", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Title":   "Layout Test",
		"Message": "With layout",
	}

	err := engine.Render(&buf, "index.html", data, "layouts/main.html")
	if err != nil {
		t.Fatalf("layout rendering failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<html>") || !strings.Contains(output, "</html>") {
		t.Errorf("layout not applied: %s", output)
	}
	if !strings.Contains(output, "Layout Test") {
		t.Errorf("base content missing: %s", output)
	}
}

// TestAddFunc checks that custom template functions are applied correctly
func TestAddFunc(t *testing.T) {
	dir := t.TempDir()

	templateContent := `{{ upper .Name }}`
	err := os.WriteFile(filepath.Join(dir, "hello.html"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	engine := New(dir, ".html")
	engine.AddFunc("upper", strings.ToUpper)

	if err := engine.Load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	var buf bytes.Buffer
	err = engine.Render(&buf, "hello.html", map[string]string{"Name": "gopher"})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(buf.String(), "GOPHER") {
		t.Errorf("expected 'GOPHER', got: %s", buf.String())
	}
}

func TestTemplateNameAliases(t *testing.T) {
	// Create temporary templates in a test directory
	dir := createTestTemplates(t)
	defer os.RemoveAll(dir)

	// Initialize the engine with the test template directory
	engine := New(dir, ".html")

	// Load templates into memory
	if err := engine.Load(); err != nil {
		t.Fatalf("failed to load templates: %v", err)
	}

	// Render using alias without extension ("index")
	buf1 := &bytes.Buffer{}
	data := map[string]any{"Title": "Alias 1", "Message": "Test"}
	err := engine.Render(buf1, "index", data)
	if err != nil {
		t.Fatalf("render index failed: %v", err)
	}

	// Render using alias with extension ("index.html")
	buf2 := &bytes.Buffer{}
	err = engine.Render(buf2, "index.html", data) // reusar o mesmo "data"!
	if err != nil {
		t.Fatalf("render index.html failed: %v", err)
	}

	// Both render outputs must be identical
	if buf1.String() != buf2.String() {
		t.Error("Expected both 'index' and 'index.html' to render the same output")
	}
}

func createTestTemplates(tb testing.TB) string {
	tmpDir := tb.TempDir()

	// Create subfolders
	layoutDir := filepath.Join(tmpDir, "layouts")
	err := os.MkdirAll(layoutDir, 0755)
	if err != nil {
		tb.Fatalf("failed to create layouts dir: %v", err)
	}

	// index.html
	writeFile(tb, filepath.Join(tmpDir, "index.html"), `
		<h1>{{ .Title }}</h1>
		<p>{{ .Message }}</p>
	`)

	// layouts/main.html
	writeFile(tb, filepath.Join(layoutDir, "main.html"), `
		<html><body>
		<header><h2>Main Layout</h2></header>
		<div>{{ .yield }}</div>
		</body></html>
	`)

	// layouts/base.html
	writeFile(tb, filepath.Join(layoutDir, "base.html"), `
		<!DOCTYPE html>
		<html><body>
		<nav>Base Layout</nav>
		<section>{{ .yield }}</section>
		</body></html>
	`)

	return tmpDir
}

// writeFile is a helper to write content to a file
func writeFile(tb testing.TB, path, content string) {
	err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0644)
	if err != nil {
		tb.Fatalf("failed to write file %s: %v", path, err)
	}
}

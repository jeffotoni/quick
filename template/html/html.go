// Package html provides a powerful and extensible HTML template engine
// for the Quick web framework, built on top of Go's html/template.
//
// It supports rendering templates with optional layout wrapping and
// allows templates to be loaded from either the local filesystem or
// an embedded filesystem (fs.FS), such as embed.FS.
//
// Key features:
//   - Dynamic HTML rendering using native html/template
//   - Support for multiple layout layers (nested layouts)
//   - Custom template functions via AddFunc
//   - Template aliasing (load templates with multiple names/paths)
//   - Lazy loading support via Render if Load is not called explicitly
package html

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Engine is a customizable HTML template engine designed for use with the Quick web framework.
// It supports parsing templates from the local filesystem or an embedded file system (fs.FS),
// as well as optional layout composition and function maps for advanced rendering.
type Engine struct {
	templates *template.Template // Parsed and compiled templates
	funcMap   template.FuncMap   // Custom template functions (e.g., "upper", "formatDate", etc.)
	Dir       string             // Root directory of templates (used when not using fs.FS)
	ext       string             // File extension for template files (e.g., ".html")
	fileSys   fs.FS              // Optional embedded filesystem (e.g., embed.FS)
}

// New returns a new Engine configured to load templates from the local filesystem.
func New(dir, ext string) *Engine {
	return &Engine{
		funcMap: make(template.FuncMap),
		Dir:     dir,
		ext:     ext,
	}
}

// NewFileSystem returns a new Engine configured to load templates from an fs.FS (e.g., embed.FS).
func NewFileSystem(fsys fs.FS, ext string) *Engine {
	return &Engine{
		funcMap: make(template.FuncMap),
		fileSys: fsys,
		ext:     ext,
	}
}

// AddFunc registers a custom function to the template engine.
// These functions can be used within the templates.
func (e *Engine) AddFunc(name string, fn interface{}) {
	e.funcMap[name] = fn
}

// Load parses and loads all templates from either the local filesystem or an embedded fs.FS.
// It recursively walks through the configured directory, loading all files with the specified extension.
// The parsed templates are stored in memory and can be rendered later via the Render method.
//
// Returns an error if any template file cannot be read or parsed.
func (e *Engine) Load() error {
	e.templates = template.New("").Funcs(e.funcMap)

	// Embedded filesystem (e.fileSys)
	if e.fileSys != nil {
		return fs.WalkDir(e.fileSys, e.Dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, e.ext) {
				return nil
			}

			content, err := fs.ReadFile(e.fileSys, path)
			if err != nil {
				return err
			}

			relPath := strings.TrimPrefix(path, e.Dir)
			relPath = strings.TrimPrefix(relPath, "/") // Normalize
			relSlash := filepath.ToSlash(relPath)
			fullSlash := filepath.ToSlash(path)
			baseName := strings.TrimSuffix(relSlash, e.ext)

			// Register all aliases
			aliases := []string{
				fullSlash, // e.g., views/index.html
				relSlash,  // e.g., index.html
				baseName,  // e.g., index
			}

			for _, alias := range aliases {
				if _, err := e.templates.New(alias).Parse(string(content)); err != nil {
					return err
				}
			}

			return nil
		})
	}

	// Local filesystem
	return filepath.Walk(e.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, e.ext) {
			return err
		}

		relPath, err := filepath.Rel(e.Dir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relSlash := filepath.ToSlash(relPath)
		baseName := strings.TrimSuffix(relSlash, e.ext)

		aliases := []string{
			relSlash, // index.html
			baseName, // index
		}

		for _, alias := range aliases {
			if _, err := e.templates.New(alias).Parse(string(content)); err != nil {
				return err
			}
		}

		return nil
	})
}

// Render renders a named template and optionally wraps it with one or more layouts.
//
// The first parameter `name` is the base template to render.
// Optional `layouts` wrap the base template, with the outermost layout listed last.
// Within each layout, the base template is inserted using the {{ .yield }} variable.
func (e *Engine) Render(w io.Writer, name string, data interface{}, layouts ...string) error {
	// Lazy loading (in case Load was not called manually)
	if e.templates == nil {
		if err := e.Load(); err != nil {
			return err
		}
	}

	tmpl := e.templates.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("template %q not found", name)
	}

	// Render the base template to string
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return err
	}
	content := sb.String()

	// Wrap with layouts from innermost to outermost
	for i := len(layouts) - 1; i >= 0; i-- {
		layout := e.templates.Lookup(layouts[i])
		if layout == nil {
			return fmt.Errorf("layout %q not found", layouts[i])
		}

		sb.Reset()
		err := layout.Execute(&sb, map[string]interface{}{
			"yield": template.HTML(content), // safely inject rendered inner content
		})
		if err != nil {
			return err
		}
		content = sb.String()
	}

	_, err := io.WriteString(w, content)
	return err
}

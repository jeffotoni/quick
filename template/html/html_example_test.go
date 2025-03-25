//go:build !exclude_test

package html

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
)

//go:embed static/views/*.html static/views/layouts/*.html
var templates embed.FS

// ExampleNew demonstrates loading templates from filesystem.
// This function is named ExampleNew()
// it with the Examples type.
func ExampleNew() {
	engine := New("./static/views", ".html")
	engine.AddFunc("upper", strings.ToUpper)
	_ = engine.Load()

	var buf bytes.Buffer
	engine.Render(&buf, "index", map[string]interface{}{
		"Title":   "From FS",
		"Message": "Hello Quick",
	})

	fmt.Println("ok")
	// Output: ok
}

// ExampleNewFileSystem demonstrates loading templates from embed.FS.
// This function is named ExampleNewFileSystem()
// it with the Examples type.
func ExampleNewFileSystem() {
	engine := NewFileSystem(templates, ".html")
	engine.Dir = "static/views"
	_ = engine.Load()

	var buf bytes.Buffer
	engine.Render(&buf, "index", map[string]interface{}{
		"Title":   "From embed",
		"Message": "Hello from embed.FS",
	})

	fmt.Println("ok")
	// Output: ok
}

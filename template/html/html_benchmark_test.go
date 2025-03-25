package html

import (
	"io"
	"testing"
)

// Benchmark rendering without any layout
func BenchmarkRenderNoLayout(b *testing.B) {
	dir := createTestTemplates(b)
	engine := New(dir, ".html")

	if err := engine.Load(); err != nil {
		b.Fatalf("failed to load templates: %v", err)
	}

	data := map[string]interface{}{
		"Title":   "Benchmark",
		"Message": "No layout",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := engine.Render(io.Discard, "index.html", data)
		if err != nil {
			b.Fatalf("render failed: %v", err)
		}
	}
}

// Benchmark rendering with a layout
func BenchmarkRenderWithLayout(b *testing.B) {
	dir := createTestTemplates(b)
	engine := New(dir, ".html")

	if err := engine.Load(); err != nil {
		b.Fatalf("failed to load templates: %v", err)
	}

	data := map[string]interface{}{
		"Title":   "Benchmark",
		"Message": "With layout",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := engine.Render(io.Discard, "index.html", data, "layouts/main.html")
		if err != nil {
			b.Fatalf("render with layout failed: %v", err)
		}
	}
}

// Benchmark rendering with nested layouts
func BenchmarkRenderWithNestedLayouts(b *testing.B) {
	dir := createTestTemplates(b)
	engine := New(dir, ".html")

	if err := engine.Load(); err != nil {
		b.Fatalf("failed to load templates: %v", err)
	}

	data := map[string]interface{}{
		"Title":   "Benchmark",
		"Message": "Nested layout",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := engine.Render(io.Discard, "index.html", data, "layouts/main.html", "layouts/base.html")
		if err != nil {
			b.Fatalf("render with nested layout failed: %v", err)
		}
	}
}

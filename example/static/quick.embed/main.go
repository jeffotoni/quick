//go:build !exclude_test

// Embed.FS allows you to include files directly into
// the binary during compilation, eliminating the need to load files
// from the file system at runtime. This means that
// static files (HTML, CSS, JS, images, etc.)
// are embedded into the executable.
package main

import (
	"embed"

	"github.com/jeffotoni/quick"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	q := quick.New()

	q.Static("/static", staticFiles)

	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/index.html")
		return nil
	})

	q.Listen("0.0.0.0:8080")
}

// $ curl --location 'http://localhost:8080/'

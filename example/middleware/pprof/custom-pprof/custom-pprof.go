package main

import (
	"strings"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	q := quick.New()

	q.Use(pprof.New(pprof.Config{
		Prefix: "/custom-pprof",
		Next: func(c *quick.Ctx) bool {
			// Skip pprof if path contains "skip"
			return strings.Contains(c.Path(), "skip")
		},
	}))

	// Declare this route first — it must come before the wildcard
	q.Get("/custom-pprof/skip-this", func(c *quick.Ctx) error {
		return c.String("This route bypassed the profiler")
	})

	// Wildcard route to let pprof handle other requests
	q.Get("/custom-pprof*", func(c *quick.Ctx) error {
		return c.Next()
	})

	q.Listen(":8080")
}

// $ curl http://localhost:8080/custom-pprof/skip-this
// $ curl http://localhost:8080/custom-pprof/

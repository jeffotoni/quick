package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	q := quick.New()

	// Trace-ID middleware
	q.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			traceID := rand.TraceID() // internal quick

			// Inject the Trace-ID into the request and response header
			r.Header.Set("X-Trace-ID", traceID)
			w.Header().Set("X-Trace-ID", traceID)

			// Print simple log with Trace-ID
			log.Printf("[Trace-ID: %s] -> Request start %s %s\n", traceID, r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
			duration := time.Since(start)

			log.Printf("[Trace-ID: %s] <- End of request duration:[(%v)]\n", traceID, duration)
		})
	})

	q.Get("/v1/user/:name", func(c *quick.Ctx) error {
		name := c.Param("name")
		c.Set("Content-Type", "application/json")
		return c.Status(200).JSON(quick.M{
			"msg": name,
		})
	})

	q.Listen("0.0.0.0:8080")
}

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/basicauth"
	"github.com/jeffotoni/quick/middleware/cache"
	"github.com/jeffotoni/quick/middleware/compress"
	"github.com/jeffotoni/quick/middleware/cors"
	"github.com/jeffotoni/quick/middleware/healthcheck"
	"github.com/jeffotoni/quick/middleware/helmet"
	"github.com/jeffotoni/quick/middleware/limiter"
	"github.com/jeffotoni/quick/middleware/logger"
	"github.com/jeffotoni/quick/middleware/maxbody"
	"github.com/jeffotoni/quick/middleware/msgid"
	"github.com/jeffotoni/quick/middleware/msguuid"
	"github.com/jeffotoni/quick/middleware/pprof"
	"github.com/jeffotoni/quick/middleware/recover"
)

// Example of Server-Sent Events (SSE) with Quick - Simple approach without loop
// This example demonstrates sending a single SSE message to the client.
// SSE is useful for real-time updates from server to client over HTTP.
//
// Key SSE Headers:
// - Content-Type: text/event-stream
// - Cache-Control: no-cache
// - Connection: keep-alive
// - Access-Control-Allow-Origin: * (for CORS)

func main() {
	q := quick.New()

	// must use at the beginning, if the path is all
	q.Use(healthcheck.New(
		healthcheck.Options{
			App: q,
		},
	))

	q.Use(recover.New())

	q.Use(pprof.New())

	q.Use(maxbody.New(50000)) // 50KB

	q.Use(helmet.Helmet())

	q.Use(msguuid.New())
	q.Use(msgid.New())

	q.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           600,
		Debug:            false,
	}))

	q.Use(compress.Gzip())

	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "DEBUG",
	}))

	q.Use(basicauth.BasicAuth("user", "adm"))

	q.Use(cache.New())

	q.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *quick.Ctx) string {
			return c.RemoteIP()
		},
		LimitReached: func(c *quick.Ctx) error {
			c.Set("Content-Type", "application/json")
			return c.Status(quick.StatusTooManyRequests).SendString(`{"msg":"Much Request #bloqued"}`)
		},
	}))

	// Simple SSE endpoint - sends a single event
	q.Get("/events/simple", func(c *quick.Ctx) error {
		// Set SSE headers manually
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Access-Control-Allow-Origin", "*")

		// Get the Flusher to send data immediately
		flusher, ok := c.Flusher()
		if !ok {
			return c.Status(500).SendString("Streaming not supported")
		}

		// Send a welcome message
		fmt.Fprintf(c.Response, "event: welcome\n")
		fmt.Fprintf(c.Response, "data: Connected to Quick SSE server\n\n")
		flusher.Flush()

		// Send a status update
		fmt.Fprintf(c.Response, "event: status\n")
		fmt.Fprintf(c.Response, "data: Server is ready\n\n")
		flusher.Flush()

		// Send a final message
		fmt.Fprintf(c.Response, "event: info\n")
		fmt.Fprintf(c.Response, "data: This is a simple SSE example\n\n")
		flusher.Flush()

		return nil
	})

	log.Println("🚀 Quick SSE Server started on :3000")
	log.Println("📡 Test with: curl -N http://localhost:3000/events/simple")
	q.Listen(":3000")
}

// Test with curl:
// curl -N http://localhost:3000/events/simple
//
// Expected output:
// event: welcome
// data: Connected to Quick SSE server
//
// event: status
// data: Server is ready
//
// event: info
// data: This is a simple SSE example

// Test in browser with JavaScript:
// const eventSource = new EventSource('http://localhost:3000/events/simple');
// eventSource.addEventListener('welcome', (e) => console.log('Welcome:', e.data));
// eventSource.addEventListener('status', (e) => console.log('Status:', e.data));
// eventSource.addEventListener('info', (e) => console.log('Info:', e.data));

// test with curl, with multiple calls
// for i in {1..10}; do
//   curl -i -N http://localhost:3000/events/simple
//   sleep 1
// done

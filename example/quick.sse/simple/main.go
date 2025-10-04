package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
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

	q.Get("/events/no", func(c *quick.Ctx) error {
		// Set SSE headers manually
		// c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Access-Control-Allow-Origin", "*")

		// Send a welcome message
		fmt.Fprintf(c.Response, "event: welcome\n")
		fmt.Fprintf(c.Response, "data: Connected to Quick SSE server\n\n")

		c.Status(400)
		c.Status(200)

		// return nil
		// c.Status(400).SendString("Streaming not supported")
		return c.Status(200).SendString("Streaming not supported")
	})

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

	// Alternative using Flush() method - simpler approach
	q.Get("/events/simple-alt", func(c *quick.Ctx) error {
		// Set SSE headers
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		// Send welcome event
		fmt.Fprintf(c.Response, "data: Welcome to Quick SSE!\n\n")
		if err := c.Flush(); err != nil {
			return err
		}

		// Send notification event
		fmt.Fprintf(c.Response, "data: Your connection is established\n\n")
		if err := c.Flush(); err != nil {
			return err
		}

		return nil
	})

	log.Println("🚀 Quick SSE Server started on :3000")
	log.Println("📡 Test with: curl -N http://localhost:3000/events/simple")
	log.Println("📡 Test alt: curl -N http://localhost:3000/events/simple-alt")
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
//   curl -i -N http://localhost:3000/events/no
//   sleep 1
// done

// test with curl, with multiple calls
// for i in {1..10}; do
//   curl -i -N http://localhost:3000/events/simple
//   sleep 1
// done

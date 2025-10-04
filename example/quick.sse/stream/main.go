package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick"
)

// Example of Server-Sent Events (SSE) with Quick - Streaming approach with loop
// This example demonstrates continuous streaming of events to the client.
// Perfect for real-time dashboards, notifications, and live updates.
//
// SSE Format:
// - event: <event-name>  (optional)
// - data: <message>
// - id: <event-id>       (optional)
// - retry: <milliseconds> (optional)
// Each message ends with double newline (\n\n)

func main() {
	q := quick.New()

	// Real-time clock stream - sends current time every second
	q.Get("/events/clock", func(c *quick.Ctx) error {
		// Set SSE headers
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Access-Control-Allow-Origin", "*")

		// Get flusher for immediate data transmission
		flusher, ok := c.Flusher()
		if !ok {
			return c.Status(500).SendString("Streaming not supported")
		}

		// Send events in a loop (30 seconds)
		for i := 0; i < 30; i++ {
			currentTime := time.Now().Format("15:04:05")

			// Send event with custom event name and data
			fmt.Fprintf(c.Response, "event: time\n")
			fmt.Fprintf(c.Response, "data: %s\n", currentTime)
			fmt.Fprintf(c.Response, "id: %d\n\n", i)
			flusher.Flush()

			time.Sleep(1 * time.Second)
		}

		// Send completion event
		fmt.Fprintf(c.Response, "event: done\n")
		fmt.Fprintf(c.Response, "data: Stream completed\n\n")
		flusher.Flush()

		return nil
	})

	// Counter stream - counts from 1 to 10
	q.Get("/events/counter", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		// Using the simplified Flush() method
		for i := 1; i <= 10; i++ {
			msg := fmt.Sprintf("data: Count: %d\n\n", i)
			fmt.Fprint(c.Response, msg)

			if err := c.Flush(); err != nil {
				return err
			}

			time.Sleep(500 * time.Millisecond)
		}

		return nil
	})

	// Progress bar simulation
	q.Get("/events/progress", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Access-Control-Allow-Origin", "*")

		flusher, ok := c.Flusher()
		if !ok {
			return c.Status(500).SendString("Streaming not supported")
		}

		// Simulate a task with progress updates
		for progress := 0; progress <= 100; progress += 10 {
			// Send progress event with JSON data
			fmt.Fprintf(c.Response, "event: progress\n")
			fmt.Fprintf(c.Response, "data: {\"percent\": %d, \"status\": \"processing\"}\n\n", progress)
			flusher.Flush()

			time.Sleep(500 * time.Millisecond)
		}

		// Send completion event
		fmt.Fprintf(c.Response, "event: complete\n")
		fmt.Fprintf(c.Response, "data: {\"percent\": 100, \"status\": \"done\"}\n\n")
		flusher.Flush()

		return nil
	})

	// Notifications stream - sends random notifications
	q.Get("/events/notifications", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		notifications := []string{
			"New user registered",
			"Payment received",
			"Order shipped",
			"System update available",
			"New message received",
		}

		for i, notification := range notifications {
			fmt.Fprintf(c.Response, "event: notification\n")
			fmt.Fprintf(c.Response, "data: %s\n", notification)
			fmt.Fprintf(c.Response, "id: %d\n\n", i+1)

			if err := c.Flush(); err != nil {
				return err
			}

			time.Sleep(2 * time.Second)
		}

		return nil
	})

	log.Println("🚀 Quick SSE Streaming Server started on :3000")
	log.Println("")
	log.Println("📡 Available endpoints:")
	log.Println("   curl -N http://localhost:3000/events/clock")
	log.Println("   curl -N http://localhost:3000/events/counter")
	log.Println("   curl -N http://localhost:3000/events/progress")
	log.Println("   curl -N http://localhost:3000/events/notifications")
	log.Println("")
	q.Listen(":3000")
}

// Test with curl:
// curl -N http://localhost:3000/events/clock
// curl -N http://localhost:3000/events/counter
// curl -N http://localhost:3000/events/progress
// curl -N http://localhost:3000/events/notifications

// Test in browser with JavaScript:
//
// // Clock example
// const clock = new EventSource('http://localhost:3000/events/clock');
// clock.addEventListener('time', (e) => {
//   console.log('Current time:', e.data);
// });
//
// // Progress example
// const progress = new EventSource('http://localhost:3000/events/progress');
// progress.addEventListener('progress', (e) => {
//   const data = JSON.parse(e.data);
//   console.log(`Progress: ${data.percent}% - ${data.status}`);
// });
// progress.addEventListener('complete', (e) => {
//   console.log('Task completed!');
//   progress.close();
// });
//
// // Notifications example
// const notifications = new EventSource('http://localhost:3000/events/notifications');
// notifications.addEventListener('notification', (e) => {
//   console.log('New notification:', e.data);
// });

package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	q := quick.New()

	// Enable Gzip middleware
	// This will automatically compress responses for clients that support Gzip
	q.Use(compress.Gzip())

	// Define a route that returns a compressed JSON response
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		// Setting response headers
		c.Set("Content-Type", "application/json")
		c.Set("Accept-Encoding", "gzip") // Enabling Gzip compression

		// Defining the response structure
		type my struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		// Returning a JSON response with headers
		return c.Status(200).JSON(&my{
			Msg:     "Quick ",
			Headers: c.Headers,
		})
	})

	// Start the HTTP server on port 8080
	// The server will listen for incoming requests at http://localhost:8080
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

//$ curl -X GET http://localhost:8080/v1/compress -H "Accept-Encoding: gzip" --compressed -i

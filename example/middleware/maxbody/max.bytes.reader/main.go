package main

import (
	"io"
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

const maxBodySize = 1024 // 1KB

// Start the server and apply request body size limit middleware
func main() {
	q := quick.New()

	// Define route with extra validation using MaxBytesReader
	q.Post("/v1/user/maxbody/max", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Apply MaxBytesReader to prevent oversized payloads
		c.Request.Body = quick.MaxBytesReader(c.Response, c.Request.Body, maxBodySize)

		// Securely read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			return c.Status(http.StatusRequestEntityTooLarge).String("Request body too large")
		}
		return c.Status(http.StatusOK).Send(body)
	})

	log.Println("Server running at http://0.0.0.0:8080")
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

// $ curl -X POST http://0.0.0.0:8080/v1/user/maxbody/max \
//      -H "Content-Type: application/json" \
//      --data-binary @<(head -c 2048 </dev/zero | tr '\0' 'A')

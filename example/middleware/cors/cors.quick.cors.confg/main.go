package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/cors"
)

func main() {
	// CORS middleware configuration
	q := quick.New()

	q.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"https://httpbin.org"},           // Allow requests only from this domain
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allowed headers
		AllowCredentials: true,                                      // Allow credentials (cookies, auth headers)
		MaxAge:           600,                                       // Cache CORS preflight request for 10 minutes
	}))

	// OPTIONS route to handle preflight CORS requests
	q.Options("/api/data", func(c *quick.Ctx) error {
		return c.Status(204).Send(nil) // Returns an empty response with status 204 (No Content)
	})

	// GET route protected by CORS
	q.Get("/api/data", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(`{"message": "Hello, CORS!"}`)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

// Example cURL requests to test CORS

// Preflight CORS request (OPTIONS method)
// This is sent by browsers before making an actual cross-origin request
//
// curl -X OPTIONS -H "Origin: https://httpbin.org" \
//      -H "Access-Control-Request-Method: GET" \
//      -H "Access-Control-Request-Headers: Content-Type, Authorization" \
//      -v http://localhost:8080/api/data

// Actual GET request with CORS headers
//
// curl -X GET -H "Origin: https://httpbin.org" \
//      -H "Content-Type: application/json" \
//      -H "Authorization: Bearer mytoken123" \
//      -v http://localhost:8080/api/data

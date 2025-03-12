package main

import (
	"encoding/json"
	"time"

	"github.com/jeffotoni/quick"
)

// Struct representing a user model
type My struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Year     int                    `json:"year"`
	Price    float64                `json:"price"`
	Big      bool                   `json:"big"`
	Car      bool                   `json:"car"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Options  []Option               `json:"options"`
	Extra    interface{}            `json:"extra"`
	Dynamic  map[string]interface{} `json:"dynamic"`
}

type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// curl --location 'http://localhost:8080/v1/user' \
// --header 'Content-Type: application/json' \
// --data '[{"id": "123", "name": "Alice", "year": 20, "price": 100.5,
// "big": true, "car": false, "tags": ["fast", "blue"], "metadata": {"brand": "Tesla"},
//
//	"options": [{"key": "color", "value": "red"}],
//
// "extra": "some data", "dynamic": {"speed": "200km/h"}}]'
func main() {
	q := quick.New(quick.Config{
		MaxBodySize:       20 * 1024 * 1024,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}) // Initialize Quick framework

	// Define a POST route at /v1/user
	q.Post("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-type", "application/json")

		bodyBytes := c.Body()

		var my []My
		if err := json.Unmarshal(bodyBytes, &my); err != nil {
			return c.Status(400).JSON(map[string]string{"error": err.Error()})
		}
		// Return the parsed JSON data as a response with 200 OK
		return c.Status(200).Send(bodyBytes)

		// Alternative:
		// return c.Status(200).String(c.BodyString())
		// Return raw request body as a string
	})

	// Start the server and listen on port 8080
	_ = q.Listen("0.0.0.0:8080")
}

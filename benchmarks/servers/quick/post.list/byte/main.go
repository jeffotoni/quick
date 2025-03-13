package main

import (
	"encoding/json"

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
		MaxBodySize: 20 * 1024 * 1024,
		// ReadTimeout:       60 * time.Second,
		// WriteTimeout:      60 * time.Second,
		// IdleTimeout:       120 * time.Second,
		// ReadHeaderTimeout: 10 * time.Second,
	}) // Initialize Quick framework

	// Define a POST route at /v1/user
	q.Post("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-type", "application/json")

		var my []My

		// Read the request body and unmarshal JSON
		err := json.NewDecoder(c.Request.Body).Decode(&my)
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		// Serialize users struct to JSON
		b, err := json.Marshal(my)
		if err != nil {
			return c.Status(500).SendString("Error encoding JSON")
		}

		// Return the serialized JSON as response
		return c.Status(200).Send(b)

		// Alternative:
		// return c.Status(200).String(c.BodyString())
		// Return raw request body as a string
	})

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}

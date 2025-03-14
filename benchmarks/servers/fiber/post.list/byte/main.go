package main

// import (
// 	"bytes"
// 	"encoding/json"

// 	"github.com/gofiber/fiber/v2"
// )

// // Struct representing a user model
// type My struct {
// 	ID       string                 `json:"id"`
// 	Name     string                 `json:"name"`
// 	Year     int                    `json:"year"`
// 	Price    float64                `json:"price"`
// 	Big      bool                   `json:"big"`
// 	Car      bool                   `json:"car"`
// 	Tags     []string               `json:"tags"`
// 	Metadata map[string]interface{} `json:"metadata"`
// 	Options  []Option               `json:"options"`
// 	Extra    interface{}            `json:"extra"`
// 	Dynamic  map[string]interface{} `json:"dynamic"`
// }

// // Struct representing an option with a key-value pair
// type Option struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }

// // curl --location 'http://localhost:8080/v1/user' \
// // --header 'Content-Type: application/json' \
// // --data '[{"id": "123", "name": "Alice", "year": 20, "price": 100.5,
// // "big": true, "car": false, "tags": ["fast", "blue"], "metadata": {"brand": "Tesla"},
// // "options": [{"key": "color", "value": "red"}],
// // "extra": "some data", "dynamic": {"speed": "200km/h"}}]'
// func main() {
// 	// Create a new Fiber instance
// 	app := fiber.New()

// 	// Define a POST route at /v1/user
// 	app.Post("/v1/user", func(c *fiber.Ctx) error {
// 		// Set the Content-Type header to application/json
// 		c.Set("Content-Type", "application/json")

// 		var my []My

// 		// NewReader
// 		bodyReader := bytes.NewReader(c.Body())

// 		// Read the request body and unmarshal JSON
// 		err := json.NewDecoder(bodyReader).Decode(&my)
// 		if err != nil {
// 			return c.Status(500).JSON(map[string]string{"error": err.Error()})
// 		}

// 		// Serialize users struct to JSON
// 		b, err := json.Marshal(my)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Error encoding JSON response")
// 		}

// 		// Return the serialized JSON as response
// 		return c.Send(b)
// 	})

// 	// Start the Fiber server on port 8080
// 	app.Listen(":8080")
// }

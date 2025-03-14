package main

// import (
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

// type Option struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }

// // curl --location 'http://localhost:8080/v1/user' \
// // --header 'Content-Type: application/json' \
// // --data '[{"id": "123", "name": "Alice", "year": 20, "price": 100.5,
// // "big": true, "car": false, "tags": ["fast", "blue"], "metadata": {"brand": "Tesla"}, "options": [{"key": "color", "value": "red"}],
// // "extra": "some data", "dynamic": {"speed": "200km/h"}}]'

// func main() {

// 	app := fiber.New(fiber.Config{
// 		BodyLimit: 20 * 1024 * 1024, // 100MB
// 	})

// 	app.Post("/v1/user", func(c *fiber.Ctx) error {
// 		c.Set("Content-Type", "application/json")
// 		// Parse the request body into the struct
// 		var my []My // Create a variable to store incoming user data

// 		// Parse the request body into the struct
// 		err := c.BodyParser(&my)
// 		if err != nil {
// 			// If parsing fails, return a 400 Bad Request response
// 			return c.Status(400).SendString(err.Error())
// 		}

// 		// Return the parsed JSON data as a response with 200 OK
// 		return c.Status(200).JSON(my)
// 	})

// 	app.Listen(":8080")
// }

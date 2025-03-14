package main

// import (
// 	"github.com/gofiber/fiber/v2"
// )

// // Struct representing a user model
// type My struct {
// 	Name string `json:"name"` // User's name
// 	Year int    `json:"year"` // User's birth year
// }

// func main() {
// 	app := fiber.New()

// 	app.Post("/v1/user", func(c *fiber.Ctx) error {
// 		c.Set("Content-Type", "application/json")
// 		// Parse the request body into the struct
// 		var my My // Create a variable to store incoming user data

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

//curl --location 'http://localhost:8080/v1/user' \
// --header 'Content-Type: application/json' \
// --data '{"name": "Alice", "year": 20}'

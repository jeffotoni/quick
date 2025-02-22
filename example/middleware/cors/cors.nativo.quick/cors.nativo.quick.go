package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/goquick/middleware/cors"
	"github.com/jeffotoni/quick"
)

func main() {
	app := quick.New()

	app.Use(cors.New(cors.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
	}), "cors")

	app.Post("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		type My struct {
			Name string `json:"name"`
			Year int    `json:"year"`
		}

		var my My
		err := c.BodyParser(&my)
		fmt.Println("byte:", c.Body())

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		fmt.Println("String:", c.BodyString())
		return c.Status(200).JSON(&my)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

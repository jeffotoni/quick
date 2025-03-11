package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/cors"
)

func main() {
	app := quick.New()

	app.Use(cors.New(cors.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
	}), "cors")

	app.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		my := struct {
			Name string `json:"name"`
			Year int    `json:"year"`
		}{
			Name: "Teste",
			Year: 2024,
		}

		fmt.Println("Enviando resposta:", my)
		return c.Status(200).JSON(my)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/cors"
	// cors "github.com/rs/cors"
)

// / Example cURL to test
// curl -X OPTIONS -v http://localhost:8080/v1/user
//
//	curl -X OPTIONS -H "Origin: http://localhost:3000/" \
//	 -H "Access-Control-Request-Method: POST" -v \
//	 http://localhost:8080/v1/user
func main() {
	q := quick.New()

	q.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	}))

	q.Post("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		type My struct {
			Name string `json:"name"`
			Year string `json:"year"`
		}

		var my My
		err := c.BodyParser(&my)
		fmt.Println("byte:", string(c.Body()))

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		fmt.Println("String:", c.BodyString())
		return c.Status(200).JSON(my)
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}

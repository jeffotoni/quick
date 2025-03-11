// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/jeffotoni/quick"
// 	//"github.com/jeffotoni/quick/middleware/cors"
// 	cors "github.com/rs/cors"
// )

// func main() {
// 	q := quick.New()

// 	q.Use(cors.New(cors.Options{
// 		AllowedOrigins:   []string{"*"}, // Allow any origin
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowedHeaders:   []string{"*"}, // Allow any header
// 		ExposedHeaders:   []string{"*"}, // Show any header
// 		AllowCredentials: true,          // Allow cookies and authentication via CORS
// 		Debug:            true,
// 	}))

// 	q.Post("/v1/user", func(c *quick.Ctx) error {
// 		c.Set("Content-Type", "application/json")
// 		type My struct {
// 			Name string `json:"name"`
// 			Year int    `json:"year"`
// 		}

// 		var my My
// 		err := c.BodyParser(&my)
// 		fmt.Println("byte:", c.Body())

// 		if err != nil {
// 			return c.Status(400).SendString(err.Error())
// 		}

// 		fmt.Println("String:", c.BodyString())
// 		return c.Status(200).JSON(&my)
// 	})

// 	log.Fatal(q.Listen("0.0.0.0:8080"))
// }

package main

import (
	"net/http"

	"github.com/rs/cors"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	http.ListenAndServe(":8080", c.Handler(handler))
}

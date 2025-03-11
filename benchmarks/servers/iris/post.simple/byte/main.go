package main

import (
	"log"
	"net/http"

	"github.com/kataras/iris/v12"
)

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {
	i := iris.New()

	i.Post("/v1/user", func(ctx iris.Context) {

		if _, err := ctx.GetBody(); err != nil {
			ctx.JSON(iris.Map{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK)

	})

	log.Fatal(i.Listen(":8080"))
}

//curl --location 'http://localhost:8080/v1/user' \
// --header 'Content-Type: application/json' \
// --data '{"name": "Alice", "year": 20}'

package main

// import (
// 	"log"

// 	"github.com/kataras/iris/v12"
// )

// // Struct representing a user model
// type My struct {
// 	Name string `json:"name"` // User's name
// 	Year int    `json:"year"` // User's birth year
// }

// // $ curl --location 'http://localhost:8080/v1/user' \
// // --header 'Content-Type: application/json' \
// // --data '{"name": "Alice", "year": 20}'
// func main() {
// 	i := iris.New()

// 	i.Post("/v1/user", func(ctx iris.Context) {
// 		var my My
// 		if err := ctx.ReadJSON(&my); err != nil {
// 			ctx.StatusCode(iris.StatusBadRequest)
// 			ctx.JSON(iris.Map{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(my)
// 	})

// 	log.Fatal(i.Listen(":8080"))
// }

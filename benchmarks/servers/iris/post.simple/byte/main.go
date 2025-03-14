package main

// import (
// 	"log"

// 	"github.com/kataras/iris/v12"
// )

// // $ curl --location 'http://localhost:8080/v1/user' \
// // --header 'Content-Type: application/json' \
// // --data '{"name": "Alice", "year": 20}'
// func main() {
// 	app := iris.New()

// 	app.Post("/v1/user", func(ctx iris.Context) {
// 		// Capture the request body as `[]byte`
// 		bodyBytes, err := ctx.GetBody()
// 		if err != nil {
// 			ctx.StatusCode(iris.StatusInternalServerError)
// 			ctx.WriteString("Error reading body")
// 			return
// 		}

// 		// Returns the same bytes received in the response
// 		ctx.ContentType("application/json")
// 		ctx.Write(bodyBytes)
// 	})

// 	log.Fatal(app.Listen(":8080"))
// }

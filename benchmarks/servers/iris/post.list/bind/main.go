package main

// import (
// 	"encoding/json"
// 	"log"

// 	"github.com/kataras/iris/v12"
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
// // "big": true, "car": false, "tags": ["fast", "blue"], "metadata": {"brand": "Tesla"},
// //
// //	"options": [{"key": "color", "value": "red"}],
// //
// // "extra": "some data", "dynamic": {"speed": "200km/h"}}]'
// func main() {
// 	app := iris.New()

// 	app.Post("/v1/user", func(ctx iris.Context) {
// 		ctx.ContentType("application/json")

// 		// Capture the request body as `[]byte`
// 		bodyBytes, err := ctx.GetBody()
// 		if err != nil {
// 			ctx.StatusCode(iris.StatusInternalServerError)
// 			ctx.WriteString("Error reading body")
// 			return
// 		}

// 		// Deserialize the JSON into a struct
// 		var users []My
// 		if err := json.Unmarshal(bodyBytes, &users); err != nil {
// 			ctx.StatusCode(iris.StatusBadRequest)
// 			ctx.JSON(iris.Map{"error": "Invalid JSON format"})
// 			return
// 		}

// 		// Serialize back to JSON
// 		responseBytes, err := json.Marshal(users)
// 		if err != nil {
// 			ctx.StatusCode(iris.StatusInternalServerError)
// 			ctx.JSON(iris.Map{"error": "Error encoding JSON response"})
// 			return
// 		}

// 		// Return the received JSON
// 		ctx.Write(responseBytes)
// 	})

// 	log.Fatal(app.Listen(":8080"))
// }

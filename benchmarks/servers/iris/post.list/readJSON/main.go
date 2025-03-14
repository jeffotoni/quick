package main

// import (
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
// 	i := iris.New()

// 	i.Post("/v1/user", func(ctx iris.Context) {
// 		var my []My
// 		if err := ctx.ReadJSON(&my); err != nil {
// 			ctx.StatusCode(iris.StatusBadRequest)
// 			ctx.JSON(iris.Map{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(my)
// 	})

// 	log.Fatal(i.Listen(":8080"))
// }

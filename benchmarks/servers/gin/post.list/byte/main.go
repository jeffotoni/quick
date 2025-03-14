package main

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
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
// // --data '[{"id": "123", "name": "Alice", "year": 20,
// // "price": 100.5, "big": true, "car": false, "tags": ["fast", "blue"],
// // "metadata": {"brand": "Tesla"}, "options": [{"key": "color", "value": "red"}],
// // "extra": "some data", "dynamic": {"speed": "200km/h"}}]'
// func main() {

// 	gin.SetMode(gin.ReleaseMode)

// 	r := gin.New()

// 	r.POST("/v1/user", func(c *gin.Context) {
// 		c.Set("Content-Type", "application/json")

// 		var my []My // Define a slice to store multiple user objects

// 		err := json.NewDecoder(c.Request.Body).Decode(&my)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Serialize users struct to JSON
// 		b, err := json.Marshal(my)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		}

// 		// Return the parsed JSON data as a response with 200 OK
// 		// c.JSON(http.StatusOK, my)
// 		c.Data(http.StatusOK, "text/plain", b)
// 	})

// 	r.Run(":8080")
// }

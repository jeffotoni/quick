package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Struct representing a user model
type My struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Year     int                    `json:"year"`
	Price    float64                `json:"price"`
	Big      bool                   `json:"big"`
	Car      bool                   `json:"car"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Options  []Option               `json:"options"`
	Extra    interface{}            `json:"extra"`
	Dynamic  map[string]interface{} `json:"dynamic"`
}

// Struct representing an option with a key-value pair
type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//curl --location 'http://localhost:8080/v1/user' \
//--header 'Content-Type: application/json' \
//--data '[{"id": "123", "name": "Alice", "year": 20, "price": 100.5, "big": true, "car": false, "tags": ["fast", "blue"], "metadata": {"brand": "Tesla"}, "options": [{"key": "color", "value": "red"}], "extra": "some data", "dynamic": {"speed": "200km/h"}}]'

func main() {
	e := echo.New()

	// Define a POST route at /v1/user

	e.POST("/v1/user", func(c echo.Context) error {
		c.Set("Content-type", "application/json")

		var my []My

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Erro ao ler o body: "+err.Error())
		}

		if err := json.Unmarshal(body, &my); err != nil {
			return c.String(http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
		}

		return c.JSON(http.StatusOK, my)
	})

	e.Start(":8080")
}

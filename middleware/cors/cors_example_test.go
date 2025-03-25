package cors

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleNew()
//
//	it with the Examples type.
func ExampleNew() {
	q := quick.New()

	q.Use(New(Config{
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

	// Send test request using Quick's built-in test utility
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodPost,
		URI:     "/v1/user",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"name":"Alice","year":"2024"}`),
	})

	fmt.Println("Response Body:", string(resp.Body()))

	// Output
	// byte: {"name":"Alice","year":"2024"}
	// String: {"name":"Alice","year":"2024"}
	// Response Body: {"name":"Alice","year":"2024"}
}

// This function is named ExampleNew_allowedOrigin()
//
//	it with the Examples type.
func ExampleNew_allowedOrigin() {
	q := quick.New()

	// Configure or middleware CORS allowing all the origins
	q.Use(New(Config{
		AllowedOrigins: []string{"https://httpbin.org"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Origin", "Content-Type"},
		Debug:          true,
	}))

	q.Get("/v1/hello", func(c *quick.Ctx) error {
		return c.Status(200).String("Hello from CORS!")
	})

	// Simulates a request for an allowed domain
	res, _ := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodGet,
		URI:     "/v1/hello",
		Headers: map[string]string{"Origin": "https://httpbin.org"},
	})

	// Print the answer
	fmt.Println(res.StatusCode()) // Expected: 200
	fmt.Println(res.BodyStr())    // Expected: "Hello from CORS!"

	// Output:
	// 200
	// Hello from CORS!
}

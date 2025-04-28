
## 🌐 CORS (Cross-Origin Resource Sharing)
Controls how your API can be accessed from different domains.

- Restricts which domains, methods, and headers are allowed.
- Helps prevent CORS errors in browsers.
- Configurable via allowed origins, headers, and credentials.


#### 🔧 CORS Example with Quick
The example below configures CORS to allow requests from any origin, method, and header.

```go
package main

import (
    "fmt"
    "log"
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cors"
)

func main() {
    // Create a new Quick instance
    app := quick.New()

    // Apply CORS middleware to allow all origins, methods, and headers
    app.Use(cors.New(cors.Config{
        AllowedOrigins: []string{"*"}, // Allows requests from any origin
        AllowedMethods: []string{"*"}, // Allows all HTTP methods (GET, POST, PUT, DELETE, etc.)
        AllowedHeaders: []string{"*"}, // Allows all headers
    }))

    // Define a POST route for creating a user
    app.Post("/v1/user", func(c *quick.Ctx) error {
        // Set response content type as JSON
        c.Set("Content-Type", "application/json")

        // Define a struct to hold incoming JSON data
        type My struct {
            Name string `json:"name"`
            Year int    `json:"year"`
        }

        var my My

        // Parse the request body into the struct
        err := c.BodyParser(&my)
        fmt.Println("byte:", c.Body()) // Print raw request body

        if err != nil {
            // Return a 400 Bad Request if parsing fails
            return c.Status(400).SendString(err.Error())
        }

        // Print the request body as a string
        fmt.Println("String:", c.BodyString())

        // Return the parsed JSON data with a 200 OK status
        return c.Status(200).JSON(&my)
    })

    // Start the server on port 8080
    log.Fatal(app.Listen("0.0.0.0:8080"))
}

```
---
### 📌 Configuring HTTP server 

```go
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
```

### 📌 Testing with cURL

#### 🔹 Making a POST request with CORS enabled

```go
$ curl --location --request POST 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json' \
--data '{"name": "John Doe", "year": 2024}'
```
---



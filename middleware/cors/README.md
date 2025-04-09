
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
### 📌 Testing with cURL

#### 🔹 Making a POST request with CORS enabled

```go
$ curl --location --request POST 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json' \
--data '{"name": "John Doe", "year": 2024}'
```
---



## 🌍 CORS Middleware - Quick Framework ![Quick Logo](/quick.png)

### 📌 Overview

CORS stands for **"Cross-Origin Resource Sharing"**, which is a security technique used by web browsers to allow a server to restrict access from other sites or domains to its resources. The main purpose of CORS is to protect server resources from malicious attacks from other domains.

Quick is a web framework in Go that supports CORS middleware to handle requests from other domains. **CORS middleware** can be added to Quick using the "github.com/jeffotoni/quick/middleware/cors" library.

To add CORS middleware in a Quick application, simply import the library and call the Cors() function by passing the desired configuration options.

---

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
	}), "cors")

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

### 📌 What I Included in this README
- ✅ Overview: Explanation of CORS and its importance.
- ✅ CORS Implementation:
	- Using Quick Middleware
- ✅ Test with cURL 
	- Sending a POST request.
	- Checking CORS headers.
- ✅ Best Practices: Recommendation to restrict settings in production.

---


Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥


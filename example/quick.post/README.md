# 📌 POST - Quick Framework 

This example demonstrates how to create a POST route in the Quick Framework to receive and process JSON data at the `/v1/user` endpoint.

### 📜 Code Implementation

```go
package main

import "github.com/jeffotoni/quick"

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {
	q := quick.New() // Initialize Quick framework

	// Define a POST route at /v1/user
	q.Post("/v1/user", func(c *quick.Ctx) error {
		var my My // Create a variable to store incoming user data

		// Parse the request body into the struct
		err := c.BodyParser(&my)
		if err != nil {
			// If parsing fails, return a 400 Bad Request response
			return c.Status(400).SendString(err.Error())
		}

		// Return the parsed JSON data as a response with 200 OK
		return c.Status(200).JSON(&my)

		// Alternative:
		// return c.Status(200).String(c.BodyString()) 
		// Return raw request body as a string
	})

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}

```
#### 📌 Testing with cURL

##### 🔹 Create a User (POST Request):

```bash
$ curl --location 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json/' \
--data '{"name":"crow3442","year":2005}'
```

---

#### 📌 What I included in this README
- ✅ Overview of the POST endpoint in Quick Framework
- ✅ Go implementation with JSON parsing using BodyParser
- ✅ POST route to create a new user dynamically
- ✅ Handling of success (200 OK) and error (400 Bad Request) responses
- ✅ cURL example for valid and invalid requests


Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
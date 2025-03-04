# 🚀 Quick Framework - Starting the Server

**Start** is a common term in web applications that refers to starting the web server and making it ready to handle HTTP requests.

In the Quick Framework, the server starts when you:

1. Create an instance of Quick.
2. Define your routes (GET, POST, etc.).
3. Call `Listen()`, specifying the port for the server.

---

### 📜 Code Example - Starting Quick Server
```go
package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New() // Initialize Quick Framework

	// Define a route for "/v1/user"
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json") // Set response type
		return c.Status(200).SendString("Quick in action com CORS ❤️!") // Response message
	})

	// Start the server on port 8080
	q.Listen("0.0.0.0:8080")
}
```

#### 📌 Testing with cURL

##### 🔹 Basic GET Request to Check If the Server is Running:

```go
$ curl --location --request GET "http://localhost:8080/v1/user" \
--header "Content-Type: application/json"
```
---

### ⚡ Expanding the Server with Multiple Routes
You can add more routes to handle different API endpoints:

```go
package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New() // Initializes the Quick framework

	// Sets a GET route for "/v1/user"
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")                       // Sets the content type as JSON
		return c.Status(200).SendString("Quick in action with Cors❤️!") // Returns a success message
	})

	// Sets a GET route for "/v2"
	q.Get("/v2", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")         // Sets the content type as JSON
		return c.Status(200).SendString("Is in the air!") // Returns a message indicating that the service is active
	})

	// Sets a GET route for "/v3"
	q.Get("/v3", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")   // Sets the content type as JSON
		return c.Status(200).SendString("Running!") // Returns a message confirming that the server is running
	})

	// Starts the server on port 8080, allowing connections from any IP
	q.Listen("0.0.0.0.0.0:8080")
}


```

#### 📌 Testing with cURL

##### 🔹 Basic GET Request to Check If the Server is Running:

```bash
# Test endpoint v1/user
curl --location --request GET "http://localhost:8080/v1/user" \
--header "Content-Type: application/json"
```

```bash
# Test endpoint v2
curl --location --request GET "http://localhost:8080/v2" \
--header "Content-Type: application/json"
```

```bash
# Test endpoint v3
curl --location --request GET "http://localhost:8080/v3" \
--header "Content-Type: application/json"
```

#### 📌 What I included in this README
- ✅ Explanation of how to start a Quick server.
- ✅ Minimal example for starting a server with a simple route.
- ✅ Expanded example with multiple API routes.
- ✅ cURL commands to test different endpoints.


Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
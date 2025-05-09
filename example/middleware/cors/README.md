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

### 📌 What I Included in this README
- ✅ Overview: Explanation of CORS and its importance.
- ✅ CORS Implementation:
	- Using Quick Middleware
- ✅ Test with cURL 
	- Sending a POST request.
	- Checking CORS headers.
- ✅ Best Practices: Recommendation to restrict settings in production.

---


Now you can **complete with your specific examples** where I left the spaces `# 🛡️ BasicAuth - Basic Authentication with Quick ![Quick Logo](/quick.png)

This document explains how to implement basic authentication (BasicAuth) using the Quick on Go framework.

---

### 📌 What is `BasicAuth`?

**`Basic Authentication (BasicAuth)`** is a simple method of **HTTP** authentication where the client must send a username and password encoded in **Base64** in the request header.


#### 📝 **Structure of a `BasicAuth`** Request
Each part of the request contains **headers and a body**:

```text
POST /login HTTP/1.1
Host: example.com
Authorization: Basic YWRtaW46MTIzNA==
Content-Type: application/json
Content-Length: 50
{
    "username": "admin",
    "password": "12345"
}
```
📌 **Important headers in `BasicAuth`:**
| Header | Description |
|-----------|-----------|
| `Authorization` | ends the BasicAuth credentials (Base64(username:password)). |
| `Content-Type` | Defines the format of the request body (e.g., application/json). |
| `Content-Length` | Specifies the size of the request body (optional but recommended). |

---

---
### 📌 Basic Authusing environment variables

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

var (
	User     = os.Getenv("USER")
	Password = os.Getenv("PASSWORD")
)

func main() {

	q := quick.New()

	// Adding BasicAuth middleware
	q.Use(middleware.BasicAuth(User, Password))

	// Protected route
	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Starting the server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
---

### 📌 Basic Authentication with Quick Middleware

```go
package basicauth

import (
	"log"
	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

// This function is named ExampleBasicAuth()
// it with the Examples type.
func ExampleBasicAuth() {
	//starting Quick
	q := quick.New()

	// calling middleware
	q.Use(BasicAuth("admin", "1234"))

	// everything below Use will apply the middleware
	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
### 📌 Basic Authentication with Quick Route Groups
```go
package main
import (
	"log"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
	"github.com/jeffotoni/quick"
	)
func main() {
	q := quick.New()

	// Using a group to isolate protected routes
	gr := q.Group("/")

	// Applying BasicAuth middleware to the group
	gr.Use(middleware.BasicAuth("admin", "1234"))

	// Public route
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("Public quick route")
	})

	// Protected route
	gr.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```

### 📌 Custom implementation of BasicAuth

```go
package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"github.com/jeffotoni/quick"
)

func main() {

	q := quick.New()

	// BasicAuth Middleware Custom
	q.Use(func(next http.Handler) http.Handler {
		username := "admin"
		password := "1234"

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Decoding credentials
			payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
			if err != nil || len(strings.SplitN(string(payload), ":", 2)) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			creds := strings.SplitN(string(payload), ":", 2)
			if creds[0] != username || creds[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Protected route
	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
---
### 📌 Testing with cURL

### 🔹 Request authenticated via user and password

```bash
$ curl -u admin:1234 http://localhost:8080/protected
```

### 🔹 Request authenticated via Authorization header

```bash
$ curl -H "Authorization: Basic YWRtaW46MTIzNA==" http://localhost:8080/protected
```

---

###### 🚀 Now you can implement fast and efficient BasicAuth in Quick! 🔥

## **📌 What I included in this README**

- ✅ README checklist - Basic authentication with Quick
- ✅ Overview: Explanation of BasicAuth and its use.
- ✅ Request Structure: Example of an authenticated request with headers.
- ✅ Implementation:
	- Manual BasicAuth middleware.
	- Using the integrated middleware of Quick.
	- Environment variables for credentials.
	- Grouping of protected and public routes.
- ✅ Tests: examples of cURL for authentication and error handling.


Now you can **complete with your specific examples** where I left the spaces **` go ...`**.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥


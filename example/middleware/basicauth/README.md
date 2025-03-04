# 🛡️ BasicAuth - Basic Authentication with Quick ![Quick Logo](/quick.png)

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
    "password": "1234"
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
### 📌 Manual implementation of BasicAuth

```go
func main() {
	q := quick.New()

	// BasicAuth Middleware manual
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
	q := quick.New()
	// Applying BasicAuth middleware
	q.Use(middleware.BasicAuth("admin", "1234"))

	// All routes below `Use` will require authentication
	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
```
### 📌 Basic Authentication with Quick Route Groups
```go

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
### 📌 Testing with cURL

### 🔹 Request authenticated via user and password

```bash
curl -u admin:1234 http://localhost:8080/protected
```

### 🔹 Request authenticated via Authorization header

```bash
curl -H "Authorization: Basic YWRtaW46MTIzNA==" http://localhost:8080/protected
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


Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

## ⚡ Dynamic Middleware - Quick Framework ![Quick Logo](/quick.png)

The **Dynamic Middleware** in Quick allows you to handle requests dynamically based on **headers, parameters**, or **any other custom validation rule**. This is useful for restricting access, modifying responses or adding specific logic before processing a request.

---

### 🚀 How does the example below work?
- Intercepts all requests and verifies the existence of the header Block.
- If the header Block is missing, returns a 400 Bad Request error.
- If Block is true, returns a lock message.
- If Block is false, it allows the request to continue normally for the corresponding route.

```go
package main

import (
	"net/http"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Block") == "" {
				w.WriteHeader(400)
				w.Write([]byte("Your call will not continue, I need to set your Header with Block: false to pass"))
				return
			}

			if r.Header.Get("Block") == "true" {
				w.WriteHeader(200)
				w.Write([]byte("Your messgem is locked, set false in your Block parameter"))
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	q.Get("/greet/:name", func(c *quick.Ctx) error {
		name := c.Param("name")
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Hello " + name + "!")
	})

	q.Listen("0.0.0.0:8080")
}
```


#### 📌 Testing with cURL

##### 🔹Request without the Header Block (Error 400)

```bash
$ curl --location --request GET 'http://localhost:8080/greet/:name' \
--header 'Content-Type: application/json/' \
--data 'Your call will not continue, need to set your Header with Block: false to pass'
```

##### 🔹Request with Block: true (Message blocked)

```bash
$ curl --location --request GET 'http://localhost:8080/greet/Ana' \
--header 'Content-Type: application/json' \
--header 'Block: true'
```

##### 🔹Request with Block: false (Success)

```bash
$ curl --location --request GET 'http://localhost:8080/greet/Ana' \
--header 'Content-Type: application/json' \
--header 'Block: false'
```
---

### 🔐 BasicAuth as a Dynamic Middleware

BasicAuth is a type of Dynamic Middleware because it:

- Intercepts the request before going to the main logic.
- Checks if the Authorization header is present.
- Decodes and validates credentials.
- Allows the request to continue only if authentication is successful.



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

	// implementing middleware directly in Use
	q.Use(func(next http.Handler) http.Handler {
		// credentials
		username := "admin"
		password := "1234"

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if it starts with "Basic"
			if !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Decode credentials
			payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			creds := strings.SplitN(string(payload), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))

}
```
#### 📌 Testing with cURL

##### 🔹 Accessing route protected with username and password

```bash
$ curl --location --request GET 'http://localhost:8080/protected'
-user admin:1234
```

##### 🔹Accessing with header Authorization manual
```bash
$ curl --location --request GET 'http://localhost:8080/protected'
-header 'Authorization: Basic YWRtaW46MTIzNA=='
```

##### 🔹Trying to log in without credentials (should fail)
```bash
$ curl --location --request GET 'http://localhost:8080/protected'
```

---

## ✅ Advantages of Dynamic Middleware

| 🔹 Benefit   | ✅ Description |
|----------------|------------|
| 📂 **Flexibility** | Allows you to modify requests before reaching the main logic. |
| 🔄 **Reuse**  | Can be applied globally for multiple routes. |
| 🔒 **Security**   | Allows restrictions such as authentication, header validation and permissions. |
| ⚡ **Performance**   | Middleware processed before routing, avoiding unnecessary runs. |

---

#### 📌 What I included in this README
- ✅ README checklist - Dynamic Middleware with Quick
- ✅ Overview: Explanation of Dynamic Middleware and its purpose.
- ✅ Request Structure: How the middleware processes requests based on headers.
- ✅ Implementation:
	- Custom middleware intercepting requests.
	- Header-based validation (Block: true/false).
	- Example of a protected route.
	- Connection between Dynamic Middleware and BasicAuth.
- ✅ Tests:
	- cURL examples for blocked, allowed, and unauthorized requests.
	- Authentication using BasicAuth with credentials and headers.
- ✅ Advantages: Table listing flexibility, security, performance, and reusability benefits.


---

Now you can **complete with your specific examples** where I left the spaces # 🛡️ BasicAuth - Basic Authentication with Quick ![Quick Logo](/quick.png)

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

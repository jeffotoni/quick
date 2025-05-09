## 📦 Compression Middleware (compress) - Quick Framework ![Quick Logo](/quick.png)

The **`Compression Middleware`** in Quick provides GZIP compression for HTTP responses, reducing the amount of data sent over the network. It helps to improve performance and bandwidth efficiency, especially for text-based content like JSON, HTML, and CSS.

---

#### 🚀 How It Works

When a client sends a request with the header Accept-Encoding: gzip, the middleware automatically compresses the response. This results in faster load times and reduced bandwidth usage.

#### 📌 Key Features

- ✅ Automatic GZIP compression for compatible clients
- ✅ Improves performance by reducing response size
- ✅ Saves bandwidth and enhances user experience
- ✅ Works seamlessly with Quick’s request-handling flow



#### 🔹 Using Quick Default Middleware (quick.Handler)
```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	q := quick.New()

	// Enable GZIP compression
	q.Use(compress.Gzip())

	// Define a compressed response route
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("Accept-Encoding", "gzip")

		type response struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		return c.Status(200).JSON(&response{
			Msg:     "Quick ❤️",
			Headers: c.Headers,
		})
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
#### 🔹 Using Quick HandlerFunc Middleware (quick.HandlerFunc)
```go
package main

import (
	"log"
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	q := quick.New()

	// Enable GZIP middleware using HandlerFunc version
	q.Use(compress.Gzip())

	// Define a compressed response route
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type response struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		return c.Status(200).JSON(&response{
			Msg:     "Quick ❤️",
			Headers: c.Headers,
		})
	})

	log.Fatal(q.Listen(":8080"))
}
```
#### 🔹 Using Pure net/http Middleware
```go
package main

import (
	"log"
	"net/http"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	mux := http.NewServeMux()

	// Route with compression enabled using the middleware
	mux.Handle("/v1/compress", compress.Gzip()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, net/http with Gzip!"}`))
	})))

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```


#### 📌 Testing with cURL

##### 🔹Request Without GZIP (Uncompressed Response):
```bash
$ curl -X GET http://localhost:8080/v1/compress
```
##### 🔹Request With GZIP:
```bash
$ curl -X GET http://localhost:8080/v1/compress -H "Accept-Encoding: gzip" --compressed
```

#### 🔍 Why Use GZIP Compression?  

| Feature                     | Benefit                                              |
|-----------------------------|------------------------------------------------------|
| 🚀 **Faster Load Times**     | Reduces response sizes, improving website speed.    |
| 💾 **Bandwidth Optimization** | Saves data usage, especially on mobile networks.   |
| 🎯 **Better User Experience** | Users receive responses faster, improving performance. |
| 🔄 **Seamless Integration**  | Works automatically when a client supports GZIP.   |


#### 🔧 When to Use GZIP?
- ✅ When serving JSON, HTML, CSS, JS, or plain text
- ❌ Avoid compressing already compressed content (e.g., images, videos, ZIP files)


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
.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

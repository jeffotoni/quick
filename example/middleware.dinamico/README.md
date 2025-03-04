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

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

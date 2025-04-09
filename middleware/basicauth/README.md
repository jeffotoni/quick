# 📂 Middleware - Quick Framework ![Quick Logo](/quick.png)



### 📝 **Structure of a `BasicAuth`** Request
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
### 📌 **Important headers in `BasicAuth`:**
| Header | Description |
|-----------|-----------|
| `Authorization` | ends the BasicAuth credentials (Base64(username:password)). |
| `Content-Type` | Defines the format of the request body (e.g., application/json). |
| `Content-Length` | Specifies the size of the request body (optional but recommended). |


---

## 📜 Middlewares Available

🔐 BasicAuth
Provides HTTP Basic Authentication, requiring a username and password to access protected routes.

- Can be applied globally or to specific routes.
- Supports authentication via environment variables.
- Allows custom implementation if needed.

### Basic Auth environment variables

This example sets up Basic Authentication using environment variables to store the credentials securely.
the routes below are affected, to isolate the route use group to apply only to routes in the group.

```go
package main

import (
    "log"
    "os"

    "github.com/jeffotoni/quick"
    middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

// Environment variables for authentication
// export USER=admin
// export PASSWORD=1234

var (
    // Retrieve the username and password from environment variables
    User     = os.Getenv("USER")
    Password = os.Getenv("PASSORD")
)

func main() {

    // Initialize a new Quick instance
    q := quick.New()

    // Apply Basic Authentication middleware
    q.Use(middleware.BasicAuth(User, Password))

    // Define a protected route
    q.Get("/protected", func(c *quick.Ctx) error {
        // Set the response content type to JSON
        c.Set("Content-Type", "application/json")

        // Return a success message
        return c.SendString("You have accessed a protected route!")
    })

    // Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/api/v1/users/123

You have accessed a protected route!
```
---

### Basic Authentication with Quick Middleware

This example uses the built-in BasicAuth middleware provided by Quick, offering a simple authentication setup.

```go
package main

import (
    "log"

    "github.com/jeffotoni/quick"
    middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

func main() {

    //starting Quick
    q := quick.New()

    // calling middleware
    q.Use(middleware.BasicAuth("admin", "1234"))

    // everything below Use will apply the middleware
    q.Get("/protected", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.SendString("You have accessed a protected route!")
    })

    // Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL
```bash
$ curl -i -X GET 'http://localhost:8080/protected' \
--header 'Authorization: Basic YWRtaW46MTIzNA=='

You have accessed a protected route!
```

### Basic Authentication with Quick Route Groups

This example shows how to apply Basic Authentication to a specific group of routes using Quick's Group functionality.
When we use group we can isolate the middleware, this works for any middleware in quick.

```go
package main

import (
    "log"

    "github.com/jeffotoni/quick"
    middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

func main() {
    //starting Quick
    q := quick.New()

    // using group to isolate routes and middlewares
    gr := q.Group("/")

    // middleware BasicAuth
    gr.Use(middleware.BasicAuth("admin", "1234"))

    // route public
    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.SendString("Public quick route")
    })

    // protected route
    gr.Get("/protected", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.SendString("You have accessed a protected route!")
    })

    // Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/v1/user

Public quick route
```


### BasicAuth Customized

This example shows a custom implementation of Basic Authentication without using any middleware. It manually verifies user credentials and applies authentication to protected routes.

In quick you are allowed to make your own custom implementation directly in q.Use(..), that is, you will be able to implement it directly if you wish.

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
    //starting Quick
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

    // Start the server on port 8080
    q.Listen("0.0.0.0:8080")

}
```
### 📌 cURL
```bash
$ curl -i -u admin:1234 -X GET http://localhost:8080/protected

You have accessed a protected route!
```
---

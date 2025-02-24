![Logo do Quick](/quick_logo.png)

🛡️ BasicAuth - Basic Authentication with Quick
This document explains how to implement basic authentication (BasicAuth) using the Quick on Go framework.

:memo: Implementation Examples

:round_pushpin: Example 1: Manual implementation of BasicAuth

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
:hammer_and_wrench: Testing an API with cURL
### Request authenticated with user and password

```bash
$ curl -u admin:1234 http://localhost:8080/protected
```
###  Or sending the Authorization header manually
```bash
$ curl -H "Authorization: Basic YWRtaW46MTIzNA==" http://localhost:8080/protected
```
:white_check_mark: Expected response (200 OK)
```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 37
{"message": "You have accessed a protected route"}

```


:round_pushpin: Example 2: Basic Authusing environment variables

```go
package main

import (
	"log"
	"os"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

// export USER=admin
// export PASSWORD=1234

// Obtaining credentials from the environment
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

:hammer_and_wrench: Testing an API with cURL
1. Requisição autenticada via usuário e senha

```bash
curl -i -XGET -u admin:1234 http://localhost:8080/protected
```
:white_check_mark: Expected response (200 OK)
```bash
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 37

{"message": "You have accessed a protected route"}

```

2. Requisição autenticada via cabeçalho Authorization

```bash
curl -i -XGET -H "Authorization: Basic YWRtaW46MTIzNA==" http://localhost:8080/protected
```
:white_check_mark: Expected response (200 OK)
```bash
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 37

{"message": "You have accessed a protected route"}

```
3. Tentar acessar sem autenticação

```bash
curl -i -XGET http://localhost:8080/protected
```
:x: Expected response
```bash
HTTP/1.1 401 Unauthorized
Content-Length: 12

Unauthorized
```
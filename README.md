
![Logo do Quick](/quick_logo.png)


[![GoDoc](https://godoc.org/github.com/jeffotoni/quick?status.svg)](https://godoc.org/github.com/jeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jeffotoni/quick/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jeffotoni/quick/tree/main) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick) [![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/jeffotoni/quick/main)  ![Coveralls](https://img.shields.io/coverallsCoverage/github/jeffotoni/quick) ![GitHub contributors](https://img.shields.io/github/contributors/jeffotoni/quick)
![GitHub stars](https://img.shields.io/github/last-commit/jeffotoni/quick) ![GitHub stars](https://img.shields.io/github/forks/jeffotoni/quick?style=social) ![GitHub stars](https://img.shields.io/github/stars/jeffotoni/quick)

<!-- [![Github Release](https://img.shields.io/github/v/release/jeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/jeffotoni/quick) -->

<h2 align="center">
    <p>
         <a href="README.md">English</a> |
          <a href="README.pt-br.md">Рortuguês</a>
    </p> 
</h2>


# Quick - a lightweight router for go ![Quick Logo](/quick.png)

🚀 Quick is a **flexible and extensible** route manager for the Go language. Its goal is to be **fast and high-performance**, as well as being **100% compatible with net/http**. Quick is a **project in constant development** and is open for **collaboration**, everyone is welcome to contribute. 😍

💡 If you’re new to coding, Quick is a great opportunity to start learning how to work with Go. With its **ease of use** and features, you can **create custom routes** and expand your knowledge of the language.

👍 I hope you can participate and enjoy **Enjoy**! 😍

🔍 The repository of examples of the Framework Quick Run [Examples](https://github.com/jeffotoni/quick/tree/main/example).


# Quick in action 💕🐧🚀😍
![Quick](quick_server.gif)

## 🎛️| Features

| Features                                         | Has  | Status | Completion |
|--------------------------------------------------|------|--------|------------|
| 🛣️ Route Manager                                  | yes  | 🟢     | 100%       |
| 📁 Server Files Static                            | yes  | 🟢     | 100%       |
| 🔗 Http Client                                    | yes  | 🟢     | 100%       |
| 📤 Upload Files (multipart/form-data)             | yes  | 🟢     | 100%       |
| 🚪 Route Group                                    | yes  | 🟢     | 100%       |
| 🛡️ Middlewares                                    | yes  | 🟡     | 50%        |
| ⚡ HTTP/2 support                                  | yes  | 🟢     | 100%       |
| 🔄 Data binding for JSON, XML and form payload    | yes  | 🟢     | 100%       |
| 🔍 Regex support                                  | yes  | 🟡     | 80%        |
| 🌎 Site                                           | yes  | 🟡     | 90%        |
| 📚 Docs                                           | yes  | 🟡     | 40%        |


## 🗺️| Development Rodmap

| Task   | Progress |
|-------------------------------------------------|-----------|
| Develop MaxBodySize method Post   | 100%   |
| Develop MaxBodySize method Put   | 100%   |
| Develop Config in New(Config{}) not required | 100%   |
| Create print function to not use fmt too much | 100% |
| Creation of own function for Concat String   | 100%   |
| Creation of benchmarking between the. Stdout and fmt.Println | 100%   |
| Develop Routes GET method   | 100%   |
| Develop Routes GET method by accepting Query String | 100%   |
| Develop Routes GET method accepting Parameters | 100%   |
| Develop Routes GET method accepting Query String and Parameters | 100% |
| Develop Routes GET method accepting regular expression | 100% |
| Develop Routes Method POST   | 100%   |
| Develop Routes POST method accepting JSON   | 100%   |
| Develop for METHOD POST the parse JSON   | 100%   |
| Develop for the POST METHOD functions to access byte or string from Parse | 100% |
| Develop for PUT METHOD   | 100%   |
| Develop for the PUT METHOD the JSON parse   | 100%   |
| Develop for the PUT METHOD the JSON parse   | 100%   |
| Develop for METHOD PUT functions to access byte or string from the Parse | 100% |
| Develop for DELETE METHOD   | 100%   |
| Develop method for ListenAndServe   | 100%   |
| Develop ServeHTTP support   | 100%   |
| Develop middleware support   | 100%   |
| Develop support for middleware compress   | 100%   |
| Develop support for middleware cors   | 100%   |
| Develop logger middleware support   | 100%   |
| Develop support for maxbody middlewares   | 100%   |
| Develop middleware support msgid   | 100%   |
| Develop middleware support msguuid   | 100%   |
| Develop support Cors   | 100%   |
| Develop Cient Get   | 100%   |
| Develop Cient Post support   | 100%   |
| Develop Cient Put support   | 100%   |
| Develop Cient support Delete   | 100%   |


## 🚧| Rodmap in progress

 | Task   | Progress |
|--------------------------------------------------|-----------|
| Develop and relate to Listen the Config   | 42%   |
| Develops support for Uploads and Uploads Multiples | 100% |
| Develops support for JWT | 10% |
| Develop method to Facilitate ResponseWriter handling | 80% |
| Develop method to Facilitate the handling of the Request | 80%  |
| Develop Standard of Unit Testing   | 90%   |


## 🚀| Rodmap for development


| Task   | Progress |
|---------------------------------------------------|-----------|
| Documentation Tests Examples PKG Go | 45% |
| Test Coverage go test -cover | 74.6% |
| Regex feature coverage, but possibilities | 0.% |
| Develop for OPTIONS METHOD | 100% |
| Develop for CONNECT METHOD [See more](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect) | 0.% |
| Develop method for ListenAndServeTLS (http2) | 0.% |
| Develop Static Files support | 100% |
| WebSocket Support | 0.% |
| Rate Limiter Support                              | 0.%       |
| Template Engines                                  | 0.%       |
| Documentation Tests Examples PKG Go   | 45%   |
| Test coverage go test -cover   | 75.5%   |
| Coverage of Regex resources, but possibilities   | 0.%   |
| Develop for METHOD OPTIONS   | 100%   |
| Develop for CONNECT METHOD [See more](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect)   | 0.%   |
| Develop method for ListenAndServeTLS (http2) | 0.%   |
| Create a CLI (Command Line Interface) Quick.   | 0.%   |



## 📊| Cover Testing Roadmap

| Archive   | Coverage | Status  |
|---------------|-----------|---------|
| Ctx   | 84.1%   | 🟡   |
| Group   | 100.0%   | 🟢   |
| Http Status  | 7.8%   | 🔴   |
| Client   | 83.3%   | 🟢   |
| Mock   | 100.0%   | 🟢   |
| Concat   | 100.0%   | 🟢   |
| Log   | 0.0%   | 🔴   |
| Print   | 66.7%   |  🟡   |
| Qos   | 0.0%   |  🔴   |
| Rand   | 0.0%   |  🔴   |
| Compressa   | 71,4%   | 🟡   |
| Cors   | 76.0%   | 🟡   |
| Logger   | 100.0%   | 🟢   |
| Maxbody   | 100.0%   | 🟢   |
| Msgid   | 100.0%   | 🟢   |
| Msguuid   | 86.4%   | 🟢   |
| Quick   | 79.5%   | 🟡   |
| QuickTest   | 100.0%   | 🟢   |


### Fast quick example
```go

package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New()

    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action ❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/user'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 23

Quick in action ❤️!

```

### Quick Get Params
```go

package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New()

    q.Get("/v1/customer/:param1/:param2", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")

        type my struct {
            Msg string `json:"msg"`
            Key string `json:"key"`
            Val string `json:"val"`
        }

        return c.Status(200).JSON(&my{
            Msg: "Quick ❤️",
            Key: c.Param("param1"),
            Val: c.Param("param2"),
        })
    })

    q.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/customer/val1/val2'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 23

{"msg":"Quick ❤️","key":"val1","val":"val2"}

```

### Quick Post Body json
```go

package main

import "github.com/jeffotoni/quick"

type My struct {
    Name string `json:"name"`
    Year int    `json:"year"`
}

func main() {
    q := quick.New()
    q.Post("/v1/user", func(c *quick.Ctx) error {
        var my My
        err := c.BodyParser(&my)
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }

        return c.Status(200).String(c.BodyString())
        // ou 
        // c.Status(200).JSON(&my)
    })

    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL

```bash

$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v1/user' \
-d '{"name":"jeffotoni", "year":1990}'
HTTP/1.1 200 OK
Date: Wed, 22 Feb 2023 08:10:06 GMT
Content-Length: 32
Content-Type: text/plain; charset=utf-8

{"name":"jeffotoni","year":1990}

```

### Uploads multipart/form-data

Quick provides a simplified API for managing uploads, allowing you to easily retrieve and manipulate files.

✅ **Main Methods and Functionalities**:
| Method | Description |
|--------|-----------|
| `c.FormFile("file")` | Returns a single file uploaded in the form. |
| `c.FormFiles("files")` | Returns a list of uploaded files (multiple uploads). |
| `c.FormFileLimit("10MB")` | Sets an upload limit (default is `1MB`). |
| `uploadedFile.FileName()` | Returns the file name. |
| `uploadedFile.Size()` | Returns the file size in bytes. |
| `uploadedFile.ContentType()` | Returns the MIME type of the file. |
| `uploadedFile.Bytes()` | Returns the bytes of the file. |
| `uploadedFile.Save("/path/")` | Saves the file to a specified directory. |
| `uploadedFile.Save("/path", "your-name-file")` | Saves the file with your name. |
| `uploadedFile.SaveAll("/path")` | Saves the file to a specified directory. |

---

##### 📌 File Upload Feature Comparison with other Frameworks

| Framework  | `FormFile()` | `FormFiles()` | Dynamic Limit | Methods (`FileName()`, `Size()`) | `Save()`, `SaveAll()` Method |
|------------|-------------|--------------|---------------|---------------------------------|------------|
| **Quick**  | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| Fiber      | ✅ Yes | ✅ Yes | ❌ No | ❌ No (uses `FileHeader` directly) | ✅ Yes |
| Gin        | ✅ Yes | ✅ Yes | ❌ No | ❌ No (uses `FileHeader` directly) | ❌ No |
| Echo       | ✅ Yes | ❌ No  | ❌ No | ❌ No | ❌ No |
| net/http   | ✅ Yes | ❌ No  | ❌ No | ❌ No | ❌ No |

---

### 📌 File Upload Example

```go
package main

import (
    "fmt"
    "github.com/jeffotoni/quick"
)

func main() {
    // start Quick
    q := quick.New()

    q.Post("/upload", func(c *quick.Ctx) error {
        // set limit upload
        c.FormFileLimit("10MB")

        uploadedFile, err := c.FormFile("file")
        if err != nil {
            return c.Status(400).JSON(Msg{
                Msg: "Upload error",
                Error: err.Error(),
             })
        }

        fmt.Println("Name:", uploadedFile.FileName())
        fmt.Println("Size:", uploadedFile.Size())
        fmt.Println("MIME Type:", uploadedFile.ContentType())

        // Save the file (optional)
        // uploadedFile.Save("/tmp/uploads")

        return c.Status(200).JSONIN(uploadedFile)

    })

     q.Listen("0.0.0.0:8080")
}
```
### 📌 Multiple Upload Example

```go
package main

import (
    "fmt"
    "github.com/jeffotoni/quick"
)

func main() {
    // start Quick
    q := quick.New()

    q.Post("/upload-multiple", func(c *quick.Ctx) error {
        // set limit upload
        c.FormFileLimit("10MB")

        // recebereceiving files
        files, err := c.FormFiles("files")
        if err != nil {
            return c.Status(400).JSON(Msg{
                Msg:   "Upload error",
                Error: err.Error(),
            })
        }

        // listing all files
        for _, file := range files {
            fmt.Println("Name:", file.FileName())
            fmt.Println("Size:", file.Size())
            fmt.Println("Type MINE:", file.ContentType())
            fmt.Println("Bytes:", file.Bytes())
        }

        // optional
        // files.SaveAll("/my-dir/uploads")

        return c.Status(200).JSONIN(files)
    })
        
    q.Listen("0.0.0.0:8080")
}
```
### 📌 Testing with cURL

##### 🔹Upload a single file:
```bash

$ curl -X POST http://localhost:8080/upload -F "file=@example.png"
```

##### 🔹 Upload multiple files:
```bash

$ curl -X POST http://localhost:8080/upload-multiple \
-F "files=@image1.jpg" -F "files=@document.pdf"
```

### Quick Post Bind json
```go

package main

import "github.com/jeffotoni/quick"

type My struct {
    Name string `json:"name"`
    Year int    `json:"year"`
}

func main() {
    q := quick.New()
    q.Post("/v2/user", func(c *quick.Ctx) error {
        var my My
        err := c.Bind(&my)
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }
        return c.Status(200).JSON(&my)
    })

    q.Listen("0.0.0.0:8080")
}

```

### 📌 cURL
```bash

$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v2/user' \
-d '{"name":"Marcos", "year":1990}'
HTTP/1.1 200 OK
Date: Wed, 22 Feb 2023 08:10:06 GMT
Content-Length: 32
Content-Type: text/plain; charset=utf-8

{"name":"Marcos","year":1990}

```

### Cors

```go

package main

import (
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cors"
)

func main() {
    q := quick.New()
    q.Use(cors.New())

    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action com Cors❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

### quick.New(quick.Config{})
```go

package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New(quick.Config{
        MaxBodySize: 5 * 1024 * 1024,
    })

    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action com Cors❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

### quick.Group()
```go
package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New(quick.Config{
        MaxBodySize: 5 * 1024 * 1024,
    })

    v1 := q.Group("/v1")
    v1.Get("/user", func(c *quick.Ctx) error {
        return c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
    })
    v1.Post("/user", func(c *quick.Ctx) error {
        return c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
    })

    v2 := q.Group("/v2")
    v2.Get("/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action com [GET] /v2/user ❤️!")
    })

    v2.Post("/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action com [POST] /v2/user ❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

### Quick Tests
```go

package main

import (
    "io"
    "strings"
    "testing"

    "github.com/jeffotoni/quick"
)

func TestQuickExample(t *testing.T) {

    // Here is a handler function Mock
    testSuccessMockHandler := func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        b, _ := io.ReadAll(c.Request.Body)
        resp := `"data":` + string(b)
        return c.Byte([]byte(resp))
    }

    q := quick.New()
    // Here you can create all routes that you want to test
    q.Post("/v1/user", testSuccessMockHandler)
    q.Post("/v1/user/:p1", testSuccessMockHandler)

    wantOutData := `"data":{"name":"jeff", "age":35}`
    reqBody := []byte(`{"name":"jeff", "age":35}`)
    reqHeaders := map[string]string{"Content-Type": "application/json"}

    data, err := q.QuickTest("POST", "/v1/user", reqHeaders, reqBody)
    if err != nil {
        t.Errorf("error: %v", err)
        return
    }

    s := strings.TrimSpace(data.BodyStr())
    if s != wantOutData {
        t.Errorf("was suppose to return %s and %s come", wantOutData, s)
        return
    }

    t.Logf("\nOutputBodyString -> %v", data.BodyStr())
    t.Logf("\nStatusCode -> %d", data.StatusCode())
    t.Logf("\nOutputBody -> %v", string(data.Body())) 
    t.Logf("\nResponse -> %v", data.Response())
}

```

### quick.regex
```go

package main

import (
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/msgid"
)

func main() {
    q := quick.New()

    q.Use(msgid.New())

    q.Get("/v1/user/{id:[0-9]+}", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).String("Quick ação total!!!")
    })

    q.Listen("0.0.0.0:8080")
}

```
### 🔑 Basic Authentication

Basic Authentication (Basic Auth) is a simple authentication mechanism defined in RFC 7617. It is commonly used for HTTP-based authentication, allowing clients to provide credentials (username and password) in the request header.

**🔹 How it Works** 
  1.	The client encodes the username and password in Base64 (username:password → dXNlcm5hbWU6cGFzc3dvcmQ=).
  2.	The encoded credentials are sent in the Authorization header:
```bash

Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=
```
  3.	The server decodes the credentials and verifies them before granting access.

**⚠️ Security Considerations**
	• Not encrypted: Basic Auth only encodes credentials in Base64, but does not encrypt them.
	• Use over HTTPS: Always use Basic Auth with TLS/SSL (HTTPS) to prevent credentials from being exposed.
	• Alternative authentication methods: For higher security, consider OAuth2, JWT, or API keys.

Basic Auth is suitable for simple use cases, but for production applications, stronger authentication mechanisms are recommended. 🚀

#### Basic Auth environment variables

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

// export USER=admin
// export PASSWORD=1234

var (
	User     = os.Getenv("USER")
	Password = os.Getenv("PASSORD")
)

func main() {

	q := quick.New()

	q.Use(middleware.BasicAuth(User, Password))

	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

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

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

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

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

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
---

#### 📂 STATIC FILES

A Static File Server is a fundamental feature in web frameworks, allowing the efficient serving of static content such as HTML, CSS, JavaScript, images, and other assets. It is useful for hosting front-end applications, providing downloadable files, or serving resources directly from the backend.


🔹 How It Works
    
1. The server listens for HTTP requests targeting static file paths.
2. If a requested file exists in the configured directory, the server reads and returns the file as a response.
3. MIME types are automatically determined based on the file extension.

:zap: Key Features
- Efficient handling: Serves files directly without additional processing.
- MIME type detection: Automatically identifies file types for proper rendering.
- Caching support: Can be configured to improve performance via HTTP headers.
- Directory listing: (Optional) Allows browsing available static files.

:warning: Security Considerations
- Restrict access to sensitive files (.env, .git, etc.).
- Configure CORS policies when necessary.
- Use a Content Security Policy (CSP) to mitigate XSS risks.

#### Serving Static Files with Quick Framework

This example sets up a basic web server that serves static files, such as HTML, CSS, or JavaScript.

```go
package main

import "github.com/jeffotoni/quick"

func main() {
    
    // Create a new Quick server instance
    q := quick.New()

    // Static Files Setup
    // Serves files from the "./static" directory 
    // under the "/static" URL path.
    q.Static("/static", "./static")

    // Route Definition
    // Defines a route to serve the "index.html" file when accessing "/".
    q.Get("/", func(c *quick.Ctx) error {
        c.File("./static/index.html")
        return nil
    })

    // Starting the Server
    // Starts the server to listen for incoming requests on port 8080.
    q.Listen("0.0.0.0:8080")
}


```
---

#### 📁 EMBED
🔹 How Embedded Static Files Work
    
1. Static assets are compiled directly into the binary at build time (e.g., using Go’s embed package).
2. The application serves these files from memory instead of reading from disk.
3. This eliminates external dependencies, making deployment easier.

:zap:  Advantages of Embedded Files
- Portability: Single binary distribution without extra files.
- Performance: Faster access to static assets as they are stored in memory.
- Security: Reduces exposure to external file system attacks.

### Embedding Files
When embedding static files into a binary executable, the server does not rely on an external file system to serve assets. This approach is useful for standalone applications, CLI tools, and cross-platform deployments where dependencies should be minimized.

This example incorporates static files into the binary using the embed package and serves them using the Quick structure.

```go
package main

import (
	"embed"

	"github.com/jeffotoni/quick"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Server Initialization
	// Creates a new instance of the Quick server
	q := quick.New()

	// Static Files Setup
	// Defines the directory for serving static files using the embedded files
	q.Static("/static", staticFiles)

	// Route Definition
	// Defines a route that serves the HTML index file
	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/index.html") // Renders the index.html file
		return nil
	})

	// Starting the Server
	// Starts the server on port 8080, listening on all addresses
	q.Listen("0.0.0.0:8080")
}

```
---
## 🌍 HTTP Client 
The HTTP Client package in Quick provides a simple and flexible way to make HTTP requests, supporting GET, POST, PUT, and DELETE operations. It is designed to handle different types of request bodies and parse responses easily.

This client abstracts low-level HTTP handling and offers:

- Convenience functions (Get, Post, Put, Delete) for making quick requests using a default client.
- Customizable requests with support for headers, authentication, and transport settings.
- Flexible body parsing, allowing users to send JSON, plain text, or custom io.Reader types.
- Automatic JSON marshaling and unmarshaling, simplifying interaction with APIs.

#### GET Request Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Use the default client
	resp, err := client.Get("https://reqres.in/api/users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}
```

#### POST Request Example (Using a Struct)
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Define a struct to send as JSON
	data := struct {
		user string `json:"user"`
	}{
		user: "Emma",
	}

	// POST request to ReqRes API
	resp, err := client.Post("https://reqres.in/api/users", data)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}
```

#### PUT Request Example (Using a String)
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Define a struct with user data
	data := struct {
		user string `json:"name"`
	}{
		user: "Jeff",
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	// PUT request to ReqRes API
	resp, err := client.Put("https://reqres.in/api/users/2", string(jsonData))
	if err != nil {
		log.Fatal("Error making request:", err)
	}

	// Print the HTTP status and response body
	fmt.Println("HTTP Status Code:", resp.StatusCode)
	fmt.Println("Raw Response Body:", string(resp.Body))
}
```

#### DELETE Request Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// DELETE request to ReqRes API
	resp, err := client.Delete("https://reqres.in/api/users/2")
	if err != nil {
		log.Fatal("Error making request:", err)
	}

	// Print the HTTP status to confirm deletion
	fmt.Println("HTTP Status Code:", resp.StatusCode)

	// Since DELETE usually returns no content, we check if it's empty
	if len(resp.Body) > 0 {
		fmt.Println("Raw Response Body:", string(resp.Body))
	} else {
		fmt.Println("Response Body is empty (expected for 204 No Content)")
	}
}
```
---

# Qtest - HTTP Testing Utility for Quick

Qtest is an **advanced HTTP testing function** designed to simplify route validation within the **Quick** framework. It enables seamless testing of simulated HTTP requests using `httptest`, supporting:

- **Custom HTTP methods** (`GET`, `POST`, `PUT`, `DELETE`, etc.).
- **Custom headers**.
- **Query parameters**.
- **Request body**.
- **Cookies**.
- **Built-in validation methods** for status codes, headers, and response bodies.

## 📌 Overview
The `Qtest` function takes a `QuickTestOptions` struct containing request parameters, executes the request, and returns a `QtestReturn` object, which provides methods for analyzing and validating the result.
 

```go
func TestQTest_Options_POST(t *testing.T) {
    // start Quick
    q := New()

    // Define the POST route
    q.Post("/v1/user/api", func(c *Ctx) error {
        c.Set("Content-Type", "application/json") // Simplified header setting
        return c.Status(StatusOK).String(`{"message":"Success"}`)
    })

    // Configure test parameters
    opts := QuickTestOptions{
        Method: "POST",
        URI:    "/v1/user/api",
        QueryParams: map[string]string{
            "param1": "value1",
            "param2": "value2",
        },
        Body: []byte(`{"key":"value"}`),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Cookies: []*http.Cookie{
            {Name: "session", Value: "abc123"},
        },
        LogDetails: true, // Enables detailed logging
    }

    // Execute test
    result, err := q.Qtest(opts)
    if err != nil {
        t.Fatalf("Error in Qtest: %v", err)
    }

    // Validations
    if err := result.AssertStatus(StatusOK); err != nil {
        t.Errorf("Status assertion failed: %v", err)
    }

    if err := result.AssertHeader("Content-Type", "application/json"); err != nil {
        t.Errorf("Header assertion failed: %v", err)
    }

    if err := result.AssertBodyContains("Success"); err != nil {
        t.Errorf("Body assertion failed: %v", err)
    }
}
```
🚀 **More details here [Qtest - Quick](https://github.com/jeffotoni/quick/tree/main/quickTest)**

---
## 🔄 Retry & Failover Mechanisms in Quick HTTP Client

The **Quick HTTP Client** now includes **built-in retry and failover support**, allowing for more resilient and reliable HTTP requests. These features are essential for handling **transient failures**, **network instability**, and **service downtime** efficiently.

## 🚀 Key Features
- **Automatic Retries**: Retries failed requests based on configurable rules.
- **Exponential Backoff**: Gradually increases the delay between retry attempts.
- **Status-Based Retries**: Retries only on specified HTTP status codes (e.g., `500`, `502`, `503`).
- **Failover Mechanism**: Switches to predefined backup URLs if the primary request fails.
- **Logging Support**: Enables detailed logs for debugging retry behavior.

---

## 🔹 How Retry & Failover Work
The retry mechanism works by **automatically resending the request** if it fails, with options to **limit retries**, **introduce backoff delays**, and **retry only for specific response statuses**. The failover system ensures **high availability** by redirecting failed requests to alternate URLs.

### ✅ Configuration Options:
| Option                | Description |
|-----------------------|-------------|
| **MaxRetries**       | Defines the number of retry attempts. |
| **Delay**            | Specifies the delay before each retry. |
| **UseBackoff**       | Enables exponential backoff to increase delay dynamically. |
| **Statuses**         | List of HTTP status codes that trigger a retry. |
| **FailoverURLs**     | List of backup URLs for failover in case of repeated failures. |
| **EnableLog**        | Enables logging for debugging retry attempts. |

---

### **Retry with Exponential Backoff**
This example demonstrates **retrying a request** with an increasing delay (`backoff`) when encountering errors.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cClient := client.New(
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 3,       // Maximum number of retries
				Delay:      1 * time.Second,  // Initial retry delay
				UseBackoff: true,    // Enables exponential backoff
				Statuses:   []int{500, 502, 503}, // Retries only on these HTTP status codes
				EnableLog:  true,    // Enables logging for retries
			}),
	)

	resp, err := cClient.Get("http://localhost:3000/v1/resource")
	if err != nil {
		log.Fatal("GET request failed:", err)
	}
	fmt.Println("GET Response:", string(resp.Body))
}

```
### **Failover to Backup URLs**
This example switches to a backup URL when the primary request fails.
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cClient := client.New(
		client.WithRetry(client.RetryConfig{
			MaxRetries:   2,  // Try the request twice before switching
			Delay:        2 * time.Second,  // Wait 2 seconds before retrying
			Statuses:     []int{500, 502, 503}, // Trigger failover on these errors
			FailoverURLs: []string{"http://backup1.com/resource", "https://reqres.in/api/users", "https://httpbin.org/post"},
			EnableLog: true, // Enable retry logs
		}),
	)

	resp, err := cClient.Get("http://localhost:3000/v1/resource")
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	fmt.Println("Response:", string(resp.Body))
}

```
---
## 📝 Form Submission with PostForm in Quick HTTP Client

The Quick HTTP Client now includes built-in support for `PostForm`, enabling seamless handling of application/`x-www-form-urlencoded` form submissions. This feature simplifies interaction with web services and APIs that require form-encoded data, making it ideal for authentication requests, data submissions, and legacy system integrations.


## 🔹 Why Use `PostForm`?

| Feature                | Benefit |
|------------------------|---------|
| **Optimized for Forms** | Simplifies sending form-encoded data (`application/x-www-form-urlencoded`). |
| **Automatic Encoding**  | Converts `url.Values` into a valid form submission payload. |
| **Header Management**   | Automatically sets `Content-Type` to `application/x-www-form-urlencoded`. |
| **Consistent API**      | Follows the same design as `Post`, `Get`, `Put`, etc. |
| **Better Compatibility** | Works with APIs that do not accept JSON payloads. |

---
## 🔹 How PostForm Works

The PostForm method encodes form parameters, adds necessary headers, and sends an HTTP POST request to the specified URL. It is specifically designed for APIs and web services that do not accept JSON payloads but require form-encoded data.


### 🔹 **Quick Server with Form Submission**
The following example demonstrates how to send form-encoded data using Quick PostForm:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Define a route to process POST form-data
	q.Post("/postform", func(c *quick.Ctx) error {
		form := c.FormValues()
		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Quick server running at http://localhost:3000")
		if err := q.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Quick server: %v", err)
		}
	}()

	// Creating an HTTP client before calling PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Define um timeout de 5s
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // Correct type for forms
		}),
	)

	// Check if the HTTP client was initialized correctly
	if cClient == nil {
		log.Fatal("Erro: cliente HTTP não foi inicializado corretamente")
	}

	// Declare Values
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Send a POST request
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Check if the response is valid
	if resp == nil || resp.Body == nil {
		log.Fatal("Erro: resposta vazia ou inválida")
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}

```
---
### 🌐 Transport Configuration 

The `Transport` setting in the Quick HTTP Client is essential for managing the network layer of HTTP communications. It allows detailed customization of how HTTP requests and responses are handled, optimizing performance, security, and reliability.

### ✅ Key Features of Transport Configuration

| Setting                  | Description |
|--------------------------|-------------|
| **Proxy Settings**       | Manages how HTTP requests handle proxy servers, using system environment settings for automatic configuration. |
| **TLS Configuration**    | Controls aspects of security, such as TLS version and certificate verification. `InsecureSkipVerify` is available for development to bypass SSL certificate verification. |
| **Connection Management**| Includes settings like `MaxIdleConns`, `MaxConnsPerHost`, and `MaxIdleConnsPerHost`, managing the number and state of connections to optimize resource use and improve scalability. |
| **DisableKeepAlives**    | Determines whether to use persistent connections, improving performance by reducing connection setup times. |
| **HTTP/2 Support**       | Enables HTTP/2 for supported servers, enhancing communication efficiency and performance. |

This configuration ensures optimal performance and security customization, making it suitable for both development and production environments.



#### 🔹 Advanced HTTP client configuration with failover mechanism

This code example showcases the setup of an HTTP client capable of handling network interruptions and server failures gracefully. It features custom transport configurations, including enhanced security settings, connection management, and a robust failover mechanism. Such a setup ensures that the application remains resilient and responsive under various network conditions. 

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	
	customTransport := &http.Transport{
		// Uses system proxy settings if available.
		Proxy: http.ProxyFromEnvironment, 
		TLSClientConfig: &tls.Config{
			// Allows insecure TLS connections (not recommended for production).
			InsecureSkipVerify: true,    
			// Enforces a minimum TLS version for security.         
			MinVersion:         tls.VersionTLS12, 
		},
		 // Maximum number of idle connections across all hosts.
		MaxIdleConns:        50,  
		// Maximum simultaneous connections per host. 
		MaxConnsPerHost:     30, 
		// Maximum number of idle connections per host.   
		MaxIdleConnsPerHost: 10, 
		// Enables persistent connections (Keep-Alive).   
		DisableKeepAlives:   false, 
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		// Sets a global timeout for all requests.
		Timeout: 5 * time.Second,
	}

	// Creating a client using both the custom transport and other configurations.
	cClient := client.New(
		// Applying the custom HTTP client.
		client.WithCustomHTTPClient(customHTTPClient), 
		// Custom context for request cancellation and deadlines.
		client.WithContext(context.Background()),      
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		// Applying the custom transport.
		client.WithTransport(customTransport), 
		// Setting a timeout for requests.
		client.WithTimeout(5*time.Second),     
		// Retry on specific status codes.
		client.WithRetry(
			client.RetryConfig{
				MaxRetries:   2,
				Delay:        1 * time.Second,
				UseBackoff:   true,
				Statuses:     []int{500},
				FailoverURLs: []string{"http://hosterror", "https://httpbin.org/post"},
				EnableLog:    true,
			}),
	)

	// call client to POST
	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"message": "Hello Post!!"})
	if err != nil {
		log.Fatal(err)
	}

	// show resp
	fmt.Println("POST response:\n", string(resp.Body))
}

```
---
#### 🔹 HTTP Client with Advanced Transport and Failover Capabilities

Explore how to set up an HTTP client that not only adheres to security best practices with TLS configurations but also ensures your application remains operational through network issues. This example includes detailed setups for handling HTTP client retries and switching to failover URLs when typical requests fail. Ideal for systems requiring high reliability and fault tolerance. 

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// Creating a custom HTTP transport with advanced settings.
	customTransport := &http.Transport{
		// Uses system proxy settings if available.
		Proxy: http.ProxyFromEnvironment, 
		TLSClientConfig: &tls.Config{
			 // Allows insecure TLS connections (not recommended for production).
			InsecureSkipVerify: true,      
			// Enforces a minimum TLS version for security.      
			MinVersion:         tls.VersionTLS12, 
		},
		// Maximum number of idle connections across all hosts.
		MaxIdleConns:        50, 
		// Maximum simultaneous connections per host.   
		MaxConnsPerHost:     30,   
		 // Maximum number of idle connections per host. 
		MaxIdleConnsPerHost: 10,  
		// Enables persistent connections (Keep-Alive). 
		DisableKeepAlives:   false, 
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		 // Sets a global timeout for all requests.
		Timeout:   5 * time.Second,
		// Uses the custom transport.
		Transport: customTransport, 
	}

	// Creating a client using both the custom transport and other configurations.
	cClient := client.New(
		// Applying the custom HTTP client.
		client.WithCustomHTTPClient(customHTTPClient), 
		 // Custom context for request cancellation and deadlines.
		client.WithContext(context.Background()),     
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		client.WithTimeout(5*time.Second), // Setting a timeout for requests.
		// Retry on specific status codes.
		client.WithRetry(
			client.RetryConfig{
				MaxRetries:   2,
				Delay:        1 * time.Second,
				UseBackoff:   true,
				Statuses:     []int{500},
				FailoverURLs: []string{"http://hosterror", "https://httpbin.org/post"},
				EnableLog:    true,
			}),
	)

	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	// show resp
	fmt.Println("POST response:", string(resp.Body))
}
```
---
### 🔹Configuring HTTP Client with Retry and Failover Mechanisms

Discover how to build an HTTP client capable of dealing with network instabilities and server failures. This setup includes detailed retry configurations and introduces failover URLs to ensure that your application can maintain communication under adverse conditions. The example demonstrates using exponential backoff for retries and provides multiple endpoints to guarantee the availability of services. 

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create a new HTTP client with specific configurations.
	cClient := client.New(
		// Set a timeout for all requests made by this client to 10 seconds.
		// This helps prevent the client from hanging indefinitely on requests.
		client.WithTimeout(10*time.Second),

		// Set default headers for all requests made by this client.
		// Here, 'Content-Type' is set to 'application/json'
		//  which is typical for API calls.
		client.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),

		// Enable automatic retry mechanism with specific configurations.
		// This is useful for handling intermittent errors and ensuring robustness.
		client.WithRetry(
			client.RetryConfig{
				 // Retry failed requests up to two times.
				MaxRetries: 2,      
				// Wait for 1 second before retrying.             
				Delay:      1 * time.Second, 
				 // Use exponential backoff strategy for retries.     
				UseBackoff: true,         
				// HTTP status codes that trigger a retry.       
				Statuses:   []int{500, 502, 503}, 
				// Alternate URLs to try if the main request fails.
				FailoverURLs: []string{ 
					"http://hosterror",
					"https://httpbin.org/post",
				},
				// Enable logging for retry operations.
				EnableLog: true, 
			}),
	)

	// Perform a POST request using the configured HTTP client.
	// Includes a JSON payload with a "name" key.
	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{
		"name": "jeffotoni in action with Quick!!!",
	})

	// Check if there was an error with the POST request.
	if err != nil {
		// If an error occurs, log the error and terminate the program.
		log.Fatalf("POST request failed: %v", err)
	}

	// Print the response from the server to the console.
	fmt.Println("POST Form Response:", string(resp.Body))
}

```
---
### 🔹Advanced HTTP Client Configuration with Transport and Retry Settings

Explore the configuration of an HTTP client designed for high reliability and security in network communications. This example includes sophisticated transport settings, featuring TLS configurations for enhanced security, and a robust retry mechanism to handle request failures gracefully. These settings are essential for applications requiring reliable data exchange with external APIs, especially in environments where network stability might be a concern. 

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create an HTTP client with custom configurations using the Quick framework.
	cClient := client.New(
		// Set a global timeout for all requests made by this client to 10 seconds.
		// This helps prevent the client from hanging indefinitely on requests.
		client.WithTimeout(10*time.Second),

		// Set default headers for all requests made by this client.
		// Here, we specify that we expect to send and receive JSON data.
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Configure the underlying transport for the HTTP client.
		client.WithTransportConfig(&http.Transport{
			// Use the system environment settings for proxy configuration.
			Proxy: http.ProxyFromEnvironment,

			// Configure TLS settings to skip verification of the server's 
			// certificate chain and hostname.
			// Warning: Setting InsecureSkipVerify to true is not recommended for
			//  production as it is insecure.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},

			// Enable HTTP/2 for supported servers.
			ForceAttemptHTTP2: true,

			// Set the maximum number of idle connections in the connection pool for all hosts.
			MaxIdleConns: 20,

			// Set the maximum number of idle connections in the connection pool per host.
			MaxIdleConnsPerHost: 10,

			// Set the maximum number of simultaneous connections per host.
			MaxConnsPerHost: 20,

			// Keep connections alive between requests. This can help improve performance.
			DisableKeepAlives: false,
		}),
	)

	// Perform a POST request with a JSON payload.
	// The payload includes a single field "name" with a value.
	resp, err := cClient.Post("https://httpbin.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		// Log the error and stop the program if the POST request fails.
		log.Fatalf("POST request failed: %v", err)
	}

	// Output the response from the POST request.
	fmt.Println("POST Form Response:", string(resp.Body))
}

```


---

## 📚| More Examples

This directory contains practical examples of the Quick Framework, a fast and lightweight web framework developed in Go. The examples are organized in separate folders, each containing a complete example of using the framework in a simple web application. If you have some interesting example of using the Quick Framework, feel free to send a pull request with your contribution. The Quick Framework example repository can be found at [here](https://github.com/jeffotoni/quick/tree/main/example).

## 🤝| Contributions

We already have several examples, and we can already test and play 😁. Of course, we are at the beginning, still has much to do. 
Feel free to do **PR** (at risk of winning a Go t-shirt ❤️ and of course recognition as a professional Go 😍 in the labor market).


## 🚀 **Quick Project Supporters** 🙏

The Quick Project aims to develop and provide quality software for the developer community. 💻 To continue improving our tools, we rely on the support of our sponsors in Patreon. 🤝

We thank all our supporters! 🙌 If you also believe in our work and want to contribute to the advancement of the development community, consider supporting Project Quick on our Patreon [here](https://www.patreon.com/jeffotoni_quick)

Together we can continue to build amazing tools! 🚀

| Avatar | User | Donation |
|--------|------|----------|
| <img src="https://avatars.githubusercontent.com/u/1092879?s=96&v=4" height=20> | [@jeffotoni](https://github.com/jeffotoni) | x 10 |
| <img src="https://avatars.githubusercontent.com/u/99341377?s=400&u=095679b08054e215561a4d4b08da764c2de619e6&v=4" height=20> | [@Crow3442](https://github.com/Crow3442) | x 5  |
| <img src="https://avatars.githubusercontent.com/u/70351793?v=4" height=20> | [@Guilherme-De-Marchi](https://github.com/Guilherme-De-Marchi) | x 5 |
| <img src="https://avatars.githubusercontent.com/u/59976892?v=4" height=20> | [@jaquelineabreu](https://github.com/jaquelineabreu) | x 1 |












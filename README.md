![Logo do Quick](./quick_logo.png)

[![GoDoc](https://godoc.org/github.com/jeffotoni/quick?status.svg)](https://godoc.org/github.com/jeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jeffotoni/quick/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jeffotoni/quick/tree/main) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick) [![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/jeffotoni/quick/main) ![Coveralls](https://img.shields.io/coverallsCoverage/github/jeffotoni/quick) ![GitHub contributors](https://img.shields.io/github/contributors/jeffotoni/quick)
![GitHub stars](https://img.shields.io/github/last-commit/jeffotoni/quick) ![GitHub stars](https://img.shields.io/github/forks/jeffotoni/quick?style=social) ![GitHub stars](https://img.shields.io/github/stars/jeffotoni/quick)

<!-- [![Github Release](https://img.shields.io/github/v/release/jeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/jeffotoni/quick) -->

<h2 align="center">
    <p>
         <a href="README.md">English</a> |
          <a href="README.pt-br.md">Рortuguês</a>
    </p> 
</h2>

```bash
   ██████╗ ██╗   ██╗██╗ ██████╗██╗  ██╗
  ██╔═══██╗██║   ██║██║██╔═══   ██║ ██╔╝
  ██║   ██║██║   ██║██║██║      █████╔╝
  ██║▄▄ ██║██║   ██║██║██║      ██╔═██╗
  ╚██████╔╝╚██████╔╝██║╚██████╔ ██║  ██╗
   ╚══▀▀═╝  ╚═════╝ ╚═╝ ╚═════╝ ╚═╝  ╚═╝

 Quick v0.0.1 🚀 Fast & Minimal Web Framework
─────────────────── ───────────────────────────────
 🌎 Host : http://127.0.0.1:0.0.0.0:8080
 📌 Port : 0.0.0.0:8080
 🔀 Routes: 4
─────────────────── ───────────────────────────────

```

# Quick - a lightweight router for go ![Quick Logo](./quick.png)

🚀 Quick is a **flexible and extensible** route manager for the Go language. Its goal is to be **fast and high-performance**, as well as being **100% compatible with net/http**. Quick is a **project in constant development** and is open for **collaboration**, everyone is welcome to contribute. 😍

💡 If you’re new to coding, Quick is a great opportunity to start learning how to work with Go. With its **ease of use** and features, you can **create custom routes** and expand your knowledge of the language.

👍 I hope you can participate and enjoy **Enjoy**! 😍

🔍 The repository of examples of the Framework Quick Run [Examples](https://github.com/jeffotoni/quick/tree/main/example).

# Quick in action 💕🐧🚀😍

![Quick](quick_server.gif)

---

## 📦 Go Packages Documentation

To access the documentation for each **Quick Framework** package, click on the links below:

| Package                | Description                                       | Go.dev                                                                                                                                    |
| ---------------------- | ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| **quick**              | Main Router and Framework Features                | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick)                        |
| **quick/http/client**  | HTTP client optimized for requests and failover   | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick/http/client.svg)](https://pkg.go.dev/githu.com/jeffotoni/quick/http/client) |
| **quick/middleware**   | Framework middlewares                             | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick/middleware.svg)](https://pkg.go.dev/github.com/jeffotoni/quick/middleware)  |
| **quick/ctx**          | HTTP request and response context                 | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#Ctx)                    |
| **quick/http/status**  | HTTP status definitions in the framework          | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#Status)                 |
| **quick/group**        | Route group manipulation                          | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#Group)                  |
| **quickTest**          | Package for unit testing and integration in Quick | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick/quickTest.svg)](https://pkg.go.dev/github.com/jeffotoni/quickTest)          |
| **quick/route**        | Route definition and management                   | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#Route)                  |
| **quick/config**       | Framework configuration structures                | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#Config)                 |
| **quick/qtest**        | Auxiliary tools for testing in the Quick          | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#QTest)                  |
| **quick/uploadedFile** | File upload management                            | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#UploadedFile)           |
| **quick/zeroth**       | Framework helpers                                 | [![GoDoc](https://pkg.go.dev/badge/github.com/jeffotoni/quick.svg)](https://pkg.go.dev/github.com/jeffotoni/quick#Zeroth)                 |

---

## 🎛️| Features

| Features                                       | Has | Status | Completion |
| ---------------------------------------------- | --- | ------ | ---------- |
| 🛣️ Route Manager                               | yes | 🟢     | 100%       |
| 📁 Server Files Static                         | yes | 🟢     | 100%       |
| 🔗 Http Client                                 | yes | 🟢     | 100%       |
| 📤 Upload Files (multipart/form-data)          | yes | 🟢     | 100%       |
| 🚪 Route Group                                 | yes | 🟢     | 100%       |
| 🛡️ Middlewares                                 | yes | 🟡     | 50%        |
| ⚡ HTTP/2 support                              | yes | 🟢     | 100%       |
| 🔄 Data binding for JSON, XML and form payload | yes | 🟢     | 100%       |
| 🔍 Regex support                               | yes | 🟡     | 80%        |
| 🌎 Site                                        | yes | 🟡     | 90%        |
| 📚 Docs                                        | yes | 🟡     | 40%        |

# 🗺️ Development Roadmap  

## ✅ Completed Tasks  

| Task                                                                      | Progress |
| ------------------------------------------------------------------------- | -------- |
| Develop MaxBodySize method Post                                           | ✅ 100%  |
| Develop MaxBodySize method Put                                            | ✅ 100%  |
| Develop Config in New(Config{}) not required                              | ✅ 100%  |
| Create print function to not use fmt too much                             | ✅ 100%  |
| Creation of own function for Concat String                                | ✅ 100%  |
| Creation of benchmarking between the Stdout and fmt.Println               | ✅ 100%  |
| Develop Routes GET method                                                 | ✅ 100%  |
| Develop Routes GET method by accepting Query String                       | ✅ 100%  |
| Develop Routes GET method accepting Parameters                            | ✅ 100%  |
| Develop Routes GET method accepting Query String and Parameters           | ✅ 100%  |
| Develop Routes GET method accepting regular expression                    | ✅ 100%  |
| Develop Routes Method POST                                                | ✅ 100%  |
| Develop Routes POST method accepting JSON                                 | ✅ 100%  |
| Develop for METHOD POST the parse JSON                                    | ✅ 100%  |
| Develop for the POST METHOD functions to access byte or string from Parse | ✅ 100%  |
| Develop for PUT METHOD                                                    | ✅ 100%  |
| Develop for the PUT METHOD the JSON parse                                 | ✅ 100%  |
| Develop for the PUT METHOD the JSON parse                                 | ✅ 100%  |
| Develop for METHOD PUT functions to access byte or string from the Parse  | ✅ 100%  |
| Develop for DELETE METHOD                                                 | ✅ 100%  |
| Develop method for ListenAndServe                                         | ✅ 100%  |
| Develop ServeHTTP support                                                 | ✅ 100%  |
| Develop middleware support                                                | ✅ 100%  |
| Develop support for middleware compress                                   | ✅ 100%  |
| Develop support for middleware cors                                       | ✅ 100%  |
| Develop logger middleware support                                         | ✅ 100%  |
| Develop support for maxbody middlewares                                   | ✅ 100%  |
| Develop middleware support msgid                                          | ✅ 100%  |
| Develop middleware support msguuid                                        | ✅ 100%  |
| Develop support Cors                                                      | ✅ 100%  |
| Develop Cient Get                                                         | ✅ 100%  |
| Develop Cient Post support                                                | ✅ 100%  |
| Develop Cient Put support                                                 | ✅ 100%  |
| Develop Cient support Delete                                              | ✅ 100%  |

---

## 🚧 Roadmap in Progress  

| Task                                                     | Progress |
| -------------------------------------------------------- | -------- |
| Develop and relate to Listen the Config                  | ⏳ 42%   |
| Develops support for Uploads and Uploads Multiples       | ✅ 100%  |
| Develops support for JWT                                 | ⏳ 10%   |
| Develop method to Facilitate ResponseWriter handling     | ⏳ 80%   |
| Develop method to Facilitate the handling of the Request | ⏳ 80%   |
| Develop Standard of Unit Testing                         | ⏳ 90%   |

---

## 🚀 Roadmap for Development  

| Task                                                                                            | Progress |
| ----------------------------------------------------------------------------------------------- | -------- |
| Documentation Tests Examples PKG Go                                                             | ⏳ 45%   |
| Test Coverage go test -cover                                                                    | ⏳ 74.6% |
| Regex feature coverage, but possibilities                                                       | 🔴 0%    |
| Develop for OPTIONS METHOD                                                                      | ✅ 100%  |
| Develop for CONNECT METHOD [See more](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect) | 🔴 0%    |
| Develop method for ListenAndServeTLS (http2)                                                    | 🔴 0%    |
| Develop Static Files support                                                                    | ✅ 100%  |
| WebSocket Support                                                                               | 🔴 0%    |
| Rate Limiter Support                                                                            | 🔴 0%    |
| Template Engines                                                                                | 🔴 0%    |
| Documentation Tests Examples PKG Go                                                             | ⏳ 45%   |
| Test coverage go test -cover                                                                    | ⏳ 75.5% |
| Coverage of Regex resources, but possibilities                                                  | 🔴 0%    |
| Develop for METHOD OPTIONS                                                                      | ✅ 100%  |
| Develop for CONNECT METHOD [See more](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect) | 🔴 0%    |
| Develop method for ListenAndServeTLS (http2)                                                    | 🔴 0%    |
| Create a CLI (Command Line Interface) Quick.                                                    | 🔴 0%    |

---

## 📊 Cover Testing Roadmap  

| Archive     | Coverage | Status |
| ----------- | -------- | ------ |
| Ctx         | 🟡 84.1% | 🟡     |
| Group       | ✅ 100.0% | 🟢     |
| Http Status | 🔴 7.8%  | 🔴     |
| Client      | 🟢 83.3% | 🟢     |
| Mock        | ✅ 100.0% | 🟢     |
| Concat      | ✅ 100.0% | 🟢     |
| Log         | 🔴 0.0%  | 🔴     |
| Print       | 🟡 66.7% | 🟡     |
| Qos         | 🔴 0.0%  | 🔴     |
| Rand        | 🔴 0.0%  | 🔴     |
| Compressa   | 🟡 71.4% | 🟡     |
| Cors        | 🟡 76.0% | 🟡     |
| Logger      | ✅ 100.0% | 🟢     |
| Maxbody     | ✅ 100.0% | 🟢     |
| Msgid       | ✅ 100.0% | 🟢     |
| Msguuid     | 🟢 86.4% | 🟢     |
| Quick       | 🟡 79.5% | 🟡     |
| QuickTest   | ✅ 100.0% | 🟢     |

### Fast quick example

When using New, you can configure global parameters such as request body limits, read/write time-out and route capacity. Below is a custom configuration:

```go
var ConfigDefault = Config{
	BodyLimit:         4 * 1024 * 1024,
	MaxBodySize:       4 * 1024 * 1024,
	MaxHeaderBytes:    2 * 1024 * 1024,
	RouteCapacity:     500,
	MoreRequests:      500,
	ReadTimeout:       5 * time.Second,
	WriteTimeout:      5 * time.Second,
	IdleTimeout:       2 * time.Second,
	ReadHeaderTimeout: 1 * time.Second,
}
```

Check out the code below:

```go
package main

import "github.com/jeffotoni/quick"
func main() {
	// Initialize a new Quick instance
    q := quick.New()

	// Define a simple GET route at the root path
    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action ❤️!")
    })

	/// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL
```bash
$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/user'

"Quick in action ❤️!"
```

### Quick Get Params
The example below defines a GET/v1/customer/:param1/:param2 endpoint, where :param1 and :param2 are variables on the route that can receive dynamic values.

```go
package main

import "github.com/jeffotoni/quick"

func main() {
    // Initialize a new Quick instance
    q := quick.New()

    // Define a GET route with two dynamic parameters (:param1 and :param2)
    q.Get("/v1/customer/:param1/:param2", func(c *quick.Ctx) error {
        // Set response content type to JSON
        c.Set("Content-Type", "application/json")

        // Define a struct for the response format
        type my struct {
            Msg string `json:"msg"` 
            Key string `json:"key"` 
            Val string `json:"val"` 
        }

        // Return a JSON response with extracted parameters
        return c.Status(200).JSON(&my{
            Msg: "Quick ❤️",
            Key: c.Param("param1"), 
            Val: c.Param("param2"), 
        })
    })

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL
```bash
$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/customer/val1/val2'

{
   "msg":"Quick ❤️",
   "key":"val1",
   "val":"val2"
}
```

### Quick Post Body json
This example shows how to create a POST endpoint to receive and process a JSON request body. The code reads the data sent in the request body and converts it to a Go structure.

```go
package main

import "github.com/jeffotoni/quick"

// Define a struct to map the expected JSON body
type My struct {
    Name string `json:"name"`
    Year int    `json:"year"`
}

func main() {
    // Initialize a new Quick instance
    q := quick.New()

    // Define a POST route to handle JSON request body
    q.Post("/v1/user", func(c *quick.Ctx) error {
        var my My

		// Parse the request body into the 'my' struct
        err := c.BodyParser(&my)
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }

        // Return the received JSON data
        return c.Status(200).String(c.BodyString())
        // Alternative:
        // return c.Status(200).JSON(&my)
    })

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL

```bash
$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v1/user' \
-d '{"name":"jeffotoni", "year":1990}'

{
   "name":"jeffotoni",
   "year":1990
}
```

## Uploads multipart/form-data

Quick provides a simplified API for managing uploads, allowing you to easily retrieve and manipulate files.

### ✅ **Main Methods and Functionalities**:
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

### 📌 File Upload Feature Comparison with other Frameworks

| Framework | `FormFile()` | `FormFiles()` | Dynamic Limit | Methods (`FileName()`, `Size()`)   | `Save()`, `SaveAll()` Method |
| --------- | ------------ | ------------- | ------------- | ---------------------------------- | ---------------------------- |
| **Quick** | ✅ Yes       | ✅ Yes        | ✅ Yes        | ✅ Yes                             | ✅ Yes                       |
| Fiber     | ✅ Yes       | ✅ Yes        | ❌ No         | ❌ No (uses `FileHeader` directly) | ✅ Yes                       |
| Gin       | ✅ Yes       | ✅ Yes        | ❌ No         | ❌ No (uses `FileHeader` directly) | ❌ No                        |
| Echo      | ✅ Yes       | ❌ No         | ❌ No         | ❌ No                              | ❌ No                        |
| net/http  | ✅ Yes       | ❌ No         | ❌ No         | ❌ No                              | ❌ No                        |

---
### File Upload Example
This example shows how to create a file upload endpoint. It allows users to send a single file via POST to the/upload route.

```go
package main

import (
    "fmt"
    "github.com/jeffotoni/quick"
)

// Define a struct for error messages
type Msg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

func main() {
    // Initialize a new Quick instance
    q := quick.New()

	// Define a route for file upload
    q.Post("/upload", func(c *quick.Ctx) error {
        // set limit upload
        c.FormFileLimit("10MB")

		// Retrieve the uploaded file
        uploadedFile, err := c.FormFile("file")
        if err != nil {
            return c.Status(400).JSON(Msg{
                Msg: "Upload error",
                Error: err.Error(),
             })
        }

		// Print file details
        fmt.Println("Name:", uploadedFile.FileName())
        fmt.Println("Size:", uploadedFile.Size())
        fmt.Println("MIME Type:", uploadedFile.ContentType())

        // Save the file (optional)
        // uploadedFile.Save("/tmp/uploads")

		// Return JSON response with file details
		// Alternative:
        //return c.Status(200).JSONIN(uploadedFile)
		return c.Status(200).JSON(map[string]interface{}{
			"name": uploadedFile.FileName(),
			"size": uploadedFile.Size(),
			"type": uploadedFile.ContentType(),
		})
		
    })

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL

```bash
$ curl -i -X POST http://localhost:8080/upload -F "file=quick.txt"

{
   "name":"quick.txt",
   "size":1109,
   "type":"text/plain; charset=utf-8"
}
```

### Multiple Upload Example
This example allows users to send multiple files via POST to the/upload-multiple route

```go
package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// Define a struct for error messages
type Msg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

func main() {
	// Initialize a new Quick instance
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

		// Alternative:
		//return c.Status(200).JSONIN(files)
		return c.Status(200).JSON("Upload successfully completed")
	})

	// Start the server on port 8080
	q.Listen("0.0.0.0:8080")
}
```

### 📌 cURL

```bash
$ curl -X POST http://localhost:8080/upload-multiple \
-F "files=@image1.jpg" -F "files=@document.pdf"

Upload successfully completed
```

### Quick Post Bind json

```go
package main

import "github.com/jeffotoni/quick"

// Define a struct to map the expected JSON body
type My struct {
    Name string `json:"name"`
    Year int    `json:"year"`
}

func main() {
	// Initialize a new Quick instance
    q := quick.New()

   // Define a POST route to handle JSON request body
	q.Post("/v2/user", func(c *quick.Ctx) error {
		var my My

		// Parse the request body into the 'my' struct
		err := c.Bind(&my)
		if err != nil {
			// Return a 400 status if JSON parsing fails
			return c.Status(400).SendString(err.Error())
		}

		// Return the received JSON data
		return c.Status(200).JSON(&my)
	})

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}
```

### 📌 cURL

```bash
$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v2/user' \
-d '{"name":"Marcos", "year":1990}'

{
   "name":"Marcos",
   "year":1990
}

```
### Cors

Using the Cors middleware, making your call in the default way, which is:

```go
var ConfigDefault = Config{
 AllowedOrigins: []string{"*"},
 AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
 AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
 ExposedHeaders: []string{"Content-Length"},
 AllowCredentials: false,
 MaxAge: 600,
 Debug: false,
}
```

Check out the code below:

```go
package main

import (
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cors"
)

func main() {

   // Initialize a new Quick instance
    q := quick.New()

    // Use the CORS middleware to allow cross-origin requests
    q.Use(cors.New())

    // Define a GET route at /v1/user
    q.Get("/v1/user", func(c *quick.Ctx) error {
        // Set the response content type to JSON
        c.Set("Content-Type", "application/json")

        // Return a response with status 200 and a message
        return c.Status(200).SendString("Quick in action com Cors❤️!")
    })

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL

```bash
$ curl -i -XGET -H "Content-Type:application/json" \
'http://localhost:8080/v1/user'

Quick in action com Cors❤️!
```


### Initializing Quick with Custom Configuration
This example demonstrates how to start a Quick server with a custom configuration.
```go
package main

import "github.com/jeffotoni/quick"

func main() {
    // Create a new Quick server instance with custom configuration
    q := quick.New(quick.Config{
        MaxBodySize: 5 * 1024 * 1024, // Set max request body size to 5MB
    })

    // Define a GET route that returns a JSON response
    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json") // Set response content type
        return c.Status(200).SendString("Quick in action com Cors❤️!") // Return response
    })

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL

```bash
$ curl -i -XGET -H "Content-Type:application/json" \
'http://localhost:8080/v1/user'

Quick in action com Cors❤️!
```

### Grouping Routes
This example demonstrates how to group routes using quick. Group(), making the code more organized
```go
package main

import "github.com/jeffotoni/quick"

func main() {
     // Create a new Quick server instance with custom configuration
    q := quick.New(quick.Config{
        MaxBodySize: 5 * 1024 * 1024, 
    })

    // Group for /v1 routes
    v1 := q.Group("/v1")

    // Define GET and POST routes for /v1/user
    v1.Get("/user", func(c *quick.Ctx) error {
        return c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
    })
    v1.Post("/user", func(c *quick.Ctx) error {
        return c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
    })

 	 // Group for /v2 routes
    v2 := q.Group("/v2")

    // Define GET and POST routes for /v2/user
    v2.Get("/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action com [GET] /v2/user ❤️!")
    })

    v2.Post("/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick in action com [POST] /v2/user ❤️!")
    })

	// Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL
1️⃣ GET /v1/user
```bash
$ curl -i -X GET http://localhost:8080/v1/user

[GET] [GROUP] /v1/user ok!!!
```
2️⃣ POST /v1/user
```bash
$ curl -i -X POST http://localhost:8080/v1/user

[POST] [GROUP] /v1/user ok!!!
```
3️⃣ GET /v2/user
```bash
$ curl -i -X GET http://localhost:8080/v2/user

Quick in action com [GET] /v2/user ❤️!
```
4️⃣ POST /v2/user
```bash
$ curl -i -X POST http://localhost:8080/v2/user

Quick in action com [POST] /v2/user ❤️!
```

### Quick Tests
This example demonstrates how to unit test routes in Quick using QuickTest().
It simulates HTTP requests and verifies if the response matches the expected output

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

    // Initialize Quick instance for testing
    q := quick.New()

    // Define test routes
    q.Post("/v1/user", testSuccessMockHandler)
    q.Post("/v1/user/:p1", testSuccessMockHandler)

    // Expected response data
    wantOutData := `"data":{"name":"jeff", "age":35}`
    reqBody := []byte(`{"name":"jeff", "age":35}`)
    reqHeaders := map[string]string{"Content-Type": "application/json"}

    // Perform test request
    data, err := q.QuickTest("POST", "/v1/user", reqHeaders, reqBody)
    if err != nil {
        t.Errorf("error: %v", err)
        return
    }

    // Compare expected and actual response
    s := strings.TrimSpace(data.BodyStr())
    if s != wantOutData {
        t.Errorf("Expected %s but got %s", wantOutData, s)
        return
    }

    // Log test results
    t.Logf("\nOutputBodyString -> %v", data.BodyStr())
    t.Logf("\nStatusCode -> %d", data.StatusCode())
    t.Logf("\nOutputBody -> %v", string(data.Body()))
    t.Logf("\nResponse -> %v", data.Response())
}

```

---

###  Regex

This example allows access only when the ID is numeric ([0-9]+).

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	// Initialize a new Quick instance
	q := quick.New()

	// Route that accepts only numeric IDs (using regex [0-9]+)
	q.Get("/users/{id:[0-9]+}", func(c *quick.Ctx) error {
		id := c.Param("id")
		return c.JSON(map[string]string{
			"message": "User found",
			"user_id": id,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}
```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/users/123

{
   "message":"User found",
   "user_id":"123"
}
```

### Accepts only lowercase letters in the slug
This example ensures that only lowercase letters ([a-z]+) are accepted in the slug

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	// Initialize a new Quick instance
	q := quick.New()

	// Route that accepts only lowercase slugs (words with lowercase letters)
		q.Get("/profile/{slug:[a-z]+}", func(c *quick.Ctx) error {
		slug := c.Param("slug")
		return c.JSON(map[string]string{
			"message": "Profile found",
			"profile": slug,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}
```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/profile/johndoe

{
   "message":"Profile found",
   "profile":"johndoe"
}
```

### Supports API version and numeric Id

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	// Initialize a new Quick instance
	q := quick.New()

	// Route that accepts an API version (v1, v2, etc.) and a numeric user ID
	q.Get("/api/{version:v[0-9]+}/users/{id:[0-9]+}", func(c *quick.Ctx) error {
		version := c.Param("version")
		id := c.Param("id")
		return c.JSON(map[string]string{
			"message": "API Versioned User",
			"version": version,
			"user_id": id,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}
```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/api/v1/users/123

{
   "message":"API Versioned User",
   "user_id":"123",
   "version":"v1"
}
```

## 🔑 Basic Authentication

**Basic Authentication (Basic Auth)** is a simple authentication mechanism defined in **RFC 7617**.  
It is commonly used for **HTTP-based authentication**, allowing clients to provide **credentials (username and password)** in the request header.


### 🔹 **How It Works**
1️⃣ The **client encodes** the username and password in Base64: 

2️⃣ The **encoded credentials** are sent in the `Authorization` header:  
```bash
Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=
```
3️⃣ The server decodes and verifies the credentials before granting access.

---

### **⚠️ Security Considerations**
- 🔒 Not encrypted – Basic Auth only encodes credentials in Base64, but does not encrypt them.
- 🔐 Use over HTTPS – Always use Basic Auth with TLS/SSL (HTTPS) to prevent credentials from being exposed.
- 🔑 Consider alternatives – For stronger security, prefer OAuth2, JWT, or API keys.
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

## 📂 **STATIC FILES**

A **Static File Server** is an essential feature in web frameworks, enabling the efficient serving of static content such as **HTML, CSS, JavaScript, images, and other assets**.  

It is particularly useful for:  
- ✅ Hosting front-end applications  
- ✅ Providing downloadable files  
- ✅ Serving resources directly from the backend  

---

### 🔹 **How It Works**
- 1️⃣ The server **listens** for HTTP requests targeting static file paths.  
- 2️⃣ If the requested file **exists** in the configured directory, the server reads and returns it as a response.  
- 3️⃣ **MIME types** are automatically determined based on the file extension for correct rendering.  

---

### ⚡ **Key Features**
- 🚀 **Efficient Handling** – Serves files directly without additional processing.  
- 🎯 **MIME Type Detection** – Automatically identifies file types for proper rendering.  
- ⚡ **Caching Support** – Can be configured to improve performance via HTTP headers.  
- 📂 **Optional Directory Listing** – Allows browsing available static files (if enabled).  

---

### ⚠️ **Security Considerations**
- 🔒 **Restrict Access** – Prevent exposure of sensitive files like `.env`, `.git`, or configuration files.  
- 🌐 **CORS Policies** – Configure **Cross-Origin Resource Sharing (CORS)** to control file access.  
- 🛡 **Content Security Policy (CSP)** – Helps mitigate **XSS (Cross-Site Scripting)** risks.  

---

By properly configuring your static file server, you can ensure **fast, efficient, and secure delivery of resources**! 🚀🔥

---

### Serving Static Files with Quick Framework

This example sets up a basic web server that serves static files, such as HTML, CSS, or JavaScript.

```go
package main

import "github.com/jeffotoni/quick"

func main() {

    //starting Quick
    q := quick.New()

    // Static Files Setup
    q.Static("/static", "./static")

    // Route Definition
    q.Get("/", func(c *quick.Ctx) error {
        c.File("./static/index.html")
        return nil
    })

    // Start the server on port 8080
    q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/

File Server Go example html
```
---

### 📁 EMBED - Embedded Static Files

### 🔹 How do embedded static files work?
1. Static assets (HTML, CSS, JS, images, etc.) are **compiled directly into the binary** at compile time, using the Go package `embed`.
2. The application **serves these files from memory**, eliminating the need to access the disk.
3. This **removes external dependencies**, making the deployment simpler and more efficient.

---
### ⚡ Advantages of using embedded files:
- ✅ **Portability** - The binary contains everything you need, no extra files.  
- ✅ **Performance** - File access is faster because files are already loaded in memory.  
- ✅ **Security** - Reduces exposure to attacks because the file system does not need to be accessible.  

---
### 🚀 How does Quick simplify this process?
The function `q. Static()` already handles the complexity of serving embedded files. Just call it with `embed.FS`.

The example shows how to serve static files with Quick and `embed.FS`

```go
package main

import (
	"embed"

	"github.com/jeffotoni/quick"
)

//go:embed static
var staticFiles embed.FS

func main() {
	// Server Initialization
	q := quick.New()

	// Static Files Setup (serves files automatically from embed.FS)
	q.Static("/static", staticFiles)

	// Defines a route that serves the HTML index file
	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/index.html") 	
		return nil
	})

	// Start the server on port 8080
	q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/

File Server Go example html
```

### 📂 Example Project Structure
```text
quick-example
│── main.go
│── static/
│   ├── index.html
│   ├── style.css
│   ├── script.js
```

---

## 🌍 HTTP Client

The **HTTP Client** package in **Quick** provides a **simple and flexible** way to make HTTP requests, supporting **GET, POST, PUT, and DELETE** operations. 🚀

It is designed to handle different types of request bodies and parse responses easily.

###  🎯 Why Use Quick's HTTP Client?
- ✅ Easy-to-Use – Simplified functions for common HTTP requests.
- ✅ Highly Customizable – Supports headers, authentication, and transport settings.
- ✅ Flexible Body Parsing – Works with JSON, plain text, and custom io.Reader types.
- ✅ Automatic JSON Handling – No need to manually marshal/unmarshal JSON.


### ⚡ Key Features

- 🔹 **Convenience Functions** – Use `Get`, `Post`, `Put`, and `Delete` to make quick requests with a default client.
- 🔹 **Customizable Requests** – Easily add headers, authentication, and request settings.
- 🔹 **Automatic JSON Processing** – Seamless encoding and decoding of JSON data.
- 🔹 **Flexible Request Body** – Send data as **JSON**, **plain text**, or any `io.Reader`.

---
### 📌 **Client Structure**
The `Client` struct represents a configurable HTTP client with advanced features:

```go
var ClientDefault = Client{
	Ctx:          context.Background(),
	ClientHTTP:   httpGoClient{},
	Headers:      map[string]string{"Content-Type": "application/json"},
	EnableLogger: true,
	Logger:       slog.Default(),
}
```
Check out the code below:

### GET Request Example
A GET request is used to retrieve data from a server.

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create a new HTTP client
	httpClient := client.New()

	// Making a GET request to fetch a list of users
	resp, err := httpClient.Get("https://reqres.in/api/users")
	if err != nil {
		log.Fatal(err)
	}

	// Alternative:
	//fmt.Println("GET response:", string(resp.Body))

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	// Extract first user
	users := result["data"].([]interface{})
	firstUser := users[0].(map[string]interface{})

	// Print only first user
	fmt.Printf("Id: %v\n", firstUser["id"])
	fmt.Printf("Name: %v %v\n", firstUser["first_name"], firstUser["last_name"])
	fmt.Printf("Email: %v\n", firstUser["email"])
	fmt.Printf("Avatar: %v\n", firstUser["avatar"])
}
```
### 📌 Response
```bash
Id: 1
Name: George Bluth
Email: george.bluth@reqres.in
Avatar: https://reqres.in/img/faces/1-image.jpg
```

### POST Request Example (Using a Struct)
A POST request is used to send data to a server, often for creating new resources.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create a new HTTP client
	httpClient := client.New()

	// Define a struct for the request body
	data := struct {
		Name string `json:"name"`
	}{
		Name: "Emma",
	}

	// Making a POST request with JSON data
	resp, err := httpClient.Post("https://reqres.in/api/users", data)
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	// Alternative:
	//fmt.Println("POST response:", result)

	// Print formatted response
	fmt.Println("Id:", result["id"])
	fmt.Println("Created_At:", result["createdAt"])
}
```
### 📌 Response
```bash
Id: 322
Created_At: 2025-03-14T14:48:24.305Z
```

### PUT Request Example (Using a String)
A PUT request is used to update an existing resource.
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create a new HTTP client
	httpClient := client.New()

	// Define a struct with updated user data
	data := struct {
		Name string `json:"name"`
	}{
		Name: "Jeff",
	}

	// PUT request to ReqRes API
	resp, err := httpClient.Put("https://reqres.in/api/users/2", data)
	if err != nil {
		log.Fatal("Error making PUT request:", err)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	// Alternative: Print the HTTP status and response body
	// fmt.Println("HTTP Status Code:", resp.StatusCode)
	// fmt.Println("Raw Response Body:", string(resp.Body))

	// Print formatted response
	fmt.Println("Updated_At:", result["updatedAt"])
}
```
### 📌 Response
```bash
Updated_At: 2025-03-14T14:56:35.202Z
```

### DELETE Request Example
A DELETE request is used to remove a resource from the server.
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// Create a new HTTP client
	httpClient := client.New()

	// DELETE request to ReqRes API
	resp, err := httpClient.Delete("https://reqres.in/api/users/2")
	if err != nil {
		log.Fatal("Error making request:", err)
	}

	// Print the HTTP status to confirm deletion
	fmt.Println("Status Code:", resp.StatusCode)

	// Since DELETE usually returns no content, we check if it's empty
	if len(resp.Body) > 0 {
		fmt.Println("Raw Response Body:", string(resp.Body))
	} else {
		fmt.Println("Response Body is empty (expected for 204 No Content)")
	}
}
```
### 📌 Response
```bash
Status Code: 204
Response Body is empty (expected for 204 No Content)
```
---

## 🚀 Qtest - HTTP Testing Utility for Quick

Qtest is an **advanced HTTP testing function** designed to simplify route validation within the **Quick** framework. It enables seamless testing of simulated HTTP requests using `httptest`, supporting:

- **Custom HTTP methods** (`GET`, `POST`, `PUT`, `DELETE`, etc.).
- ✅ **Custom headers**.
- ✅ **Query parameters**.
- ✅ **Request body**.
- ✅ **Cookies**.
- ✅ **Built-in validation methods** for status codes, headers, and response bodies.

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
### 📌 Usage Reference

| Function                        | Description                          |
|---------------------------------|--------------------------------------|
| `Qtest(opts QuickTestOptions)`  | Executes an HTTP test request       |
| `AssertStatus(expected int)`    | Asserts expected HTTP status code   |
| `AssertHeader(key, value string)` | Checks response header value        |
| `AssertBodyContains(substr string)` | Verifies if body contains a string |


### 📖 More Details
🔗 **Check out the full documentation:** [Qtest - Quick](https://github.com/jeffotoni/quick/tree/main/quickTest)

---

## 🔄 Retry & Failover Mechanisms in Quick HTTP Client

The **Quick HTTP Client** now includes **built-in retry and failover support**, allowing for more resilient and reliable HTTP requests. These features are essential for handling **transient failures**, **network instability**, and **service downtime** efficiently.

### 🚀 Key Features

- 🔄 **Automatic Retries**: Retries failed requests based on configurable rules.
- ⏳ **Exponential Backoff**: Gradually increases the delay between retry attempts.
- 📡 **Status-Based Retries**: Retries only on specified HTTP status codes (e.g., `500`, `502`, `503`).
- 🌍 **Failover Mechanism**: Switches to predefined backup URLs if the primary request fails.
- 📑 **Logging Support**: Enables detailed logs for debugging retry behavior.

---

### 🔹 How Retry & Failover Work

The **retry mechanism** automatically **resends requests** when they fail, with configurable options to:

- **Limit the number of retries** to avoid excessive attempts.
- **Introduce backoff delays** to prevent overwhelming the server.
- **Retry only on specific HTTP status codes** (e.g., `500`, `502`, `503`).

The **failover system** ensures **high availability** by redirecting failed requests to **predefined backup URLs**, reducing downtime and improving system resilience.

### ⚙️ Configuration Options
These options allow fine-grained control over retry and failover behavior:

| Option          | Description 🚀 |
|----------------|----------------------------------------------------------------|
| **MaxRetries**  | Sets the **maximum number of retry attempts** before failure. |
| **Delay**       | Defines the **initial delay** before retrying a request. |
| **UseBackoff**  | Enables **exponential backoff**, increasing delay dynamically after each retry. |
| **Statuses**    | List of **HTTP status codes** (e.g., `500`, `502`, `503`) that trigger a retry. |
| **FailoverURLs** | List of **backup URLs** used if the primary request repeatedly fails. |
| **EnableLog**   | Enables **detailed logging** for debugging retry behavior. |

---

### ⚙️ **Configuring Retries & Failover in Quick HTTP Client**


The **Quick HTTP Client** provides built-in support for **retrying failed requests** and **switching to failover URLs** when necessary.  

You can configure these behaviors using the **`WithRetry`** option, which accepts a `RetryConfig` struct.

### **🛠 Creating a Custom Client with Retries**
The following example shows how to create a Quick client with **retry and failover mechanisms**.

```go
// Creating a Quick client using a custom HTTP client and retry settings.
cClient := client.New(
    client.WithCustomHTTPClient(customHTTPClient), 
    client.WithContext(context.Background()),      
 	client.WithHeaders(map[string]string{
			"Content-Type": "application/json", 
		}),
    client.WithRetry(client.RetryConfig{
        MaxRetries:   3,                      
        Delay:        2 * time.Second,        
        UseBackoff:   true,                
        Statuses:     []int{500, 502, 503, 504}, 
        FailoverURLs: []string{"http://hosterror", "https://httpbin.org/post"},
        EnableLog:    true,                  
    }),
)
```
### 🔄 Smart Retries: Exponential Backoff & Failover to Backup URLs
This example demonstrates retrying a request with an increasing delay (backoff) when encountering errors.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Create a new Quick HTTP client with retry settings
	cClient := client.New(
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 3,       
				Delay:      1 * time.Second,  
				UseBackoff: true,   
				Statuses:   []int{500, 502, 503}, 
				FailoverURLs: []string{"http://backup1.com/resource", "https://httpbin_error.org/get", "https://httpbin.org/get"},
				EnableLog:  true,   
			}),
	)

	// Send a GET request to the specified URL
	resp, err := cClient.Get("https://httpbin_error.org/get")
	if err != nil {
		log.Fatal("GET request failed:", err)
	}

	// Print the response body
	fmt.Println("GET Response:", string(resp.Body))
}
```
### 📌 Response
```bash
{"time":"2025-03-14T14:27:02.069237664-03:00","level":"WARN","msg":"Retrying request","url":"https://httpbin_error.org/get","method":"GET","attempt":1,"failover":1}
{"time":"2025-03-14T14:27:13.076907091-03:00","level":"WARN","msg":"Retrying request","url":"http://backup1.com/resource","method":"GET","attempt":2,"failover":2}
{"time":"2025-03-14T14:27:15.258544931-03:00","level":"WARN","msg":"Retrying request","url":"https://httpbin_error.org/get","method":"GET","attempt":3,"failover":3}
GET Response: {
  "args": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Host": "httpbin_error.org", 
    "User-Agent": "Go-http-client/1.1", 
    "X-Amzn-Trace-Id": "Root=1-67d466f8-1aafed0512167ac32426bc9f"
  }, 
  "origin": "179.216.110.129", 
  "url": "https://httpbin_error.org/get"
}
```

### 🔥 How Each Feature Works in the Code  

### 🔄 **Retry with Exponential Backoff**  
- The retry mechanism is triggered because `MaxRetries: 3` allows the request to be retried up to **three times**.  
- The wait time between attempts **automatically increases** due to `UseBackoff: true`.  
- A retry **only occurs** if the response contains an HTTP error listed in `Statuses: []int{500, 502, 503}`.  

### 🌍 **Failover to Backup URLs**  
- If **all retry attempts on the primary URL fail**, the client will try the **alternative URLs** listed in `FailoverURLs`.  
- In this example, if `https://httpbin.org/status/500` keeps failing, it will switch to `https://httpbin.org/get`.  



---

## 📝 Form Submission with `PostForm` in Quick HTTP Client  

The **Quick HTTP Client** now supports **`PostForm`**, making it easier to send **form-encoded data (`application/x-www-form-urlencoded`)**.  
This feature is particularly useful for:  

✅ **Authentication requests**  
✅ **Submitting data to web services**  
✅ **Integrations with legacy systems that do not accept JSON**  

---

## 🔹 Why Use `PostForm`?  

| 🚀 **Feature**              | 💡 **Benefit** |
|-----------------------------|------------------------------------------------|
| **📑 Optimized for Forms**    | Makes it easy to send form-encoded data (`application/x-www-form-urlencoded`). |
| **⚙️ Automatic Encoding**    | Converts `url.Values` into a valid form submission format. |
| **📌 Header Management**     | Automatically sets `Content-Type: application/x-www-form-urlencoded`. |
| **🔄 Consistent API**        | Follows the same design as `Post`, `Get`, `Put`, ensuring ease of use. |
| **🌍 Better Compatibility**  | Works seamlessly with APIs that do not accept JSON payloads. |


---

### 🔹 How PostForm Works

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
	// Initialize Quick framework
	q := quick.New()

	// Define a POST route to handle form data submission
	q.Post("/postform", func(c *quick.Ctx) error {
		// Retrieve form values from the request
		form := c.FormValues()

		// Return the received form data as JSON response
		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Start the Quick server in a separate goroutine
	go func() {
		fmt.Println("Quick server running at http://localhost:3000")
		if err := q.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Quick server: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	// Create an HTTP client before calling PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second),
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}),
	)

	// Declare form data (key-value pairs)
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Send a POST request with form data
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request failed: %v", err)
	}

	// Unmarshal the JSON response from the server
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal("Failed to parse JSON response:", err)
	}

	// Print the formatted JSON response
	// Alternative:fmt.Println("POST Response:", result)

	// Print the formatted JSON response
	formattedResponse, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("Failed to format JSON response:", err)
	}

	fmt.Println("POST Response:")
	fmt.Println(string(formattedResponse))

}
```
### 📌 Response
```bash
POST Response:
{
  "data": {
    "password": [
      "supersecret"
    ],
    "username": [
      "quick_user"
    ]
  },
  "message": "Received form data"
}
```
---

## 🌐 Transport Configuration in HTTP Client

The `Transport` setting in the Quick HTTP Client **manages the network layer**, ensuring efficient, secure, and reliable HTTP communications.
It provides **fine-grained** control over **connection management, security settings, and protocol optimizations** for both **development and production environments**.

### ✅ Key Features of Transport Configuration  

| ⚙️ **Setting**              | 🔍 **Description**  |
|-----------------------------|----------------------------------------------------------------------------------------------------------------|
| **🛠 Proxy Settings**       | Handles **proxy servers** using system environment settings for automatic configuration.  |
| **🔒 TLS Configuration**    | Controls **security settings**, such as **TLS version** and **certificate verification**. `InsecureSkipVerify` can be enabled for development to bypass SSL verification.  |
| **📡 Connection Management** | Optimizes resource usage with settings like `MaxIdleConns`, `MaxConnsPerHost`, and `MaxIdleConnsPerHost`, improving **scalability**.  |
| **🚀 Persistent Connections** | Enables or disables **Keep-Alives**, reducing connection setup time and improving performance. |
| **⚡ HTTP/2 Support**        | Enables **HTTP/2** for faster, more efficient communication when supported by the server. |
 

#### Advanced HTTP client configuration with failover mechanism

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
			// Allows insecure TLS connections .
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
	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{"message": "Quick in action"})
	if err != nil {
		log.Fatal(err)
	}

	// show resp
	fmt.Println("POST response:\n", string(resp.Body))
}
```
### 📌 Response
```bash
{"time":"2025-03-14T15:31:11.027180616-03:00","level":"WARN","msg":"Retrying request","url":"https://httpbin_error.org/post","method":"POST","attempt":1,"failover":1}
{"time":"2025-03-14T15:31:12.028294877-03:00","level":"WARN","msg":"Retrying request","url":"http://hosterror","method":"POST","attempt":2,"failover":2}
POST response:
 {
  "args": {}, 
  "data": "{\"message\":\"Quick in action\"}", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Authorization": "Bearer YOUR_ACCESS_TOKEN", 
    "Content-Length": "29", 
    "Content-Type": "application/json", 
    "Host": "httpbin_error.org", 
    "User-Agent": "Go-http-client/1.1", 
    "X-Amzn-Trace-Id": "Root=1-67d475f2-713b9e4c2fff65d413fcd097"
  }, 
  "json": {
    "message": "Quick in action"
  }, 
  "origin": "179.216.110.129", 
  "url": "https://httpbin_error.org/post"
}
```
---

#### HTTP Client with Advanced Transport and Failover Capabilities

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
### 📌 Response
```bash
{"time":"2025-03-14T15:37:43.481220287-03:00","level":"WARN","msg":"Retrying request","url":"https://httpbin_error.org/post","method":"POST","attempt":1,"failover":1}
{"time":"2025-03-14T15:37:44.482388761-03:00","level":"WARN","msg":"Retrying request","url":"http://hosterror","method":"POST","attempt":2,"failover":2}
POST response: {
  "args": {}, 
  "data": "{\"name\":\"jeffotoni\"}", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Authorization": "Bearer YOUR_ACCESS_TOKEN", 
    "Content-Length": "20", 
    "Content-Type": "application/json", 
    "Host": "httpbin_error.org", 
    "User-Agent": "Go-http-client/1.1", 
    "X-Amzn-Trace-Id": "Root=1-67d4777b-50d494284d3d242224dc62c0"
  }, 
  "json": {
    "name": "jeffotoni"
  }, 
  "origin": "179.216.110.129", 
  "url": "https://httpbin_error.org/post"
}
```

---

### Configuring HTTP Client with Retry and Failover Mechanisms

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
### 📌 Response
```bash
{"time":"2025-03-14T15:40:30.617507958-03:00","level":"WARN","msg":"Retrying request","url":"https://httpbin_error.org/post","method":"POST","attempt":1,"failover":1}
{"time":"2025-03-14T15:40:31.618144855-03:00","level":"WARN","msg":"Retrying request","url":"http://hosterror","method":"POST","attempt":2,"failover":2}
POST Form Response: {
  "args": {}, 
  "data": "{\"name\":\"jeffotoni in action with Quick!!!\"}", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "44", 
    "Content-Type": "application/json", 
    "Host": "httpbin_error.org", 
    "User-Agent": "Go-http-client/1.1", 
    "X-Amzn-Trace-Id": "Root=1-67d47822-5c80648f5a30c75c6a500470"
  }, 
  "json": {
    "name": "jeffotoni in action with Quick!!!"
  }, 
  "origin": "179.216.110.129", 
  "url": "https://httpbin_error.org/post"
}
```

---

### Advanced HTTP Client Configuration with Transport and Retry Settings

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
### 📌 Response
```bash
POST Form Response: {
  "args": {}, 
  "data": "{\"name\":\"jeffotoni\"}", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "20", 
    "Content-Type": "application/json", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/2.0", 
    "X-Amzn-Trace-Id": "Root=1-67d4786f-61ddc079287866e673f4f584"
  }, 
  "json": {
    "name": "jeffotoni"
  }, 
  "origin": "179.216.110.129", 
  "url": "https://httpbin.org/post"
}
```

---

## 🔐 TLS (Transport Layer Security) in Quick HTTP Server

`TLS (Transport Layer Security)` is a cryptographic protocol that provides **secure communication** over a network.
It is widely used to encrypt data transmitted between clients and servers, ensuring **confidentiality, integrity, and authentication**. 
TLS is the **successor to SSL (Secure Sockets Layer)** and is used in **HTTPS, email security, and many other applications**.

---

### 🚀 Key TLS Features

| 🔹 **Feature**        | 🔍 **Description**  |
|----------------------|--------------------------------------------------------------|
| 🔐 **Encryption**     | Protects data from being intercepted during transmission. |
| 🔑 **Authentication** | Ensures the server (and optionally the client) is legitimate. |
| 🔄 **Data Integrity** | Prevents data from being modified or tampered with in transit. |
| ⚡ **Performance**    | Modern TLS versions (1.2, 1.3) offer strong security with minimal overhead. |

---

### 🌍 Running a Secure HTTPS Server with Quick and TLS
This example demonstrates how to set up an HTTPS server using Quick with TLS encryption, ensuring secure communication between clients and the server.

```go
package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

func main() {
	// Initialize Quick instance
	q := quick.New()

	// Print a message indicating that the server is starting on port 8443
	fmt.Println("Run Server port:8443")

	// Start the HTTPS server with TLS encryption
	// - The server will listen on port 8443 (non-privileged port)
	// - cert.pem: SSL/TLS certificate file
	// - key.pem: Private key file for SSL/TLS encryption
	err := q.ListenTLS(":8443", "cert.pem", "key.pem", false)
	if err != nil {
		// Log an error message if the server fails to start
		fmt.Printf("Error when trying to connect with TLS: %v\n", err)
	}
}

```

### ⚠️Ports & Permissions Considerations

This example **uses port 8443** so that it runs on **any operating system without requiring extra permissions**.

However, in production, you may want to use the **standard HTTPS port 443**.

- **Port 443** (default for HTTPS) is a **privileged port** (below 1024).
- On **Linux**, running a service on port 443 requires **superuser privileges**.

To run on **port 443** on Linux, use:

```bash
$ sudo go run main.go
```
---

## 🚦 Rate Limiter - Request Limiting Middleware

The **Rate Limiter** is a middleware for the Quick framework that controls the number of requests allowed in a given time period. It helps prevent API abuse and improves system stability by preventing server overload.

### 🚀 Features

| Feature                         | Description |
|----------------------------------|-------------|
| 🎯 **Request Rate Limiting**     | Configurable maximum number of requests per client within a time window. |
| ⏳ **Automatic Expiration**      | Resets the request counter automatically after the configured time. |
| 🔑 **Custom Client Identification** | Uses a `KeyGenerator` function to define a unique client key (e.g., IP-based). |
| ⚠️ **Custom Response on Limit Exceeded** | Allows defining a custom response when the request limit is reached. |
| ⚡ **Efficient Performance**     | Implements sharding and optimizations to reduce concurrency issues. |

### 🌍 Global Rate Limiter 

The example below shows how to apply the Rate Limiter as global middleware.
```go
package main

import (
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/limiter"
)

func main() {
	q := quick.New()

	// Apply the rate limiter middleware
	q.Use(limiter.New(limiter.Config{
		// Maximum 10 requests allowed per IP
		Max: 10,
		// The limit resets after 5 seconds
		Expiration: 5 * time.Second,
		KeyGenerator: func(c *quick.Ctx) string {
			// Uses the client's IP address as the key
			return c.RemoteIP()
		},
		LimitReached: func(c *quick.Ctx) error {
			c.Set("Content-Type", "application/json")
			// The client should wait 10 seconds before retrying
			c.Set("Retry-After", "10")
			return c.Status(quick.StatusTooManyRequests).JSON(map[string]string{
				"error":   "Too many requests",
				"message": "You have exceeded the request limit. 
				Please wait 1 second and try again.",
				"retry_after": "10s",
			})
		},
	}))

	// Define a simple GET route
	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).JSON(map[string]string{"msg": "Quick in action ❤️!"})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}


```
### 📌 cURL

```bash
$ curl -i -X GET http://localhost:8080/
```
If the same IP makes more than 10 requests in 5 seconds, the middleware returns:

```bash
{
    "error": "Too many requests",
    "message": "You have exceeded the request limit. 
	 Please wait 1 second and try again.",
    "retry_after": "10s"
}
```

### 🌐 Using Rate Limiter for Specific Route Groups

The example below shows how to apply the Rate Limiter with route group.

```go
package main

import (
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/limiter"
)

func main() {
	// Create a new Quick instance
	q := quick.New()

	// Rate Limiter Middleware
	limiterMiddleware := limiter.New(limiter.Config{
		// Maximum 3 requests allowed per IP address within a 10-second window
		Max: 3,
		// The limit resets every 10 seconds
		Expiration: 10 * time.Second,
		// Use the client's IP address as the unique key to track rate limits
		KeyGenerator: func(c *quick.Ctx) string {
			return c.RemoteIP()
		},
		// If the rate limit is exceeded, send an error message and instructions
		LimitReached: func(c *quick.Ctx) error {
			// Set content type to JSON
			c.Set("Content-Type", "application/json")
			c.Set("Retry-After", "10") 
			// Response structure
			response := map[string]string{
				"error":       "Too many requests",
				"message":     "You have exceeded the request limit. 
				Please wait 10 seconds and try again.",
				"retry_after": "10s",
			}

			// Log to verify that the rate limit exceeded response is being sent
			log.Println("Rate Limit Exceeded:", response)

			// Return the response with HTTP status 429 (Too Many Requests)
			return c.Status(quick.StatusTooManyRequests).JSON(response)
		},
	})

	// Create an API group with rate limit middleware
	api := q.Group("/v1")
	// Apply the rate limiter middleware to the /api group
	api.Use(limiterMiddleware)

	// Define route /api/users that responds with a list of users
	api.Get("/users", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "list of users"})
	})

	// Define route /api/posts that responds with a list of posts
	api.Get("/posts", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "list of posts"})
	})

	// Define route without rate limit
	q.Get("/", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "Quick in action ❤️!"})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}
```
### 📌 cURL

```bash
$ curl -i -X GET http://localhost:8080/users
```
If an IP makes more than 3 requests within 10 seconds, the response is blocked and returns a 429 Too Many Requests error:

```bash
{
    "error": "Too many requests",
    "message": "You have exceeded the request limit. 
	Please wait 10 seconds and try again.",    "retry_after": "10s"
}

```
---
## 🚀 Benchmarks

**Benchmarking** is a performance evaluation technique that measures response time, resource usage, and processing capacity. It helps developers identify bottlenecks and optimize their code. This approach is widely used in various areas, including software testing and hardware evaluations.


### 📊 Why Benchmarking?

| ✅ Benefit   | 🔍 Description |
|--------------------------------------|-------------|
| 📏 Measure performance   | Evaluates how a system responds under different workloads. |
| 🔄 Compare technologies   | Allows you to analyze different frameworks, libraries or implementations. |
| 🔍 Identify bottlenecks   | Helps detect critical points that need optimization. |
| 📈 Ensure scalability   | Test system behavior with multiple simultaneous requests. |
| 🎭 Simulate real-world scenarios   | Reproduces heavy use situations, such as thousands of users accessing a service at the same time. |


### 🛠️ Load Testing with Quick and k6

To evaluate the performance of our API, we conducted a benchmark test using the [Quick](https://github.com/jeffotoni/quick) framework along with [k6](https://k6.io/) for load testing.


### 🔹 API Code for Benchmarking

The following Go API was used for benchmarking. It provides a POST endpoint at /v1/user, which:

- Accepts large JSON payloads.
- Parses the incoming JSON into a Go struct.
- Returns the parsed JSON as a response.

```go
package main

import (
	"github.com/jeffotoni/quick"
)

// Struct representing a user model
type My struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Year     int                    `json:"year"`
	Price    float64                `json:"price"`
	Big      bool                   `json:"big"`
	Car      bool                   `json:"car"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Options  []Option               `json:"options"`
	Extra    interface{}            `json:"extra"`
	Dynamic  map[string]interface{} `json:"dynamic"`
}

type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	// Initialize Quick framework with a 20MB body limit
	q := quick.New(quick.Config{
		MaxBodySize: 20 * 1024 * 1024,
	})

	// Define a POST route at /v1/user
	q.Post("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		var users []My // Store incoming user data

		// Parse the request body into the struct
		err := c.Bind(&users)
		if err != nil {
			// If parsing fails, return a 400 Bad Request response
			return c.Status(400).SendString(err.Error())
		}

		// Return the parsed JSON data as a response with 200 OK
		return c.Status(200).JSON(users)
	})

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}
```
### 🏆 k6 Load Test Script

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

// Load the JSON from the environment variable
const payloadData = open('./data_1k_list.json');

// K6 configuration
export let options = {
    stages: [
        { duration: '40s', target: 1000 }, // Ramp-up para 500 VUs
        { duration: '7s', target: 500 },  // Mantém 500 VUs
        { duration: '5s', target: 0 },   // Ramp-down
    ],
 };


export default function () {
let url = 'http://localhost:8080/v1/user';

// Always use the same list for sending
// let payload = JSON.stringify(payloadData);

let params = {
headers: { 'Content-Type': 'application/json' },
};

let res = http.post(url, payloadData, params);

// Check if the response is correct
check(res, {
'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
'response contains JSON': (r) => r.headers['Content-Type'] === 'application/json',
});

}
```

### 📈 Running the Tests  

- 1️⃣ Start the Quick API - Run the Quick server:  
```bash
$ go run main.go
```
- 2️⃣ Execute the Load Test
```bash
$ k6 run benchmark.js
```
---
## 📦 Compression Middleware (compress)
The compress middleware in Quick enables automatic GZIP compression for HTTP responses, reducing the size of data transferred over the network. This improves performance and bandwidth efficiency, especially for text-based content like JSON, HTML, and CSS.

### 🚀 Benefits of Compression
- ✅ Reduced response size – improves loading speed.
- ✅ Bandwidth savings – ideal for mobile or limited connections.
- ✅ Seamless integration – works automatically for compatible clients.
- ✅ Better user experience – faster response times.

🔹 Ways to Use
Quick provides three different ways to enable GZIP compression:

### 🌟 Available Usage Methods

Quick provides three different ways to enable GZIP compression:

- Using quick.Handler (Default) – Follows Quick’s native syntax.
- Using quick.HandlerFunc – Alternative method for direct function-based handlers.
- Using net/http standard implementation – For applications using Go’s native HTTP package.

#### For more details on using compress, check the documentation:

<h4 align="left"> 
	<p> 
		<a href="middleware/README.md">
			<strong>📖 Compress Documentation</strong>
		</a>
	</p> 
</h4>

### 🚀 Usage Example (Default)
Here is a practical example of enabling the GZIP middleware in Quick using the default approach (quick.Handler)

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	q := quick.New()

	// Enable Gzip middleware
	q.Use(compress.Gzip())

	// Define a route that returns a compressed JSON response
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		// Setting response headers
		c.Set("Content-Type", "application/json")
		// Enabling Gzip compression
		c.Set("Accept-Encoding", "gzip") 
		// Defining the response structure
		type response struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		// Returning a JSON response with headers
		return c.Status(200).JSON(&response{
			Msg:     "Quick ❤️",
			Headers: c.Headers,
		})
	})

	// Start the HTTP server on port 8080
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
### 📌 cURL
```bash
$ curl -X GET http://localhost:8080/v1/compress -H 
"Accept-Encoding: gzip" --compressed -i
```
### 📌 Response
```bash
{
   "msg":"Quick ❤️",
   "headers":{
      "Accept":[
         "*/*"
      ],
      "Accept-Encoding":[
         "gzip"
      ],
      "Cache-Control":[
         "no-cache"
      ],
      "Connection":[
         "keep-alive"
      ],
      "Postman-Token":[
         "e0b65cfe-9516-4803-96df-d443d7e6a95a"
      ],
      "User-Agent":[
         "PostmanRuntime/7.43.2"
      ]
   }
}
```
---

## 📏 Maxbody (Request Size Limiter)

The maxbody middleware restricts the maximum request body size to prevent clients from sending excessively large payloads. This helps optimize memory usage, enhance security, and avoid unnecessary processing of oversized requests.

### 🔹 Why Use maxbody?
- ✅ Prevents excessive memory usage and improves performance.
- ✅ Mitigates DoS (Denial-of-Service) attacks by limiting large payloads.
- ✅ Automatically returns 413 Payload Too Large when the limit is exceeded.

---
### 🚀 Ways to Use maxbody
There are two primary ways to enforce request body size limits in Quick:

**maxbody.New()** – Enforces a global request body size limit across all middleware layers.

**MaxBytesReader()** – Adds an extra layer of validation inside a specific request handler

### 🔹 Using maxbody.New()
This example applies a global request body limit of 50KB for all incoming requests.

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/maxbody"
)

func main() {
	q := quick.New()

	// Middleware to enforce a 50KB request body limit globally
	q.Use(maxbody.New(50 * 1024)) // 50KB

	// Define a route that accepts a request body
	q.Post("/v1/user/maxbody/any", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		log.Printf("Body received: %s", c.BodyString())
		return c.Status(200).Send(c.Body())
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
### 📌 cURL
Request within limit (Valid request)

```bash
$ curl -X POST http://0.0.0.0:8080/v1/user/maxbody/any \
     -H "Content-Type: application/json" \
     --data-binary @<(head -c 48000 </dev/zero | tr '\0' 'A')
```
Request exceeding limit (Should return 413)
```bash
$ curl -X POST http://0.0.0.0:8080/v1/user/maxbody/any \
     -H "Content-Type: application/json" \
     --data-binary @<(head -c 51000 </dev/zero | tr '\0' 'A')
```
---
### 🔹 Using MaxBytesReader()

This example adds extra protection by applying MaxBytesReader() inside the request handler, ensuring an enforced limit at the application layer.

```go
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

const maxBodySize = 1024 // 1KB

func main() {
	q := quick.New()

	// Define a route with additional MaxBytesReader validation
	q.Post("/v1/user/maxbody/max", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Limit request body size to 1KB
		c.Request.Body = quick.MaxBytesReader(c.Response, c.Request.Body, maxBodySize)

		// Securely read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			return c.Status(http.StatusRequestEntityTooLarge).String("Request body too large")
		}
		return c.Status(http.StatusOK).Send(body)
	})

	log.Println("Server running at http://0.0.0.0:8080")
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
### 📌 cURL 

Request within limit (Valid request)
```bash
$ curl -X POST http://0.0.0.0:8080/v1/user/maxbody/max \
     -H "Content-Type: application/json" \
     --data-binary @<(head -c 800 </dev/zero | tr '\0' 'A')
```

Request exceeding limit (Should return 413)
```bash
$ curl -X POST http://0.0.0.0:8080/v1/user/maxbody/max \
     -H "Content-Type: application/json" \
     --data-binary @<(head -c 2048 </dev/zero | tr '\0' 'A')
```

---

### 📌 Key Differences Between `maxbody.New()` and `MaxBytesReader()`

| Implementation      | Description                                      |
|--------------------|--------------------------------------------------|
| `maxbody.New()`    | Enforces a **global** request body size limit **before processing** the request. |
| `MaxBytesReader()` | Adds **extra validation inside the request handler**, restricting only specific endpoints. |

---

## 📜 Logger (Request Logging)
The `logger` middleware captures HTTP request details, helping with monitoring, debugging, and analytics.

### 🚀 Key Features:
- ✅ Logs request method, path, response time, and status code.
- ✅ Supports multiple formats: text, json, and slog (structured logging).
- ✅ Helps track API usage and debugging.
- ✅ Customizable log patterns and additional fields.

### 📝 Default Logging 
This example applies simple logging.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()
	q.Use(logger.New())

	q.Use(logger.New(logger.Config{
		Level: "DEGUB",
	}))

	q.Use(logger.New(logger.Config{
		Level: "WARN",
	}))

	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "Quick ❤️",
		})
	})

	q.Listen("0.0.0.0:8080")
}
```

### 📌 cURL 

Text Logging
```bash
$ curl -i -XGET http://localhost:8080/v1/logger
```
### Console:
![Quick Logger Example](readmeLogs/log.simple.png)

---

### 📝 Structured Logging(Text Format)
This example applies logging in text format with custom log fields.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	q.Use(logger.New(logger.Config{
		Format:  "text",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG",
		CustomFields: map[string]string{
			"user_id": "usr-001",
			"trace":   "trace-debug",
		},
	}))

	q.Use(logger.New(logger.Config{
		Format:  "text",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "INFO",
		CustomFields: map[string]string{
			"user_id": "usr-002",
			"trace":   "trace-info",
		},
	}))

	q.Use(logger.New(logger.Config{
		Format:  "text",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "WARN",
		CustomFields: map[string]string{
			"user_id": "usr-003",
			"trace":   "trace-warn",
		},
	}))

	// Definir rota GET para gerar logs
	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Retornar resposta JSON
		return c.Status(200).JSON(quick.M{
			"msg": "Quick ❤️",
		})
	})

	// Iniciar o servidor na porta 8080
	q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL 

Text Logging
```bash
$ curl -i -XGET http://localhost:8080/v1/logger
```
### Console:
![Quick Logger Example](readmeLogs/log.format.text.png)

---
### 🛠️ Structured Logging (Slog Format)

This example uses structured logging (slog) for better log parsing.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)


func main() {

	q := quick.New()

	// Apply the logger middleware with structured logging (slog)
	q.Use(logger.New(logger.Config{
		Format: "slog",
		Level:  "DEBUG",
		Pattern: "[${level}] ${ip} ${method} ${path} - ${latency} " +
			"user=${user_id} trace=${trace}\n",
		CustomFields: map[string]string{
			"user_id": "99999",
			"trace":   "trace-debug",
		},
	}))

	// Apply the logger middleware with structured logging (slog)
	q.Use(logger.New(logger.Config{
		Format: "slog",
		Level:  "INFO",
		Pattern: "[${level}] ${ip} ${method} ${path} - ${latency} " +
			"user=${user_id} trace=${trace}\n",
		CustomFields: map[string]string{
			"user_id": "99999",
			"trace":   "trace-info",
		},
	}))

	// Apply the logger middleware with structured logging (slog)
	q.Use(logger.New(logger.Config{
		Format: "slog",
		Level:  "WARN",
		Pattern: "[${level}] ${ip} ${method} ${path} - ${latency} " +
			"user=${user_id} trace=${trace}\n",
		CustomFields: map[string]string{
			"user_id": "99999",
			"trace":   "trace-warn",
		},
	}))

	// Define a test route
	q.Get("/v1/logger/slog", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "Structured logging with slog",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL 

 Structured Logging (Slog)
```bash
$ curl -i -XGET http://localhost:8080/v1/logger/slog
```
### Console:
![Quick Logger Example](readmeLogs/log.format.slog.png)

---
### 📦 JSON Logging (Machine-Readable)

Ideal for log aggregation systems, this example logs in JSON format.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Apply logger with JSON format
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "INFO",
	}))

	q.Use(logger.New(logger.Config{
		Format:  "json",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG",
		CustomFields: map[string]string{
			"user_id": "usr-001",
			"trace":   "trace-debug",
		},
	}))

	// Apply the logger middleware with structured logging (slog)
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "WARN",
		Pattern: "[${level}] ${ip} ${method} ${path} - ${latency} " +
			"user=${user_id} trace=${trace}\n",
		CustomFields: map[string]string{
			"user_id": "usr-001",
			"trace":   "trace-warn",
		},
	}))

	// Define an endpoint that triggers logging
	q.Get("/v1/logger/json", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "JSON logging example",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL 

JSON Logging
```bash
$ curl -i -XGET http://localhost:8080/v1/logger/json
```
### Console:
![Quick Logger Example](readmeLogs/log.format.json.png)

---
## 🆔 MsgUUID Middleware 

### 📌 Overview

The MsgUUID Middleware in Quick is responsible for automatically generating a unique identifier (UUID) for each incoming HTTP request. This identifier is added to the response headers, allowing better tracking, debugging, and log correlation in distributed systems.

---

### 🚀 How It Works
The MsgUUID Middleware works by:

- Intercepts each incoming HTTP request before processing.
- Generating a unique UUID for each request.
- Attaching the generated UUID to the response headers for tracking.
- Helping log correlation and debugging across distributed systems.

---

### ✅ Key Features  

| Feature                    | Benefit                                                     |
|----------------------------|-------------------------------------------------------------|
| 🆔 **Unique Identifier**   | Adds a UUID to each request for tracking and correlation.  |
| 🔄 **Automatic Generation** | No need for manual UUID creation, added seamlessly.       |
| 📊 **Enhanced Debugging**   | Makes log analysis easier by attaching request identifiers. |
| 🚀 **Lightweight & Fast**   | Does not impact performance, operates efficiently.         |

---

This example generates a unique request identifier with the MsgUUUID middleware.
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msguuid"
)

func main() {
	q := quick.New()

	// Apply MsgUUID Middleware globally
	q.Use(msguuid.New())

	// Define an endpoint that responds with a UUID
	q.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log headers to validate UUID presence
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status
		return c.Status(200).JSON(nil)
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
### 📌 cURL 

```bash
$ curl -i -XGET http://localhost:8080/v1/msguuid/default
```
### 📌 Response 
```bash
"Headers":"map"[
   "Content-Type":["application/json"],
   "Msguuid":[5f49cf4d-b62e-4d81-b46e-5125b52058a6]
]
```
---
## 📩 MsgID Middleware - Quick Framework 

The `MsgID Middleware`  automatically assigns a unique identifier (MsgID) to each request. This helps with tracking, debugging, and log correlation in distributed systems.

### 🚀 Overview
- Automatically generates a unique MsgID for every incoming request.
- Ensures traceability across microservices and distributed applications.
- Adds the MsgID to both request and response headers.
- Lightweight & fast, with minimal performance overhead.

---

## ✅ Key Features

| Feature                     | Benefit                                                       |
|-----------------------------|---------------------------------------------------------------|
| 🆔 **Unique Identifier**    | Adds a MsgID to each request for tracking and correlation.   |
| 🔄 **Automatic Generation** | No need for manual MsgID creation, added seamlessly.         |
| 📊 **Enhanced Debugging**   | Makes log analysis easier by attaching request identifiers.  |
| 🚀 **Lightweight & Fast**   | Minimal performance impact, operates efficiently.            |

---
### ⚙️ How It Works
The MsgID Middleware intercepts each incoming HTTP request.
It checks if the request already has a MsgID in the headers.
If not present, it generates a new MsgID and attaches it to:
- The request headers (Msgid)
- The response headers (Msgid)

The next middleware or handler processes the request with the assigned MsgID.

Here is an example of how to use the `MsgID Middleware` with Quick:
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msguuid"
)

func main() {
	q := quick.New()

	// Apply MsgUUID Middleware globally
	q.Use(msguuid.New())

	// Define an endpoint that responds with a UUID
	q.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log headers to validate UUID presence
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status
		return c.Status(200).JSON(nil)
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
### 📌 cURL
```bash
$ curl -i -X GET http://localhost:8080/v1/msguuid/default
```

### 📌 Response
```bash
{
  "msgid": "974562398"
}
```

---
## 📚| More Examples

This directory contains **practical examples** of the **Quick Framework**, a **fast and lightweight web framework** developed in Go. 

The examples are organized in **separate folders**, each containing a complete example of using the framework in a simple web application.

### 📌 Want to contribute?
If you have some interesting example of using the **Quick Framework**, feel free to send a **Pull Request(PR)** with your contribution. 

### 👉 Check out the official example repository: 
🔗 [Quick Framework Examples](https://github.com/jeffotoni/quick/tree/main/example)



### 🤝| Contributions & Community Support

We already have several examples, and we can already test and play 😁. Of course, we are at the beginning, still has much to do.
Feel free to do **PR** (at risk of winning a Go t-shirt ❤️ and of course recognition as a professional Go 😍 in the labor market).

### 🚀 **Quick Project Supporters** 🙏

The Quick Project aims to develop and provide quality software for the developer community. 💻 To continue improving our tools, we rely on the support of our sponsors in Patreon. 🤝

We thank all our supporters! 🙌 If you also believe in our work and want to contribute to the advancement of the development community, consider supporting Project Quick on our Patreon [here](https://www.patreon.com/jeffotoni_quick)

Together we can continue to build amazing tools! 🚀



| 👤 **Avatar**                                                                                                         | 🔥 **User**                                                   | 💰 **Donation** |
|----------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------|----------------|
| <img src="https://avatars.githubusercontent.com/u/1092879?s=96&v=4" height="40">                                     | [@jeffotoni](https://github.com/jeffotoni)                     | ⭐ x 10       |
| <img src="https://avatars.githubusercontent.com/u/99341377?s=400&u=095679b08054e215561a4d4b08da764c2de619e6&v=4" height="40"> | [@Crow3442](https://github.com/Crow3442)                       | ⭐ x 5        |
| <img src="https://avatars.githubusercontent.com/u/70351793?v=4" height="40">                                         | [@Guilherme-De-Marchi](https://github.com/Guilherme-De-Marchi) | ⭐ x 5        |
| <img src="https://avatars.githubusercontent.com/u/59976892?v=4" height="40">                                         | [@jaquelineabreu](https://github.com/jaquelineabreu)           | ⭐ x 1        |

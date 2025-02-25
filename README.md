
![Logo do Quick](/quick_logo.png)

[![GoDoc](https://godoc.org/github.com/jeffotoni/quick?status.svg)](https://godoc.org/github.com/jeffotoni/quick) [![Github Release](https://img.shields.io/github/v/release/jeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/jeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jeffotoni/quick/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jeffotoni/quick/tree/main) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick) [![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/jeffotoni/quick/main) ![Coveralls](https://img.shields.io/coverallsCoverage/github/jeffotoni/quick)

# Quick - a lightweight router for go ![Quick Logo](/quick.png)
🚀 Quick is a **flexible and extensible** route manager for the Go language. Its goal is to be **fast and high-performance**, as well as being **100% compatible with net/http**. Quick is a **project in constant development** and is open for **collaboration**, everyone is welcome to contribute. 😍

💡 If you’re new to coding, Quick is a great opportunity to start learning how to work with Go. With its **ease of use** and features, you can **create custom routes** and expand your knowledge of the language.

👍 I hope you can participate and enjoy **Enjoy**! 😍

🔍 The repository of examples of the Framework Quick Run [Examples](https://github.com/jeffotoni/quick/tree/main/example).


## 🎛️| Features

| Features                                         | Has  | Status | Completion |
|--------------------------------------------------|------|--------|------------|
| 🛣️ Route Manager                                  | yes  | 🟢     | 100%       |
| 📁 Server Files Static                           | yes  | 🟢     | 100%       |
| 🚪 Route Group                                   | yes  | 🟢     | 100%       |
| 🌐 Middlewares                                   | yes  | 🟢     | 50%       |
| 🚀 HTTP/2 support                                | yes  | 🟢     | 100%       |
| 🧬 Data binding for JSON, XML and form payload   | yes  | 🟢     | 100%       |
| 🔍 Regex support                                 | yes  | 🟢     | 80%       |
| 🌐 Site                                             | yes  | 🟢     | 90%        |
| 📚 Docs                                             | yes  | 🟡     | 40%        |




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
| Develop Routes GET method accepting regular expression | 100. % |
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
| Develop support Cors   | 100. %   |
| Develop Cient Get   | 100. %   |
| Develop Cient Post support   | 100. %   |
| Develop Cient Put support   | 100. %   |
| Develop Cient support Delete   | 100. %   |


## 🚧| Rodmap in progress

 | Task   | Progress |
|--------------------------------------------------|-----------|
| Develop and relate to Listen the Config   | 42%   |
| Develops support for Group Routes - Group Get and Post | 70% |
| Develops support for Uploads and Uploads Multiples | 100% |
| Develops support for JWT | 10% |
| Develop method to Facilitate ResponseWriter handling | 80% |
| Develop method to Facilitate the handling of the Request | 80%  |
| Develop Standard of Unit Testing   | 90%   |


## 🚀| Rodmap for development


| Task   | Progress |
|---------------------------------------------------|-----------|
| Documentation Tests Examples PKG Go | 45.% |
| Test Coverage go test -cover | 74.6% |
| Regex feature coverage, but possibilities | 0.% |
| Develop for OPTIONS METHOD | 0.% |
| Develop for CONNECT METHOD [See more](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect) | 0.% |
| Develop method for ListenAndServeTLS (http2) | 0.% |
| Develop Static Files support | 0.% |
| Create a CLI (Command Line Interface) Quick. | 0.% |
| WebSocket Support | 0.% |
| Rate Limiter Support                              | 0.%       |
| Template Engines                                  | 0.%       |
| Documentation Tests Examples PKG Go   | 45. %   |
| Test coverage go test -cover   | 75.5%   |
| Coverage of Regex resources, but possibilities   | 0.%   |
| Develop for METHOD OPTIONS   | 0.%   |
| Develop for CONNECT METHOD [See more](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect)   | 0.%   |
| Develop method for ListenAndServeTLS (http2) | 0.%   |
| Develops Static Files support   | 0.%   |
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


### Quick provides a simplified API for managing uploads, allowing you to easily retrieve and manipulate files.

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

## 📚| More Examples

This directory contains practical examples of the Quick Framework, a fast and lightweight web framework developed in Go. The examples are organized in separate folders, each containing a complete example of using the framework in a simple web application. If you have some interesting example of using the Quick Framework, feel free to send a pull request with your contribution. The Quick Framework example repository can be found at [here](https://github.com/jeffotoni/quick/tree/main/example).


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
    t.Logf("\nOutputBody -> %v", string(data.Body())) // I have converted in this example to string but comes []byte as default
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












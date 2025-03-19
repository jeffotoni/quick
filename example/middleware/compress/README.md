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


Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

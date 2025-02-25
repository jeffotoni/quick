## 📂 FileSystem Middleware - Quick Framework![Quick Logo](/quick.png)
The FileSystem Middleware in Quick simplifies serving static files while providing security controls and customization options. It allows developers to serve static assets such as HTML, CSS, JavaScript, images, and other resources efficiently.

---

#### 🚀 What is FileSystem Middleware?
**`FileSystem Middleware`** helps serve static files securely and efficiently, ensuring:

- Optimized performance for loading static assets.
- Access control to prevent exposure of sensitive files.
- Easy integration with Quick’s HTTP server.

---

## ✅ Key Features

| Feature                | Benefit                                                    |
|------------------------|------------------------------------------------------------|
| 📂 **Static File Serving** | Serves HTML, CSS, JS, images, etc.                      |
| 🚀 **Efficient Handling**   | Uses `http.FileServer` for optimized performance.      |
| 🔒 **Access Restriction**   | Blocks access to sensitive files like `.git` or `.env`. |
| 🔄 **Flexible Middleware**       | Can be integrated with middlewares for authentication and logging. |

---

#### 🛠️ How It Works
- The middleware checks the Block header in each request.
- If Block is missing or set to "false", the request is rejected with a 400 Bad Request.
- If Block is "true", the request proceeds normally.
- The response includes a JSON message containing request details.

---
```go
package main

import (
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

func main() {
	app := quick.New()

	// Middleware to block requests missing "Block: true"
	app.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Block") == "" || r.Header.Get("Block") == "false" {
				w.WriteHeader(400)
				w.Write([]byte(`{"Message": "Please send Block:true in your header! :("}`))
				return
			}

			if r.Header.Get("Block") == "true" {
				w.WriteHeader(200)
				w.Write([]byte(""))
				return
			}

			h.ServeHTTP(w, r)
		})
	})

	// Protected Route
	app.Get("/v1/blocked", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type Response struct {
			Msg   string `json:"msg"`
			Block string `json:"block_message"`
		}

		return c.Status(200).JSON(&Response{
			Msg:   "Quick ❤️",
			Block: c.Headers["Block"][0],
		})
	})

	// Start the server
	log.Fatal(app.Listen("0.0.0.0:8080"))
}

```
---
#### 📌 Testing with cURL

##### 🔹 Blocked Request (Missing Block Header)
```bash
$ curl -i -X GET http://localhost:8080/v1/blocked
```

##### 🔹 Allowed Request (Block: true in Header)
```bash
$ curl -i -H "Block:true" -X GET http://localhost:8080/v1/blocked
```

---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

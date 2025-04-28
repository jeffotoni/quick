# 📏 MaxBody Middleware - Quick Framework ![Quick Logo](/quick.png)

### 📌 Overview
The **MaxBody Middleware** in Quick sets a limit on the **maximum request body size**, preventing excessively large payloads from overloading the server.
It is useful for **controlling resource usage, enhancing security, and avoiding denial-of-service (DoS) attacks.**

---
### 🚀 How It Works
The middleware intercepts incoming requests and checks the body size.
If the request body exceeds the defined limit, the request is blocked.
If the body size is within the limit, the request proceeds as usual.


### ✅ Key Features  
| Feature                  | Benefit                                                 |
|--------------------------|---------------------------------------------------------|
| 📏 **Request Size Limit**   | Restricts maximum body size to prevent large payloads. |
| 🔄 **Configurable Limit**   | Customizable size (default, defined in bytes).         |
| 🔒 **Security**             | Helps mitigate DoS attacks and excessive memory usage. |
| ⚡ **Efficient Processing** | Blocks large requests before further processing.       |

---

### 🔹Custom Limit (50KB)

This example restricts the body size to 50,000 bytes (50KB).

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/maxbody"
)

func main() {
	q := quick.New()

	// Set maximum request body size to 50KB
	q.Use(maxbody.New(50000))

	q.Post("/v1/user/maxbody/any", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		log.Printf("body: %s", c.BodyString())
		return c.Status(200).Send(c.Body())
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```

### 📌 Testing with cURL
```bash
$ curl -i -X POST http://localhost:8080/v1/user/maxbody/any \
   -H "Content-Type: application/json" \
   -d '{"data":"quick is awesome!"}'
```
--- 
### 🔹Default Limit

This example applies the default request body size limit.

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/maxbody"
)

func main() {
	q := quick.New()

	// Use default max body size
	q.Use(maxbody.New())

	q.Post("/v1/user/maxbody", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		log.Printf("body: %s", c.BodyString())
		return c.Status(200).Send(c.Body())
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
### 📌 Testing with cURL

#### 🔹Sending a valid request (Within Limit)
```bash
$ curl -i -X POST http://localhost:8080/v1/user/maxbody \
   -H "Content-Type: application/json" \
   -d '{"data":"quick is awesome!"}'
```
---

### 🔹Disabling the request body limit (Set to 0)

This example disables the body size limit by setting it to 0.

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/maxbody"
)

func main() {
	q := quick.New()

	// No limit on request body size
	q.Use(maxbody.New(0))

	q.Post("/v1/user/maxbody/large", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		log.Printf("body: %s", c.BodyString())
		return c.Status(200).Send(c.Body())
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
### 📌 Testing with cURL

#### 🔹Sending a valid request (Within Limit)
```bash
$ curl -i -X POST http://localhost:8080/v1/user/maxbody/large \
   -H "Content-Type: application/json" \
   -d '{"data":"quick is awesome!"}'
```
---
#### 📌 What I included in this README
- ✅ Overview: Explanation of what MaxBody Middleware does.
- ✅ Key Features: Table summarizing request size control, security, and efficiency.
- ✅ How It Works: Breakdown of request interception and body size validation.
- ✅ Code Examples: Different configurations for custom limits, default limit, and unlimited body size.
- ✅ Testing with cURL: Commands to send valid and oversized requests.

---

Now you can **complete with your specific examples** where I left the spaces 
.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

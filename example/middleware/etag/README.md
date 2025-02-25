## 📌 ETag Middleware in Quick ![Quick Logo](/quick.png)

#### 🔍 What is ETag?

**`ETag (Entity Tag)`** is an HTTP response header used for **caching and conditional requests**.  
It allows the client to determine if the requested resource has changed since the last request, reducing **bandwidth usage** and **improving performance**.

---

#### 🚀 How it Works
1. The server generates a unique identifier (ETag) based on the resource content.
2. The client stores the ETag and sends it in the **`If-None-Match`** header for future requests.
3. If the resource **has not changed**, the server responds with **304 Not Modified** instead of re-sending the full content.
4. If the resource **has changed**, the server sends the new content along with an updated ETag.

---

## ⚡ Features
| Feature                     | Benefit                                              |
|-----------------------------|------------------------------------------------------|
| 🚀 **Improves Performance**  | Reduces unnecessary data transfer, speeding up responses. |
| 📉 **Saves Bandwidth**       | Avoids re-downloading unchanged resources.         |
| 🔄 **Automatic Validation**  | Clients only receive updates when content changes. |
| 🔐 **Prevents Data Stale Issues** | Ensures clients always receive the latest version. |

---

## 📌 How does it work in Quick?

The **ETag Middleware** automatically generates and validates ETag headers for responses.  
It helps optimize API performance by **reducing redundant data transfers**.

✅ **Main Methods and Functionalities**
| Method | Description |
|--------|-----------|
| `c.Set("ETag", "unique-id")` | Manually sets an ETag for the response. |
| `If-None-Match` | Clients send this header to check if content has changed. |
| `304 Not Modified` | Returned when the resource hasn't changed, saving bandwidth. |

---

## 📌 ETag Middleware Example

```go
package main

import (
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

// curl -i -H "Block:true" -XGET localhost:8080/v1/blocked
func main() {

	app := quick.New()

	app.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//This middleware will block your request if it does not pass header Block:true
			if r.Header.Get("Block") == "" || r.Header.Get("Block") == "false" {
				w.WriteHeader(400)
				w.Write([]byte(`{"Message": "Envia block em seu header, por favor! :("}`))
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

	app.Get("/v1/blocked", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type my struct {
			Msg   string `json:"msg"`
			Block string `json:"block_message"`
		}

		log.Println(c.Headers["Messageid"])

		return c.Status(200).JSON(&my{
			Msg:   "Quick ❤️",
			Block: c.Headers["Block"][0],
		})
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

```
---
#### 📌 Testing with cURL

##### 🔹No Header Block (Answer 400 - Bad Request)
```bash
$ curl -i -X GET http://localhost:8080/v1/blocked
```
##### 🔹 With Block: true (Answer 200 - OK)
```bash
$ curl -i -H "Block: true" -X GET http://localhost:8080/v1/blocked
```
---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

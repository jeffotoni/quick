## 🚀 Request Limiter Middleware - Quick Framework![Quick Logo](/quick.png)

#### 📌 Overview
The Limiter Middleware in Quick helps control the number of requests a client can make within a certain timeframe. This is useful for rate limiting, preventing abuse, and enhancing API security.

#### ✅ Key Features

| Feature          | Benefit                                              |
|-----------------|------------------------------------------------------|
| ⏳ **Rate Limiting**  | Prevents excessive requests from a single client.  |
| 🔄 **Configurable**   | Allows custom request limits per user/IP.          |
| 🔒 **Security**       | Helps mitigate **DoS (Denial-of-Service) attacks**. |
| ⚡ **Efficient**      | Uses lightweight in-memory tracking for performance. |
---

#### 🛠️ How It Works
- The middleware tracks client requests using headers or IP addresses.
- If a client exceeds the allowed request limit, a **`429 Too Many Requests`** response is returned.
- Otherwise, the request proceeds normally to the next handler.

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

	app.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Este middleware, irá bloquear sua requisicao se não passar header Block:true
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

##### 🔹 Successful Request (Within Limit)
```bash
$ curl -i -H "Block:true" -XGET localhost:8080/v1/blocked

```

---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

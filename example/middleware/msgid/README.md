# 📩 MsgID Middleware - Quick Framework ![Quick Logo](/quick.png)

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

	// Aplica o Middleware MsgID globalmente
	q.Use(msgid.New())

	// Define uma rota que retorna o MsgID gerado
	q.Get("/v1/msgid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Obtém o MsgID do header da requisição
		msgId := c.Request.Header.Get("Msgid")

		// Log para depuração
		fmt.Printf("Generated MsgID: %s\n", msgId)

		// Retorna o MsgID no JSON da resposta
		return c.Status(200).JSON(map[string]string{"msgid": msgId})
	})

	// Inicia o servidor
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
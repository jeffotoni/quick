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

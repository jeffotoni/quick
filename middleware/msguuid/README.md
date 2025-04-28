
## 🆔 MsgUUID
Assigns a UUID (Universally Unique Identifier) to each request.

- Allows easy tracking of requests in logs.
- Useful for distributed systems where tracing requests across services is required.
- Adds a unique identifier to every request automatically.


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

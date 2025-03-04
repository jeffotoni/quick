# 🆔 MsgUUID Middleware - Quick Framework ![Quick Logo](/quick.png)

## 📌 Overview
The MsgUUID Middleware in Quick automatically generates a unique identifier (UUID) for each request and adds it to the response headers.
This helps track requests, improve debugging, and correlate logs.

---

### 🚀 How It Works
The MsgUUID Middleware works by:

- Intercepting all incoming HTTP requests.
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

### 🔹Attaching a UUID to responses
This example ensures that all responses from the /v1/msguuid/default route contain a unique **UUID.**


```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msguuid"
)

func main() {
	app := quick.New()

	// Enable MsgUUID Middleware
	app.Use(msguuid.New())

	// Define an endpoint that includes MsgUUID
	app.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log the response headers to check the UUID
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status with no body
		return c.Status(200).JSON(nil)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

```

### 📌 Testing with cURL

#### 🔹Sending a GET request
```bash
$ curl -i -X GET http://localhost:8080/v1/user
```
---

### 🔹Default MsgUUID Usage
The following example demonstrates how MsgUUID automatically attaches a **UUID** to every response.
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msguuid"
)

func main() {
	app := quick.New()

	// Apply MsgUUID Middleware globally
	app.Use(msguuid.New())

	// Define an endpoint that responds with a UUID
	app.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log headers to validate UUID presence
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status
		return c.Status(200).JSON(nil)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

```

### 📌 Testing with cURL

#### 🔹Sending a GET request to the default MsgUUID route

```bash
$ curl -i -X GET http://localhost:8080/v1/msguuid/default
```
---

### 📌 What I Included in This README
- ✅ Overview: Explanation of MsgUUID and its purpose.
- ✅ Key Features: Why use MsgUUID middleware?
- ✅ How It Works: Breakdown of request interception and UUID assignment.
- ✅ Code Examples: How to apply MsgUUID to endpoints.
- ✅ Testing with cURL: Examples of expected response headers.
- ✅ Advantages Table: Summary of benefits.

---

Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

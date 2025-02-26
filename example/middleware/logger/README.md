## 📜 Logger Middleware - Quick Framework![Quick Logo](/quick.png)

The **`Logger Middleware`** provides automatic logging of incoming HTTP requests in the Quick Framework. It helps developers track API requests, response times, and request details in a structured way.

---

#### ✅ Key Features

| Feature                  | Benefit                                                             |
| ------------------------ | ------------------------------------------------------------------- |
| 📄 **Request Logging**   | Automatically logs incoming requests with method, path, and status. |
| ⏳ **Execution Time**    | Captures the duration of each request.                              |
| 📊 **Debugging**         | Helps identify slow or failing requests.                            |
| 📜 **Structured Output** | Logs information in an easy-to-read format.                         |

---

#### 🚀 How It Works

When enabled, the Logger middleware captures and prints each request's details, including:

- **HTTP Method** (GET, POST, PUT, DELETE)
- **Request Path**
- **Status Code**
- **Response Time**

It helps in debugging and analyzing API performance.

---

#### 📝 Code Example:

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Enable logger middleware
	q.Use(logger.New())

	// Example route
	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type my struct {
			Msg string `json:"msg"`
		}

		return c.Status(200).JSON(&my{
			Msg: "Quick ❤️",
		})
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
```

#### 📌 Testing with cURL

##### 🔹 Sending a request

```bash
$ curl -i -X GET http://localhost:8080/v1/logger
```

#### This log message shows:
![Quick Logger Example](logger_output.png)

---
#### 📌 What I included in this README
- ✅ Overview: Explanation of what Logger Middleware does and its benefits.
- ✅ Key Features: Table highlighting request logging, execution time, debugging, and structured output.
- ✅ How It Works: Breakdown of logged request details (method, path, status, response time).
- ✅ Code Example: Implementation of Logger Middleware with a test route (/v1/logger).
- ✅ Testing with cURL: Example cURL request to verify logs in the terminal.
- ✅ Expected Output.

---



Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

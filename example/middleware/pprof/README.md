## 📌 PPROF Middleware in Quick ![Quick Logo](/quick.png)

**pprof** provides profiling endpoints for your Quick application. It helps you to identify 
performance bottlenecks, monitor resource usage, and ensure that the code runs efficiently

---
### 🔻 Environment
Profiling is only enabled in development mode (APP_ENV=development). 

We strongly recommend to use it only in development mode because in production it can introduce
unwanted overhead and potentially degrade performance.

---
### 🧩 Example Usage
```go
package main

import (
	"errors"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	q := quick.New()

	// Apply the Profiling middleware
	q.Use(pprof.New())

	// Define a test route
	q.Get("/", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(quick.StatusOK).String("OK")
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
```

---
### Routes

Profiling middleware registers a set of routes for profiling:

- `/debug/pprof`
- `/debug/cmdline`
- `/debug/profile`
- `/debug/symbol`
- `/debug/pprof/trace`
- `/debug/goroutine`
- `/debug/heap`
- `/debug/threadcreate`
- `/debug/mutex`
- `/debug/allocs`
- `/debug/block`

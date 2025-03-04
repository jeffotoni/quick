# 📌 GET - Quick Framework ![Quick Logo](/quick.png)

The method 'GET' can be used to fetch values of different types, such as strings, integers and JSON responses. 


### 📜 Code Implementation
In this example, we show how to set up different routes using the Quick Framework.

```go
package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New() // Initialize Quick framework

	// Route to greet a user by name (dynamic route parameter)
	q.Get("/v1/user/:name", func(c *quick.Ctx) error {
		name := c.Param("name")                              // Retrieve the 'name' parameter from the URL
		c.Set("Content-Type", "text/plain")                  // Set response content type as plain text
		return c.Status(200).SendString("Olá " + name + "!") // Return greeting message
	})

	// Simple route returning a static message
	q.Get("/v2/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")            // Set response content type as JSON
		return c.Status(200).SendString("Opa, funcionando!") // Return confirmation message
	})

	// Route to return an ID from the URL
	q.Get("/v3/user/:id", func(c *quick.Ctx) error {
		id := c.Param("id")                         // Retrieve the 'id' parameter from the URL
		c.Set("Content-Type", "application/json")   // Set response content type as JSON
		return c.Status(200).SendString("Id:" + id) // Return the ID in the response
	})

	// Complex route with multiple parameters
	q.Get("/v1/userx/:p1/:p2/cust/:p3/:p4", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")              // Set response content type as JSON
		return c.Status(200).SendString("Quick in action ❤️!") // Return a success message
	})

	// Print all registered routes
	for k, v := range q.GetRoute() {
		fmt.Println(k, "[", v, "]")
	}

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}

```

### 📌 Testing with cURL

### 🔹Returns a greeting message with the given name
```bash
$ curl --location --request GET "http://localhost:8080/v1/user/Jeff"
```

### 🔹Basic GET request
```bash
$ curl --location --request GET "http://localhost:8080/v2/user"
```

### 🔹Get user by ID
```bash
$ curl --location --request GET "http://localhost:8080/v1/user/123"
```

### 🔹Complex route with multiple parameters
```bash
$ curl --location --request GET "http://localhost:8080/v1/userx/test1/test2/cust/test3/test4"
```
---
#### 📌 What I included in this
- ✅ GET method in Quick Framework
- ✅ Go implementation with dynamic parameters (:name, :id)
- ✅ GET routes for user retrieval and static responses
- ✅ Handling of dynamic and complex route patterns
- ✅ cURL examples for all GET endpoints

---

Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
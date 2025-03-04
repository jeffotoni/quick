# 📌 Using Regex in Routes - Quick Framework ![Quick Logo](/quick.png)

`Regex` (or "Regular Expressions") is a powerful technique used in programming to match and manipulate text patterns.

The Quick Framework supports dynamic routes but does not support inline regex definitions like `{id:[0-9]+}` in route parameters. Instead, Quick uses colon-prefixed `(:param)` parameters for dynamic route matching.

### 📜 Code Implementation

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msgid"
)

func main() {
	q := quick.New()

	// Adding middleware msgid
	q.Use(msgid.New())

	// Corrected route using :id instead of {id:[0-9]+}
	q.Get("/v1/user/:id", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).String("Quick ação total!!!")
	})

	q.Listen(":8080")
}

```
#### 📌 Testing with cURL

##### 🔹 Retrieve user details by ID:

```go
$ curl --location --request GET "http://localhost:8080/v1/user/123" \
--header "Content-Type: application/json"
```
---

#### 📌 What I included in this README
- ✅ Introduction to Regex and its usage in Quick Framework.
- ✅ Explanation of Quick’s dynamic route handling.
- ✅ Example of defining a dynamic route with :param.
- ✅ cURL example to test dynamic route behavior.

---

Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
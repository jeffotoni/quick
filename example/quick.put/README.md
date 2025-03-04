# 📌 PUT - Quick Framework ![Quick Logo](/quick.png)

The `Put()` **method** in the **Quick Framework** is used to define an HTTP route that handles **PUT requests**. The PUT method allows clients to send data to a server for **updating existing resources.**

In simple terms, it lets a client **modify** an existing resource on a server.

The example defines a PUT route using Quick. To create a PUT route, simply call the Put() method in a Quick instance.

1. The route path (`/users/:id` or `/types/:id`)
2. A handler function that executes when the route is matched.

---
### 📜 Code Implementation

```go
package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New() // Initialize Quick framework

	// PUT route to update a user by ID
	q.Put("/users/:id", func(c *quick.Ctx) error {
		userID := c.Param("id") // Retrieve the user ID from the URL parameter
		// Logic to update user data would go here
		return c.Status(200).SendString("User " + userID + " updated successfully!")
	})

	// PUT route to update a specific type by ID
	q.Put("/types/:id", func(c *quick.Ctx) error {
		typeID := c.Param("id") // Retrieve the type ID from the URL parameter
		// Logic to update the type would go here
		return c.Status(200).SendString("Types " + typeID + " type updated successfully!")
	})

	// Start the server and listen on port 8080
	q.Listen(":8080")
}

```

#### 📌 Testing with cURL

##### 🔹 Updates the details of a user based on the provided id:

```bash
$ curl --location --request PUT "http://localhost:8080/users/123" \
--header "Content-Type: application/json" \
--data '{"name":"UpdatedUser","year":2024}'
```

##### 🔹 Update a Type:

```bash
$ curl --location --request PUT "http://localhost:8080/types/456" \
--header "Content-Type: application/json" \
--data '{"type":"admin"}'
```


---

#### 📌 What I included in this README
- ✅ Overview of the PUT method in Quick Framework
- ✅ Go implementation of dynamic parameter handling (:id)
- ✅ PUT routes for updating users and types
- ✅ Handling of success (200 OK) responses
- ✅ cURL examples for different PUT endpoints


Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥
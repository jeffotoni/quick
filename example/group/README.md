# 📂 Group - Grouping Routes in Quick ![Quick Logo](/quick.png)


**Group** is a feature of the **Quick Framework** that allows you to group routes and apply middleware to them.

For example, if you have a set of routes that need authentication before they are accessed, instead of adding the authentication middleware individually for each route, can group them using the Group feature and apply middleware to all routes in the group at once. This can make the code more readable and organized, as well as avoiding code repetition.

---

## 📌 Why use Group?

| 🔹 **Advantage**   | ✅ **Benefit**   |
|--------------------|---------------------------------------------------------|
| 📂 **Organization** | Structure routes into logical groups.   |
| 🔄 **Reuse** | Avoids code repetition when applying middlewares.   |
| 🔒 **Security**   | Allows you to protect routes with authentication or validation.  |
| ⚡ **Performance** | Middleware processed once for all routes in the group. |

---

#### 📝 Example of Use


#### Group 1 - Creating a Route Group (/v1/user)

```go
package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024,
	})

	// Creating a route group for /v1
	group := q.Group("/v1")

	// Definition of routes within the group
	group.Get("/user", func(c *quick.Ctx) error {
		return c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
	})

	group.Post("/user", func(c *quick.Ctx) error {
		return c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
	})

 q.Listen("0.0.0.0:8080")
}
```
### 📌 Testing with cURL
```bash
$ curl --location --request GET 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json/' \
--data '[GET] [GROUP] /v1/user ok!!!'
```

--- 

#### Group 2 - Creating a Second Group (/v2/user)

```go
package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024, //Defines the limit of the request body
	})

	// Creating a second group of routes for /v2
	group2 := q.Group("/v2")

	group2.Get("/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick in action com [GET] /v2/user ❤️!")
	})

	group2.Post("/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick in action com [POST] /v2/user ❤️!")
	})


 q.Listen("0.0.0.0:8080")
}
```
### 📌 Testing with cURL

```bash
$ curl --location --request GET 'http://localhost:8080/v2/user' \
--header 'Content-Type: application/json/' \
--data 'Quick in action com [POST] /v2/user ❤️!'
```
---

#### **📝 What I included in this README**

- ✅ What is **Group** in Quick  
- ✅ How to group routes and apply **middlewares**  
- ✅ **Implementation example** with `/v1', and `/v2'  
- ✅ **Testing with cURL** to check requests  


🚀 Now you know how to use **Group** in Quick to structure your routes!  


Now you can **complete with your specific examples** where I left the spaces ` ```go ... ``` `.

🚀 **If you need adjustments or improvements, just let me know!** 😃🔥

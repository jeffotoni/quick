package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New()
	q.Get("/v1/user",
		func(c *quick.Ctx) error {
			return c.Status(200).
				String("Quick! ❤️")
		})
	_ = q.Listen(":8080")
}

// ----------

/**
$ curl -i XGET localhost:8080/v1/user
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 27 Mar 2025 19:06:20 GMT
Content-Length: 13

Quick! ❤️
*/



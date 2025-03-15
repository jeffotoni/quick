package main

import "github.com/jeffotoni/quick"

func main() {

	q := quick.New()

	q.Get("/v1/user", func(c *quick.Ctx) error {
		userAgent := c.GetHeader("User-Agent")
		ip := c.RemoteIP()
		method := c.Method()
		path := c.Path()
		queryValue := c.QueryParam("search")

		return c.Status(200).JSON(map[string]string{
			"user_agent": userAgent,
			"ip":         ip,
			"method":     method,
			"path":       path,
			"search":     queryValue,
		})
	})

	q.Listen(":8080")
}

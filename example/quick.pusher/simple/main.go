package main

import (
	"time"

	"github.com/jeffotoni/quick"
)

// you need to create key.pem and cert.pem
// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
// https://localhost:8443
func main() {
	q := quick.New()

	q.Get("/*", func(c *quick.Ctx) error {
		pusher, ok := c.Pusher()
		if ok {
			pusher.Push("/public/style.css", nil)
			pusher.Push("/public/app.js", nil)
		}

		time.Sleep(2 * time.Second)

		html := `<!DOCTYPE html>
		<html>
		<head><link rel="stylesheet" href="/public/style.css"></head>
		<body>
			<h1>HTTP/2 Push Test</h1>
			<script src="/public/app.js"></script>
		</body>
		</html>`

		c.Set("Content-Type", "text/html; charset=utf-8")

		return c.Status(200).SendString(html)
	})

	q.Static("/public", "./public")

	// true = enable HTTP/2
	q.ListenTLS(":8443", "cert.pem", "key.pem", true)
}

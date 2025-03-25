package helmet

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

func ExampleHelmet() {

	q := quick.New()

	// Use Helmet middleware with default security headers
	q.Use(Helmet())

	// Simple route to test headers
	q.Get("/v1/user", func(c *quick.Ctx) error {

		// list all headers
		headers := make(map[string]string)
		for k, v := range c.Response.Header() {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		return c.Status(200).JSONIN(headers)
	})

	// Send test request using Quick's built-in test utility
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodGet,
		URI:     "/v1/user",
		Headers: map[string]string{"Content-Type": "application/json"},
	})

	fmt.Println("Response Body:", string(resp.Body()))

	// Output:
	// 	Response Body: {
	//   "Cache-Control": "no-cache, no-store, must-revalidate",
	//   "Content-Security-Policy": "default-src 'self'",
	//   "Cross-Origin-Embedder-Policy": "require-corp",
	//   "Cross-Origin-Opener-Policy": "same-origin",
	//   "Cross-Origin-Resource-Policy": "same-origin",
	//   "Origin-Agent-Cluster": "?1",
	//   "Referrer-Policy": "no-referrer",
	//   "X-Content-Type-Options": "nosniff",
	//   "X-Dns-Prefetch-Control": "off",
	//   "X-Download-Options": "noopen",
	//   "X-Frame-Options": "SAMEORIGIN",
	//   "X-Permitted-Cross-Domain-Policies": "none",
	//   "X-Xss-Protection": "0"
	// }

}

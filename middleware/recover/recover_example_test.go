package recover

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// ExampleNew_defaultBehavior demonstrates the default behavior of the Recover middleware.
func ExampleNew_defaultBehavior() {
	q := quick.New()

	// Use the default Recover middleware
	q.Use(New())

	// Define a test route
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		panic("Panicking!")
	})

	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 500
	// Body: Internal Server Error
}

// ExampleNew_withNextSkipping demonstrates how to use the Next() function to skip the middleware.
func ExampleNew_withNextSkipping() {
	q := quick.New()

	// Use the Recover middleware with Next() function
	q.Use(New(Config{
		Next: func(c *quick.Ctx) bool {
			return true // Always skip
		},
	}))

	// Define a test route with panic
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		panic("Panicking!")
	})

	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Out put:
	// Status: 200
	// Body:
}

// ExampleNew_withStacktraceDisabled demonstrates how to disable the stacktrace.
func ExampleNew_withStacktraceDisabled() {
	q := quick.New()

	// Use the Recover middleware with stacktrace disabled
	q.Use(New(Config{
		EnableStacktrace: false,
	}))

	// Define a test route with panic
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		panic("Panicking!")
	})

	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 500
	// Body: Internal Server Error
}

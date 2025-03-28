package recover

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// ExampleNew_defaultBehavior demonstrates the default behavior of the Recover middleware.
// This function is named ExampleNew_defaultBehavior()
// it with the Examples type.
func ExampleNew_defaultBehavior() {
	q := quick.New()

	// Uses the default recover middleware
	q.Use(New())

	// Route that generates a panic
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		// panic("Panicking!") // This would cause a status 500, handled by recover middleware.
		return c.String("Panicking!")
	})
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	if err := resp.AssertString("Panicking!"); err != nil {
		fmt.Println("body error:", err)
	}

	fmt.Println(resp.BodyStr())

	// Output:
	// Panicking!
}

// ExampleNew_withNextSkipping demonstrates how to use the Next() function to skip the middleware.
// This function is named ExampleNew_withNextSkipping()
// it with the Examples type.
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

	if err := resp.AssertStatus(500); err != nil {
		fmt.Println("status error:", err)
	}

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 500
	// Body: Internal Server Error
}

// ExampleNew_withStacktraceDisabled demonstrates how to disable the stacktrace.
// This function is named ExampleNew_withStacktraceDisabled()
// it with the Examples type.
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

	if err := resp.AssertStatus(500); err != nil {
		fmt.Println("status error:", err)
	}

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 500
	// Body: Internal Server Error
}

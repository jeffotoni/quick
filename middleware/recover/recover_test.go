package recover

import (
	"testing"

	"github.com/jeffotoni/quick"
)

// TestWithStacktraceDisabled tests when stacktrace is disabled.
func TestWithStacktraceDisabled(t *testing.T) {
	q := quick.New()
	q.Use(New(Config{
		App:              q,
		EnableStacktrace: false,
	}))

	// Define a test route
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		panic("Panicking!")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := resp.AssertStatus(quick.StatusInternalServerError); err != nil {
		t.Error(err)
	}
}

// TestWithStacktraceEnabled tests when stacktrace is enabled.
func TestWithStacktraceEnabled(t *testing.T) {
	q := quick.New()
	q.Use(New(Config{
		App:              q,
		EnableStacktrace: true,
	}))

	// Define a test route
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		panic("Panicking!")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	if err != nil {
		t.Fatal(err)
	}
	if err := resp.AssertStatus(quick.StatusInternalServerError); err != nil {
		t.Error(err)
	}
}

// TestWithNextSkipping tests when Next() returns true (skips route logic).
func TestWithNextSkipping(t *testing.T) {
	q := quick.New()

	// Use the Recover middleware with Next() function
	q.Use(New(Config{
		App: q,
		Next: func(c *quick.Ctx) bool {
			return true // Always skip
		},
	}))

	// Define a test route with panic
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		panic("Panicking!")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/recover",
	})

	if err != nil {
		t.Fatal(err)
	}
	if err := resp.AssertStatus(quick.StatusOK); err != nil {
		t.Error(err)
	}
}

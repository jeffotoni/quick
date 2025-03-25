package timeout

import (
	"context"
	"errors"
	"time"

	"github.com/jeffotoni/quick"
)

type Options struct {
	// Duration is the timeout duration
	Duration time.Duration

	// Next is an optional function. If it returns true, the middleware is skipped.
	// This is useful to conditionally disable timeout for specific requests.
	Next func(c *quick.Ctx) bool
}

func defaultOptions(opt ...Options) Options {
	if len(opt) == 0 {
		return Options{
			Duration: 5 * time.Second,
			Next:     func(c *quick.Ctx) bool { return false },
		}
	}

	return opt[0]
}

// New creates a new timeout middleware. It returns a middleware function that
// can be used with quick. Use it like:
//
//	app := quick.New()
//	app.Use(timeout.New(timeout.Options{
//		Duration: 5 * time.Second,
//	}))
func New(opt ...Options) func(next quick.Handler) quick.Handler {
	option := defaultOptions(opt...)

	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// Skip middleware logic if Next returns true
			if option.Next != nil && option.Next(c) {
				return next.ServeQuick(c)
			}

			// skip middleware logic if duration is 0
			if option.Duration <= 0 {
				return next.ServeQuick(c)
			}
			// Create a timeout context
			ctx, cancel := context.WithTimeout(c.Request.Context(), option.Duration)
			defer cancel()

			// Attach timeout context to request
			c.Request = c.Request.WithContext(ctx)

			// Execute the next handler
			return runHandler(c, next, context.DeadlineExceeded)

		})
	}
}

// runHandler executes the handler and checks for timeout
// If the context deadline is exceeded, it returns a timeout response.
func runHandler(c *quick.Ctx, h quick.Handler, errs ...error) error {
	err := h.ServeQuick(c)
	if errors.Is(err, context.DeadlineExceeded) {
		return c.Status(quick.StatusRequestTimeout).SendString("Request Timeout")
	}
	return err
}

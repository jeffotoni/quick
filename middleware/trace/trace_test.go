package trace

import (
	"context"
	"testing"
	"time"

	"github.com/jeffotoni/quick"
)

func TestNewMiddleware_InjectionAndRetrieval(t *testing.T) {
	calledTrace := false
	calledUser := false

	mw := New(Config{
		Timeout: 5 * time.Second,
		Fields: map[string]func(*quick.Ctx) string{
			"X-Trace-ID": func(c *quick.Ctx) string {
				calledTrace = true
				return "trace-123"
			},
			"X-User-ID": func(c *quick.Ctx) string {
				calledUser = true
				return "user-999"
			},
		},
	})

	// Mock next handler
	nextCalled := false
	next := quick.HandlerFunc(func(c *quick.Ctx) error {
		nextCalled = true
		ctx := c.Ctx()

		if GetCtx(ctx, "X-Trace-ID") != "trace-123" {
			t.Errorf("Expected X-Trace-ID to be 'trace-123', got '%s'", GetCtx(ctx, "X-Trace-ID"))
		}

		if GetCtx(ctx, "X-User-ID") != "user-999" {
			t.Errorf("Expected X-User-ID to be 'user-999', got '%s'", GetCtx(ctx, "X-User-ID"))
		}

		// all := GetMap(ctx)
		// if len(all) != 2 {
		// 	t.Errorf("Expected 2 fields in context map, got %d", len(all))
		// }

		return nil
	})

	ctx := context.Background()
	c := &quick.Ctx{}
	c.SetCtx(ctx)

	h := mw(next)
	err := h.ServeQuick(c)

	if err != nil {
		t.Errorf("Middleware returned unexpected error: %v", err)
	}

	if !nextCalled {
		t.Errorf("Next handler was not called")
	}
	if !calledTrace || !calledUser {
		t.Errorf("Expected both field functions to be called")
	}
}

func TestGet_DefaultEmptyValues(t *testing.T) {
	ctx := context.Background()

	if val := GetCtx(ctx, "X-Trace-ID"); val != "" {
		t.Errorf("Expected empty value for missing key, got '%s'", val)
	}

	// m := GetMap(ctx)
	// if len(m) != 0 {
	// 	t.Errorf("Expected empty context map, got %v", m)
	// }
}

func TestGet_NilContext(t *testing.T) {
	var ctx context.Context

	if val := GetCtx(ctx, "X-Trace-ID"); val != "" {
		t.Errorf("Expected empty string on nil context, got '%s'", val)
	}

	// m := GetMap(ctx)
	// if len(m) != 0 {
	// 	t.Errorf("Expected empty map on nil context, got %v", m)
	// }
}

// func TestGetMap_RetrievesAllKeys(t *testing.T) {
// 	ctx, cancel := NewCtx().
// 		Set("X-Trace-ID", "trace-abc").
// 		Set("X-User-ID", "user-xyz").
// 		Build()
// 	defer cancel()

// 	vals := GetMap(ctx)

// 	if len(vals) != 2 {
// 		t.Errorf("Expected 2 values in map, got %d", len(vals))
// 	}

// 	if !strings.HasPrefix(vals["X-Trace-ID"], "trace") {
// 		t.Errorf("Expected 'X-Trace-ID' to start with 'trace', got '%s'", vals["X-Trace-ID"])
// 	}

// 	if !strings.HasPrefix(vals["X-User-ID"], "user") {
// 		t.Errorf("Expected 'X-User-ID' to start with 'user', got '%s'", vals["X-User-ID"])
// 	}
// }

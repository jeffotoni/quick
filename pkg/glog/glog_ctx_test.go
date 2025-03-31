package glog

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewCtx_DefaultKey(t *testing.T) {
	ctx, cancel := NewCtx().Set("TraceID", "abc123").Build()
	defer cancel()

	val := GetCtx(ctx)
	if val != "abc123" {
		t.Errorf("Expected default key value 'abc123', got '%s'", val)
	}
}

func TestNewCtx_CustomKeyName(t *testing.T) {
	ctx, cancel := NewCtx().Set("X-User-ID", "user42").Build()
	defer cancel()

	val := GetCtx(ctx, "X-User-ID")
	if val != "user42" {
		t.Errorf("Expected custom key value 'user42', got '%s'", val)
	}
}

func TestNewCtx_MultipleKeys(t *testing.T) {
	ctx, cancel := NewCtx().
		Set("X-Trace-ID", "abc123").
		Set("User", "jeff").
		Set("env", "prod").
		Build()
	defer cancel()

	got := GetCtxMap(ctx)
	if len(got) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(got))
	}
	if got["X-Trace-ID"] != "abc123" {
		t.Errorf("Expected 'X-Trace-ID' = 'abc123', got '%s'", got["X-Trace-ID"])
	}
	if got["User"] != "jeff" {
		t.Errorf("Expected 'User' = 'jeff', got '%s'", got["User"])
	}
	if got["env"] != "prod" {
		t.Errorf("Expected 'env' = 'prod', got '%s'", got["env"])
	}
}

func TestNewCtx_NilContext(t *testing.T) {
	var ctx context.Context
	val := GetCtx(ctx)
	if val != "" {
		t.Errorf("Expected empty value for nil context, got '%s'", val)
	}
}

func TestNewCtx_KeyNotFound(t *testing.T) {
	ctx := context.Background()
	val := GetCtx(ctx)
	if val != "" {
		t.Errorf("Expected empty value for missing key, got '%s'", val)
	}
}

func TestNewCtx_Timeout(t *testing.T) {
	ctx, cancel := NewCtx().
		Set("TraceID", "with-timeout").
		Timeout(10 * time.Millisecond).
		Build()
	defer cancel()

	select {
	case <-ctx.Done():
		if !strings.Contains(ctx.Err().Error(), "context deadline") {
			t.Errorf("Expected deadline exceeded, got: %v", ctx.Err())
		}
	case <-time.After(20 * time.Millisecond):
		// OK
	}
}

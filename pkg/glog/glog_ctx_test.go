package glog_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jeffotoni/quick/pkg/glog"
)

// TestNewCtx_DefaultKey verifies that a context built with a "TraceID" key
// correctly stores and retrieves its value.
//
// To run:
//
//	go test -v -run ^TestNewCtx_DefaultKey$
func TestNewCtx_DefaultKey(t *testing.T) {
	ctx, cancel := glog.CreateCtx().Set("TraceID", "abc123").Build()
	defer cancel()

	val := glog.GetCtx(ctx, "TraceID")
	if val != "abc123" {
		t.Errorf("Expected default key value 'abc123', got '%s'", val)
	}
}

// TestNewCtx_CustomKeyName verifies that custom key names (e.g., "X-User-ID")
// can be used and correctly retrieved from the context.
//
// To run:
//
//	go test -v -run ^TestNewCtx_CustomKeyName$
func TestNewCtx_CustomKeyName(t *testing.T) {
	ctx, cancel := glog.CreateCtx().Set("X-User-ID", "user42").Build()
	defer cancel()

	val := glog.GetCtx(ctx, "X-User-ID")
	if val != "user42" {
		t.Errorf("Expected custom key value 'user42', got '%s'", val)
	}
}

// TestNewCtx_NilContext ensures that passing a nil context to GetCtx
// returns an empty string instead of causing a panic or error.
//
// To run:
//
//	go test -v -run ^TestNewCtx_NilContext$
func TestNewCtx_NilContext(t *testing.T) {
	var ctx context.Context
	val := glog.GetCtx(ctx)
	if val != "" {
		t.Errorf("Expected empty value for nil context, got '%s'", val)
	}
}

// TestNewCtx_KeyNotFound verifies that requesting a non-existent key
// from a context returns an empty string.
//
// To run:
//
//	go test -v -run ^TestNewCtx_KeyNotFound$
func TestNewCtx_KeyNotFound(t *testing.T) {
	ctx := context.Background()
	val := glog.GetCtx(ctx)
	if val != "" {
		t.Errorf("Expected empty value for missing key, got '%s'", val)
	}
}

// TestNewCtx_Timeout ensures that contexts built with a short timeout
// correctly expire after the given duration.
//
// To run:
//
//	go test -v -run ^TestNewCtx_Timeout$
func TestNewCtx_Timeout(t *testing.T) {
	ctx, cancel := glog.CreateCtx().Set("TraceID", "with-timeout").Timeout(10 * time.Millisecond).Build()
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

// TestNewCtx_GetCtxAll validates that GetCtxAll returns all
// keys and values previously set in the context.
//
// To run:
//
//	go test -v -run ^TestNewCtx_GetCtxAll$
func TestNewCtx_GetCtxAll(t *testing.T) {
	ctx, cancel := glog.CreateCtx().
		Set("TraceID", "abc-123").
		Set("X-User-ID", "user-456").
		Set("X-Session-ID", "sess-789").
		Build()
	defer cancel()

	values := glog.GetCtxAll(ctx)

	if len(values) != 3 {
		t.Errorf("Expected 3 context values, got %d", len(values))
	}

	if values["TraceID"] != "abc-123" {
		t.Errorf("Expected TraceID to be 'abc-123', got '%s'", values["TraceID"])
	}
	if values["X-User-ID"] != "user-456" {
		t.Errorf("Expected X-User-ID to be 'user-456', got '%s'", values["X-User-ID"])
	}
	if values["X-Session-ID"] != "sess-789" {
		t.Errorf("Expected X-Session-ID to be 'sess-789', got '%s'", values["X-Session-ID"])
	}
}

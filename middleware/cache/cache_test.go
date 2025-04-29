package cache

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jeffotoni/quick"
)

func TestCacheMiddleware(t *testing.T) {
	t.Run("Basic Cache Hit and Miss", func(t *testing.T) {
		q := quick.New()

		// Add cache middleware with default config
		q.Use(New())

		// Counter to track handler executions
		var counter int

		// Add a test handler that increments the counter
		q.Get("/test", func(c *quick.Ctx) error {

			counter++
			return c.String(fmt.Sprintf("Response #%d", counter))
		})

		// First request should miss cache
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		if counter != 1 {
			t.Errorf("Expected counter to be 1, got %d", counter)
		}

		// Second request should hit cache
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			t.Errorf("Expected X-Cache-Status to be HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		if counter != 1 {
			t.Errorf("Expected counter to still be 1, got %d", counter)
		}
	})

	t.Run("Custom Key Generator", func(t *testing.T) {
		q := quick.New()

		// Add cache middleware with custom key generator
		q.Use(New(Config{
			KeyGenerator: func(c *quick.Ctx) string {
				return c.Path() + "?lang=" + c.Query["lang"]
			},
		}))

		// Counter to track handler executions
		var counter int

		// Add a test handler that increments the counter
		q.Get("/test", func(c *quick.Ctx) error {

			counter++
			return c.String(fmt.Sprintf("Response #%d", counter))
		})

		// First request with lang=en
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test?lang=en",
		})
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		if counter != 1 {
			t.Errorf("Expected counter to be 1, got %d", counter)
		}

		// Second request with lang=fr (should miss cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test?lang=fr",
		})
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		if counter != 2 {
			t.Errorf("Expected counter to be 2, got %d", counter)
		}

		// Third request with lang=en again (should hit cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test?lang=en",
		})
		if err != nil {
			t.Fatalf("Failed to make third request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			t.Errorf("Expected X-Cache-Status to be HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		if counter != 2 {
			t.Errorf("Expected counter to still be 2, got %d", counter)
		}
	})

	t.Run("Cache Invalidation", func(t *testing.T) {
		q := quick.New()

		// Add cache middleware with invalidation function
		q.Use(New(Config{
			CacheInvalidator: func(c *quick.Ctx) bool {
				return c.Query["clear"] == "1"
			},
		}))

		// Counter to track handler executions
		var counter int

		// Add a test handler that increments the counter
		q.Get("/test", func(c *quick.Ctx) error {

			counter++
			return c.String(fmt.Sprintf("Response #%d", counter))
		})

		// First request
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		// Second request (should hit cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			t.Errorf("Expected X-Cache-Status to be HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		// Third request with clear=1 (should invalidate cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test?clear=1",
		})
		if err != nil {
			t.Fatalf("Failed to make third request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "INVALIDATED" {
			t.Errorf("Expected X-Cache-Status to be INVALIDATED, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		// Fourth request (should miss cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make fourth request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}
	})

	t.Run("HTTP Method Filtering", func(t *testing.T) {
		q := quick.New()

		// Add cache middleware with method filtering
		q.Use(New(Config{
			Methods: []string{quick.MethodGet},
		}))

		// Counter to track handler executions
		var counter int

		// Add test handlers for GET and POST
		q.Get("/test", func(c *quick.Ctx) error {
			counter++
			return c.String(fmt.Sprintf("GET Response #%d", counter))
		})

		q.Post("/test", func(c *quick.Ctx) error {
			counter++
			return c.String(fmt.Sprintf("POST Response #%d", counter))
		})

		// First GET request
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make first GET request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		initialCounter := counter

		// POST request (should not be cached)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodPost,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}

		// POST should increment counter
		if counter != initialCounter+1 {
			t.Errorf("Expected counter to be %d, got %d", initialCounter+1, counter)
		}

		// Second GET request (should hit cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make second GET request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			t.Errorf("Expected X-Cache-Status to be HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		// Counter should not have incremented again
		if counter != initialCounter+1 {
			t.Errorf("Expected counter to still be %d, got %d", initialCounter+1, counter)
		}
	})

	t.Run("Custom Expiration", func(t *testing.T) {
		q := quick.New()

		// Add cache middleware with very short expiration
		q.Use(New(Config{
			Expiration: 50 * time.Millisecond,
		}))

		// Counter to track handler executions
		var counter int

		// Add a test handler that increments the counter
		q.Get("/test", func(c *quick.Ctx) error {

			counter++
			return c.String(fmt.Sprintf("Response #%d", counter))
		})

		// First request
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		// Second request (should hit cache)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			t.Errorf("Expected X-Cache-Status to be HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}

		// Wait for cache to expire
		time.Sleep(100 * time.Millisecond)

		// Third request (should miss cache due to expiration)
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/test",
		})
		if err != nil {
			t.Fatalf("Failed to make third request: %v", err)
		}

		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			t.Errorf("Expected X-Cache-Status to be MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}
	})
}

func TestCacheStorage(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		// Set a value
		cache.Set("test_key", "test_value", DefaultExpiration)

		// Get the value
		value, found := cache.Get("test_key")
		if !found {
			t.Error("Expected to find key, but it was not found")
		}

		if value != "test_value" {
			t.Errorf("Expected value to be 'test_value', got %v", value)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		// Set a value
		cache.Set("test_key", "test_value", DefaultExpiration)

		// Delete the value
		cache.Delete("test_key")

		// Try to get the deleted value
		_, found := cache.Get("test_key")
		if found {
			t.Error("Expected key to be deleted, but it was found")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		// Set a value with a short expiration
		cache.Set("test_key", "test_value", 50*time.Millisecond)

		// Get the value before expiration
		_, found := cache.Get("test_key")
		if !found {
			t.Error("Expected to find key before expiration, but it was not found")
		}

		// Wait for the value to expire
		time.Sleep(100 * time.Millisecond)

		// Try to get the expired value
		_, found = cache.Get("test_key")
		if found {
			t.Error("Expected key to be expired, but it was found")
		}
	})
}

// //---
type capture struct {
	count int
}

func TestCacheMiddleware_Complementary(t *testing.T) {
	t.Run("Cache-Control Header Should Bypass", func(t *testing.T) {
		q := quick.New()
		q.Use(New())

		var counter int
		q.Get("/nocache", func(c *quick.Ctx) error {
			counter++
			return c.String("Bypass test")
		})

		resp, err := q.Qtest(quick.QuickTestOptions{
			Method:  quick.MethodGet,
			URI:     "/nocache",
			Headers: map[string]string{"Cache-Control": "no-cache"},
		})
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.Response().Header.Get("X-Cache-Status") != "BYPASS" {
			t.Errorf("Expected X-Cache-Status BYPASS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}
		if counter != 1 {
			t.Errorf("Expected counter to be 1, got %d", counter)
		}
	})

	t.Run("Next Skips Middleware", func(t *testing.T) {
		q := quick.New()
		q.Use(New(Config{
			Next: func(c *quick.Ctx) bool {
				return true
			},
		}))

		var counter int
		q.Get("/skip", func(c *quick.Ctx) error {
			counter++
			return c.String("Next skipped")
		})

		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/skip",
		})
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if status := resp.Response().Header.Get("X-Cache-Status"); status != "" {
			t.Errorf("Expected no X-Cache-Status, got %s", status)
		}
		if counter != 1 {
			t.Errorf("Expected counter to be 1, got %d", counter)
		}
	})

	t.Run("MaxBytes Exceeded", func(t *testing.T) {
		q := quick.New()
		q.Use(New(Config{
			MaxBytes: 10,
		}))

		var counter int
		q.Get("/big", func(c *quick.Ctx) error {
			counter++
			return c.String(strings.Repeat("x", 100))
		})

		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/big",
		})
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if status := resp.Response().Header.Get("X-Cache-Status"); status != "MISS" {
			t.Errorf("Expected MISS, got %s", status)
		}

		// segunda chamada, ainda é MISS pq nunca foi cacheado
		resp, err = q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/big",
		})
		if status := resp.Response().Header.Get("X-Cache-Status"); status != "MISS" {
			t.Errorf("Expected MISS on second request, got %s", status)
		}
	})

	t.Run("Handler Error Not Cached", func(t *testing.T) {
		q := quick.New()
		q.Use(New())

		var counter int
		q.Get("/fail", func(c *quick.Ctx) error {
			counter++
			return fmt.Errorf("error forced")
		})

		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/fail",
		})
		if err != nil {
			t.Fatalf("Expected no Qtest error, got %v", err)
		}

		if status := resp.Response().Header.Get("X-Cache-Status"); status != "MISS" {
			t.Errorf("Expected MISS, got %s", status)
		}
	})

	t.Run("X-Cache-Expires-At Is Set on HIT", func(t *testing.T) {
		q := quick.New()
		q.Use(New(Config{
			Expiration: 2 * time.Second,
		}))

		q.Get("/expires", func(c *quick.Ctx) error {
			return c.String("Expires test")
		})

		// first request: MISS
		q.Qtest(quick.QuickTestOptions{Method: quick.MethodGet, URI: "/expires"})
		// second request: HIT
		resp, _ := q.Qtest(quick.QuickTestOptions{Method: quick.MethodGet, URI: "/expires"})

		if exp := resp.Response().Header.Get("X-Cache-Expires-At"); exp == "" {
			t.Error("Expected X-Cache-Expires-At to be set on HIT")
		}
	})

	t.Run("Content-Type fallback without StoreResponseHeaders", func(t *testing.T) {
		q := quick.New()
		q.Use(New(Config{
			StoreResponseHeaders: false,
		}))

		q.Get("/type", func(c *quick.Ctx) error {
			return c.String("text plain fallback")
		})

		// miss
		q.Qtest(quick.QuickTestOptions{Method: quick.MethodGet, URI: "/type"})
		// hit
		resp, _ := q.Qtest(quick.QuickTestOptions{Method: quick.MethodGet, URI: "/type"})
		if ct := resp.Response().Header.Get("Content-Type"); ct == "" {
			t.Error("Expected Content-Type to be set on fallback")
		}
	})
}

// ----
func TestResponseCapture_HeaderHelpers(t *testing.T) {
	rec := httptest.NewRecorder()
	rc := &responseCapture{
		ResponseWriter: rec,
		buffer:         bytes.NewBuffer(nil),
		headers:        make(http.Header),
	}

	// Set
	rc.Set("X-One", "1")
	if got := rc.headers.Get("X-One"); got != "1" {
		t.Fatalf("Set failed, got %s", got)
	}
	if rec.Header().Get("X-One") != "1" {
		t.Fatalf("Set did not propagate to ResponseWriter")
	}

	// Add
	rc.Add("X-One", "2")
	want := []string{"1", "2"}
	if !reflect.DeepEqual(rc.headers["X-One"], want) {
		t.Fatalf("Add failed, got %v", rc.headers["X-One"])
	}

	// Del
	rc.Del("X-One")
	if _, ok := rc.headers["X-One"]; ok {
		t.Fatalf("Del failed, header still present")
	}
	if rec.Header().Get("X-One") != "" {
		t.Fatalf("Del failed on ResponseWriter")
	}
}

func TestResponseCapture_WriteHeaderAndImplicit(t *testing.T) {
	// explicit WriteHeader path
	rec := httptest.NewRecorder()
	rc := &responseCapture{
		ResponseWriter: rec,
		buffer:         bytes.NewBuffer(nil),
		headers:        make(http.Header),
	}
	rc.Set("X-Test", "ok")
	rc.WriteHeader(201)

	if rc.statusCode != 201 || !rc.headerWritten {
		t.Fatalf("WriteHeader did not set status or flag")
	}
	if rec.Result().StatusCode != 201 {
		t.Fatalf("WriteHeader not propagated, got %d", rec.Result().StatusCode)
	}
	if rec.Header().Get("X-Test") != "ok" {
		t.Fatalf("Headers not copied to ResponseWriter")
	}

	// implicit path via Write()
	rec2 := httptest.NewRecorder()
	rc2 := &responseCapture{
		ResponseWriter: rec2,
		buffer:         bytes.NewBuffer(nil),
		headers:        make(http.Header),
	}
	rc2.Set("Content-Type", "text/plain")
	rc2.Write([]byte("hi"))

	if rc2.statusCode != http.StatusOK || rec2.Code != http.StatusOK {
		t.Fatalf("implicit WriteHeader failed")
	}
	if rec2.Body.String() != "hi" {
		t.Fatalf("body not written")
	}
}

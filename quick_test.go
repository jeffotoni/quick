package quick

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestNew ensures that creating a new Quick instance returns a valid object
// The result will TestNew(expected any) error
func TestNew(t *testing.T) {
	t.Run("New instance with default config", func(t *testing.T) {
		q := New()
		if q == nil {
			t.Fatal("Expected non-nil Quick instance")
		}
	})

	t.Run("New instance with custom config", func(t *testing.T) {
		config := Config{RouteCapacity: 500}
		q := New(config)
		if q.config.RouteCapacity != 500 {
			t.Errorf("Expected RouteCapacity to be 500, got %d", q.config.RouteCapacity)
		}
	})
}

// TestUseCorsMiddleware test if CORS middleware is applied correctly.
// The result will TestUseCorsMiddleware(expected any) error
func TestUseCorsMiddleware(t *testing.T) {
	t.Run("Apply CORS middleware", func(t *testing.T) {
		q := New()
		middleware := func(h http.Handler) http.Handler { return h }
		q.Use(middleware, "cors")

		if !q.Cors {
			t.Errorf("Expected q.Cors to be true")
		}
		if q.CorsSet == nil {
			t.Errorf("Expected q.CorsSet to be set")
		}
	})
}

// TestExtractHandler test if an unknown HTTP method returns nil
// The result will TestExtractHandler(expected any) error
func TestExtractHandler(t *testing.T) {
	t.Run("Returns nil for unknown method", func(t *testing.T) {
		handler := extractHandler(New(), "UNKNOWN", "/path", "", func(c *Ctx) error { return nil })
		if handler != nil {
			t.Errorf("Expected nil handler for unknown method")
		}
	})
}

// TestExtractParamsPattern test if extracting parameters from a route works correctly
// The result will TestExtractParamsPattern(expected any) error
func TestExtractParamsPattern(t *testing.T) {
	t.Run("Index == 1 case", func(t *testing.T) {
		path, _, _ := extractParamsPattern("/:id")
		if path != "/" {
			t.Errorf("Expected '/', but got '%s'", path)
		}
	})
}

// TestExtractParamsGet Test handling of a GET request with nil context
// The result will TestExtractParamsGet(expected any) error
func TestExtractParamsGet(t *testing.T) {
	t.Run("Handles nil context value", func(t *testing.T) {
		q := New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		handler := extractParamsGet(q, "/test", "", func(c *Ctx) error { return nil })

		handler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", w.Code)
		}
	})
}

// TestExtractParamsPost test handling of a POST request with oversized body
// The result will TestExtractParamsPost(expected any) error
func TestExtractParamsPost(t *testing.T) {
	t.Run("Handles nil context value", func(t *testing.T) {
		q := New()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		w := httptest.NewRecorder()
		handler := extractParamsPost(q, func(c *Ctx) error { return nil })

		handler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", w.Code)
		}
	})

	t.Run("Rejects oversized body", func(t *testing.T) {
		q := New(Config{MaxBodySize: 10})
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(make([]byte, 20)))
		req.ContentLength = 20
		w := httptest.NewRecorder()
		handler := extractParamsPost(q, func(c *Ctx) error { return nil })

		handler(w, req)

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status 413 Request Entity Too Large, got %d", w.Code)
		}
	})
}

// TestExtractParamsOptions test handling of an OPTIONS request with an error in the handler
// The result will TestExtractParamsOptions(expected any) error
func TestExtractParamsOptions(t *testing.T) {
	t.Run("Handles error in handler function", func(t *testing.T) {
		q := New()
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		w := httptest.NewRecorder()
		handler := extractParamsOptions(q, http.MethodOptions, "/test", func(c *Ctx) error {
			return errors.New("internal server error")
		})

		handler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 Internal Server Error, got %d", w.Code)
		}
	})

	t.Run("Handles empty handler function", func(t *testing.T) {
		q := New()
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		w := httptest.NewRecorder()
		handler := extractParamsOptions(q, http.MethodOptions, "/test", nil)

		handler(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204 No Content, got %d", w.Code)
		}
	})
}

// TestExtractParamsPut test if a PUT request rejects an oversized body.
// The result will TestExtractParamsPut(expected any) error
func TestExtractParamsPut(t *testing.T) {
	t.Run("Handles nil context value", func(t *testing.T) {
		q := New()
		req := httptest.NewRequest(http.MethodPut, "/test", nil)
		w := httptest.NewRecorder()
		handler := extractParamsPut(q, func(c *Ctx) error { return nil })

		handler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", w.Code)
		}
	})

	t.Run("Rejects oversized body", func(t *testing.T) {
		q := New(Config{MaxBodySize: 10})
		req := httptest.NewRequest(http.MethodPut, "/test", bytes.NewReader(make([]byte, 20)))
		req.ContentLength = 20
		w := httptest.NewRecorder()
		handler := extractParamsPut(q, func(c *Ctx) error { return nil })

		handler(w, req)

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status 413 Request Entity Too Large, got %d", w.Code)
		}
	})
}

// TestExtractParamsDelete test if a DELETE request with nil context is handled properly.
// The result will TestExtractParamsDelete(expected any) error
func TestExtractParamsDelete(t *testing.T) {
	t.Run("Handles nil context value", func(t *testing.T) {
		q := New()
		req := httptest.NewRequest(http.MethodDelete, "/test", nil)
		w := httptest.NewRecorder()
		handler := extractParamsDelete(q, func(c *Ctx) error { return nil })

		handler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", w.Code)
		}
	})
}

// TestExecHandleFunc test if a TLS server starts and responds correctly.
// The result will TestExecHandleFunc(expected any) error
func TestExecHandleFunc(t *testing.T) {
	t.Run("Executes handler function successfully", func(t *testing.T) {
		c := &Ctx{Response: httptest.NewRecorder()}
		execHandleFunc(c, func(ctx *Ctx) error {
			ctx.Response.WriteHeader(http.StatusOK)
			return nil
		})
		if c.Response.(*httptest.ResponseRecorder).Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", c.Response.(*httptest.ResponseRecorder).Code)
		}
	})

	t.Run("Handles handler function error", func(t *testing.T) {
		c := &Ctx{Response: httptest.NewRecorder()}
		execHandleFunc(c, func(ctx *Ctx) error {
			return errors.New("test error")
		})
		if c.Response.(*httptest.ResponseRecorder).Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 Internal Server Error, got %d", c.Response.(*httptest.ResponseRecorder).Code)
		}
	})
}

// TestStatic tests if invalid parameters cause a panic
// The result will TestStatic(expected any) error
func TestStatic(t *testing.T) {
	t.Run("Panic on invalid parameter", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic but got none")
			}
		}()
		q := New()
		q.Static("/static", 123) // Invalid parameter
	})
}

// TestMWWrapper_Cover tests if middleware is correctly applied
// The result will TestMWWrapper_Cover(expected any) error
func TestMWWrapper_Cover(t *testing.T) {
	t.Run("Applies middleware", func(t *testing.T) {
		q := New()
		mw := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Test", "true")
				h.ServeHTTP(w, r)
			})
		}
		q.Use(mw)
		handler := q.mwWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(w, req)

		if w.Header().Get("X-Test") != "true" {
			t.Errorf("Expected X-Test header to be 'true', got '%s'", w.Header().Get("X-Test"))
		}
	})
}

// TestListenWithShutdown_Cover tests if the server starts and shuts down correctly
// The result will TestListenWithShutdown_Cover(expected any) error
func TestListenWithShutdown_Cover(t *testing.T) {
	t.Run("Starts and shuts down server", func(t *testing.T) {
		q := New()
		server, shutdown, err := q.ListenWithShutdown(":0")
		if err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}
		if server == nil {
			t.Fatal("Expected non-nil server")
		}

		shutdown() // Ensure server shuts down
	})
}

// TestListenWithShutdown tests if the server starts and shuts down with PRINT_SERVER=true
// The result will TestListenWithShutdown(expected any) error
func TestListenWithShutdown(t *testing.T) {
	t.Run("Starts and shuts down server with PRINT_SERVER=true", func(t *testing.T) {
		t.Setenv("PRINT_SERVER", "true")

		q := New()
		server, shutdown, err := q.ListenWithShutdown(":0")
		if err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}
		if server == nil {
			t.Fatal("Expected non-nil server")
		}

		shutdown() // Ensure server shuts down
	})
}

// TestExecHandler tests if a wrapped handler executes successfully
// The result will TestExecHandler(expected any) error
func TestExecHandler(t *testing.T) {
	t.Run("Executes wrapped handler", func(t *testing.T) {
		q := New()
		handler := q.execHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", w.Code)
		}
	})
}

// TestCorsHandler tests if the CORS middleware is applied correctly
// The result will TestCorsHandler(expected any) error
func TestCorsHandler(t *testing.T) {
	t.Run("Applies CORS middleware", func(t *testing.T) {
		q := New()
		q.CorsSet = func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				h.ServeHTTP(w, r)
			})
		}

		handler := q.corsHandler()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected CORS header to be '*', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
		}
	})
}

// TestHttpServer tests if the server is correctly created with the expected handler
// The result will TestHttpServer(expected any) error
func TestHttpServer(t *testing.T) {
	t.Run("Creates server with execHandler", func(t *testing.T) {
		q := New()
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		server := q.httpServer(":8080", h)

		if server.Handler == nil {
			t.Errorf("Expected handler to be set")
		}
	})

	t.Run("Creates server with CORS handler", func(t *testing.T) {
		q := New()
		q.Cors = true
		q.CorsSet = func(h http.Handler) http.Handler { return h }

		server := q.httpServer(":8080")

		if server.Handler == nil {
			t.Errorf("Expected handler to be set")
		}
	})
}

// TestCreateParamsAndValid tests if regex parameters are correctly extracted from URLs
// The result will TestCreateParamsAndValid(expected any) error
func TestCreateParamsAndValid(t *testing.T) {
	t.Run("Handles regex parameter", func(t *testing.T) {
		params, valid := createParamsAndValid("/users/123", "/users/{id:[0-9]+}")
		if !valid {
			t.Errorf("Expected valid match")
		}
		if params["id"] != "123" {
			t.Errorf("Expected extracted param '123', got '%s'", params["id"])
		}
	})
}

// errorReader is an io.Reader that always returns an error.
type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

// TestMWWrapper tests if a middleware function is correctly applied
// The result will TestMWWrapper(expected any) error
func TestMWWrapper(t *testing.T) {
	t.Run("Applies middleware function", func(t *testing.T) {
		q := New()
		mw := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Test", "true")
				h.ServeHTTP(w, r)
			})
		}
		q.Use(mw)
		handler := q.mwWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(w, req)

		if w.Header().Get("X-Test") != "true" {
			t.Errorf("Expected X-Test header to be 'true', got '%s'", w.Header().Get("X-Test"))
		}
	})
}

// TestExtractBodyBytes tests if an error is handled correctly when reading the request body
// The result will TestExtractBodyBytes(expected any) error
func TestExtractBodyBytes(t *testing.T) {
	t.Run("Handles error reading body", func(t *testing.T) {
		// Force read error
		r := io.NopCloser(errorReader{})

		result, _ := extractBodyBytes(r)

		if result != nil {
			t.Errorf("Expected nil, but got %v", result)
		}
	})
}

// TestListen verifies if the server starts and runs successfully
// The result will TestListen(expected any) error
func TestListen(t *testing.T) {
	t.Run("Listen starts and runs", func(t *testing.T) {
		q := New()
		go func() {
			_ = q.Listen(":0")
		}()

		time.Sleep(500 * time.Millisecond) // Give time for the server to start
	})
}

// TestCreateParamsAndValid_RegexScenarios tests various regex matching cases for URL parameters
// The result will TestCreateParamsAndValid_RegexScenarios(expected any) error
func TestCreateParamsAndValid_RegexScenarios(t *testing.T) {
	t.Run("Fails on non-matching numeric regex", func(t *testing.T) {
		// It should fail because 'abc' does not match [0-9]+
		_, valid := createParamsAndValid("/users/abc", "/users/{id:[0-9]+}")
		if valid {
			t.Error("Expected invalid match for '/users/abc' with '/users/{id:[0-9]+}'")
		}
	})

	t.Run("Matches alphabetic regex", func(t *testing.T) {
		// {slug:[a-z]+} must match 'golang'
		params, valid := createParamsAndValid("/profile/golang", "/profile/{slug:[a-z]+}")
		if !valid {
			t.Error("Expected valid match for '/profile/golang' with '/profile/{slug:[a-z]+'")
		}
		if params["slug"] != "golang" {
			t.Errorf("Expected 'golang', got '%s'", params["slug"])
		}
	})

	t.Run("Fails alphabetic regex if uppercase is present", func(t *testing.T) {
		// 'Golang' has a capital letter, it should not match [a-z]+
		_, valid := createParamsAndValid("/profile/Golang", "/profile/{slug:[a-z]+}")
		if valid {
			t.Error("Expected invalid match, but it was valid")
		}
	})

	t.Run("Handles multiple regex segments", func(t *testing.T) {
		// Example: /api/v1/users/123 => /api/{version:v[0-9]+}/users/{id:[0-9]+}
		params, valid := createParamsAndValid("/api/v1/users/123", "/api/{version:v[0-9]+}/users/{id:[0-9]+}")
		if !valid {
			t.Error("Expected valid match for multiple segments")
		}
		if params["version"] != "v1" {
			t.Errorf("Expected 'v1', got '%s'", params["version"])
		}
		if params["id"] != "123" {
			t.Errorf("Expected '123', got '%s'", params["id"])
		}
	})

	t.Run("Handles empty param name but valid regex", func(t *testing.T) {
		// Example: /number/123 => /number/{:[0-9]+}
		// Note: paramName is empty before the ':', can we use it?
		_, valid := createParamsAndValid("/number/123", "/number/{:[0-9]+}")
		// We expect failure because the logic requires name and regex
		if valid {
			t.Error("Expected invalid match because there's no param name before the colon")
		}
	})
}

// TestMWWrapper_CustomMiddleware tests if a custom middleware function is correctly applied
// The result will TestMWWrapper_CustomMiddleware(expected any) error
func TestMWWrapper_CustomMiddleware(t *testing.T) {
	q := New()

	// Middleware in the format `func(http.ResponseWriter, *http.Request, http.Handler)`
	middleware := func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		w.Header().Set("X-Test", "middleware-invoked")
		next.ServeHTTP(w, r)
	}

	// Add middleware to Quick
	q.Use(middleware)

	// Involves a base handler with the middleware
	handler := q.mwWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Simulate a request
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, req)

	// Checks if the middleware actually added the header
	if w.Header().Get("X-Test") != "middleware-invoked" {
		t.Errorf("Expected header 'X-Test' to be set, got '%s'", w.Header().Get("X-Test"))
	}

	// Checks if the answer was correct
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

// TestListen_ErrorCase tests if an error is returned when trying to listen on an occupied port
// The result will TestListen_ErrorCase(expected any) error
func TestListen_ErrorCase(t *testing.T) {
	q := New()

	// Criar um servidor na mesma porta para simular erro
	occupiedListener, _ := net.Listen("tcp", ":8089")
	defer occupiedListener.Close()

	err := q.Listen(":8089") // Deve falhar pois a porta já está em uso

	if err == nil {
		t.Errorf("Expected an error when trying to listen on an occupied port, but got nil")
	}
}

// TestCreateParamsAndValid_NoMatch tests if non-matching routes return nil and false
// The result will TestCreateParamsAndValid_NoMatch(expected any) error
func TestCreateParamsAndValid_NoMatch(t *testing.T) {
	t.Run("Returns nil and false for non-matching routes", func(t *testing.T) {
		params, valid := createParamsAndValid("/wrong/path", "/expected/{id:[0-9]+}")

		if valid {
			t.Errorf("Expected valid to be false, but got true")
		}
		if params != nil {
			t.Errorf("Expected params to be nil, but got %v", params)
		}
	})
}

// TestCreateParamsAndValid_EmptyParamName tests if an empty param name returns nil and false
// The result will TestCreateParamsAndValid_EmptyParamName(expected any) error
func TestCreateParamsAndValid_EmptyParamName(t *testing.T) {
	t.Run("Returns nil and false for empty param name", func(t *testing.T) {
		params, valid := createParamsAndValid("/users/123", "/users/:")

		if valid {
			t.Errorf("Expected valid to be false, but got true")
		}
		if params != nil {
			t.Errorf("Expected params to be nil, but got %v", params)
		}
	})
}

// TestCreateParamsAndValid_BuilderMismatch tests if mismatched reconstructed paths return nil and false
// The result will TestCreateParamsAndValid_BuilderMismatch(expected any) error
func TestCreateParamsAndValid_BuilderMismatch(t *testing.T) {
	t.Run("Returns nil and false if reconstructed path does not match", func(t *testing.T) {
		params, valid := createParamsAndValid("/users/123", "/users/456")

		if valid {
			t.Errorf("Expected valid to be false, but got true")
		}
		if params != nil {
			t.Errorf("Expected params to be nil, but got %v", params)
		}
	})
}

// TestCreateParamsAndValid_PathMismatch tests if mismatched paths return nil and false
// The result will TestCreateParamsAndValid_PathMismatch(expected any) error
func TestCreateParamsAndValid_PathMismatch(t *testing.T) {
	t.Run("Returns nil and false when the reconstructed path does not match the request URI", func(t *testing.T) {
		params, valid := createParamsAndValid("/users/123/profile", "/users/:id/settings")

		if valid {
			t.Errorf("Expected valid to be false, but got true")
		}
		if params != nil {
			t.Errorf("Expected params to be nil, but got %v", params)
		}
	})
}

// TestRegisterRoute_WithRegexParamPanic tests if the server correctly registers a route with regex parameters
// The result will TestRegisterRoute_WithRegexParamPanic(expected any) error
func TestRegisterRoute_WithRegexParamPanic(t *testing.T) {
	q := New()

	// This record should trigger the panic
	q.Get("/v1/user/{id:[0-9]+}", func(c *Ctx) error {
		return c.Status(200).String("Hello Quick!")
	})

	// Simulates the call from /v1/user/123
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/user/123", nil)

	// Execute the route
	q.ServeHTTP(w, req)

	// Checks if the status is 200
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Checks if the body is "hello Quick!"
	if w.Body.String() != "Hello Quick!" {
		t.Errorf("Expected response body 'hello Quick!', got '%s'", w.Body.String())
	}
}

// $ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
// $ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
// TestListenTLS tests if a TLS server starts and responds correctly
// The result will TestListenTLS(expected any) error
func TestListenTLS(t *testing.T) {
	q := New()

	// Define a simple handler for the root "/"
	q.Get("/", func(c *Ctx) error {
		return c.Status(200).String("OK")
	})

	go func() {
		err := q.ListenTLS(":8443", "cert.pem", "key.pem", false)
		if err != nil {
			t.Errorf("error starting TLS server: %v", err)
		}
	}()

	// Short timeout to make sure the server is up
	time.Sleep(500 * time.Millisecond)

	// Create an HTTP client to test the TLS connection
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip SSL verification
		},
	}

	resp, err := client.Get("https://localhost:8443/")
	if err != nil {
		t.Fatalf("failed to connect to TLS server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Shut down the server at the end of the test
	_ = q.Shutdown()
}

// TestListenTLSH2 tests if a TLS server starts and responds correctly
// $ curl -i -X GET -H "Content-Type: application/json" --http2 -v https://localhost:443/v1/user
// $ curl -i -X GET -H "Content-Type: application/json" --http2 -k https://localhost:443/v1/user
// The result will TestListenTLS(expected any) error
func TestListenTLSH2(t *testing.T) {
	q := New()

	// Define a simple handler for the root "/"
	q.Get("/", func(c *Ctx) error {
		return c.Status(200).String("OK")
	})

	go func() {
		err := q.ListenTLS(":8443", "cert.pem", "key.pem", true)
		if err != nil {
			t.Errorf("error starting TLS server: %v", err)
		}
	}()

	// Short timeout to make sure the server is up
	time.Sleep(500 * time.Millisecond)

	// Create an HTTP client to test the TLS connection
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip SSL verification
		},
	}

	resp, err := client.Get("https://localhost:8443/")
	if err != nil {
		t.Fatalf("failed to connect to TLS server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Shut down the server at the end of the test
	_ = q.Shutdown()
}

// The result will TestResponseWriterPool
func TestResponseWriterPool(t *testing.T) {
	w := httptest.NewRecorder()

	rw := acquireResponseWriter(w)
	if rw.ResponseWriter != w {
		t.Errorf("Expected ResponseWriter to match original")
	}
	if rw.buf == nil {
		t.Errorf("Buffer should be initialized")
	}

	rw.buf.WriteString("test")
	if rw.buf.Len() == 0 {
		t.Errorf("Expected buffer length > 0 after WriteString")
	}

	releaseResponseWriter(rw)

	// Acquire again to ensure buffer reset
	rw2 := acquireResponseWriter(w)
	if rw2.buf.Len() != 0 {
		t.Errorf("Buffer was not reset, got len %d", rw2.buf.Len())
	}
	releaseResponseWriter(rw2)
}

// The result will TestHttpServerTLS
func TestHttpServerTLS(t *testing.T) {
	q := New()
	q.Cors = true

	// Dummy handler to test coverage
	srv := q.httpServerTLS(":8080", &tls.Config{}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	if srv, ok := srv.Handler.(http.Handler); !ok || srv == nil {
		t.Errorf("Expected non-nil handler, got nil")
	}
}

// The result will TestQuick_ListenTLS_GCPercent
func TestQuick_ListenTLS_GCPercent(t *testing.T) {
	q := New()
	q.config.GCPercent = 100

	go func() {
		err := q.ListenTLS(":8081", "cert.pem", "key.pem", false)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			t.Errorf("Expected file error, got: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	// Check GC percentage
	if gc := debug.SetGCPercent(-1); gc != q.config.GCPercent {
		t.Errorf("Expected GCPercent %d, got %d", q.config.GCPercent, gc)
	}
}

// The result will TestQuick_StartServerWithGracefulShutdown
func TestQuick_StartServerWithGracefulShutdown(t *testing.T) {
	listener, _ := net.Listen("tcp", "localhost:0")
	defer listener.Close()

	q := &Quick{
		server: &http.Server{},
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	err := q.startServerWithGracefulShutdown(listener, "cert.pem", "key.pem")
	if err != nil {
		t.Errorf("Graceful shutdown failed: %v", err)
	}
}

// The result will TestQuick_Shutdown
func TestQuick_Shutdown(t *testing.T) {
	q := New()
	err := q.Shutdown()
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
}

// The result will TestExtractParamsBind_UnsupportedContentType
func TestExtractParamsBind_UnsupportedContentType(t *testing.T) {
	// Prepare the request with an unsupported content-type
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Quick"}`))
	req.Header.Set("Content-Type", "text/plain") // unsupported content-type

	// Prepare context
	c := &Ctx{
		Request: req,
	}

	var target map[string]string
	err := extractParamsBind(c, &target)
	if err == nil {
		t.Fatal("Expected an error due to unsupported content type, got nil")
	}

	expectedError := "unsupported content type"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error '%s', got '%v'", expectedError, err)
	}
}

// The result will TestExtractParamsBind_InvalidXMLContent
func TestExtractParamsBind_InvalidXMLContent(t *testing.T) {
	// Prepare the request with invalid XML content
	body := strings.NewReader("<invalid><xml>")
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "application/xml")

	ctx := &Ctx{
		Request: req,
	}

	var v interface{}
	err := extractParamsBind(ctx, &v)
	if err == nil {
		t.Fatal("Expected an XML parsing error, got nil")
	}
}

type errReader struct{}

// The result will Read(p []byte) (n int, err error)
func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

// The result will TestExtractParamsBind_BodyReadError
func TestExtractParamsBind_BodyReadError(t *testing.T) {
	errReadCloser := ioutil.NopCloser(&errReader{})
	req := httptest.NewRequest("POST", "/", errReadCloser)
	req.Header.Set("Content-Type", ContentTypeAppJSON)

	ctx := &Ctx{Request: req}

	var v interface{}
	err := extractParamsBind(ctx, &v)
	if err == nil {
		t.Fatal("Expected an error due to body read failure, got nil")
	}
}

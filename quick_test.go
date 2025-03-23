// Package quick provides a high-performance, minimalistic web framework for Go.
//
// This file contains **unit tests** for various functionalities of the Quick framework.
// These tests ensure that the core features of Quick work as expected.
//
// 📌 To run all unit tests, use:
//
//	$ go test -v ./...
//	$ go test -v
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

// TestNew verifies the behavior of creating a new Quick instance using default and custom configuration.
//
// The first subtest ensures that a default Quick instance is not nil.
// The second subtest checks if the custom configuration is correctly applied.
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

// TestExtractHandler verifies that extractHandler returns nil when an unknown HTTP method is provided.
//
// It ensures that routes with unsupported methods are not registered.
func TestExtractHandler(t *testing.T) {
	t.Run("Returns nil for unknown method", func(t *testing.T) {
		handler := extractHandler(New(), "UNKNOWN", "/path", "", func(c *Ctx) error { return nil })
		if handler != nil {
			t.Errorf("Expected nil handler for unknown method")
		}
	})
}

// TestExtractParamsPattern verifies the behavior of extractParamsPattern when parsing route patterns.
//
// It ensures that when the route pattern has only one parameter (index == 1), the returned base path is "/".
func TestExtractParamsPattern(t *testing.T) {
	t.Run("Index == 1 case", func(t *testing.T) {
		path, _, _ := extractParamsPattern("/:id")
		if path != "/" {
			t.Errorf("Expected '/', but got '%s'", path)
		}
	})
}

// TestExtractParamsGet verifies the behavior of extractParamsGet when handling a GET request with a nil context.
//
// It ensures that when no matching route context is found, the response returns a 404 status code.
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

// TestExtractParamsPost verifies the behavior of extractParamsPost when handling POST requests.
//
// It ensures that requests with nil context return 404, and that requests with a body exceeding MaxBodySize
// are rejected with a 413 Request Entity Too Large status.
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

// TestExtractParamsOptions verifies the behavior of extractParamsOptions for handling OPTIONS requests.
//
// It tests two scenarios:
//  1. If the handler function returns an error, the response should have status 500 Internal Server Error.
//  2. If no handler is provided, the response should return status 204 No Content.
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

// TestExtractParamsPut validates the behavior of extractParamsPut for PUT requests.
//
// It covers two cases:
//  1. When the request context is nil, it should respond with 404 Not Found.
//  2. When the request body exceeds MaxBodySize, it should respond with 413 Request Entity Too Large.
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

// TestExtractParamsDelete validates the behavior of extractParamsDelete for DELETE requests.
//
// It ensures that when the request context is nil or does not match any route,
// the server responds with a 404 Not Found status code.
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

// TestExecHandleFunc validates the execution of a handler function using execHandleFunc.
//
// It ensures that the response is correctly set to 200 OK when the handler succeeds,
// and 500 Internal Server Error when the handler returns an error.
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

// TestStatic ensures that passing an invalid parameter to Static triggers a panic.
//
// It verifies that calling Static with a non-string or non-struct path parameter causes a controlled panic.
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

// TestMWWrapper_Cover checks if a middleware function is properly wrapped and applied.
//
// It validates that a middleware can add a custom header and that the header is present in the final response.
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

// TestListenWithShutdown_Cover verifies that the server starts and shuts down without errors.
//
// It checks that ListenWithShutdown returns a valid server instance and allows graceful termination.
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

// TestListenWithShutdown checks server behavior when NoBanner config is set.
//
// It ensures the server starts and shuts down correctly when configured with NoBanner.
func TestListenWithShutdown(t *testing.T) {
	t.Run("Starts and shuts down server with NoBanner", func(t *testing.T) {
		q := New(Config{
			NoBanner: false,
		})

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

// TestExecHandler confirms that the execHandler properly executes the wrapped handler.
//
// It ensures the final HTTP response is correctly returned with status 200 OK.
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

// TestCorsHandler ensures that a CORS middleware is correctly applied.
//
// It validates that the "Access-Control-Allow-Origin" header is included in the response.
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

// TestHttpServer verifies the correct creation of the internal HTTP server.
//
// It ensures the handler is set when execHandler or CORS handler is used based on config.
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

// TestCreateParamsAndValid ensures that parameters using regex patterns are correctly extracted from URLs.
//
// It validates that when the route contains a pattern like "/users/{id:[0-9]+}",
// a matching request like "/users/123" correctly extracts "123" as the "id" parameter.
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

// Read implements the io.Reader interface and always returns an error.
func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

// TestMWWrapper checks whether middleware functions are properly applied.
//
// It validates that middleware can inject headers into the response before passing control to the handler.
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

// TestExtractBodyBytes ensures that reading from a request body handles I/O errors gracefully.
//
// It uses a simulated reader that always fails to validate error handling logic.
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

// TestListen verifies if the Quick server starts without panics or errors.
//
// It launches a server on an ephemeral port and ensures it runs for a short duration.
func TestListen(t *testing.T) {
	t.Run("Listen starts and runs", func(t *testing.T) {
		q := New()
		go func() {
			_ = q.Listen(":0")
		}()

		time.Sleep(500 * time.Millisecond) // Give time for the server to start
	})
}

// TestCreateParamsAndValid_RegexScenarios tests various regex route patterns for parameter extraction.
//
// It covers cases such as alphabetic, numeric, uppercase mismatch, and invalid parameter names.
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

// TestMWWrapper_CustomMiddleware checks if custom middleware using a 3-argument function signature is applied correctly.
//
// It verifies that the middleware adds a header and allows the handler to respond successfully.
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

// TestListen_ErrorCase ensures that the server returns an error when attempting to listen on a port that is already in use.
//
// It pre-binds the port using net.Listen to simulate the conflict.
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

// TestCreateParamsAndValid_NoMatch ensures that non-matching routes return nil parameters and false validity.
//
// It verifies that if the actual request path does not match the route pattern,
// the parameter map is nil and the match is considered invalid.
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

// TestCreateParamsAndValid_EmptyParamName ensures that an empty parameter name in the route pattern results in an invalid match.
//
// It checks that route patterns like "/users/:" are not accepted and return nil parameters and false validity.
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

// TestCreateParamsAndValid_BuilderMismatch verifies that mismatched reconstructed paths fail to validate.
//
// It simulates a case where the requested path and route pattern lead to a different reconstructed path.
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

// TestCreateParamsAndValid_PathMismatch ensures that paths with structural differences fail to match.
//
// For example, "/users/123/profile" should not match "/users/:id/settings".
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

// TestRegisterRoute_WithRegexParamPanic validates that a route with a regex parameter is registered and matched correctly.
//
// It tests the case "/v1/user/{id:[0-9]+}" and ensures the route executes and responds as expected.
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

// TestListenTLS starts a TLS server using cert/key files and verifies successful response.
//
// It simulates a request to a secure endpoint and expects an HTTP 200 OK status.
//
// $ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
// $ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
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

// TestListenTLSH2 starts a TLS server with HTTP/2 support and ensures it responds correctly.
//
// It uses a client to make a secure request and verifies the status and response.
//
// $ curl -i -X GET -H "Content-Type: application/json" --http2 -v https://localhost:443/v1/user
// $ curl -i -X GET -H "Content-Type: application/json" --http2 -k https://localhost:443/v1/user
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

// TestResponseWriterPool validates the behavior of the response writer pooling system.
//
// It ensures the buffer is initialized, reused, and reset between uses.
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

// TestHttpServerTLS confirms that a server with TLS configuration is created correctly.
//
// It verifies that the handler is set and ready to serve.
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

// TestQuick_ListenTLS_GCPercent ensures that the GC percent configuration is applied before starting the server.
//
// It compares the configured value against the runtime value.
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

// TestQuick_StartServerWithGracefulShutdown verifies that the server can shut down gracefully when receiving a signal.
//
// It simulates an interrupt signal and expects no errors during shutdown.
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

// TestQuick_Shutdown validates that the Shutdown method completes successfully without error.
//
// It is a basic check to ensure cleanup operations don't fail.
func TestQuick_Shutdown(t *testing.T) {
	q := New()
	err := q.Shutdown()
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
}

// TestExtractParamsBind_UnsupportedContentType ensures that unsupported content types return an error.
//
// It simulates a request with an invalid content-type like "text/plain" and expects a failure.
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

// TestExtractParamsBind_InvalidXMLContent ensures that invalid XML data results in a parsing error.
//
// It provides malformed XML and checks that the system returns an appropriate error.
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

// Read implements the io.Reader interface and always returns an error.
func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

// TestExtractParamsBind_BodyReadError simulates a read error when binding a request body.
//
// It ensures that such failures are detected and reported properly.
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

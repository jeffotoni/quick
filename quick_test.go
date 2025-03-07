package quick

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

func TestExtractHandler(t *testing.T) {
	t.Run("Returns nil for unknown method", func(t *testing.T) {
		handler := extractHandler(New(), "UNKNOWN", "/path", "", func(c *Ctx) error { return nil })
		if handler != nil {
			t.Errorf("Expected nil handler for unknown method")
		}
	})
}

func TestExtractParamsPattern(t *testing.T) {
	t.Run("Index == 1 case", func(t *testing.T) {
		path, _, _ := extractParamsPattern("/:id")
		if path != "/" {
			t.Errorf("Expected '/', but got '%s'", path)
		}
	})
}

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

func TestListen(t *testing.T) {
	t.Run("Listen starts and runs", func(t *testing.T) {
		q := New()
		go func() {
			_ = q.Listen(":0")
		}()

		time.Sleep(500 * time.Millisecond) // Give time for the server to start
	})
}

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

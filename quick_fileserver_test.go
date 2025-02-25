package quick

import (
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// traditional test
///

// TestQuickStatic Tests if the static/* server functionality redirects correctly to index.html
// The will test TestQuickStatic(t *testing.T)
//
// Run:
//
//	$ go test -v -run ^TestQuickStatic
func TestQuickStatic(t *testing.T) {
	q := New()

	// Configure static file server from the "./static" directory
	q.Static("/static", "./static")

	// Define a route that serves static files
	q.Get("/", func(c *Ctx) error {
		c.File("static/*") // Testing if `static/index.html` is found
		return nil
	})

	// Creating a test server
	server := httptest.NewServer(q)
	defer server.Close()

	// Makes a GET request to "/"
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Checks if the response is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, but received: %d", resp.StatusCode)
	}

	// Read the response content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	// Check if the response contains any content expected from index.html
	expectedContent := "<h1>File Server Go example html</h1>" // Example: if index.html has a <title> tag
	if !strings.Contains(string(body), expectedContent) {
		t.Errorf("Expected to find '%s' in the content, but did not find it", expectedContent)
	}
}

// Table-driven test
// /
//
//go:embed static/*
var staticFiles embed.FS

// TestQuickStaticDriven Tests if the static/* server functionality redirects correctly to index.html
// The will test TestQuickStaticDriven(t *testing.T)
//
// Run:
//
//	$ go test -v -run ^TestQuickStaticDriven
func TestQuickStaticDriven(t *testing.T) {
	tests := []struct {
		name       string // Test case description
		useEmbed   bool   // Whether to use embedded files or local file system
		path       string // Path to test
		statusCode int    // Expected HTTP status code
		expectBody string // Expected content in the response
	}{
		{"Serve index.html from file system", false, "/", http.StatusOK, "<h1>File Server Go example html</h1>"},
		{"Serve static/index.html directly from file system", false, "/static/index.html", StatusNotFound, "404"},
		{"Arquivo not found from file system", false, "/static/missing.html", http.StatusNotFound, "404"},
		{"Serve index.html from embed FS", true, "/", http.StatusOK, "<h1>File Server Go example html</h1>"},
		{"Serve static/index.html directly from embed FS", true, "/static/index.html", http.StatusNotFound, "404"},
		{"Arquivo not found from embed FS", true, "/static/missing.html", http.StatusNotFound, "404"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := New()

			// Choose between embedded FS or local file system
			if tc.useEmbed {
				q.Static("/static", staticFiles)
			} else {
				q.Static("/static", "./static")
			}

			// Define a route for serving files
			q.Get("/", func(c *Ctx) error {
				c.File("static/*") // Must find `static/index.html`
				return nil
			})

			// Creating a test server
			server := httptest.NewServer(q)
			defer server.Close()

			// Making test request
			resp, err := http.Get(server.URL + tc.path)
			if err != nil {
				t.Fatalf("Error making request to %s: %v", tc.path, err)
			}
			defer resp.Body.Close()

			// Check the status code
			if resp.StatusCode != tc.statusCode {
				t.Errorf("Expected status %d, but received %d", tc.statusCode, resp.StatusCode)
			}

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response: %v", err)
			}

			// Checks if the response contains the expected content
			if tc.expectBody != "" && !strings.Contains(string(body), tc.expectBody) {
				t.Errorf("Expected to find '%s' in the response body, but did not find it", tc.expectBody)
			}
		})
	}
}

// ExampleServeStaticIndex demonstrates how to start the Quick server and serve static files correctly.
// The will return func ExampleServeStaticIndex()
//
// Run:
//
//	$ go run main.go
func ExampleQuick_Static() {
	//Quick Start
	q := New()

	/**
	//go:embed static/*
	var staticFiles embed.FS
	*/

	// start FileServer
	// or
	// q.Static("/static", staticFiles)
	q.Static("/static", "./static")

	// send ServeFile
	q.Get("/", func(c *Ctx) error {
		c.File("./static/index.html")
		return nil
	})
}

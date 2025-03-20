package msgid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleNew()
//
//	it with the Examples type.
func ExampleNew() {
	q := quick.New()

	// Apply MsgUUID Middleware globally
	q.Use(New())

	// Define an endpoint that responds with a UUID
	q.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		msgId := c.Request.Header.Get("Msgid")

		// Return 200 OK status
		return c.Status(200).JSON(map[string]string{"msgid": msgId})
	})

	// Send a test request using Quick's built-in test utility
	resp, err := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodGet,
		URI:     "/v1/msguuid/default",
		Headers: map[string]string{"Content-Type": "application/json"},
	})

	// Handle potential errors in test execution
	if err != nil {
		fmt.Println("Test execution error:", err)
		return
	}

	// Print the response body for verification
	fmt.Println("Response:", string(resp.Body()))

	// Out put:
	// Response: "msgid":"f299b00d-875e-4502-966e-22e16767eb13"
}

// This function is named ExampleNew_withCustomConfig()
//
//	it with the Examples type.
func ExampleNew_withCustomConfig() {
	mux := http.NewServeMux()

	// Custom configuration with a different MsgID range
	customConfig := Config{
		Name:  "X-Custom-MsgID",
		Start: 500000000,
		End:   600000000,
	}

	// Apply the MsgID middleware with custom settings
	mux.Handle("/", New(customConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgId := r.Header.Get("X-Custom-MsgID")
		w.Write([]byte("Custom MsgID: " + msgId))
	})))

	// Simulate an HTTP request for testing
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Capture the response body
	responseBody := rec.Body.String()

	// Validate the response dynamically
	fmt.Println(responseBody)

	// Out put:
	// Custom MsgID: 5XXXXXXXX (random number within range)
}

// This function is named ExampleAlgoDefault()
//
//	it with the Examples type.
func ExampleAlgoDefault() {
	// Define the expected range
	start := 900000000
	end := 1000000000

	// Generate a MsgID
	msgID := AlgoDefault(start, end)

	// Check if the MsgID starts with "9"
	if !strings.HasPrefix(msgID, "9") {
		fmt.Println("Test failed: MsgID does not start with 9")
		return
	}

	// Print the generated MsgID for debugging
	fmt.Println("Generated MsgID:", msgID)

	// Out put:
	// Generated MsgID: 9XXXXXXXX (random number in range)
}

// This function is named ExampleNew_withCustomAlgo()
//
//	it with the Examples type.
func ExampleNew_withCustomAlgo() {
	mux := http.NewServeMux()

	// Custom MsgID generator function
	customAlgo := func() string {
		return "custom-msg-12345"
	}

	// Custom configuration using the custom MsgID generator
	customConfig := Config{
		Name: "X-Custom-Trace",
		Algo: customAlgo,
	}

	// Apply the middleware with custom algorithm
	mux.Handle("/", New(customConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgId := r.Header.Get("X-Custom-Trace")
		w.Write([]byte("Generated MsgID: " + msgId))
	})))

	// Simulate an HTTP request for testing
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Print the response for verification
	fmt.Println("Response:", rec.Body.String())

	// Out put:
	// Response: Generated MsgID: custom-msg-12345
}

package main

import (
	"net/http"

	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {
	// Default: Logger in text mode (colored output for terminal)
	loggerMiddleware := logger.New()

	// Logger in JSON mode (structured logs for production)
	jsonLoggerMiddleware := logger.New(logger.Config{Format: "json"})

	// Using middleware in an HTTP server
	mux := http.NewServeMux()
	mux.Handle("/v1/user", loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Quick 💕!"))
		})))

	mux.Handle("/v1/logger", jsonLoggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Quick 💕!"))
		})))

	http.ListenAndServe(":8080", mux)
}

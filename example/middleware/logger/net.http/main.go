package main

import (
	"net/http"

	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {
	loggerMiddleware := logger.New()
	jsonLoggerMiddleware := logger.New(logger.Config{Format: "json"})

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

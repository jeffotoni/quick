package main

import (
	"github.com/jeffotoni/quick/middleware/logger"
	"net/http"
)

func main() {
	loggerMiddleware := logger.New()

	slogLoggerMiddleware := logger.New(logger.Config{
		Format: "slog",
	})

	jsonLoggerMiddleware := logger.New(logger.Config{
		Format: "json",
	})

	CustomLoggerMiddleware := logger.New(logger.Config{
		Format:  "slog",
		Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency}",
	})

	mux := http.NewServeMux()

	mux.Handle("/v1/logger", loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Quick middleware default 💕!"))
		})))

	mux.Handle("/v1/logger/json", jsonLoggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Quick middleware json💕!"))
		})))

	mux.Handle("/v1/logger/slog", slogLoggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Quick middleware slog 💕!"))
		})))

	mux.Handle("/v1/logger/custom", CustomLoggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Quick middleware custom💕!"))
		})))

	http.ListenAndServe(":8080", mux)
}

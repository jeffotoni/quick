package main

import (
	"net/http"

	"github.com/jeffotoni/quick/middleware/cors"
)

func main() {
	c := cors.New(cors.Config{
		AllowedOrigins: []string{"*"},
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"name\": \"quick\"}"))
	})

	wrappedHandler := c(handler)

	http.ListenAndServe(":8080", wrappedHandler)
}

package main

import (
	"net/http"

	"github.com/jeffotoni/quick/middleware/cors"
)

func main() {

	corsMiddleware := cors.New(cors.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"message\": \"User endpoint\"}"))
	})

	wrappedMux := corsMiddleware(mux)

	http.ListenAndServe(":8080", wrappedMux)
}

package cors

import (
	"net/http"

	"github.com/rs/cors"
)

func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
			MaxAge:         86400,
		})
		cors.ServeHTTP(w, r, next.ServeHTTP)
	})
}

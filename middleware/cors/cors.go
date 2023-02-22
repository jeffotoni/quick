package cors

import (
	"net/http"

	"github.com/rs/cors"
)

// func Cors(h http.Handler) http.Handler {
// 	return cors.Default().Handler(h)
// }

func Cors(next http.Handler) http.Handler {
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

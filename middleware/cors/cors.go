package cors

import (
	"net/http"
	"strings"
	// "github.com/rs/cors"
)

// type Cors struct {
// 	handler *cors.Cors
// }

// func start(AllowedOrigins []string) *Cors {
// 	return &Cors{
// 		handler: cors.New(cors.Options{
// 			AllowedOrigins: AllowedOrigins,
// 		}),
// 	}
// }

// func (c *Cors) Handler(next http.Handler) http.Handler {
// 	return c.handler.Handler(next)
// }

type Config struct {
	AllowedOrigins   []string
	AllowCredentials string
}

var ConfigDefault = Config{
	AllowedOrigins:   []string{"*"},
	AllowCredentials: "true",
}

func New(config ...Config) func(http.Handler) http.Handler {
	cfd := ConfigDefault
	if len(config) > 0 {
		cfd = config[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//cor := start(cfd.AllowedOrigins)
			//cor.Handler(next).ServeHTTP(w, r)
			if r.Method != "POST" {
				return
			}

			// r.Header.Set("Access-Control-Allow-Origin", strings.Join(cfd.AllowedOrigins, ","))
			// r.Header.Set("Access-Control-Allow-Credentials", cfd.AllowCredentials)
			// r.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			// r.Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
			// r.Header.Set("Content-Type", "application/json")

			w.Header().Set("Access-Control-Allow-Origin", strings.Join(cfd.AllowedOrigins, ","))
			w.Header().Set("Access-Control-Allow-Credentials", cfd.AllowCredentials)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

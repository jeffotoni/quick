package cors

import (
	"net/http"

	"github.com/rs/cors"
)

type Cors struct {
	handler *cors.Cors
}

func start(AllowedOrigins []string) *Cors {
	return &Cors{
		handler: cors.New(cors.Options{
			AllowedOrigins: AllowedOrigins,
		}),
	}
}

func (c *Cors) Handler(next http.Handler) http.Handler {
	return c.handler.Handler(next)
}

type Config struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

var ConfigDefault = Config{
	AllowedOrigins:   []string{"*"},
	AllowCredentials: true,
}

func New(config ...Config) func(http.Handler) http.Handler {
	cfd := ConfigDefault
	if len(config) > 0 {
		cfd = config[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cor := start(cfd.AllowedOrigins)
			cor.Handler(next).ServeHTTP(w, r)
		})
	}
}

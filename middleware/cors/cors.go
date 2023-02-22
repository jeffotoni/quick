package cors

import (
	"net/http"

	"github.com/rs/cors"
)

type Cors struct {
	handler *cors.Cors
}

func New() *Cors {
	return &Cors{
		handler: cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
		}),
	}
}

func (c *Cors) Handler(next http.Handler) http.Handler {
	return c.handler.Handler(next)
}

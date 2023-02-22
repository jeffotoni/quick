package middleware

import (
	"net/http"

	apitoken "github.com/jeffotoni/quick/middleware/api_token"
)

type Middleware interface {
	Auth(apiTokenKey string, value string)
}

type middleware struct {
	Route *http.HandlerFunc
}

func NewMiddleware(route *http.HandlerFunc) Middleware {
	return &middleware{
		Route: route,
	}
}

func (m *middleware) Auth(apiTokenKey string, value string) {
	*m.Route = apitoken.Auth(apiTokenKey, value)
}

package middleware

import (
	"net/http"

	apitoken "github.com/jeffotoni/quick/middleware/api_token"
)

type Middleware interface {
	Auth(apiTokenKey string, value string)
	Route() http.HandlerFunc
}

type middleware struct {
	route http.HandlerFunc
}

func NewMiddleware(route http.HandlerFunc) Middleware {
	return &middleware{
		route: route,
	}
}

func (m *middleware) Auth(apiTokenKey string, value string) {
	m.route = apitoken.Auth(m.route, apiTokenKey, value)
}

func (m *middleware) Route() http.HandlerFunc {
	return m.route
}

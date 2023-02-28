package cors

import (
	"net/http"

	"github.com/rs/cors"
)

type Config struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool
	// AllowOriginRequestFunc is a custom function to validate the origin. It takes the HTTP Request object and the origin as
	// argument and returns true if allowed or false otherwise. If this option is set, the content of `AllowedOrigins`
	// and `AllowOriginFunc` is ignored.
	AllowOriginRequestFunc func(r *http.Request, origin string) bool
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	// AllowPrivateNetwork indicates whether to accept cross-origin requests over a
	// private network.
	AllowPrivateNetwork bool
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
	// Provides a status code to use for successful OPTIONS requests.
	// Default value is http.StatusNoContent (204).
	OptionsSuccessStatus int
	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool
}

var ConfigDefault = Config{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	AllowCredentials: true,
	Debug:            true,
}

func New(options ...Config) func(next http.Handler) http.Handler {
	cfd := ConfigDefault
	if len(options) > 0 {
		cfd = options[0]
	}
	op := cors.Options{}
	op = cors.Options(cfd)

	c := cors.New(op)
	return c.Handler
}

// func New(config ...Config) func(http.Handler) http.Handler {
// 	cfd := ConfigDefault
// 	if len(config) > 0 {
// 		cfd = config[0]
// 	}
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			op := cors.Options{}
// 			op = cors.Options(cfd)
// 			cors.New(op).Handler(next)
// 			//cors.New(op).Handler(next)
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

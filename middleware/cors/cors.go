package cors

import (
	"net/http"
	"strings"
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
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{
		"POST",
		"GET",
		"PUT",
		"DELETE",
		"PATH",
		"HEAD",
		"OPTIONS",
	},
	AllowCredentials: true,
	AllowedHeaders:   []string{"Origin", "Content-Type"},
	Debug:            false,
	MaxAge:           0,
}

func New(config ...Config) func(http.Handler) http.Handler {
	cfd := ConfigDefault
	if len(config) > 0 {
		cfd = config[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Cors", "true")

			if len(cfd.AllowedOrigins) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", strings.Join(cfd.AllowedOrigins, ", "))
			}
			if len(cfd.AllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfd.AllowedMethods, ", "))
			}
			if len(cfd.AllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfd.AllowedHeaders, ", "))
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

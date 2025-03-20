// Package cors provides middleware for handling Cross-Origin Resource Sharing (CORS)
// in HTTP servers. It allows fine-grained control over cross-origin requests,
// defining which origins, headers, and methods are permitted.
//
// CORS is a security mechanism implemented by browsers to restrict
// how resources on a web page can be requested from another domain.
// This middleware helps manage CORS policies efficiently.
//
// Features:
// - Configurable allowed origins, headers, and methods.
// - Supports credentials (cookies, authentication).
// - Handles CORS preflight (`OPTIONS`) requests automatically.
// - Allows wildcard (`*`) for flexible domain matching.
// - Provides private network access support.
package cors

import (
	"net/http"
	"strconv"
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

// ConfigDefault is the default CORS configuration.
var ConfigDefault = Config{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
	ExposedHeaders:   []string{"Content-Length"},
	AllowCredentials: false,
	MaxAge:           600,
	Debug:            false,
}

// New creates a new CORS middleware with the given configuration.
//
// If no configuration is provided, it defaults to `ConfigDefault`.
//
// Parameters:
//   - config ...Config (optional): Custom configuration. If not provided, `ConfigDefault` is used.
//
// Returns:
//   - A middleware function that wraps an `http.Handler`, ensuring CORS rules are applied.
func New(config ...Config) func(http.Handler) http.Handler {
	c := ConfigDefault
	if len(config) > 0 {
		c = config[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// aply rules CORS
			applyCORSHeaders(c, w, r)

			// If it is an OPTIONS (preflight) request, respond directly
			if r.Method == http.MethodOptions {
				if r.Header.Get("Access-Control-Request-Method") != "" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// applyCORSHeaders applies the necessary CORS headers to an HTTP response.
// This function is called during request handling to ensure cross-origin requests
// are properly validated and allowed based on the provided configuration.
//
// The function dynamically sets `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`,
// `Access-Control-Allow-Headers`, and other necessary headers to comply with CORS policies.
//
// If `AllowCredentials` is `true`, it ensures the correct Origin is set instead of "*",
// since browsers restrict using "*" when credentials are enabled.
//
// Headers Managed by This Function:
// - Access-Control-Allow-Origin: Specifies the allowed origin(s) based on configuration.
// - Access-Control-Allow-Credentials: Determines whether credentials (cookies, HTTP authentication) are allowed.
// - Access-Control-Allow-Methods: Lists allowed HTTP methods.
// - Access-Control-Allow-Headers: Specifies allowed request headers, supporting `*` dynamically.
// - Access-Control-Expose-Headers: Defines headers that are accessible from the frontend.
// - Access-Control-Max-Age: Specifies the cache duration for preflight responses.
//
// Parameters:
// - c: Config – CORS configuration settings.
// - w: http.ResponseWriter – The response writer to send headers.
// - r: *http.Request – The incoming HTTP request.
//
// Example Usage:
// This function is used inside a middleware for handling CORS.
//
//	func New(config Config) func(http.Handler) http.Handler {
//	    return func(next http.Handler) http.Handler {
//	        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	            applyCORSHeaders(config, w, r)
//	            next.ServeHTTP(w, r)
//	        })
//	    }
//	}
func applyCORSHeaders(c Config, w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// If there is no Origin in the request, there's no need to apply CORS.
	if origin == "" {
		return
	}

	// Handle Access-Control-Allow-Origin based on configuration
	if contains(c.AllowedOrigins, "*") && c.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	} else if contains(c.AllowedOrigins, "*") && !c.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else if contains(c.AllowedOrigins, origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if c.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
	}

	// Set Allowed Methods
	if len(c.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ", "))
	}

	// Set Allowed Headers
	if len(c.AllowedHeaders) == 1 && c.AllowedHeaders[0] == "*" {
		reqHeaders := r.Header.Get("Access-Control-Request-Headers")
		if reqHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
		}
	} else if len(c.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowedHeaders, ", "))
	}

	// Set Exposed Headers
	if len(c.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(c.ExposedHeaders, ", "))
	}

	// Set Max-Age if specified
	if c.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(c.MaxAge))
	}
}

// contains checks if a slice contains a specific value.
func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

// Default returns the default CORS configuration.
//
// If a custom configuration is provided, it overrides `ConfigDefault`.
func Default(config ...Config) Config {
	cfd := ConfigDefault
	if len(config) > 0 {
		cfd = config[0]
	}
	return cfd
}

// Handler applies CORS rules to the given `http.Handler`.
//
// This function ensures that CORS policies are applied before passing
// the request to the next handler.
func (c Config) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// aply rules CORS
		applyCORSHeaders(c, w, r)

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			// w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

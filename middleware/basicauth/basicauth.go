// Package basicauth provides HTTP Basic Authentication middleware for Go web servers.
// The middleware implements RFC 7617 (Basic Authentication) to protect HTTP endpoints
// by requiring valid credentials in the Authorization header.

// # Features

//   - Simple integration with net/http handlers
//   - RFC-compliant Basic Authentication
//   - Secure credential validation
//   - WWW-Authenticate header with realm support
//   - Clear unauthorized responses

// Example of how to use middleware in Quick

//	$ curl -H "Authorization: Basic $(echo -n 'wronguser:wrongpass' | base64)" http://localhost:8080/protected
//	$ curl -H "Authorization: Basic $(echo -n 'admin:1234' | base64)" http://localhost:8080/protected
//	$ curl http://localhost:8080/protected
//	$ curl -u admin:1234 http://localhost:8080/protected

// # Response Behavior

//   - Valid credentials: Proceeds to the next handler
//   - Missing header: Returns 401 with WWW-Authenticate header
//   - Invalid format: Returns 401 Unauthorized
//   - Wrong credentials: Returns 401 Unauthorized

package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuth creates middleware that enforces HTTP Basic Authentication.

// Parameters:
//   - username: The required username for authentication
//   - password: The required password for authentication

// Returns:
//   - A middleware function that wraps http.Handler with authentication

// Example:
//   // Protect a handler with basic auth
//   authMiddleware := BasicAuth("admin", "s3cr3t")
//   protectedHandler := authMiddleware(yourHandler)

// Note:
// For more advanced configuration (multiple users, custom responses),
// consider implementing a Config struct as shown in more complete examples.

func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				//fmt.Fprintf(w, "error in logger: %v", err)
				return
			}

			// Check if the header starts with "Basic"
			if !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Decodifica as credenciais base64
			payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Decode base64 credentials
			creds := strings.SplitN(string(payload), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// If authenticated, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

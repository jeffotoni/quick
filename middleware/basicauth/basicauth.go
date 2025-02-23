// The BasicAuth middleware implements HTTP Basic Authentication
// to secure specific routes on an HTTP server. It requires clients
// to send an Authorization header with a Base64-encoded username
// and password to access specific endpoints.
// This middleware follows the RFC 7617 authentication standard,
// allowing secure applications to authenticate users easily,
// without the need for tokens or external authentication systems.
// Example of how to use middleware in Quick
//
//	$ curl -H "Authorization: Basic $(echo -n 'wronguser:wrongpass' | base64)" http://localhost:8080/protected
//	$ curl -H "Authorization: Basic $(echo -n 'admin:1234' | base64)" http://localhost:8080/protected
//	$ curl http://localhost:8080/protected
//	$ curl -u admin:1234 http://localhost:8080/protected
package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuth returns a middleware for basic authentication
// The result will BasicAuth(username, password string) func(http.Handler) http.Handler
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

// basic auth
//
// $ curl -u admin:1234 http://localhost:8080/protected
// $ curl -H "Authorization: Basic YWRtaW46MTIzNA==" http://localhost:8080/protected
package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/jeffotoni/quick"
)

func main() {

	q := quick.New()

	// implementing middleware directly in Use
	q.Use(func(next http.Handler) http.Handler {
		// credentials
		username := "admin"
		password := "1234"

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if it starts with "Basic"
			if !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Decode credentials
			payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			creds := strings.SplitN(string(payload), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))

}

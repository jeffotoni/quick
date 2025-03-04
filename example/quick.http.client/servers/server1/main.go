//go:build !exclude_test

package main

import (
	"encoding/json"
	"net/http"

	"github.com/jeffotoni/quick"
)

// Middleware that processes the form before passing it to the next handler
func postform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensures that the request is of type POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Create a map to store the received data
		formData := make(map[string]string)
		for key, values := range r.Form {
			formData[key] = values[0] // Get the first value for each key
		}

		// Convert to JSON to respond to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Form received successfully",
			"data":    formData,
		})

		// Pass the request to the next handler (if any)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func adaptHandler(h http.Handler) quick.HandleFunc {
	return func(c *quick.Ctx) error {
		h.ServeHTTP(c.Response, c.Request)
		return nil
	}
}

func main() {

	h := func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		if c.Body() == nil {
			return c.Status(200).Send([]byte(`{"data":"quick is awesome!"}`))
		}
		return c.Status(200).Send(c.Body())
	}

	q := quick.New()

	q.Get("/v1/user/:id", h)
	q.Post("/v1/user", h)
	q.Put("/v1/user/:id", h)
	q.Delete("/v1/user/:id", h)

	q.Post("/postform", adaptHandler(postform(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}))))

	q.Listen("0.0.0.0:3000")
}

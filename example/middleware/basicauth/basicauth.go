package main

import (
	"log"
	"net/http"

	"github.com/jeffotoni/goquick"
)

// curl -i -H "Block:true" -XGET localhost:8080/v1/blocked
func main() {

	q := quick.New()

	q.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Este middleware irá bloquear sua requisição se não passar o header Block:true
			if r.Header.Get("Block") == "" || r.Header.Get("Block") == "false" {
				w.WriteHeader(403)
				w.Write([]byte(`{"Message": "Envie block no seu header, por favor! :("}`))
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	q.Get("/v1/blocked", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type my struct {
			Msg   string `json:"msg"`
			Block string `json:"block_message"`
		}

		log.Println(c.Headers["Messageid"])

		return c.Status(200).JSON(&my{
			Msg:   "Quick ❤️",
			Block: c.Headers["Block"][0],
		})
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))

}

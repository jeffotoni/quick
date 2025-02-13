package main

import (
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

// curl -i -H "Block:true" -XGET localhost:8080/v1/blocked
func main() {

	app := quick.New()

	app.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Este middleware, irá bloquear sua requisicao se não passar header Block:true
			if r.Header.Get("Block") == "" || r.Header.Get("Block") == "false" {
				w.WriteHeader(400)
				w.Write([]byte(`{"Message": "Envia block em seu header, por favor! :("}`))
				return
			}

			if r.Header.Get("Block") == "true" {
				w.WriteHeader(200)
				w.Write([]byte(""))
				//h.ServeHTTP(w, r)

				return
			}
			h.ServeHTTP(w, r)
		})
	})

	app.Get("/v1/blocked", func(c *quick.Ctx) error {
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

	log.Fatal(app.Listen("0.0.0.0:8080"))

}

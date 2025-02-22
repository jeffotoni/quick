# CORS

CORS significa **"Cross-Origin Resource Sharing"**, que é uma técnica de segurança usada pelos navegadores da web para permitir que um servidor restrinja o acesso de outros sites ou domínios aos seus recursos. O objetivo principal do CORS é proteger os recursos do servidor de ataques maliciosos de outros domínios.

O Quick é um framework web em Go que suporta o middleware de CORS para lidar com solicitações de outros domínios. O middleware de CORS pode ser adicionado ao Quick usando a biblioteca "github.com/jeffotoni/goquick/middleware/cors".

Para adicionar o middleware de CORS em um aplicativo Quick, basta importar a biblioteca e chamar a função Cors() passando as opções de configuração desejadas.

## cors.nativo

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/goquick"
	"github.com/jeffotoni/goquick/middleware/cors"
)

func main() {
	app := quick.New()

	app.Use(cors.New(cors.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
	}), "cors")

	app.Post("/v1/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		type My struct {
			Name string `json:"name"`
			Year int    `json:"year"`
		}

		var my My
		err := c.BodyParser(&my)
		fmt.Println("byte:", c.Body())

		if err != nil {
			c.Status(400).SendString(err.Error())
			return
		}

		fmt.Println("String:", c.BodyString())
		c.Status(200).JSON(&my)
		return
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}
```
```go
curl --location 'http://localhost:8080/v1/user'
--header 'Content-Type", "application/json' \
--data '
```

## cors.blocked

```go
package main

import (
	"log"
	"net/http"

	"github.com/jeffotoni/goquick"
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
```
```go
curl --location 'http://localhost:8080/v1/blocked'
--header 'Content-Type", "application/json' \
--data '
```

## cors.rs

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/goquick"
	"github.com/rs/cors"
)

func main() {

	app := quick.New()

	app.Post("/v1/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		type My struct {
			Name string `json:"name"`
			Year int    `json:"year"`
		}

		var my My
		err := c.BodyParser(&my)
		fmt.Println("byte:", c.Body())

		if err != nil {
			c.Status(400).SendString(err.Error())
			return
		}

		fmt.Println("String:", c.BodyString())
		c.Status(200).JSON(&my)
		return
	})

	mux := cors.Default().Handler(app)
	log.Fatal(app.Listen("0.0.0.0:8080", mux))
}
```
```go
curl --location 'http://localhost:8080/v1/user'
--header 'Content-Type", "application/json' \
--data '
```

## usando net/http

```go
package main

import (
    "io"
    "net/http"

    "github.com/jeffotoni/goquick/middleware/cors"
)

type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json")
    b, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(400)
        w.Write([]byte(`{"msg":"error"}`))
        return
    }
    w.WriteHeader(200)
    w.Write(b)
}

func OtherHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write([]byte("Outro endpoint!"))
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/v1/user", &MyHandler{})
    mux.HandleFunc("/outro", OtherHandler)

    newmux := cors.Default().Handler(mux)
    println("server: :8080")
    http.ListenAndServe(":8080", newmux)
}
```
```go
curl --location 'http://localhost:8080/v1/user'
--header 'Content-Type", "application/json' \
--data '
```



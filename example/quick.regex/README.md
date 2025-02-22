# Regex

Regex, ou "Expressões Regulares", é uma técnica utilizada em programação para buscar e manipular padrões de texto. O Framework Quick suporta o uso de regex em rotas, permitindo que desenvolvedores criem rotas dinâmicas e flexíveis.

Para usar regex em rotas no Quick, o desenvolvedor precisa definir uma rota usando uma string que contenha um padrão de expressão regular válido. Isso pode ser feito usando o método HTTP apropriado (Get(), Post(), Put(), etc.) no objeto de aplicativo Quick.

```go
package main

import (
	"github.com/jeffotoni/goquick"
	"github.com/jeffotoni/goquick/middleware/msgid"
)

func main() {
	q := quick.New()

	// adicionando middleware msgid
	q.Use(msgid.New())

	q.Get("/v1/user/{id:[0-9]+}", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).String("Quick ação total!!!")
	})

	q.Listen(":8080")
}

```
```go
$ curl --location -g --request GET 'http://localhost:8080/v1/user/{id:[0-9]+}' \
--header 'Content-Type: application/json/' \
--data 'Quick ação total!!!'
```
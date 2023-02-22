# quick - Route Go
<!-- [![GoDoc](https://godoc.org/github.com/jeffotoni/quick?status.svg)](https://godoc.org/github.com/jeffotoni/quick) [![Github Release](https://img.shields.io/github/v/release/jeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/jeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jeffotoni/quick/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jeffotoni/quick/tree/master) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick) [![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/jeffotoni/quick/master) ![Coveralls](https://img.shields.io/coverallsCoverage/github/jeffotoni/quick) -->

O quick é um gerenciador de rotas para Go, onde sua sintaxe foi inspirado no framework fiber.

É um gerenciador de rotas minimalistico está nascendo e está em desenvolvimento, é para ser ráṕido e com alto desempenho 100% compatível com net/http.

O objetivo é didático, e colaboração, todos são bem vindos a ajudar.

O foco será o desempenho, otimizações e muito testes unitários.

Vou deixar um sumário simples do que teremos para desenvolver:

- [50%] Desenvolver Routes Método GET
- [90%] Desenvolver Routes Método GET aceitando Query String
- [90%] Desenvolver Routes Método GET aceitando Parametros 
- [90%] Desenvolver Routes Método GET aceitando Query String e Parametros
- [0.%] Desenvolver Routes Método GET aceitando expressão regular
- [80%] Desenvolver Routes Método POST
- [90%] Desenvolver Routes Método POST aceitando JSON
- [90%] Desenvolver para o MÉTODO POST o parse JSON
- [90%] Desenvolver para o MÉTODO POST funções para acessar byte ou string do Parse
- [0.%] Desenvolver para o MÉTODO PUT
- [0.%] Desenvolver para o MÉTODO DELETE
- [90%] Desenvolver método para ListenAndServe
- [0.%] Desenvolver método para ListenAndServeTLS (http2)
- [70%] Desenvolver método para Facilitar a manipulação do ResponseWriter
- [70%] Desenvolver método para Facilitar a manipulação do Request
- [70%] Desenvolver suporte a ServeHTTP
- [10%] Desenvolver suporte a middlewares
- [80%] Desenvolve suporte a Grupo de Rotas


#### Contribuição 
Jà temos um exemplo, e já podemmos testar e brincar 😁, é claro estamos no inicio ainda tem muito para fechar e fiquem a vontade em fazerem *PR* (com risco de ganhar uma camiseta Go ❤️ e é claro notoriedade para trabalhar com Go 😍 no mercado de trabalho)

##### Quick
```go

package main

import "github.com/jeffotoni/quick"

func main() {
	app := quick.New()

	app.Get("/v1/user", func(c *quick.Ctx) {
		c.Status(200).String("Quick em ação ❤️!")
		return
	})

	app.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" 'localhost:8080/v1/user'

```



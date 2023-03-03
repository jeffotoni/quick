# quick - Route Go
[![GoDoc](https://godoc.org/github.com/gojeffotoni/quick?status.svg)](https://godoc.org/github.com/gojeffotoni/quick) [![Github Release](https://img.shields.io/github/v/release/gojeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/gojeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/gojeffotoni/quick/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/gojeffotoni/quick/tree/master) [![Go Report](https://goreportcard.com/badge/github.com/gojeffotoni/quick)](https://goreportcard.com/badge/github.com/gojeffotoni/quick) [![License](https://img.shields.io/github/license/gojeffotoni/quick)](https://img.shields.io/github/license/gojeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/gojeffotoni/quick/master) ![Coveralls](https://img.shields.io/coverallsCoverage/github/gojeffotoni/quick)

O **Quick** é um gerenciador de rotas para Go bem flexível e extensível com diversas funcionalidades, onde sua sintaxe foi inspirado no **framework fiber**.

É um gerenciador de rotas que encontra em **desenvolvimento** constante, projetado para ser **rápido e de alto desempenho, 100% compatível com net/http**. **O objetivo é didático e colaborativo, todos são bem-vindos a ajudar. 😍** 

Espero que possam participar e que gostem de Godar!!! 😍

## 📃| Rodmap do desenvolvimento

| Tarefa                                          | Progresso |
|-------------------------------------------------|-----------|
| Desenvolver MaxBodySize metodos Post e Put       | 100%     |
| Desenvolver Padrão de Testes Unitários       	   | 90%      |
| Desenvolver Config em New(Config{}) não obrigatório | 100%   |
| Desenvolve suporte a Grupo de Rotas - Group Get e Post | 70% |
| Desenvolver e relacionar ao Listen o Config      | 30%       |
| Criação de função print para não usar fmt de forma demasiada | 100% |
| Criação de função própria para Concat String     | 100%      |
| Criação de benchmarking entre os.Stdout e fmt.Println | 100%   |
| Desenvolver Routes Método GET                    | 50%       |
| Desenvolver Routes Método GET aceitando Query String | 90%    |
| Desenvolver Routes Método GET aceitando Parametros | 90%      |
| Desenvolver Routes Método GET aceitando Query String e Parametros | 90% |
| Desenvolver Routes Método GET aceitando expressão regular | 90.% |
| Desenvolver Routes Método POST                   | 80%       |
| Desenvolver Routes Método POST aceitando JSON    | 90%       |
| Desenvolver para o MÉTODO POST o parse JSON       | 90%       |
| Desenvolver para o MÉTODO POST funções para acessar byte ou string do Parse | 90% |
| Desenvolver para o MÉTODO PUT                    | 80%       |
| Desenvolver para o MÉTODO PUT o parse JSON        | 90%       |
| Desenvolver para o MÉTODO PUT o parse JSON        | 90%       |
| Desenvolver para o MÉTODO PUT funções para acessar byte ou string do Parse | 90% |
| Desenvolver para o MÉTODO DELETE                  | 0.%       |
| Desenvolver para o MÉTODO OPTIONS                 | 0.%       |
| Desenvolver método para ListenAndServe           | 90%       |
| Desenvolver método para ListenAndServeTLS (http2) | 0.%       |
| Desenvolver método para Facilitar a manipulação do ResponseWriter | 70% |
| Desenvolver método para Facilitar a manipulação do Request | 70%  |
| Desenvolver suporte a ServeHTTP                  | 70%       |
| Desenvolver suporte a middlewares                 | 10%       |
| Desenvolve suporte Static Files                   | 0.%       |
| Desenvolver suporte Cors                          | 0.%       |

##### Primeiro exemplo Quick
```go

package main

import "github.com/gojeffotoni/quick"

func main() {
	app := quick.New()

	app.Get("/v1/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		c.Status(200).SendString("Quick em ação ❤️!")
	})

	app.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/user'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 23

Quick em ação ❤️!

```

##### Quick Get Params
```go

package main

import "github.com/gojeffotoni/quick"

func main() {
	app := quick.New()

	app.Get("/v1/customer/:param1/:param2", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")

		type my struct {
			Msg string `json:"msg"`
			Key string `json:"key"`
			Val string `json:"val"`
		}

		c.Status(200).JSON(&my{
			Msg: "Quick ❤️",
			Key: c.Param("param1"),
			Val: c.Param("param2"),
		})
	})

	app.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/customer/val1/val2'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 23

{"msg":"Quick ❤️","key":"val1","val":"val2"}

```

##### Quick Post Body json
```go

package main

import "github.com/gojeffotoni/quick"

type My struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func main() {
	app := quick.New()
	app.Post("/v1/user", func(c *quick.Ctx) {
		var my My
		err := c.Body(&my)
		if err != nil {
			c.Status(400).SendString(err.Error())
			return
		}

		c.Status(200).String(c.BodyString())
		// ou 
		// c.Status(200).JSON(&my)
	})

	app.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v1/user' \
-d '{"name":"jeffotoni", "year":1990}'
HTTP/1.1 200 OK
Date: Wed, 22 Feb 2023 08:10:06 GMT
Content-Length: 32
Content-Type: text/plain; charset=utf-8

{"name":"jeffotoni","year":1990}

```

## ⚙️| Funcionalidades

| Funcionalidades                                 | Possui    |
|-------------------------------------------------|-----------|
| Gerenciador de Rotas                            |   sim     |
| Server Files Static                             |   sim     |
| Grupo de Rotas                                  |   sim     |
| Middlewares                                     |   sim     |
| HTTP/2 support                                  |   sim     |
| Data binding for JSON, XML and form payload     |   sim     |


## 👁‍🗨| Examples


##### Quick Post Bind json
```go

package main

import "github.com/gojeffotoni/quick"

type My struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func main() {
	app := quick.New()
	app.Post("/v2/user", func(c *quick.Ctx) {
		var my My
		err := c.Bind(&my)
		if err != nil {
			c.Status(400).SendString(err.Error())
			return
		}
		c.Status(200).JSON(&my)
	})

	app.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v2/user' \
-d '{"name":"Marcos", "year":1990}'
HTTP/1.1 200 OK
Date: Wed, 22 Feb 2023 08:10:06 GMT
Content-Length: 32
Content-Type: text/plain; charset=utf-8

{"name":"Marcos","year":1990}

```

##### Cors
```go

package main

import "github.com/gojeffotoni/quick"
import "github.com/gojeffotoni/quick/middleware/cors"

func main() {
	app := quick.New()
	app.Use(cors.New().Handler)

	app.Get("/v1/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		c.Status(200).SendString("Quick em ação com Cors❤️!")
	})

	app.Listen("0.0.0.0:8080")
}

```

##### quick.Config{}
```go

package main

import "github.com/gojeffotoni/quick"

func main() {
	app := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024,
	})

	app.Get("/v1/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		c.Status(200).SendString("Quick em ação com Cors❤️!")
	})

	app.Listen("0.0.0.0:8080")
}

```

##### quick.Group()
```go
package main

import "github.com/gojeffotoni/quick"

func main() {
	app := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024,
	})

	v1 := app.Group("/v1")
	v1.Get("/user", func(c *quick.Ctx) {
		c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
		return
	})
	v1.Post("/user", func(c *quick.Ctx) {
		c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
		return
	})

	v2 := app.Group("/v2")
	v2.Get("/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		c.Status(200).SendString("Quick em ação com [GET] /v2/user ❤️!")
	})

	v2.Post("/user", func(c *quick.Ctx) {
		c.Set("Content-Type", "application/json")
		c.Status(200).SendString("Quick em ação com [POST] /v2/user ❤️!")
	})

	app.Listen("0.0.0.0:8080")
}

```

##### Quick Tests
```go

package main

import "github.com/gojeffotoni/quick"

func TestQuickExample(t *testing.T) {

    // Here is a handler function Mock
	testSuccessMockHandler := func(c *Ctx) {
		c.Set("Content-Type", "application/json")
		b, _ := io.ReadAll(c.Request.Body)
		resp := ConcatStr(`"data":`, string(b))
		c.Byte([]byte(resp))
	}

	app := quick.New()
	// Here you can create all routes that you want to test
	app.Post("/v1/user", testSuccessMockHandler)
	app.Post("/v1/user/:p1", testSuccessMockHandler)

	wantOutData := `"data":{"name":"jeff", "age":35}`
	reqBody := []byte(`{"name":"jeff", "age":35}`)
    reqHeaders := map[string]string{"Content-Type": "application/json"}

	data, err := app.QuickTest("POST", "/v1/user", reqHeaders, reqBody)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	s := strings.TrimSpace(data.BodyStr())
	if s != wantOutData {
		t.Errorf("was suppose to return %s and %s come", wantOutData, s)
		return
	}

	t.Logf("\nOutputBodyString -> %v", data.BodyStr())
    t.Logf("\nStatusCode -> %d", data.StatusCode())
    t.Logf("\nOutputBody -> %v", string(data.Body())) // I have converted in this example to string but comes []byte as default
    t.Logf("\nResponse -> %v", data.Response())
}

```


## 🤍| Contribuições

Já temos diversos exemplos, e já podemos testar e brincar 😁. É claro, estamos no início, ainda tem muito para fazer. 
Fiquem à vontade em fazer **PR** (com risco de ganhar uma camiseta Go ❤️ e claro reconhecimento como profissional Go 😍 no mercado de trabalho).


### ☕| Apoiadores

**Quick** é um projeto open source, estamos desenvendo nos tempos livres e é claro nas madrugadas, e você está convidado a particpar e fique a vontade em deixar um café para ajudar e incentivar nosso trabalho 😍 
**<img src="https://github.githubassets.com/images/icons/emoji/unicode/2615.png" height=20 alt="Stargazers over time"> [pode comprar um café aqui](https://www.buymeacoffee.com/ZlD8rtr1cq)**

| Avatar | User | Donation |
|--------|------|----------|
| <img src="https://avatars.githubusercontent.com/u/1092879?s=96&v=4" height=20> | @jeffotoni | ☕ x 10 |
| <img src="https://avatars.githubusercontent.com/u/7107528?v=4" height=20> | @crow34 | ☕ x 5  |
| <img src="https://avatars.githubusercontent.com/u/70351793?v=4" height=20> | @Guilherme-De-Marchi | ☕ x 5 |

## Repositório de Exemplos do Framework Quick

Este repositório contém exemplos práticos do Framework Quick, um framework web rápido e leve, desenvolvido em Go.

### Como Usar

Os exemplos estão organizados em pastas separadas, cada uma contendo um exemplo completo de uso do framework Quick em uma aplicação web simples. 

### Contribuindo

Se você tem algum exemplo interessante de uso do Framework Quick, sinta-se à vontade para enviar uma solicitação de pull request com sua contribuição.

### Licença

Este repositório está sob a licença MIT. Consulte o arquivo LICENSE para mais informações.

### Contato

Se você tiver alguma dúvida ou sugestão sobre o uso do Framework Quick ou sobre este repositório de exemplos, sinta-se à vontade para entrar em contato com o desenvolvedor do Framework Quick, Jeff Otoni, através do email gojeff@protonmail.com.

### Link para o Repositório

O repositório de exemplos do Framework Quick pode ser encontrado em https://github.com/gojeffotoni/examples.



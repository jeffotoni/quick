# Group

O Group é uma funcionalidade do Framework Quick que permite agrupar rotas e aplicar middleware a elas.

Por exemplo, se você tiver um conjunto de rotas que precisam de autenticação antes de serem acessadas, em vez de adicionar o middleware de autenticação individualmente para cada rota, pode agrupá-las usando a funcionalidade Group e aplicar o middleware a todas as rotas do grupo de uma só vez. Isso pode tornar o código mais legível e organizado, além de evitar a repetição de código.

#### Group01

```go
package main

import "github.com/jeffotoni/goquick"

func main() {
	q := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024,
	})

	group := q.Group("/v1")
	group.Get("/user", func(c *quick.Ctx) error {
		return c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
	})

	group.Post("/user", func(c *quick.Ctx) error {
		return c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
	})

 q.Listen("0.0.0.0:8080")
}
```
```go
$ curl --location --request GET 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json/' \
--data '[GET] [GROUP] /v1/user ok!!!'
```

#### Group02

```go
package main

import "github.com/jeffotoni/goquick"

func main() {
	q := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024,
	})

	group2 := q.Group("/v2")

	group2.Get("/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("GoQuick em ação com [GET] /v2/user ❤️!")
	})

	group2.Post("/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("GoQuick em ação com [POST] /v2/user ❤️!")
	})


 q.Listen("0.0.0.0:8080")
}
```
```go
$ curl --location --request GET 'http://localhost:8080/v2/user' \
--header 'Content-Type: application/json/' \
--data 'GoQuick em ação com [POST] /v2/user ❤️!'
```




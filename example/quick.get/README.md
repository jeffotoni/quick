# GET

O método get pode ser usado para gerar valores aleatórios de diferentes tipos de dados, como inteiros, strings, booleanos, slices, entre outros. Por exemplo, se você estiver testando uma função que recebe um número inteiro como argumento, você pode usar o método get para gerar números aleatórios e verificar se a função está funcionando corretamente para diferentes valores de entrada.

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Define a rota HTTP GET "/greet/:name"
	q.Get("/greet/:name", func(c *quick.Ctx) error {
		name := c.Param("name")
		c.Set("Content-Type", "text/plain")
		return c.Status(200).SendString("Olá " + name + "!")
	})

	q.Listen("0.0.0.0:8080")
}
```
```go
$ curl --location --request GET 'http://localhost:8080/greet/:name' \
--header 'Content-Type: application/json/' \
--data 'Olá !'
```

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Get("/v2/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Opa, funcionando!")
	})

	q.Listen("0.0.0.0:8080")
}
```
```go
$ curl --location --request GET 'http://localhost:8080/v2/user' \
--header 'Content-Type: application/json/' \
--data 'Opa, funcionando!'
```

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Get("/v1/userx/:p1/:p2/cust/:p3/:p4", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick in action ❤️!")
	})

	q.Listen("0.0.0.0:8080")
}
```
```go
curl --location --request GET 'http://localhost:8080/v1/userx/:p1/:p2/cust/:p3/:p4' \
--header 'Content-Type: application/json/' \
--data 'Quick in action ❤️!'
```

```go
package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Get("/hello/:name", func(c *quick.Ctx) error {
		name := c.Param("name")
		message := "Olá, " + name + "!"
		return c.Status(200).SendString(message)
	})

	q.Listen("0.0.0.0:8080")
}
```
```go
curl --location --request GET 'http://localhost:8080/hello/:name' \
--header 'Content-Type: application/json/' \
--data 'Olá, !'
```
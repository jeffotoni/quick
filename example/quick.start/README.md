# Start

"Start" é um termo geralmente usado em aplicativos web para indicar que o servidor está sendo iniciado e está pronto para lidar com solicitações HTTP. No contexto do Framework Quick, "start" geralmente se refere à ação de iniciar o servidor web Quick.

Para iniciar o servidor web Quick, é preciso criar uma instância do aplicativo Quick, definir as rotas do servidor e chamar a função Listen() passando a porta em que o servidor deve ser executado.

```go
package main

import "github.com/jeffotoni/goquick"

func main() {
	q := quick.New()
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("GoQuick em ação com Cors❤️!")
	})
	q.Listen("0.0.0.0:8080")
}

```

```go
$ curl --location --request GET 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json/' \
--data 'GoQuick em ação com Cors❤️!'
```
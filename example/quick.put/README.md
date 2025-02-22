# Put

O método Put() no Framework Quick é usado para definir uma rota HTTP que responde a solicitações PUT. O método Put() é um método HTTP que permite que os clientes enviem dados para um servidor para atualização de recursos. Em termos simples, ele permite que um cliente envie uma atualização para um recurso existente em um servidor.

Para definir uma rota PUT usando o Quick, basta chamar o método Put() em uma instância de aplicativo Quick. O método Put() aceita dois argumentos: o primeiro é o caminho da rota, e o segundo é o manipulador de rota, que é uma função que é chamada sempre que a rota é correspondida.

```go
package main

import (
	"github.com/jeffotoni/goquick"
)

func main() {
	q := quick.New()

	q.Put("/users/:id", func(c *quick.Ctx) error {
		userID := c.Param("id")
		// Lógica de atualização do usuário
		return c.Status(200).SendString("Usuário " + userID + " atualizado com sucesso!")
	})

	q.Listen(":8080")
}
```
```go
$ curl --location --request PUT 'http://localhost:8080/users/:id' \
--header 'Content-Type: application/json/' \
--data 'Usuário  atualizado com sucesso!'
```


package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msgid"
)

func main() {
	q := quick.New()

	// Aplica o Middleware MsgID globalmente
	q.Use(msgid.New())

	// Define uma rota que retorna o MsgID gerado
	q.Get("/v1/msgid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Obtém o MsgID do header da requisição
		msgId := c.Request.Header.Get("Msgid")

		// Log para depuração
		fmt.Printf("Generated MsgID: %s\n", msgId)

		// Retorna o MsgID no JSON da resposta
		return c.Status(200).JSON(map[string]string{"msgid": msgId})
	})

	// Inicia o servidor
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

// $ curl -i -X GET http://localhost:8080/v1/msgid/default

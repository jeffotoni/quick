package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New()

	q.Put("/users/:id", func(c *quick.Ctx) error {
		userID := c.Param("id")
		// Lógica de atualização do usuário
		return c.Status(200).SendString("Usuário " + userID + " atualizado com sucesso!")
	})

	q.Put("/tipos/:id", func(c *quick.Ctx) error {
		tiposID := c.Param("id")
		// Lógica de atualização do usuário
		return c.Status(200).SendString("Usuário " + tiposID + " tipo atualizado com sucesso!")
	})

	q.Listen(":8080")
}

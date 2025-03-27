package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	q.Use(logger.New(logger.Config{
		Format:  "text",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG",
		CustomFields: map[string]string{
			"user_id": "usr-001",
			"trace":   "trace-debug",
		},
	}))

	q.Use(logger.New(logger.Config{
		Format:  "text",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "INFO",
		CustomFields: map[string]string{
			"user_id": "usr-002",
			"trace":   "trace-info",
		},
	}))

	q.Use(logger.New(logger.Config{
		Format:  "text",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "WARN",
		CustomFields: map[string]string{
			"user_id": "usr-003",
			"trace":   "trace-warn",
		},
	}))

	// Definir rota GET para gerar logs
	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Retornar resposta JSON
		return c.Status(200).JSON(quick.M{
			"msg": "Quick",
		})
	})

	// Iniciar o servidor na porta 8080
	q.Listen("0.0.0.0:8080")
}

// $ curl -i -XGET localhost:8080/v1/logger

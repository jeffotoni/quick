package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

// curl -X POST http://localhost:3000/postform \
//      -H "Content-Type: application/x-www-form-urlencoded" \
//      -d "username=quick_user&password=supersecret"

func main() {
	q := quick.New()

	// Define a route to process POST form-data
	q.Post("/postform", func(c *quick.Ctx) error {
		// Get form values
		form := c.FormValues()

		// Log received form data for debugging
		log.Println("Received form data:", form)

		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Start the server BEFORE making the request
	fmt.Println("Quick server running at http://localhost:3000")
	if err := q.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start Quick server: %v", err)
	}

	// Criando um cliente HTTP antes de chamar PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Define um timeout de 5s
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // Tipo correto para forms
		}),
	)

	// Verifica se o cliente foi inicializado corretamente
	if cClient == nil {
		log.Fatal("Erro: cliente HTTP não foi inicializado corretamente")
	}

	// Declara os valores do formulário
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Envia uma requisição POST com form-data
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Verifica se a resposta é válida
	if resp == nil || resp.Body == nil {
		log.Fatal("Erro: resposta vazia ou inválida")
	}

	// Decodifica a resposta JSON
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	// Exibe a resposta no terminal
	fmt.Println("POST response:", result)
}

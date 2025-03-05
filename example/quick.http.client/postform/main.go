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

func main() {
	q := quick.New()

	// Define a route to process POST form-data
	q.Post("/postform", func(c *quick.Ctx) error {
		form := c.FormValues()
		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Quick server running at http://localhost:3000")
		if err := q.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Quick server: %v", err)
		}
	}()

	// Criando um cliente HTTP antes de chamar PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Define um timeout de 5s
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // Tipo correto para forms
		}),
	)

	// Check if the HTTP client was initialized correctly
	if cClient == nil {
		log.Fatal("Erro: cliente HTTP não foi inicializado corretamente")
	}

	// Declare Values
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Send a POST request
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Check if the response is valid
	if resp == nil || resp.Body == nil {
		log.Fatal("Erro: resposta vazia ou inválida")
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}

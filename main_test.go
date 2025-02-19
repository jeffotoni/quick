package quick

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ExampleGet demonstra como registrar uma rota GET no Quick.
//
// Este exemplo cria uma nova instância do Quick, adiciona uma rota `GET /`
// e retorna uma resposta JSON.
func ExampleGet() {
	q := New()

	q.Get("/", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick em ação com Cors❤️!")
	})

	fmt.Println("Servidor rodando...")
	fmt.Println("Quick em ação com Cors❤️!")

	// Output:
	// Servidor rodando...
	// Quick em ação com Cors❤️!
}

// TestExampleGet verifica se a resposta da rota GET está correta.
// go test -v -count=1 -cover -failfast -run ^TestExampleGet
func TestExampleGet(t *testing.T) {
	q := New()

	q.Get("/", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick em ação com Cors❤️!")
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	q.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 200, mas obteve %d", resp.StatusCode)
	}

	expectedBody := "Quick em ação com Cors❤️!"
	if string(body) != expectedBody {
		t.Errorf("Esperado '%s', mas obteve '%s'", expectedBody, string(body))
	}
}

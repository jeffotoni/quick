package quick

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuickRegistersDistinctDynamicRoutes(t *testing.T) {
	q := New()
	q.Get("/wallets/:walletId", func(c *Ctx) error {
		return c.String("wallet:" + c.Param("walletId"))
	})
	q.Get("/wallets/:walletId/ledger", func(c *Ctx) error {
		return c.String("ledger:" + c.Param("walletId"))
	})
	q.Post("/wallets/:walletId/reconciliation", func(c *Ctx) error {
		return c.String("reconciliation:" + c.Param("walletId"))
	})

	if got := q.HandlersCount(); got != 3 {
		t.Fatalf("expected 3 routes, got %d", got)
	}

	tests := []struct {
		method string
		path   string
		want   string
	}{
		{method: http.MethodGet, path: "/wallets/w-1", want: "wallet:w-1"},
		{method: http.MethodGet, path: "/wallets/w-1/ledger", want: "ledger:w-1"},
		{method: http.MethodPost, path: "/wallets/w-1/reconciliation", want: "reconciliation:w-1"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			response := httptest.NewRecorder()

			q.ServeHTTP(response, req)

			if response.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", response.Code)
			}
			if response.Body.String() != tt.want {
				t.Fatalf("expected body %q, got %q", tt.want, response.Body.String())
			}
		})
	}
}

func TestQuickRejectsExactDuplicateRoute(t *testing.T) {
	q := New()
	handler := func(c *Ctx) error { return c.String("ok") }

	q.Get("/wallets/:walletId", handler)
	q.Get("/wallets/:walletId", handler)

	if got := q.HandlersCount(); got != 1 {
		t.Fatalf("expected 1 route, got %d", got)
	}
}

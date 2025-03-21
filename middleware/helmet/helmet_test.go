package helmet

import (
	"net/http"
	"testing"

	"github.com/jeffotoni/quick"
)

// TestHelmet verifies that the Helmet middleware adds all security headers by default.
func TestHelmet(t *testing.T) {
	q := quick.New()
	q.Use(Helmet())

	q.Get("/v1/health", func(c *quick.Ctx) error {
		return c.Status(http.StatusOK).String("OK")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/health",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := resp.AssertStatus(http.StatusOK); err != nil {
		t.Error(err)
	}

	// Assert default security headers
	resp.AssertHeader("X-XSS-Protection", "0")
	resp.AssertHeader("X-Content-Type-Options", "nosniff")
	resp.AssertHeader("X-Frame-Options", "SAMEORIGIN")
	resp.AssertHeader("Content-Security-Policy", "default-src 'self'")
	resp.AssertHeader("Referrer-Policy", "no-referrer")
	resp.AssertHeader("Permissions-Policy", "")
	resp.AssertHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	resp.AssertHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	resp.AssertHeader("Cross-Origin-Opener-Policy", "same-origin")
	resp.AssertHeader("Cross-Origin-Embedder-Policy", "require-corp")
	resp.AssertHeader("Cross-Origin-Resource-Policy", "same-origin")
	resp.AssertHeader("Origin-Agent-Cluster", "?1")
	resp.AssertHeader("X-DNS-Prefetch-Control", "off")
	resp.AssertHeader("X-Download-Options", "noopen")
	resp.AssertHeader("X-Permitted-Cross-Domain-Policies", "none")
}

// TestHelmetWithoutMiddleware confirms that no security headers are set when Helmet is not used.
func TestWithoutHelmet(t *testing.T) {
	q := quick.New()

	q.Get("/v1/health", func(c *quick.Ctx) error {
		return c.Status(http.StatusOK).String("OK")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/health",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := resp.AssertStatus(http.StatusOK); err != nil {
		t.Error(err)
	}

	// Headers should not be present
	resp.AssertNoHeader("X-XSS-Protection")
	resp.AssertNoHeader("Content-Security-Policy")
	resp.AssertNoHeader("Strict-Transport-Security")
}

// TestHelmetWithCSPReportOnly checks if CSP is correctly set as report-only.
func TestHelmetWithCSPReportOnly(t *testing.T) {
	q := quick.New()
	q.Use(Helmet(Options{
		ContentSecurityPolicy: "default-src 'self'",
		CSPReportOnly:         true,
	}))

	q.Get("/", func(c *quick.Ctx) error {
		return c.String("ok")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/",
	})
	if err != nil {
		t.Error(err)
	}

	resp.AssertHeader("Content-Security-Policy-Report-Only", "default-src 'self'")
	resp.AssertNoHeader("Content-Security-Policy")
}

// TestHelmetWithCustomHSTS ensures custom HSTS settings are applied.
func TestHelmetWithCustomHSTS(t *testing.T) {
	q := quick.New()
	q.Use(Helmet(Options{
		HSTSMaxAge:            86400,
		HSTSExcludeSubdomains: true,
		HSTSPreloadEnabled:    false,
	}))

	q.Get("/", func(c *quick.Ctx) error {
		return c.String("ok")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/",
		// TLS:    true, // Simulate HTTPS
	})
	if err != nil {
		t.Error(err)
	}

	resp.AssertHeader("Strict-Transport-Security", "max-age=86400")
}

// TestHelmetWithNextFunc ensures middleware is skipped when Next returns true.
func TestHelmetWithNextFunc(t *testing.T) {
	q := quick.New()
	q.Use(Helmet(Options{
		Next: func(c *quick.Ctx) bool {
			return c.Path() == "/skip"
		},
	}))

	q.Get("/skip", func(c *quick.Ctx) error {
		return c.String("skipped")
	})

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/skip",
	})
	if err != nil {
		t.Error(err)
	}

	resp.AssertNoHeader("Content-Security-Policy")
	resp.AssertString("skipped")
}

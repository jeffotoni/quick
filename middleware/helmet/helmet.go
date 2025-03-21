// Package helmet provides middleware for the Quick framework that sets various
// HTTP headers to help secure your application.
//
// Inspired by Helmet in the Node.js ecosystem, this package includes protections
// against well-known web vulnerabilities by configuring headers such as:
//
//   - X-XSS-Protection
//   - X-Content-Type-Options
//   - X-Frame-Options
//   - Content-Security-Policy
//   - Referrer-Policy
//   - Permissions-Policy
//   - Cross-Origin-Embedder-Policy
//   - Cross-Origin-Opener-Policy
//   - Cross-Origin-Resource-Policy
//   - Origin-Agent-Cluster
//   - X-DNS-Prefetch-Control
//   - X-Download-Options
//   - X-Permitted-Cross-Domain-Policies
//   - Strict-Transport-Security
//   - Cache-Control
//
// It provides secure defaults, but allows customization via the Options struct.
// You can skip the middleware for specific requests by providing a Next function.
package helmet

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// Options defines the configuration for the Helmet middleware.
// All fields map to specific HTTP headers that enhance security.
// These options can override the default behavior provided by the middleware.
type Options struct {
	// Next defines a function to skip the middleware
	Next func(c *quick.Ctx) bool

	// XSSProtection sets the X-XSS-Protection header
	XSSProtection string

	// ContentTypeNosniff sets the X-Content-Type-Options header
	ContentTypeNosniff string

	// XFrameOptions sets the X-Frame-Options header
	XFrameOptions string

	// ContentSecurityPolicy sets the Content-Security-Policy header
	ContentSecurityPolicy string

	// CSPReportOnly determines whether to use Content-Security-Policy-Report-Only (CSP)
	CSPReportOnly bool

	// ReferrerPolicy sets the Referrer-Policy header
	ReferrerPolicy string

	// PermissionsPolicy sets the Permissions-Policy header
	PermissionsPolicy string

	// Cross-Origin headers
	CrossOriginEmbedderPolicy string
	CrossOriginOpenerPolicy   string
	CrossOriginResourcePolicy string

	// Origin-Agent-Cluster header
	OriginAgentCluster string

	// X-DNS-Prefetch-Control header
	XDNSPrefetchControl string

	// X-Download-Options header
	XDownloadOptions string

	// X-Permitted-Cross-Domain-Policies header
	XPermittedCrossDomain string

	// HSTSMaxAge defines Strict-Transport-Security max-age (HSTS)
	HSTSMaxAge int

	// HSTSExcludeSubdomains omits includeSubDomains
	HSTSExcludeSubdomains bool

	// HSTSPreloadEnabled adds preload directive
	HSTSPreloadEnabled bool

	// CacheControl sets Cache-Control header
	CacheControl string
}

// Helmet returns a Quick-compatible middleware that adds security-related HTTP headers.
//
// Usage:
//
//	q.Use(helmet.Helmet(helmet.Options{
//	    XSSProtection:         "1; mode=block",
//	    ContentTypeNosniff:    "nosniff",
//	    XFrameOptions:         "DENY",
//	    ContentSecurityPolicy: "default-src 'self';",
//	    HSTSMaxAge:            63072000,
//	    HSTSPreloadEnabled:    true,
//	}))
//
// The middleware adds the following headers (depending on the configuration):
//   - X-XSS-Protection
//   - X-Content-Type-Options
//   - X-Frame-Options
//   - Content-Security-Policy / Content-Security-Policy-Report-Only
//   - Referrer-Policy
//   - Permissions-Policy
//   - Cross-Origin-Embedder-Policy
//   - Cross-Origin-Opener-Policy
//   - Cross-Origin-Resource-Policy
//   - Origin-Agent-Cluster
//   - X-DNS-Prefetch-Control
//   - X-Download-Options
//   - X-Permitted-Cross-Domain-Policies
//   - Strict-Transport-Security (only for HTTPS requests)
//   - Cache-Control
//
// You can override default values by passing a custom Options struct.
// Usage Example:
//
//	func secureApp() {
//	    q := quick.New()
//	    q.Use(helmet.Helmet()) // Use with defaults
//
//	    // Or customize:
//	    q.Use(helmet.Helmet(helmet.Options{
//	        XFrameOptions: "DENY",
//	        HSTSMaxAge:    31536000,
//	    }))
//
//	    q.Get("/", func(c *quick.Ctx) error {
//	        return c.SendString("Hello, secure world!")
//	    })
//	}
//
// If the Next function is defined and returns true, the middleware is skipped.
func Helmet(opt ...Options) func(next quick.Handler) quick.Handler {
	return func(next quick.Handler) quick.Handler {
		// Apply default options
		options := defaultOptions()
		if len(opt) > 0 {
			options = opt[0]
		}

		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// Skip middleware if Next function returns true
			if options.Next != nil && options.Next(c) {
				return next.ServeQuick(c)
			}

			// X-XSS-Protection
			if options.XSSProtection != "" {
				c.Set("X-XSS-Protection", options.XSSProtection)
			}

			// X-Content-Type-Options
			if options.ContentTypeNosniff != "" {
				c.Set("X-Content-Type-Options", options.ContentTypeNosniff)
			}

			// X-Frame-Options
			if options.XFrameOptions != "" {
				c.Set("X-Frame-Options", options.XFrameOptions)
			}

			// Content-Security-Policy
			if options.ContentSecurityPolicy != "" {
				if options.CSPReportOnly {
					c.Set("Content-Security-Policy-Report-Only", options.ContentSecurityPolicy)
				} else {
					c.Set("Content-Security-Policy", options.ContentSecurityPolicy)
				}
			}

			// Referrer-Policy
			if options.ReferrerPolicy != "" {
				c.Set("Referrer-Policy", options.ReferrerPolicy)
			}

			// Permissions-Policy
			if options.PermissionsPolicy != "" {
				c.Set("Permissions-Policy", options.PermissionsPolicy)
			}

			// Cross-Origin headers
			setIfNotEmpty(c, "Cross-Origin-Embedder-Policy", options.CrossOriginEmbedderPolicy)
			setIfNotEmpty(c, "Cross-Origin-Opener-Policy", options.CrossOriginOpenerPolicy)
			setIfNotEmpty(c, "Cross-Origin-Resource-Policy", options.CrossOriginResourcePolicy)

			// Origin-Agent-Cluster
			setIfNotEmpty(c, "Origin-Agent-Cluster", options.OriginAgentCluster)

			// DNS Prefetch Control
			setIfNotEmpty(c, "X-DNS-Prefetch-Control", options.XDNSPrefetchControl)

			// Download Options
			setIfNotEmpty(c, "X-Download-Options", options.XDownloadOptions)

			// Permitted Cross Domain Policies
			setIfNotEmpty(c, "X-Permitted-Cross-Domain-Policies", options.XPermittedCrossDomain)

			// Strict-Transport-Security
			if c.Request.TLS != nil && options.HSTSMaxAge > 0 {
				hsts := fmt.Sprintf("max-age=%d", options.HSTSMaxAge)
				if !options.HSTSExcludeSubdomains {
					hsts += "; includeSubDomains"
				}
				if options.HSTSPreloadEnabled {
					hsts += "; preload"
				}
				c.Set("Strict-Transport-Security", hsts)
			}

			// Cache-Control
			if options.CacheControl != "" {
				c.Set("Cache-Control", options.CacheControl)
			}

			return next.ServeQuick(c)
		})
	}
}

// defaultOptions returns a set of secure default values for the Helmet middleware.
// These defaults aim to provide sensible protection out-of-the-box.
func defaultOptions() Options {
	return Options{
		XSSProtection:             "0",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		ContentSecurityPolicy:     "default-src 'self'",
		CSPReportOnly:             false,
		ReferrerPolicy:            "no-referrer",
		PermissionsPolicy:         "",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		OriginAgentCluster:        "?1",
		XDNSPrefetchControl:       "off",
		XDownloadOptions:          "noopen",
		XPermittedCrossDomain:     "none",
		HSTSMaxAge:                31536000,
		HSTSPreloadEnabled:        true,
		CacheControl:              "no-cache, no-store, must-revalidate",
	}
}

// setIfNotEmpty sets a response header only if the provided value is not empty.
// It's used internally to avoid setting headers with blank values.
func setIfNotEmpty(c *quick.Ctx, key, value string) {
	if value != "" {
		c.Set(key, value)
	}
}

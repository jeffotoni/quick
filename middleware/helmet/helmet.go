package helmet

import (
	"net/http"
)

// Options defines the configuration options for the helmet middleware.
type Options struct {
	PoweredBy      bool   // Remove X-Powered-By and Server headers
	CSP            string // Content Security Policy
	ReferrerPolicy string // Referrer Policy
	XFrameOptions  string // X-Frame-Options
	MIMESniffing   bool   // MIME Sniffing protection
	STSetting      string // Strict Transport Security
	CacheControl   string // Cache-Control
}

// Helmet secures HTTP responses by adding security headers.
func Helmet(option ...Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Default options
			opt := defaultOptions()
			if len(option) > 0 {
				opt = option[0]
			}

			// Remove sensitive headers
			if opt.PoweredBy {
				deleteHeader(w, "X-Powered-By")
				deleteHeader(w, "Server")
			}

			// MIME Sniffing protection
			if opt.MIMESniffing {
				addHeader(w, "X-Content-Type-Options", "nosniff")
			}

			// X-Frame-Options to prevent clickjacking
			if opt.XFrameOptions != "" {
				addHeader(w, "X-Frame-Options", opt.XFrameOptions)
			}

			// Referrer Policy
			if opt.ReferrerPolicy != "" {
				addHeader(w, "Referrer-Policy", opt.ReferrerPolicy)
			}

			// Content Security Policy
			if opt.CSP != "" {
				addHeader(w, "Content-Security-Policy", opt.CSP)
			}

			// Strict Transport Security
			if opt.STSetting != "" {
				addHeader(w, "Strict-Transport-Security", opt.STSetting)
			}

			// Cache-Control
			if opt.CacheControl != "" {
				addHeader(w, "Cache-Control", opt.CacheControl)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// defaultOptions returns the default security header settings.
func defaultOptions() Options {
	return Options{
		PoweredBy:      true,
		MIMESniffing:   true,
		XFrameOptions:  "DENY",
		ReferrerPolicy: "strict-origin-when-cross-origin",
		CSP:            "default-src 'self'",
		STSetting:      "max-age=31536000; includeSubDomains; preload",
		CacheControl:   "no-cache, no-store, must-revalidate",
	}
}

func deleteHeader(w http.ResponseWriter, header string) {
	if header == "" {
		return
	}
	w.Header().Del(header)
}

// addHeader adds a header to the response.
func addHeader(w http.ResponseWriter, header, value string) {
	if value == "" {
		return
	}
	w.Header().Set(header, value)
}

package helmet

import (
	"net/http"
)

// Options defines the configuration options for the helmet middleware.
type Options struct {
	PoweredBy      bool // Enable or disable X-Powered-By and Server headers
	CSP            bool // Content Security Policy
	ReferrerPolicy bool // Referrer Policy
	XFrameOptions  bool // X-Frame-Options
	MIMESniffing   bool // MIME Sniffing
	XSSProtection  bool // X-XSS-Protection
	STSetting      bool // Strict Transport Security
	CacheControl   bool // Cache-Control
}

// Helmet helps to secure HTTP responses by adding security headers.
//
// /Parameters:
//   - options: Optional slice of `Options` containing custom security header settings.
//
// /Return:
//   - func(http.Handler) http.Handler: A middleware function that applies security headers to the response.
func Helmet(option ...Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Default options
			opt := defaultOptions()
			if len(option) > 0 {
				opt = option[0]
			}

			// PoweredBy removes X-Powered-By and Server headers
			if opt.PoweredBy {
				deleteHeader(w, "X-Powered-By")
				deleteHeader(w, "Server")
			}

			// MIME Sniffing used to protect against MIME sniffing vulnerabilities.
			if opt.MIMESniffing {
				addHeader(w, "X-Content-Type-Options", "nosniff")
			}

			// X-XSS-Protection used to prevent cross-site scripting attacks
			if opt.XSSProtection {
				addHeader(w, "X-XSS-Protection", "1; mode=block")
			}

			// X-Frame-Options used to prevent clickjacking attacks
			if opt.XFrameOptions {
				addHeader(w, "X-Frame-Options", "DENY")
			}

			// Referrer Policy used to prevent referrer information disclosure
			if opt.ReferrerPolicy {
				addHeader(w, "Referrer-Policy", "strict-origin-when-cross-origin")
			}

			// Content Security Policy instructs to load only those contents that are mentioned in the policy
			if opt.CSP {
				addHeader(w, "Content-Security-Policy", "default-src 'self'")
			}

			// Strict Transport Security force HTTPS to prevent man-in-the-middle attacks
			if opt.STSetting {
				addHeader(w, "Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}

			// Cache-Control remove cache to prevent caching
			if opt.CacheControl {
				addHeader(w, "Cache-Control", "no-cache, no-store, must-revalidate")
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
		XFrameOptions:  true,
		XSSProtection:  false,
		ReferrerPolicy: true,
		CSP:            true,
		STSetting:      true,
		CacheControl:   true,
	}
}

func deleteHeader(w http.ResponseWriter, header string) {
	if header == "" {
		return
	}
	w.Header().Del(header)
}

// addHeader adds a header to the response.
// If the header already exists, it will be overwritten.
func addHeader(w http.ResponseWriter, header string, value string) {
	if value == "" {
		return
	}
	w.Header().Set(header, value)
}

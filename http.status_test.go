package quick

import "testing"

// TestQuick_HttpStatus verifies that the HTTP statuses and HTTP methods are registered correctly.
// This test only prints the values into the log for visual inspection.
// The will test func TestQuick_HttpStatus(t *testing.T)
//
//	$ go test -v -run ^func TestQuick_HttpStatus(t *testing.T)
func TestQuick_HttpStatus(t *testing.T) {

	// Log of informative status codes
	t.Log(StatusContinue)
	t.Log(StatusSwitchingProtocols)
	t.Log(StatusProcessing)
	t.Log(StatusEarlyHints)

	// Log of success status codes
	t.Log(StatusOK)
	t.Log(StatusCreated)
	t.Log(StatusAccepted)
	t.Log(StatusNonAuthoritativeInfo)
	t.Log(StatusNoContent)
	t.Log(StatusResetContent)
	t.Log(StatusPartialContent)

	// Log of supported HTTP methods
	t.Log(MethodGet)
	t.Log(MethodHead)
	t.Log(MethodPost)
	t.Log(MethodPut)
	t.Log(MethodPatch)
	t.Log(MethodDelete)
	t.Log(MethodConnect)
	t.Log(MethodTrace)
}

// TestStatusText checks if the StatusText function returns the correct description
// for each HTTP status code.
// The will test func TestStatusText(t *testing.T)
//
//	$ go test -v -run ^func TestStatusText(t *testing.T)
func TestStatusText(t *testing.T) {
	// Set of tests with HTTP status codes and their expected descriptions.
	tests := []struct {
		code     int    // HTTP status code
		expected string // Expected description
	}{
		// Tests for 1xx status codes (Informative)
		{StatusContinue, "Continue"},
		{StatusSwitchingProtocols, "Switching Protocols"},
		{StatusProcessing, "Processing"},
		{StatusEarlyHints, "Early Hints"},

		// Tests for 2xx status codes (Success)
		{StatusOK, "OK"},
		{StatusCreated, "Created"},
		{StatusAccepted, "Accepted"},
		{StatusNonAuthoritativeInfo, "Non-Authoritative Information"},
		{StatusNoContent, "No Content"},
		{StatusResetContent, "Reset Content"},
		{StatusPartialContent, "Partial Content"},
		{StatusMultiStatus, "Multi-Status"},
		{StatusAlreadyReported, "Already Reported"},
		{StatusIMUsed, "IM Used"},

		// Tests for 3xx status codes (Redirection)
		{StatusMultipleChoices, "Multiple Choices"},
		{StatusMovedPermanently, "Moved Permanently"},
		{StatusFound, "Found"},
		{StatusSeeOther, "See Other"},
		{StatusNotModified, "Not Modified"},
		{StatusUseProxy, "Use Proxy"},
		{StatusTemporaryRedirect, "Temporary Redirect"},
		{StatusPermanentRedirect, "Permanent Redirect"},

		// Testing for 4xx status codes (Customer Error)
		{StatusBadRequest, "Bad Request"},
		{StatusUnauthorized, "Unauthorized"},
		{StatusPaymentRequired, "Payment Required"},
		{StatusForbidden, "Forbidden"},
		{StatusNotFound, "Not Found"},
		{StatusMethodNotAllowed, "Method Not Allowed"},
		{StatusNotAcceptable, "Not Acceptable"},
		{StatusProxyAuthRequired, "Proxy Authentication Required"},
		{StatusRequestTimeout, "Request Timeout"},
		{StatusConflict, "Conflict"},
		{StatusGone, "Gone"},
		{StatusLengthRequired, "Length Required"},
		{StatusPreconditionFailed, "Precondition Failed"},
		{StatusRequestEntityTooLarge, "Request Entity Too Large"},
		{StatusRequestURITooLong, "Request URI Too Long"},
		{StatusUnsupportedMediaType, "Unsupported Media Type"},
		{StatusRequestedRangeNotSatisfiable, "Requested Range Not Satisfiable"},
		{StatusExpectationFailed, "Expectation Failed"},
		{StatusTeapot, "I'm a teapot"},
		{StatusMisdirectedRequest, "Misdirected Request"},
		{StatusUnprocessableEntity, "Unprocessable Entity"},
		{StatusLocked, "Locked"},
		{StatusFailedDependency, "Failed Dependency"},
		{StatusTooEarly, "Too Early"},
		{StatusUpgradeRequired, "Upgrade Required"},
		{StatusPreconditionRequired, "Precondition Required"},
		{StatusTooManyRequests, "Too Many Requests"},
		{StatusRequestHeaderFieldsTooLarge, "Request Header Fields Too Large"},

		// Tests for 5xx status codes (Server Error)
		{StatusUnavailableForLegalReasons, "Unavailable For Legal Reasons"},
		{StatusInternalServerError, "Internal Server Error"},
		{StatusNotImplemented, "Not Implemented"},
		{StatusBadGateway, "Bad Gateway"},
		{StatusServiceUnavailable, "Service Unavailable"},
		{StatusGatewayTimeout, "Gateway Timeout"},
		{StatusHTTPVersionNotSupported, "HTTP Version Not Supported"},
		{StatusVariantAlsoNegotiates, "Variant Also Negotiates"},
		{StatusInsufficientStorage, "Insufficient Storage"},
		{StatusLoopDetected, "Loop Detected"},
		{StatusNotExtended, "Not Extended"},
		{StatusNetworkAuthenticationRequired, "Network Authentication Required"},
		{999, ""}, // Test for an unknown code
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := StatusText(test.code)
			if result != test.expected {
				t.Errorf("For code %d, expected '%s', but got '%s'", test.code, test.expected, result)
			}
		})
	}
}

package quick

import "testing"

func TestQuick_HttpStatus(t *testing.T) {
	t.Log(StatusContinue)
	t.Log(StatusSwitchingProtocols)
	t.Log(StatusProcessing)
	t.Log(StatusEarlyHints)

	t.Log(StatusOK)
	t.Log(StatusCreated)
	t.Log(StatusAccepted)
	t.Log(StatusNonAuthoritativeInfo)
	t.Log(StatusNoContent)
	t.Log(StatusResetContent)
	t.Log(StatusPartialContent)

	t.Log(MethodGet)
	t.Log(MethodHead)
	t.Log(MethodPost)
	t.Log(MethodPut)
	t.Log(MethodPatch)
	t.Log(MethodDelete)
	t.Log(MethodConnect)
	t.Log(MethodTrace)
}

func TestStatusText(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{StatusContinue, "Continue"},
		{StatusSwitchingProtocols, "Switching Protocols"},
		{StatusProcessing, "Processing"},
		{StatusEarlyHints, "Early Hints"},
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
		{StatusMultipleChoices, "Multiple Choices"},
		{StatusMovedPermanently, "Moved Permanently"},
		{StatusFound, "Found"},
		{StatusSeeOther, "See Other"},
		{StatusNotModified, "Not Modified"},
		{StatusUseProxy, "Use Proxy"},
		{StatusTemporaryRedirect, "Temporary Redirect"},
		{StatusPermanentRedirect, "Permanent Redirect"},
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
		{999, ""}, // Teste para um código desconhecido
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := StatusText(test.code)
			if result != test.expected {
				t.Errorf("Para código %d, esperado '%s', mas obteve '%s'", test.code, test.expected, result)
			}
		})
	}
}

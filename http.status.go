// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goquick

const (
	// MethodGet represents the HTTP GET method.
	MethodGet = "GET"

	// MethodHead represents the HTTP HEAD method.
	MethodHead = "HEAD"

	// MethodPost represents the HTTP POST method.
	MethodPost = "POST"

	// MethodPut represents the HTTP PUT method.
	MethodPut = "PUT"

	// MethodPatch represents the HTTP PATCH method (RFC 5789).
	MethodPatch = "PATCH"

	// MethodDelete represents the HTTP DELETE method.
	MethodDelete = "DELETE"

	// MethodConnect represents the HTTP CONNECT method.
	MethodConnect = "CONNECT"

	// MethodOptions represents the HTTP OPTIONS method.
	MethodOptions = "OPTIONS"

	// MethodTrace represents the HTTP TRACE method.
	MethodTrace = "TRACE"
)

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (

	// StatusContinue indicates that the client should continue with its request.
	StatusContinue = 100 // RFC 9110, 15.2.1

	// StatusSwitchingProtocols indicates that the server is switching protocols as requested.
	StatusSwitchingProtocols = 101 // RFC 9110, 15.2.2

	// StatusProcessing indicates that the server has received the request but has not completed processing.
	StatusProcessing = 102 // RFC 2518, 10.1

	// StatusEarlyHints provides early hints to help a client start preloading resources.
	StatusEarlyHints = 103 // RFC 8297

	// StatusOK indicates that the request has succeeded.
	StatusOK = 200 // RFC 9110, 15.3.1

	// StatusCreated indicates that the request has been fulfilled and a new resource is created.
	StatusCreated = 201 // RFC 9110, 15.3.2

	// StatusAccepted indicates that the request has been accepted for processing but is not complete.
	StatusAccepted = 202 // RFC 9110, 15.3.3

	// StatusNonAuthoritativeInfo indicates that the response contains information from another source.
	StatusNonAuthoritativeInfo = 203 // RFC 9110, 15.3.4

	// StatusNoContent indicates that the server successfully processed the request but has no content to return.
	StatusNoContent = 204 // RFC 9110, 15.3.5

	// StatusResetContent indicates that the client should reset the document view.
	StatusResetContent = 205 // RFC 9110, 15.3.6

	// StatusPartialContent indicates that only part of the requested resource is returned.
	StatusPartialContent = 206 // RFC 9110, 15.3.7

	// StatusMultiStatus indicates that multiple status codes might be returned.
	StatusMultiStatus = 207 // RFC 4918, 11.1

	// StatusAlreadyReported indicates that the request has already been reported in a previous response.
	StatusAlreadyReported = 208 // RFC 5842, 7.1

	// StatusIMUsed indicates that the response is a result of an instance manipulation.
	StatusIMUsed = 226 // RFC 3229, 10.4.1

	// StatusMultipleChoices indicates that multiple possible resources could be returned.
	StatusMultipleChoices = 300 // RFC 9110, 15.4.1

	// StatusMovedPermanently indicates that the resource has moved permanently to a new URI.
	StatusMovedPermanently = 301 // RFC 9110, 15.4.2

	// StatusFound indicates that the requested resource has been temporarily moved.
	StatusFound = 302 // RFC 9110, 15.4.3

	// StatusSeeOther indicates that the response is available at a different URI.
	StatusSeeOther = 303 // RFC 9110, 15.4.4

	// StatusNotModified indicates that the resource has not been modified since the last request.
	StatusNotModified = 304 // RFC 9110, 15.4.5

	// StatusUseProxy indicates that the requested resource must be accessed through a proxy.
	StatusUseProxy = 305 // RFC 9110, 15.4.6

	// 306 Unused (was previously defined in an earlier version of the HTTP specification).
	_ = 306 // RFC 9110, 15.4.7 (Unused)

	// StatusTemporaryRedirect indicates that the request should be repeated with a different URI.
	StatusTemporaryRedirect = 307 // RFC 9110, 15.4.8

	// StatusPermanentRedirect indicates that the resource has been permanently moved.
	StatusPermanentRedirect = 308 // RFC 9110, 15.4.9

	// StatusBadRequest indicates that the server cannot process the request due to client error.
	StatusBadRequest = 400 // RFC 9110, 15.5.1

	// StatusUnauthorized indicates that authentication is required and has failed or not been provided.
	StatusUnauthorized = 401 // RFC 9110, 15.5.2

	// StatusPaymentRequired is reserved for future use (typically related to digital payments).
	StatusPaymentRequired = 402 // RFC 9110, 15.5.3

	// StatusForbidden indicates that the request is valid, but the server is refusing to process it.
	StatusForbidden = 403 // RFC 9110, 15.5.4

	// StatusNotFound indicates that the requested resource could not be found.
	StatusNotFound = 404 // RFC 9110, 15.5.5

	// StatusMethodNotAllowed indicates that the request method is not allowed for the resource.
	StatusMethodNotAllowed = 405 // RFC 9110, 15.5.6

	// StatusNotAcceptable indicates that the server cannot return a response that meets the client's requirements.
	StatusNotAcceptable = 406 // RFC 9110, 15.5.7

	// StatusProxyAuthRequired indicates that authentication is required for a proxy server.
	StatusProxyAuthRequired = 407 // RFC 9110, 15.5.8

	// StatusRequestTimeout indicates that the server timed out waiting for the request.
	StatusRequestTimeout = 408 // RFC 9110, 15.5.9

	// StatusConflict indicates that the request could not be completed due to a conflict with the current resource state.
	StatusConflict = 409 // RFC 9110, 15.5.10

	// StatusGone indicates that the requested resource is no longer available and will not return.
	StatusGone = 410 // RFC 9110, 15.5.11

	// StatusLengthRequired indicates that the request must include a valid `Content-Length` header.
	StatusLengthRequired = 411 // RFC 9110, 15.5.12

	// StatusPreconditionFailed indicates that a precondition in the request headers was not met.
	StatusPreconditionFailed = 412 // RFC 9110, 15.5.13

	// StatusRequestEntityTooLarge indicates that the request body is too large for the server to process.
	StatusRequestEntityTooLarge = 413 // RFC 9110, 15.5.14

	// StatusRequestURITooLong indicates that the request URI is too long for the server to process.
	StatusRequestURITooLong = 414 // RFC 9110, 15.5.15

	// StatusUnsupportedMediaType indicates that the request body format is not supported by the server.
	StatusUnsupportedMediaType = 415 // RFC 9110, 15.5.16

	// StatusRequestedRangeNotSatisfiable indicates that the range specified in the request cannot be fulfilled.
	StatusRequestedRangeNotSatisfiable = 416 // RFC 9110, 15.5.17

	// StatusExpectationFailed indicates that the server cannot meet the expectations set in the request headers.
	StatusExpectationFailed = 417 // RFC 9110, 15.5.18

	// StatusTeapot is an Easter egg from RFC 9110, originally from April Fools' Day (RFC 2324).
	StatusTeapot = 418 // RFC 9110, 15.5.19 (Unused)

	// StatusMisdirectedRequest indicates that the request was directed to a server that cannot respond appropriately.
	StatusMisdirectedRequest = 421 // RFC 9110, 15.5.20

	// StatusUnprocessableEntity indicates that the request was well-formed but contains semantic errors.
	StatusUnprocessableEntity = 422 // RFC 9110, 15.5.21

	// StatusLocked indicates that the requested resource is currently locked.
	StatusLocked = 423 // RFC 4918, 11.3

	// StatusFailedDependency indicates that the request failed due to a failed dependency.
	StatusFailedDependency = 424 // RFC 4918, 11.4

	// StatusTooEarly indicates that the request was sent too early and should be retried later.
	StatusTooEarly = 425 // RFC 8470, 5.2.

	// StatusUpgradeRequired indicates that the client should switch to a different protocol (e.g., HTTPS).
	StatusUpgradeRequired = 426 // RFC 9110, 15.5.22

	// StatusPreconditionRequired indicates that a precondition header is required for the request.
	StatusPreconditionRequired = 428 // RFC 6585, 3

	// StatusTooManyRequests indicates that the client has sent too many requests in a given period.
	StatusTooManyRequests = 429 // RFC 6585, 4

	// StatusRequestHeaderFieldsTooLarge indicates that the request headers are too large for the server to process.
	StatusRequestHeaderFieldsTooLarge = 431 // RFC 6585, 5

	// StatusUnavailableForLegalReasons indicates that the resource is unavailable for legal reasons (e.g., censorship).
	StatusUnavailableForLegalReasons = 451 // RFC 7725, 3

	// StatusInternalServerError indicates that the server encountered an unexpected condition.
	StatusInternalServerError = 500 // RFC 9110, 15.6.1

	// StatusNotImplemented indicates that the server does not support the requested functionality.
	StatusNotImplemented = 501 // RFC 9110, 15.6.2

	// StatusBadGateway indicates that the server, acting as a gateway or proxy, received an invalid response.
	StatusBadGateway = 502 // RFC 9110, 15.6.3

	// StatusServiceUnavailable indicates that the server is temporarily unable to handle the request (e.g., overloaded or under maintenance).
	StatusServiceUnavailable = 503 // RFC 9110, 15.6.4

	// StatusGatewayTimeout indicates that the server, acting as a gateway or proxy, did not receive a timely response.
	StatusGatewayTimeout = 504 // RFC 9110, 15.6.5

	// StatusHTTPVersionNotSupported indicates that the server does not support the HTTP version used in the request.
	StatusHTTPVersionNotSupported = 505 // RFC 9110, 15.6.6

	// StatusVariantAlsoNegotiates indicates that the server has an internal configuration error preventing negotiation.
	StatusVariantAlsoNegotiates = 506 // RFC 2295, 8.1

	// StatusInsufficientStorage indicates that the server cannot store the representation needed to complete the request.
	StatusInsufficientStorage = 507 // RFC 4918, 11.5

	// StatusLoopDetected indicates that the server detected an infinite loop while processing the request.
	StatusLoopDetected = 508 // RFC 5842, 7.2

	// StatusNotExtended indicates that further extensions to the request are required for the server to fulfill it.
	StatusNotExtended = 510 // RFC 2774, 7

	// StatusNetworkAuthenticationRequired indicates that the client must authenticate to gain network access.
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	switch code {
	case StatusContinue:
		return "Continue"
	case StatusSwitchingProtocols:
		return "Switching Protocols"
	case StatusProcessing:
		return "Processing"
	case StatusEarlyHints:
		return "Early Hints"
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNonAuthoritativeInfo:
		return "Non-Authoritative Information"
	case StatusNoContent:
		return "No Content"
	case StatusResetContent:
		return "Reset Content"
	case StatusPartialContent:
		return "Partial Content"
	case StatusMultiStatus:
		return "Multi-Status"
	case StatusAlreadyReported:
		return "Already Reported"
	case StatusIMUsed:
		return "IM Used"
	case StatusMultipleChoices:
		return "Multiple Choices"
	case StatusMovedPermanently:
		return "Moved Permanently"
	case StatusFound:
		return "Found"
	case StatusSeeOther:
		return "See Other"
	case StatusNotModified:
		return "Not Modified"
	case StatusUseProxy:
		return "Use Proxy"
	case StatusTemporaryRedirect:
		return "Temporary Redirect"
	case StatusPermanentRedirect:
		return "Permanent Redirect"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusPaymentRequired:
		return "Payment Required"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusNotAcceptable:
		return "Not Acceptable"
	case StatusProxyAuthRequired:
		return "Proxy Authentication Required"
	case StatusRequestTimeout:
		return "Request Timeout"
	case StatusConflict:
		return "Conflict"
	case StatusGone:
		return "Gone"
	case StatusLengthRequired:
		return "Length Required"
	case StatusPreconditionFailed:
		return "Precondition Failed"
	case StatusRequestEntityTooLarge:
		return "Request Entity Too Large"
	case StatusRequestURITooLong:
		return "Request URI Too Long"
	case StatusUnsupportedMediaType:
		return "Unsupported Media Type"
	case StatusRequestedRangeNotSatisfiable:
		return "Requested Range Not Satisfiable"
	case StatusExpectationFailed:
		return "Expectation Failed"
	case StatusTeapot:
		return "I'm a teapot"
	case StatusMisdirectedRequest:
		return "Misdirected Request"
	case StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case StatusLocked:
		return "Locked"
	case StatusFailedDependency:
		return "Failed Dependency"
	case StatusTooEarly:
		return "Too Early"
	case StatusUpgradeRequired:
		return "Upgrade Required"
	case StatusPreconditionRequired:
		return "Precondition Required"
	case StatusTooManyRequests:
		return "Too Many Requests"
	case StatusRequestHeaderFieldsTooLarge:
		return "Request Header Fields Too Large"
	case StatusUnavailableForLegalReasons:
		return "Unavailable For Legal Reasons"
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusGatewayTimeout:
		return "Gateway Timeout"
	case StatusHTTPVersionNotSupported:
		return "HTTP Version Not Supported"
	case StatusVariantAlsoNegotiates:
		return "Variant Also Negotiates"
	case StatusInsufficientStorage:
		return "Insufficient Storage"
	case StatusLoopDetected:
		return "Loop Detected"
	case StatusNotExtended:
		return "Not Extended"
	case StatusNetworkAuthenticationRequired:
		return "Network Authentication Required"
	default:
		return ""
	}
}

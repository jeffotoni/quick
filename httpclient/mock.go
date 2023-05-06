package httpclient

import "net/http"

type httpMock struct {
	response *http.Response
	err      error
}

func (h *httpMock) Do(*http.Request) (*http.Response, error) {
	return h.response, h.err
}

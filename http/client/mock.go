package client

import (
	"bytes"
	"net/http"
)

type httpMock struct {
	response *http.Response
	err      error
}

func (h *httpMock) Do(*http.Request) (*http.Response, error) {
	return h.response, h.err
}

var (
	letsgoquickOutMock = `<html>  <head>    <title>Quick - Go</title>  </head>  <body>    <br/>    <br/>    <br/>    <br/>    <h1 style="text-align: center;">Quick - route 100% net/http</h1>  </body></html>`
)

func removeSpaces(b *[]byte) {
	*b = bytes.ReplaceAll(*b, []byte("\t"), []byte(""))
	*b = bytes.ReplaceAll(*b, []byte("\n"), []byte(""))
}

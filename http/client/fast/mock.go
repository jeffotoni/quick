package fast

import (
	"bytes"
	"time"

	f "github.com/valyala/fasthttp"
)

type httpMock struct {
	err error
}

func (h *httpMock) Do(*f.Request, *f.Response) error {
	return h.err
}

func (h *httpMock) DoTimeout(*f.Request, *f.Response, time.Duration) error {
	return h.err
}

var (
	letsgoquickOutMock = `<html>  <head>    <title>Quick - Go</title>  </head>  <body>    <br/>    <br/>    <br/>    <br/>    <h1 style="text-align: center;">Quick - route 100% net/http</h1>  </body></html>`
)

func removeSpaces(b *[]byte) {
	*b = bytes.ReplaceAll(*b, []byte("\t"), []byte(""))
	*b = bytes.ReplaceAll(*b, []byte("\n"), []byte(""))
}

package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type loggerRespWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *loggerRespWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *loggerRespWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			ip, port, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error in logger: %v", err)
				return
			}

			var bodySize int64
			if req.Body != nil {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "error in logger: %v", err)
					return
				}
				bodySize = int64(len(body))
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			lrw := &loggerRespWriter{ResponseWriter: w}
			next.ServeHTTP(lrw, req)

			elapsed := time.Since(start)
			log.Printf("[%s]:%s %d - %s %s %v %d\n", ip, port, lrw.status, req.Method, req.URL.Path, elapsed, bodySize)
		})
	}
}

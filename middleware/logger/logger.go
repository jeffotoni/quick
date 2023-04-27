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

func New() func(http.Handler) http.Handler {
	start := time.Now()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ip, port, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error in logger: %v", err)
				return
			}

			var bodySize int64
			if req.Body != nil {
				body, errReadBody := io.ReadAll(req.Body)

				if errReadBody != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "error in logger: %v", errReadBody)
					return
				}

				bodySize = int64(len(body))
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			log.Printf("[%s]:%s %d - %s %s %v %d\n", ip, port, http.StatusOK, req.Method, req.URL.Path, time.Since(start), bodySize)

			w.WriteHeader(http.StatusOK)
			next.ServeHTTP(w, req)
		})
	}
}

package maxbody

import (
	"fmt"
	"net/http"
)

const defaultMaxBytes int64 = 1024 * 1024 * 5

func New(maxBytes ...int64) func(http.Handler) http.Handler {
	mb := defaultMaxBytes
	if len(maxBytes) > 0 {
		mb = maxBytes[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > mb {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				fmt.Fprint(w, "Request body too large")
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, mb)
			next.ServeHTTP(w, r)
		})
	}
}

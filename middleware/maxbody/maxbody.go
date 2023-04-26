package maxbody

import (
	"net/http"
)

const DefaultMaxBytes int64 = 1024 * 1024 * 5

func New(maxBytes ...int64) func(http.Handler) http.Handler {
	mb := DefaultMaxBytes
	if len(maxBytes) > 0 && maxBytes != nil {
		mb = maxBytes[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > mb {
				http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, mb)
			next.ServeHTTP(w, r)
		})
	}
}

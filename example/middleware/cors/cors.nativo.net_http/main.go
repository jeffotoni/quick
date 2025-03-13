package main

import (
	"io"
	"net/http"

	"github.com/jeffotoni/quick/middleware/cors"
)

type MyHandler struct{}

// / Example cURL to test
// curl -X OPTIONS -v http://localhost:8080/v1/user
//
//	curl -X OPTIONS -H "Origin: http://localhost:3000/" \
//	 -H "Access-Control-Request-Method: POST" -v \
//	 http://localhost:8080/v1/user
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"msg":"error"}`))
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func OtherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("Outro endpoint!"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/v1/user", &MyHandler{})
	mux.HandleFunc("/outro", OtherHandler)

	newmux := cors.Default().Handler(mux)
	println("server: :8080")
	http.ListenAndServe(":8080", newmux)
}

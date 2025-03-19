package main

import (
	"log"
	"net/http"
	//"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	mux := http.NewServeMux()

	// Route with compression enabled using the Gzip middleware
	// mux.Handle("/v1/compress", compress.Gzip()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`{"message": "Hello, net/http with Gzip!"}`))
	// })))

	// Starting the HTTP server on port 8080
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// $ curl -X GET http://localhost:8080/v1/compress -H "Accept-Encoding: gzip" --compressed -i

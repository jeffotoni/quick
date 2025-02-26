//go:build !exclude_test

package main

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	fs := http.FileServer(http.FS(staticFiles))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	port := 8080
	fmt.Printf("Server Run http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

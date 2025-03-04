package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, POST Quick New!!!",
	}
	resp, err := client.Put("https://httpbin.org/put", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PUT response:", string(resp.Body))
}

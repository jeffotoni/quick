package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	resp, err := client.Delete("http://localhost:3000/v1/user/1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DELETE response:", string(resp.Body))
}

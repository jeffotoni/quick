package main

import (
	"fmt"
	"github.com/jeffotoni/quick"
)

func main() {
	// start quick
	q := quick.New()

	// instance Listen
	fmt.Println("Run Server port:403")
	err := q.ListenTLS(":443", "cert.pem", "key.pem")
	if err != nil {
		fmt.Printf("error when trying to connect with TLS: %v\n", err)
	}
}

package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// $ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
//
// $ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
func main() {
	// Initialize Quick instance
	q := quick.New()

	// Print a message indicating that the server is starting on port 443
	fmt.Println("Run Server port:443")

	// Start the HTTPS server with TLS encryption
	// - The server will listen on port 443
	// - cert.pem: SSL/TLS certificate file
	// - key.pem: Private key file for SSL/TLS encryption
	// - http2 [true or false]
	err := q.ListenTLS(":443", "cert.pem", "key.pem", false)
	if err != nil {
		// Log an error message if the server fails to start
		fmt.Printf("Error when trying to connect with TLS: %v\n", err)
	}
}

//**Note for Linux Users**
//By default, this example runs on **port 443**, which is a privileged port (below 1024).
//On **Linux**, running on this port requires **superuser privileges**.

//To run this example on Linux, use:
//sudo go run main.go

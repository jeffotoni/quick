package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cClient := client.New(
		client.WithTimeout(10*time.Second),
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),
		client.WithRetryTransport(
			100,                       // MaxIdleConns
			10,                        // MaxIdleConnsPerHost
			50,                        // MaxConnsPerHost
			false,                     // DisableKeepAlives
			true,                      // Force HTTP/2
			http.ProxyFromEnvironment, // System default proxy
			&tls.Config{ // Custom TLS configuration
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		),
	)

	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	fmt.Println("POST Form Response:", string(resp.Body))
}

package main

import (
	"log"
	"time"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New(quick.Config{
		GCPercent:         500, // GC more aggressive for high load
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1024 * 1024 * 20, // 20MB
		// cors
		CorsConfig: &quick.CorsConfig{
			Enabled: true,
			Options: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST",
			},
		},
	})

	q.Get("/v1/user/:id", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		// time.Sleep(2 * time.Second)
		// fmt.Println("log:", c.Request.URL.Path)
		myuser := struct {
			Name string `json:"name"`
		}{
			Name: c.Param("id"),
		}
		return c.Status(200).JSON(myuser)
	})

	//config: quick.Config{
	//	GCPercent:         500,  // GC mais agressivo para alta carga
	//	ReadTimeout:       10 * time.Second,
	//	WriteTimeout:      10 * time.Second,
	//	IdleTimeout:       30 * time.Second,
	//	ReadHeaderTimeout: 2 * time.Second,
	//	MaxHeaderBytes:    1 << 20, // 1MB
	//},
	//CorsConfig: &CorsConfig{
	//	Enabled: true,
	//	Options: map[string]string{
	//		"Access-Control-Allow-Origin":  "*",
	//		"Access-Control-Allow-Methods": "GET,POST",
	//	},
	//}
	if err := q.ListenTLS(":443", "cert.pem", "key.pem", false); err != nil {
		log.Fatal(err)
	}
}

//
//func main() {
//	q := quick.New()
//
//	// Route to respond with a simulated delay
//	q.Get("/v1/user/:id", func(c *quick.Ctx) error {
//		c.Set("Content-Type", "application/json")
//		// time.Sleep(2 * time.Second)
//		// fmt.Println("log:", c.Request.URL.Path)
//		return c.Status(200).String(c.Request.URL.Path)
//	})
//
//	// Start the HTTPS server with HTTP/2 enabled
//	fmt.Println("Servidor HTTP/2 rodando em https://localhost:443...")
//
//	// Start the HTTPS server with TLS encryption
//	// - The server will listen on port 443
//	// - cert.pem: SSL/TLS certificate file
//	// - key.pem: Private key file for SSL/TLS encryption
//	//err := q.ListenTLS(":443", "cert.pem", "key.pem", false)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	if err := q.Listen("0.0.0.0:8080"); err != nil {
//		log.Fatal(err)
//	}
//}

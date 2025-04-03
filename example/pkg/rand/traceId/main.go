package main

import (
	"fmt"

	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	traceID := rand.TraceID()
	fmt.Printf("Generated Trace ID: %s\n", traceID)
}

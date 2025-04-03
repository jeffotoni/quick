package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	min := 10
	max := 100
	randomNumber, err := rand.RandomInt(min, max)
	if err != nil {
		log.Fatalf("Failed to generate random integer: %v", err)
	}
	fmt.Printf("Random integer between %d and %d: %d\n", min, max, randomNumber)
}

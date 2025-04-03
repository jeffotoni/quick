package main

import (
	"fmt"

	"github.com/jeffotoni/quick/pkg/uuid"
)

func main() {
	uuid := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	fmt.Println("String:", uuid.String())
	fmt.Println("URN:", uuid.URN())
}

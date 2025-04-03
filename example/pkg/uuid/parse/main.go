package main

import (
	"fmt"

	"github.com/jeffotoni/quick/pkg/uuid"
)

func main() {
	s := "123e4567-e89b-12d3-a456-426614174000"
	uuid, err := uuid.Parse(s)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	fmt.Println("UUID:", uuid)
}

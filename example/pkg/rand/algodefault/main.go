package main

import (
	"fmt"

	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	start := 1
	end := 500
	randomMsgID := rand.AlgoDefault(start, end)
	fmt.Printf("Generated MsgID: %s\n", randomMsgID)
}

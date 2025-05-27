package rand

import (
	"fmt"
	"log"
)

// This function is named ExampleRandomInt()
// it with the Examples type.
func ExampleRandomInt() {
	_, err := RandomInt(10, 20)
	if err != nil {
		log.Fatal(err)
	}

	// Simulated value for documentation purposes
	//fmt.Println("Random number generated: 17")

	fmt.Println("Random number generated")
	// Output: Random number generated
}

// This function is named ExampleTraceID()
// it with the Examples type.
func ExampleTraceID() {
	_ = TraceID()

	// Simulated trace ID format
	//fmt.Println("Trace ID generated: a8b7XkP1ZcDqWmE2")

	fmt.Println("Trace ID generated")
	// Output: Trace ID generated
}

// This function is named ExampleAlgoDefault()
// it with the Examples type.
func ExampleAlgoDefault() {
	_ = AlgoDefault(1000, 9999)

	// Simulated Msg ID for example purposes
	//fmt.Println("Msg Id generated: 3842")

	fmt.Println("Msg Id generated")
	// Output: Msg Id generated
}

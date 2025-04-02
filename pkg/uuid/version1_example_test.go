package uuid

import "fmt"

// This function is named ExampleNewUUID()
// it with the Examples type.
func ExampleNewUUID() {
	u, err := NewUUID()
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(u.Version())
	// Output: VERSION_1
}

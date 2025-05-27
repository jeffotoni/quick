package uuid

import (
	"fmt"
	"time"
)

// This function is named ExampleNewV6()
// it with the Examples type.
func ExampleNewV6() {
	u, err := NewV6()
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(u.Version(), u.Variant())
	// Output: VERSION_6 RFC4122
}

// This function is named ExampleNewV6WithTime()
// it with the Examples type.
func ExampleNewV6WithTime() {
	t := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	u, err := NewV6WithTime(&t)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(u.Version(), u.Variant())
	// Output: VERSION_6 RFC4122
}

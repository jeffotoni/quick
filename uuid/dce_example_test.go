package uuid

import (
	"fmt"
)

// This function is named ExampleNewDCESecurity()
// it with the Examples type.
func ExampleNewDCESecurity() {
	u, err := NewDCESecurity(Person, 1234)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Printf("Domain: %s, ID: %d, Version: %s\n", u.Domain(), u.ID(), u.Version())
	// Output: Domain: Person, ID: 1234, Version: VERSION_2
}

// This function is named ExampleNewDCEPerson()
// it with the Examples type.
func ExampleNewDCEPerson() {
	u, err := NewDCEPerson()
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(u.Domain())
	// Output: Person
}

// This function is named ExampleNewDCEGroup()
// it with the Examples type.
func ExampleNewDCEGroup() {
	u, err := NewDCEGroup()
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(u.Domain())
	// Output: Group
}

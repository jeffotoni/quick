package quick

import (
	"fmt"
)

// This function is named ExampleStatusText()
// it with the Examples type.
func ExampleStatusText() {
	fmt.Println(StatusText(200))
	fmt.Println(StatusText(404))
	fmt.Println(StatusText(500))

	// Output:
	// OK
	// Not Found
	// Internal Server Error
}

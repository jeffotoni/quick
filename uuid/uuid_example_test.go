package uuid

import (
	"fmt"
	"log"
)

// This function is named ExampleParse()
// it with the Examples type.
func ExampleParse() {
	u, err := Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.String())
	// Output: 6ba7b810-9dad-11d1-80b4-00c04fd430c8
}

// This function is named ExampleParseBytes()
// it with the Examples type.
func ExampleParseBytes() {
	b := []byte("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	u, err := ParseBytes(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.String())
	// Output: f47ac10b-58cc-0372-8567-0e02b2c3d479
}

// This function is named ExampleValidate()
// it with the Examples type.
func ExampleValidate() {
	uuidStr := "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}"
	if err := Validate(uuidStr); err == nil {
		fmt.Println("UUID is valid")
	} else {
		fmt.Println("UUID is invalid")
	}
	// Output: UUID is valid
}

// This function is named ExampleFromBytes()
// it with the Examples type.
func ExampleFromBytes() {
	b := []byte{
		0x7d, 0x44, 0x48, 0x40,
		0x9d, 0xc0,
		0x11, 0xd1,
		0xb2, 0x45,
		0x5f, 0xfd, 0xce, 0x74, 0xfa, 0xd2,
	}
	u, err := FromBytes(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.String())
	// Output: 7d444840-9dc0-11d1-b245-5ffdce74fad2
}

// This function is named ExampleUUID_String()
// it with the Examples type.
func ExampleUUID_String() {
	u := MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	fmt.Println(u.String())
	// Output: f47ac10b-58cc-0372-8567-0e02b2c3d479
}

// This function is named ExampleUUID_URN()
// it with the Examples type.
func ExampleUUID_URN() {
	u := MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	fmt.Println(u.URN())
	// Output: urn:uuid:f47ac10b-58cc-0372-8567-0e02b2c3d479
}

// This function is named ExampleMust()
// it with the Examples type.
func ExampleMust() {
	u := Must(Parse("f47ac10b-58cc-0372-8567-0e02b2c3d479"))
	fmt.Println(u.URN())
	// Output: urn:uuid:f47ac10b-58cc-0372-8567-0e02b2c3d479
}

// This function is named Example_versionAndVariant()
// it with the Examples type.
func Example_versionAndVariant() {
	u := MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479") // version 4
	fmt.Printf("Version: %s, Variant: %s\n", u.Version(), u.Variant())
	// Output: Version: VERSION_4, Variant: RFC4122
}

// This function is named ExampleUUID_Version()
// it with the Examples type.
func ExampleUUID_Version() {
	u := MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	fmt.Println(u.Version())
	// Output: VERSION_4
}

// This function is named ExampleUUID_Variant()
// it with the Examples type.
func ExampleUUID_Variant() {
	u := MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	fmt.Println(u.Variant())
	// Output: RFC4122
}

// This function is named ExampleUUIDs_Strings()
// it with the Examples type.
func ExampleUUIDs_Strings() {
	u1 := MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	u2 := MustParse("7d444840-9dc0-11d1-b245-5ffdce74fad2")
	us := UUIDs{u1, u2}
	fmt.Println(us.Strings())
	// Output: [f47ac10b-58cc-0372-8567-0e02b2c3d479 7d444840-9dc0-11d1-b245-5ffdce74fad2]
}

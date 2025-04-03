package uuid

import (
	"fmt"
)

// This function is named ExampleUUID_MarshalText()
// it with the Examples type.
func ExampleUUID_MarshalText() {
	uuid, err := Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		fmt.Println("error parsing UUID:", err)
		return
	}
	text, err := uuid.MarshalText()
	if err != nil {
		fmt.Println("error marshaling to text:", err)
		return
	}
	fmt.Printf("%s\n", text)
	// Output: 123e4567-e89b-12d3-a456-426614174000
}

// This function is named ExampleUUID_UnmarshalText()
// it with the Examples type.
func ExampleUUID_UnmarshalText() {
	var uuid UUID
	data := []byte("123e4567-e89b-12d3-a456-426614174000")
	err := uuid.UnmarshalText(data)
	if err != nil {
		fmt.Println("error unmarshaling text:", err)
		return
	}
	fmt.Printf("UUID: %v\n", uuid)
	// Output: UUID: 123e4567-e89b-12d3-a456-426614174000
}

// This function is named ExampleUUID_MarshalBinary()
// it with the Examples type.
func ExampleUUID_MarshalBinary() {
	// Criando um UUID fixo para consistência nos testes
	uuid, err := Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		fmt.Printf("failed to parse UUID: %v\n", err)
		return
	}
	data, err := uuid.MarshalBinary()
	if err != nil {
		fmt.Printf("failed to marshal UUID: %v\n", err)
		return
	}
	fmt.Printf("Binary data: %x\n", data)
	// Output: Binary data: 123e4567e89b12d3a456426614174000
}

// This function is named ExampleUUID_UnmarshalBinary()
// it with the Examples type.
func ExampleUUID_UnmarshalBinary() {
	var uuid UUID
	// Representação binária exemplo de um UUID
	data := []byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00}
	err := uuid.UnmarshalBinary(data)
	if err != nil {
		fmt.Println("error unmarshaling binary:", err)
		return
	}
	fmt.Printf("UUID: %v\n", uuid)
	// Output: UUID: 123e4567-e89b-12d3-a456-426614174000
}

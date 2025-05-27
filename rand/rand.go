// Package rand provides utility functions for generating random values using only the Go standard library.
//
// It includes support for secure cryptographic randomness via crypto/rand
// and pseudo-random generators via math/rand for lightweight use cases.
//
// Functions:
//
//   - RandomInt: securely generates a random int between a min and max range
//   - TraceID: generates a random string (alphanumeric trace ID)
//   - AlgoDefault: generates a secure random MsgID between a given range
package rand

import (
	"crypto/rand"
	randx "crypto/rand"
	"log"
	"math/big"
	randm "math/rand"
	"strconv"
	"time"
)

// RandomInt returns a cryptographically secure random integer in the range [min, max).
//
// It uses crypto/rand and big.Int for security and supports large ranges.
func RandomInt(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max))
	minBigInt := big.NewInt(int64(min))
	diffBigInt := new(big.Int).Sub(maxBigInt, minBigInt)

	randomBytes := make([]byte, diffBigInt.BitLen()/8+1)
	_, err := randx.Read(randomBytes)
	if err != nil {
		return 0, err
	}

	randomInt := new(big.Int).SetBytes(randomBytes)
	randomInt.Mod(randomInt, diffBigInt)
	randomInt.Add(randomInt, minBigInt)
	return int(randomInt.Int64()), nil
}

// TraceID generates a pseudo-random alphanumeric trace ID of fixed length (16).
//
// It uses math/rand seeded with time.Now() for simplicity.
// The output is suitable for temporary identifiers (non-cryptographic).
func TraceID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16

	// Initialize the seed once per run
	src := randm.NewSource(time.Now().UnixNano())
	r := randm.New(src)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// AlgoDefault generates a secure random integer between Start and End
// using crypto/rand and returns it as a string.
//
// If an error occurs during generation, an empty string is returned.
func AlgoDefault(Start, End int) string {
	max := big.NewInt(int64(End))              // Define the upper limit for the random number
	randInt, err := rand.Int(rand.Reader, max) // Generate a secure random number
	if err != nil {
		log.Printf("error: %v", err)
		return "" // Return an empty string if an error occurs
	}
	// Return the generated MsgID as a string
	return strconv.Itoa(Start + int(randInt.Int64()))
}

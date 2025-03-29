package rand

import (
	randx "crypto/rand"
	"math/big"
	randm "math/rand"
	"time"
)

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

// Generate Trace ID with letters and numbers, 100% stdlib
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

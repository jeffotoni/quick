package rand

import (
	"crypto/rand"
	"math/big"
)

func RandomInt(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max))
	minBigInt := big.NewInt(int64(min))
	diffBigInt := new(big.Int).Sub(maxBigInt, minBigInt)

	randomBytes := make([]byte, diffBigInt.BitLen()/8+1)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return 0, err
	}

	randomInt := new(big.Int).SetBytes(randomBytes)
	randomInt.Mod(randomInt, diffBigInt)
	randomInt.Add(randomInt, minBigInt)
	return int(randomInt.Int64()), nil
}

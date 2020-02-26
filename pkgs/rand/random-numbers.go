package rand

import (
	"crypto/sha256"
	"encoding/binary"
)

const MAX_UINT_32 = 4294967295

func GenerateRandomNumbers(input []byte) []uint32 {
	sum := sha256.Sum256(input)

	numbers := []uint32{}
	for i := 0; i < 32; i += 4 {
		numbers = append(numbers, binary.BigEndian.Uint32(sum[i:i+4]))
	}

	return numbers
}

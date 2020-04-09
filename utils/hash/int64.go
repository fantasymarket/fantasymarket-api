package hash

import (
	"crypto/sha256"
	"encoding/binary"
	"strconv"
)

// Int64HashRange hashes a string to an int64 in [min,max).
// It panics if min < max.
func Int64HashRange(min int64, max int64, seed string) int64 {
	if min > max {
		panic("invalid argument to Int64HashRange: min cannot be greater than max")
	}

	if min == max {
		return min
	}

	x := Int64Hashn(max-min, seed) + min
	return x
}

// Int64Hashn hashes a string to an int64 in [0,n).
// It panics if n <= 0.
// Based on https://golang.org/src/math/rand/rand.go?s=10421:10447#L317
// originlly written and licensed by The Go Autor under The 3-Clause BSD License
func Int64Hashn(n int64, seed string) int64 {
	if n <= 0 {
		panic("invalid argument to Int64Hashn: n can't be less than 0")
	}

	// Performance optimisation, if n is a power of 2, we can mask
	if n&(n-1) == 0 {
		return Int64Hash(seed) & (n - 1)
	}

	// `max` is the maximum size our v can be for `v % n` to be smaller than n
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := Int64Hash(seed)

	// We generate new random values until they are smaller than the defined max
	for i := 0; v > max; i++ {
		v = Int64Hash(seed + strconv.Itoa(i))
	}

	return v % n
}

// Int64Hash returns a non-negative int64 from a string
// (based on the sha256 hash-function)
func Int64Hash(seed string) int64 {

	// While sha256 distribution isn't 100% perfectly random,
	// it is random enough for our pourposes.

	hash := sha256.Sum256([]byte(seed))
	// we use a uint64 since we only care about positive values
	n := binary.LittleEndian.Uint32(hash[:4])

	return int64(n)
}

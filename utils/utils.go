package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"strconv"
)

// RandInt64Range returns a random int64 as an int64 in [min,max).
// It panics if min < max.
func RandInt64Range(min int64, max int64, seed string) int64 {
	if min > max {
		panic("invalid argument to RandInt64: Min cannot be greater than Max")
	}
	x := randInt64n(max-min, seed) + min
	return x
}

// randInt64n returns a random int64 as an int64 in [0,n).
// It panics if n <= 0.
// Based on https://golang.org/src/math/rand/rand.go?s=10421:10447#L317
func randInt64n(n int64, seed string) int64 {
	if n <= 0 {
		panic("invalid argument to max")
	}

	if n&(n-1) == 0 { // n is power of two, can mask
		return randInt64(seed) & (n - 1)
	}

	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := randInt64(seed)

	// keep generating random numbers until it's within range
	for i := 0; v > max; i++ {
		v = randInt64(seed + strconv.Itoa(i))
	}

	return v % n
}

// randInt64 returns, as an int, a non-negative pseudo-random number in [min, max].
// if we plan to pass only one value n will be set to 0 and therefore won't affect the value we pass to m (the value we care about)
func randInt64(seed string) int64 {

	var n int64
	hash := sha256.Sum256([]byte(seed))
	buf := bytes.NewBuffer(hash[:8])
	binary.Read(buf, binary.LittleEndian, &n)

	return n
}

package utils

import (
	"math/rand"
)

// RandInt64 returns, as an int, a non-negative pseudo-random number in [min, max].
// if we plan to pass only one value n will be set to 0 and therefore won't affect the value we pass to m (the value we care about)
func RandInt64(min int64, max int64, seed int64) int64 {
	rand.Seed(seed)
	if min > max {
		panic("invalid argument to RandInt64: Min cannot be greater than Max")
	}

	x := rand.Int63n(max-min) + min
	return x
}

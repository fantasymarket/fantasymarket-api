package utils

import (
	"math/rand"
	"time"
)

// if we plan to pass only one value n will be set to 0 and therefore won't affect the value we pass to m (the value we care about)
func RandInt64(n int, m int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(m-n) + n)
}

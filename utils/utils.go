package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// if we plan to pass only one value n will be set to 0 and therefore won't affect the value we pass to m (the value we care about)
func RandInt64(min int, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	x := int64(rand.Intn(max-min) + min)
	fmt.Println("RANDOM: ", x)
	return x
}

package hash

// Float64HashRange hashes a string to an float64 in [min,max).
// precision is the
// It panics if min < max.
func Float64HashRange(min float64, max float64, seed string) float64 {
	if min > max {
		panic("invalid argument to Int64HashRange: min cannot be greater than max")
	}

	if min == max {
		return min
	}

	return min + Float64Hash(seed)*(max-min)
}

// Float64Hash returns, as a float64, a pseudo-random number in [0.0,1.0).
// Based on https://golang.org/src/math/rand/rand.go?s=5359:5391#L180
// originlly written and licensed by The Go Autor under The 3-Clause BSD License
func Float64Hash(seed string) float64 {
	return float64(Int64Hashn(1<<53, seed)) / (1 << 53)
}

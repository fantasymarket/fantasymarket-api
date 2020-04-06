package utils

import "sort"

// Includes checks if a []string includes a string
func Includes(slice []string, value string) bool {
	i := sort.SearchStrings(slice, value)
	return i < len(slice) && slice[i] == value
}

// Some checks if a []string includes at least one string from a slice
func Some(slice []string, values []string) bool {
	some := false
	for _, v := range values {
		if Includes(slice, v) {
			some = true
		}
	}
	return some
}

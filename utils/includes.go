package utils

import "sort"

// Includes checks if a []string includes a string
func Includes(slice []string, value string) bool {
	i := sort.SearchStrings(slice, value)
	return i < len(slice) && slice[i] == value
}

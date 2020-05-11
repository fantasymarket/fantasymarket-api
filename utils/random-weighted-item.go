package utils

import (
	"errors"
	"fmt"
	"sort"
)

var (
	// ErrInvalidRandomNumber happens when a random number is invalid
	ErrInvalidRandomNumber = errors.New("utils: invalid random number")
	// ErrInvalidChance happens when an item's chance is invalid
	ErrInvalidChance = errors.New("utils: invalid chance for item")
)

// SelectRandomWeightedItem selects a random item from a map[item]chance
// Needs a map of items and their weight and a random Number 0 < n < 1
// Returns an empty string if no item was selected
func SelectRandomWeightedItem(items map[string]float64, randomNumber float64) (string, error) {
	// r := hash.Int64HashRange(0, 100, seed)

	if randomNumber <= 0 || randomNumber >= 1 {
		return "", ErrInvalidRandomNumber
	}

	var lowerBound float64
	sortedKeys := sortItemMap(items)

	for _, item := range sortedKeys {
		chance := items[item]

		if chance < 0 || chance > 1 {
			return "", fmt.Errorf("item %v: %w", item, ErrInvalidChance)
		}

		if chance == 0 { // 0 is the nil value of Chance
			return item, nil
		}

		if lowerBound <= randomNumber && randomNumber <= chance {
			return item, nil
		}

		lowerBound = chance
	}

	return "", nil
}

// sortItemMap sorts a map's keys by it's float64 value
func sortItemMap(items map[string]float64) []string {
	// reverse the map to make sorting easier
	itemsReversed := make(map[float64]string)
	for k, v := range items {
		itemsReversed[v] = k
	}

	// sort the map since chances need to be ordered
	// for our algorithm to work
	var keys []float64
	for k := range itemsReversed {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	var values []string
	for _, key := range keys {
		values = append(values, itemsReversed[key])
	}

	return values
}

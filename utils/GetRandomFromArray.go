package utils

import (
	"errors"
	"fantasymarket/utils/hash"
	"fmt"
	"sort"
)

// SelectRandomWeightedItem selects a random item from a map[item]chance
// Returns an empty string if no item was selected
func SelectRandomWeightedItem(items map[string]float64, seed string) (string, error) {
	r := hash.Int64HashRange(0, 100, seed)
	randomFloat := float64(r) / 100 // Get the float for computation (.2f)
	var lowerBound float64 = 0

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

	for _, chance := range keys {
		item := itemsReversed[chance]

		if chance < 0 || chance > 1 {
			return "", errors.New("error: invalid chance for item `" + item + "`")
		}

		if chance == 0 { // 0 is the nil value of Chance
			return item, nil
		}

		fmt.Println("lowerBound: ", lowerBound)
		if lowerBound <= randomFloat && randomFloat <= chance {
			return item, nil
		}

		fmt.Println("set lower bound to", chance)
		lowerBound = chance
	}

	return "", nil
}

package utils_test

import (
	"fantasymarket/utils"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

type getRandomFromArrayTestData struct {
	items        map[string]float64
	randomNumber float64
	expectation  string
}

var getRandomFromArrayData = []getRandomFromArrayTestData{
	{
		items: map[string]float64{
			"event1":  0.1,
			"event2":  0.2,
			"event3":  0.3,
			"event4":  0.69,
			"event5 ": 0.92,
			"event6 ": 0.93,
			"event7 ": 1,
		},
		randomNumber: 0.22,
		expectation:  "event3",
	},
	{
		items: map[string]float64{
			"event1": 0,
			"event2": 0.5,
			"event3": 1,
		},
		randomNumber: 0.22,
		expectation:  "event1",
	},
	{
		items: map[string]float64{
			"TestCase0.001": 0.11,
		},
		randomNumber: 0.22,
		expectation:  "",
	},
	{
		items: map[string]float64{
			"TestCase0": 0,
		},
		randomNumber: 0.22,
		expectation:  "TestCase0",
	},
	{
		items:        map[string]float64{},
		randomNumber: 0.22,
		expectation:  "",
	},
}

func TestSelectRandomWeightedItem(t *testing.T) {
	for _, test := range getRandomFromArrayData {
		val, err := utils.SelectRandomWeightedItem(test.items, test.randomNumber)
		if assert.NoError(t, err) {
			assert.Equal(t, test.expectation, val)
		}
	}
}

func TestSelectRandomWeightedItemInvalidChance(t *testing.T) {
	_, err := utils.SelectRandomWeightedItem(map[string]float64{
		"test": -1,
	}, 0.22)

	assert.True(t, errors.Is(err, utils.ErrInvalidChance))
}

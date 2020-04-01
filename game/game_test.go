package game

import (
	"fantasymarket/database/models"
	"fantasymarket/game/structs"

	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	stock        models.Stock
	affectedness int64
	expectation  map[string]int64
}

var data = []TestData{
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: -1,
		expectation: map[string]int64{
			"TEST0": 1,
			"TEST1": 5,
			"TEST2": 21,
			"TEST3": -5,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: -100,
		expectation: map[string]int64{
			"TEST0": -494,
			"TEST1": -490,
			"TEST2": -474,
			"TEST3": -500,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: 100,
		expectation: map[string]int64{
			"TEST0": 506,
			"TEST1": 510,
			"TEST2": 526,
			"TEST3": 500,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: 0,
		expectation: map[string]int64{
			"TEST0": 6,
			"TEST1": 10,
			"TEST2": 26,
			"TEST3": 0,
		},
	},
}

var stockSettings = map[string]structs.StockSettings{
	"TEST0": {Stability: 1, Trend: 1},
	"TEST1": {Stability: 5, Trend: 1},
	"TEST2": {Stability: 1, Trend: 5},
	"TEST3": {Stability: 0, Trend: 0},
}

func TestGetTendency(t *testing.T) {
	s := Service{
		StockSettings: stockSettings,
	}

	for _, test := range data {
		for i := range stockSettings {
			test.stock.Symbol = string(i)

			result := s.GetTendency(test.stock, test.affectedness)
			assert.Equal(t, test.expectation[i], result)
		}
	}
}

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
			"TEST0": 10,
			"TEST1": 34,
			"TEST2": 30,
			"TEST3": -1,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: -100,
		expectation: map[string]int64{
			"TEST0": -89,
			"TEST1": -65,
			"TEST2": -69,
			"TEST3": -100,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: 100,
		expectation: map[string]int64{
			"TEST0": 111,
			"TEST1": 135,
			"TEST2": 131,
			"TEST3": 100,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: 0,
		expectation: map[string]int64{
			"TEST0": 11,
			"TEST1": 35,
			"TEST2": 31,
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
	s := GameService{
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

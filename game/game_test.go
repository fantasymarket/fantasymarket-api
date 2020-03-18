package game

import (
	"fantasymarket/database/models"
	"fmt"
	"testing"
	"time"
)

type TestData struct {
	stock        models.Stock
	affectedness int64
	expectation  map[string]int64
}

var data = []TestData{
	{
		stock:        models.Stock{Index: 10000, StockID: "TEST"},
		affectedness: -1,
		expectation: map[string]int64{
			"TEST0": 10,
			"TEST1": 34,
			"TEST2": 30,
			"TEST3": -1,
		},
	},
	{
		stock:        models.Stock{Index: 10000, StockID: "TEST"},
		affectedness: -100,
		expectation: map[string]int64{
			"TEST0": -89,
			"TEST1": -65,
			"TEST2": -69,
			"TEST3": -100,
		},
	},
	{
		stock:        models.Stock{Index: 10000, StockID: "TEST"},
		affectedness: 100,
		expectation: map[string]int64{
			"TEST0": 111,
			"TEST1": 135,
			"TEST2": 131,
			"TEST3": 100,
		},
	},
	{
		stock:        models.Stock{Index: 10000, StockID: "TEST"},
		affectedness: 0,
		expectation: map[string]int64{
			"TEST0": 11,
			"TEST1": 35,
			"TEST2": 31,
			"TEST3": 0,
		},
	},
}

var stockSettings = map[string]StockSettings{
	"TEST0": {Stability: 1, Trend: 1},
	"TEST1": {Stability: 5, Trend: 1},
	"TEST2": {Stability: 1, Trend: 5},
	"TEST3": {Stability: 0, Trend: 0},
}

func TestGetTendency(t *testing.T) {
	fmt.Println("Testing getTendancy")

	s := GameService{
		StockSettings: stockSettings,
	}

	for _, test := range data {
		for i, _ := range stockSettings {
			test.stock.StockID = string(i)

			if result := s.GetTendency(test.stock, test.affectedness, time.Unix(5, 0)); result != test.expectation[i] {
				t.Fatal("Expected ", test.expectation[i], ", got ", result)
			}
		}
	}
}

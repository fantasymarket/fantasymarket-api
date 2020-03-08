package game

import "time"

type EventSettings struct {
	EventID string

	Title string
	Text  string

	MinTimeBetweenEvents time.Duration
	Chance               float64 // 0 - 1

	Tags map[string]TagOptions
}

// StockSettings is the Stock "Class"
type StockSettings struct {
	StockID   string          // Stock Symbol e.g GOOG
	Name      string          // Stock Name e.g Alphabet Inc.
	Index     int64           // Price per share
	Shares    int64           // Number per share
	Tags      map[string]bool //A stock can have up to 5 tags
	Stability int64           // Shows how many fluctuations the stock will have
	Trend     int64           // Shows the generall trend of the Stock
}

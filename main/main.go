package main

import (
	"fmt"
	"time"
)

/// Was wir noch nehmen können
/// SQL:			https://github.com/jmoiron/sqlx
/// Decimal:	https://github.com/shopspring/decimal


//FantasyMarketOptions manages the Options of the programm
type FantasyMarketOptions struct {
	TicksPerSecond  float64
	GameTimePerTick time.Duration
}

// Stock is the Stock "Class"
type Stock struct {
	ID    string						// Stock Symbol e.g GOOG
	Name  string						// Stock NAme e.g Alphabet Inc.
	Index string						// Price per share
	Shares int64						// Price per share
	Tags  map[string]bool
}

// Events happen randomly every game tick
type Event struct {
	MinTimeBetweenEvents time.Duration
	Chance               float64 // 0 - 1

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	Duration time.Duration // Time during which the event is the event is run every tick

	// We use tags if we only want to affect only specific stocks
	Tags map[string]TagOptions
}

//TagOptions are more indepth settings for specific event tags only
type TagOptions struct {
	Trend float64
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	Duration time.Duration
}

const (
	// Hour is the duration of 60 minutes
	Hour = time.Second * 60 * 60
	//Day is the duration o 24 hours
	Day = Hour * 24
)

func main() {
	stockArray := []Stock{
		{
			ID:    "GOOG",
			Index: float64(100.00),
			Tags:  map[string]bool{"tech": true, "intl": true},
		},
		{
			ID:    "FRIZ",
			Index: float64(139.69),
			Tags:  map[string]bool{"food": true, "local": true},
		},
		{
			ID:    "LMAO",
			Index: float64(139.69),
			Tags:  map[string]bool{"arthur": true, "henry": true},
		},
	}
	e := Event{Tags: map[string]bool{"tech": TagOptions{Trend: .2, Duration: Day}, "china": TagOptions{Trend: .2, Duration: Day}}}

	fmt.Println(s, e)
}

func tick() {
	// Get currently Running Events
	// Stop Events that are over the max duration
	// Randomly add new Events to the list of running events that are currently valid (e.g min time between events)
	// Filter Only Currently relevant events
	// Run all events on the stocks
	// Update Database
}

func getAffectedStocks(e Event, s []Stock) (stocks []Stock) {
	for _, stock := range s {
		for tag := range stock.Tags {
			if _, ok := e.Tags[tag]; ok {
				stocks = append(stocks, stock)
				break
			}
		}
	}

	return stocks
}

func computeStockNumbers(stocks []Stock, e Event) number int {
for _, stock := range stocks {
stock.index += 1
print()
}

}

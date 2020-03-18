package game

import "time"

type EventSettings struct {
	EventID     string
	Title       string
	Description string

	// Type says when the event is run
	// can be:
	//	- recurring		- events like ellections that happen every x years
	// 	- fixed				- events that happen on a fixed date
	//	- random			- events that can happen randomly
	Type string

	FixedDate         time.Time // ONLY EVENT TYPE FIXED: date the event has to be run
	RandomChance      float64   // the chance for the event to occur in a tick [0, 1] (thould be less than .01%)
	ReccuringDuration time.Duration

	Dependencies []string // these eventIDs have to have run before this event
	Effects      []string // these eventIDs have to be run after this event is over

	RunBefore time.Time // The event has to be run before this date
	RunAfter  time.Time // The event has to be run after this date

	Duration time.Duration // Time during which the event is the event is run every tick
	Tags     map[string]TagOptions
}

// TagOptions are more indepth settings for specific event tags only
type TagOptions struct {
	AffectsTag     string
	AffectsStockID string
	Trend          int64 // Note: .2 would be 20 and .02 would be 2
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	////Duration time.Duration
}

// StockSettings is the Stock Data Type for storing Stocks ("Class")
type StockSettings struct {
	StockID   string          `json:"stockID"`    // Stock Symbol e.g GOOG
	Name      string          `json:"name"`       // Stock Name e.g Alphabet Inc.
	Index     int64           `json:"startPrice"` // Price per share
	Shares    int64           `json:"stockCount"` // Number per share
	Tags      map[string]bool `json:"tags"`       // A stock can have up to 5 tags
	Stability int64           `json:"stability"`  // Shows how many fluctuations the stock will have
	Trend     int64           `json:"trend"`      // Shows the generall trend of the Stock
}

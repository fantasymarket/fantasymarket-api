package structs

import (
	"fantasymarket/utils/timeutils"
)

// EventSettings is the type for storing information about events
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

	FixedDate         timeutils.Time // ONLY EVENT TYPE FIXED: date the event has to be run
	RandomChance      float64        // the chance for the event to occur in a tick [0, 1] (thould be less than .01%)
	ReccuringDuration timeutils.Duration

	Dependencies []string // these eventIDs have to have run before this event
	Effects      []string // these eventIDs have to be run after this event is over

	RunBefore timeutils.Time // The event has to be run before this date
	RunAfter  timeutils.Time // The event has to be run after this date

	Duration timeutils.Duration // Time during which the event is the event is run every tick
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

package structs

import (
	"fantasymarket/utils/timeutils"
)

// EventSettings is the type for storing information about events
//
//   It supports the following date/duration formats when parsing json/yaml:
//
//   - Dates:
//      Our custom format:
//      2006-01-02 15:04:05
//      YYYY-MM-DD
//
//      RFC3339:
//      2006-01-02T15:04:05Z07:00
//
//   - Durations:
//      ISO8601:
//      P3Y6M4DT12H30M5S
//      => three years, six months, four days, twelve hours, thirty minutes, and five seconds
//
type EventSettings struct {
	// A unique string ID for an event (e.g `quantum-breakthrough`)
	EventID string

	// Used for the news-feed
	Title       string
	Description string

	// Type says when the event is run
	// can be:
	//	- recurring		- events like ellections that happen every x years
	// 	- fixed				- events that happen on a fixed date
	//	- random			- events that can happen randomly
	//
	// Properties like `FixedDate` can only be used when the
	// type is actually set to `fixed`
	Type string

	FixedDate         timeutils.Time     // date the event has to be run
	RandomChance      float64            // the chance for the event to occur in a tick [0, 1] (thould be less than .01%)
	ReccuringDuration timeutils.Duration // When an event has to happen e.g yearly

	// these eventIDs have to have run before this event
	// can be prefixed with `!` if the event should only be run if an event hasn't happened yet
	Dependencies []string

	// Effects are eventIDs that have to be run after this event is over
	Effects []EventEffect

	RunBefore timeutils.Time // The event has to be run before this date
	RunAfter  timeutils.Time // The event has to be run after this date

	Duration timeutils.Duration // Time during which the event is the event is run every tick
	Tags     map[string]TagOptions
}

// EventEffect is a effect that runs after an event is finished
type EventEffect struct {
	// The chance this event is run (0 < chance < 1)
	Chance float64

	// Events that are run directly after the parent event is over
	Effects []string
}

// TagOptions are more indepth settings for specific event tags only
type TagOptions struct {
	AffectsTag     string
	AffectsStockID string
	Trend          int64 // Note: .2 would be 20 and .02 would be 2
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	////Duration time.Duration
}

package events

import (
	"fantasymarket/utils/hash"
	"fantasymarket/utils/timeutils"
	"strconv"
)

// EventDetails is the type for storing information about events
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
type EventDetails struct {
	// A unique string ID for an event (e.g `quantum-breakthrough`)
	EventID string

	// Used for the news-feed
	Title       string
	Description string

	// Type says when the event is run
	// can be:
	// 	- fixed				- events that happen on a fixed date
	//	- recurring		- events like ellections that happen every x years
	//                  note: these also need the same options as a fixed date
	//	- random			- events that can happen randomly
	//  - custom			- event isn't run automatically
	//
	// Properties like `FixedDate` can only be used when the
	// type is actually set to `fixed`
	Type string

	FixedDate             timeutils.Time     // date the event has to be run
	FixedDateRandomOffset timeutils.Duration // Offset from 0-n added to the fixed date

	randomChancePerDay      float64            // the chance for the event to occur in a day [0, 1]
	RecurringDuration timeutils.Duration // When an event has to happen e.g yearly

	// these eventIDs have to have run before this event
	// can be prefixed with `!` if the event should only be run if an event hasn't happened yet
	Dependencies []string

	// Effects are eventIDs that have to be run after this event is over
	Effects []EventEffect

	RunBefore timeutils.Time // The event has to be run before this date
	RunAfter  timeutils.Time // The event has to be run after this date

	Duration timeutils.Duration // Time during which the event is the event is run every tick

	Tags []TagOptions
}

// EventEffect is a effect that runs after an event is finished
type EventEffect struct {
	// The chance this event is run (0 < chance < 1) (decimal percentage)
	Chance float64

	// Event that is run directly after the parent event is over
	EventID string
}

// TagOptions are more indepth settings for specific event tags only
type TagOptions struct {
	// Can be a list of stock tags tags
	// "all" is a shortcut for all stocks
	// You can ignore stocks with `!`, e.g ["all", "!healthcare"]
	AffectsTags   []string
	AffectsStocks []string

	// NOTE: An event can only be affected while an event is active
	// So Offset + Duration cant be larger than EventSettings.Duration
	Offset   timeutils.Duration // The tag only effects the stock after this duration
	Duration timeutils.Duration // The tag only effects the stock for this duration

	// How the much the tags are affected, as a decimal percentage
	Trend float64 // 0.02 => +2% every game-tick

	// If these are set, Trend is ignored and a number in between MinTrend and MaxTrend is chosen
	MinTrend float64
	MaxTrend float64
}

// CalculateTrend calculates the trend for a tag
// If Min and MaxTrend are defined, these take precedent over Trend
func (t *TagOptions) CalculateTrend(ticksSinceStart int64, stockSymbol string) float64 {
	if t.MinTrend == t.MaxTrend {
		if t.MinTrend != 0 {
			return t.MinTrend
		}
		return t.Trend
	}

	seed := stockSymbol + strconv.FormatInt(ticksSinceStart, 10)
	trend := hash.Int64HashRange(int64(t.MinTrend*1000), int64(t.MaxTrend*1000), seed)
	return float64(trend) / 1000
}

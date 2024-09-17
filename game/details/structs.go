package details

import (
	"fantasymarket/utils/hash"
	"fantasymarket/utils/timeutils"
	"strconv"
	"time"
)

// AssetDetails is the type for storing information about assets
type AssetDetails struct {
	Symbol      string   `yaml:"symbol"`      // Asset Symbol e.g GOOG
	Name        string   `yaml:"name"`        // Asset Name e.g Alphabet Inc.
	Type        string   `yaml:"type"`        // Asset Name e.g Alphabet Inc.
	Description string   `yaml:"description"` // Asset Name e.g Alphabet Inc.
	Index       int64    `yaml:"startPrice"`  // Price per share
	Shares      int64    `yaml:"assetCount"`  // Number per share
	Tags        []string `yaml:"tags"`        // A asset can have up to 5 tags
	Stability   float64  `yaml:"stability"`   // Shows how many fluctuations the asset will have
	Trend       float64  `yaml:"trend"`       // Shows the generall trend of the Asset
}

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
	EventID string `yaml:"eventID"`

	// Used for the news-feed
	Title       string `yaml:"title"`
	Description string `yaml:"description"`

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
	Type string `yaml:"type"`

	FixedDate             timeutils.Time     `yaml:"fixedDate"`             // date the event has to be run
	FixedDateRandomOffset timeutils.Duration `yaml:"fixedDateRandomOffset"` // Offset from 0-n added to the fixed date

	RandomChancePerDay float64            `yaml:"randomChancePerDay"` // the chance for the event to occur in a day [0, 1]
	RecurringDuration  timeutils.Duration `yaml:"recurringDuration"`  // When an event has to happen e.g yearly

	// these eventIDs have to have run before this event
	// can be prefixed with `!` if the event should only be run if an event hasn't happened yet
	Dependencies []string `yaml:"dependencies"`

	// Effects are eventIDs that have to be run after this event is over
	Effects []EventEffect `yaml:"events"`

	Duration timeutils.Duration `yaml:"duration"` // Time during which the event is the event is run every tick

	MinTimeBetweenEvents time.Duration  `yaml:"minTimeBetweenEvents"` // The event can only run again after this time has passed
	RunBefore            timeutils.Time `yaml:"runBefore"`            // The event has to be run before this date
	RunAfter             timeutils.Time `yaml:"runAfter"`             // The event has to be run after this date

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
	// Can be a list of asset tags tags
	// "all" is a shortcut for all assets
	// You can ignore assets with `!`, e.g ["all", "!healthcare"]
	AffectsTags   []string
	AffectsAssets []string

	// NOTE: An event can only be affected while an event is active
	// So Offset + Duration cant be larger than EventSettings.Duration
	Offset   timeutils.Duration // The tag only effects the asset after this duration
	Duration timeutils.Duration // The tag only effects the asset for this duration

	// How the much the tags are affected, as a decimal percentage
	Trend float64 // 0.02 => +2% every game-tick

	// If these are set, Trend is ignored and a number in between MinTrend and MaxTrend is chosen
	MinTrend float64
	MaxTrend float64
}

// CalculateTrend calculates the trend for a tag
// If Min and MaxTrend are defined, these take precedent over Trend
func (t *TagOptions) CalculateTrend(ticksSinceStart int64, assetSymbol string) float64 {
	if t.MinTrend == t.MaxTrend {
		if t.MinTrend != 0 {
			return t.MinTrend
		}
		return t.Trend
	}

	seed := assetSymbol + strconv.FormatInt(ticksSinceStart, 10)
	trend := hash.Int64HashRange(int64(t.MinTrend*1000), int64(t.MaxTrend*1000), seed)
	return float64(trend) / 1000
}

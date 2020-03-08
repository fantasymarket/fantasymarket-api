package game

// TagOptions are more indepth settings for specific event tags only
type TagOptions struct {
	Affects string
	Trend   int64 // Note: .2 would be 20 and .02 would be 2
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	////Duration time.Duration
}

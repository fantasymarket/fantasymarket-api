package models

import (
	"time"
) 

// Events happen randomly every game tick
type Event struct {
	EventID string

	Title string
	Text  string

	Duration  time.Duration // Time during which the event is the event is run every tick
	CreatedAt time.Time

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
}

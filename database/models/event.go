package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Events happen randomly every game tick
type Event struct {
	ID      string // A Unique ID for every event (since the same event might happen multiple times)
	EventID string

	Title string
	Text  string

	Duration  time.Duration // Time during which the event is the event is run every tick
	CreatedAt time.Time

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
}

func (e Event) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

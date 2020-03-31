package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Event are Events that happen randomly every game tick
type Event struct {
	ID      uuid.UUID `gorm:"primary_key"` // A Unique ID for every event (since the same event might happen multiple times)
	EventID string
	Active  bool

	Title string
	Text  string

	Duration  time.Duration // Time during which the event is the event is run every tick
	CreatedAt time.Time

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
}

// BeforeCreate runs before an event is created in the database
func (e Event) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User created Orders
type Order struct {
	ID         string `gorm:"primary_key"` // A Unique ID for every order (since the same event might happen multiple times)
	CreatedAt  time.Time
	CanceledAt time.Time
	Type       string
	Side       string
	Symbol     string
	status     string

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
}

func (o Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}

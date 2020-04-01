package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Order is a Order Struct
type Order struct {
	OrderID uuid.UUID `gorm:"primary_key"` // A Unique ID for every order (since the same event might happen multiple times)

	UserID uuid.UUID
	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:UserID"`

	CreatedAt  time.Time
	CanceledAt time.Time
	Type       string
	Side       string
	Symbol     string
	status     string

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
}

// BeforeCreate runs before a order is created in the database
func (o Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("OrderID", uuid.NewV4())
	return nil
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User is the User "Class"
type User struct {
	UserID    uuid.UUID `gorm:"primary_key"` // A Unique ID for every stock data point (since theres a new entry for each stock ID every tick)
	CreatedAt time.Time

	Portfolio Portfolio `gorm:"foreignkey:UserID;association_foreignkey:UserID"`

	// Stock Name e.g Alphabet Inc.
	Username string
}

// BeforeCreate runs before a user is created in the database
func (u User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("userID", uuid.NewV4())
	return nil
}

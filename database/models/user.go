package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User is the User "Class"
type User struct {
	UserID    uuid.UUID `gorm:"primary_key"` // A Unique ID for every stock data point (since theres a new entry for each stock ID every tick)
	CreatedAt time.Time

	Portfolio Portfolio `gorm:"foreignkey:UserID;association_foreignkey:UserID"`

	// Stock Name e.g Alphabet Inc.
	Username string `gorm:"not null;unique" valid:"required;stringlength(3"`
	Password string
}

// BeforeSave runs before User is saved to the database
func (user *User) BeforeSave(scope *gorm.Scope) error {
	// Validate all fields
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return errors.New("validation failed")
	}

	return nil
}

// BeforeCreate runs before a user is created in the database
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("userID", uuid.NewV4())
	return nil
}

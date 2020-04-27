package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Order is a Order Struct
type Order struct {
	OrderID uuid.UUID `gorm:"primary_key"`

	UserID    uuid.UUID
	User      User `gorm:"ForeignKey:UserID;AssociationForeignKey:UserID"`
	CreatedAt time.Time

	Type   string `json:"type"`// the type of the asset (so only stock for now)
	Symbol string `json:"symbol"`// the symbol of the stock
	Side   string // buy or sell
	Status string // filled, waiting or canceled
	Amount int64  // amount of stocks in the order

	CanceledAt time.Time
	FilledAt   time.Time
}

// BeforeCreate runs before a order is created in the database
func (o Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("OrderID", uuid.NewV4())
	return nil
}

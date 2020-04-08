package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Stock is the Stock "Class"
type Stock struct {
	StockID   uuid.UUID `gorm:"primary_key"` // A Unique ID for every stock data point (since theres a new entry for each stock ID every tick)
	CreatedAt time.Time
	Tick      int64

	// Stock Symbol e.g GOOG
	Symbol string

	// Stock Name e.g Alphabet Inc.
	Name string

	// Price per share
	Index int64

	// Volume since last tick, we'll have to invent this
	//    Calculated based on
	// 		- the change of the index from the last tick,
	// 		- total index (so expensive stocks have larger volume than cheaper ones)
	// 		- random fluctuation
	Volume int64
}

// BeforeCreate runs before a stock is created in the database
func (s Stock) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("StockID", uuid.NewV4())
	return nil
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Stock is the Stock "Class"
type Stock struct {
	ID        string  `gorm:"primary_key"` // A Unique ID for every stock data point (since theres a new entry for each stock ID every tick)
	CreatedAt time.Time

	// Stock Symbol e.g GOOG
	StockID string

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

func (s Stock) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}

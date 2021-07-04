package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Asset is the Asset "Class"
type Asset struct {
	AssetID   uuid.UUID `gorm:"primary_key"` // A Unique ID for every asset data point (since theres a new entry for each asset ID every tick)
	CreatedAt time.Time
	Tick      int64

	// Asset Symbol e.g GOOG
	Symbol string

	// Asset Name e.g Alphabet Inc.
	Name string

	// Price per share
	Index int64

	// Volume since last tick, we'll have to invent this
	//    Calculated based on
	// 		- the change of the index from the last tick,
	// 		- total index (so expensive assets have larger volume than cheaper ones)
	// 		- random fluctuation
	Volume int64
}

// BeforeCreate runs before a asset is created in the database
func (s Asset) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("AssetID", uuid.NewV4())
	return nil
}

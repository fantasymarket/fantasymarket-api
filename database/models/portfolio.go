package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Portfolio is the Portfolio "Class"
type Portfolio struct {
	PortfolioID uuid.UUID `gorm:"primary_key"`
	UserID      uuid.UUID

	Balance int64
	Items   []PortfolioItem
}

// BeforeCreate runs before a portfolio is created in the database
func (p Portfolio) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("PortfolioID", uuid.NewV4())
	return nil
}

// PortfolioItem tracks an item (like a stock) in a specific portfolio
type PortfolioItem struct {
	PortfolioItemID uuid.UUID `gorm:"primary_key"`
	PortfolioID     uuid.UUID `gorm:"not null;unique"`

	Type   string // only stock for now
	Symbol string
	Amount int64
}

// BeforeCreate runs before a portfolioItem is created in the database
func (p PortfolioItem) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("PortfolioItemID", uuid.NewV4())
	return nil
}

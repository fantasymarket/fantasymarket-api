package database

import (
	"fmt"

	// "fantasymarket/game"

	"fantasymarket/database/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// StockSettings TEST STRUCT - DELETE
type StockSettings struct {
	StockID   string          // Stock Symbol e.g GOOG
	Name      string          // Stock Name e.g Alphabet Inc.
	Index     int64           // Price per share
	Shares    int64           // Number per share
	Tags      map[string]bool // A stock can have up to 5 tags
	Stability int64           // Shows how many fluctuations the stock will have
	Trend     int64           // Shows the generall trend of the Stock
	Volume    int64
}

// Connect connects to the database and returns thedatabase object
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "database.db")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&models.Stock{})
	db.AutoMigrate(&models.Event{})

	// db.Create(&models.Stock{StockID: "GOOG"})

	// hier steht alles wie man daten kriegt http://gorm.io/docs/query.html

	var stock models.Stock
	db.First(&stock, "stock_id = ?", "GOOG") // find product with code l1212

	return db, nil
}

// AddStockToTable takes the stock as input and adds it to the StockDB
func AddStockToTable(db *gorm.DB, stock models.Stock) {
	db.Create(models.Stock{
		StockID: Stock.StockID,
		Name: Stock.Name
		Index: Stock.Index
		Volume: Stock.Volume
	})
}

// GetStocks returns all the stocks in the DB as a list
func GetStocks(db *gorm.DB) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := db.Find(&stocks).Error; err != nil {
		return nil, err
	}

	stocks = db.Find(&stock)
	return stocks
}

// USE DBName;
// GO
// DECLARE @MyMsg VARCHAR(50)
// SELECT @MyMsg = 'Hello, World.'
// GO -- @MyMsg is not valid after this GO ends the batch.
// https://sqliteonline.com

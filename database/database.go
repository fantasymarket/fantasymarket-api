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

var database *gorm.DB

// Connect connects to the database and returns thedatabase object
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "database.db")
	//ja eben. jetzt gehts ja instant zu dann
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&models.Stock{})
	db.AutoMigrate(&models.Event{})

	// db.Exec("DROP TABLE users;")

	// db.Create(&models.Stock{StockID: "GOOG"})

	// hier steht alles wie man daten kriegt http://gorm.io/docs/query.html
	// so kann man shit lesen
	// Was mach ich wenn ich alles haben will? - .Find
	var stock models.Stock
	db.First(&stock, "stock_id = ?", "GOOG") // find product with code l1212

	return db, nil

	//addStocksToDB("GOOG", "GOOGLE", 10000, 2)
	//getDatabaseData()
}

// USE DBName;
// GO
// DECLARE @MyMsg VARCHAR(50)
// SELECT @MyMsg = 'Hello, World.'
// GO -- @MyMsg is not valid after this GO ends the batch.
// https://sqliteonline.com

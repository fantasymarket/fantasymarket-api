package database

import (
	"fmt"

	"fantasymarket/database/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DatabaseService struct {
	DB *gorm.DB // gorm database instance
}

// Connect connects to the database and returns thedatabase object
func Connect() (*DatabaseService, error) {
	db, err := gorm.Open("sqlite3", "database.db")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&models.Stock{})
	db.AutoMigrate(&models.Event{})

	fmt.Println("connected to da database my doods D:")

	return &DatabaseService{
		DB: db,
	}, nil
}

func (s *DatabaseService) CreateStockForTest(stockID string, name string, index int64, volume int64) models.Stock {
	stock := models.Stock{StockID: stockID, Name: name, Index: index, Volume: volume, Tick: 0}
	return stock
}

func (s *DatabaseService) AddStockToTable(stock models.Stock, tick int64) error {
	return s.DB.Create(&models.Stock{
		StockID: stock.StockID,
		Name:    stock.Name,
		Index:   stock.Index,
		Volume:  stock.Volume,
		Tick:    tick,
	}).Error
}

func (s *DatabaseService) GetEvents() ([]models.Event, error) {
	var events []models.Event
	if err := s.DB.Where(models.Event{Active: true}).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (s *DatabaseService) RemoveEvent(eventID string) error {
	return s.DB.Where(models.Event{Active: true, EventID: eventID}).Update("active", false).Error
}

func (s *DatabaseService) GetStocksAtTick(lastTick int64) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := s.DB.Where(models.Stock{Tick: lastTick}).Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

// USE DBName;
// GO
// DECLARE @MyMsg VARCHAR(50)
// SELECT @MyMsg = 'Hello, World.'
// GO -- @MyMsg is not valid after this GO ends the batch.
// https://sqliteonline.com

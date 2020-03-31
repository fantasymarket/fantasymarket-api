package database

import (
	"fmt"

	"fantasymarket/database/models"
	gameStructs "fantasymarket/game/structs"

	"github.com/jinzhu/gorm"

	// load sqlite dialect
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

	db.AutoMigrate(
		&models.Stock{},
		&models.Event{},
		&models.Order{},
	)
	fmt.Println("connected to da database my doods D:")

	return &DatabaseService{
		DB: db,
	}, nil
}

func (s *DatabaseService) CreateInitialStocks(stocks map[string]gameStructs.StockSettings) error {

	for _, stock := range stocks {
		if err := s.DB.FirstOrCreate(
			&models.Stock{},
			&models.Stock{
				Symbol: stock.Symbol,
				Index:  stock.Index,
				Name:   stock.Name,
				Tick:   0,
				Volume: 0,
			},
		).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *DatabaseService) AddStockToTable(stock models.Stock, tick int64) error {
	return s.DB.Create(&models.Stock{
		Symbol: stock.Symbol,
		Name:   stock.Name,
		Index:  stock.Index,
		Volume: stock.Volume,
		Tick:   tick,
	}).Error
}

func (s *DatabaseService) GetEvents() ([]models.Event, error) {
	var events []models.Event
	if err := s.DB.Where(models.Event{Active: true}).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// RemoveEvent marks an event as inactive so it won't affect stocks in the GameLoop anymore
func (s *DatabaseService) RemoveEvent(eventID string) error {
	return s.DB.Where(models.Event{Active: true, EventID: eventID}).Update("active", false).Error
}

// GetNextTick retrieves the tick number for the next tick from the database,
// this is used to initialize our GameService when the program restarts
func (s *DatabaseService) GetNextTick() (int64, error) {
	var stock models.Stock
	if err := s.DB.Table("stocks").Select("tick").Order("tick desc").First(&stock).Error; err != nil {
		return 0, err
	}

	fmt.Println("Next Tick: ", stock.Tick+1)
	return stock.Tick + 1, nil
}

func (s *DatabaseService) GetStocksAtTick(lastTick int64) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := s.DB.Where(models.Stock{Tick: lastTick}).Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

//func (s *DatabaseService) AddOrder(order map[string]string) error {

//}

//func (s *DatabaseService) GetOrder(id int) error {

//}

//func (s *DatabaseService) DeleteOrder(id int) error {

//}

// USE DBName;
// GO
// DECLARE @MyMsg VARCHAR(50)
// SELECT @MyMsg = 'Hello, World.'
// GO -- @MyMsg is not valid after this GO ends the batch.
// https://sqliteonline.com

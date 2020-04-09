package database

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"fantasymarket/database/models"
	"fantasymarket/game/stocks"
	"fantasymarket/utils/config"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	// load sqlite dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Service is the Database Service
type Service struct {
	DB     *gorm.DB // gorm database instance
	Config *config.Config
}

// Connect connects to the database and returns thedatabase object
func Connect(config *config.Config) (*Service, error) {
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

	log.Info().Msg("successfully connected to the database")

	return &Service{
		DB:     db,
		Config: config,
	}, nil
}

// CreateInitialStocks takes a list of initial stocks and uses them to initialize the database
func (s *Service) CreateInitialStocks(stockDetails map[string]stocks.StockDetails) error {

	for _, stock := range stockDetails {
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

// AddStock adds a stock to the stock table
func (s *Service) AddStock(stock models.Stock, tick int64) error {
	return s.DB.Create(&models.Stock{
		Symbol: stock.Symbol,
		Name:   stock.Name,
		Index:  stock.Index,
		Volume: stock.Volume,
		Tick:   tick,
	}).Error
}

// AddStocks adds a slice of stocks to the stock table
func (s *Service) AddStocks(stocks []models.Stock, tick int64) error {
	for _, stock := range stocks {
		if err := s.AddStock(stock, tick); err != nil {
			return err
		}
	}
	return nil
}

// GetEvents fetches all currently active events
func (s *Service) GetEvents() ([]models.Event, error) {
	var events []models.Event
	// TODO: createdAt > currentGameTime
	if err := s.DB.Where(models.Event{Active: true}).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// RemoveEvent marks an event as inactive so it won't affect stocks in the GameLoop anymore
func (s *Service) RemoveEvent(uniqueEventID uuid.UUID) error {
	return s.DB.Where(models.Event{Active: true, ID: uniqueEventID}).Update("active", false).Error
}

// GetNextTick retrieves the tick number for the next tick from the database,
// this is used to initialize our Service when the program restarts
func (s *Service) GetNextTick() (int64, error) {
	var stock models.Stock
	if err := s.DB.Table("stocks").Select("tick").Order("tick desc").First(&stock).Error; err != nil {
		return 0, err
	}

	return stock.Tick + 1, nil
}

// GetStocksAtTick fetches the value of all stocks at a specific tick
func (s *Service) GetStocksAtTick(lastTick int64) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := s.DB.Where(models.Stock{Tick: lastTick}).Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

//func (s *Service) AddOrder(order map[string]string) error {

//}

//func (s *Service) GetOrder(id int) error {

//}

//func (s *Service) DeleteOrder(id int) error {

//}

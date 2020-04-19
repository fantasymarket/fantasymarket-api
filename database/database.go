package database

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"fantasymarket/database/models"
	"fantasymarket/game/events"
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
		&models.User{},
		&models.Portfolio{},
		&models.PortfolioItem{},
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
func (s *Service) GetEvents(currentDate time.Time) ([]models.Event, error) {
	var events []models.Event

	if err := s.DB.Where(models.Event{
		Active: true,
	}).Where("created_at > ?", currentDate).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// AddEvent adds an event to the event table
func (s *Service) AddEvent(event events.EventDetails, createdAt time.Time) error {
	return s.DB.Create(&models.Event{
		EventID:   event.EventID,
		Title:     event.Title,
		Text:      event.Description,
		Active:    true,
		CreatedAt: createdAt,
	}).Error
}

// RemoveEvent marks an event as inactive so it won't affect stocks in the GameLoop anymore
func (s *Service) RemoveEvent(uniqueEventID uuid.UUID) error {
	return s.DB.Where(models.Event{Active: true, ID: uniqueEventID}).Update("active", false).Error
}

func (s *Service) GetEventHistory() (map[string][]time.Time, error) {
	eventHistory := map[string][]time.Time{}

	var events []models.Event
	if err := s.DB.Find(&events).Error; err != nil {
		return nil, err
	}

	for _, event := range events {
		eventID := event.EventID
		createdAt := event.CreatedAt

		if _, exists := eventHistory[eventID]; !exists {
			eventHistory[eventID] = []time.Time{}
		}
		eventHistory[eventID] = append(eventHistory[eventID], createdAt)

	}

	return eventHistory, nil
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

func (s *Service) AddOrder(order models.Order, userID uuid.UUID, currentDate time.Time) error {
	return s.DB.Create(&models.Order{
		UserID:    userID,
		CreatedAt: currentDate,
		Type:      order.Type,
		Side:      order.Side,
		Symbol:    order.Symbol,
		Status:    order.Status,
	}).Error
}

//func (s *Service) GetOrderForUser(userID uuid.UUID) error {

//}

//func (s *Service) GetOrderForUserByID(orderID uuid.UUID, userID uuid.UUID) {


//}

func (s *Service) CancelOrder(orderID uuid.UUID, currentDate time.Time) error {

	var order models.Order
	if err := s.DB.Where(models.Order{OrderID: orderID}).First(&order).Error; err != nil {
		return err
	}

	// TODO check if the order is still active

	return s.DB.Model(&order).Updates(models.Order{Status: "cancelled", FilledAt: currentDate}).Error
}

func (s *Service) FillOrder(orderID uuid.UUID, userID uuid.UUID, currentDate time.Time) error {

	var order models.Order
	// TODO: update users portfolio
	if err := s.DB.Where(models.Order{OrderID: orderID}).Find(&order).Error; err != nil {
		return err
	}

	var user models.User
	if err := s.DB.Where(models.User{UserID: userID}).Preload("Portfolio.Items").Find(&user).Error; err != nil {
		return err
	}

	// user.Portfolio is the portfolio
	// user.Portfolio.Items are all items

	// create new PortfolioItem if it doesn't exist yet
	// update the amount

	// if order.Type == "stock" {
	// 	if Portfolio.Items.contains(order.Symbol) {
	// 		Portfolio.Items.amount += order.Amount
	// 	}
	// 	else {
	// 		Portfolio.Items = Append(Portfolio.items, new PortfolioItem{Type: order.Type, Symbol: order.Symbol, Amount: order.Amount})
	// 	}
	// }

	// TODO: update users portfolio balance

	return s.DB.Where(models.Order{OrderID: orderID}).Updates(models.Order{Status: "filled", FilledAt: currentDate}).Error
}


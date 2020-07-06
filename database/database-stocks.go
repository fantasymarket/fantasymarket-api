package database

import (
	"fantasymarket/database/models"
	"fantasymarket/game/details"

	uuid "github.com/satori/go.uuid"
)

// CreateInitialStocks takes a list of initial stocks and uses them to initialize the database
func (s *Service) CreateInitialStocks(stockDetails map[string]details.StockDetails) error {

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

// GetStockMapAtTick fetches the value of all stocks at a tick as a Map
func (s *Service) GetStockMapAtTick(lastTick int64) (map[string]models.Stock, error) {
	stocks, err := s.GetStocksAtTick(lastTick)
	if err != nil {
		return nil, err
	}

	stockMap := map[string]models.Stock{}
	for _, stock := range stocks {
		stockMap[stock.Symbol] = stock
	}

	return stockMap, nil
}

// GetStockAtTick fetches the value of a specific stock at a specific tick
func (s *Service) GetStockAtTick(stockID uuid.UUID, lastTick int64) (*models.Stock, error) {
	var stock models.Stock
	if err := s.DB.Where(models.Stock{Tick: lastTick, StockID: stockID}).First(&stock).Error; err != nil {
		return nil, err
	}

	return &stock, nil
}

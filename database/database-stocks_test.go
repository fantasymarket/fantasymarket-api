package database_test

import (
	"fantasymarket/database/models"
	"fantasymarket/game/details"

	"github.com/stretchr/testify/assert"
)

type CreateInitialStocksTestData struct {
	stock       map[string]details.StockDetails
	expectation models.Stock
}

var testCreateInitialStocksData = []CreateInitialStocksTestData{
	{
		stock: map[string]details.StockDetails{"HELLO": details.StockDetails{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Stock",
		},
		},
		expectation: models.Stock{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Stock",
		},
	},
	{
		stock: map[string]details.StockDetails{"": details.StockDetails{
			Symbol: "",
			Index:  401,
			Name:   "Not Hello Stock",
		},
		},
		expectation: models.Stock{
			Symbol: "",
			Index:  401,
			Name:   "Not Hello Stock",
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestCreateInitialStocks() {

	stocks := []models.Stock{}
	for i, test := range testCreateInitialStocksData {
		err := suite.dbService.CreateInitialStocks(test.stock)
		assert.Equal(suite.T(), nil, err)
		err = suite.dbService.DB.Find(&stocks).Error
		assert.Equal(suite.T(), nil, err)

		if test.expectation.Symbol != "" {
			//Again.., I hate it
			assert.Equal(suite.T(), test.expectation.Symbol, stocks[i].Symbol)
			assert.Equal(suite.T(), test.expectation.Index, stocks[i].Index)
			assert.Equal(suite.T(), test.expectation.Name, stocks[i].Name)
		}
	}

	suite.dbService.DB.Close()
}

type AddStockTestData struct {
	stock       models.Stock
	expectation models.Stock
}

var testAddStockData = []AddStockTestData{
	{
		stock: models.Stock{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Stock",
		},
		expectation: models.Stock{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Stock",
		},
	},
	{
		stock: models.Stock{
			Symbol: "",
			Index:  100,
			Name:   "Hello Stock",
		},
		expectation: models.Stock{
			Symbol: "",
			Index:  100,
			Name:   "Hello Stock",
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestAddStock() {
	for _, test := range testAddStockData {
		err := suite.dbService.AddStock(test.stock, 1)
		assert.Equal(suite.T(), nil, err)
		assert.Equal(suite.T(), false, suite.dbService.DB.Where("symbol = ?", test.stock.Symbol).Find(&models.Stock{}).RecordNotFound())
	}
}

// AddStock adds a stock to the stock table
// func (s *Service) AddStock(stock models.Stock, tick int64) error {
// 	return s.DB.Create(&models.Stock{
// 		Symbol: stock.Symbol,
// 		Name:   stock.Name,
// 		Index:  stock.Index,
// 		Volume: stock.Volume,
// 		Tick:   tick,
// 	}).Error
// }

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
		stock: map[string]details.StockDetails{"HELLO": {
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
		stock: map[string]details.StockDetails{"": {
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
	suite.dbService.DB.Close()
}

type GetNextTickTestData struct {
	stock       models.Stock
	expectation int64
}

var testGetNextTickData = []GetNextTickTestData{
	{
		stock: models.Stock{
			Tick: 0,
		},
		expectation: 1,
	},
	{
		stock:       models.Stock{},
		expectation: 1,
	},
	{
		expectation: 1,
	},
}

func (suite *DatabaseTestSuite) TestGetNextTick() {
	for _, test := range testGetNextTickData {
		err := suite.dbService.DB.Create(&test.stock).Error
		assert.Equal(suite.T(), nil, err)
		result, err := suite.dbService.GetNextTick()
		assert.Equal(suite.T(), nil, err)
		assert.Equal(suite.T(), test.expectation, result)
	}
	suite.dbService.DB.Close()
}

type GetStocksAtTickTestData struct {
	tick        int64
	stock       models.Stock
	expectation []models.Stock
}

var testGetStocksAtTickData = []GetStocksAtTickTestData{
	{
		tick: 1,
		stock: models.Stock{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Stock",
			Tick:   1,
		},
		expectation: []models.Stock{
			{
				Symbol: "HELLO",
				Index:  100,
				Name:   "Hello Stock",
				Tick:   1,
			},
		},
	},
	{
		tick: 2,
		stock: models.Stock{
			Symbol: "NOTHEL",
			Index:  100,
			Name:   "Not Hello Stock",
			Tick:   2,
		},
		expectation: []models.Stock{
			{
				Symbol: "NOTHEL",
				Index:  100,
				Name:   "Not Hello Stock",
				Tick:   2,
			},
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestGetStockAtTick() {
	for _, test := range testGetStocksAtTickData {
		err := suite.dbService.DB.Create(&test.stock).Error
		assert.Equal(suite.T(), nil, err)
		result, err := suite.dbService.GetStocksAtTick(test.tick)
		assert.Equal(suite.T(), nil, err)

		for j := 0; j < len(test.expectation); j++ {
			//I need to find a way to fix this
			assert.Equal(suite.T(), test.expectation[j].Symbol, result[j].Symbol)
			assert.Equal(suite.T(), test.expectation[j].Index, result[j].Index)
			assert.Equal(suite.T(), test.expectation[j].Name, result[j].Name)
			assert.Equal(suite.T(), test.expectation[j].Tick, result[j].Tick)
		}
	}
	suite.dbService.DB.Close()
}

// // GetStocksAtTick fetches the value of all stocks at a specific tick
// func (s *Service) GetStocksAtTick(lastTick int64) ([]models.Stock, error) {
// 	var stocks []models.Stock
// 	if err := s.DB.Where(models.Stock{Tick: lastTick}).Find(&stocks).Error; err != nil {
// 		return nil, err
// 	}

// 	return stocks, nil
// }

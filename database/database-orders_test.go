package database_test

import (
	"errors"
	"fantasymarket/database/models"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type AddOrderTestData struct {
	input  models.Order
	expect models.Order
}

var testAddOrderData = []AddOrderTestData{
	{
		input: models.Order{
			Type:   "stock",
			Symbol: "APPL",
			Status: "waiting",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "APPL",
			Status: "waiting",
		},
	},
	{
		input: models.Order{
			Type:   "stock",
			Symbol: "GOOG",
			Side:   "buy",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "GOOG",
			Side:   "buy",
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestAddOrder() {
	userID := uuid.NewV4()
	currentDate := parseTime("2019-12-30T15:00:05Z")
	for _, test := range testAddOrderData {
		result := models.Order{}
		test.input.UserID = userID
		test.input.CreatedAt = currentDate
		err := suite.dbService.AddOrder(test.input, test.input.UserID, currentDate)
		if test.input.OrderID == uuid.Nil {
			assert.Equal(suite.T(), errors.New("you cant add an empty order"), err)
		} else {
			assert.Equal(suite.T(), nil, err)

			err = suite.dbService.DB.Where(test.input).First(&result).Error
			assert.Equal(suite.T(), nil, err)

			assert.Equal(suite.T(), test.expect.Type, result.Type)
			assert.Equal(suite.T(), test.expect.Symbol, result.Symbol)
			assert.Equal(suite.T(), test.expect.Side, result.Side)
		}
	}

	suite.dbService.DB.Close()
}

type GetOrderTestData struct {
	orderDetails models.Order
	expect       []models.Order
	limit        int
	offset       int
}

var testGetOrderData = []GetOrderTestData{
	{
		orderDetails: models.Order{
			Type:   "stock",
			Symbol: "GOOG",
		},
		expect: []models.Order{
			{
				Type:   "stock",
				Symbol: "GOOG",
				Price:  100,
			},
		},
		limit:  1,
		offset: -1,
	},
	{
		orderDetails: models.Order{
			Type: "stock",
		},
		expect: []models.Order{
			{
				Type:   "stock",
				Symbol: "GOOG",
				Price:  100,
			},
			{
				Type:   "stock",
				Symbol: "APPL",
				Price:  50,
			},
		},
		limit:  2,
		offset: -1,
	},
	{
		orderDetails: models.Order{
			Symbol: "GOOG",
		},
		expect: []models.Order{
			{
				Type:   "stock",
				Symbol: "GOOG",
				Price:  100,
			},
			{
				Type:   "stock",
				Symbol: "GOOG",
				Price:  200,
			},
		},
		limit:  -1,
		offset: -1,
	},
	{
		orderDetails: models.Order{
			Type: "stock",
		},
		expect: []models.Order{
			{
				Type:   "stock",
				Symbol: "AMZN",
				Price:  5000,
			},
			{
				Type:   "stock",
				Symbol: "GOOG",
				Price:  200,
			},
		},
		limit:  2,
		offset: 2,
	},
	{
		orderDetails: models.Order{
			Type: "commoditites",
		},
		expect: []models.Order{
			{
				Type:   "commoditites",
				Symbol: "GOLD",
				Price:  100,
			},
		},
		limit:  1,
		offset: -1,
	},
}
var initialOrdersInDB = []models.Order{
	{
		Type:   "stock",
		Symbol: "GOOG",
		Price:  100,
	},
	{
		Type:   "commoditites",
		Symbol: "GOLD",
		Price:  100,
	},
	{
		Type:   "stock",
		Symbol: "APPL",
		Price:  50,
	},
	{
		Type:   "commoditites",
		Symbol: "SILV",
		Price:  1,
	},
	{
		Type:   "stock",
		Symbol: "AMZN",
		Price:  5000,
	},
	{
		Type:   "stock",
		Symbol: "GOOG",
		Price:  200,
	},
}

func (suite *DatabaseTestSuite) TestGetOrder() {
	userID := uuid.NewV4()
	for i, order := range initialOrdersInDB {
		if i%2 == 0 {
			order.UserID = userID
		}
		err := suite.dbService.DB.Create(&order).Error
		assert.Equal(suite.T(), nil, err)
	}
	for index, test := range testGetOrderData {
		if index == 5 {
			test.orderDetails.UserID = userID
		}
		result, err := suite.dbService.GetOrders(test.orderDetails, test.limit, test.offset)
		assert.Equal(suite.T(), nil, err)

		for i, r := range result {
			assert.Equal(suite.T(), test.expect[i].Type, r.Type)
			assert.Equal(suite.T(), test.expect[i].Symbol, r.Symbol)
			if index == 5 {
				assert.Equal(suite.T(), test.expect[i].UserID, r.UserID)
			}
		}
	}

	suite.dbService.DB.Close()
}

type GetOrderByIDTestData struct {
	input  models.Order
	expect models.Order
}

var testGetOrderByIDData = []GetOrderByIDTestData{
	{
		input: models.Order{
			Type:   "stock",
			Symbol: "AMZN",
			Price:  5000,
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "AMZN",
			Price:  5000,
		},
	},
	{
		input: models.Order{
			Type:   "commodities",
			Symbol: "GOLD",
			Price:  100,
		},
		expect: models.Order{
			Type:   "commodities",
			Symbol: "GOLD",
			Price:  100,
		},
	},
}

func (suite *DatabaseTestSuite) TestGetorderByID() {
	for _, test := range testGetOrderByIDData {
		err := suite.dbService.DB.Create(&test.input).Error
		assert.Equal(suite.T(), nil, err)

		result, err := suite.dbService.GetOrderByID(test.input.OrderID)
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), test.input.OrderID, result.OrderID)
	}

	suite.dbService.DB.Close()
}

type CancelOrderTestData struct {
	input  models.Order
	expect string
}

var testCancelOrderData = []CancelOrderTestData{
	{
		input: models.Order{
			Symbol: "ASDF",
			Status: "waiting",
		},
		expect: "cancelled",
	},
	{
		input: models.Order{
			Symbol: "YES",
			Status: "waiting",
		},
		expect: "cancelled",
	},
	{
		input: models.Order{
			Symbol: "KMS",
			Status: "filled",
		},
		expect: "cancelled",
	},
	{
		input: models.Order{
			Symbol: "WHAT",
			Status: "filled",
		},
		expect: "cancelled",
	},
}

func (suite *DatabaseTestSuite) TestCancelOrder() {
	var ErrOrderFilledOrCancelled = errors.New("can't cancel order, as its already filled or cancelled")
	currentDate := parseTime("2020-08-05T10:58:00Z")
	for _, test := range testCancelOrderData {
		err := suite.dbService.DB.Create(&test.input).Error
		assert.Equal(suite.T(), nil, err)

		err = suite.dbService.CancelOrder(test.input.OrderID, currentDate)

		if test.input.Status == "filled" || test.input.Status == "cancelled" {
			assert.Equal(suite.T(), ErrOrderFilledOrCancelled, err)
		} else {
			assert.Equal(suite.T(), nil, err)
			err = suite.dbService.DB.Where(models.Order{OrderID: test.input.OrderID}).First(&test.input).Error
			assert.Equal(suite.T(), nil, err)

			assert.Equal(suite.T(), test.expect, test.input.Status)
		}
	}

	suite.dbService.DB.Close()
}

type UpdateOrderTestData struct {
	input   models.Order
	updated models.Order
}

var testUpdateOrderData = []UpdateOrderTestData{
	{
		input: models.Order{
			Symbol: "TSLA",
			Status: "waiting",
		},
		updated: models.Order{
			Symbol: "TSLA",
			Status: "filled",
		},
	},
}

func (suite *DatabaseTestSuite) TestUpdateOrder() {

	for _, test := range testUpdateOrderData {
		var result models.Order
		err := suite.dbService.DB.Create(&test.input).Error
		assert.Equal(suite.T(), nil, err)

		err = suite.dbService.UpdateOrder(test.input.OrderID, test.updated)
		assert.Equal(suite.T(), nil, err)

		err = suite.dbService.DB.Where(models.Order{OrderID: test.input.OrderID}).First(&result).Error
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), test.updated.Symbol, result.Symbol)
		assert.Equal(suite.T(), test.updated.Status, result.Status)
	}
}

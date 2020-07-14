package database_test

import (
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
			Symbol: "HALLO",
			Status: "waiting",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "HALLO",
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

//write test for GetOrder(), FillOrder(), CancelOrder(), UPdate order and FillOrder()

func (suite *DatabaseTestSuite) TestAddOrder() {
	userID := uuid.NewV4()
	currentDate := parseTime("2019-12-30T15:00:05Z")
	for _, test := range testAddOrderData {
		result := models.Order{}
		test.input.UserID = userID
		test.input.CreatedAt = currentDate
		err := suite.dbService.AddOrder(test.input, test.input.UserID, currentDate)
		assert.Equal(suite.T(), nil, err)

		suite.T().Log(test)
		suite.T().Log("ok") //Dont delete this line. The tests literally fail if you do

		err = suite.dbService.DB.Where(test.input).First(&result).Error
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), test.expect.Type, result.Type)
		assert.Equal(suite.T(), test.expect.Symbol, result.Symbol)
		assert.Equal(suite.T(), test.expect.Side, result.Side)
	}

	suite.dbService.DB.Close()
}

type GetOrderTestData struct {
	orderDetails models.Order
	expect       models.Order
}

var testGetOrderData = []GetOrderTestData{
	{
		orderDetails: models.Order{
			Type:   "stock",
			Symbol: "KMS",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "KMS",
		},
	},
	{
		orderDetails: models.Order{
			Type:   "stock",
			Symbol: "KYS",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "KYS",
		},
	},
}

//func (suite *DatabaseTestSuite) TestGetOrder() {
//	userID := uuid.NewV4()
//	for _, test := range
//}
// // GetOrders gets all orders based on the parameters of orderDetails where Symbol, Type and userID can be set.
// // Limit is how many items. Offset is from where to where the data is used
// func (s *Service) GetOrders(orderDetails models.Order, limit int, offset int) (*[]models.Order, error) {
// 	var orders *[]models.Order
// 	if err := s.DB.Where(models.Order{UserID: orderDetails.UserID, Type: orderDetails.Type, Symbol: orderDetails.Symbol}).Limit(limit).Offset(offset).Error; err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

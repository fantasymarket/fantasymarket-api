package database_test

import (
	"fantasymarket/database/models"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type AddOrderTestData struct {
	userID uuid.UUID
	input  models.Order
	expect models.Order
}

var testAddOrderData = []AddOrderTestData{
	{
		input: models.Order{
			Type:   "stock",
			Symbol: "GOOG",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "HALLO",
		},
	},
	{
		input: models.Order{
			Type:   "stock",
			Symbol: "GOOG",
		},
		expect: models.Order{
			Type:   "stock",
			Symbol: "GOOG",
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestAddOrderAndGetOrder() {

	result := *models.Order{}
	currentDate := parseTime("2019-12-30T15:00:05Z")
	for order := range testAddOrderData {
		err := suite.dbService.AddOrder(order, userID, currentDate)
		assert.Equal(suite.T(), nil, err)
		err = s.DB.Where(models.Order{UserID: order.UserID, Type: order.Type, Symbol: order.Symbol}).Limit(1).First(&result).Offset(-1).Error
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), order.expect, result)
	}

	suite.dbService.DB.Close()
}

// // AddOrder adds an Order to the database
// func (s *Service) AddOrder(order models.Order, userID uuid.UUID, currentDate time.Time) error {
// 	return s.DB.Create(&models.Order{
// 		UserID:    userID,
// 		CreatedAt: currentDate,
// 		Type:      order.Type,
// 		Side:      order.Side,
// 		Symbol:    order.Symbol,
// 		Status:    order.Status,
// 	}).Error
// }

// // GetOrders gets all orders based on the parameters of orderDetails where Symbol, Type and userID can be set.
// // Limit is how many items. Offset is from where to where the data is used
// func (s *Service) GetOrders(orderDetails models.Order, limit int, offset int) (*[]models.Order, error) {
// 	var orders *[]models.Order
// 	if err := s.DB.Where(models.Order{UserID: orderDetails.UserID, Type: orderDetails.Type, Symbol: orderDetails.Symbol}).Limit(limit).Offset(offset).Error; err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }
// }

package database

import (
	"errors"
	"fantasymarket/database/models"
	"fantasymarket/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	// ErrOrderFilledOrCancelled means a order is already filled or cancelled
	ErrOrderFilledOrCancelled = errors.New("can't cancel order, as its already filled or cancelled")
	// ErrInvalidAmount means the order's amount is invalid
	ErrInvalidAmount = errors.New("amount cannot be less than 0")
	// ErrNotEnoughMoney means the order can't be executed
	ErrNotEnoughMoney = errors.New("insufficient balance")
	// ErrInvalidType means the type given in the order is not found
	ErrInvalidType = errors.New("404 type not found")
	// ErrCantSellMoreThanYouHave meansyou can't sell more than
	ErrCantSellMoreThanYouHave = errors.New("cant sell more than you have")
	// ErrOrderCantBeNil is if the user wants to add an order that is empty
	ErrOrderCantBeNil = errors.New("you cant add an empty order")
	// ListOFValidTypes is the list of types accepted for trading
	ListOFValidTypes = [3]string{"stock", "crypto", "commodities"}
)

// AddOrder adds an Order to the database
func (s *Service) AddOrder(order models.Order, userID uuid.UUID, currentDate time.Time) error {
	if order.OrderID == uuid.Nil {
		return ErrOrderCantBeNil
	}
	return s.DB.Create(&models.Order{
		UserID:    userID,
		CreatedAt: currentDate,
		Type:      order.Type,
		Side:      order.Side,
		Symbol:    order.Symbol,
		Status:    order.Status,
	}).Error
}

// GetOrders gets all orders based on the parameters of orderDetails where Symbol, Type and userID can be set.
// Limit is how many items. Offset is from where to where the data is used
func (s *Service) GetOrders(orderDetails models.Order, limit int, offset int) ([]models.Order, error) {
	var orders []models.Order
	if err := s.DB.Where(models.Order{UserID: orderDetails.UserID, Type: orderDetails.Type, Symbol: orderDetails.Symbol}).Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID gets an order by the orderID
func (s *Service) GetOrderByID(orderID uuid.UUID) (models.Order, error) {
	var order models.Order
	if err := s.DB.Where(models.Order{OrderID: orderID}).First(&order).Error; err != nil {
		return models.Order{}, err
	}

	return order, nil
}

// CancelOrder cancels an order in the database
func (s *Service) CancelOrder(orderID uuid.UUID, currentDate time.Time) error {

	var order models.Order
	if err := s.DB.Where(models.Order{OrderID: orderID}).First(&order).Error; err != nil {
		return err
	}

	// TODO check if the order is still active
	if order.Status == "waiting" {
		return s.DB.Model(&order).Updates(models.Order{Status: "cancelled", FilledAt: currentDate}).Error
	}

	return ErrOrderFilledOrCancelled
}

// UpdateOrder updates a order
func (s *Service) UpdateOrder(orderID uuid.UUID, update models.Order) error {
	return s.DB.Table("orders").Where(models.Order{OrderID: orderID}).Updates(&update).Error
}

// FillOrder "completes" the order
// Note: this is only supposed to be called from the game code
// This should ONLY be called if the order is actually supposed to run
// and things like the order type, limit prices etc. have been checked
// (since we only check if the user has enough money in his portfolio)
func (s *Service) FillOrder(orderID uuid.UUID, userID uuid.UUID, currentIndex int64, currentDate time.Time) error {
	var order models.Order
	var user models.User

	if err := s.DB.Where(models.Order{
		OrderID: orderID,
	}).Find(&order).Error; err != nil {
		return err
	}

	if err := s.DB.Where(models.User{
		UserID: userID,
	}).Preload("Portfolio").Find(&user).Error; err != nil {
		return err
	}

	if order.Amount < 0 {
		return ErrInvalidAmount
	}

	if !utils.Includes(ListOFValidTypes[:], order.Type) {
		return ErrInvalidType
	}

	if order.Side == "sell" {
		order.Amount = order.Amount * -1
	}

	price := order.Amount * currentIndex
	if (user.Portfolio.Balance + price) < 0 {
		s.CancelOrder(order.OrderID, currentDate)
		return ErrNotEnoughMoney
	}

	var affectedPortfolioItem models.PortfolioItem
	if err := s.DB.Where(models.PortfolioItem{
		PortfolioID: user.Portfolio.PortfolioID,
		Symbol:      order.Symbol,
	}).Attrs(models.PortfolioItem{
		PortfolioID: user.Portfolio.PortfolioID,
		Symbol:      order.Symbol,
		Type:        order.Type,
		Amount:      0,
	}).FirstOrCreate(&affectedPortfolioItem).Error; err != nil {
		s.CancelOrder(order.OrderID, currentDate)
		return err
	}

	newAmount := order.Amount + affectedPortfolioItem.Amount
	if newAmount < 0 {
		s.CancelOrder(order.OrderID, currentDate)
		return ErrCantSellMoreThanYouHave
	}

	newBalance := user.Portfolio.Balance - price
	if err := s.updatePortfolioItem(user.Portfolio.PortfolioID, affectedPortfolioItem.PortfolioItemID, newAmount, newBalance); err != nil {
		s.CancelOrder(order.OrderID, currentDate)
		return err
	}

	return s.DB.Where(models.Order{OrderID: orderID}).Updates(models.Order{Status: "filled", FilledAt: currentDate}).Error
}

func (s *Service) updatePortfolioItem(portfolioID uuid.UUID, itemID uuid.UUID, newAmount int64, newBalance int64) error {
	// we use a transactions since if updating
	// balance fails, we also need to rollback
	// the portfolioItem's amount
	tx := s.DB.Begin()

	if err := tx.Where(&models.PortfolioItem{
		PortfolioItemID: itemID,
	}).Updates(models.PortfolioItem{
		Amount: newAmount,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where(&models.Portfolio{
		PortfolioID: portfolioID,
	}).Updates(models.Portfolio{
		Balance: newBalance,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

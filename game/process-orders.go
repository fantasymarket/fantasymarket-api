package game

import (
	"fantasymarket/database/models"
)

//ProcessOrders processes the order if it matches any of the order types.
func (s *Service) ProcessOrders(orders []models.Order) error {
	getCurrentDate := s.GetCurrentDate
	currentStocks, err := s.DB.GetStockMapAtTick(s.TicksSinceStart)
	if err != nil {
		return err
	}

	lastStocks, _ := s.DB.GetStockMapAtTick(s.TicksSinceStart - 1)

	for _, order := range orders {
		switch order.Type {
		case "Market":
			s.DB.FillOrder(order.OrderID, order.UserID, currentStocks[order.Symbol].Index, getCurrentDate())
		case "Stop Loss":
			if currentStocks[order.Symbol].Index <= order.StopLossValue {
				s.DB.FillOrder(order.OrderID, order.UserID, currentStocks[order.Symbol].Index, getCurrentDate())
			}
		case "Limit":
			if currentStocks[order.Symbol].Index <= order.BuyAtValue {
				s.DB.FillOrder(order.OrderID, order.UserID, currentStocks[order.Symbol].Index, getCurrentDate())
			} else if currentStocks[order.Symbol].Index >= order.SellAtValue {
				s.DB.FillOrder(order.OrderID, order.UserID, currentStocks[order.Symbol].Index, getCurrentDate())
			}
		case "Trailing Stop":
			if currentStocks[order.Symbol].Index <= lastStocks[order.Symbol].Index {

			}

		}
	}
	return nil
}

package game

import (
	"fantasymarket/database/models"
)

//ProcessOrders processes the order if it matches any of the order types.
func (s *Service) ProcessOrders(orders []models.Order) error {
	CurrentStockPrice := map[string]models.Stock
	if CurrentStockPrice, err := s.DB.GetStockMapAtTick(s.TicksSinceStart); err != nil {
		return err
	}

	for _, order := range orders {
		switch order.Type {
		case "Market":
			s.DB.FillOrder(order.OrderID, order.UserID, CurrentStockPrice[order.Symbol].Index, getCurrentDate())
		case "Stop Loss":
			if CurrentStockPrice <= order.StopLossValue {
				s.DB.FillOrder(order.OrderID, order.UserID, CurrentStockPrice[order.Symbol].Index, getCurrentDate())
			}
		case "Limit":
			if CurrentStockPrice <= order.BuyAtValue {
				s.DB.FillOrder(order.OrderID, order.UserID, CurrentStockPrice[order.Symbol].Index, getCurrentDate())
			} else if CurrentStockPrice >= order.SellAtValue {
				s.DB.FillOrder(order.OrderID, order.UserID, CurrentStockPrice[order.Symbol].Index, getCurrentDate())
			}
		case ""
		}
	}
	return nil
}

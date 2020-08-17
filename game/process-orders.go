package game

import (
	"fantasymarket/database/models"

	"github.com/shopspring/decimal"
)

//ProcessOrder is an instant of an order that needs to be processed
type ProcessOrder struct {
	order        models.Order
	currentStock models.Stock
	service      Service
}

func (p ProcessOrder) cancelOrder() {
	p.service.DB.CancelOrder(p.order.OrderID, p.service.GetCurrentDate())
}
func (p ProcessOrder) fillOrder() {
	if err := p.service.DB.FillOrder(p.order.OrderID, p.order.UserID, p.currentStock.Index, p.service.GetCurrentDate()); err != nil {
		p.cancelOrder()
	}
}
func (p ProcessOrder) updateOrder(update models.Order) {
	if err := p.service.DB.UpdateOrder(p.order.OrderID, update); err != nil {
		p.cancelOrder()
	}
}

//ProcessOrders processes the order if it matches any of the order types.
func (s *Service) ProcessOrders(orders []models.Order) error {

	currentStocks, err := s.DB.GetStockMapAtTick(s.TicksSinceStart)
	if err != nil {
		return err
	}

	for _, order := range orders {

		currentStock := currentStocks[order.Symbol]

		processOrder := ProcessOrder{order: order, service: *s, currentStock: currentStock}

		// Decimalizes the prices of the stock and the trailing percentage for the stop-loss order to improve precision
		currentStockIndex := decimal.NewFromInt(currentStock.Index)
		trailingPercentage := decimal.NewFromInt(order.TrailingPercentage).Div(decimal.NewFromInt(100))

		switch order.Type {
		case "market":

			processOrder.fillOrder()

		case "stop":
			if order.Side == "sell" && currentStock.Index <= order.Price {
				processOrder.fillOrder()
			}
			if order.Side == "buy" && currentStock.Index >= order.Price {
				processOrder.fillOrder()
			}

		case "limit":
			if order.Side == "sell" && currentStock.Index >= order.Price {
				processOrder.fillOrder()
			}

			if order.Side == "buy" && currentStock.Index <= order.Price {
				processOrder.fillOrder()
			}

		case "trailing-stop":

			if order.TrailingPercentage < 0 || order.TrailingPercentage > 1 {
				processOrder.cancelOrder()
			}

			var newPrice int64

			if order.Side == "sell" {
				if currentStock.Index <= order.Price {
					processOrder.fillOrder()
					break
				}

				newPrice = currentStockIndex.Sub(currentStockIndex.Mul(trailingPercentage)).Round(0).IntPart()

				if newPrice <= order.Price {
					break
				}
			}

			if order.Side == "buy" {
				if currentStock.Index >= order.Price {
					processOrder.fillOrder()
					break
				}

				newPrice = currentStockIndex.Add(currentStockIndex.Mul(trailingPercentage)).Round(0).IntPart()

				if newPrice >= order.Price {
					break
				}
			}

			processOrder.updateOrder(models.Order{
				Price: newPrice,
			})
		}
	}

	return nil
}

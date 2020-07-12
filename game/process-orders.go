package game

import (
	"fantasymarket/database/models"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

//ProcessOrders processes the order if it matches any of the order types.
func (s *Service) ProcessOrders(orders []models.Order) error {

	currentDate := s.GetCurrentDate()

	currentStocks, err := s.DB.GetStockMapAtTick(s.TicksSinceStart)
	if err != nil {
		return err
	}

	for _, order := range orders {

		currentStock := currentStocks[order.Symbol]

		// Decimalizes the prices of the stock and the trailing percentage for the stop-loss order to improve precision
		currentStockIndex := decimal.NewFromInt(currentStock.Index)
		trailingPercentage := decimal.NewFromInt(order.TrailingPercentage).Div(decimal.NewFromInt(100))

		cancelOrder := func() {
			log.Error().Err(err).Str("orderID", order.OrderID.String()).Msg("error filling order")
			s.DB.CancelOrder(order.OrderID, currentDate)
		}

		fillOrder := func() {
			if err := s.DB.FillOrder(order.OrderID, order.UserID, currentStock.Index, currentDate); err != nil {
				cancelOrder()
				log.Error().Err(err).Str("orderID", order.OrderID.String()).Msg("failed to execute order")
			}
		}

		updateOrder := func(update models.Order) {
			if err := s.DB.UpdateOrder(order.OrderID, update); err != nil {
				cancelOrder()
				log.Error().Err(err).Str("orderID", order.OrderID.String()).Msg("failed to update order")
			}
		}

		switch order.Type {
		case "market":
			// the order will sell at the next best available price.
			fillOrder()

		case "stop":
			if order.Side == "sell" {
				if currentStock.Index <= order.Price {
					fillOrder()
				}
			}
			//Buy stop order: set a stop price above the current price of the stock. If stock rises to stop price buy stop order becomes buy market order.
			if order.Side == "buy" {
				if currentStock.Index >= order.Price {
					fillOrder()
				}
			}

		case "limit":
			// Limit orders specify the minimum amount you are willing to receive when selling a stock.
			if order.Side == "sell" {
				if currentStock.Index >= order.Price {
					fillOrder()
				}
			}

			// Limit orders specify the maximum amount you are willing to pay for a stock.
			if order.Side == "buy" {
				if currentStock.Index <= order.Price {
					fillOrder()
				}
			}

		case "trailing-stop":

			if order.TrailingPercentage < 0 || order.TrailingPercentage > 1 {
				cancelOrder()
			}

			if order.Side == "sell" {
				if currentStock.Index <= order.Price {
					fillOrder()
					break
				}

				newPrice := currentStockIndex.Sub(currentStockIndex.Mul(trailingPercentage)).Round(0).IntPart()

				// If the newPrice stays between the stop value
				// the limit can never go down
				if newPrice <= order.Price {
					break
				}

				updateOrder(models.Order{
					Price: newPrice,
				})
			}

			if order.Side == "buy" {
				if currentStock.Index >= order.Price {
					fillOrder()
					break
				}

				newPrice := currentStockIndex.Add(currentStockIndex.Mul(trailingPercentage)).Round(0).IntPart()

				if newPrice >= order.Price {
					break
				}

				updateOrder(models.Order{
					Price: newPrice,
				})
			}
		}
	}

	return nil
}

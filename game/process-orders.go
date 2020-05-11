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

			// trailing stop for sell: Say Stock A is at $100. A user puts the trailing stop of 10% meaning that the stop value is $90. As long
			// as CurrentStock.Index is between 90 and 100, the stop value doesn't change. If it goes above $100, the stop value goes up by the percentage.
			// If the CurrentStock.Index reaches $90 or below, the order gets executed.

			// current price: 100
			// trailing stop: 10%
			// starting stop value: 90$

			// scenario 1
			// stock price changes to 102
			// -> newPrice 91.8
			// we update

			// scenario 2
			// stock price changes to 98
			// -> newPrice88.2
			// 88.2 < 90 => we dont updates

			// scanerio 3
			// stock price changes to <=90
			// fillOrder() is called

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

			// trailing stop for buy: Say Stock A is priced at $110. Trailing percentage is 5%, starting stop value : $115.50

			// scenario 1:
			// stock prices fall to 100
			//	-> newPrice: 5% higher than lowest price (100), becomes $105
			//	we update

			// scenario 2:
			// stock prices rise to 113
			//	->newPrice : 105.65, greater than 105, don't update

			// scenario 3:
			// stock prices rise to >= 115.50
			// fillOrder() is called

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

package game

import (
	"fantasymarket/database/models"

	"github.com/rs/zerolog/log"
)

//ProcessOrders processes the order if it matches any of the order types.
func (s *Service) ProcessOrders(orders []models.Order) error {

	currentDate := s.GetCurrentDate()

	currentStocks, err := s.DB.GetStockMapAtTick(s.TicksSinceStart)
	if err != nil {
		return err
	}

	lastStocks, err := s.DB.GetStockMapAtTick(s.TicksSinceStart - 1)
	if err != nil {
		return err
	}

	for _, order := range orders {
		currentStock := currentStocks[order.Symbol]
		lastStock := lastStocks[order.Symbol]

		cancelOrder := func() {
			log.Error().Err(err).Str("orderID", order.OrderID.String()).Msg("error filling order")
			s.DB.CancelOrder(order.OrderID, currentDate)
		}

		fillOrder := func() error {
			err := s.DB.FillOrder(order.OrderID, order.UserID, currentStock.Index, currentDate)
			if err != nil {
				cancelOrder()
			}
		}

		switch order.Type {
		case "market":
			// the order will sell at the next best available price.
			s.DB.FillOrder(order.OrderID, order.UserID, currentStock.Index, currentDate)

		case "stop-loss":

			// Stop Loss sell orders trigger a market order to sell when the stop price is met.
			// Ex. XYZ stock is trading at $25. A stop order can be placed at $20 to trigger
			// a market sell order when a trade executes at $20 or lower.
			if order.Side == "sell" {
				if currentStocks[order.Symbol].Index >= order.StopLossPrice {
					fillOrder()
				}
			}

			// Stop Loss buy orders trigger a market order to buy when the stop price is reached.
			// Stop loss orders are sent as stop limit orders with the limit price collared up to 5% above the stop price.
			// Ex. ABC stock is trading at $10. A stop order can be placed at $11 to trigger a market buy order
			// when a trade executes at $11 or higher.
			if order.Side == "buy" {
				if currentStocks[order.Symbol].Index <= order.StopLossPrice {
					fillOrder()
				}
			}

		case "limit":
			// Limit orders specify the minimum amount you are willing to receive when selling a stock.
			if order.Side == "sell" {
				if currentStocks[order.Symbol].Index >= order.Price {
					fillOrder()
				}
			}

			// Limit orders specify the maximum amount you are willing to pay for a stock.
			if order.Side == "buy" {
				if currentStocks[order.Symbol].Index <= order.Price {
					fillOrder()
				}

			}

		case "trailing-stop":

			difference := currentStock.Index - lastStock.Index
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

				var newPrice float64 = currentStock.Index - currentStock.Index*order.TrailingPercentage

				//if the newPrice stays between the stop value
				// the limit can never go down
				if newPrice > order.Price {
					// increase the limit price
				}
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

				newPrice := currentStock.Index + currentStock.Index*order.TrailingPercentage
				if newPrice < order.Price {
					// Update price if some shit is true
				}
			}
		}
	}
	return nil
}

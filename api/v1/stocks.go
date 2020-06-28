package v1

import (
	"fantasymarket/game/details"
	"fantasymarket/utils/http/responses"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func (api *APIHandler) getAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := details.StocksYamlBytes()
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}

	m := make([]details.StockDetails, 30)
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}

	responses.CustomResponse(w, m, 200)
}

func (api *APIHandler) getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	givenTime := chi.URLParam(r, "time")

	stock, ok := api.Game.StockDetails[symbol]
	if !ok {
		responses.ErrorResponse(w, http.StatusInternalServerError, stockNotFoundError.Error())
		return
	}

	if givenTime != "" {
		tick, err := api.getTickAtTime(givenTime)
		if err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, stockNotFoundError.Error())
			return
		}
		stockMapAtTick, err := api.DB.GetStockMapAtTick(tick)
		if err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, stockNotFoundError.Error())
			return
		}
		desiredStock := stockMapAtTick[stock.Symbol]
		responses.CustomResponse(w, desiredStock, http.StatusOK)
		return
	}

	stockMapAtTick, err := api.DB.GetStockMapAtTick(api.Game.TicksSinceStart)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, stockNotFoundError.Error())
		return
	}
	responses.CustomResponse(w, stockMapAtTick[symbol], http.StatusOK)
}

func (api *APIHandler) getTickAtTime(timestamp string) (int64, error) {
	startTime, err := time.Parse(time.RFC3339, api.Config.Game.StartDate.String())
	if err != nil {
		return 0, err
	}
	currentTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, err
	}
	difference := int64(currentTime.Sub(startTime).Hours())

	return difference, nil
}

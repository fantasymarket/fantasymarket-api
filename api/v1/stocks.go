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
		responses.ErrorResponse(w, http.StatusInternalServerError, errFetchingData.Error())
		return
	}

	m := make([]details.StockDetails, 30)
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errFetchingData.Error())
		return
	}

	responses.CustomResponse(w, m, 200)
}

func (api *APIHandler) getStockValue(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	givenTime := chi.URLParam(r, "time")

	tick := api.Game.TicksSinceStart

	if givenTime != "" {
		var err error
		tick, err = api.getTickAtTime(givenTime)
		if err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, errStockNotFound.Error())
			return
		}
	}

	stock, err := api.DB.GetStockAtTick(symbol, tick)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errStockNotFound.Error())
		return
	}

	responses.CustomResponse(w, stock, http.StatusOK)
}

func (api *APIHandler) getStockHistory(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	fromTime := chi.URLParam(r, "from")
	toTime := chi.URLParam(r, "to")

	from := int64(0)
	to := api.Game.TicksSinceStart

	var err error
	if fromTime != "" {
		from, err = api.getTickAtTime(fromTime)
	}

	if toTime != "" {
		to, err = api.getTickAtTime(toTime)
	}

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errInvalidParameters.Error())
		return
	}

	stock, err := api.DB.GetStockData(symbol, from, to)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errStockNotFound.Error())
		return
	}

	responses.CustomResponse(w, stock, http.StatusOK)
}

func (api *APIHandler) getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	stock, ok := api.Game.StockDetails[symbol]
	if !ok {
		responses.ErrorResponse(w, http.StatusInternalServerError, errStockNotFound.Error())
		return
	}

	responses.CustomResponse(w, stock, http.StatusOK)
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

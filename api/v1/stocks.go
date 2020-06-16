package v1

import (
	"errors"
	"fantasymarket/game/details"
	"fantasymarket/utils/http/responses"
	"fantasymarket/utils/timeutils"
	"net/http"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func (api *APIHandler) getAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := details.StocksYamlBytes()
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}

	m := []details.StockDetails{}
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}

	responses.CustomResponse(w, m, 200)
}

func (api *APIHandler) getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	time := chi.URLParam(r, "time")

	yamlData, err := details.StocksYamlBytes()
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "Error getting Stock Details")
	}

	var myStocks []details.StockDetails

	if err := yaml.Unmarshal(yamlData, &myStocks); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "Error parsing the stock")
	}

	stock, err := getStockHelper(myStocks, symbol)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "Error getting the Stock Detail")
		return
	}
	if time != "" {
		tick := timeutils.GetTickAtTime(time)
		// TODO: Implement the GetTickAtTime function and then grab the corresponding
		// TODO: Stock at that tick from db.GetStockAtTick and then return that stock
		// TODO: to the client.
		responses.CustomResponse(w, stockAtTime, http.StatusOK)
		return
	}

	responses.CustomResponse(w, stock, http.StatusOK)
}

func getStockHelper(stocks []details.StockDetails, symbol string) (*details.StockDetails, error){
	for i := range stocks {
		if stocks[i].Symbol == symbol {
			return &stocks[i], nil
		}
	}

	return nil, stockNotFoundError
}

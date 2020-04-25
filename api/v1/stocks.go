package v1

import (
	"fantasymarket/game/details"
	"fantasymarket/utils/http/responses"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func (api *APIHandler) getAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := ioutil.ReadFile("game/stocks.yaml")
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error getting list of stocks")
		return
	}

	m := []details.StockDetails{} // Test again if this works, I checked the docs and in theory when the map is already initialized it should still work. If not change it back again
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error parsing stocks")
		return
	}

	responses.CustomResponse(w, m, 200)
}

func (api *APIHandler) getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	yamlData, err := ioutil.ReadFile("game/stocks.yaml")

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error getting list of stocks")
	}

	var stocks []details.StockDetails

	if err := yaml.Unmarshal(yamlData, &stocks); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error parsing stock")
	}

	for i := range stocks {
		if stocks[i].Symbol == symbol {
			responses.CustomResponse(w, stocks[i], 200)
			return
		}
	}

	responses.ErrorResponse(w, http.StatusNotFound, "no stock with symbol available")
}

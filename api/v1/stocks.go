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
	yamlData, err := ioutil.ReadFile("game/stocks.yaml")

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
	}

	var myStocks []details.StockDetails

	if err := yaml.Unmarshal(yamlData, &myStocks); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
	}

	for i := range myStocks {
		if myStocks[i].Symbol == symbol {
			responses.CustomResponse(w, myStocks[i], 200)
			return
		}
	}

	responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
}

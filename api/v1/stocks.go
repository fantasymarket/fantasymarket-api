package v1

import (
	"fantasymarket/game/stocks"
	"fantasymarket/utils/http/responses"
	"io/ioutil"
	"net/http"
	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func (api *APIHandler) getAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := ioutil.ReadFile("game/stocks.yaml")
	if err != nil {
		responses.ErrorResponse(w, "Error getting list of stocks", http.StatusInternalServerError)
		return
	}

	m := []stocks.StockDetails{}
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, "Error parsing the stocks", http.StatusInternalServerError)
		return
	}

	responses.CustomResponse(w, m, 200)
}

func (api *APIHandler) getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	yamlData, err := ioutil.ReadFile("game/stocks.yaml")

	if err != nil {
		responses.ErrorResponse(w, "Error getting Stock Details", http.StatusInternalServerError)
	}

	var myStocks []stocks.StockDetails

	if err := yaml.Unmarshal(yamlData, &myStocks); err != nil {
		responses.ErrorResponse(w, "Error parsing the stock", http.StatusInternalServerError)
	}

	for i := range myStocks {
		if myStocks[i].Symbol == symbol {
			responses.CustomResponse(w, myStocks[i], 200)
			return
		}
	}

	responses.ErrorResponse(w, "Error getting the Stock Detail", http.StatusInternalServerError)
}

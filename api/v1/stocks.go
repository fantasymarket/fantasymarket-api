package v1

import (
	"fantasymarket/game/details"
	"fantasymarket/utils/http/responses"
	"fantasymarket/utils/timeutils"
	"fmt"
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

	m := make([]details.StockDetails, 20)
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
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}

	var myStocks []details.StockDetails
	if err := yaml.Unmarshal(yamlData, &myStocks); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, decodingError.Error())
		return
	}

	stock, err := getStockHelper(myStocks, symbol)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, stockNotFoundError.Error())
		return
	}
	if time != "" {
		tick := timeutils.GetTickAtTime(time)
		stockMapAtTick, err := api.DB.GetStockMapAtTick(tick)
		if err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, stockNotFoundError.Error())
			return
		}
		desiredStock := stockMapAtTick[stock.Name]
		responses.CustomResponse(w, desiredStock, http.StatusOK)
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

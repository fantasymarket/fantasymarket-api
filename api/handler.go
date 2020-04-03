package api

import (
	"fantasymarket/game/stocks"
	"fantasymarket/utils/http/responses"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func getAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := ioutil.ReadFile("game/stocks.yaml")
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error getting list of stocks")
		return
	}

	m := []stocks.StockDetails{} // Test again if this works, I checked the docs and in theory when the map is already initialized it should still work. If not change it back again
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error parsing stocks")
		return
	}

	responses.CustomResponse(w, m, 200)
}

func getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	yamlData, err := ioutil.ReadFile("game/stocks.yaml")

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error getting list of stocks")
	}

	var stocks []stocks.StockDetails

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

func getPortfolio(w http.ResponseWriter, r *http.Request) {

}

func orders(w http.ResponseWriter, r *http.Request) {

}

func getEvents(w http.ResponseWriter, r *http.Request) {
	events, err := db.GetEvents()

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're ginormously fucked")
	}

	responses.CustomResponse(w, events, 200)

}

func getOverview(w http.ResponseWriter, r *http.Request) {

}

func getTime(w http.ResponseWriter, r *http.Request) {
	t := g.Options.StartDate

	if !t.IsZero() {
		responses.CustomResponse(w, t.Format(time.RFC3339), 200)
	} else {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're absolutely fucked")
	}

}

func addOrder(w http.ResponseWriter, r *http.Request) {

}

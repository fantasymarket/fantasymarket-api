package api

import (
	"fantasymarket/game"
	"fantasymarket/utils/http/responses"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func getAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := ioutil.ReadFile("game/stocks.json")
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're really fucked")
	}

	m := make(map[string]game.StockSettings) // Test again if this works, I checked the docs and in theory when the map is already initialized it should still work. If not change it back again
	err = yaml.Unmarshal(allStocks, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're royally fucked")
	}

	responses.CustomResponse(w, allStocks, 200)
}

func getStockDetails(w http.ResponseWriter, r *http.Request) {
	stockID := chi.URLParam(r, "stockID")
	yamlData, err := ioutil.ReadFile("game/stocks.json")

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're majorly fucked")
	}

	var stockDetail map[string]game.StockSettings
	err = yaml.Unmarshal(yamlData, &stockDetail)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're hugely fucked")
	}

	responses.CustomResponse(w, stockDetail[stockID], 200)

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

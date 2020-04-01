package api

import (
	"fantasymarket/game"
	"fantasymarket/utils/http_responses"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	allStocks, err := ioutil.ReadFile("game/stocks.json")

	if err != nil {
		http_responses.ErrorResponse(w, http.StatusInternalServerError, "we're really fucked")
		return
	}

	stockMap := make(map[string]game.StockSettings) // Test again if this works, I checked the docs and in theory when the map is already initialized it should still work. If not change it back again
	err = yaml.Unmarshal(allStocks, &stockMap)

	if err != nil {
		http_responses.ErrorResponse(w, http.StatusInternalServerError, "we're royally fucked")
		return
	}

	http_responses.CustomResponse(w, stockMap, 200)
}

func GetStockDetails(w http.ResponseWriter, r *http.Request) {
	urlparam := chi.URLParam(r, "stockID")
	stockID := strings.SplitN(urlparam, "=", 2)[1]
	yamlData, err := ioutil.ReadFile("game/stocks.json")


	if err != nil {
		http_responses.ErrorResponse(w, http.StatusInternalServerError, "we're majorly fucked")
		return
	}

	stockMap := make(map[string]game.StockSettings)
	err = yaml.Unmarshal(yamlData, &stockMap)

	if err != nil {
		http_responses.ErrorResponse(w, http.StatusInternalServerError, "we're hugely fucked")
		return
	}

	http_responses.CustomResponse(w, stockMap[stockID], 200)

}

func GetPortfolio(w http.ResponseWriter, r *http.Request) {

}

func Orders(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {

	}

}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := db.GetEvents()

	if err != nil {
		http_responses.ErrorResponse(w, http.StatusInternalServerError, "we're ginormously fucked")
		return
	}

	http_responses.CustomResponse(w, events, 200)

}

func GetOverview(w http.ResponseWriter, r *http.Request) {

}

func GetTime(w http.ResponseWriter, r *http.Request) {
	t := g.Options.StartDate

	if !t.IsZero() {
		http_responses.CustomResponse(w, t.Format(time.RFC3339), 200)
	} else {
		http_responses.ErrorResponse(w, http.StatusInternalServerError, "we're absolutely fucked")
	}

}

func AddOrder(w http.ResponseWriter, r *http.Request) {

}

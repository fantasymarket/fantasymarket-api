package api

import (
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"fantasymarket/database"
	"fantasymarket/game"
	"gopkg.in/yaml.v3"
)

func CheckInternalError(w http.ResponseWriter, err... error) error {
	for _, error := range err {
		if error != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error Occured"))
			return error
		}
	}
	return nil
}


func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	yamlData, errOne := ioutil.ReadFile("game/stocks.json")

	m := make(map[string]game.StockSettings) // Test again if this works, I checked the docs and in theory when the map is already initialized it should still work. If not change it back again
	errTwo := yaml.Unmarshal(yamlData, &m)

	internalError := CheckInternalError(w, errOne, errTwo)

	if internalError == nil {
		w.WriteHeader(200)
		w.Write(yamlData)
	}


}

func GetStockDetails(w http.ResponseWriter, r *http.Request) {
	stockID := chi.URLParam(r, "stockID")
	yamlData, errOne := ioutil.ReadFile("game/stocks.json")

	var m map[string]game.StockSettings
	errTwo := yaml.Unmarshal(yamlData, &m)

	stockDetail, errThree := yaml.Marshal(m[stockID])
	internalError := CheckInternalError(w, errOne, errTwo, errThree)

	if internalError == nil {
		w.WriteHeader(200)
		w.Write(stockDetail)
	}

}

func GetPortfolio(w http.ResponseWriter, r *http.Request) {

}

func Orders(w http.ResponseWriter, r *http.Request) {

}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	events, errOne := database.GetEvents() // Check why you can't access events method

	internalError := CheckInternalError(errOne)

	if internalError == nil {
		w.WriteHeader(200)
		w.Write([]byte(events))
	}
}

func GetOverview(w http.ResponseWriter, r *http.Request) {

}

func GetTime(w http.ResponseWriter, r *http.Request) {

}

func AddOrder(w http.ResponseWriter, r *http.Request) {

}

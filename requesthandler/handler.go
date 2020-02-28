package requesthandler

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func GetPortfolioNumbers(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	}


func GetStockNumbers(w http.ResponseWriter, r *http.Request) {


}

func BuyStock(w http.ResponseWriter, r *http.Request) {
	stock_name := chi.URLParam(r, "name")


}

func SellStock(w http.ResponseWriter, r *http.Request) {
	stock_name := chi.URLParam(r, "name")


}

func UpdateStocks(w http.ResponseWriter, r *http.Request) {
	stock_name := chi.URLParam(r, "name")


}
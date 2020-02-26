package requesthandler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func GetPortfolioNumbers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	}
}

func GetStockNumber(w http.ResponseWriter, r *http.Request) {
	stockName := chi.URLParam(r, "name")

	fmt.Println(stockName)
	return
}

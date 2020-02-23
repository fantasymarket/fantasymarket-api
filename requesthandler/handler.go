package requesthandler

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func GetPortfolioNumbers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alex = Gay"))
	}
}

func GetStockRequest(w http.ResponseWriter, r *http.Request) {
	stock_name := chi.URLParam(r, "name")


}



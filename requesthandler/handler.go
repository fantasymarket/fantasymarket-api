package requesthandler

import (
	"encoding/json"
	"fantasymarket/mock-data"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)



func GetStockNumbers(w http.ResponseWriter, r *http.Request) {
	 googStock := mock_data.Stocks{Name: "Google", Index: int64(100000), Trend: int64(1)}
	 msftStock := mock_data.Stocks{Name: "Microsoft", Index: int64(100050), Trend: int64(2)}
	 structArray := []interface{}{googStock, msftStock}

	 w.Header().Set("Content-Type", "application/json")
	 w.WriteHeader(200)
	 json.NewEncoder(w).Encode(structArray)

}

func GetStockDetails(w http.ResponseWriter, r *http.Request) {


}

func UserStats(w http.ResponseWriter, r *http.Request) {
	stock_name := chi.URLParam(r, "name")
	fmt.Println(stock_name)


}

func Orders(w http.ResponseWriter, r *http.Request) {


}

func News(w http.ResponseWriter, r *http.Request) {


}
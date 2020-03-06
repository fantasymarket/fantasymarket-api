package requesthandler

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)



func GetStockNumbers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("works")

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
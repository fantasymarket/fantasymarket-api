package api

import (
	"encoding/json"
	"fantasymarket/mock-data"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	//"fantasymarket/mock-data"
	"net/http"
	"net/http/httptest"
	"fantasymarket/requesthandler"
)

func TestApiStockData(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	requesthandler.GetStockNumbers(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	var data []interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(data)
	}
	googStock := mock_data.Stocks{Name: "Google", Index: int64(100000), Trend: int64(1)}
	msftStock := mock_data.Stocks{Name: "Microsoft", Index: int64(100050), Trend: int64(2)}
	mockdata := []interface{}{googStock, msftStock}

	if data == mockdata {

	}
}
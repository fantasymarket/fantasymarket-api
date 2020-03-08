package api

import (
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

	if string(body) != "[{\"Name\":\"Google\",\"Index\":100000,\"Trend\":1},{\"Name\":\"Microsoft\",\"Index\":100050,\"Trend\":2}]" {
		fmt.Println(string(body))
		t.Errorf("want body to equal %q", body)
	}


}
package v1

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

func (api *APIHandler) orders(w http.ResponseWriter, r *http.Request) {

}

func (api *APIHandler) addOrder(w http.ResponseWriter, r *http.Request) {
	type Order struct {
		Type string	`json:"type"`
		Side string `json:"side"`
		Symbol string `json:"symbol"`
		Quantity int `json:"quantity"`
		LimitPrice int `json:"limitPrice"`
		StopPrice int `json:"stopPrice"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	newOrder := &Order{}
	err = yaml.Unmarshal(body, newOrder)

	if err != nil {
		log.Println(err)
	}
	// TODO: Call the fitting Database function here
	fmt.Println(newOrder)

}

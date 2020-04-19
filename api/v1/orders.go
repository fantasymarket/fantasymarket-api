package v1

import (
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/database/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

func (api *APIHandler) orders(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(jwt.UserKey).(jwt.Claims)
	if ok != true {
		log.Println(ok)
	}

	orders := api.DB.GetOrdersForUser(user.UserID)
	w.Write(orders)
}

func (api *APIHandler) addOrder(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(jwt.UserKey).(jwt.Claims)
	if ok != true {
		log.Println(ok)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	newOrder := &models.Order{UserID: user.UserID}
	err = yaml.Unmarshal(body, newOrder)

	if err != nil {
		log.Println(err)
	}

}

package v1

import (
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/database/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"fantasymarket/utils/http/responses"
)

func (api *APIHandler) orders(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(jwt.UserKey).(jwt.Claims)
	if ok != true {
		log.Println(ok)
	}

	// TODO: Add corresponding logic here to get all orders
	orders := api.DB.GetOrdersForUser(user.UserID)
	responses.CustomResponse(w, orders, 200)
}

func (api *APIHandler) ordersID(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(jwt.UserKey).(jwt.Claims)
	if ok != true {
		log.Println(ok)
	}

	// TODO: Add corresponding logic here to get a specific order
	orders := api.DB.GetOrderForUserByID(user.UserID)
	responses.CustomResponse(w, orders, 200)
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
	time := api.Config.Game.StartDate
	err = api.DB.AddOrder(*newOrder, user.UserID, time.Time)

	if err != nil {
		responses.ErrorResponse(w, 500, "Order couldn't be added")
	}

}

func (api *APIHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	newOrder := &models.Order{}
	err = yaml.Unmarshal(body, newOrder)

	if err != nil {
		log.Println(err)
	}
	time := api.Config.Game.StartDate
	err = api.DB.CancelOrder(newOrder.OrderID, time.Time)

	if err != nil {
		responses.ErrorResponse(w, 500, "Order couldn't be deleted")
	}
}

func (api *APIHandler) fillOrder(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(jwt.UserKey).(jwt.Claims)
	if ok != true {
		log.Println(ok)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	newOrder := &models.Order{}
	err = yaml.Unmarshal(body, newOrder)

	if err != nil {
		log.Println(err)
	}
	time := api.Config.Game.StartDate
	err = api.DB.FillOrder(newOrder.OrderID, user.UserID, time.Time)

	if err != nil {
		responses.ErrorResponse(w, 500, "Order couldn't be deleted")
	}


}


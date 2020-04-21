package v1

import (
	"fantasymarket/database/models"
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/utils/http/responses"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

type customOrder struct {
	Order models.Order `json:"order"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

func (api *APIHandler) orders(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(jwt.UserKey).(jwt.Claims)
	if ok != true {
		log.Println(ok)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	// TODO: This initialized models.Order instance might get overwritten
	// TODO: after unmarshaling the request body into it, therefore
	// TODO: "deleting" the userID. Test that
	allOrders := &customOrder{Order: models.Order{UserID: user.UserID}}
	err = yaml.Unmarshal(body, allOrders)

	if err != nil {
		log.Println(err)
	}

	orders, err := api.DB.GetOrders(allOrders.Order, allOrders.Limit, allOrders.Offset)
	if err != nil {
		responses.ErrorResponse(w, err.Error(), 500)
		return
	}
	responses.CustomResponse(w, orders, 200)

}


func (api *APIHandler) ordersID(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	requestOrder := &models.Order{}
	err = yaml.Unmarshal(body, requestOrder)

	if err != nil {
		log.Println(err)
	}

	orders, err := api.DB.GetOrderByID(requestOrder.OrderID)
	if err != nil {
		responses.ErrorResponse(w, err.Error(), 500)
		return
	}
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

	requestOrder := &models.Order{UserID: user.UserID}
	err = yaml.Unmarshal(body, requestOrder)

	if err != nil {
		log.Println(err)
	}
	time := api.Config.Game.StartDate
	err = api.DB.AddOrder(*requestOrder, user.UserID, time.Time)

	if err != nil {
		responses.ErrorResponse(w, "Order couldn't be added" , 500)
	}
}

func (api *APIHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	requestOrder := &models.Order{}
	err = yaml.Unmarshal(body, requestOrder)

	if err != nil {
		log.Println(err)
	}
	time := api.Config.Game.StartDate
	err = api.DB.CancelOrder(requestOrder.OrderID, time.Time)

	if err != nil {
		responses.ErrorResponse(w, "Order couldn't be deleted", 500)
		return
	}

	responses.CustomResponse(w, requestOrder, 200)
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
		responses.ErrorResponse(w, "Order couldn't be deleted", 500)
	}

	responses.CustomResponse(w, newOrder.OrderID, 200)

}


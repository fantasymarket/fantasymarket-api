package v1

import (
	"encoding/json"
	"fantasymarket/database/models"
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/utils/http/responses"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v3"
)

type customOrder struct {
	Order  models.Order `json:"order"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

func (api *APIHandler) ordersForUser(w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserFromContext(r.Context())

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

	if allOrders.Limit <= 0 || allOrders.Limit > 20 {
		allOrders.Limit = 20
	}

	orders, err := api.DB.GetOrders(allOrders.Order, allOrders.Limit, allOrders.Offset)
	if err != nil {
		responses.ErrorResponse(w, 500, err.Error())
		return
	}
	responses.CustomResponse(w, orders, 200)

}

func (api *APIHandler) ordersID(w http.ResponseWriter, r *http.Request) {
	var requestOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "Error Parsing Request")
		return
	}

	orders, err := api.DB.GetOrderByID(requestOrder.OrderID)
	if err != nil {
		responses.ErrorResponse(w, 500, err.Error())
		return
	}
	responses.CustomResponse(w, orders, 200)
}

func (api *APIHandler) addOrder(w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserFromContext(r.Context())

	var requestOrder *models.Order
	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "Error Parsing Request")
		return
	}

	time := api.Config.Game.StartDate
	err := api.DB.AddOrder(*requestOrder, user.UserID, time.Time)
	if err != nil {
		responses.ErrorResponse(w, 500, "Order couldn't be added")
	}
}

func (api *APIHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	var requestOrder *models.Order
	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "Error Parsing Request")
		return
	}

	time := api.Config.Game.StartDate
	err := api.DB.CancelOrder(requestOrder.OrderID, time.Time)
	if err != nil {
		responses.ErrorResponse(w, 500, "Order couldn't be deleted")
		return
	}

	responses.CustomResponse(w, requestOrder, 200)
}

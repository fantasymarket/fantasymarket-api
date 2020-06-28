package v1

import (
	"encoding/json"
	"fantasymarket/database/models"
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/utils/http/responses"
	"io/ioutil"
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
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}

	allOrders := &customOrder{Order: models.Order{UserID: user.UserID}}
	err = yaml.Unmarshal(body, allOrders)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, decodingError.Error())
		return
	}

	if allOrders.Limit <= 0 || allOrders.Limit > 20 {
		allOrders.Limit = 20
	}

	orders, err := api.DB.GetOrders(allOrders.Order, allOrders.Limit, allOrders.Offset)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
		return
	}
	responses.CustomResponse(w, orders, 200)

}

func (api *APIHandler) getOrdersID(w http.ResponseWriter, r *http.Request) {
	var requestOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, decodingError.Error())
		return
	}

	orders, err := api.DB.GetOrderByID(requestOrder.OrderID)
	if err != nil {
		responses.ErrorResponse(w, 500, fetchingError.Error())
		return
	}
	responses.CustomResponse(w, orders, 200)
}

func (api *APIHandler) addOrder(w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserFromContext(r.Context())

	var requestOrder *models.Order
	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, decodingError.Error())
		return
	}

	time := api.Config.Game.StartDate
	err := api.DB.AddOrder(*requestOrder, user.UserID, time)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, orderUpdateError.Error())
	}
}

func (api *APIHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	var requestOrder *models.Order
	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, decodingError.Error())
		return
	}

	time := api.Config.Game.StartDate
	err := api.DB.CancelOrder(requestOrder.OrderID, time)
	if err != nil {
		responses.ErrorResponse(w, 500, orderDeletionError.Error())
		return
	}

	responses.CustomResponse(w, requestOrder, 200)
}

package v1

import (
	"fantasymarket/utils/http/responses"
	"net/http"
)

func (api *APIHandler) getEvents(w http.ResponseWriter, r *http.Request) {
	t := api.Config.Game.StartDate
	events, err := api.DB.GetEvents(t)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, fetchingError.Error())
	}

	responses.CustomResponse(w, events, 200)

}

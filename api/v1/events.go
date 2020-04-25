package v1

import (
	"fantasymarket/utils/http/responses"
	"net/http"
)

func (api *APIHandler) getEvents(w http.ResponseWriter, r *http.Request) {
	t := api.Config.Game.StartDate
	events, err := api.DB.GetEvents(t.Time)

	if err != nil {
		responses.ErrorResponse(w, "Events couldn't be fetched", http.StatusInternalServerError)
	}

	responses.CustomResponse(w, events, 200)

}
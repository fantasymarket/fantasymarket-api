package v1

import (
	"fantasymarket/utils/http/responses"
	"net/http"
	"time"
)

func (api *APIHandler) getEvents(w http.ResponseWriter, r *http.Request) {
	events, err := api.DB.GetEvents()

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're ginormously fucked")
	}

	responses.CustomResponse(w, events, 200)

}

func (api *APIHandler) getOverview(w http.ResponseWriter, r *http.Request) {

}

func (api *APIHandler) getTime(w http.ResponseWriter, r *http.Request) {
	t := api.Game.Options.StartDate

	if !t.IsZero() {
		responses.CustomResponse(w, t.Format(time.RFC3339), 200)
	} else {
		responses.ErrorResponse(w, http.StatusInternalServerError, "we're absolutely fucked")
	}

}

package v1

import (
	"fantasymarket/utils/http/responses"
	"net/http"
	"time"
)


func (api *APIHandler) getTime(w http.ResponseWriter, r *http.Request) {
	t := api.Config.Game.StartDate

	if !t.IsZero() {
		responses.CustomResponse(w, t.Format(time.RFC3339), 200)
	} else {
		responses.ErrorResponse(w, "we're absolutely fucked", http.StatusInternalServerError)
	}
}
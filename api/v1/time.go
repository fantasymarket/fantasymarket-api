package v1

import (
	"fantasymarket/utils/http/responses"
	"net/http"
	"time"
)

func (api *APIHandler) getTime(w http.ResponseWriter, r *http.Request) {
	t := api.Config.Game.StartDate

	if !t.IsZero() {
		responses.CustomResponse(w, map[string]interface{}{
			"timestamp":          api.Game.GetCurrentDate().Format(time.RFC3339),
			"ticksPerSecond":     api.Config.Game.TicksPerSecond,
			"gameSecondsPerTick": api.Config.Game.GameTimePerTick.Seconds(),
			"startDate":          api.Config.Game.StartDate,
			"ticksSinceStart":    api.Game.TicksSinceStart,
		}, 200)
		return
	}

	responses.ErrorResponse(w, http.StatusInternalServerError, errFetchingData.Error())
}

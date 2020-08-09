package v1

import (
	"fantasymarket/game/details"
	"fantasymarket/utils/http/responses"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v3"
)

func (api *APIHandler) getAllAssets(w http.ResponseWriter, r *http.Request) {
	allAssets, err := details.AssetsYamlBytes()
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errFetchingData.Error())
		return
	}

	m := make([]details.AssetDetails, 30)
	err = yaml.Unmarshal(allAssets, &m)

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errFetchingData.Error())
		return
	}

	responses.CustomResponse(w, m, 200)
}

func (api *APIHandler) getAsset(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	givenTime := chi.URLParam(r, "time")

	tick := api.Game.TicksSinceStart

	if givenTime != "" {
		var err error
		tick, err = api.getTickAtTime(givenTime)
		if err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
			return
		}
	}

	assetData, err := api.DB.GetAssetAtTick(symbol, tick)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
		return
	}

	asset, ok := api.Game.AssetDetails[symbol]
	if !ok {
		responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
		return
	}

	responses.CustomResponse(w, map[string]interface{}{
		"symbol":      asset.Symbol,
		"type":        asset.Type,
		"name":        asset.Name,
		"description": asset.Description,
		"tick":        tick,
		"date":        api.Game.TickToTime(tick),
		"price":       assetData.Index,
		"price24h":    0,
		"volume":      0,
		"volume24h":   0,
	}, http.StatusOK)
}

func (api *APIHandler) getAssetHistory(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	fromTime := chi.URLParam(r, "from")
	toTime := chi.URLParam(r, "to")

	from := int64(0)
	to := api.Game.TicksSinceStart

	var err error
	if fromTime != "" {
		from, err = api.getTickAtTime(fromTime)
	}

	if toTime != "" {
		to, err = api.getTickAtTime(toTime)
	}

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errInvalidParameters.Error())
		return
	}

	asset, err := api.DB.GetAssetData(symbol, from, to)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
		return
	}

	responses.CustomResponse(w, asset, http.StatusOK)
}

func (api *APIHandler) getTickAtTime(timestamp string) (int64, error) {
	startTime, err := time.Parse(time.RFC3339, api.Config.Game.StartDate.String())
	if err != nil {
		return 0, err
	}
	currentTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, err
	}
	difference := int64(currentTime.Sub(startTime).Hours())

	return difference, nil
}

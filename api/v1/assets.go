package v1

import (
	"fantasymarket/game/details"
	"fantasymarket/utils/http/responses"
	"net/http"
	"strconv"
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

type assetResponse struct {
	Symbol      string              `json:"symbol"`
	Type        string              `json:"type"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Tick        string              `json:"tick,omitempty"`
	Date        string              `json:"date,omitempty"`
	Price       string              `json:"price,omitempty"`
	Price24h    string              `json:"price24h,omitempty"`
	Volume      string              `json:"volume,omitempty"`
	Volume24h   string              `json:"volume24h,omitempty"`
	From        string              `json:"from,omitempty"`
	To          string              `json:"to,omitempty"`
	Prices      []map[string]string `json:"prices,omitempty"`
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

	responses.CustomResponse(w,
		assetResponse{
			Symbol:      asset.Symbol,
			Type:        asset.Type,
			Name:        asset.Name,
			Description: asset.Description,
			Tick:        strconv.FormatInt(tick, 10),
			Date:        api.Game.TickToTime(tick).Format(time.RFC3339),
			Price:       strconv.FormatInt(assetData.Index, 10),
			Price24h:    "",
			Volume:      "",
			Volume24h:   "",
		},
		http.StatusOK,
	)
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

	assetData, err := api.DB.GetAssetData(symbol, from, to)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
		return
	}

	prices := []map[string]string{}
	for _, a := range *assetData {
		prices = append(prices, map[string]string{
			"date":  api.Game.TickToTime(a.Tick).Format(time.RFC3339),
			"index": strconv.FormatInt(a.Index, 10),
		})
	}

	asset, ok := api.Game.AssetDetails[symbol]
	if !ok {
		responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
		return
	}

	responses.CustomResponse(w,
		assetResponse{
			Symbol:      asset.Symbol,
			Type:        asset.Type,
			Name:        asset.Name,
			Description: asset.Description,
			From:        api.Game.TickToTime(from).Format(time.RFC3339),
			To:          api.Game.TickToTime(to).Format(time.RFC3339),
			Prices:      prices,
		},
		http.StatusOK,
	)
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

package v1

import (
	"fantasymarket/utils/http/responses"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (api *APIHandler) getAllAssets(w http.ResponseWriter, r *http.Request) {

	givenTime := r.URL.Query().Get("time")
	tick := api.Game.TicksSinceStart

	fmt.Println(givenTime)
	fmt.Println(givenTime)
	fmt.Println(givenTime)
	fmt.Println(givenTime)
	fmt.Println(givenTime)
	fmt.Println(givenTime)

	if givenTime != "" {
		var err error
		tick, err = api.Game.TimeStringToTick(givenTime)

		if err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, errAssetNotFound.Error())
			return
		}
	}

	assetData, err := api.DB.GetAssetsAtTick(tick)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, errFetchingData.Error())
		return
	}

	var assets []assetResponse
	for _, asset := range assetData {
		assetDetails := api.Game.AssetDetails[asset.Symbol]

		assets = append(assets, assetResponse{
			Symbol:      asset.Symbol,
			Type:        assetDetails.Type,
			Name:        assetDetails.Name,
			Description: assetDetails.Description,
			Tick:        strconv.FormatInt(tick, 10),
			Date:        api.Game.TickToTime(tick).Format(time.RFC3339),
			Price:       strconv.FormatInt(asset.Index, 10),
			Volume:      "",
		})
	}

	responses.CustomResponse(w, assets, 200)
}

type assetResponse struct {
	Symbol      string              `json:"symbol"`
	Type        string              `json:"type"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Tick        string              `json:"tick,omitempty"`
	Date        string              `json:"date,omitempty"`
	Price       string              `json:"price,omitempty"`
	Volume      string              `json:"volume,omitempty"`
	From        string              `json:"from,omitempty"`
	To          string              `json:"to,omitempty"`
	Prices      []map[string]string `json:"prices,omitempty"`
}

func (api *APIHandler) getAsset(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	givenTime := r.URL.Query().Get("time")

	tick := api.Game.TicksSinceStart

	if givenTime != "" {
		var err error
		tick, err = api.Game.TimeStringToTick(givenTime)
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
			Volume:      "",
		},
		http.StatusOK,
	)
}

func (api *APIHandler) getAssetHistory(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	fromTime := r.URL.Query().Get("from")
	toTime := r.URL.Query().Get("to")

	from := int64(0)
	to := api.Game.TicksSinceStart

	var err error
	if fromTime != "" {
		from, err = api.Game.TimeStringToTick(fromTime)
	}

	if toTime != "" {
		to, err = api.Game.TimeStringToTick(toTime)
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

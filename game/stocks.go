package game

import (
	"fantasymarket/game/structs"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func loadStocks() (map[string]structs.StockSettings, error) {
	stockSettings := []structs.StockSettings{}

	stockData, err := ioutil.ReadFile("./game/stocks.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(stockData, &stockSettings); err != nil {
		return nil, err
	}

	stockSettingsMap := map[string]structs.StockSettings{}
	for _, stock := range stockSettings {
		stockSettingsMap[stock.Symbol] = stock
	}

	return stockSettingsMap, nil
}

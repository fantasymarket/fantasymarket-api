package details

import (
	"gopkg.in/yaml.v2"
)

// LoadStockDetails loads all stock details from the stocks.yaml
func LoadStockDetails() (map[string]StockDetails, error) {
	stockSettings := []StockDetails{}

	stockData, err := StocksYamlBytes()
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(stockData, &stockSettings); err != nil {
		return nil, err
	}

	stockSettingsMap := map[string]StockDetails{}
	for _, stock := range stockSettings {
		stockSettingsMap[stock.Symbol] = stock
	}

	return stockSettingsMap, nil
}

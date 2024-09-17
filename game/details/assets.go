package details

import (
	"gopkg.in/yaml.v2"
)

// LoadAssetDetails loads all asset details from the assets.yaml
func LoadAssetDetails() (map[string]AssetDetails, error) {
	assetSettings := []AssetDetails{}

	assetData, err := AssetsYamlBytes()
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(assetData, &assetSettings); err != nil {
		return nil, err
	}

	assetSettingsMap := map[string]AssetDetails{}
	for _, asset := range assetSettings {
		assetSettingsMap[asset.Symbol] = asset
	}

	return assetSettingsMap, nil
}

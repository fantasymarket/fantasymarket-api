package database

import (
	"fantasymarket/database/models"
	"fantasymarket/game/details"

	uuid "github.com/satori/go.uuid"
)

// CreateInitialAssets takes a list of initial assets and uses them to initialize the database
func (s *Service) CreateInitialAssets(assetDetails map[string]details.AssetDetails) error {

	for _, asset := range assetDetails {
		if err := s.DB.FirstOrCreate(
			&models.Asset{},
			&models.Asset{
				Symbol: asset.Symbol,
				Index:  asset.Index,
				Name:   asset.Name,
				Tick:   0,
				Volume: 0,
			},
		).Error; err != nil {
			return err
		}
	}
	return nil
}

// AddAsset adds a asset to the asset table
func (s *Service) AddAsset(asset models.Asset, tick int64) error {
	return s.DB.Create(&models.Asset{
		Symbol: asset.Symbol,
		Name:   asset.Name,
		Index:  asset.Index,
		Volume: asset.Volume,
		Tick:   tick,
	}).Error
}

// AddAssets adds a slice of assets to the asset table
func (s *Service) AddAssets(assets []models.Asset, tick int64) error {
	for _, asset := range assets {
		if err := s.AddAsset(asset, tick); err != nil {
			return err
		}
	}
	return nil
}

// GetNextTick retrieves the tick number for the next tick from the database,
// this is used to initialize our Service when the program restarts
func (s *Service) GetNextTick() (int64, error) {
	var asset models.Asset
	if err := s.DB.Table("assets").Select("tick").Order("tick desc").First(&asset).Error; err != nil {
		return 0, err
	}

	return asset.Tick + 1, nil
}

// GetAssetsAtTick fetches the value of all assets at a specific tick
func (s *Service) GetAssetsAtTick(lastTick int64) ([]models.Asset, error) {
	var assets []models.Asset
	if err := s.DB.Where(models.Asset{Tick: lastTick}).Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

// GetAssetMapAtTick fetches the value of all assets at a tick as a Map
func (s *Service) GetAssetMapAtTick(lastTick int64) (map[string]models.Asset, error) {
	assets, err := s.GetAssetsAtTick(lastTick)
	if err != nil {
		return nil, err
	}

	assetMap := map[string]models.Asset{}
	for _, asset := range assets {
		assetMap[asset.Symbol] = asset
	}

	return assetMap, nil
}

// GetAssetAtTick fetches the value of a specific asset at a specific tick
// `asset` can either be a asset symbol or assetID
func (s *Service) GetAssetAtTick(asset string, tick int64) (*models.Asset, error) {
	query := models.Asset{Tick: tick}

	if uuid, err := uuid.FromString(asset); err == nil {
		query.AssetID = uuid
	} else {
		query.Symbol = asset
	}

	var assetData models.Asset
	if err := s.DB.Where(query).First(&assetData).Error; err != nil {
		return nil, err
	}

	return &assetData, nil
}

// GetAssetData fetches all ticks between two ticks
func (s *Service) GetAssetData(asset string, from int64, to int64) (*[]models.Asset, error) {
	query := models.Asset{}

	if uuid, err := uuid.FromString(asset); err == nil {
		query.AssetID = uuid
	} else {
		query.Symbol = asset
	}

	assets := []models.Asset{}
	if err := s.DB.Where(query).Where(
		"tick BETWEEN ? AND ?",
		from,
		to,
	).Find(&assets).Error; err != nil {
		return nil, err
	}

	return &assets, nil
}

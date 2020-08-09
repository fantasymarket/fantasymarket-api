package database_test

import (
	"fantasymarket/database/models"
	"fantasymarket/game/details"

	"github.com/stretchr/testify/assert"
)

type CreateInitialAssetsTestData struct {
	asset       map[string]details.AssetDetails
	expectation models.Asset
}

var testCreateInitialAssetsData = []CreateInitialAssetsTestData{
	{
		asset: map[string]details.AssetDetails{"HELLO": {
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
		},
		},
		expectation: models.Asset{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
		},
	},
	{
		asset: map[string]details.AssetDetails{"": {
			Symbol: "",
			Index:  401,
			Name:   "Not Hello Asset",
		},
		},
		expectation: models.Asset{
			Symbol: "",
			Index:  401,
			Name:   "Not Hello Asset",
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestCreateInitialAssets() {

	var assets []models.Asset
	for i, test := range testCreateInitialAssetsData {
		err := suite.dbService.CreateInitialAssets(test.asset)
		assert.Equal(suite.T(), nil, err)
		err = suite.dbService.DB.Find(&assets).Error
		assert.Equal(suite.T(), nil, err)

		if test.expectation.Symbol != "" {
			//Again.., I hate it
			assert.Equal(suite.T(), test.expectation.Symbol, assets[i].Symbol)
			assert.Equal(suite.T(), test.expectation.Index, assets[i].Index)
			assert.Equal(suite.T(), test.expectation.Name, assets[i].Name)
		}
	}

	suite.dbService.DB.Close()
}

type AddAssetTestData struct {
	asset       models.Asset
	expectation models.Asset
}

var testAddAssetData = []AddAssetTestData{
	{
		asset: models.Asset{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
		},
		expectation: models.Asset{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
		},
	},
	{
		asset: models.Asset{
			Symbol: "",
			Index:  100,
			Name:   "Hello Asset",
		},
		expectation: models.Asset{
			Symbol: "",
			Index:  100,
			Name:   "Hello Asset",
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestAddAsset() {
	for _, test := range testAddAssetData {
		err := suite.dbService.AddAsset(test.asset, 1)
		assert.Equal(suite.T(), nil, err)
		assert.Equal(suite.T(), false, suite.dbService.DB.Where("symbol = ?", test.asset.Symbol).Find(&models.Asset{}).RecordNotFound())
	}
	suite.dbService.DB.Close()
}

type GetNextTickTestData struct {
	asset       models.Asset
	expectation int64
}

var testGetNextTickData = []GetNextTickTestData{
	{
		asset: models.Asset{
			Tick: 0,
		},
		expectation: 1,
	},
	{
		asset:       models.Asset{},
		expectation: 1,
	},
	{
		expectation: 1,
	},
}

func (suite *DatabaseTestSuite) TestGetNextTick() {
	for _, test := range testGetNextTickData {
		err := suite.dbService.DB.Create(&test.asset).Error
		assert.Equal(suite.T(), nil, err)
		result, err := suite.dbService.GetNextTick()
		assert.Equal(suite.T(), nil, err)
		assert.Equal(suite.T(), test.expectation, result)
	}
	suite.dbService.DB.Close()
}

type GetAssetsAtTickTestData struct {
	tick        int64
	asset       models.Asset
	expectation []models.Asset
}

var testGetAssetsAtTickData = []GetAssetsAtTickTestData{
	{
		tick: 1,
		asset: models.Asset{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
			Tick:   1,
		},
		expectation: []models.Asset{
			{
				Symbol: "HELLO",
				Index:  100,
				Name:   "Hello Asset",
				Tick:   1,
			},
		},
	},
	{
		tick: 2,
		asset: models.Asset{
			Symbol: "NOTHEL",
			Index:  100,
			Name:   "Not Hello Asset",
			Tick:   2,
		},
		expectation: []models.Asset{
			{
				Symbol: "HELLO",
				Index:  100,
				Name:   "Hello Asset",
				Tick:   2,
			},
			{
				Symbol: "NOTHEL",
				Index:  100,
				Name:   "Not Hello Asset",
				Tick:   2,
			},
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestGetAssetsAtTick() {
	for _, test := range testGetAssetsAtTickData {
		err := suite.dbService.DB.Create(&test.asset).Error
		assert.Equal(suite.T(), nil, err)
		test.asset.Tick++
		err = suite.dbService.DB.Create(&test.asset).Error
		assert.Equal(suite.T(), nil, err)
		result, err := suite.dbService.GetAssetsAtTick(test.tick)
		assert.Equal(suite.T(), nil, err)

		for j := 0; j < len(test.expectation); j++ {
			assert.Equal(suite.T(), test.expectation[j].Symbol, result[j].Symbol)
			assert.Equal(suite.T(), test.expectation[j].Index, result[j].Index)
			assert.Equal(suite.T(), test.expectation[j].Name, result[j].Name)
			assert.Equal(suite.T(), test.expectation[j].Tick, result[j].Tick)
		}
	}
	suite.dbService.DB.Close()
}

type GetAssetMapAtTickTestData struct {
	tick        int64
	asset       models.Asset
	expectation map[string]models.Asset
}

var testGetAssetMapAtTickData = []GetAssetMapAtTickTestData{
	{
		tick: 1,
		asset: models.Asset{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
			Tick:   1,
		},
		expectation: map[string]models.Asset{
			"HELLO": {
				Symbol: "HELLO",
				Index:  100,
				Name:   "Hello Asset",
				Tick:   1,
			},
		},
	},
	{
		tick: 2,
		asset: models.Asset{
			Symbol: "NOTHEL",
			Index:  100,
			Name:   "Not Hello Asset",
			Tick:   2,
		},
		expectation: map[string]models.Asset{
			"HELLO": {
				Symbol: "HELLO",
				Index:  100,
				Name:   "Hello Asset",
				Tick:   2,
			},
			"NOTHEL": {
				Symbol: "NOTHEL",
				Index:  100,
				Name:   "Not Hello Asset",
				Tick:   2,
			},
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestGetAssetMapAtTick() {
	for _, test := range testGetAssetMapAtTickData {
		err := suite.dbService.DB.Create(&test.asset).Error
		assert.Equal(suite.T(), nil, err)
		test.asset.Tick++
		err = suite.dbService.DB.Create(&test.asset).Error
		assert.Equal(suite.T(), nil, err)
		result, err := suite.dbService.GetAssetMapAtTick(test.tick)
		assert.Equal(suite.T(), nil, err)

		for j := 0; j < len(test.expectation); j++ {
			assert.Equal(suite.T(), test.expectation[test.asset.Symbol].Symbol, result[test.asset.Symbol].Symbol)
			assert.Equal(suite.T(), test.expectation[test.asset.Symbol].Index, result[test.asset.Symbol].Index)
			assert.Equal(suite.T(), test.expectation[test.asset.Symbol].Name, result[test.asset.Symbol].Name)
			assert.Equal(suite.T(), test.expectation[test.asset.Symbol].Tick, result[test.asset.Symbol].Tick)
		}
	}
	suite.dbService.DB.Close()
}

type GetAssetAtTickTestData struct {
	tick        int64
	asset       models.Asset
	assetName   string
	needsToFail bool
}

var testGetAssetAtTickData = []GetAssetAtTickTestData{
	{
		tick:      1,
		assetName: "HELLO",
		asset: models.Asset{
			Symbol: "HELLO",
			Index:  100,
			Name:   "Hello Asset",
			Tick:   1,
		},
	}, {
		tick:      99999,
		assetName: "hi",
		asset: models.Asset{
			Symbol: "hi",
			Index:  22,
			Name:   "Hello Asset",
			Tick:   0,
		},
		needsToFail: true,
	},
	{
		tick:      2,
		assetName: "insert-uuid",
		asset: models.Asset{
			Symbol: "NOTHEL",
			Index:  100,
			Name:   "Not Hello Asset",
			Tick:   2,
		},
	},
}

func (suite *DatabaseTestSuite) TestGetAssetAtTick() {
	assert := suite.Assert()

	for _, test := range testGetAssetAtTickData {

		if err := suite.dbService.DB.Create(&test.asset).Error; err != nil {
			assert.Fail(err.Error())
			return
		}

		assetName := test.assetName
		if assetName == "insert-uuid" {
			assetName = test.asset.AssetID.String()
		}

		result, err := suite.dbService.GetAssetAtTick(assetName, test.tick)
		if test.needsToFail && assert.Error(err, "returned asset instead of error") {
			return
		}

		if err != nil {
			assert.Fail(err.Error())
			return
		}

		assert.Equal(test.asset.Symbol, result.Symbol)
		assert.Equal(test.asset.Index, result.Index)
		assert.Equal(test.asset.Name, result.Name)
		assert.Equal(test.asset.Tick, result.Tick)
	}

	suite.dbService.DB.Close()
}

package database_test

import (
	"fantasymarket/database"
	"fantasymarket/database/models"
	"fantasymarket/game/details"
	"fantasymarket/utils/config"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Connect connects to the database and returns thedatabase object
func Connect(config *config.Config) (*database.Service, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	database.AutoMigrate(db)

	log.Info().Msg("successfully connected to the database")

	return &database.Service{
		DB:     db,
		Config: config,
	}, nil
}

type DatabaseTestSuite struct {
	suite.Suite
	dbService *database.Service
}

var testAddEventData = []details.EventDetails{
	{
		EventID: "testEvent1",
	},
	{},
}

func (suite *DatabaseTestSuite) SetupTest() {
	var err error
	suite.dbService, err = Connect(&config.Config{})
	if err != nil {
		panic(err)
	}
}

func (suite *DatabaseTestSuite) TestAddEvent() {

	createdAt := parseTime("2019-12-30T15:04:05Z")
	currentDate := parseTime("2020-12-30T15:04:05Z")
	for _, event := range testAddEventData {
		suite.dbService.AddEvent(event, createdAt)
		err := suite.dbService.DB.Where(models.Event{
			Active: true,
		}).Where("created_at < ?", currentDate).Find(&models.Event{}).Error
		assert.Equal(suite.T(), nil, err)

		suite.dbService.DB.Delete(models.Event{})
	}

	suite.dbService.DB.Close()
}

var testRemoveEventData = []details.EventDetails{
	{
		EventID: "testEvent1",
	},
	{},
}

func (suite *DatabaseTestSuite) TestRemoveEvent() {

	createdAt := parseTime("2019-12-30T15:04:05Z")

	for _, event := range testAddEventData {
		suite.dbService.DB.Create(&models.Event{
			EventID:   event.EventID,
			Title:     event.Title,
			Text:      event.Description,
			Active:    true,
			CreatedAt: createdAt,
		})
		err := suite.dbService.RemoveEvent(event.EventID)
		assert.Equal(suite.T(), nil, err)
	}

	suite.dbService.DB.Close()
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

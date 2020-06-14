package database_test

import (
	"fantasymarket/database"
	"fantasymarket/utils/config"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
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

func (suite *DatabaseTestSuite) SetupTest() {
	var err error
	suite.dbService, err = Connect(&config.Config{})
	if err != nil {
		panic(err)
	}
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

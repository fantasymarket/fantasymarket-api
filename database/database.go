package database

import (
	"github.com/rs/zerolog/log"

	"fantasymarket/database/models"
	"fantasymarket/utils/config"

	"github.com/jinzhu/gorm"

	// load sqlite dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Service is the Database Service
type Service struct {
	DB     *gorm.DB // gorm database instance
	Config *config.Config
}

// Connect connects to the database and returns thedatabase object
func Connect(config *config.Config) (*Service, error) {
	db, err := gorm.Open("sqlite3", "database.db")

	if err != nil {
		return nil, err
	}

	AutoMigrate(db)

	log.Info().Msg("successfully connected to the database")

	return &Service{
		DB:     db,
		Config: config,
	}, nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.Stock{},
		&models.Event{},
		&models.Order{},
		&models.User{},
		&models.Portfolio{},
		&models.PortfolioItem{},
	)
}

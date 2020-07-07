package database

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"fantasymarket/database/models"
	"fantasymarket/utils/config"

	"github.com/jinzhu/gorm"

	// load sqlite dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Service is the Database Service
type Service struct {
	DB     *gorm.DB // gorm database instance
	Config *config.Config
}

// Connect connects to the database and returns thedatabase object
func Connect(config *config.Config) (*Service, error) {

	var db *gorm.DB
	var err error

	if config.Database.Type == "sqlite" {
		// SQLITE
		db, err = gorm.Open("sqlite3", "database.db")

	} else if config.Database.Type == "postgres" {
		// POSTGRESQL

		addr := config.Database.URL
		if addr == "" {
			format := "host=%s port=%s dbname=%s user=%s sslmode=%s"
			addr = fmt.Sprintf(
				format,
				config.Database.Host,
				config.Database.Port,
				config.Database.Database,
				config.Database.Username,
				config.Database.SSL,
			)
		}
		db, err = gorm.Open("postgres", addr)
	} else {
		panic("error: unknown database type: " + config.Database.Type)
	}

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

// AutoMigrate migrates the database tables
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

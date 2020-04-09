package main

import (
	"fantasymarket/api"
	"fantasymarket/database"
	"fantasymarket/game"
	"fantasymarket/utils/config"

	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("error while loading the configuration")
	}

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal().Err(err).Msg("error while connecting to the database")
	}

	game, err := game.Start(db, config)
	if err != nil {
		log.Fatal().Err(err).Msg("error while connecting starting the game")
	}

	api.Start(db, game, config)
}

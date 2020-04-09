package main

import (
	"fantasymarket/api"
	"fantasymarket/database"
	"fantasymarket/game"
)

func main() {
	db, _ := database.Connect()

	game, err := game.Start(db)
	if err != nil {
		panic(err)
	}

	api.Start(db, game)
}

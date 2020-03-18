package main

import (
	"fantasymarket/api"
	"fantasymarket/database"
	"fantasymarket/game"

	"time"
)

const (
	// Minute is the duration of 60 seconds
	Minute = time.Second * 60
	// Hour is the duration of 60 minutes
	Hour = time.Second * 60 * 60
	//Day is the duration o 24 hours
	Day = Hour * 24
)

func main() {
	db, _ := database.Connect()

	api.Start(db)
	game.Start(db)
}

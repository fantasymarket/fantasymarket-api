package main

import (
	// "fantasymarket/database"
	"fantasymarket/game"
	"time"
)

/// Was wir noch nehmen können
/// SQL:			https://github.com/jmoiron/sqlx
/// Decimal:	https://github.com/shopspring/decimal

const (
	// Minute is the duration of 60 seconds
	Minute = time.Second * 60
	// Hour is the duration of 60 minutes
	Hour = time.Second * 60 * 60
	//Day is the duration o 24 hours
	Day = Hour * 24
)

func main() {
	game.MainStocks()
	// database.DatabaseMain()
}

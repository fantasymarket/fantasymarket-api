package mock_data

import (
	"time"
)

type UserStats struct {
	name string
	ownedStocks []string
}

type Stocks struct {
	Name string
	Index int64
	Trend int64
}

type Orders struct {
	orders []string

}

type News struct {
	headline string
	date time.Time
}
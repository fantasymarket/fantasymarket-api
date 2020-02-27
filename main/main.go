package main

import (
	"math/rand"
	"time"
)

/// Was wir noch nehmen können
/// SQL:			https://github.com/jmoiron/sqlx
/// Decimal:		https://github.com/shopspring/decimal

//FantasyMarketOptions manages the Options of the programm
type FantasyMarketOptions struct {
	TicksPerSecond  float64       // How many times the game updates per second
	GameTimePerTick time.Duration // How much ingame time passes between updates
	StartDate       time.Time     // The initial ingame time
}

type Service struct {
	//DB Database
	Options FantasyMarketOptions
	Stocks  []Stock
	Events  []Event
}

// Stock is the Stock "Class"
type Stock struct {
	ID     string // Stock Symbol e.g GOOG
	Name   string // Stock NAme e.g Alphabet Inc.
	Index  int64  // Price per share
	Shares int64  // Number per share
	Tags   map[string]bool
}

// Events happen randomly every game tick
type Event struct {
	Name                 string
	MinTimeBetweenEvents time.Duration
	Chance               float64 // 0 - 1

	// Stuff that affects all tags
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	Duration    time.Duration // Time during which the event is the event is run every tick
	TimeCreated time.Time

	// We use tags if we only want to affect only specific stocks
	Tags map[string]TagOptions
}

//TagOptions are more indepth settings for specific event tags only
type TagOptions struct {
	Trend int64 // Note: .2 would be 20 and .02 would be 2
	//// TimeOffset time.Duration // Optionally offset the event to e.g only affect a tag after x time
	////Duration time.Duration
}

const (
	// Minute is the duration of 60 seconds
	Minute = time.Second * 60
	// Hour is the duration of 60 minutes
	Hour = time.Second * 60 * 60
	//Day is the duration o 24 hours
	Day = Hour * 24
)

func MainStocks() {
	stocks := []Stock{
		{
			ID:    "GOOG",
			Index: int64(10000),
			Tags:  map[string]bool{"tech": true, "intl": true},
		},
		{
			ID:    "FRIZ",
			Index: int64(13969),
			Tags:  map[string]bool{"food": true, "local": true},
		},
		{
			ID:    "LMAO",
			Index: int64(13969),
			Tags:  map[string]bool{"arthur": true, "henry": true},
		},
	}

	events := []Event{
		{Name: "Virus in Seattle", Tags: map[string]TagOptions{"tech": {Trend: 20}, "china": {Trend: 20}}},
		{Name: ".com bubble Crash", Tags: map[string]TagOptions{"tech": {Trend: 20}, "china": {Trend: 20}}},
	}

	s := Service{
		Options: FantasyMarketOptions{
			TicksPerSecond:  0.1,
			StartDate:       time.Now(),
			GameTimePerTick: time.Hour,
		},
		Stocks: stocks,
		Events: events,
	}
	rand.Seed(time.Now().UnixNano())

	go StartLoop(s)
}

// startLoop startsrunningticks indefinitly
func StartLoop(s Service) {

	// We need to calculatre the current game date
	startDate := s.Options.StartDate
	gameTimePerTick := s.Options.GameTimePerTick
	ticksSinceStart := 0 // TODO persist this number so it doesnt reset after restarting the program
	dateNow := startDate.Add(gameTimePerTick * time.Duration(ticksSinceStart))

	for {
		Tick(s, dateNow)

		// Sleep for the duration of a single tick (Since we want 1 tick in 10 Seconds)
		time.Sleep(time.Duration(1 / s.Options.TicksPerSecond))

		// Adding 1 hour every tick(Update) (10 seconds when TicksPerSecond=0.1 ) onto the previously defined Date time
		dateNow = dateNow.Add(gameTimePerTick)
		ticksSinceStart++
	}
}

// tick is updating the current state of our system
func Tick(s Service, dateNow time.Time) {
	// TODO: Get currently Running Events
	// TODO: Stop Events that are over the max duration
	e := s.Events
	for i := 0; i < len(e); i++ {

		endDate := e[i].TimeCreated.Add(e[i].Duration)
		if !dateNow.Before(endDate) {
			// TODO: remove event
		}
	}
	// TODO: Randomly add new Events to the list of running events that are currently valid (e.g min time between events) @Andre
	// TODO: Filter Only Currently relevant events @Andre
	// TODO: Run all events on the stocks @Arthur
	// TODO: Update Database @Andre
	// TODO: Update Orderbook @Arthur Andre
}

func isAffected(e Event, stock Stock) (string, bool) {
	for tag := range stock.Tags {
		if _, ok := e.Tags[tag]; ok {

			return tag, true
		}
	}
	return "", false
}

func ComputeStockNumbers(stocks []Stock, e Event) {
	max := 1
	min := -1
	for i := 0; i < len(stocks); i++ {
		eventTendency := int64(1)
		tag, flag := isAffected(e, stocks[i])
		if flag {
			eventTendency = int64(e.Tags[tag].Trend)

		}

		tendency := int64(rand.Intn(max-min) + min) // Range of -2 to 2
		stocks[i].Index = tendency * eventTendency
	}
}

// |func mockGraph() {
// 	count = 100,
// 	stock = { index: 173.43, trend: 1.5, fluctuation: 2 },
// 	startDate = new Date('2019-04-11'), } = {}) => {
// 	const a = [];
// 	const s = stock;

// 	for (let i = 0; i < count; i++) {
// 		const x = (Math.random() - 0.5) * s.fluctuation + (i / 1000) * s.trend;
// 		s.index += x;
// 		a.push({
// 			time: addDays(startDate, i).toDateString(),
// 			value: roundTo(0.01, s.index),
// 		});
// 	}

// 	return a;
// }

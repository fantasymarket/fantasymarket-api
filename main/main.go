package main

import (
	"bufio"
	"fantasymarket/utils"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

/// Was wir noch nehmen können
/// SQL:			https://github.com/jmoiron/sqlx
/// Decimal:	https://github.com/shopspring/decimal

//FantasyMarketOptions manages the Options of the programm
type FantasyMarketOptions struct {
	TicksPerSecond  float64       // How many times the game updates per second
	GameTimePerTick time.Duration // How much ingame time passes between updates
	StartDate       time.Time     // The initial ingame time
}

type Product struct {
	gorm.Model
	//db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
	Code  string
	Price uint
}
type Service struct {
	Options FantasyMarketOptions
	Stocks  []Stock
	Events  []Event
}

// Stock is the Stock "Class"
type Stock struct {
	ID        string // Stock Symbol e.g GOOG
	Name      string // Stock NAme e.g Alphabet Inc.
	Index     int64  // Price per share
	Shares    int64  // Number per share
	Tags      map[string]bool
	Stability int64 // Shows how many fluctuations the stock will have
	Trend     int64 // Shows the generall trend of the Stock

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

func main() {
	MainStocks()

}

func MainStocks() {
	stocks := []Stock{
		{
			ID:        "GOOG",
			Index:     int64(10530),
			Tags:      map[string]bool{"tech": true, "intl": true},
			Stability: 1,
			Trend:     1,
		},
		{
			ID:        "FRIZ",
			Index:     int64(13969),
			Tags:      map[string]bool{"food": true, "local": true},
			Stability: 1,
			Trend:     1,
		},
		{
			ID:        "LMAO",
			Index:     int64(12001),
			Tags:      map[string]bool{"arthur": true, "henry": true},
			Stability: 1,
			Trend:     1,
		},
	}

	events := []Event{
		{Name: "Virus in Seattle", Tags: map[string]TagOptions{"tech": {Trend: 2}, "china": {Trend: 2}}},
		{Name: ".com bubble Crash", Tags: map[string]TagOptions{"tech": {Trend: 2}, "china": {Trend: 2}}},
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

	go startLoop(s)
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// startLoop startsrunningticks indefinitly
func startLoop(s Service) {

	// We need to calculatre the current game date
	startDate := s.Options.StartDate
	gameTimePerTick := s.Options.GameTimePerTick
	ticksSinceStart := 0 // TODO persist this number so it doesnt reset after restarting the program
	dateNow := startDate.Add(gameTimePerTick * time.Duration(ticksSinceStart))

	for {
		tick(s, dateNow)

		// Sleep for the duration of a single tick (Since we want 1 tick in 10 Seconds)
		time.Sleep(time.Duration(1/s.Options.TicksPerSecond) * time.Second)

		// Adding 1 hour every tick(Update) (10 seconds when TicksPerSecond=0.1 ) onto the previously defined Date time
		dateNow = dateNow.Add(gameTimePerTick)
		ticksSinceStart++
	}
}

// tick is updating the current state of our system
func tick(s Service, dateNow time.Time) {

	// TODO: Get currently Running Events
	// TODO: Stop Events that are over the max duration
	e := s.Events
	for i := 0; i < len(e); i++ {

		endDate := e[i].TimeCreated.Add(e[i].Duration) //Calculate the endDate by adding the Duration to the time created
		if !dateNow.Before(endDate) {                  //Check if the current date is after the end date.
			// TODO: remove event
		}
	}

	ComputeStockNumbers(s.Stocks, s.Events)
	// TODO: Randomly add new Events to the list of running events that are currently valid (e.g min time between events) @Andre
	// TODO: Filter Only Currently relevant events @Andre
	// TODO: Run all events on the stocks @Arthur
	// TODO: Update Database @Andre
	// TODO: Update Orderbook @Arthur Andre

	//API:  /user/stats 			 => Worst/Best performer, Total Value, Portfolio Trend
	//		/news?from=XXX&to=XXXX 	 => Active Events
	//		/stocks					 => Overview of stocks
	//		/stocks?from=XXX&to=XXXX => All stocks in Time range
	//		/stocks/GOOG			 => Specific Stock request

	//Events:	Events have tags: Fixed, Recurring, Random
	//			Hardcoded Events => Elections, Olympic Games etc
	//			Definate Date Events (Moon Landing 1969)?
}

func isAffected(e []Event, stock Stock) (int64, bool) {
	for i := 0; i < len(e); i++ {
		for tag := range stock.Tags {
			if _, ok := e[i].Tags[tag]; ok {
				return e[i].Tags[tag].Trend, true
			}
		}
	}

	return 1, false
}

func ComputeStockNumbers(stocks []Stock, e []Event) {
	eventTendency := int64(1)
	for i := 0; i < len(stocks); i++ {
		trend, flag := isAffected(e, stocks[i])
		if flag {
			eventTendency = int64(trend)
		}

		tendency := getTendancy(stocks[i]) // Range of -2 to 2
		stocks[i].Index = tendency * eventTendency
		fmt.Println("Name: ", stocks[i].ID, "Index: ", stocks[i].Index)
	}

}

func getTendancy(s Stock) int64 {
	// 	const x = (Math.random() - 0.5) * s.fluctuation + (i / 1000) * s.trend;
	n := 5
	return utils.RandInt64(-n, n)*s.Stability + (s.Index/1000)*s.Trend
}

//func mockGraph() {
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

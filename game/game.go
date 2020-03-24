package game

import (
	"fantasymarket/database"
	"fantasymarket/database/models"
	"fantasymarket/utils"
	"fmt"
	"strconv"
	"time"
)

type GameService struct {
	DB              *database.DatabaseService
	Options         FantasyMarketOptions
	StockSettings   map[string]StockSettings
	EventSettings   map[string]EventSettings
	TicksSinceStart int64
}

//FantasyMarketOptions manages the Options of the programm
type FantasyMarketOptions struct {
	TicksPerSecond  float64       // How many times the game updates per second
	GameTimePerTick time.Duration // How much ingame time passes between updates
	StartDate       time.Time     // The initial ingame time
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// mention this for the assessment - clean code plus points
func Start(db *database.DatabaseService) (*GameService, error) {

	stockSettings := map[string]StockSettings{}
	eventSettings := map[string]EventSettings{}

	db.AddStockToTable(db.CreateStockForTest("GOOG", "Google", 10000, 1), 0)
	db.AddStockToTable(db.CreateStockForTest("APPL", "Apple Inc", 10000, 2), 0)

	s := &GameService{
		Options: FantasyMarketOptions{
			TicksPerSecond:  0.1,
			StartDate:       time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			GameTimePerTick: time.Hour,
		},
		StockSettings: stockSettings,
		EventSettings: eventSettings,
		DB:            db,
	}

	go startLoop(s)
	fmt.Println("stated game loop :o")

	return s, nil
}

// startLoop startsrunningticks indefinitly
func startLoop(s *GameService) {

	// We need to calculatre the current game date
	startDate := s.Options.StartDate
	gameTimePerTick := s.Options.GameTimePerTick
	s.TicksSinceStart = int64(0) // TODO persist this number so it doesnt reset after restarting the program
	dateNow := startDate.Add(gameTimePerTick * time.Duration(s.TicksSinceStart))

	for {
		s.TicksSinceStart++
		s.tick(dateNow)

		// Sleep for the duration of a single tick (Since we want 1 tick in 10 Seconds)
		time.Sleep(time.Duration(1/s.Options.TicksPerSecond) * time.Second)

		// Adding 1 hour every tick(Update) (10 seconds when TicksPerSecond=0.1 ) onto the previously defined Date time
		dateNow = dateNow.Add(gameTimePerTick)
	}
}

// tick is updating the current state of our system
func (s *GameService) tick(dateNow time.Time) {
	// TODO: Get currently Running Events from the database (models.Event)

	fmt.Println("[running tick:  " + strconv.FormatInt(s.TicksSinceStart, 10) + "]")

	currentlyRunningEvents, _ := s.DB.GetEvents()                      // Sub this for the DB query results
	lastStockIndexes, _ := s.DB.GetStocksAtTick(s.TicksSinceStart - 1) // Sub this for the DB query results

	fmt.Println("last stocks: ", lastStockIndexes)

	s.checkEventStillGoing(currentlyRunningEvents, dateNow)

	// TODO: Stop Events that are over the max duration

	// TODO: Randomly add new Events to the list of running events that are currently valid (e.g min time between events) @Andre
	// TODO: Filter Only Currently relevant events @Andre
	// TODO: Run all events on the stocks @Arthur
	// TODO: Update Orderbook @Arthur Andre

	// Events: Events have tags: Fixed, Recurring, Random
	//    Hardcoded Events => Elections, Olympic Games etc
	//    Definate Date Events (Moon Landing 1969)?

	s.ComputeStockNumbers(lastStockIndexes, currentlyRunningEvents, dateNow)
	// saveNextStockIndexes()
}

// checkEventStillGoing calculates if the duration of the event is over and then removes the event
func (s GameService) checkEventStillGoing(e []models.Event, dateNow time.Time) {
	for i := 0; i < len(e); i++ {
		endDate := e[i].CreatedAt.Add(e[i].Duration) // Calculate the endDate by adding the Duration to the time created
		if !dateNow.Before(endDate) {                // Check if the current date is after the end date.
			s.DB.RemoveEvent(e[i].EventID)
		}
	}
}

func (s GameService) getEventAffectedness(e []models.Event, stock models.Stock) int64 {

	affectedness := int64(0)
	for _, event := range e {

		eventSettings := s.EventSettings[event.EventID]

		for tag := range eventSettings.Tags {
			if _, ok := eventSettings.Tags[tag]; ok {
				affectedness += eventSettings.Tags[tag].Trend
			}
		}
	}

	return affectedness
}

func (s GameService) ComputeStockNumbers(stocks []models.Stock, e []models.Event, dateNow time.Time) {

	// This computes the random and own stock, not taking into account other peoples selling
	// As a stock drops to a % of its value, theres gonna be more buyers or more sellers
	for _, stock := range stocks {
		stock.Index += s.GetTendency(stock, s.getEventAffectedness(e, stock), dateNow) // Range of -2 to 2
		fmt.Println("Name: ", stock.StockID, "Index: ", stock.Index)
		s.DB.AddStockToTable(stock, s.TicksSinceStart)
	}
	fmt.Println("-----------------------------")
}

func (s GameService) GetTendency(stock models.Stock, affectedness int64, dateNow time.Time) int64 {
	const n int64 = 10
	stockSettings := s.StockSettings[stock.StockID]
	// Old Index: 10000, Stability: 1, Trend: -1
	// Rand(-10,10) * 1 + (10000/2000)*1 + (10000/10000)*-1
	// (3)*1 + (5)*1 + (1)*-1
	// 3 + 5 - 1
	// 7
	// 10000 + 7
	// New Index: 10007

	randomModifier := utils.RandInt64(-n, n, dateNow.UnixNano())
	fmt.Println("randomModifier: " + strconv.FormatInt(randomModifier, 10))
	stockTrend := (stock.Index / 2000) * stockSettings.Trend
	fmt.Println("stockTrend: " + strconv.FormatInt(stockTrend, 10))
	eventTrend := (stock.Index / 10000) * affectedness
	fmt.Println("eventTrend: " + strconv.FormatInt(eventTrend, 10))
	return randomModifier*stockSettings.Stability + stockTrend + eventTrend
	// Stability indicates how strong the random aspect is evaluated in comparison to the trend
}

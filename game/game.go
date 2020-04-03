package game

import (
	"errors"
	"fantasymarket/database"
	"fantasymarket/database/models"
	"fantasymarket/game/structs"
	"fantasymarket/utils/hash"
	"fmt"
	"strconv"
	"time"
)

// Service is the GameService
type Service struct {
	EventList       map[string]structs.EventSettings
	StockList       map[string]structs.StockSettings
	DB              *database.Service
	Options         FantasyMarketOptions
	TicksSinceStart int64
}

// FantasyMarketOptions manages the Options of the programm
type FantasyMarketOptions struct {
	TicksPerSecond  float64       // How many times the game updates per second
	GameTimePerTick time.Duration // How much ingame time passes between updates
	StartDate       time.Time     // The initial ingame time
}

// Start starts the game loop
func Start(db *database.Service) (*Service, error) {

	loadedStocks, err := loadStocks()
	if err != nil {
		return nil, err
	}

	if err := db.CreateInitialStocks(loadedStocks); err != nil {
		return nil, err
	}

	// TODO: right now, this map is empty
	loadedEvents, err := loadEvents()
	if err != nil {
		return nil, err
	}

	// TODO: Take all Fixed events and map them where Key is startDate and value true?

	s := &Service{
		Options: FantasyMarketOptions{
			TicksPerSecond:  0.1,
			StartDate:       time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			GameTimePerTick: time.Hour,
		},
		StockList: loadedStocks,
		EventList: loadedEvents,
		DB:        db,
	}

	go startLoop(s)
	fmt.Println("stated game loop :o")

	return s, nil
}

// startLoop startsrunningticks indefinitly
func startLoop(s *Service) {

	// We need to calculatre the current game date
	startDate := s.Options.StartDate
	gameTimePerTick := s.Options.GameTimePerTick
	s.TicksSinceStart, _ = s.DB.GetNextTick()
	dateNow := startDate.Add(gameTimePerTick * time.Duration(s.TicksSinceStart))

	for {
		s.tick(dateNow)

		// Sleep for the duration of a single tick (Since we want 1 tick in 10 Seconds)
		time.Sleep(time.Duration(1/s.Options.TicksPerSecond) * time.Second)

		// Adding 1 hour every tick(Update) (10 seconds when TicksPerSecond=0.1 ) onto the previously defined Date time
		dateNow = dateNow.Add(gameTimePerTick)
		s.TicksSinceStart++
	}
}

// tick is updating the current state of our system
func (s *Service) tick(dateNow time.Time) {
	// TODO: Get currently Running Events from the database (models.Event)

	fmt.Println("[running tick:  " + strconv.FormatInt(s.TicksSinceStart, 10) + "]")

	currentlyRunningEvents, _ := s.DB.GetEvents()                      // Sub this for the DB query results
	lastStockIndexes, _ := s.DB.GetStocksAtTick(s.TicksSinceStart - 1) // Sub this for the DB query results

	s.checkEventStillGoing(currentlyRunningEvents, dateNow)

	// Events: Events have tags: Fixed, Recurring, Random
	//    Hardcoded Events => Elections, Olympic Games etc
	//    Definate Date Events (Moon Landing 1969)?

	s.ComputeStockNumbers(lastStockIndexes, currentlyRunningEvents)
	// saveNextStockIndexes()
}

// checkEventStillGoing calculates if the duration of the event is over and then removes the event
func (s Service) checkEventStillGoing(e []models.Event, dateNow time.Time) {
	for i := 0; i < len(e); i++ {
		endDate := e[i].CreatedAt.Add(e[i].Duration) // Calculate the endDate by adding the Duration to the time created
		if !dateNow.Before(endDate) {                // Check if the current date is after the end date.
			s.DB.RemoveEvent(e[i].EventID)
		}
	}
}

func (s Service) getNextEventFromProbability(e models.Event) (string, error) {
	event := s.EventList[e.EventID]
	r := hash.Int64HashRange(0, 10, event.EventID)
	randomFloat := float64(r / 10) // Get the float for computation

	// [EventEffect, EventEffect, EventEffect]
	// EventEffect.Chance
	lowerBound := float64(0)
	for _, effect := range event.Effects { // 0.4 :: 0.6 :: 1 so r is between 0 - 0.4 (effect.Chance inclusive) and 0.4 - 0.6 and 0.6 - 1
		if lowerBound < randomFloat && randomFloat <= effect.Chance || randomFloat == 0 {
			return effect.NextEventID, nil
			// TODO: Ask alex again how the details work
		}
		lowerBound += effect.Chance
	}
	return "", errors.New("Empty Event Effect Error")
}

func (s Service) getEventAffectedness(activeEvents []models.Event, stock models.Stock) int64 {

	affectedness := int64(0)
	for _, activeEvent := range activeEvents {

		fullEventDetails := s.EventList[activeEvent.EventID]
		fullStockDetails := s.StockList[stock.Symbol]
		// models.Stock is what we get from the database and is a "lite" version of the "full" stock struct
		// Hence we take the stock symbol as the key to extract the stock with the full details from the
		// stock list and call it fullStockDetails

		for _, tagOptions := range fullEventDetails.Tags {

			_, affectedByTag := fullStockDetails.Tags[tagOptions.AffectsTag]
			affectedBySymbol := stock.Symbol == tagOptions.AffectsStock

			if affectedByTag || affectedBySymbol {
				affectedness += tagOptions.Trend
			}

		}

	}

	return affectedness
}

// ComputeStockNumbers computes the index at the next tick for a list of stocks
func (s Service) ComputeStockNumbers(stocks []models.Stock, e []models.Event) {

	// This computes the random and own stock, not taking into account other peoples selling
	// As a stock drops to a % of its value, theres gonna be more buyers or more sellers
	for _, stock := range stocks {
		stock.Index += s.GetTendency(stock, s.getEventAffectedness(e, stock)) // Range of -2 to 2
		fmt.Println("Name: ", stock.Symbol, "Index: ", stock.Index)
		s.DB.AddStock(stock, s.TicksSinceStart)
	}
	fmt.Println("-----------------------------")
}

// GetTendency calculates the tendency of a stock to go up or down
func (s Service) GetTendency(stock models.Stock, EventAffectedness int64) int64 {
	const rangeValue int64 = 10
	const weighConst = 2000
	stockSettings := s.StockList[stock.Symbol]

	seed := stock.Symbol + strconv.FormatInt(s.TicksSinceStart, 10)
	// weightOfTrends is the value that determines how prominent the trends are in the calculation.
	// the higher the weighConst, the less prominent the trend is
	weightOfTrends := (stock.Index / weighConst)

	// randomModifier needs a random value in a range multiplied by a variable (depending on the stock) to create spikes in the stock chart
	randomModifier := hash.Int64HashRange(-rangeValue, rangeValue, seed) * stockSettings.Stability
	// stockTrend and eventTrend are the inputs to calculate the stock graph trends
	stockTrend := weightOfTrends * stockSettings.Trend
	eventTrend := weightOfTrends * EventAffectedness
	tendency := randomModifier + stockTrend + eventTrend

	return tendency
	// Stability indicates how strong the random aspect is evaluated in comparison to the trend
}

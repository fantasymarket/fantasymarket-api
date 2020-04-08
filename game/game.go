package game

import (
	"fantasymarket/database"
	"fantasymarket/database/models"
	"fantasymarket/game/events"
	"fantasymarket/game/stocks"
	"fantasymarket/utils"
	"fantasymarket/utils/hash"
	"fmt"
	"strconv"
	"time"
)

// Service is the GameService
type Service struct {
	EventDetails    map[string]events.EventDetails
	StockDetails    map[string]stocks.StockDetails
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

// GetCurrentDate returns the current in-game date
func (s *Service) GetCurrentDate() time.Time {
	timeSinceStart := time.Duration(s.TicksSinceStart) * s.Options.GameTimePerTick
	return s.Options.StartDate.Add(timeSinceStart)
}

// Start starts the game loop
func Start(db *database.Service) (*Service, error) {

	loadedStocks, err := stocks.LoadStockDetails()
	if err != nil {
		return nil, err
	}

	if err := db.CreateInitialStocks(loadedStocks); err != nil {
		return nil, err
	}

	// TODO: right now, this map is empty
	loadedEvents, err := events.LoadEventDetails()
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
		StockDetails: loadedStocks,
		EventDetails: loadedEvents,
		DB:           db,
	}

	go startLoop(s)
	fmt.Println("stated game loop")

	return s, nil
}

// startLoop startsrunningticks indefinitly
func startLoop(s *Service) {
	s.TicksSinceStart, _ = s.DB.GetNextTick()
	fmt.Println("loaded ticksSinceStart from database:", s.TicksSinceStart, "ticks")

	for {
		s.tick()

		// Sleep for the duration of a single tick (Since we want 1 tick in 10 Seconds)
		time.Sleep(time.Duration(1/s.Options.TicksPerSecond) * time.Second)

		s.TicksSinceStart++
	}
}

// GetRandomEventEffect selects a random effect from an event
func (s *Service) GetRandomEventEffect(e models.Event) (string, error) {
	event := s.EventDetails[e.EventID]

	effects := make(map[string]float64)
	for _, e := range event.Effects {
		effects[e.EventID] = e.Chance
	}

	seed := e.EventID + strconv.FormatInt(s.TicksSinceStart, 10)
	return utils.SelectRandomWeightedItem(effects, seed)
}

// tick is updating the current state of our system
func (s *Service) tick() error {
	fmt.Println("\n> tick: " + strconv.FormatInt(s.TicksSinceStart, 10))

	currentlyRunningEvents, _ := s.DB.GetEvents(s.GetCurrentDate())    // Sub this for the DB query results
	lastStockIndexes, _ := s.DB.GetStocksAtTick(s.TicksSinceStart - 1) // Sub this for the DB query results

	s.checkEventStillActive(currentlyRunningEvents)

	// TODO: add new events to database:
	//    - fixed events that need to be added at a fixed date
	//		- random events
	// 		- reccuring events

	newStocks := s.ComputeStockNumbers(lastStockIndexes, currentlyRunningEvents)
	if err := s.DB.AddStocks(newStocks, s.TicksSinceStart); err != nil {
		return err
	}

	// TODO: process current orderbook

	return nil
}

// checkEventStillActive calculates if the duration of the event is over and then removes the event
func (s Service) checkEventStillActive(events []models.Event) {

	currentDate := s.GetCurrentDate()

	for _, event := range events {

		eventDetails := s.EventDetails[event.EventID]
		endDate := eventDetails.Duration.Shift(event.CreatedAt)

		if !currentDate.Before(endDate) {
			s.DB.RemoveEvent(event.ID)
		}
	}
}

// CalcutaleAffectedness calculates how much a stock is affected by all currently running events
func (s Service) CalcutaleAffectedness(stocks []models.Stock, activeEvents []models.Event) map[string]float64 {
	var affectedness map[string]float64

	for _, activeEvent := range activeEvents {
		for _, tagOptions := range s.EventDetails[activeEvent.EventID].Tags {

			for _, stock := range stocks {
				stockDetails := s.StockDetails[stock.Symbol]

				affectedByTag := utils.Some(stockDetails.Tags, tagOptions.AffectsTags)
				affectedBySymbol := utils.Includes(tagOptions.AffectsStocks, stock.Symbol)

				if affectedByTag || affectedBySymbol {

					if tagOptions.MinTrend != tagOptions.MaxTrend {
						seed := stock.Symbol + strconv.FormatInt(s.TicksSinceStart, 10)
						trend := hash.Int64HashRange(int64(tagOptions.MinTrend*1000), int64(tagOptions.MaxTrend*1000), seed)
						affectedness[stock.Symbol] += float64(trend) / 1000
						continue
					}

					affectedness[stock.Symbol] += tagOptions.Trend
				}
			}

		}
	}

	return affectedness
}

// ComputeStockNumbers computes the index at the next tick for a list of stocks
func (s Service) ComputeStockNumbers(stocks []models.Stock, events []models.Event) []models.Stock {

	affectedness := s.CalcutaleAffectedness(stocks, events)

	// This computes the random and own stock, not taking into account other peoples selling
	// As a stock drops to a % of its value, theres gonna be more buyers or more sellers
	for _, stock := range stocks {
		affectedness := affectedness[stock.Symbol]
		stock.Index += s.GetTendency(stock, affectedness)

		fmt.Println("Name: ", stock.Symbol, "Index: ", stock.Index)
	}

	return stocks
}

// GetTendency calculates the tendency of a stock to go up or down
func (s Service) GetTendency(stock models.Stock, eventAffectedness float64) int64 {
	const rangeValue int64 = 10
	const weighConst = 2000
	stockSettings := s.StockDetails[stock.Symbol]

	seed := stock.Symbol + strconv.FormatInt(s.TicksSinceStart, 10)
	// weightOfTrends is the value that determines how prominent the trends are in the calculation.
	// the higher the weighConst, the less prominent the trend is
	weightOfTrends := (float64(stock.Index) / weighConst)

	// randomModifier needs a random value in a range multiplied by a variable (depending on the stock) to create spikes in the stock chart
	randomModifier := hash.Int64HashRange(-rangeValue, rangeValue, seed)
	// stockTrend and eventTrend are the inputs to calculate the stock graph trends
	stockTrend := weightOfTrends * stockSettings.Trend
	eventTrend := weightOfTrends * eventAffectedness
	tendency := float64(randomModifier)*stockSettings.Stability + stockTrend + eventTrend

	return int64(tendency)
}

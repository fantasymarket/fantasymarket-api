package game

import (
	"fantasymarket/database"
	"fantasymarket/database/models"
	"fantasymarket/game/events"
	"fantasymarket/game/stocks"
	"fantasymarket/utils"
	"fantasymarket/utils/config"
	"fantasymarket/utils/hash"
	"fantasymarket/utils/timeutils"

	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// Service is the GameService
type Service struct {
	EventDetails    map[string]events.EventDetails
	StockDetails    map[string]stocks.StockDetails
	DB              *database.Service
	Config          *config.Config
	TicksSinceStart int64
}

// GetCurrentDate returns the current in-game date
func (s *Service) GetCurrentDate() time.Time {
	timeSinceStart := time.Duration(s.TicksSinceStart) * s.Config.Game.GameTimePerTick
	return s.Config.Game.StartDate.Add(timeSinceStart)
}

// Start starts the game loop
func Start(db *database.Service, config *config.Config) (*Service, error) {

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
		Config:       config,
		StockDetails: loadedStocks,
		EventDetails: loadedEvents,
		DB:           db,
	}

	go startLoop(s)
	log.Info().Msg("successfully started the game loop")

	return s, nil
}

// startLoop startsrunningticks indefinitly
func startLoop(s *Service) {
	s.TicksSinceStart, _ = s.DB.GetNextTick()
	log.Debug().Int64("ticksSinceStart", s.TicksSinceStart).Msg("loaded loaded ticksSinceStart from database")

	for {
		s.tick()

		timePerTick := time.Duration(1/s.Config.Game.TicksPerSecond) * time.Second
		time.Sleep(timePerTick)

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
	log.Debug().Int64("tick", s.TicksSinceStart).Msg("running tick")

	currentlyRunningEvents, _ := s.DB.GetEvents()                      // Sub this for the DB query results
	lastStockIndexes, _ := s.DB.GetStocksAtTick(s.TicksSinceStart - 1) // Sub this for the DB query results

	// TODO: add new events to database:
	//      - fixed events that need to be added at a fixed date
	//		- random events
	// 		- reccuring events

	s.removeInactiveEvents(currentlyRunningEvents)
	newStocks := s.ComputeStockNumbers(lastStockIndexes, currentlyRunningEvents)
	if err := s.DB.AddStocks(newStocks, s.TicksSinceStart); err != nil {
		return err
	}

	// TODO: process current orderbook

	return nil
}

// removeInactiveEvents calculates if the duration of the event is over and then removes the event
func (s Service) removeInactiveEvents(events []models.Event) {

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

	activeTags := s.GetActiveEventTags(activeEvents)

	for _, tag := range activeTags {
		for _, stock := range stocks {
			stockDetails := s.StockDetails[stock.Symbol]

			affectedByTag := utils.Some(stockDetails.Tags, tag.AffectsTags)
			affectedBySymbol := utils.Includes(tag.AffectsStocks, stock.Symbol)

			if affectedByTag || affectedBySymbol {
				affectedness[stock.Symbol] += tag.CalculateTrend(s.TicksSinceStart, stock.Symbol)
			}
		}
	}

	return affectedness
}

// GetActiveEventTags returns a list of event tags that
// should currently be affecting all stocks
func (s Service) GetActiveEventTags(activeEvents []models.Event) []events.TagOptions {
	var activeTags []events.TagOptions
	for _, activeEvent := range activeEvents {
		eventDetails := s.EventDetails[activeEvent.EventID]
		for _, tag := range eventDetails.Tags {

			startDate := activeEvent.CreatedAt
			currentDate := s.GetCurrentDate()

			if !timeutils.Duration.IsZero(tag.Offset) {
				startDate = tag.Offset.Shift(startDate)
			}

			if startDate.After(currentDate) {
				continue
			}

			if !timeutils.Duration.IsZero(tag.Duration) && tag.Duration.Shift(startDate).Before(currentDate) {
				continue
			}

			activeTags = append(activeTags, tag)
		}
	}

	return activeTags
}

// ComputeStockNumbers computes the index at the next tick for a list of stocks
func (s Service) ComputeStockNumbers(stocks []models.Stock, events []models.Event) []models.Stock {

	affectedness := s.CalcutaleAffectedness(stocks, events)

	// This computes the random and own stock, not taking into account other peoples selling
	// As a stock drops to a % of its value, theres gonna be more buyers or more sellers
	for _, stock := range stocks {
		affectedness := affectedness[stock.Symbol]
		stock.Index += s.GetTendency(stock, affectedness)

		log.Debug().Str("name", stock.Symbol).Int64("index", stock.Index).Msg("updated stock")
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

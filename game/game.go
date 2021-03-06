package game

import (
	"fantasymarket/database"
	"fantasymarket/database/models"
	"fantasymarket/game/details"
	"fantasymarket/utils"
	"fantasymarket/utils/config"
	"fantasymarket/utils/hash"
	"fantasymarket/utils/timeutils"
	"fmt"

	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// Service is the GameService
type Service struct {
	EventDetails    map[string]details.EventDetails
	StockDetails    map[string]details.StockDetails
	DB              *database.Service
	Config          *config.Config
	TicksSinceStart int64

	// a history of all events that have run in the past
	// map[eventID][]createdAt
	EventHistory map[string][]time.Time
}

// GetCurrentDate returns the current in-game date
func (s *Service) GetCurrentDate() time.Time {
	timeSinceStart := time.Duration(s.TicksSinceStart) * s.Config.Game.GameTimePerTick
	return s.Config.Game.StartDate.Add(timeSinceStart)
}

// Start starts the game loop
func Start(db *database.Service, config *config.Config) (*Service, error) {

	loadedStocks, err := details.LoadStockDetails()
	if err != nil {
		return nil, fmt.Errorf("game: failed to Load stock Details: %w", err)
	}

	if err := db.CreateInitialStocks(loadedStocks); err != nil {
		return nil, fmt.Errorf("game: failed to initialize stocks: %w", err)
	}

	loadedEvents, err := details.LoadEventDetails()
	if err != nil {
		return nil, fmt.Errorf("game: failed to load events: %w", err)
	}

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
	s.EventHistory, _ = s.DB.GetEventHistory()

	log.Debug().Int64("ticksSinceStart", s.TicksSinceStart).Msg("loaded loaded ticksSinceStart from database")

	for {
		if err := s.tick(); err != nil {
			log.Error().Err(err).Int64("ticksSinceStart", s.TicksSinceStart).Msg("error while running tick")
		}

		timePerTick := time.Duration(1/s.Config.Game.TicksPerSecond) * time.Second
		time.Sleep(timePerTick)

		s.TicksSinceStart++
	}
}

// GetRandomEventEffect selects a random effect from an event
func (s *Service) GetRandomEventEffect(eventID string) (string, error) {
	event := s.EventDetails[eventID]

	effects := make(map[string]float64)
	for _, e := range event.Effects {
		effects[e.EventID] = e.Chance
	}

	seed := eventID + strconv.FormatInt(s.TicksSinceStart, 10)
	randomNumber := hash.Float64Hash(seed)
	return utils.SelectRandomWeightedItem(effects, randomNumber)
}

// tick is updating the current state of our system
func (s *Service) tick() error {
	log.Debug().Int64("tick", s.TicksSinceStart).Str("date", s.GetCurrentDate().String()).Msg("running tick")

	currentlyRunningEvents, err := s.DB.GetEvents(s.GetCurrentDate())
	if err != nil {
		return fmt.Errorf("game: failed to get events from DB: %w", err)
	}

	lastStockIndexes, err := s.DB.GetStocksAtTick(s.TicksSinceStart - 1)
	if err != nil {
		return fmt.Errorf("game: failed to get stocks indexes: %w", err)
	}

	if err := s.startEvents(); err != nil {
		return fmt.Errorf("game: failed to start events: %w", err)
	}

	s.removeInactiveEvents(currentlyRunningEvents)
	newStocks := s.ComputeStockNumbers(lastStockIndexes, currentlyRunningEvents)

	if err := s.DB.AddStocks(newStocks, s.TicksSinceStart); err != nil {
		return fmt.Errorf("game: failed to add stocks: %w", err)
	}

	// TODO: process current orderbook
	// s.processOrders()

	return nil
}

// removeInactiveEvents calculates if the duration of the event is over and then removes the event
func (s Service) removeInactiveEvents(events []models.Event) {

	currentDate := s.GetCurrentDate()

	for _, event := range events {

		eventDetails := s.EventDetails[event.EventID]
		endDate := eventDetails.Duration.Shift(event.CreatedAt)

		if !currentDate.Before(endDate) {
			s.DB.RemoveEvent(event.EventID)
		}
	}
}

// CalculateAffectedness calculates how much a stock is affected by all currently running events
func (s Service) CalculateAffectedness(stocks []models.Stock, activeEvents []models.Event) map[string]float64 {
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
func (s Service) GetActiveEventTags(activeEvents []models.Event) []details.TagOptions {
	var activeTags []details.TagOptions

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

			activeTags = append(activeTags, details.TagOptions{
				AffectsTags:   tag.AffectsTags,
				AffectsStocks: tag.AffectsStocks,
				Trend:         tag.Trend,
			})
		}
	}

	return activeTags
}

// ComputeStockNumbers computes the index at the next tick for a list of stocks
func (s Service) ComputeStockNumbers(stocks []models.Stock, events []models.Event) []models.Stock {

	affectedness := s.CalculateAffectedness(stocks, events)

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

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
	AssetDetails    map[string]details.AssetDetails
	DB              *database.Service
	Config          *config.Config
	TicksSinceStart int64

	// a history of all events that have run in the past
	// map[eventID][]createdAt
	EventHistory map[string][]time.Time
}

// GetCurrentDate returns the current in-game date
func (s *Service) GetCurrentDate() time.Time {
	return s.TickToTime(s.TicksSinceStart)
}

// TickToTime converts a tick to a timestamp
func (s *Service) TickToTime(ticks int64) time.Time {
	timeSinceStart := time.Duration(ticks) * s.Config.Game.GameTimePerTick
	return s.Config.Game.StartDate.Add(timeSinceStart)
}

// Start starts the game loop
func Start(db *database.Service, config *config.Config) (*Service, error) {

	loadedAssets, err := details.LoadAssetDetails()
	if err != nil {
		return nil, fmt.Errorf("game: failed to Load asset Details: %w", err)
	}

	if err := db.CreateInitialAssets(loadedAssets); err != nil {
		return nil, fmt.Errorf("game: failed to initialize assets: %w", err)
	}

	loadedEvents, err := details.LoadEventDetails()
	if err != nil {
		return nil, fmt.Errorf("game: failed to load events: %w", err)
	}

	s := &Service{
		Config:       config,
		AssetDetails: loadedAssets,
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

	lastAssetIndexes, err := s.DB.GetAssetsAtTick(s.TicksSinceStart - 1)
	if err != nil {
		return fmt.Errorf("game: failed to get assets indexes: %w", err)
	}

	if err := s.startEvents(); err != nil {
		return fmt.Errorf("game: failed to start events: %w", err)
	}

	s.removeInactiveEvents(currentlyRunningEvents)
	newAssets := s.ComputeAssetNumbers(lastAssetIndexes, currentlyRunningEvents)

	if err := s.DB.AddAssets(newAssets, s.TicksSinceStart); err != nil {
		return fmt.Errorf("game: failed to add assets: %w", err)
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

// CalculateAffectedness calculates how much a asset is affected by all currently running events
func (s Service) CalculateAffectedness(assets []models.Asset, activeEvents []models.Event) map[string]float64 {
	var affectedness map[string]float64

	activeTags := s.GetActiveEventTags(activeEvents)

	for _, tag := range activeTags {
		for _, asset := range assets {
			assetDetails := s.AssetDetails[asset.Symbol]

			affectedByTag := utils.Some(assetDetails.Tags, tag.AffectsTags)
			affectedBySymbol := utils.Includes(tag.AffectsAssets, asset.Symbol)

			if affectedByTag || affectedBySymbol {
				affectedness[asset.Symbol] += tag.CalculateTrend(s.TicksSinceStart, asset.Symbol)
			}
		}
	}

	return affectedness
}

// GetActiveEventTags returns a list of event tags that
// should currently be affecting all assets
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
				AffectsAssets: tag.AffectsAssets,
				Trend:         tag.Trend,
			})
		}
	}

	return activeTags
}

// ComputeAssetNumbers computes the index at the next tick for a list of assets
func (s Service) ComputeAssetNumbers(assets []models.Asset, events []models.Event) []models.Asset {

	affectedness := s.CalculateAffectedness(assets, events)

	// This computes the random and own asset, not taking into account other peoples selling
	// As a asset drops to a % of its value, theres gonna be more buyers or more sellers
	for _, asset := range assets {
		affectedness := affectedness[asset.Symbol]
		asset.Index += s.GetTendency(asset, affectedness)

		log.Debug().Str("name", asset.Symbol).Int64("index", asset.Index).Msg("updated asset")
	}

	return assets
}

// GetTendency calculates the tendency of a asset to go up or down
func (s Service) GetTendency(asset models.Asset, eventAffectedness float64) int64 {
	const rangeValue int64 = 10
	const weighConst = 2000
	assetSettings := s.AssetDetails[asset.Symbol]

	seed := asset.Symbol + strconv.FormatInt(s.TicksSinceStart, 10)
	// weightOfTrends is the value that determines how prominent the trends are in the calculation.
	// the higher the weighConst, the less prominent the trend is
	weightOfTrends := (float64(asset.Index) / weighConst)

	// randomModifier needs a random value in a range multiplied by a variable (depending on the asset) to create spikes in the asset chart
	randomModifier := hash.Int64HashRange(-rangeValue, rangeValue, seed)
	// assetTrend and eventTrend are the inputs to calculate the asset graph trends
	assetTrend := weightOfTrends * assetSettings.Trend
	eventTrend := weightOfTrends * eventAffectedness
	tendency := float64(randomModifier)*assetSettings.Stability + assetTrend + eventTrend

	return int64(tendency)
}

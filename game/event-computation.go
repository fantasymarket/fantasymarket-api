package game

import (
	"bytes"
	"fantasymarket/game/details"
	"fantasymarket/utils/hash"
	"fantasymarket/utils/timeutils"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// StartEvents checks an event if it should run or not depending on the event type
func (s *Service) startEvents() error {
	currentDate := s.GetCurrentDate()
	events := s.EventDetails

	for _, event := range events {

		createdAt := event.FixedDate.Time
		eventID := event.EventID
		seed := eventID + strconv.FormatInt(s.TicksSinceStart, 10)

		if event.Type == "recurring" {
			date := event.FixedDate.Time
			for date.Before(currentDate) {
				date = event.RecurringDuration.Shift(date)
			}

			date = timeutils.ShiftBack(event.RecurringDuration, date)
			createdAt = date
			event.FixedDate = timeutils.Time{Time: createdAt}
		}

		if event.Type == "random" {
			ticksPerDay := time.Hour * 24 / s.Config.Game.GameTimePerTick
			chancePerTick := event.RandomChancePerDay * float64(ticksPerDay)
			createdAt = s.GetCurrentDate()

			// we can just skip to the next loop
			if chancePerTick < hash.Float64Hash(seed) {
				continue
			}
		}

		if !event.FixedDateRandomOffset.IsZero() {
			offset := calculateRandomOffset(event.FixedDateRandomOffset, seed)
			createdAt.Add(offset)
		}

		if s.eventNeedsToBeRun(event) {
			if err := s.addEventToRun(event, createdAt); err != nil {
				return fmt.Errorf("event-computation: failed to start event: %w", err)
			}
		}
	}

	return nil
}

func (s *Service) addEventToRun(event details.EventDetails, createdAt time.Time) error {
	eventID := event.EventID

	if err := s.DB.AddEvent(event, createdAt); err != nil {
		return err
	}

	if _, ok := s.EventHistory[eventID]; !ok {
		s.EventHistory[eventID] = []time.Time{}
	}
	s.EventHistory[eventID] = append(s.EventHistory[eventID], createdAt)

	log.Debug().Str("eventID", event.EventID).Msg("starting event")
	return nil
}

func calculateRandomOffset(randomOffset timeutils.Duration, seed string) time.Duration {
	date := time.Time{}
	shiftedDate := randomOffset.Shift(date)
	difference := shiftedDate.Sub(shiftedDate)

	offset := hash.Int64HashRange(0, int64(difference), seed)

	return time.Duration(offset)
}

func (s *Service) eventNeedsToBeRun(event details.EventDetails) bool {
	currentDate := s.GetCurrentDate()

	lengthOfEventHistorySlice := len(s.EventHistory[event.EventID])

	timeStampOfLastEvent := time.Time{}
	if lengthOfEventHistorySlice != 0 {
		timeStampOfLastEvent = s.EventHistory[event.EventID][lengthOfEventHistorySlice-1]
	}

	eventHistory, ok := s.EventHistory[event.EventID]

	eventHasNeverRun := !ok || len(eventHistory) == 0
	eventDateInPast := currentDate.After(event.FixedDate.Time)

	randomEventShouldRun := currentDate.After(timeStampOfLastEvent.Add(event.MinTimeBetweenEvents))

	return eventHasNeverRun && eventDateInPast || event.Type == "random" && randomEventShouldRun
}

// ChangeDescriptionPlaceholder fills the templates of a description string
func (s *Service) ChangeDescriptionPlaceholder(description string) (string, error) {
	currentDate := s.GetCurrentDate()

	data := struct {
		Year  int
		Month string
		Day   int
	}{
		Year:  currentDate.Year(),
		Month: currentDate.Month().String(),
		Day:   currentDate.Day(),
	}

	tmpl, err := template.New("description").Parse(description)
	if err != nil {
		return "", fmt.Errorf("event-computation: failed to change description: %w", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("event-computation: failed to execute description substitution: %w", err)
	}

	return result.String(), nil
}

package game

import (
	"bytes"
	"fantasymarket/game/events"
	"fantasymarket/utils/hash"
	"fantasymarket/utils/timeutils"
	"html/template"
	"strconv"
	"time"
)

func (s *Service) startEvents() {
	currentDate := s.GetCurrentDate()
	events := s.EventDetails

	for _, event := range events {

		eventNeedsToBeRun := false
		createdAt := event.FixedDate.Time
		eventID := event.EventID
		seed := eventID + strconv.FormatInt(s.TicksSinceStart, 10)

		switch event.Type {
		case "fixed":
			eventNeedsToBeRun = s.eventNeedsToBeRun(event)
		case "recurring":

			date := event.FixedDate.Time
			for date.Before(currentDate) {
				date = event.RecurringDuration.Shift(date)
			}

			date = timeutils.ShiftBack(event.RecurringDuration, date)
			createdAt = date
			event.FixedDate = timeutils.Time{Time: createdAt}
			eventNeedsToBeRun = s.eventNeedsToBeRun(event)
		case "random":

			ticksPerDay := time.Hour * 24 / s.Config.Game.GameTimePerTick
			chancePerTick := event.RandomChancePerDay * float64(ticksPerDay)

			if chancePerTick > (float64(hash.Int64HashRange(0, 1e6, seed)) / 1e6) {
				eventNeedsToBeRun = s.eventNeedsToBeRun(event)
			}

			createdAt = s.GetCurrentDate()
		}

		if !event.FixedDateRandomOffset.IsZero() {
			offset := calculateRandomOffset(event.FixedDateRandomOffset, seed)
			createdAt.Add(offset)
		}

		if eventNeedsToBeRun {
			s.DB.AddEvent(event, createdAt)

			if _, ok := s.EventHistory[eventID]; !ok {
				s.EventHistory[eventID] = []time.Time{}
			}
			s.EventHistory[eventID] = append(s.EventHistory[eventID], createdAt)
		}
	}
}

func calculateRandomOffset(randomOffset timeutils.Duration, seed string) time.Duration {
	date := time.Time{}
	shiftedDate := randomOffset.Shift(date)
	difference := shiftedDate.Sub(shiftedDate)

	offset := hash.Int64HashRange(0, int64(difference), seed)

	return time.Duration(offset)
}

func (s *Service) eventNeedsToBeRun(event events.EventDetails) bool {
	currentDate := s.GetCurrentDate()

	lengthOfEventHistorySlice := len(s.EventHistory[event.EventID])
	timeStampOfLastEvent := s.EventHistory[event.EventID][lengthOfEventHistorySlice-1]

	eventHistory, ok := s.EventHistory[event.EventID]

	eventHasNeverRun := !ok || len(eventHistory) == 0
	eventDateInPast := currentDate.After(event.FixedDate.Time)

	// TODO: handle events that can happening multiple times
	// check if MinTimeBetween events is long enough
	// eventShouldRun :=  currentDate after MinTimeBetween events + eventHistory[len(eventHistory)-1]
	// eventShouldRun := checks if the current date is after the combined time of last time the event ran and how long
	//	it has to wait until the event can run again.

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
		return "", err
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}

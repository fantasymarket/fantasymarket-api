package game

import (
	"fantasymarket/database/models"
	"fantasymarket/game/events"
	"fantasymarket/utils/timeutils"
	"fantasymarket/utils/hash"
	"time"
)

func (s *Service) startEvents(currentlyRunningEvents []models.Event) {
	currentDate := s.GetCurrentDate()
	events := s.EventDetails

	var currentlyRunningEventsMap = make(map[string]models.Event)
	for _, e := range currentlyRunningEvents {
		currentlyRunningEventsMap[e.EventID] = e
	}

	for _, event := range events {

		eventNeedsToBeRun := false
		createdAt := event.FixedDate.Time

		switch event.Type {
		case "fixed":
			eventNeedsToBeRun = s.eventNeedsToBeRun(event, currentlyRunningEventsMap)
		case "recurring":

			// TODO remove this comment never
			// event is 2020.1.1 00:00
			// today is 2028.1.1 01:00 // so 1 tick in
			// so we keep adding 4 years until date > today
			// date = 2032.1.1 00:00
			// then - 4 years
			// 2028.1.1 00:00
			// this date is before the current time so the event will run in the next tick

			// and then run the same logic as fixed
			date := event.FixedDate.Time
			for date.Before(currentDate) {
				date = event.RecurringDuration.Shift(date)
			}
			date = timeutils.ShiftBack(event.RecurringDuration, date)
			createdAt = date
			event.FixedDate = timeutils.Time{Time: createdAt}
			eventNeedsToBeRun = s.eventNeedsToBeRun(event, currentlyRunningEventsMap)
		case "random":
			eventNeedsToBeRun = false // TODO
			createdAt = s.GetCurrentDate()
		}

		if !eventNeedsToBeRun {
			continue
		}

		// TODO Add Random Offset to createdAt
		s.DB.AddEvent(event, createdAt)
	}
}

func (s *Service) addRandomOffset(time time.Time, randomOffset timeutils.Duration, seed string) time.Duration {
	date := time.Time{}
	shiftedDate :== randomOffset.Shift(date)
	difference := shiftedDate.Sub(shiftedDate)

	offset := hash.Int64HashRange(0, int64(difference), seed)

	return time.Duration(offset)
}

func (s *Service) eventNeedsToBeRun(event events.EventDetails, currentlyRunningEventsMap map[string]models.Event) bool {
	date := event.FixedDate.Time
	currentDate := s.GetCurrentDate()
	_, currentlyRunning := currentlyRunningEventsMap[event.EventID]

	// notRunInThePast means the event can't have been run in another tick before this
	notRunInThePast := event.FixedDate.After(timeutils.ShiftBack(event.Duration, date))

	// Event hasn't happened yet since it is not in currentlyRunning events
	// and since fixedDate is after (currentTime - event.Duration),
	// it can't have been run in the past
	eventHasNotRunYet := !currentlyRunning && notRunInThePast

	// fixed Date is in the past so we'll have to run the event
	fixedDateInPast := event.FixedDate.Before(date)

	if fixedDateInPast && eventHasNotRunYet {
		return true
	}
	return false
	//TODO: Add an example to understand it better
}

func (s *Service) getRandomDateinRange(date timeutils.Time, duration timeutils.Duration) {

	start := date
	end := duration.Shift(date.Time)

	// createdAt needs to be set to the date in our db
}

func (s *Service) changeDescriptionPlaceholder() {
	// parse templates with
	// https://golang.org/pkg/text/template/

	date := s.GetCurrentDate()
	//TODO: Get the description of an event and find {year}, then sub {year} for data.year
}

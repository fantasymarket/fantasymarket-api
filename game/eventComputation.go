package game

import (
	"fantasymarket/database/models"
	"fantasymarket/utils/timeutils"
)

func (s *Service) startEvents(currentlyRunningEvents []models.Event) {
	events := s.EventDetails
	date := s.GetCurrentDate()

	var currentlyRunningEventsMap = make(map[string]models.Event)
	for _, e := range currentlyRunningEvents {
		currentlyRunningEventsMap[e.EventID] = e
	}

	for _, event := range events {
		switch event.Type {
		case "fixed":
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
				// RUN THE FING EVENT
				// Add Random Offset
			}
			//TODO: Add an example to understand it better
		case "recurring":

		case "random":

		default:
			continue
		}
	}
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

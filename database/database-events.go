package database

import (
	"errors"
	"fantasymarket/database/models"
	"fantasymarket/game/details"
	"time"
)

// GetEvents fetches all currently active events
func (s *Service) GetEvents(currentDate time.Time) ([]models.Event, error) {
	var events []models.Event

	if err := s.DB.Where(models.Event{
		Active: true,
	}).Where("created_at < ?", currentDate).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// AddEvent adds an event to the event table
func (s *Service) AddEvent(event details.EventDetails, createdAt time.Time) error {
	if event.EventID == "" {
		return errors.New("error: attempting to add empty event")
	}
	return s.DB.Create(&models.Event{
		EventID:   event.EventID,
		Title:     event.Title,
		Text:      event.Description,
		Active:    true,
		CreatedAt: createdAt,
	}).Error
}

// RemoveEvent marks an event as inactive so it won't affect stocks in the GameLoop anymore
func (s *Service) RemoveEvent(eventID string) error {
	return s.DB.Model(&models.Event{}).Where(models.Event{Active: true, EventID: eventID}).Update("active", false).Error
}

// GetEventHistory returns all the events that ran at some point as a map
func (s *Service) GetEventHistory() (map[string][]time.Time, error) {
	eventHistory := map[string][]time.Time{}

	var events []models.Event
	if err := s.DB.Find(&events).Error; err != nil {
		return nil, err
	}

	for _, event := range events {
		eventID := event.EventID
		createdAt := event.CreatedAt

		if _, exists := eventHistory[eventID]; !exists {
			eventHistory[eventID] = []time.Time{}
		}
		eventHistory[eventID] = append(eventHistory[eventID], createdAt)

	}

	return eventHistory, nil
}

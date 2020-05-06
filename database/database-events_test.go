package database_test

import (
	"fantasymarket/database/models"
	"fantasymarket/game/details"
	"time"

	"github.com/stretchr/testify/assert"
)

type GetEventTestData struct {
	eventID     string
	expectation models.Event
}

var testGetEventData = []GetEventTestData{
	{
		eventID: "TestEvent1",
		expectation: models.Event{
			EventID:   "TestEvent1",
			Title:     "TestEvent1TestEvent1",
			Active:    true,
			CreatedAt: parseTime("2019-12-30T15:01:05Z"),
		},
	},
	{
		eventID: "TestEvent2",
		expectation: models.Event{
			EventID:   "TestEvent2",
			Title:     "TestEvent2TestEvent2",
			Active:    true,
			CreatedAt: parseTime("2019-12-30T15:02:05Z"),
		},
	},
	{},
}

func (suite *DatabaseTestSuite) TestGetEvents() {
	createdAt := parseTime("2019-12-30T15:00:05Z")
	for i, event := range testGetEventData {
		if event.eventID != "" {
			createdAt = createdAt.Add(time.Minute)
			suite.dbService.DB.Create(&models.Event{
				EventID:   event.eventID,
				Title:     event.eventID + event.eventID,
				Text:      "",
				Active:    true,
				CreatedAt: createdAt,
			})
			newEvent, err := suite.dbService.GetEvents(parseTime("2020-12-30T15:00:05Z"))
			assert.Equal(suite.T(), nil, err)

			// I hate my life
			// Cos newEvent is an array, I need to get correct model for each test case with the index
			assert.EqualValues(suite.T(), event.expectation.EventID, newEvent[i].EventID)
			assert.EqualValues(suite.T(), event.expectation.Title, newEvent[i].Title)
			assert.EqualValues(suite.T(), event.expectation.Active, newEvent[i].Active)
			assert.EqualValues(suite.T(), event.expectation.CreatedAt, newEvent[i].CreatedAt)
		}
	}
	suite.dbService.DB.Close()
}

var testAddEventData = []details.EventDetails{
	{
		EventID: "testEvent1",
	},
	{},
}

func (suite *DatabaseTestSuite) TestAddEvent() {

	createdAt := parseTime("2019-12-30T15:04:05Z")
	currentDate := parseTime("2020-12-30T15:04:05Z")
	for _, event := range testAddEventData {
		suite.dbService.AddEvent(event, createdAt)
		err := suite.dbService.DB.Where(models.Event{
			Active: true,
		}).Where("created_at < ?", currentDate).Find(&models.Event{}).Error
		assert.Equal(suite.T(), nil, err)
	}

	suite.dbService.DB.Close()
}

var testRemoveEventData = []details.EventDetails{
	{
		EventID: "testEvent1",
	},
	{},
}

func (suite *DatabaseTestSuite) TestRemoveEvent() {

	createdAt := parseTime("2019-12-30T15:04:05Z")

	for _, event := range testRemoveEventData {
		suite.dbService.DB.Create(&models.Event{
			EventID:   event.EventID,
			Title:     event.Title,
			Text:      event.Description,
			Active:    true,
			CreatedAt: createdAt,
		})
		err := suite.dbService.RemoveEvent(event.EventID)
		assert.Equal(suite.T(), nil, err)
		err = suite.dbService.DB.Where("event_id = ?", event.EventID).Find(&models.Event{}).Error
		assert.Equal(suite.T(), nil, err)
		assert.Equal(suite.T(), true, suite.dbService.DB.Where("event_id = ? AND active = ?", event.EventID, true).Find(&models.Event{}).RecordNotFound())
	}

	suite.dbService.DB.Close()
}

type GetEventHistoryTestData struct {
	eventID     string
	expectation map[string][]time.Time
}

var testGetEventHistory = []GetEventHistoryTestData{
	{
		eventID:     "testEvent1",
		expectation: map[string][]time.Time{"testEvent1": {parseTime("2019-12-30T15:01:05Z")}},
	},
	{
		eventID: "testEvent2",
		expectation: map[string][]time.Time{
			"testEvent1": {parseTime("2019-12-30T15:01:05Z")},
			"testEvent2": {parseTime("2019-12-30T15:02:05Z")},
		},
	},
	{
		eventID: "testEvent1",
		expectation: map[string][]time.Time{
			"testEvent1": {parseTime("2019-12-30T15:01:05Z"), parseTime("2019-12-30T15:03:05Z")},
			"testEvent2": {parseTime("2019-12-30T15:02:05Z")},
		},
	},
	{
		expectation: map[string][]time.Time{
			"testEvent1": {parseTime("2019-12-30T15:01:05Z"), parseTime("2019-12-30T15:03:05Z")},
			"testEvent2": {parseTime("2019-12-30T15:02:05Z")},
		},
	},
}

func (suite *DatabaseTestSuite) TestGetEventHistory() {
	createdAt := parseTime("2019-12-30T15:00:05Z")

	for _, test := range testGetEventHistory {
		if test.eventID != "" {
			createdAt = createdAt.Add(time.Minute)
			suite.dbService.DB.Create(&models.Event{
				EventID:   test.eventID,
				Title:     "",
				Text:      "",
				Active:    true,
				CreatedAt: createdAt,
			})
			eventHistory, err := suite.dbService.GetEventHistory()
			assert.Equal(suite.T(), nil, err)
			assert.Equal(suite.T(), test.expectation, eventHistory)
		}
	}

	suite.dbService.DB.Close()
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

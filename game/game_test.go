package game

import (
	"fantasymarket/database/models"
	"fantasymarket/game/details"
	"fantasymarket/utils/config"
	"fantasymarket/utils/timeutils"
	"time"

	"github.com/senseyeio/duration"

	"testing"

	"github.com/stretchr/testify/assert"
)

type TestGetTendencyData struct {
	stock        models.Stock
	affectedness float64
	expectation  map[string]int64
}

var GetTendencyData = []TestGetTendencyData{
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: -1,
		expectation: map[string]int64{
			"TEST0": 1,
			"TEST1": 5,
			"TEST2": 21,
			"TEST3": -5,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: -100,
		expectation: map[string]int64{
			"TEST0": -494,
			"TEST1": -490,
			"TEST2": -474,
			"TEST3": -500,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: 100,
		expectation: map[string]int64{
			"TEST0": 506,
			"TEST1": 510,
			"TEST2": 526,
			"TEST3": 500,
		},
	},
	{
		stock:        models.Stock{Index: 10000, Symbol: "TEST"},
		affectedness: 0,
		expectation: map[string]int64{
			"TEST0": 6,
			"TEST1": 10,
			"TEST2": 26,
			"TEST3": 0,
		},
	},
}

var stockDetails = map[string]details.StockDetails{
	"TEST0": {Stability: 1, Trend: 1},
	"TEST1": {Stability: 5, Trend: 1},
	"TEST2": {Stability: 1, Trend: 5},
	"TEST3": {Stability: 0, Trend: 0},
}

func TestGetTendency(t *testing.T) {
	s := Service{
		StockDetails: stockDetails,
	}

	for _, test := range GetTendencyData {
		for i := range stockDetails {
			test.stock.Symbol = string(i)

			result := s.GetTendency(test.stock, test.affectedness)
			assert.Equal(t, test.expectation[i], result)
		}
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func parseDuration(s string) timeutils.Duration {
	d, _ := duration.ParseISO8601(s)
	return timeutils.Duration{Duration: d}
}

type activeEvent struct {
	EventID   string
	CreatedAt time.Time
	Tags      []details.TagOptions
}

type activeEventsTestData struct {
	activeEvents []activeEvent // map[eventID]createdAt
	result       []details.TagOptions
	expectation  []details.TagOptions
}

var testActiveTagsData = activeEventsTestData{
	activeEvents: []activeEvent{
		{
			EventID:   "event1",
			CreatedAt: parseTime("2019-12-30T15:04:05Z"),
			Tags: []details.TagOptions{
				{
					AffectsTags: []string{"some-type-only-event1-affects"},
					Offset:      parseDuration("P1D"),
					Duration:    parseDuration("P2M"),
				},
			},
		},
		{
			EventID:   "event2",
			CreatedAt: parseTime("2019-01-02T15:04:05Z"),
			Tags: []details.TagOptions{
				{
					AffectsStocks: []string{"IDEXX", "ANTM"},
					Offset:        parseDuration("P1D"),
					Duration:      parseDuration("P2M"),
				},
			},
		},
		{
			EventID:   "event3",
			CreatedAt: parseTime("2019-01-02T15:04:05Z"),
			Tags: []details.TagOptions{
				{
					AffectsStocks: []string{"GOOG", "APPL"},
					Offset:        parseDuration("P1D"),
					Duration:      parseDuration("P12M"),
				},
			},
		},
	},
	expectation: []details.TagOptions{
		{
			AffectsTags: []string{"some-type-only-event1-affects"},
		},
		{
			AffectsStocks: []string{"GOOG", "APPL"},
		},
	},
}

func TestGetActiveEventTags(t *testing.T) {

	events := []models.Event{}
	eventDetails := map[string]details.EventDetails{}

	for _, event := range testActiveTagsData.activeEvents {
		events = append(events, models.Event{
			Active:    true,
			CreatedAt: event.CreatedAt,
			EventID:   event.EventID,
		})

		eventDetails[event.EventID] = details.EventDetails{
			EventID: event.EventID,
			Tags:    event.Tags,
		}
	}

	startDate, _ := time.Parse(time.RFC3339, "2020-01-02T15:04:05Z")

	s := Service{
		Config: &config.Config{
			Game: config.GameConfig{
				StartDate: startDate,
			},
		},
		EventDetails: eventDetails,
	}

	result := s.GetActiveEventTags(events)
	assert.Equal(t, testActiveTagsData.expectation, result)
}

type getRandomEventEffectTestData struct {
	EventID      string
	eventDetails map[string]details.EventDetails
	expectation  string
}

var getRandomEventEffectData = []getRandomEventEffectTestData{
	{
		EventID:     "testEvent1",
		expectation: "newEvent1",
		eventDetails: map[string]details.EventDetails{
			"testEvent1": {
				Effects: []details.EventEffect{
					{
						Chance:  0.2,
						EventID: "newEvent1",
					},
					{
						Chance:  0.9,
						EventID: "newEvent2",
					},
					{
						Chance:  1,
						EventID: "newEvent3",
					},
				},
			},
		},
	},
	{
		EventID:     "testEvent2",
		expectation: "newEvent1",
		eventDetails: map[string]details.EventDetails{
			"testEvent2": {
				Effects: []details.EventEffect{
					{
						Chance:  0.001,
						EventID: "newEvent1",
					},
					{
						Chance:  0.002,
						EventID: "newEvent2",
					},
					{
						Chance:  1,
						EventID: "newEvent3",
					},
				},
			},
		},
	},
	{
		EventID:     "testEvent3",
		expectation: "",
		eventDetails: map[string]details.EventDetails{
			"testEvent2": {},
		},
	},
}

func TestGetRandomEventEffect(t *testing.T) {
	s := Service{
		TicksSinceStart: 100,
		EventDetails:    map[string]details.EventDetails{},
	}

	for _, test := range getRandomEventEffectData {
		s.EventDetails = test.eventDetails
		result, err := s.GetRandomEventEffect(test.EventID)
		if assert.NoError(t, err) {
			assert.Equal(t, test.expectation, result)
		}
	}
}

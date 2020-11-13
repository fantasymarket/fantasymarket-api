package game_test

import (
	"fantasymarket/game"
	"fantasymarket/game/details"
	"fantasymarket/utils/config"
	"fantasymarket/utils/timeutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type changeDescriptionPlaceholderTestData struct {
	input  string
	expect string
}

var changeDescriptionPlaceholderData = []changeDescriptionPlaceholderTestData{
	{
		input:  "In {{.Year}}, there will be many houses build",
		expect: "In 2006, there will be many houses build",
	},
	{
		input:  "There will be 100 more in {{.Year}}!",
		expect: "There will be 100 more in 2006!",
	},
}

func TestChangeDescriptionPlaceholder(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	s := game.Service{
		Config: &config.Config{
			Game: config.GameConfig{
				StartDate: startDate,
			},
		},
	}

	for _, test := range changeDescriptionPlaceholderData {
		result, err := s.ChangeDescriptionPlaceholder(test.input)
		if assert.NoError(t, err) {
			assert.Equal(t, test.expect, result)
		}
	}
}

type eventsNeedToBeRunTestData struct {
	event  details.EventDetails
	expect bool
}

var eventsNeedToBeRunData = []eventsNeedToBeRunTestData{
	{
		event: details.EventDetails{
			EventID:   "event1",
			Type:      "fixed",
			FixedDate: parseTimeAsTimeUtils("2006-06-01T14:06:05Z"),
		},
		expect: true,
	},
	{
		event: details.EventDetails{
			EventID: "event2",
			Type:    "random",
		},
		expect: true,
	},
	{
		event: details.EventDetails{
			EventID: "event3",
			Type:    "recurring",
		},
		expect: true,
	},
}

func TestEventNeedsToBeRun(t *testing.T) {
	startDate := parseTime("2006-06-01T15:04:05Z")

	s := game.Service{
		Config: &config.Config{
			Game: config.GameConfig{
				StartDate: startDate,
			},
		},
	}

	for i := 0; i < 3; i++ {
		s.EventHistory = map[string][]time.Time{
			"event1": {},
			"event2": {parseTime("2004-01-01T15:04:05Z"), parseTime("2005-05-31T15:04:05Z")},
			"event3": {parseTime("2004-06-01T15:04:05Z"), parseTime("2005-06-01T15:04:05Z")},
		}
	}

	for _, test := range eventsNeedToBeRunData {
		result := s.EventNeedsToBeRun(test.event)
		assert.Equal(t, test.expect, result)
	}
}

func parseTimeAsTimeUtils(s string) timeutils.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return timeutils.Time{t}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

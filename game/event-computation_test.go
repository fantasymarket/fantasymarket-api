package game_test

import (
	"fantasymarket/game"
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
				StartDate: timeutils.Time{Time: startDate},
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
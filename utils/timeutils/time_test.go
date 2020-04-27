package timeutils

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Date struct {
	Date Time `json:"date"`
}

func parseTime(s string) int64 {
	t, _ := time.Parse(time.RFC3339, s)
	return t.UnixNano()
}

var unmarshallJSON = map[string]int64{
	"{\"date\": \"2006-01-02 15:04:05\"}":  parseTime("2006-01-02T15:04:05Z"),
	"{\"date\": \"2006-01-02T15:04:05Z\"}": parseTime("2006-01-02T15:04:05Z"),
}

func TestUnmarshalJSONIntegration(t *testing.T) {
	for jsonData, expectedDate := range unmarshallJSON {
		var result Date
		err := json.Unmarshal([]byte(jsonData), &result)

		assert.NoError(t, err)
		assert.Equal(t, expectedDate, result.Date.UnixNano())
	}
}

func TestUnmarshalJSON(t *testing.T) {
	for jsonData, expectedDate := range unmarshallJSON {

		date := strings.TrimPrefix(jsonData, "{\"date\": ")
		date = strings.TrimSuffix(date, "}")

		timeInstance := Time{}
		err := timeInstance.UnmarshalJSON([]byte(date))
		assert.NoError(t, err)
		assert.Equal(t, expectedDate, timeInstance.UnixNano())
	}
}

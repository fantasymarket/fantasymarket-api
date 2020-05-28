package timeutils

import (
	"fantasymarket/utils/config"
	"fmt"
	"strings"
	"time"
)

// Time is our custom time type
type Time struct {
	time.Time
}

// MarshalJSON implements `json.Marshaler`
func (t Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))
	return []byte(stamp), nil
}

// UnmarshalJSON implements `json.Unmarshaler`
func (t *Time) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}

	if strings.ContainsAny(s, "TtZz") {
		t.Time, err = time.Parse(time.RFC3339, s)
	} else {
		t.Time, err = time.Parse("2006-01-02 15:04:05", s)
	}

	return
}

// UnmarshalYAML implements `yaml.Unmarshaler`
func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var b string
	err := unmarshal(&b)

	if err != nil {
		return err
	}

	return t.UnmarshalJSON([]byte(b))
}

func ConvertToTicks(timestamp string) (int, error) {
	// Vergiss nicht das Defaultconfig in config.go erstmal wieder private zu machen
	startDate := config.DefaultConfig.Game.StartDate.Time
	currentDate, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return -1, err
	}

	return int(currentDate.Sub(startDate).Hours()), nil
}

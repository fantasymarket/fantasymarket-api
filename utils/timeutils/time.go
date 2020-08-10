package timeutils

import (
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

// ParseTime parses a RFC3339 timestamp
func ParseTime(s string) (time.Time, error) {
	if strings.ContainsAny(s, "TtZz") {
		return time.Parse(time.RFC3339, s)
	}
	return time.Parse("2006-01-02 15:04:05", s)
}

// UnmarshalJSON implements `json.Unmarshaler`
func (t *Time) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = ParseTime(s)

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

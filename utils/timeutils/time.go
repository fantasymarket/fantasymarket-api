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

// UnmarshalJSON implements `json.Unmarshaler`
func (t *Time) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}

	if strings.ContainsAny(s, "TtZz") {
		fmt.Println("tfff", s)
		t.Time, err = time.Parse(time.RFC3339, s)
	} else {
		fmt.Println("naaa")
		t.Time, err = time.Parse("2006-01-02 15:04:05", s)
	}

	return
}

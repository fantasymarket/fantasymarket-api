package timeutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/senseyeio/duration"
)

// Duration is our custom duration type
type Duration struct {
	duration.Duration
}

// ShiftBack shifts time back by a duration
// Based on https://github.com/senseyeio/duration/blob/7c2a214ada4602c1d0638fb1abdbf8c6f25d0967/duration.go#L92
func ShiftBack(d Duration, t time.Time) time.Time {
	if d.Y != 0 || d.M != 0 || d.W != 0 || d.D != 0 {
		days := d.W*7 + d.D
		t = t.AddDate(-d.Y, -d.M, -days)
	}
	t = t.Add(timeDuration(d) * time.Duration(-1))
	return t
}

func timeDuration(d Duration) time.Duration {
	var dur time.Duration
	dur = dur + (time.Duration(d.TH) * time.Hour)
	dur = dur + (time.Duration(d.TM) * time.Minute)
	dur = dur + (time.Duration(d.TS) * time.Second)
	return dur
}

// UnmarshalJSON implements `json.Unmarshaler`
func (d *Duration) UnmarshalJSON(b []byte) (err error) {

	fmt.Println("fuck")

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Duration = duration.Duration{}
		return
	}

	d.Duration, err = duration.ParseISO8601(s)

	return
}

// UnmarshalYAML implements `yaml.Unmarshaler`
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var b = ""
	err := unmarshal(&b)
	if err != nil {
		return err
	}

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Duration = duration.Duration{}
		return nil
	}

	d.Duration, err = duration.ParseISO8601(s)

	return nil
}

// MarshalJSON implements `json.Marshaler`
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(timeDuration(d).String()), nil
}

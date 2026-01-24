package model

import (
	"strings"
	"time"
)

// Time wraps time.Time to handle timestamps without timezone suffix
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		return nil
	}

	// Try RFC3339 first
	parsed, err := time.Parse(time.RFC3339, s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Try without timezone (backend format)
	parsed, err = time.Parse("2006-01-02T15:04:05.999999", s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Try without microseconds
	parsed, err = time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t Time) Format(layout string) string {
	return t.Time.Format(layout)
}

func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

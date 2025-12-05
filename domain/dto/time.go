package dto

import (
	"encoding/json"
	"time"
)

const (
	// LocalTimezone is the timezone used for displaying timestamps
	LocalTimezone = "Asia/Phnom_Penh"
	// TimeFormat is the format used for displaying timestamps (e.g., "2025-12-05 6:50pm")
	TimeFormat = "2006-01-02 3:04pm"
)

// LocalTime represents a time that will be formatted in local timezone for JSON
type LocalTime struct {
	time.Time
}

// NewLocalTime creates a new LocalTime from a time.Time
func NewLocalTime(t time.Time) LocalTime {
	return LocalTime{Time: t}
}

// NewLocalTimePtr creates a pointer to LocalTime from a time.Time pointer
func NewLocalTimePtr(t *time.Time) *LocalTime {
	if t == nil {
		return nil
	}
	lt := LocalTime{Time: *t}
	return &lt
}

// MarshalJSON implements json.Marshaler interface
func (lt LocalTime) MarshalJSON() ([]byte, error) {
	// Load the local timezone
	loc, err := time.LoadLocation(LocalTimezone)
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}

	// Convert to local timezone
	localTime := lt.Time.In(loc)
	
	// Format as readable string
	formatted := localTime.Format(TimeFormat)
	
	return json.Marshal(formatted)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (lt *LocalTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Try to parse the formatted time
	loc, err := time.LoadLocation(LocalTimezone)
	if err != nil {
		loc = time.UTC
	}

	parsed, err := time.ParseInLocation(TimeFormat, str, loc)
	if err != nil {
		// Fallback to RFC3339 format
		parsed, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
	}

	lt.Time = parsed.UTC()
	return nil
}

// String returns the formatted time string
func (lt LocalTime) String() string {
	loc, err := time.LoadLocation(LocalTimezone)
	if err != nil {
		loc = time.UTC
	}
	return lt.Time.In(loc).Format(TimeFormat)
}

package entity

import "time"

// Setting represents user preferences and configuration
type Setting struct {
	ID                     uint
	UserID                 uint
	Language               string
	Theme                  string
	NotifyOnBooking        bool
	NotifyOnMusic          bool
	NotifyOnEvent          bool
	EnablePushNotifications bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// NewSetting creates a new setting entity with default values
func NewSetting(userID uint) *Setting {
	return &Setting{
		UserID:                 userID,
		Language:               "en",
		Theme:                  "light",
		NotifyOnBooking:        true,
		NotifyOnMusic:          false,
		NotifyOnEvent:          true,
		EnablePushNotifications: true,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

// IsValid validates setting entity
func (s *Setting) IsValid() bool {
	validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
	return s.UserID > 0 && validThemes[s.Theme]
}


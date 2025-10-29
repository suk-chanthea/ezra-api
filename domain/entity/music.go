package entity

import "time"

// Music represents a music track
type Music struct {
	ID          uint
	Title       string
	Artist      string
	Album       string
	Genre       string
	Duration    int    // in seconds
	BPM         int    // beats per minute
	Key         string // musical key (C, Am, G, etc.)
	Cover       string
	Lyrics      string
	Description string
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMusic creates a new music entity
func NewMusic(title, cover string, userID uint) *Music {
	return &Music{
		Title:     title,
		Cover:     cover,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsValid validates music entity
func (m *Music) IsValid() bool {
	return m.Title != "" && m.UserID > 0
}
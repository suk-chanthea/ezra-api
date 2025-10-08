package entity

import "time"

// Music represents a music track
type Music struct {
	ID        uint
	Title     string
	Cover     string
	Audio     string
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMusic creates a new music entity
func NewMusic(title, cover, audio string, userID uint) *Music {
	return &Music{
		Title:     title,
		Cover:     cover,
		Audio:     audio,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsValid validates music entity
func (m *Music) IsValid() bool {
	return m.Title != "" && m.Cover != ""
}
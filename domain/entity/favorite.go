package entity

import "time"

// Favorite represents a user's favorite music
type Favorite struct {
	ID        uint
	UserID    uint
	MusicID   uint
	CreatedAt time.Time
}

// NewFavorite creates a new favorite entity
func NewFavorite(userID, musicID uint) *Favorite {
	return &Favorite{
		UserID:    userID,
		MusicID:   musicID,
		CreatedAt: time.Now(),
	}
}

// IsValid validates favorite entity
func (f *Favorite) IsValid() bool {
	return f.UserID > 0 && f.MusicID > 0
}


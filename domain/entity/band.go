package entity

import "time"

// Band represents a music collection or group/organization
type Band struct {
	ID          uint
	Name        string
	Description string
	Cover       string
	IsPublic    bool
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewBand creates a new band entity
func NewBand(name, description, cover string, isPublic bool, userID uint) *Band {
	return &Band{
		Name:        name,
		Description: description,
		Cover:       cover,
		IsPublic:    isPublic,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsValid validates band entity
func (b *Band) IsValid() bool {
	return b.Name != "" && b.UserID > 0
}


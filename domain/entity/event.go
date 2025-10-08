package entity

import "time"

// Event represents the core event entity
type Event struct {
	ID        uint
	Title     string
	Content   string
	Cover     string
	Location  string
	StartTime time.Time
	EndTime   time.Time
	UserID    uint
	Musics    []*Music  // Add this field for related musics
	MusicIDs  []uint    // Add this for input
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewEvent creates a new event entity
func NewEvent(title, content, cover, location string, startTime, endTime time.Time, userID uint, musicIDs []uint) *Event {
	return &Event{
		Title:     title,
		Content:   content,
		Cover:     cover,
		Location:  location,
		StartTime: startTime,
		EndTime:   endTime,
		UserID:    userID,
		MusicIDs:  musicIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsValid validates event entity
func (e *Event) IsValid() bool {
	if e.Title == "" || e.Cover == "" || e.Location == "" {
		return false
	}
	if e.EndTime.Before(e.StartTime) {
		return false
	}
	return true
}
package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// EventRepository defines the interface for event data operations
type EventRepository interface {
	Save(event *entity.Event) error
	FindAll() ([]*entity.Event, error)
	FindAllPaginated(offset, limit int) ([]*entity.Event, int64, error)
	FindByID(id uint) (*entity.Event, error)
	FindByUserID(userID uint) ([]*entity.Event, error)
	Update(event *entity.Event) error
	Delete(id uint) error
	
	// Methods for managing event-music relationships
	AddMusicsToEvent(eventID uint, musicIDs []uint) error
	RemoveMusicsFromEvent(eventID uint, musicIDs []uint) error
	GetEventMusics(eventID uint) ([]*entity.Music, error)
}
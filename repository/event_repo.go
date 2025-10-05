package repository

import (
	"github.com/suk-chanthea/ezra/domain"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *domain.Event) error
	FindAll() ([]domain.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db}
}

func (r *eventRepository) Create(event *domain.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindAll() ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Find(&events).Error
	return events, err
}

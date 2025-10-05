package service

import (
	"github.com/suk-chanthea/ezra/domain"
	"github.com/suk-chanthea/ezra/repository"
)

type EventService interface {
	CreateEvent(event domain.Event) error
	GetEvents() ([]domain.Event, error)
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(r repository.EventRepository) EventService {
	return &eventService{r}
}

func (s *eventService) CreateEvent(event domain.Event) error {
	return s.repo.Create(&event)
}

func (s *eventService) GetEvents() ([]domain.Event, error) {
	return s.repo.FindAll()
}

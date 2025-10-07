package usecase

import (
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type EventUseCase interface {
	CreateEvent(req *dto.CreateEventRequest, userID uint) error
	GetAllEvents() ([]*dto.EventResponse, error)
	GetEventByID(id uint) (*dto.EventResponse, error)
	GetEventsByUserID(userID uint) ([]*dto.EventResponse, error)
}

type eventUseCase struct {
	eventRepo repository.EventRepository
}

func NewEventUseCase(repo repository.EventRepository) EventUseCase {
	return &eventUseCase{
		eventRepo: repo,
	}
}

func (uc *eventUseCase) CreateEvent(req *dto.CreateEventRequest, userID uint) error {
	// Create entity
	event := entity.NewEvent(
		req.Title,
		req.Content,
		req.Cover,
		req.Location,
		req.StartTime,
		req.EndTime,
		userID,
	)

	// Validate
	if !event.IsValid() {
		return errors.New("invalid event data")
	}

	// Save
	return uc.eventRepo.Save(event)
}

func (uc *eventUseCase) GetAllEvents() ([]*dto.EventResponse, error) {
	events, err := uc.eventRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(events), nil
}

func (uc *eventUseCase) GetEventByID(id uint) (*dto.EventResponse, error) {
	event, err := uc.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.entityToResponse(event), nil
}

func (uc *eventUseCase) GetEventsByUserID(userID uint) ([]*dto.EventResponse, error) {
	events, err := uc.eventRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(events), nil
}

func (uc *eventUseCase) entityToResponse(event *entity.Event) *dto.EventResponse {
	return &dto.EventResponse{
		ID:        event.ID,
		Title:     event.Title,
		Content:   event.Content,
		Cover:     event.Cover,
		Location:  event.Location,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		UserID:    event.UserID,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

func (uc *eventUseCase) entitiesToResponses(events []*entity.Event) []*dto.EventResponse {
	responses := make([]*dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = uc.entityToResponse(event)
	}
	return responses
}
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
	UpdateEvent(id uint, req *dto.UpdateEventRequest, userID uint) error
	DeleteEvent(id uint, userID uint) error
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
		req.MusicIDs,  // Add music IDs
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
	response := &dto.EventResponse{
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
	
	// Convert musics
	if len(event.Musics) > 0 {
		response.Musics = make([]*dto.MusicResponse, len(event.Musics))
		for i, music := range event.Musics {
			response.Musics[i] = &dto.MusicResponse{
				ID:        music.ID,
				Title:     music.Title,
				Cover:     music.Cover,
				Audio:     music.Audio,
				UserID:    music.UserID,
				CreatedAt: music.CreatedAt,
				UpdatedAt: music.UpdatedAt,
			}
		}
	}
	
	return response
}

func (uc *eventUseCase) entitiesToResponses(events []*entity.Event) []*dto.EventResponse {
	responses := make([]*dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = uc.entityToResponse(event)
	}
	return responses
}

func (uc *eventUseCase) UpdateEvent(id uint, req *dto.UpdateEventRequest, userID uint) error {
	// Check if event exists and belongs to user
	existingEvent, err := uc.eventRepo.FindByID(id)
	if err != nil {
		return errors.New("event not found")
	}

	if existingEvent.UserID != userID {
		return errors.New("unauthorized to update this event")
	}

	// Update event fields
	existingEvent.Title = req.Title
	existingEvent.Content = req.Content
	existingEvent.Cover = req.Cover
	existingEvent.Location = req.Location
	existingEvent.StartTime = req.StartTime
	existingEvent.EndTime = req.EndTime

	// Validate
	if !existingEvent.IsValid() {
		return errors.New("invalid event data")
	}

	// Update the event
	if err := uc.eventRepo.Update(existingEvent); err != nil {
		return err
	}

	// Update music associations if provided
	if req.MusicIDs != nil {
		// Remove old music associations
		oldMusics, err := uc.eventRepo.GetEventMusics(id)
		if err == nil && len(oldMusics) > 0 {
			oldMusicIDs := make([]uint, len(oldMusics))
			for i, music := range oldMusics {
				oldMusicIDs[i] = music.ID
			}
			uc.eventRepo.RemoveMusicsFromEvent(id, oldMusicIDs)
		}

		// Add new music associations
		if len(req.MusicIDs) > 0 {
			if err := uc.eventRepo.AddMusicsToEvent(id, req.MusicIDs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *eventUseCase) DeleteEvent(id uint, userID uint) error {
	// Check if event exists and belongs to user
	existingEvent, err := uc.eventRepo.FindByID(id)
	if err != nil {
		return errors.New("event not found")
	}

	if existingEvent.UserID != userID {
		return errors.New("unauthorized to delete this event")
	}

	return uc.eventRepo.Delete(id)
}
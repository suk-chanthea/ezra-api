package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type EventUseCase interface {
	CreateEvent(req *dto.CreateEventRequest, userID uint) error
	GetAllEvents() ([]*dto.EventResponse, error)
	GetAllEventsPaginated(page, pageSize int) ([]*dto.EventResponse, *dto.PaginationMetadata, error)
	GetEventByID(id uint) (*dto.EventResponse, error)
	GetEventsByUserID(userID uint) ([]*dto.EventResponse, error)
	UpdateEvent(id uint, req *dto.UpdateEventRequest, userID uint) error
	DeleteEvent(id uint, userID uint) error
}

type eventUseCase struct {
	eventRepo        repository.EventRepository
	musicRepo        repository.MusicRepository
	notificationRepo repository.NotificationRepository
}

func NewEventUseCase(eventRepo repository.EventRepository, musicRepo repository.MusicRepository, notificationRepo repository.NotificationRepository) EventUseCase {
	return &eventUseCase{
		eventRepo:        eventRepo,
		musicRepo:        musicRepo,
		notificationRepo: notificationRepo,
	}
}

func (uc *eventUseCase) CreateEvent(req *dto.CreateEventRequest, userID uint) error {
	// Validate music IDs if provided
	if len(req.MusicIDs) > 0 {
		musics, err := uc.musicRepo.FindByIDs(req.MusicIDs)
		
		if err != nil {
			return errors.New("failed to validate music IDs")
		}
		
		if len(musics) != len(req.MusicIDs) {
			return errors.New("one or more music IDs do not exist")
		}
	}

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

	// Save event
	if err := uc.eventRepo.Save(event); err != nil {
		return err
	}

	// Create broadcast notification for all users about the new event
	notification := entity.NewBroadcastNotification(
		fmt.Sprintf("New Event: %s", event.Title),
		fmt.Sprintf("%s at %s on %s", event.Title, event.Location, event.StartTime.Format("Jan 02, 2006 3:04 PM")),
		"event",
	)
	notification.SenderID = &userID
	notification.RelatedType = "event"
	notification.RelatedID = &event.ID
	
	// Send notification (don't fail event creation if notification fails)
	// Note: FCM push notifications are now handled by the notification system
	if err := uc.notificationRepo.Create(context.Background(), notification); err != nil {
		// Log the error but don't return it
		// In production, you might want to use a proper logger here
		// log.Printf("Failed to send event notification: %v", err)
	}

	return nil
}

func (uc *eventUseCase) GetAllEvents() ([]*dto.EventResponse, error) {
	events, err := uc.eventRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(events), nil
}

func (uc *eventUseCase) GetAllEventsPaginated(page, pageSize int) ([]*dto.EventResponse, *dto.PaginationMetadata, error) {
	offset := (page - 1) * pageSize
	events, total, err := uc.eventRepo.FindAllPaginated(offset, pageSize)
	if err != nil {
		return nil, nil, err
	}
	
	pagination := dto.NewPaginationMetadata(page, pageSize, total)
	return uc.entitiesToResponses(events), pagination, nil
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
				ID:          music.ID,
				Title:       music.Title,
				Artist:      music.Artist,
				Album:       music.Album,
				Genre:       music.Genre,
				Duration:    music.Duration,
				BPM:         music.BPM,
				Key:         music.Key,
				Cover:       music.Cover,
				Lyrics:      music.Lyrics,
				Description: music.Description,
				UserID:      music.UserID,
				CreatedAt:   music.CreatedAt,
				UpdatedAt:   music.UpdatedAt,
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
		// Validate music IDs if provided
		if len(req.MusicIDs) > 0 {
			musics, err := uc.musicRepo.FindByIDs(req.MusicIDs)
			
			if err != nil {
				return errors.New("failed to validate music IDs")
			}
			
			if len(musics) != len(req.MusicIDs) {
				return errors.New("one or more music IDs do not exist")
			}
		}
		
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
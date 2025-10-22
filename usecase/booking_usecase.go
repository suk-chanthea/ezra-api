package usecase

import (
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type BookingUseCase interface {
	CreateBooking(req *dto.CreateBookingRequest, userID uint) error
	GetAllBookings() ([]*dto.BookingResponse, error)
	GetBookingByID(id uint) (*dto.BookingResponse, error)
	GetBookingsByEventID(eventID uint) ([]*dto.BookingResponse, error)
	GetBookingsByUserID(userID uint) ([]*dto.BookingResponse, error)
	UpdateBooking(id uint, req *dto.UpdateBookingRequest, userID uint) error
	DeleteBooking(id uint, userID uint) error
}

type bookingUseCase struct {
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
}

func NewBookingUseCase(
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
) BookingUseCase {
	return &bookingUseCase{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func (uc *bookingUseCase) CreateBooking(req *dto.CreateBookingRequest, userID uint) error {
	// Check if event exists
	event, err := uc.eventRepo.FindByID(req.EventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Prevent users from booking their own events
	if event.UserID == userID {
		return errors.New("cannot book your own event")
	}

	// Check if user already booked this event
	existingBooking, _ := uc.bookingRepo.FindByEventAndUser(req.EventID, userID)
	if existingBooking != nil {
		return errors.New("you have already booked this event")
	}

	// Create booking entity
	booking := entity.NewBooking(req.EventID, userID, req.Notes)

	// Validate
	if !booking.IsValid() {
		return errors.New("invalid booking data")
	}

	// Save
	return uc.bookingRepo.Save(booking)
}

func (uc *bookingUseCase) GetAllBookings() ([]*dto.BookingResponse, error) {
	bookings, err := uc.bookingRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(bookings), nil
}

func (uc *bookingUseCase) GetBookingByID(id uint) (*dto.BookingResponse, error) {
	booking, err := uc.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.entityToResponse(booking), nil
}

func (uc *bookingUseCase) GetBookingsByEventID(eventID uint) ([]*dto.BookingResponse, error) {
	bookings, err := uc.bookingRepo.FindByEventID(eventID)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(bookings), nil
}

func (uc *bookingUseCase) GetBookingsByUserID(userID uint) ([]*dto.BookingResponse, error) {
	bookings, err := uc.bookingRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(bookings), nil
}

func (uc *bookingUseCase) UpdateBooking(id uint, req *dto.UpdateBookingRequest, userID uint) error {
	// Check if booking exists
	booking, err := uc.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("booking not found")
	}

	// Check authorization - only the user who made the booking can update it
	if booking.UserID != userID {
		return errors.New("unauthorized to update this booking")
	}

	// Update status
	booking.Status = entity.BookingStatus(req.Status)
	booking.Notes = req.Notes

	// Save
	return uc.bookingRepo.Update(booking)
}

func (uc *bookingUseCase) DeleteBooking(id uint, userID uint) error {
	// Check if booking exists
	booking, err := uc.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("booking not found")
	}

	// Check authorization - only the user who made the booking can delete it
	if booking.UserID != userID {
		return errors.New("unauthorized to delete this booking")
	}

	return uc.bookingRepo.Delete(id)
}

func (uc *bookingUseCase) entityToResponse(booking *entity.Booking) *dto.BookingResponse {
	response := &dto.BookingResponse{
		ID:        booking.ID,
		EventID:   booking.EventID,
		UserID:    booking.UserID,
		Status:    string(booking.Status),
		Notes:     booking.Notes,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}

	// Convert event if available
	if booking.Event != nil {
		response.Event = &dto.EventResponse{
			ID:        booking.Event.ID,
			Title:     booking.Event.Title,
			Content:   booking.Event.Content,
			Cover:     booking.Event.Cover,
			Location:  booking.Event.Location,
			StartTime: booking.Event.StartTime,
			EndTime:   booking.Event.EndTime,
			UserID:    booking.Event.UserID,
			CreatedAt: booking.Event.CreatedAt,
			UpdatedAt: booking.Event.UpdatedAt,
		}

		// Convert musics if available
		if len(booking.Event.Musics) > 0 {
			response.Event.Musics = make([]*dto.MusicResponse, len(booking.Event.Musics))
			for i, music := range booking.Event.Musics {
				response.Event.Musics[i] = &dto.MusicResponse{
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
	}

	// Convert user if available
	if booking.User != nil {
		response.User = &dto.UserResponse{
			ID:        booking.User.ID,
			Username:  booking.User.Username,
			Fullname:  booking.User.Fullname,
			Profile:   booking.User.Profile,
			Email:     booking.User.Email,
			Role:      booking.User.Role,
			CreatedAt: booking.User.CreatedAt,
			UpdatedAt: booking.User.UpdatedAt,
		}
	}

	return response
}

func (uc *bookingUseCase) entitiesToResponses(bookings []*entity.Booking) []*dto.BookingResponse {
	responses := make([]*dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = uc.entityToResponse(booking)
	}
	return responses
}


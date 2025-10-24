package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	Save(booking *entity.Booking) error
	FindAll() ([]*entity.Booking, error)
	FindAllPaginated(offset, limit int) ([]*entity.Booking, int64, error)
	FindByID(id uint) (*entity.Booking, error)
	FindByEventID(eventID uint) ([]*entity.Booking, error)
	FindByUserID(userID uint) ([]*entity.Booking, error)
	FindByEventAndUser(eventID, userID uint) (*entity.Booking, error)
	Update(booking *entity.Booking) error
	Delete(id uint) error
}

